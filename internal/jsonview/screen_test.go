// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package jsonview

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/help"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestReadKeyDecodesTerminalInput(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"q":         "q",
		"é":         "é",
		" ":         " ",
		"\r":        "enter",
		"\n":        "enter",
		"\x03":      "ctrl+c",
		"\x04":      "ctrl+d",
		"\x15":      "ctrl+u",
		"\x7f":      "backspace",
		"\x1b":      "esc",
		"\x1bx":     "alt+x",
		"\x1b[A":    "up",
		"\x1b[B":    "down",
		"\x1bOC":    "right",
		"\x1b[D":    "left",
		"\x1b[1;5A": "up",
		"\x1b[H":    "home",
		"\x1b[F":    "end",
		"\x1b[1~":   "home",
		"\x1b[4~":   "end",
		"\x1b[5~":   "pgup",
		"\x1b[6~":   "pgdown",
		"\x1b[3~":   "delete",
		"\x1b[9~":   "",
		"\x1b[Z":    "",
	}
	for input, want := range tests {
		pressed, err := readKey(bufio.NewReader(strings.NewReader(input)))
		require.NoError(t, err, "%q", input)
		assert.Equal(t, want, pressed.String(), "%q", input)
	}

	for _, input := range []string{"", "\x1b[", "\x1b[12"} {
		_, err := readKey(bufio.NewReader(strings.NewReader(input)))
		require.Error(t, err, "%q", input)
	}
}

func TestScreenRedrawsInPlace(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	terminal := &screen{output: output}
	require.NoError(t, terminal.draw("a\nb"))
	require.NoError(t, terminal.draw("c"))
	require.NoError(t, terminal.draw("d"))
	assert.Equal(t, "\r\x1b[Ja\r\nb\x1b[1A\r\x1b[Jc\r\x1b[Jd", output.String())
	width, height := terminal.size()
	assert.Equal(t, 80, width)
	assert.Equal(t, 24, height)
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestExploreJSONReportsOutputErrors(t *testing.T) {
	t.Parallel()

	err := ExploreJSON("Account", gjson.Parse(`{"name":"xquik"}`), WithInput(strings.NewReader("q")), WithOutput(failingWriter{}))
	require.ErrorContains(t, err, "write failed")
}

func TestExploreJSONEndsAtEndOfInput(t *testing.T) {
	t.Parallel()

	screen := &bytes.Buffer{}
	require.NoError(t, ExploreJSON("Account", gjson.Parse(`{"name":"xquik"}`), WithInput(strings.NewReader("")), WithOutput(screen)))
	assert.Contains(t, screen.String(), "Account")
	assert.True(t, strings.HasSuffix(screen.String(), "\r\n"+showCursor))
}

func TestExploreJSONMovesWithArrowKeys(t *testing.T) {
	stdout := os.Stdout
	printed, err := os.CreateTemp(t.TempDir(), "printed")
	require.NoError(t, err)
	os.Stdout = printed
	t.Cleanup(func() { os.Stdout = stdout })

	screen := &bytes.Buffer{}
	require.NoError(t, ExploreJSON("Accounts", gjson.Parse(`[{"name":"first"},{"name":"second"}]`), WithInput(strings.NewReader("\x1b[Bp")), WithOutput(screen)))

	content, err := os.ReadFile(printed.Name())
	require.NoError(t, err)
	assert.Equal(t, "\n{\"name\":\"second\"}\n", string(content))
}

func testViewer(t *testing.T) *JSONViewer {
	t.Helper()
	view, err := newView("", gjson.Parse(`{"name":"xquik"}`), false)
	require.NoError(t, err)
	return &JSONViewer{stack: []JSONView{view}, root: "Account", help: help.New()}
}

func TestScreenStopsOnSignalAndRestoresTheCursor(t *testing.T) {
	t.Parallel()

	input, inputWriter := io.Pipe()
	defer inputWriter.Close()
	signals := make(chan os.Signal, 1)
	signals <- syscall.SIGTERM
	output := &bytes.Buffer{}
	terminal := &screen{input: input, output: output, signals: signals, ticks: make(chan time.Time)}

	require.ErrorContains(t, terminal.run(testViewer(t)), "the explorer stopped on SIGTERM")
	assert.True(t, strings.HasSuffix(output.String(), "\r\n"+showCursor))
}

func TestScreenRedrawsWhenTheTerminalResizes(t *testing.T) {
	t.Parallel()

	input, inputWriter := io.Pipe()
	ticks := make(chan time.Time)
	var width atomic.Int64
	width.Store(80)
	output := &lockedBuffer{}
	terminal := &screen{
		input:   input,
		output:  output,
		signals: make(chan os.Signal),
		ticks:   ticks,
		measure: func() (int, int) { return int(width.Load()), 24 },
	}
	finished := make(chan error, 1)
	go func() { finished <- terminal.run(testViewer(t)) }()

	ticks <- time.Now()
	width.Store(120)
	ticks <- time.Now()
	ticks <- time.Now()
	require.NoError(t, inputWriter.Close())
	require.NoError(t, <-finished)
	assert.Equal(t, 2, strings.Count(output.String(), clearBelow))
}

type lockedBuffer struct {
	mutex  sync.Mutex
	buffer bytes.Buffer
}

func (b *lockedBuffer) Write(data []byte) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.Write(data)
}

func (b *lockedBuffer) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.String()
}
