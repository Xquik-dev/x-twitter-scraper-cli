// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var extractionsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get extraction results",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Previous nextCursor.",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "field-style",
			Usage:     "Preserve source keys or convert result field names.",
			Default:   "source",
			QueryPath: "fieldStyle",
		},
		&requestflag.Flag[bool]{
			Name:      "include-raw",
			Usage:     "Use outputMode=raw instead.",
			Default:   false,
			QueryPath: "includeRaw",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of results to return (1-1000, default 100)",
			Default:   100,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "output-mode",
			Usage:     "Select compact, full, or raw-compatible result fields.",
			Default:   "full",
			QueryPath: "outputMode",
		},
		&requestflag.Flag[string]{
			Name:      "output-preset",
			Usage:     "Keep enrichment nested or merge it into each result.",
			Default:   "nested",
			QueryPath: "outputPreset",
		},
	},
	Action:          handleExtractionsRetrieve,
	HideHelpCommand: true,
}

var extractionsList = cli.Command{
	Name:    "list",
	Usage:   "List extraction jobs",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Previous nextCursor.",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of items to return (1-100, default 50). For paid per-result endpoints, the returned count may be lower when remaining credits cannot cover the requested page. If zero paid results are affordable, the endpoint returns 402 insufficient_credits.\n",
			Default:   50,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "status",
			Usage:     "Filter by job status",
			QueryPath: "status",
		},
		&requestflag.Flag[string]{
			Name:      "tool-type",
			Usage:     "Filter by extraction tool type",
			QueryPath: "toolType",
		},
	},
	Action:          handleExtractionsList,
	HideHelpCommand: true,
}

