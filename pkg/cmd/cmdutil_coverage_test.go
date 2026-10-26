// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

type failingReadCloser struct{}

func (failingReadCloser) Read([]byte) (int, error) { return 0, errors.New("read failed") }
func (failingReadCloser) Close() error             { return nil }

func TestCommandValidatorsAndManpages(t *testing.T) {
	command := newCommand()
	command.Writer = io.Discard
	command.ErrWriter = io.Discard
	require.Error(t, command.Run(context.Background(), []string{"x-twitter-scraper", "--format", "invalid"}))

	command = newCommand()
	command.Writer = io.Discard
	command.ErrWriter = io.Discard
	require.Error(t, command.Run(context.Background(), []string{"x-twitter-scraper", "--format-error", "invalid"}))

	command = newCommand()
	command.Writer = io.Discard
	command.ErrWriter = io.Discard
	require.Error(t, command.Run(context.Background(), []string{"x-twitter-scraper", "--base-url", "invalid"}))

	output := t.TempDir()
	command = newCommand()
	command.Writer = io.Discard
	command.ErrWriter = io.Discard
	require.NoError(t, command.Run(context.Background(), []string{
		"x-twitter-scraper",
		"@manpages",
		"--output", output,
		"--text",
		"--gzip=false",
	}))
	_, err := os.Stat(filepath.Join(output, "man1", "x-twitter-scraper.1"))
	require.NoError(t, err)
}

func TestPipedDataIsDetected(t *testing.T) {
	t.Parallel()

	idle, idleWriter, err := os.Pipe()
	require.NoError(t, err)
	defer idle.Close()
	defer idleWriter.Close()
	_, ready := stdinProbe.check(idle)
	assert.False(t, ready)

	filled, filledWriter, err := os.Pipe()
	require.NoError(t, err)
	defer filled.Close()
	_, err = filledWriter.WriteString("piped data")
	require.NoError(t, err)
	require.NoError(t, filledWriter.Close())
	input, ready := stdinProbe.check(filled)
	require.True(t, ready)
	content, err := io.ReadAll(input)
	require.NoError(t, err)
	assert.Equal(t, "piped data", string(content))
}

func TestReadWithinKeepsEveryByte(t *testing.T) {
	t.Parallel()

	idle, idleWriter, err := os.Pipe()
	require.NoError(t, err)
	defer idle.Close()
	input, ready := readWithin(idle, 10*time.Millisecond)
	assert.False(t, ready)
	_, err = idleWriter.WriteString("late data")
	require.NoError(t, err)
	require.NoError(t, idleWriter.Close())
	content, err := io.ReadAll(input)
	require.NoError(t, err)
	assert.Equal(t, "late data", string(content))

	filled, filledWriter, err := os.Pipe()
	require.NoError(t, err)
	defer filled.Close()
	_, err = filledWriter.WriteString(strings.Repeat("x", 40*1024))
	require.NoError(t, err)
	require.NoError(t, filledWriter.Close())
	input, ready = readWithin(filled, time.Second)
	require.True(t, ready)
	content, err = io.ReadAll(input)
	require.NoError(t, err)
	assert.Len(t, content, 40*1024)

	failing := &controlledReader{err: errors.New("read failed")}
	input, ready = readWithin(failing, time.Second)
	require.True(t, ready)
	_, err = io.ReadAll(input)
	require.EqualError(t, err, "read failed")
}

type controlledReader struct {
	err error
}

func (r *controlledReader) Read([]byte) (int, error) { return 0, r.err }

