// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package jsonview

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/terminaltest"
)

// TestMain keeps these tests away from the terminal, so they give the same
// result from a maintainer's terminal as under go test.
func TestMain(m *testing.M) {
	terminaltest.Run(m)
}
