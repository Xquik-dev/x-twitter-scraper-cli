// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package cmd

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	procPeekNamedPipe = kernel32.NewProc("PeekNamedPipe")
)

func isPipedDataAvailableOSSpecific() bool {
	// Peek without consuming input because Windows lacks unix.Poll.
	var available uint32
	r, _, _ := procPeekNamedPipe.Call(
		os.Stdin.Fd(),
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&available)),
		0,
	)
	return r != 0 && available > 0
}

func streamOutputOSSpecific(label string, generateOutput func(w *os.File) error) error {
	// Windows uses pipes because Unix socket paging is unavailable.
	return streamToPagerWithPipe(label, generateOutput)
}
