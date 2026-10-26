// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package jsonview

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

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
