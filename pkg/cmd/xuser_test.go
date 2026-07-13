// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXUsersRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"x:users", "retrieve",
			"--id", "id",
		)
	})
}

func TestXUsersRemoveFollower(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"x:users", "remove-follower",
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
			"--cookie-session", "string",
			"x:users", "remove-follower",
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
			"--cookie-session", "string",
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
			"--cookie-session", "string",
			"x:users", "retrieve-followers",
			"--id", "id",
			"--after", "after",
			"--cursor", "cursor",
			"--limit", "0",
			"--page-size", "20",
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
			"--cookie-session", "string",
			"x:users", "retrieve-followers-you-know",
			"--id", "id",
			"--cursor", "cursor",
			"--page-size", "20",
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
			"--cookie-session", "string",
			"x:users", "retrieve-following",
			"--id", "id",
			"--after", "after",
			"--cursor", "cursor",
			"--limit", "0",
			"--page-size", "20",
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
			"--cookie-session", "string",
			"x:users", "retrieve-likes",
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
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
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
			"--cookie-session", "string",
			"x:users", "retrieve-media",
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
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
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
			"--cookie-session", "string",
			"x:users", "retrieve-mentions",
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

func TestXUsersRetrieveReplies(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"x:users", "retrieve-replies",
			"--id", "id",
			"--any-words", "anyWords",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--hashtags", "hashtags",
			"--include-parent-tweet=true",
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
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
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
			"--cookie-session", "string",
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
			"--cookie-session", "string",
			"x:users", "retrieve-tweets",
			"--id", "id",
			"--any-words", "anyWords",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--hashtags", "hashtags",
			"--include-parent-tweet=true",
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
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
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
			"--cookie-session", "string",
			"x:users", "retrieve-verified-followers",
			"--id", "id",
			"--cursor", "cursor",
			"--page-size", "20",
		)
	})
}
