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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type failingReadCloser struct{}

func (failingReadCloser) Read([]byte) (int, error) { return 0, errors.New("read failed") }
func (failingReadCloser) Close() error             { return nil }

type errorIterator[T any] struct {
	item    T
	current bool
	err     error
}

func (iterator *errorIterator[T]) Next() bool {
	if iterator.current {
		return false
	}
	iterator.current = true
	return true
}

func (iterator *errorIterator[T]) Current() T { return iterator.item }
func (iterator *errorIterator[T]) Err() error { return iterator.err }

type rawJSONValue string

func (value rawJSONValue) RawJSON() string { return string(value) }

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

func TestPagerAndStreamingBranches(t *testing.T) {
	t.Setenv("PAGER", "definitely-not-an-installed-pager")
	require.Error(t, streamToPagerWithPipe("test", func(*os.File) error { return nil }))
	require.Error(t, streamOutputOSSpecific("test", func(*os.File) error { return nil }))

	t.Setenv("PAGER", "cat")
	require.NoError(t, streamToPagerWithPipe("test", func(writer *os.File) error {
		_, err := writer.WriteString("pager output\n")
		return err
	}))
	require.EqualError(t, streamToPagerWithPipe("test", func(*os.File) error {
		return errors.New("generation failed")
	}), "generation failed")
	require.NoError(t, streamToPagerWithPipe("test", func(*os.File) error {
		return errors.New("broken pipe")
	}))
	require.NoError(t, streamOutputOSSpecific("test", func(writer *os.File) error {
		_, err := writer.WriteString("socket output\n")
		return err
	}))
	_ = isPipedDataAvailableOSSpecific()
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

	channelIterator := &errorIterator[any]{item: make(chan int)}
	require.Error(t, ShowJSONIterator[any](channelIterator, -1, ShowJSONOpts{
		Format: "json",
		Stdout: os.Stdout,
	}))

	expected := errors.New("iterator failed")
	failingIterator := &errorIterator[rawJSONValue]{
		item: rawJSONValue(`{"a":1}`),
		err:  expected,
	}
	require.ErrorIs(t, ShowJSONIterator[rawJSONValue](failingIterator, -1, ShowJSONOpts{
		Format: "raw",
		Stdout: os.Stdout,
	}), expected)
	assert.Greater(t, countTerminalLines([]byte("first\nsecond\n"), 80), 0)
}
