// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXCommunitiesTweetsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities:tweets", "list",
			"--community-id", "321669910225",
			"--q", "q",
			"--cursor", "cursor",
			"--page-size", "1",
			"--query-type", "Latest",
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
			"--id", "id",
			"--cursor", "cursor",
			"--page-size", "1",
		)
	})
}