var extractionsEstimateCost = requestflag.WithInnerFlags(cli.Command{
	Name:    "estimate-cost",
	Usage:   "Estimate extraction cost",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "tool-type",
			Usage:    "Identifier for the extraction tool used to run a job.",
			Required: true,
			BodyPath: "toolType",
		},
		&requestflag.Flag[string]{
			Name:     "advanced-query",
			Usage:    "Raw advanced search query appended as-is (tweet_search_extractor)",
			BodyPath: "advancedQuery",
		},
		&requestflag.Flag[string]{
			Name:     "any-words",
			Usage:    "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "anyWords",
		},
		&requestflag.Flag[string]{
			Name:     "bio-contains",
			Usage:    "Bio terms separated by commas or lines.",
			BodyPath: "bioContains",
		},
		&requestflag.Flag[bool]{
			Name:     "blue-verified-only",
			Usage:    "Return only Blue-verified Tweet authors.",
			Default:  false,
			BodyPath: "blueVerifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:     "bounding-box",
			Usage:    "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8 (tweet_search_extractor)",
			BodyPath: "boundingBox",
		},
		&requestflag.Flag[string]{
			Name:     "card-name",
			Usage:    "Match the Tweet card name.",
			BodyPath: "cardName",
		},
		&requestflag.Flag[string]{
			Name:     "cashtags",
			Usage:    "Cashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:     "collection-strategy",
			Usage:    "Reply collection strategy.",
			Default:  "auto",
			BodyPath: "collectionStrategy",
		},
		&requestflag.Flag[string]{
			Name:     "conversation-id",
			Usage:    "Conversation ID filter (tweet_search_extractor)",
			BodyPath: "conversationId",
		},
		&requestflag.Flag[bool]{
			Name:     "dedupe-across-targets",
			Usage:    "Merge duplicate results across collection targets.",
			Default:  true,
			BodyPath: "dedupeAcrossTargets",
		},
		&requestflag.Flag[string]{
			Name:     "dedupe-mode",
			Usage:    "Keep target duplicates, first rows, or merged overlap.",
			BodyPath: "dedupeMode",
		},
		&requestflag.Flag[string]{
			Name:     "exact-phrase",
			Usage:    "Exact phrase to match (tweet_search_extractor)",
			BodyPath: "exactPhrase",
		},
		&requestflag.Flag[bool]{
			Name:     "exclude-original-author",
			Usage:    "Exclude replies from the source author.",
			Default:  false,
			BodyPath: "excludeOriginalAuthor",
		},
		&requestflag.Flag[string]{
			Name:     "exclude-source",
			Usage:    "Exclude a source application.",
			BodyPath: "excludeSource",
		},
		&requestflag.Flag[string]{
			Name:     "exclude-words",
			Usage:    "Words or quoted phrases to exclude. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:     "from-user",
			Usage:    "Filter by author username (tweet_search_extractor)",
			BodyPath: "fromUser",
		},
		&requestflag.Flag[string]{
			Name:     "geocode",
			Usage:    "Match latitude, longitude, and radius.",
			BodyPath: "geocode",
		},
		&requestflag.Flag[string]{
			Name:     "hashtags",
			Usage:    "Hashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "hashtags",
		},
		&requestflag.Flag[bool]{
			Name:     "has-location",
			Usage:    "Require a profile location.",
			Default:  false,
			BodyPath: "hasLocation",
		},
		&requestflag.Flag[bool]{
			Name:     "has-media-only",
			Usage:    "Return only replies with media.",
			Default:  false,
			BodyPath: "hasMediaOnly",
		},
		&requestflag.Flag[bool]{
			Name:     "has-website",
			Usage:    "Require a profile website.",
			Default:  false,
			BodyPath: "hasWebsite",
		},
		&requestflag.Flag[bool]{
			Name:     "include-original-post",
			Usage:    "Include the source post in reply results.",
			Default:  false,
			BodyPath: "includeOriginalPost",
		},
		&requestflag.Flag[bool]{
			Name:     "include-search-terms",
			Usage:    "Add matching search terms to collection metadata.",
			Default:  false,
			BodyPath: "includeSearchTerms",
		},
		&requestflag.Flag[bool]{
			Name:     "include-target-metadata",
			Usage:    "Add source target metadata to each result.",
			Default:  true,
			BodyPath: "includeTargetMetadata",
		},
		&requestflag.Flag[string]{
			Name:     "in-reply-to-tweet-id",
			Usage:    "Only replies to this tweet ID (tweet_search_extractor)",
			BodyPath: "inReplyToTweetId",
		},
		&requestflag.Flag[string]{
			Name:     "language",
			Usage:    "Language code filter (tweet_search_extractor)",
			BodyPath: "language",
		},
		&requestflag.Flag[string]{
			Name:     "list-id",
			Usage:    "Search within a list ID (tweet_search_extractor)",
			BodyPath: "listId",
		},
		&requestflag.Flag[string]{
			Name:     "location-contains",
			Usage:    "Required profile location text.",
			BodyPath: "locationContains",
		},
		&requestflag.Flag[int64]{
			Name:     "max-depth",
			Usage:    "Maximum nested reply depth.",
			BodyPath: "maxDepth",
		},
		&requestflag.Flag[int64]{
			Name:     "max-followers",
			Usage:    "Maximum follower count for profile results.",
			BodyPath: "maxFollowers",
		},
		&requestflag.Flag[int64]{
			Name:     "max-following",
			Usage:    "Maximum following count for profile results.",
			BodyPath: "maxFollowing",
		},
		&requestflag.Flag[string]{
			Name:     "max-id",
			Usage:    "Return Tweets older than this Tweet ID.",
			BodyPath: "maxId",
		},
		&requestflag.Flag[int64]{
			Name:     "max-items-per-target",
			Usage:    "Maximum results collected for each target.",
			BodyPath: "maxItemsPerTarget",
		},
		&requestflag.Flag[int64]{
			Name:     "max-likes",
			Usage:    "Maximum Tweet like count.",
			BodyPath: "maxLikes",
		},
		&requestflag.Flag[int64]{
			Name:     "max-pages-per-target",
			Usage:    "Reply pages collected for each target.",
			BodyPath: "maxPagesPerTarget",
		},
		&requestflag.Flag[int64]{
			Name:     "max-posts",
			Usage:    "Maximum post count for profile results.",
			BodyPath: "maxPosts",
		},
		&requestflag.Flag[int64]{
			Name:     "max-quotes",
			Usage:    "Maximum Tweet quote count.",
			BodyPath: "maxQuotes",
		},
		&requestflag.Flag[int64]{
			Name:     "max-replies",
			Usage:    "Maximum Tweet reply count.",
			BodyPath: "maxReplies",
		},
		&requestflag.Flag[int64]{
			Name:     "max-retweets",
			Usage:    "Maximum Tweet repost count.",
			BodyPath: "maxRetweets",
		},
		&requestflag.Flag[string]{
			Name:     "media-type",
			Usage:    "Media type filter (tweet_search_extractor)",
			BodyPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:     "mentioning",
			Usage:    "Filter tweets mentioning a username (tweet_search_extractor)",
			BodyPath: "mentioning",
		},
		&requestflag.Flag[int64]{
			Name:     "min-account-age-days",
			Usage:    "Minimum profile age in days.",
			BodyPath: "minAccountAgeDays",
		},
		&requestflag.Flag[int64]{
			Name:     "min-bookmarks",
			Usage:    "Minimum Tweet bookmark count.",
			BodyPath: "minBookmarks",
		},
		&requestflag.Flag[int64]{
			Name:     "min-faves",
			Usage:    "Minimum likes threshold (tweet_search_extractor)",
			BodyPath: "minFaves",
		},
		&requestflag.Flag[int64]{
			Name:     "min-followers",
			Usage:    "Minimum follower count for profile results.",
			BodyPath: "minFollowers",
		},
		&requestflag.Flag[int64]{
			Name:     "min-following",
			Usage:    "Minimum following count for profile results.",
			BodyPath: "minFollowing",
		},
		&requestflag.Flag[int64]{
			Name:     "min-posts",
			Usage:    "Minimum post count for profile results.",
			BodyPath: "minPosts",
		},
		&requestflag.Flag[int64]{
			Name:     "min-quotes",
			Usage:    "Minimum quote count threshold (tweet_search_extractor)",
			BodyPath: "minQuotes",
		},
		&requestflag.Flag[int64]{
			Name:     "min-replies",
			Usage:    "Minimum replies threshold (tweet_search_extractor)",
			BodyPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:     "min-retweets",
			Usage:    "Minimum retweets threshold (tweet_search_extractor)",
			BodyPath: "minRetweets",
		},
		&requestflag.Flag[int64]{
			Name:     "min-views",
			Usage:    "Minimum Tweet view count.",
			BodyPath: "minViews",
		},
		&requestflag.Flag[bool]{
			Name:     "native-retweets",
			Usage:    "Only return native reposts.",
			Default:  false,
			BodyPath: "nativeRetweets",
		},
		&requestflag.Flag[string]{
			Name:     "near",
			Usage:    "Match a place name.",
			BodyPath: "near",
		},
		&requestflag.Flag[bool]{
			Name:     "news",
			Usage:    "Only return news results.",
			Default:  false,
			BodyPath: "news",
		},
		&requestflag.Flag[bool]{
			Name:     "overlap-mode",
			Usage:    "Shortcut for dedupeMode=merge.",
			Default:  false,
			BodyPath: "overlapMode",
		},
		&requestflag.Flag[string]{
			Name:     "place",
			Usage:    "Search within a place ID (tweet_search_extractor)",
			BodyPath: "place",
		},
		&requestflag.Flag[string]{
			Name:     "place-country",
			Usage:    "Search within a country code (tweet_search_extractor)",
			BodyPath: "placeCountry",
		},
		&requestflag.Flag[string]{
			Name:     "point-radius",
			Usage:    "Geo point radius, e.g. -73.99 40.73 25mi (tweet_search_extractor)",
			BodyPath: "pointRadius",
		},
		&requestflag.Flag[string]{
			Name:     "query-type",
			Usage:    "Search ranking applied to every query.",
			Default:  "Latest",
			BodyPath: "queryType",
		},
		&requestflag.Flag[string]{
			Name:     "quotes",
			Usage:    "Quote mode (tweet_search_extractor)",
			BodyPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:     "quotes-of-tweet-id",
			Usage:    "Only quotes of this tweet ID (tweet_search_extractor)",
			BodyPath: "quotesOfTweetId",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "relation-target",
			Usage:    "Profile relations processed within one job.",
			BodyPath: "relationTargets",
		},
		&requestflag.Flag[string]{
			Name:     "replies",
			Usage:    "Reply mode (tweet_search_extractor)",
			BodyPath: "replies",
		},
		&requestflag.Flag[int64]{
			Name:     "results-limit",
			Usage:    "Maximum number of results to extract. When set, the extraction stops after reaching this limit.",
			BodyPath: "resultsLimit",
		},
		&requestflag.Flag[string]{
			Name:     "retweets",
			Usage:    "Retweet mode (tweet_search_extractor)",
			BodyPath: "retweets",
		},
		&requestflag.Flag[string]{
			Name:     "retweets-of-tweet-id",
			Usage:    "Only retweets of this tweet ID (tweet_search_extractor)",
			BodyPath: "retweetsOfTweetId",
		},
		&requestflag.Flag[bool]{
			Name:     "safe",
			Usage:    "Enable the safe-search filter.",
			Default:  false,
			BodyPath: "safe",
		},
		&requestflag.Flag[string]{
			Name:     "scope",
			Usage:    "Reply depth scope.",
			Default:  "all",
			BodyPath: "scope",
		},
		&requestflag.Flag[[]string]{
			Name:     "search-query",
			Usage:    "Search queries processed as one collection job.",
			BodyPath: "searchQueries",
		},
		&requestflag.Flag[string]{
			Name:     "search-query",
			Usage:    "Required for tweet_search_extractor & community_search.",
			BodyPath: "searchQuery",
		},
		&requestflag.Flag[any]{
			Name:     "since-date",
			Usage:    "Start date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:     "since-id",
			Usage:    "Return Tweets newer than this Tweet ID.",
			BodyPath: "sinceId",
		},
		&requestflag.Flag[any]{
			Name:     "since-time",
			Usage:    "Reply start time as ISO 8601 or Unix seconds.",
			BodyPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:     "sort",
			Usage:    "Reply result order.",
			Default:  "relevance",
			BodyPath: "sort",
		},
		&requestflag.Flag[string]{
			Name:     "source",
			Usage:    "Match the source application.",
			BodyPath: "source",
		},
		&requestflag.Flag[string]{
			Name:     "start-cursor",
			Usage:    "Resume one reply target from this cursor.",
			BodyPath: "startCursor",
		},
		&requestflag.Flag[string]{
			Name:     "target-community-id",
			Usage:    "Required for community_post_extractor & community_search.",
			BodyPath: "targetCommunityId",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-community-id",
			Usage:    "Community IDs processed as one collection job.",
			BodyPath: "targetCommunityIds",
		},
		&requestflag.Flag[string]{
			Name:     "target-list-id",
			Usage:    "Required for list_follower_explorer, list_member_extractor & list_post_extractor.",
			BodyPath: "targetListId",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-list-id",
			Usage:    "List IDs processed as one collection job.",
			BodyPath: "targetListIds",
		},
		&requestflag.Flag[[]any]{
			Name:     "target",
			Usage:    "Mixed targets auto-routed within one job.",
			BodyPath: "targets",
		},
		&requestflag.Flag[string]{
			Name:     "target-space-id",
			Usage:    "Required for space_explorer.",
			BodyPath: "targetSpaceId",
		},
		&requestflag.Flag[string]{
			Name:     "target-tweet-id",
			BodyPath: "targetTweetId",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-tweet-id",
			Usage:    "Tweet IDs processed as one collection job.",
			BodyPath: "targetTweetIds",
		},
		&requestflag.Flag[string]{
			Name:     "target-username",
			BodyPath: "targetUsername",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-username",
			Usage:    "Usernames processed as one collection job.",
			BodyPath: "targetUsernames",
		},
		&requestflag.Flag[string]{
			Name:     "to-user",
			Usage:    "Filter replies sent to a username (tweet_search_extractor)",
			BodyPath: "toUser",
		},
		&requestflag.Flag[any]{
			Name:     "until-date",
			Usage:    "End date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "untilDate",
		},
		&requestflag.Flag[any]{
			Name:     "until-time",
			Usage:    "Reply end time as ISO 8601 or Unix seconds.",
			BodyPath: "untilTime",
		},
		&requestflag.Flag[string]{
			Name:     "url",
			Usage:    "URL substring or domain filter (tweet_search_extractor)",
			BodyPath: "url",
		},
		&requestflag.Flag[string]{
			Name:     "username-contains",
			Usage:    "Required username text.",
			BodyPath: "usernameContains",
		},
		&requestflag.Flag[bool]{
			Name:     "verified-only",
			Usage:    "Only verified authors (tweet_search_extractor)",
			BodyPath: "verifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:     "verified-type",
			Usage:    "Exact profile verification type.",
			BodyPath: "verifiedType",
		},
		&requestflag.Flag[string]{
			Name:     "within",
			Usage:    "Set the radius for the near filter.",
			BodyPath: "within",
		},
		&requestflag.Flag[string]{
			Name:     "within-time",
			Usage:    "Match Tweets inside a recent time window.",
			BodyPath: "withinTime",
		},
	},
	Action:          handleExtractionsEstimateCost,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"relation-target": {
		&requestflag.InnerFlag[string]{
			Name:       "relation-target.relation",
			Usage:      `Allowed values: "community_members", "followers", "following", "list_followers", "list_members", "verified_followers".`,
			InnerField: "relation",
		},
		&requestflag.InnerFlag[string]{
			Name:       "relation-target.value",
			InnerField: "value",
		},
	},
})

