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

func pipedDataOSSpecific(stdin *os.File) (io.Reader, bool) {
	reader, writer := io.Pipe()
	ready := make(chan struct{})
	go func() {
		buffer := make([]byte, 32*1024)
		n, err := stdin.Read(buffer)
		close(ready)
		if n > 0 {
			if _, writeErr := writer.Write(buffer[:n]); writeErr != nil {
				return
			}
		}
		if err == nil {
			_, err = io.Copy(writer, stdin)
		}
		writer.CloseWithError(err)
	}()
	select {
	case <-ready:
		return reader, true
	case <-time.After(10 * time.Millisecond):
		return reader, false
	}
}
