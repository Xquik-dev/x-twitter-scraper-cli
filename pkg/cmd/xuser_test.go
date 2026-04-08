// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestXUsersRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve",
			"--id", "id",
		)
	})
}

func TestXUsersRetrieveBatch(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-batch",
			"--ids", "ids",
		)
	})
}

func TestXUsersRetrieveFollowers(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-followers",
			"--id", "id",
			"--cursor", "cursor",
			"--page-size", "0",
		)
	})
}

func TestXUsersRetrieveFollowersYouKnow(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-followers-you-know",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXUsersRetrieveFollowing(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-following",
			"--id", "id",
			"--cursor", "cursor",
			"--page-size", "0",
		)
	})
}

func TestXUsersRetrieveLikes(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-likes",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXUsersRetrieveMedia(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-media",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXUsersRetrieveMentions(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-mentions",
			"--id", "id",
			"--cursor", "cursor",
			"--since-time", "sinceTime",
			"--until-time", "untilTime",
		)
	})
}

func TestXUsersRetrieveSearch(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-search",
			"--q", "q",
			"--cursor", "cursor",
		)
	})
}

func TestXUsersRetrieveTweets(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-tweets",
			"--id", "id",
			"--cursor", "cursor",
			"--include-parent-tweet=true",
			"--include-replies=true",
		)
	})
}

func TestXUsersRetrieveVerifiedFollowers(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-verified-followers",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}