var extractionsExportResults = cli.Command{
	Name:    "export-results",
	Usage:   "Export extraction results",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "format",
			Usage:     "Export file format",
			Required:  true,
			QueryPath: "format",
		},
		&requestflag.Flag[bool]{
			Name:      "has-description",
			Usage:     "Require a non-empty description.",
			QueryPath: "hasDescription",
		},
		&requestflag.Flag[bool]{
			Name:      "has-location",
			Usage:     "Require a non-empty location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Flag[bool]{
			Name:      "has-media",
			Usage:     "Require media.",
			QueryPath: "hasMedia",
		},
		&requestflag.Flag[string]{
			Name:      "lang",
			Usage:     "Filter by language code.",
			QueryPath: "lang",
		},
		&requestflag.Flag[int64]{
			Name:      "max-followers",
			Usage:     "Maximum follower count.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Flag[int64]{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Flag[int64]{
			Name:      "max-posts",
			Usage:     "Maximum post count.",
			QueryPath: "maxPosts",
		},
		&requestflag.Flag[int64]{
			Name:      "min-followers",
			Usage:     "Minimum follower count.",
			QueryPath: "minFollowers",
		},
		&requestflag.Flag[int64]{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Flag[int64]{
			Name:      "min-likes",
			Usage:     "Minimum like count.",
			QueryPath: "minLikes",
		},
		&requestflag.Flag[int64]{
			Name:      "min-posts",
			Usage:     "Minimum post count.",
			QueryPath: "minPosts",
		},
		&requestflag.Flag[int64]{
			Name:      "min-replies",
			Usage:     "Minimum reply count.",
			QueryPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "min-retweets",
			Usage:     "Minimum repost count.",
			QueryPath: "minRetweets",
		},
		&requestflag.Flag[int64]{
			Name:      "min-views",
			Usage:     "Minimum view count.",
			QueryPath: "minViews",
		},
		&requestflag.Flag[string]{
			Name:      "search",
			Usage:     "Search exported result text.",
			QueryPath: "search",
		},
		&requestflag.Flag[any]{
			Name:      "since-date",
			Usage:     "Include results on or after this date.",
			QueryPath: "sinceDate",
		},
		&requestflag.Flag[any]{
			Name:      "until-date",
			Usage:     "Include results on or before this date.",
			QueryPath: "untilDate",
		},
		&requestflag.Flag[bool]{
			Name:      "verified",
			Usage:     "Filter by verified status.",
			QueryPath: "verified",
		},
		&requestflag.Flag[string]{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The file where the response contents will be stored. Use the value '-' to force output to stdout.",
		},
	},
	Action:          handleExtractionsExportResults,
	HideHelpCommand: true,
}

