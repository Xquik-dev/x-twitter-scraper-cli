// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func isPipedDataAvailableOSSpecific() bool {
	fds := []unix.PollFd{{Fd: int32(os.Stdin.Fd()), Events: unix.POLLIN}}
	n, _ := unix.Poll(fds, 10)
	return n > 0
}

func streamOutputOSSpecific(label string, generateOutput func(w *os.File) error) error {
	pagerInput, pid, err := openSocketPairPager(label)
	if err != nil || pagerInput == nil {
		return streamToPagerWithPipe(label, generateOutput)
	}
	defer pagerInput.Close()

	if isTerminal(os.Stdout) && os.Getenv("FORCE_COLOR") == "" {
		os.Setenv("FORCE_COLOR", "1")
	}

	if err := generateOutput(pagerInput); err != nil &&
		!strings.Contains(err.Error(), "broken pipe") {
		return err
	}

	pagerInput.Close()

	var wstatus syscall.WaitStatus
	_, err = syscall.Wait4(pid, &wstatus, 0, nil)
	if wstatus.ExitStatus() != 0 {
		return fmt.Errorf("pager failed with exit status %d. Check PAGER", wstatus.ExitStatus())
	}
	return err
}

func openSocketPairPager(label string) (*os.File, int, error) {
	fds, err := unix.Socketpair(unix.AF_UNIX, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, 0, err
	}

	parentFd, childFd := fds[0], fds[1]
	defer unix.Close(childFd)

	if err := unix.SetsockoptInt(parentFd, unix.SOL_SOCKET, unix.SO_SNDBUF, 128); err != nil {
		unix.Close(parentFd)
		return nil, 0, err
	}
	if err := unix.SetsockoptInt(childFd, unix.SOL_SOCKET, unix.SO_RCVBUF, 128); err != nil {
		unix.Close(parentFd)
		return nil, 0, err
	}

	syscall.CloseOnExec(parentFd)

	parentConn := os.NewFile(uintptr(parentFd), "parent-socket")

	pagerProgram := os.Getenv("PAGER")
	if pagerProgram == "" {
		pagerProgram = "less"
	}

	pagerPath, err := exec.LookPath(pagerProgram)
	if err != nil {
		parentConn.Close()
		return nil, 0, err
	}

	env := os.Environ()
	env = append(env, "LESS=-r -P "+label)
	env = append(env, "MORE=-r -P "+label)

	procAttr := &syscall.ProcAttr{
		Dir: "",
		Env: env,
		Files: []uintptr{
			uintptr(childFd),
			uintptr(syscall.Stdout),
			uintptr(syscall.Stderr),
		},
	}

	pid, err := syscall.ForkExec(pagerPath, []string{pagerProgram}, procAttr)
	if err != nil {
		parentConn.Close()
		return nil, 0, err
	}

	return parentConn, pid, nil
}
