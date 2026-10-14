// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXWriteActionsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:write-actions", "retrieve",
			"--id", "id",
		)
	})
}