func TestReadingProbeLeavesStdinAloneWithoutRequestInput(t *testing.T) {
	previousStdin, previousProbe := os.Stdin, stdinProbe
	t.Cleanup(func() { os.Stdin, stdinProbe = previousStdin, previousProbe })
	stdinProbe = pipeProbe{
		check:    func(stdin *os.File) (io.Reader, bool) { return readWithin(stdin, time.Second) },
		consumes: true,
	}
	pipe := func(t *testing.T, data string) *os.File {
		t.Helper()
		reader, writer, err := os.Pipe()
		require.NoError(t, err)
		t.Cleanup(func() { reader.Close() })
		_, err = writer.WriteString(data)
		require.NoError(t, err)
		require.NoError(t, writer.Close())
		return reader
	}

	untouched := pipe(t, "id: 42\n")
	os.Stdin = untouched
	_, err := flagOptions(&cli.Command{Name: "probe", Flags: []cli.Flag{&requestflag.String{Name: "note"}}}, EmptyBody)
	require.NoError(t, err)
	left, err := io.ReadAll(untouched)
	require.NoError(t, err)
	assert.Equal(t, "id: 42\n", string(left))

	pathInput := pipe(t, "id: 42\n")
	os.Stdin = pathInput
	command := &cli.Command{Name: "probe", Flags: []cli.Flag{&requestflag.String{Name: "id", PathParam: "id"}}}
	_, err = flagOptions(command, EmptyBody)
	require.NoError(t, err)
	assert.Equal(t, "42", command.String("id"))
	left, err = io.ReadAll(pathInput)
	require.NoError(t, err)
	assert.Empty(t, left)

	bodyInput := pipe(t, "text: hi\n")
	os.Stdin = bodyInput
	body := &cli.Command{Name: "probe", Flags: []cli.Flag{&requestflag.String{Name: "note"}}}
	_, err = flagOptions(body, ApplicationJSON)
	require.NoError(t, err)
	left, err = io.ReadAll(bodyInput)
	require.NoError(t, err)
	assert.Empty(t, left)
}

func TestBinaryResponseFailuresAndExtensions(t *testing.T) {
	response := &http.Response{Body: failingReadCloser{}}
	_, err := writeBinaryResponse(response, io.Discard, "-")
	require.EqualError(t, err, "read failed")

	closed, err := os.CreateTemp(t.TempDir(), "closed")
	require.NoError(t, err)
	require.NoError(t, closed.Close())
	response = &http.Response{Body: io.NopCloser(strings.NewReader("data"))}
	_, err = writeBinaryResponse(response, closed, "-")
	require.Error(t, err)

	response = &http.Response{Body: io.NopCloser(strings.NewReader("data"))}
	_, err = writeBinaryResponse(response, io.Discard, filepath.Join(t.TempDir(), "missing", "file"))
	require.Error(t, err)

	cases := [][]byte{
		[]byte("plain text"),
		{0x00, 0x01, 0x02, 0x03},
		[]byte("%PDF-1.7"),
		{0x1f, 0x8b, 0x08, 0x00},
		{'P', 'K', 0x03, 0x04},
		[]byte("GIF89a"),
		{0xff, 0xd8, 0xff, 0xe0},
		{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'},
	}
	for _, data := range cases {
		assert.NotEmpty(t, guessExtension(data))
	}
}

func TestColorAndJSONFormattingBranches(t *testing.T) {
	t.Setenv("FORCE_COLOR", "1")
	assert.True(t, shouldUseColors(io.Discard))
	colored, err := formatJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{Format: "json", Stdout: os.Stdout})
	require.NoError(t, err)
	assert.NotEmpty(t, colored)
	colored, err = formatJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{Format: "jsonl", Stdout: os.Stdout})
	require.NoError(t, err)
	assert.NotEmpty(t, colored)

	t.Setenv("FORCE_COLOR", "0")
	assert.False(t, shouldUseColors(os.Stdout))
	uncolored, err := formatJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{Format: "auto", Stdout: os.Stdout})
	require.NoError(t, err)
	assert.NotEmpty(t, uncolored)
	uncolored, err = formatJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{Format: "pretty", Stdout: os.Stdout, Title: "test"})
	require.NoError(t, err)
	assert.Contains(t, string(uncolored), "test")

	yaml, err := formatJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{Format: "yaml", Stdout: os.Stdout})
	require.NoError(t, err)
	assert.Nil(t, yaml)

	_, err = formatJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{Format: "invalid", Stdout: os.Stdout})
	require.Error(t, err)
	assert.False(t, shouldUseColors(bytes.NewBuffer(nil)))
}

func TestShowJSONAndIteratorErrors(t *testing.T) {
	output, err := os.CreateTemp(t.TempDir(), "output")
	require.NoError(t, err)
	require.NoError(t, ShowJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{
		Format: "auto",
		Stdout: output,
	}))
	require.NoError(t, output.Close())
	require.Error(t, ShowJSON(gjson.Parse(`{"a":1}`), ShowJSONOpts{
		Format: "raw",
		Stdout: output,
	}))

}
