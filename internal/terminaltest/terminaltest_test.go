// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package terminaltest

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	Run(m)
}

func useStandardStreams(t *testing.T) (output *os.File) {
	t.Helper()
	stdin, stdout := os.Stdin, os.Stdout
	t.Cleanup(func() { os.Stdin, os.Stdout = stdin, stdout })

	input, err := os.CreateTemp(t.TempDir(), "input")
	require.NoError(t, err)
	_, err = input.WriteString("typed by a person")
	require.NoError(t, err)
	_, err = input.Seek(0, io.SeekStart)
	require.NoError(t, err)
	output, err = os.CreateTemp(t.TempDir(), "output")
	require.NoError(t, err)
	os.Stdin, os.Stdout = input, output
	return output
}

func readAll(t *testing.T, file *os.File) string {
	t.Helper()
	content, err := os.ReadFile(file.Name())
	require.NoError(t, err)
	return string(content)
}

func TestTestsOnATerminalSeeNoTerminalAndStillPrintToIt(t *testing.T) {
	terminal := useStandardStreams(t)
	onTerminal := func(file *os.File) bool { return file == terminal }

	restore := detach(onTerminal)

	assert.False(t, onTerminal(os.Stdout))
	typed, err := io.ReadAll(os.Stdin)
	require.NoError(t, err)
	assert.Empty(t, typed)
	_, err = os.Stdout.WriteString("printed by a test\n")
	require.NoError(t, err)

	require.NoError(t, restore())
	assert.Same(t, terminal, os.Stdout)
	assert.Equal(t, "printed by a test\n", readAll(t, terminal))
}

func TestTestsOnAPipeKeepTheirOutputAndReadNoInput(t *testing.T) {
	pipe := useStandardStreams(t)

	restore := detach(func(*os.File) bool { return false })

	assert.Same(t, pipe, os.Stdout)
	typed, err := io.ReadAll(os.Stdin)
	require.NoError(t, err)
	assert.Empty(t, typed)

	require.NoError(t, restore())
	assert.Same(t, pipe, os.Stdout)
}

func TestStyledTextHasNoColorCodes(t *testing.T) {
	restore := detach(func(*os.File) bool { return false })
	t.Cleanup(func() { require.NoError(t, restore()) })

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("113")).Bold(true)
	adaptive := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})

	assert.Equal(t, "value", style.Render("value"))
	assert.Empty(t, style.Render(""))
	assert.Equal(t, "value", adaptive.Render("value"))
}

func TestAStreamThatCannotMoveStopsTheTests(t *testing.T) {
	assert.PanicsWithValue(
		t,
		"terminaltest could not open a pipe. too many open files",
		func() { must(errors.New("too many open files"), "open a pipe") },
	)
	assert.NotPanics(t, func() { must(nil, "open a pipe") })
}
