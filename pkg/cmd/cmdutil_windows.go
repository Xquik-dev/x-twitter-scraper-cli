// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package cmd

import (
	"io"
	"os"
	"time"
)

var stdinProbe = pipeProbe{
	check: func(stdin *os.File) (io.Reader, bool) {
		return readWithin(stdin, 10*time.Millisecond)
	},
	consumes: true,
}