var extractionsRun = requestflag.WithInnerFlags(cli.Command{
	Name:    "run",
	Usage:   "Run extraction",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "tool-type",
			Usage:    "Identifier for the extraction tool used to run a job.",
			Required: true,
			BodyPath: "toolType",
		},
		&requestflag.Flag[bool]{
			Name:      "dry-run",
			Usage:     "Estimate cost without creating an extraction.",
			Default:   false,
			QueryPath: "dry_run",
		},
		&requestflag.Flag[string]{
			Name:     "advanced-query",
			Usage:    "Raw advanced search query appended as-is (tweet_search_extractor)",
			BodyPath: "advancedQuery",
		},
		&requestflag.Flag[string]{
			Name:     "any-words",
			Usage:    "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "anyWords",
		},
		&requestflag.Flag[string]{
			Name:     "bio-contains",
			Usage:    "Bio terms separated by commas or lines.",
			BodyPath: "bioContains",
		},
		&requestflag.Flag[bool]{
			Name:     "blue-verified-only",
			Usage:    "Return only Blue-verified Tweet authors.",
			Default:  false,
			BodyPath: "blueVerifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:     "bounding-box",
			Usage:    "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8 (tweet_search_extractor)",
			BodyPath: "boundingBox",
		},
		&requestflag.Flag[string]{
			Name:     "card-name",
			Usage:    "Match the Tweet card name.",
			BodyPath: "cardName",
		},
		&requestflag.Flag[string]{
			Name:     "cashtags",
			Usage:    "Cashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:     "collection-strategy",
			Usage:    "Reply collection strategy.",
			Default:  "auto",
			BodyPath: "collectionStrategy",
		},
		&requestflag.Flag[string]{
			Name:     "conversation-id",
			Usage:    "Conversation ID filter (tweet_search_extractor)",
			BodyPath: "conversationId",
		},
		&requestflag.Flag[bool]{
			Name:     "dedupe-across-targets",
			Usage:    "Merge duplicate results across collection targets.",
			Default:  true,
			BodyPath: "dedupeAcrossTargets",
		},
		&requestflag.Flag[string]{
			Name:     "dedupe-mode",
			Usage:    "Keep target duplicates, first rows, or merged overlap.",
			BodyPath: "dedupeMode",
		},
		&requestflag.Flag[string]{
			Name:     "exact-phrase",
			Usage:    "Exact phrase to match (tweet_search_extractor)",
			BodyPath: "exactPhrase",
		},
		&requestflag.Flag[bool]{
			Name:     "exclude-original-author",
			Usage:    "Exclude replies from the source author.",
			Default:  false,
			BodyPath: "excludeOriginalAuthor",
		},
		&requestflag.Flag[string]{
			Name:     "exclude-source",
			Usage:    "Exclude a source application.",
			BodyPath: "excludeSource",
		},
		&requestflag.Flag[string]{
			Name:     "exclude-words",
			Usage:    "Words or quoted phrases to exclude. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:     "from-user",
			Usage:    "Filter by author username (tweet_search_extractor)",
			BodyPath: "fromUser",
		},
		&requestflag.Flag[string]{
			Name:     "geocode",
			Usage:    "Match latitude, longitude, and radius.",
			BodyPath: "geocode",
		},
		&requestflag.Flag[string]{
			Name:     "hashtags",
			Usage:    "Hashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "hashtags",
		},
		&requestflag.Flag[bool]{
			Name:     "has-location",
			Usage:    "Require a profile location.",
			Default:  false,
			BodyPath: "hasLocation",
		},
		&requestflag.Flag[bool]{
			Name:     "has-media-only",
			Usage:    "Return only replies with media.",
			Default:  false,
			BodyPath: "hasMediaOnly",
		},
		&requestflag.Flag[bool]{
			Name:     "has-website",
			Usage:    "Require a profile website.",
			Default:  false,
			BodyPath: "hasWebsite",
		},
		&requestflag.Flag[bool]{
			Name:     "include-original-post",
			Usage:    "Include the source post in reply results.",
			Default:  false,
			BodyPath: "includeOriginalPost",
		},
		&requestflag.Flag[bool]{
			Name:     "include-search-terms",
			Usage:    "Add matching search terms to collection metadata.",
			Default:  false,
			BodyPath: "includeSearchTerms",
		},
		&requestflag.Flag[bool]{
			Name:     "include-target-metadata",
			Usage:    "Add source target metadata to each result.",
			Default:  true,
			BodyPath: "includeTargetMetadata",
		},
		&requestflag.Flag[string]{
			Name:     "in-reply-to-tweet-id",
			Usage:    "Only replies to this tweet ID (tweet_search_extractor)",
			BodyPath: "inReplyToTweetId",
		},
		&requestflag.Flag[string]{
			Name:     "language",
			Usage:    "Language code filter (tweet_search_extractor)",
			BodyPath: "language",
		},
		&requestflag.Flag[string]{
			Name:     "list-id",
			Usage:    "Search within a list ID (tweet_search_extractor)",
			BodyPath: "listId",
		},
		&requestflag.Flag[string]{
			Name:     "location-contains",
			Usage:    "Required profile location text.",
			BodyPath: "locationContains",
		},
		&requestflag.Flag[int64]{
			Name:     "max-depth",
			Usage:    "Maximum nested reply depth.",
			BodyPath: "maxDepth",
		},
		&requestflag.Flag[int64]{
			Name:     "max-followers",
			Usage:    "Maximum follower count for profile results.",
			BodyPath: "maxFollowers",
		},
		&requestflag.Flag[int64]{
			Name:     "max-following",
			Usage:    "Maximum following count for profile results.",
			BodyPath: "maxFollowing",
		},
		&requestflag.Flag[string]{
			Name:     "max-id",
			Usage:    "Return Tweets older than this Tweet ID.",
			BodyPath: "maxId",
		},
		&requestflag.Flag[int64]{
			Name:     "max-items-per-target",
			Usage:    "Maximum results collected for each target.",
			BodyPath: "maxItemsPerTarget",
		},
		&requestflag.Flag[int64]{
			Name:     "max-likes",
			Usage:    "Maximum Tweet like count.",
			BodyPath: "maxLikes",
		},
		&requestflag.Flag[int64]{
			Name:     "max-pages-per-target",
			Usage:    "Reply pages collected for each target.",
			BodyPath: "maxPagesPerTarget",
		},
		&requestflag.Flag[int64]{
			Name:     "max-posts",
			Usage:    "Maximum post count for profile results.",
			BodyPath: "maxPosts",
		},
		&requestflag.Flag[int64]{
			Name:     "max-quotes",
			Usage:    "Maximum Tweet quote count.",
			BodyPath: "maxQuotes",
		},
		&requestflag.Flag[int64]{
			Name:     "max-replies",
			Usage:    "Maximum Tweet reply count.",
			BodyPath: "maxReplies",
		},
		&requestflag.Flag[int64]{
			Name:     "max-retweets",
			Usage:    "Maximum Tweet repost count.",
			BodyPath: "maxRetweets",
		},
		&requestflag.Flag[string]{
			Name:     "media-type",
			Usage:    "Media type filter (tweet_search_extractor)",
			BodyPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:     "mentioning",
			Usage:    "Filter tweets mentioning a username (tweet_search_extractor)",
			BodyPath: "mentioning",
		},
		&requestflag.Flag[int64]{
			Name:     "min-account-age-days",
			Usage:    "Minimum profile age in days.",
			BodyPath: "minAccountAgeDays",
		},
		&requestflag.Flag[int64]{
			Name:     "min-bookmarks",
			Usage:    "Minimum Tweet bookmark count.",
			BodyPath: "minBookmarks",
		},
		&requestflag.Flag[int64]{
			Name:     "min-faves",
			Usage:    "Minimum likes threshold (tweet_search_extractor)",
			BodyPath: "minFaves",
		},
		&requestflag.Flag[int64]{
			Name:     "min-followers",
			Usage:    "Minimum follower count for profile results.",
			BodyPath: "minFollowers",
		},
		&requestflag.Flag[int64]{
			Name:     "min-following",
			Usage:    "Minimum following count for profile results.",
			BodyPath: "minFollowing",
		},
		&requestflag.Flag[int64]{
			Name:     "min-posts",
			Usage:    "Minimum post count for profile results.",
			BodyPath: "minPosts",
		},
		&requestflag.Flag[int64]{
			Name:     "min-quotes",
			Usage:    "Minimum quote count threshold (tweet_search_extractor)",
			BodyPath: "minQuotes",
		},
		&requestflag.Flag[int64]{
			Name:     "min-replies",
			Usage:    "Minimum replies threshold (tweet_search_extractor)",
			BodyPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:     "min-retweets",
			Usage:    "Minimum retweets threshold (tweet_search_extractor)",
			BodyPath: "minRetweets",
		},
		&requestflag.Flag[int64]{
			Name:     "min-views",
			Usage:    "Minimum Tweet view count.",
			BodyPath: "minViews",
		},
		&requestflag.Flag[bool]{
			Name:     "native-retweets",
			Usage:    "Only return native reposts.",
			Default:  false,
			BodyPath: "nativeRetweets",
		},
		&requestflag.Flag[string]{
			Name:     "near",
			Usage:    "Match a place name.",
			BodyPath: "near",
		},
		&requestflag.Flag[bool]{
			Name:     "news",
			Usage:    "Only return news results.",
			Default:  false,
			BodyPath: "news",
		},
		&requestflag.Flag[bool]{
			Name:     "overlap-mode",
			Usage:    "Shortcut for dedupeMode=merge.",
			Default:  false,
			BodyPath: "overlapMode",
		},
		&requestflag.Flag[string]{
			Name:     "place",
			Usage:    "Search within a place ID (tweet_search_extractor)",
			BodyPath: "place",
		},
		&requestflag.Flag[string]{
			Name:     "place-country",
			Usage:    "Search within a country code (tweet_search_extractor)",
			BodyPath: "placeCountry",
		},
		&requestflag.Flag[string]{
			Name:     "point-radius",
			Usage:    "Geo point radius, e.g. -73.99 40.73 25mi (tweet_search_extractor)",
			BodyPath: "pointRadius",
		},
		&requestflag.Flag[string]{
			Name:     "query-type",
			Usage:    "Search ranking applied to every query.",
			Default:  "Latest",
			BodyPath: "queryType",
		},
		&requestflag.Flag[string]{
			Name:     "quotes",
			Usage:    "Quote mode (tweet_search_extractor)",
			BodyPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:     "quotes-of-tweet-id",
			Usage:    "Only quotes of this tweet ID (tweet_search_extractor)",
			BodyPath: "quotesOfTweetId",
		},
		&requestflag.Flag[[]map[string]any]{
			Name:     "relation-target",
			Usage:    "Profile relations processed within one job.",
			BodyPath: "relationTargets",
		},
		&requestflag.Flag[string]{
			Name:     "replies",
			Usage:    "Reply mode (tweet_search_extractor)",
			BodyPath: "replies",
		},
		&requestflag.Flag[int64]{
			Name:     "results-limit",
			Usage:    "Maximum number of results to extract. When set, the extraction stops after reaching this limit.",
			BodyPath: "resultsLimit",
		},
		&requestflag.Flag[string]{
			Name:     "retweets",
			Usage:    "Retweet mode (tweet_search_extractor)",
			BodyPath: "retweets",
		},
		&requestflag.Flag[string]{
			Name:     "retweets-of-tweet-id",
			Usage:    "Only retweets of this tweet ID (tweet_search_extractor)",
			BodyPath: "retweetsOfTweetId",
		},
		&requestflag.Flag[bool]{
			Name:     "safe",
			Usage:    "Enable the safe-search filter.",
			Default:  false,
			BodyPath: "safe",
		},
		&requestflag.Flag[string]{
			Name:     "scope",
			Usage:    "Reply depth scope.",
			Default:  "all",
			BodyPath: "scope",
		},
		&requestflag.Flag[[]string]{
			Name:     "search-query",
			Usage:    "Search queries processed as one collection job.",
			BodyPath: "searchQueries",
		},
		&requestflag.Flag[string]{
			Name:     "search-query",
			Usage:    "Required for tweet_search_extractor & community_search.",
			BodyPath: "searchQuery",
		},
		&requestflag.Flag[any]{
			Name:     "since-date",
			Usage:    "Start date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:     "since-id",
			Usage:    "Return Tweets newer than this Tweet ID.",
			BodyPath: "sinceId",
		},
		&requestflag.Flag[any]{
			Name:     "since-time",
			Usage:    "Reply start time as ISO 8601 or Unix seconds.",
			BodyPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:     "sort",
			Usage:    "Reply result order.",
			Default:  "relevance",
			BodyPath: "sort",
		},
		&requestflag.Flag[string]{
			Name:     "source",
			Usage:    "Match the source application.",
			BodyPath: "source",
		},
		&requestflag.Flag[string]{
			Name:     "start-cursor",
			Usage:    "Resume one reply target from this cursor.",
			BodyPath: "startCursor",
		},
		&requestflag.Flag[string]{
			Name:     "target-community-id",
			Usage:    "Required for community_post_extractor & community_search.",
			BodyPath: "targetCommunityId",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-community-id",
			Usage:    "Community IDs processed as one collection job.",
			BodyPath: "targetCommunityIds",
		},
		&requestflag.Flag[string]{
			Name:     "target-list-id",
			Usage:    "Required for list_follower_explorer, list_member_extractor & list_post_extractor.",
			BodyPath: "targetListId",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-list-id",
			Usage:    "List IDs processed as one collection job.",
			BodyPath: "targetListIds",
		},
		&requestflag.Flag[[]any]{
			Name:     "target",
			Usage:    "Mixed targets auto-routed within one job.",
			BodyPath: "targets",
		},
		&requestflag.Flag[string]{
			Name:     "target-space-id",
			Usage:    "Required for space_explorer.",
			BodyPath: "targetSpaceId",
		},
		&requestflag.Flag[string]{
			Name:     "target-tweet-id",
			BodyPath: "targetTweetId",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-tweet-id",
			Usage:    "Tweet IDs processed as one collection job.",
			BodyPath: "targetTweetIds",
		},
		&requestflag.Flag[string]{
			Name:     "target-username",
			BodyPath: "targetUsername",
		},
		&requestflag.Flag[[]string]{
			Name:     "target-username",
			Usage:    "Usernames processed as one collection job.",
			BodyPath: "targetUsernames",
		},
		&requestflag.Flag[string]{
			Name:     "to-user",
			Usage:    "Filter replies sent to a username (tweet_search_extractor)",
			BodyPath: "toUser",
		},
		&requestflag.Flag[any]{
			Name:     "until-date",
			Usage:    "End date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "untilDate",
		},
		&requestflag.Flag[any]{
			Name:     "until-time",
			Usage:    "Reply end time as ISO 8601 or Unix seconds.",
			BodyPath: "untilTime",
		},
		&requestflag.Flag[string]{
			Name:     "url",
			Usage:    "URL substring or domain filter (tweet_search_extractor)",
			BodyPath: "url",
		},
		&requestflag.Flag[string]{
			Name:     "username-contains",
			Usage:    "Required username text.",
			BodyPath: "usernameContains",
		},
		&requestflag.Flag[bool]{
			Name:     "verified-only",
			Usage:    "Only verified authors (tweet_search_extractor)",
			BodyPath: "verifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:     "verified-type",
			Usage:    "Exact profile verification type.",
			BodyPath: "verifiedType",
		},
		&requestflag.Flag[string]{
			Name:     "within",
			Usage:    "Set the radius for the near filter.",
			BodyPath: "within",
		},
		&requestflag.Flag[string]{
			Name:     "within-time",
			Usage:    "Match Tweets inside a recent time window.",
			BodyPath: "withinTime",
		},
	},
	Action:          handleExtractionsRun,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"relation-target": {
		&requestflag.InnerFlag[string]{
			Name:       "relation-target.relation",
			Usage:      `Allowed values: "community_members", "followers", "following", "list_followers", "list_members", "verified_followers".`,
			InnerField: "relation",
		},
		&requestflag.InnerFlag[string]{
			Name:       "relation-target.value",
			InnerField: "value",
		},
	},
})

func handleExtractionsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ExtractionGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Extractions.Get(
		ctx,
		cmd.Value("id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "extractions retrieve",
		Transform:      transform,
	})
}

func handleExtractionsList(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ExtractionListParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Extractions.List(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "extractions list",
		Transform:      transform,
	})
}

func handleExtractionsEstimateCost(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ExtractionEstimateCostParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Extractions.EstimateCost(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "extractions estimate-cost",
		Transform:      transform,
	})
}

func handleExtractionsExportResults(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ExtractionExportResultsParams{}

	response, err := client.Extractions.ExportResults(
		ctx,
		cmd.Value("id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}
	message, err := writeBinaryResponse(response, os.Stdout, cmd.String("output"))
	if message != "" {
		fmt.Println(message)
	}
	return err
}

func handleExtractionsRun(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ExtractionRunParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Extractions.Run(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "extractions run",
		Transform:      transform,
	})
}
