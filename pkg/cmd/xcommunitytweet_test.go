// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestXCommunitiesTweetsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities:tweets", "list",
			"--max-items", "10",
			"--q", "q",
			"--cursor", "cursor",
			"--query-type", "queryType",
		)
	})
}

func TestXCommunitiesTweetsListByCommunity(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities:tweets", "list-by-community",
			"--max-items", "10",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}
