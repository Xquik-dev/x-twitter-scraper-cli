// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestXTweetsCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "create",
			"--account", "@elonmusk",
			"--text", "Just launched our new feature!",
			"--attachment-url", "https://x.com/elonmusk/status/1234567890",
			"--community-id", "1500000000000000000",
			"--is-note-tweet=false",
			"--media-id", "1234567890123456789",
			"--reply-to-tweet-id", "1234567890",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"text: Just launched our new feature!\n" +
			"attachment_url: https://x.com/elonmusk/status/1234567890\n" +
			"community_id: '1500000000000000000'\n" +
			"is_note_tweet: false\n" +
			"media_ids:\n" +
			"  - '1234567890123456789'\n" +
			"reply_to_tweet_id: '1234567890'\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "create",
		)
	})
}

func TestXTweetsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "retrieve",
			"--id", "id",
		)
	})
}

func TestXTweetsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "list",
			"--ids", "ids",
		)
	})
}

func TestXTweetsDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "delete",
			"--id", "id",
			"--account", "@elonmusk",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("account: '@elonmusk'")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "delete",
			"--id", "id",
		)
	})
}

func TestXTweetsGetFavoriters(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "get-favoriters",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXTweetsGetQuotes(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "get-quotes",
			"--id", "id",
			"--cursor", "cursor",
			"--include-replies=true",
			"--since-time", "sinceTime",
			"--until-time", "untilTime",
		)
	})
}

func TestXTweetsGetReplies(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "get-replies",
			"--id", "id",
			"--cursor", "cursor",
			"--since-time", "sinceTime",
			"--until-time", "untilTime",
		)
	})
}

func TestXTweetsGetRetweeters(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "get-retweeters",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXTweetsGetThread(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "get-thread",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXTweetsSearch(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "search",
			"--q", "q",
			"--cursor", "cursor",
			"--limit", "200",
			"--query-type", "Latest",
			"--since-time", "sinceTime",
			"--until-time", "untilTime",
		)
	})
}
