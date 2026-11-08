// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package binaryparam opens a binary parameter from a file path or from standard
// input.
package binaryparam

import (
	"io"
	"os"
)

const stdinGlyph = "-"

func FileOrStdin(stdin io.ReadCloser, path string) (io.ReadCloser, bool, error) {
	switch path {
	case "", stdinGlyph, "/dev/fd/0", "/dev/stdin":
		return stdin, true, nil
	}

	readCloser, err := os.Open(path)
	if err != nil {
		return nil, false, err
	}

	return readCloser, false, err
}
