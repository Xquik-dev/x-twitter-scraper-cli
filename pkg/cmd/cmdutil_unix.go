// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package cmd

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func pipedDataOSSpecific(stdin *os.File) (io.Reader, bool) {
	fds := []unix.PollFd{{Fd: int32(stdin.Fd()), Events: unix.POLLIN}}
	n, _ := unix.Poll(fds, 10)
	return stdin, n > 0
}
