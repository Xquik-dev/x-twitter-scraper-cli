// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestXListsRetrieveFollowers(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"x:lists", "retrieve-followers",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXListsRetrieveMembers(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"x:lists", "retrieve-members",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXListsRetrieveTweets(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"x:lists", "retrieve-tweets",
			"--id", "id",
			"--cursor", "cursor",
			"--include-replies=true",
			"--since-time", "sinceTime",
			"--until-time", "untilTime",
		)
	})
}
