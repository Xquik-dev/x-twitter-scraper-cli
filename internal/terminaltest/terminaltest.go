// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package terminaltest gives a package's tests the same standard streams and
// the same text styling whether a terminal, a pipe or go test started them.
//
// The commands ask whether standard output is a terminal before they color
// text, open a pager or read the window size. A test binary that a release
// runs from a maintainer's terminal would get other answers than go test
// gives, so its tests would color text or wait in a pager.
package terminaltest

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/termenv"
)

// Run runs a package's tests away from the terminal. Call it from TestMain
// and return, so the test runner sets the exit code and writes coverage data.
func Run(m *testing.M) {
	restore := detach(isTerminal)
	m.Run()
	must(restore(), "put the standard streams back")
}

func isTerminal(file *os.File) bool {
	return term.IsTerminal(file.Fd())
}

// must stops the tests when the standard streams cannot be moved, because a
// test that runs anyway would read the terminal.
func must(err error, step string) {
	if err != nil {
		panic("terminaltest could not " + step + ". " + err.Error())
	}
}

// detach points standard input at the null device and, when standard output
// is a terminal, routes it through a pipe that still ends at that terminal.
// It also fixes styled text to no colors on a dark background, so no test
// reads the terminal's colors. The function it returns puts both streams
// back and waits until the terminal has every byte. The styling stays fixed.
func detach(isTerminal func(*os.File) bool) func() error {
	lipgloss.SetColorProfile(termenv.Ascii)
	lipgloss.SetHasDarkBackground(true)

	stdin, stdout := os.Stdin, os.Stdout
	null, err := os.Open(os.DevNull)
	must(err, "open the null device")
	os.Stdin = null
	if !isTerminal(stdout) {
		return func() error {
			os.Stdin = stdin
			return null.Close()
		}
	}

	reader, writer, err := os.Pipe()
	must(err, "open a pipe")
	os.Stdout = writer
	copied := make(chan error, 1)
	go func() {
		_, copyErr := io.Copy(stdout, reader)
		copied <- copyErr
	}()
	return func() error {
		os.Stdin, os.Stdout = stdin, stdout
		closeErr := writer.Close()
		return errors.Join(closeErr, <-copied, reader.Close(), null.Close())
	}
}
