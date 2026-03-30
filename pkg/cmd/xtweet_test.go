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
			"--account", "account",
			"--text", "text",
			"--attachment-url", "attachment_url",
			"--community-id", "community_id",
			"--is-note-tweet=true",
			"--media-id", "string",
			"--reply-to-tweet-id", "reply_to_tweet_id",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: account\n" +
			"text: text\n" +
			"attachment_url: attachment_url\n" +
			"community_id: community_id\n" +
			"is_note_tweet: true\n" +
			"media_ids:\n" +
			"  - string\n" +
			"reply_to_tweet_id: reply_to_tweet_id\n")
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
			"--tweet-id", "tweetId",
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
			"--tweet-id", "tweetId",
			"--account", "account",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("account: account")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "delete",
			"--tweet-id", "tweetId",
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
