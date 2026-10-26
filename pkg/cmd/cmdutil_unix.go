// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package cmd

import (
	"os"

	"golang.org/x/sys/unix"
)

func isPipedDataAvailableOSSpecific() bool {
	fds := []unix.PollFd{{Fd: int32(os.Stdin.Fd()), Events: unix.POLLIN}}
	n, _ := unix.Poll(fds, 10)
	return n > 0
}
