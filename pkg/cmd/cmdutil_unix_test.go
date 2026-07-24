// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

//go:build !windows

package cmd

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFailedPagerLookupDoesNotCloseReusedDescriptor(t *testing.T) {
	t.Setenv("PAGER", "definitely-not-an-installed-pager")

	_, _, err := openSocketPairPager("test")
	require.Error(t, err)

	file, err := os.CreateTemp(t.TempDir(), "descriptor")
	require.NoError(t, err)
	defer file.Close()

	runtime.GC()
	_, err = file.WriteString("still open")
	require.NoError(t, err)
}
