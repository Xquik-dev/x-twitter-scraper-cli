// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestExtractionsRetrieve(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "retrieve",
			"--id", "id",
			"--cursor", "cursor",
			"--limit", "1",
		)
	})
}

func TestExtractionsList(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "list",
			"--cursor", "cursor",
			"--limit", "1",
			"--status", "running",
			"--tool-type", "follower_explorer",
		)
	})
}

func TestExtractionsEstimateCost(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "estimate-cost",
			"--tool-type", "follower_explorer",
			"--advanced-query", "min_faves:100",
			"--any-words", "ChatGPT AI model",
			"--bounding-box", "-74.1 40.6 -73.9 40.8",
			"--cashtags", "$TSLA $NVDA",
			"--conversation-id", "1234567890",
			"--exact-phrase", "artificial intelligence",
			"--exclude-words", "spam",
			"--from-user", "nasa",
			"--hashtags", "#AI startups",
			"--in-reply-to-tweet-id", "1234567890",
			"--language", "en",
			"--list-id", "1234567890",
			"--media-type", "images",
			"--mentioning", "example_user",
			"--min-faves", "10",
			"--min-quotes", "2",
			"--min-replies", "3",
			"--min-retweets", "5",
			"--place", "96683cc9126741d1",
			"--place-country", "US",
			"--point-radius", "-73.99 40.73 25mi",
			"--quotes", "include",
			"--quotes-of-tweet-id", "1234567890",
			"--replies", "include",
			"--results-limit", "1000",
			"--retweets", "exclude",
			"--retweets-of-tweet-id", "1234567890",
			"--search-query", "AI trends 2025",
			"--since-date", "'2025-01-01'",
			"--target-community-id", "1500000000000000000",
			"--target-list-id", "1234567890",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-username", "elonmusk",
			"--to-user", "openai",
			"--until-date", "'2025-12-31'",
			"--url", "example.com",
			"--verified-only=false",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: follower_explorer\n" +
			"advancedQuery: min_faves:100\n" +
			"anyWords: ChatGPT AI model\n" +
			"boundingBox: '-74.1 40.6 -73.9 40.8'\n" +
			"cashtags: $TSLA $NVDA\n" +
			"conversationId: '1234567890'\n" +
			"exactPhrase: artificial intelligence\n" +
			"excludeWords: spam\n" +
			"fromUser: nasa\n" +
			"hashtags: '#AI startups'\n" +
			"inReplyToTweetId: '1234567890'\n" +
			"language: en\n" +
			"listId: '1234567890'\n" +
			"mediaType: images\n" +
			"mentioning: example_user\n" +
			"minFaves: 10\n" +
			"minQuotes: 2\n" +
			"minReplies: 3\n" +
			"minRetweets: 5\n" +
			"place: 96683cc9126741d1\n" +
			"placeCountry: US\n" +
			"pointRadius: '-73.99 40.73 25mi'\n" +
			"quotes: include\n" +
			"quotesOfTweetId: '1234567890'\n" +
			"replies: include\n" +
			"resultsLimit: 1000\n" +
			"retweets: exclude\n" +
			"retweetsOfTweetId: '1234567890'\n" +
			"searchQuery: AI trends 2025\n" +
			"sinceDate: '2025-01-01'\n" +
			"targetCommunityId: '1500000000000000000'\n" +
			"targetListId: '1234567890'\n" +
			"targetSpaceId: 1vOGwMdBqpwGB\n" +
			"targetTweetId: '1234567890'\n" +
			"targetUsername: elonmusk\n" +
			"toUser: openai\n" +
			"untilDate: '2025-12-31'\n" +
			"url: example.com\n" +
			"verifiedOnly: false\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "estimate-cost",
		)
	})
}

func TestExtractionsExportResults(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "export-results",
			"--id", "id",
			"--format", "csv",
			"--output", "/dev/null",
		)
	})
}

func TestExtractionsRun(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "run",
			"--tool-type", "follower_explorer",
			"--advanced-query", "min_faves:100",
			"--any-words", "ChatGPT AI model",
			"--bounding-box", "-74.1 40.6 -73.9 40.8",
			"--cashtags", "$TSLA $NVDA",
			"--conversation-id", "1234567890",
			"--exact-phrase", "artificial intelligence",
			"--exclude-words", "spam",
			"--from-user", "nasa",
			"--hashtags", "#AI startups",
			"--in-reply-to-tweet-id", "1234567890",
			"--language", "en",
			"--list-id", "1234567890",
			"--media-type", "images",
			"--mentioning", "example_user",
			"--min-faves", "10",
			"--min-quotes", "2",
			"--min-replies", "3",
			"--min-retweets", "5",
			"--place", "96683cc9126741d1",
			"--place-country", "US",
			"--point-radius", "-73.99 40.73 25mi",
			"--quotes", "include",
			"--quotes-of-tweet-id", "1234567890",
			"--replies", "include",
			"--results-limit", "1000",
			"--retweets", "exclude",
			"--retweets-of-tweet-id", "1234567890",
			"--search-query", "AI trends 2025",
			"--since-date", "'2025-01-01'",
			"--target-community-id", "1500000000000000000",
			"--target-list-id", "1234567890",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-username", "elonmusk",
			"--to-user", "openai",
			"--until-date", "'2025-12-31'",
			"--url", "example.com",
			"--verified-only=false",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: follower_explorer\n" +
			"advancedQuery: min_faves:100\n" +
			"anyWords: ChatGPT AI model\n" +
			"boundingBox: '-74.1 40.6 -73.9 40.8'\n" +
			"cashtags: $TSLA $NVDA\n" +
			"conversationId: '1234567890'\n" +
			"exactPhrase: artificial intelligence\n" +
			"excludeWords: spam\n" +
			"fromUser: nasa\n" +
			"hashtags: '#AI startups'\n" +
			"inReplyToTweetId: '1234567890'\n" +
			"language: en\n" +
			"listId: '1234567890'\n" +
			"mediaType: images\n" +
			"mentioning: example_user\n" +
			"minFaves: 10\n" +
			"minQuotes: 2\n" +
			"minReplies: 3\n" +
			"minRetweets: 5\n" +
			"place: 96683cc9126741d1\n" +
			"placeCountry: US\n" +
			"pointRadius: '-73.99 40.73 25mi'\n" +
			"quotes: include\n" +
			"quotesOfTweetId: '1234567890'\n" +
			"replies: include\n" +
			"resultsLimit: 1000\n" +
			"retweets: exclude\n" +
			"retweetsOfTweetId: '1234567890'\n" +
			"searchQuery: AI trends 2025\n" +
			"sinceDate: '2025-01-01'\n" +
			"targetCommunityId: '1500000000000000000'\n" +
			"targetListId: '1234567890'\n" +
			"targetSpaceId: 1vOGwMdBqpwGB\n" +
			"targetTweetId: '1234567890'\n" +
			"targetUsername: elonmusk\n" +
			"toUser: openai\n" +
			"untilDate: '2025-12-31'\n" +
			"url: example.com\n" +
			"verifiedOnly: false\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "run",
		)
	})
}
