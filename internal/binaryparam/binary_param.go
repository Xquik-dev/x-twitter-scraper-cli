// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package binaryparam

import (
	"io"
	"os"
)

const stdinGlyph = "-"

// FileOrStdin opens a path or returns stdin for "-".
// The caller closes files and reserves stdin when the boolean is true.
func FileOrStdin(stdin io.ReadCloser, path string) (io.ReadCloser, bool, error) {
	// Support both portable "-" and Unix stdin paths.
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
