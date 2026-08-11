// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXUsersRetrieve(t *testing.T) {
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

func TestXUsersRemoveFollower(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "remove-follower",
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
			"x:users", "remove-follower",
			"--id", "id",
			"--idempotency-key", "Idempotency-Key",
		)
	})
}

func TestXUsersRetrieveBatch(t *testing.T) {
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
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-followers",
			"--id", "id",
			"--after", "after",
			"--bio-contains", "bioContains",
			"--cursor", "cursor",
			"--has-location=true",
			"--has-website=true",
			"--limit", "1",
			"--location-contains", "locationContains",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-statuses", "0",
			"--min-account-age-days", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-statuses", "0",
			"--mode", "standard",
			"--page-size", "20",
			"--username-contains", "usernameContains",
			"--verified-only=true",
			"--verified-type", "verifiedType",
		)
	})
}

func TestXUsersRetrieveFollowersYouKnow(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-followers-you-know",
			"--id", "id",
			"--bio-contains", "bioContains",
			"--cursor", "cursor",
			"--has-location=true",
			"--has-website=true",
			"--location-contains", "locationContains",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-statuses", "0",
			"--min-account-age-days", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-statuses", "0",
			"--page-size", "20",
			"--username-contains", "usernameContains",
			"--verified-only=true",
			"--verified-type", "verifiedType",
		)
	})
}

func TestXUsersRetrieveFollowing(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-following",
			"--id", "id",
			"--after", "after",
			"--bio-contains", "bioContains",
			"--cursor", "cursor",
			"--has-location=true",
			"--has-website=true",
			"--limit", "1",
			"--location-contains", "locationContains",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-statuses", "0",
			"--min-account-age-days", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-statuses", "0",
			"--mode", "standard",
			"--page-size", "20",
			"--username-contains", "usernameContains",
			"--verified-only=true",
			"--verified-type", "verifiedType",
		)
	})
}

func TestXUsersRetrieveLikes(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-likes",
			"--id", "id",
			"--any-words", "anyWords",
			"--blue-verified-only=true",
			"--card-name", "cardName",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-source", "excludeSource",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--geocode", "geocode",
			"--hashtags", "hashtags",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--max-faves", "0",
			"--max-id", "maxId",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-bookmarks", "0",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--page-size", "1",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--safe=true",
			"--since-date", "'2019-12-27'",
			"--since-id", "sinceId",
			"--source", "source",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})
}

func TestXUsersRetrieveMedia(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-media",
			"--id", "id",
			"--any-words", "anyWords",
			"--blue-verified-only=true",
			"--card-name", "cardName",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-source", "excludeSource",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--geocode", "geocode",
			"--hashtags", "hashtags",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--max-faves", "0",
			"--max-id", "maxId",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-bookmarks", "0",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--page-size", "1",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--safe=true",
			"--since-date", "'2019-12-27'",
			"--since-id", "sinceId",
			"--source", "source",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})
}

func TestXUsersRetrieveMentions(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-mentions",
			"--id", "id",
			"--any-words", "anyWords",
			"--blue-verified-only=true",
			"--card-name", "cardName",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-source", "excludeSource",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--geocode", "geocode",
			"--hashtags", "hashtags",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--max-faves", "0",
			"--max-id", "maxId",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-bookmarks", "0",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--page-size", "1",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--safe=true",
			"--since-date", "'2019-12-27'",
			"--since-id", "sinceId",
			"--since-time", "sinceTime",
			"--source", "source",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--until-time", "untilTime",
			"--url", "url",
			"--verified-only=true",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})
}

func TestXUsersRetrieveReplies(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-replies",
			"--id", "id",
			"--any-words", "anyWords",
			"--blue-verified-only=true",
			"--card-name", "cardName",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-source", "excludeSource",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--geocode", "geocode",
			"--hashtags", "hashtags",
			"--include-parent-tweet=true",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--max-faves", "0",
			"--max-id", "maxId",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-bookmarks", "0",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--page-size", "1",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--safe=true",
			"--since-date", "'2019-12-27'",
			"--since-id", "sinceId",
			"--source", "source",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})
}

func TestXUsersRetrieveSearch(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-search",
			"--q", "q",
			"--bio-contains", "bioContains",
			"--cursor", "cursor",
			"--has-location=true",
			"--has-website=true",
			"--location-contains", "locationContains",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-statuses", "0",
			"--min-account-age-days", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-statuses", "0",
			"--username-contains", "usernameContains",
			"--verified-only=true",
			"--verified-type", "verifiedType",
		)
	})
}

func TestXUsersRetrieveTweets(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-tweets",
			"--id", "id",
			"--any-words", "anyWords",
			"--blue-verified-only=true",
			"--card-name", "cardName",
			"--cashtags", "cashtags",
			"--conversation-id", "conversationId",
			"--cursor", "cursor",
			"--exact-phrase", "exactPhrase",
			"--exclude-source", "excludeSource",
			"--exclude-words", "excludeWords",
			"--from-user", "fromUser",
			"--geocode", "geocode",
			"--hashtags", "hashtags",
			"--include-parent-tweet=true",
			"--include-replies=true",
			"--in-reply-to-tweet-id", "inReplyToTweetId",
			"--language", "language",
			"--max-faves", "0",
			"--max-id", "maxId",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "mentioning",
			"--min-bookmarks", "0",
			"--min-faves", "0",
			"--min-quotes", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--page-size", "1",
			"--quotes", "include",
			"--quotes-of-tweet-id", "quotesOfTweetId",
			"--replies", "include",
			"--retweets", "include",
			"--retweets-of-tweet-id", "retweetsOfTweetId",
			"--safe=true",
			"--since-date", "'2019-12-27'",
			"--since-id", "sinceId",
			"--source", "source",
			"--to-user", "toUser",
			"--until-date", "'2019-12-27'",
			"--url", "url",
			"--verified-only=true",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})
}

func TestXUsersRetrieveVerifiedFollowers(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:users", "retrieve-verified-followers",
			"--id", "id",
			"--after", "after",
			"--bio-contains", "bioContains",
			"--cursor", "cursor",
			"--has-location=true",
			"--has-website=true",
			"--limit", "1",
			"--location-contains", "locationContains",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-statuses", "0",
			"--min-account-age-days", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-statuses", "0",
			"--mode", "standard",
			"--page-size", "20",
			"--username-contains", "usernameContains",
			"--verified-only=true",
			"--verified-type", "verifiedType",
		)
	})
}
