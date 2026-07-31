// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
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
			"--idempotency-key", "Idempotency-Key",
			"--community-id", "1500000000000000000",
			"--is-note-tweet=false",
			"--media", "https://example.com/video.mp4",
			"--reply-to-tweet-id", "1234567890",
			"--text", "Just launched our new feature!",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"community_id: '1500000000000000000'\n" +
			"is_note_tweet: false\n" +
			"media:\n" +
			"  - https://example.com/video.mp4\n" +
			"reply_to_tweet_id: '1234567890'\n" +
			"text: Just launched our new feature!\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:tweets", "create",
			"--idempotency-key", "Idempotency-Key",
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
			"--idempotency-key", "Idempotency-Key",
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
			"--idempotency-key", "Idempotency-Key",
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
			"--page-size", "20",
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
			"--any-words", "anyWords",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--hashtags", "hashtags",
			"--include-replies=true",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--page-size", "1",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--since-date", "'2019-12-27'",
			"--since-time", "sinceTime",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--until-time", "untilTime",
			"--url", "url",
			"--verified-only=true",
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
			"--any-words", "anyWords",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--hashtags", "hashtags",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--limit", "1",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--mode", "complete",
			"--page-size", "1",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--since-date", "'2019-12-27'",
			"--since-time", "sinceTime",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--until-time", "untilTime",
			"--url", "url",
			"--verified-only=true",
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
			"--page-size", "20",
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
			"--page-size", "1",
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
			"--advanced-query", "advancedQuery",
			"--any-words", "anyWords",
			"--bounding-box", "boundingBox",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--hashtags", "hashtags",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--limit", "200",
			"--list-id", "listId",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--place", "place",
			"--place-country", "placeCountry",
			"--point-radius", "pointRadius",
			"--query-type", "Latest",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--since-date", "'2019-12-27'",
			"--since-time", "sinceTime",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--until-time", "untilTime",
			"--url", "url",
			"--verified-only=true",
		)
	})
}
