// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
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
			"--field-style", "source",
			"--include-raw=true",
			"--limit", "1",
			"--output-mode", "compact",
			"--output-preset", "nested",
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
			"--bio-contains", "bioContains",
			"--blue-verified-only=true",
			"--bounding-box", "-74.1 40.6 -73.9 40.8",
			"--card-name", "cardName",
			"--cashtags", "$TSLA $NVDA",
			"--collection-strategy", "auto",
			"--conversation-id", "1234567890",
			"--dedupe-across-targets=true",
			"--dedupe-mode", "none",
			"--exact-phrase", "artificial intelligence",
			"--exclude-original-author=true",
			"--exclude-source", "excludeSource",
			"--exclude-words", "spam",
			"--from-user", "nasa",
			"--geocode", "geocode",
			"--hashtags", "#AI startups",
			"--has-location=true",
			"--has-media-only=true",
			"--has-website=true",
			"--include-original-post=true",
			"--include-search-terms=true",
			"--include-target-metadata=true",
			"--in-reply-to-tweet-id", "1234567890",
			"--language", "en",
			"--list-id", "1234567890",
			"--location-contains", "locationContains",
			"--max-depth", "1",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-id", "maxId",
			"--max-items-per-target", "1",
			"--max-likes", "0",
			"--max-pages-per-target", "1",
			"--max-posts", "0",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "example_user",
			"--min-account-age-days", "0",
			"--min-bookmarks", "0",
			"--min-faves", "10",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-posts", "0",
			"--min-quotes", "2",
			"--min-replies", "3",
			"--min-retweets", "5",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--overlap-mode=true",
			"--place", "96683cc9126741d1",
			"--place-country", "US",
			"--point-radius", "-73.99 40.73 25mi",
			"--query-type", "Latest",
			"--quotes", "include",
			"--quotes-of-tweet-id", "1234567890",
			"--relation-target", "{relation: community_members, value: x}",
			"--replies", "include",
			"--results-limit", "1000",
			"--retweets", "exclude",
			"--retweets-of-tweet-id", "1234567890",
			"--safe=true",
			"--scope", "all",
			"--search-query", "string",
			"--search-query", "AI trends 2025",
			"--since-date", "'2025-01-01'",
			"--since-id", "sinceId",
			"--since-time", "'2019-12-27T18:11:19.117Z'",
			"--sort", "relevance",
			"--source", "source",
			"--start-cursor", "x",
			"--target-community-id", "1500000000000000000",
			"--target-community-id", "string",
			"--target-list-id", "1234567890",
			"--target-list-id", "string",
			"--target", "string",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-tweet-id", "string",
			"--target-username", "elonmusk",
			"--target-username", "string",
			"--to-user", "openai",
			"--until-date", "'2025-12-31'",
			"--until-time", "'2019-12-27T18:11:19.117Z'",
			"--url", "example.com",
			"--username-contains", "usernameContains",
			"--verified-only=false",
			"--verified-type", "verifiedType",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(extractionsEstimateCost)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "estimate-cost",
			"--tool-type", "follower_explorer",
			"--advanced-query", "min_faves:100",
			"--any-words", "ChatGPT AI model",
			"--bio-contains", "bioContains",
			"--blue-verified-only=true",
			"--bounding-box", "-74.1 40.6 -73.9 40.8",
			"--card-name", "cardName",
			"--cashtags", "$TSLA $NVDA",
			"--collection-strategy", "auto",
			"--conversation-id", "1234567890",
			"--dedupe-across-targets=true",
			"--dedupe-mode", "none",
			"--exact-phrase", "artificial intelligence",
			"--exclude-original-author=true",
			"--exclude-source", "excludeSource",
			"--exclude-words", "spam",
			"--from-user", "nasa",
			"--geocode", "geocode",
			"--hashtags", "#AI startups",
			"--has-location=true",
			"--has-media-only=true",
			"--has-website=true",
			"--include-original-post=true",
			"--include-search-terms=true",
			"--include-target-metadata=true",
			"--in-reply-to-tweet-id", "1234567890",
			"--language", "en",
			"--list-id", "1234567890",
			"--location-contains", "locationContains",
			"--max-depth", "1",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-id", "maxId",
			"--max-items-per-target", "1",
			"--max-likes", "0",
			"--max-pages-per-target", "1",
			"--max-posts", "0",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "example_user",
			"--min-account-age-days", "0",
			"--min-bookmarks", "0",
			"--min-faves", "10",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-posts", "0",
			"--min-quotes", "2",
			"--min-replies", "3",
			"--min-retweets", "5",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--overlap-mode=true",
			"--place", "96683cc9126741d1",
			"--place-country", "US",
			"--point-radius", "-73.99 40.73 25mi",
			"--query-type", "Latest",
			"--quotes", "include",
			"--quotes-of-tweet-id", "1234567890",
			"--relation-target.relation", "community_members",
			"--relation-target.value", "x",
			"--replies", "include",
			"--results-limit", "1000",
			"--retweets", "exclude",
			"--retweets-of-tweet-id", "1234567890",
			"--safe=true",
			"--scope", "all",
			"--search-query", "string",
			"--search-query", "AI trends 2025",
			"--since-date", "'2025-01-01'",
			"--since-id", "sinceId",
			"--since-time", "'2019-12-27T18:11:19.117Z'",
			"--sort", "relevance",
			"--source", "source",
			"--start-cursor", "x",
			"--target-community-id", "1500000000000000000",
			"--target-community-id", "string",
			"--target-list-id", "1234567890",
			"--target-list-id", "string",
			"--target", "string",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-tweet-id", "string",
			"--target-username", "elonmusk",
			"--target-username", "string",
			"--to-user", "openai",
			"--until-date", "'2025-12-31'",
			"--until-time", "'2019-12-27T18:11:19.117Z'",
			"--url", "example.com",
			"--username-contains", "usernameContains",
			"--verified-only=false",
			"--verified-type", "verifiedType",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: follower_explorer\n" +
			"advancedQuery: min_faves:100\n" +
			"anyWords: ChatGPT AI model\n" +
			"bioContains: bioContains\n" +
			"blueVerifiedOnly: true\n" +
			"boundingBox: '-74.1 40.6 -73.9 40.8'\n" +
			"cardName: cardName\n" +
			"cashtags: $TSLA $NVDA\n" +
			"collectionStrategy: auto\n" +
			"conversationId: '1234567890'\n" +
			"dedupeAcrossTargets: true\n" +
			"dedupeMode: none\n" +
			"exactPhrase: artificial intelligence\n" +
			"excludeOriginalAuthor: true\n" +
			"excludeSource: excludeSource\n" +
			"excludeWords: spam\n" +
			"fromUser: nasa\n" +
			"geocode: geocode\n" +
			"hashtags: '#AI startups'\n" +
			"hasLocation: true\n" +
			"hasMediaOnly: true\n" +
			"hasWebsite: true\n" +
			"includeOriginalPost: true\n" +
			"includeSearchTerms: true\n" +
			"includeTargetMetadata: true\n" +
			"inReplyToTweetId: '1234567890'\n" +
			"language: en\n" +
			"listId: '1234567890'\n" +
			"locationContains: locationContains\n" +
			"maxDepth: 1\n" +
			"maxFollowers: 0\n" +
			"maxFollowing: 0\n" +
			"maxId: maxId\n" +
			"maxItemsPerTarget: 1\n" +
			"maxLikes: 0\n" +
			"maxPagesPerTarget: 1\n" +
			"maxPosts: 0\n" +
			"maxQuotes: 0\n" +
			"maxReplies: 0\n" +
			"maxRetweets: 0\n" +
			"mediaType: images\n" +
			"mentioning: example_user\n" +
			"minAccountAgeDays: 0\n" +
			"minBookmarks: 0\n" +
			"minFaves: 10\n" +
			"minFollowers: 0\n" +
			"minFollowing: 0\n" +
			"minPosts: 0\n" +
			"minQuotes: 2\n" +
			"minReplies: 3\n" +
			"minRetweets: 5\n" +
			"minViews: 0\n" +
			"nativeRetweets: true\n" +
			"near: near\n" +
			"news: true\n" +
			"overlapMode: true\n" +
			"place: 96683cc9126741d1\n" +
			"placeCountry: US\n" +
			"pointRadius: '-73.99 40.73 25mi'\n" +
			"queryType: Latest\n" +
			"quotes: include\n" +
			"quotesOfTweetId: '1234567890'\n" +
			"relationTargets:\n" +
			"  - relation: community_members\n" +
			"    value: x\n" +
			"replies: include\n" +
			"resultsLimit: 1000\n" +
			"retweets: exclude\n" +
			"retweetsOfTweetId: '1234567890'\n" +
			"safe: true\n" +
			"scope: all\n" +
			"searchQueries:\n" +
			"  - string\n" +
			"searchQuery: AI trends 2025\n" +
			"sinceDate: '2025-01-01'\n" +
			"sinceId: sinceId\n" +
			"sinceTime: '2019-12-27T18:11:19.117Z'\n" +
			"sort: relevance\n" +
			"source: source\n" +
			"startCursor: x\n" +
			"targetCommunityId: '1500000000000000000'\n" +
			"targetCommunityIds:\n" +
			"  - string\n" +
			"targetListId: '1234567890'\n" +
			"targetListIds:\n" +
			"  - string\n" +
			"targets:\n" +
			"  - string\n" +
			"targetSpaceId: 1vOGwMdBqpwGB\n" +
			"targetTweetId: '1234567890'\n" +
			"targetTweetIds:\n" +
			"  - string\n" +
			"targetUsername: elonmusk\n" +
			"targetUsernames:\n" +
			"  - string\n" +
			"toUser: openai\n" +
			"untilDate: '2025-12-31'\n" +
			"untilTime: '2019-12-27T18:11:19.117Z'\n" +
			"url: example.com\n" +
			"usernameContains: usernameContains\n" +
			"verifiedOnly: false\n" +
			"verifiedType: verifiedType\n" +
			"within: within\n" +
			"withinTime: withinTime\n")
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
			"--has-description=true",
			"--has-location=true",
			"--has-media=true",
			"--lang", "lang",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-posts", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-likes", "0",
			"--min-posts", "0",
			"--min-replies", "0",
			"--min-retweets", "0",
			"--min-views", "0",
			"--search", "search",
			"--since-date", "'2019-12-27'",
			"--until-date", "'2019-12-27'",
			"--verified=true",
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
			"--dry-run=true",
			"--advanced-query", "min_faves:100",
			"--any-words", "ChatGPT AI model",
			"--bio-contains", "bioContains",
			"--blue-verified-only=true",
			"--bounding-box", "-74.1 40.6 -73.9 40.8",
			"--card-name", "cardName",
			"--cashtags", "$TSLA $NVDA",
			"--collection-strategy", "auto",
			"--conversation-id", "1234567890",
			"--dedupe-across-targets=true",
			"--dedupe-mode", "none",
			"--exact-phrase", "artificial intelligence",
			"--exclude-original-author=true",
			"--exclude-source", "excludeSource",
			"--exclude-words", "spam",
			"--from-user", "nasa",
			"--geocode", "geocode",
			"--hashtags", "#AI startups",
			"--has-location=true",
			"--has-media-only=true",
			"--has-website=true",
			"--include-original-post=true",
			"--include-search-terms=true",
			"--include-target-metadata=true",
			"--in-reply-to-tweet-id", "1234567890",
			"--language", "en",
			"--list-id", "1234567890",
			"--location-contains", "locationContains",
			"--max-depth", "1",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-id", "maxId",
			"--max-items-per-target", "1",
			"--max-likes", "0",
			"--max-pages-per-target", "1",
			"--max-posts", "0",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "example_user",
			"--min-account-age-days", "0",
			"--min-bookmarks", "0",
			"--min-faves", "10",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-posts", "0",
			"--min-quotes", "2",
			"--min-replies", "3",
			"--min-retweets", "5",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--overlap-mode=true",
			"--place", "96683cc9126741d1",
			"--place-country", "US",
			"--point-radius", "-73.99 40.73 25mi",
			"--query-type", "Latest",
			"--quotes", "include",
			"--quotes-of-tweet-id", "1234567890",
			"--relation-target", "{relation: community_members, value: x}",
			"--replies", "include",
			"--results-limit", "1000",
			"--retweets", "exclude",
			"--retweets-of-tweet-id", "1234567890",
			"--safe=true",
			"--scope", "all",
			"--search-query", "string",
			"--search-query", "AI trends 2025",
			"--since-date", "'2025-01-01'",
			"--since-id", "sinceId",
			"--since-time", "'2019-12-27T18:11:19.117Z'",
			"--sort", "relevance",
			"--source", "source",
			"--start-cursor", "x",
			"--target-community-id", "1500000000000000000",
			"--target-community-id", "string",
			"--target-list-id", "1234567890",
			"--target-list-id", "string",
			"--target", "string",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-tweet-id", "string",
			"--target-username", "elonmusk",
			"--target-username", "string",
			"--to-user", "openai",
			"--until-date", "'2025-12-31'",
			"--until-time", "'2019-12-27T18:11:19.117Z'",
			"--url", "example.com",
			"--username-contains", "usernameContains",
			"--verified-only=false",
			"--verified-type", "verifiedType",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(extractionsRun)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "run",
			"--tool-type", "follower_explorer",
			"--dry-run=true",
			"--advanced-query", "min_faves:100",
			"--any-words", "ChatGPT AI model",
			"--bio-contains", "bioContains",
			"--blue-verified-only=true",
			"--bounding-box", "-74.1 40.6 -73.9 40.8",
			"--card-name", "cardName",
			"--cashtags", "$TSLA $NVDA",
			"--collection-strategy", "auto",
			"--conversation-id", "1234567890",
			"--dedupe-across-targets=true",
			"--dedupe-mode", "none",
			"--exact-phrase", "artificial intelligence",
			"--exclude-original-author=true",
			"--exclude-source", "excludeSource",
			"--exclude-words", "spam",
			"--from-user", "nasa",
			"--geocode", "geocode",
			"--hashtags", "#AI startups",
			"--has-location=true",
			"--has-media-only=true",
			"--has-website=true",
			"--include-original-post=true",
			"--include-search-terms=true",
			"--include-target-metadata=true",
			"--in-reply-to-tweet-id", "1234567890",
			"--language", "en",
			"--list-id", "1234567890",
			"--location-contains", "locationContains",
			"--max-depth", "1",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-id", "maxId",
			"--max-items-per-target", "1",
			"--max-likes", "0",
			"--max-pages-per-target", "1",
			"--max-posts", "0",
			"--max-quotes", "0",
			"--max-replies", "0",
			"--max-retweets", "0",
			"--media-type", "images",
			"--mentioning", "example_user",
			"--min-account-age-days", "0",
			"--min-bookmarks", "0",
			"--min-faves", "10",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-posts", "0",
			"--min-quotes", "2",
			"--min-replies", "3",
			"--min-retweets", "5",
			"--min-views", "0",
			"--native-retweets=true",
			"--near", "near",
			"--news=true",
			"--overlap-mode=true",
			"--place", "96683cc9126741d1",
			"--place-country", "US",
			"--point-radius", "-73.99 40.73 25mi",
			"--query-type", "Latest",
			"--quotes", "include",
			"--quotes-of-tweet-id", "1234567890",
			"--relation-target.relation", "community_members",
			"--relation-target.value", "x",
			"--replies", "include",
			"--results-limit", "1000",
			"--retweets", "exclude",
			"--retweets-of-tweet-id", "1234567890",
			"--safe=true",
			"--scope", "all",
			"--search-query", "string",
			"--search-query", "AI trends 2025",
			"--since-date", "'2025-01-01'",
			"--since-id", "sinceId",
			"--since-time", "'2019-12-27T18:11:19.117Z'",
			"--sort", "relevance",
			"--source", "source",
			"--start-cursor", "x",
			"--target-community-id", "1500000000000000000",
			"--target-community-id", "string",
			"--target-list-id", "1234567890",
			"--target-list-id", "string",
			"--target", "string",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-tweet-id", "string",
			"--target-username", "elonmusk",
			"--target-username", "string",
			"--to-user", "openai",
			"--until-date", "'2025-12-31'",
			"--until-time", "'2019-12-27T18:11:19.117Z'",
			"--url", "example.com",
			"--username-contains", "usernameContains",
			"--verified-only=false",
			"--verified-type", "verifiedType",
			"--within", "within",
			"--within-time", "withinTime",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: follower_explorer\n" +
			"advancedQuery: min_faves:100\n" +
			"anyWords: ChatGPT AI model\n" +
			"bioContains: bioContains\n" +
			"blueVerifiedOnly: true\n" +
			"boundingBox: '-74.1 40.6 -73.9 40.8'\n" +
			"cardName: cardName\n" +
			"cashtags: $TSLA $NVDA\n" +
			"collectionStrategy: auto\n" +
			"conversationId: '1234567890'\n" +
			"dedupeAcrossTargets: true\n" +
			"dedupeMode: none\n" +
			"exactPhrase: artificial intelligence\n" +
			"excludeOriginalAuthor: true\n" +
			"excludeSource: excludeSource\n" +
			"excludeWords: spam\n" +
			"fromUser: nasa\n" +
			"geocode: geocode\n" +
			"hashtags: '#AI startups'\n" +
			"hasLocation: true\n" +
			"hasMediaOnly: true\n" +
			"hasWebsite: true\n" +
			"includeOriginalPost: true\n" +
			"includeSearchTerms: true\n" +
			"includeTargetMetadata: true\n" +
			"inReplyToTweetId: '1234567890'\n" +
			"language: en\n" +
			"listId: '1234567890'\n" +
			"locationContains: locationContains\n" +
			"maxDepth: 1\n" +
			"maxFollowers: 0\n" +
			"maxFollowing: 0\n" +
			"maxId: maxId\n" +
			"maxItemsPerTarget: 1\n" +
			"maxLikes: 0\n" +
			"maxPagesPerTarget: 1\n" +
			"maxPosts: 0\n" +
			"maxQuotes: 0\n" +
			"maxReplies: 0\n" +
			"maxRetweets: 0\n" +
			"mediaType: images\n" +
			"mentioning: example_user\n" +
			"minAccountAgeDays: 0\n" +
			"minBookmarks: 0\n" +
			"minFaves: 10\n" +
			"minFollowers: 0\n" +
			"minFollowing: 0\n" +
			"minPosts: 0\n" +
			"minQuotes: 2\n" +
			"minReplies: 3\n" +
			"minRetweets: 5\n" +
			"minViews: 0\n" +
			"nativeRetweets: true\n" +
			"near: near\n" +
			"news: true\n" +
			"overlapMode: true\n" +
			"place: 96683cc9126741d1\n" +
			"placeCountry: US\n" +
			"pointRadius: '-73.99 40.73 25mi'\n" +
			"queryType: Latest\n" +
			"quotes: include\n" +
			"quotesOfTweetId: '1234567890'\n" +
			"relationTargets:\n" +
			"  - relation: community_members\n" +
			"    value: x\n" +
			"replies: include\n" +
			"resultsLimit: 1000\n" +
			"retweets: exclude\n" +
			"retweetsOfTweetId: '1234567890'\n" +
			"safe: true\n" +
			"scope: all\n" +
			"searchQueries:\n" +
			"  - string\n" +
			"searchQuery: AI trends 2025\n" +
			"sinceDate: '2025-01-01'\n" +
			"sinceId: sinceId\n" +
			"sinceTime: '2019-12-27T18:11:19.117Z'\n" +
			"sort: relevance\n" +
			"source: source\n" +
			"startCursor: x\n" +
			"targetCommunityId: '1500000000000000000'\n" +
			"targetCommunityIds:\n" +
			"  - string\n" +
			"targetListId: '1234567890'\n" +
			"targetListIds:\n" +
			"  - string\n" +
			"targets:\n" +
			"  - string\n" +
			"targetSpaceId: 1vOGwMdBqpwGB\n" +
			"targetTweetId: '1234567890'\n" +
			"targetTweetIds:\n" +
			"  - string\n" +
			"targetUsername: elonmusk\n" +
			"targetUsernames:\n" +
			"  - string\n" +
			"toUser: openai\n" +
			"untilDate: '2025-12-31'\n" +
			"untilTime: '2019-12-27T18:11:19.117Z'\n" +
			"url: example.com\n" +
			"usernameContains: usernameContains\n" +
			"verifiedOnly: false\n" +
			"verifiedType: verifiedType\n" +
			"within: within\n" +
			"withinTime: withinTime\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "run",
			"--dry-run=true",
		)
	})
}
