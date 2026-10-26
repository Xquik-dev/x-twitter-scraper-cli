// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package terminaltest runs tests with the standard streams detached from the
// terminal that https://pkg.go.dev/github.com/charmbracelet/x/term detects.
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

func Run(m *testing.M) {
	restore := detach(isTerminal)
	m.Run()
	must(restore(), "put the standard streams back")
}

func isTerminal(file *os.File) bool {
	return term.IsTerminal(file.Fd())
}

func must(err error, step string) {
	if err != nil {
		panic("terminaltest could not " + step + ". " + err.Error())
	}
}

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
