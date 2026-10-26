// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package jsonview

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

const (
	hideCursor    = "\x1b[?25l"
	showCursor    = "\x1b[?25h"
	clearBelow    = "\r\x1b[J"
	defaultWidth  = 80
	defaultHeight = 24
)

type screen struct {
	input  io.Reader
	output io.Writer
	lines  int
}

type Option func(*screen)

func WithInput(input io.Reader) Option {
	return func(s *screen) { s.input = input }
}

func WithOutput(output io.Writer) Option {
	return func(s *screen) { s.output = output }
}

var finalKeys = map[byte]tea.KeyType{
	'A': tea.KeyUp,
	'B': tea.KeyDown,
	'C': tea.KeyRight,
	'D': tea.KeyLeft,
	'H': tea.KeyHome,
	'F': tea.KeyEnd,
}

var tildeKeys = map[string]tea.KeyType{
	"1": tea.KeyHome,
	"7": tea.KeyHome,
	"4": tea.KeyEnd,
	"8": tea.KeyEnd,
	"5": tea.KeyPgUp,
	"6": tea.KeyPgDown,
	"3": tea.KeyDelete,
}

var controlKeys = map[rune]tea.KeyType{
	0x03: tea.KeyCtrlC,
	0x04: tea.KeyCtrlD,
	0x15: tea.KeyCtrlU,
	'\r': tea.KeyEnter,
	'\n': tea.KeyEnter,
	0x7f: tea.KeyBackspace,
	0x08: tea.KeyBackspace,
}

func (s *screen) size() (int, int) {
	if file, ok := s.output.(*os.File); ok {
		if width, height, err := term.GetSize(file.Fd()); err == nil {
			return width, height
		}
	}
	return defaultWidth, defaultHeight
}

func (s *screen) draw(view string) error {
	var frame strings.Builder
	if s.lines > 1 {
		fmt.Fprintf(&frame, "\x1b[%dA", s.lines-1)
	}
	frame.WriteString(clearBelow)
	frame.WriteString(strings.ReplaceAll(view, "\n", "\r\n"))
	s.lines = strings.Count(view, "\n") + 1
	_, err := io.WriteString(s.output, frame.String())
	return err
}

func terminalName() string {
	if runtime.GOOS == "windows" {
		return "CONIN$"
	}
	return "/dev/tty"
}

func (s *screen) run(viewer *JSONViewer) (err error) {
	if s.input == nil {
		if !term.IsTerminal(os.Stdin.Fd()) {
			tty, openErr := os.OpenFile(terminalName(), os.O_RDWR, 0)
			if openErr != nil {
				return openErr
			}
			defer func() { err = errors.Join(err, tty.Close()) }()
			s.input = tty
		} else {
			s.input = os.Stdin
		}
	}
	if file, ok := s.input.(*os.File); ok && term.IsTerminal(file.Fd()) {
		state, rawErr := term.MakeRaw(file.Fd())
		if rawErr != nil {
			return rawErr
		}
		defer func() { err = errors.Join(err, term.Restore(file.Fd(), state)) }()
	}
	if _, writeErr := io.WriteString(s.output, hideCursor); writeErr != nil {
		return writeErr
	}
	defer func() {
		_, writeErr := io.WriteString(s.output, "\r\n"+showCursor)
		err = errors.Join(err, writeErr)
	}()

	width, height := s.size()
	viewer.resize(width-borderPadding, height)
	reader := bufio.NewReader(s.input)
	for {
		if drawErr := s.draw(viewer.View()); drawErr != nil {
			return drawErr
		}
		pressed, readErr := readKey(reader)
		if errors.Is(readErr, io.EOF) {
			return nil
		}
		if readErr != nil {
			return readErr
		}
		if newWidth, newHeight := s.size(); newWidth != width || newHeight != height {
			width, height = newWidth, newHeight
			viewer.resize(width-borderPadding, height)
		}
		quit, keyErr := viewer.handleKey(pressed)
		if keyErr != nil || quit {
			return keyErr
		}
	}
}

func readKey(reader *bufio.Reader) (tea.KeyMsg, error) {
	r, _, err := reader.ReadRune()
	if err != nil {
		return tea.KeyMsg{}, err
	}
	if keyType, ok := controlKeys[r]; ok {
		return tea.KeyMsg{Type: keyType}, nil
	}
	switch r {
	case ' ':
		return tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}, nil
	case 0x1b:
		return readEscape(reader)
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}, nil
}

func readEscape(reader *bufio.Reader) (tea.KeyMsg, error) {
	if reader.Buffered() == 0 {
		return tea.KeyMsg{Type: tea.KeyEsc}, nil
	}
	introducer, err := reader.ReadByte()
	if err != nil {
		return tea.KeyMsg{}, err
	}
	if introducer != '[' && introducer != 'O' {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(introducer)}, Alt: true}, nil
	}
	var parameter strings.Builder
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return tea.KeyMsg{}, err
		}
		if b < 0x40 || b > 0x7e {
			parameter.WriteByte(b)
			continue
		}
		keyType, ok := finalKeys[b]
		if b == '~' {
			keyType, ok = tildeKeys[parameter.String()]
		}
		if !ok {
			return tea.KeyMsg{Type: tea.KeyRunes}, nil
		}
		return tea.KeyMsg{Type: keyType}, nil
	}
}
