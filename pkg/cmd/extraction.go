// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"os"

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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Previous nextCursor.",
			QueryPath: "cursor",
		},
		&requestflag.String{
			Name:      "field-style",
			Usage:     "Preserve source keys or convert result field names.",
			Default:   "source",
			QueryPath: "fieldStyle",
		},
		&requestflag.Bool{
			Name:      "include-raw",
			Usage:     "Use outputMode=raw instead.",
			Default:   false,
			QueryPath: "includeRaw",
		},
		&requestflag.Int{
			Name:      "limit",
			Usage:     "Maximum number of results to return (1-1000, default 100)",
			Default:   100,
			QueryPath: "limit",
		},
		&requestflag.String{
			Name:      "output-mode",
			Usage:     "Select compact, full, or raw-compatible result fields.",
			Default:   "full",
			QueryPath: "outputMode",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Previous nextCursor.",
			QueryPath: "cursor",
		},
		&requestflag.Int{
			Name:      "limit",
			Usage:     "Maximum number of items to return (1-100, default 50). For paid per-result endpoints, the returned count may be lower when remaining credits cannot cover the requested page. If zero paid results are affordable, the endpoint returns 402 insufficient_credits.\n",
			Default:   50,
			QueryPath: "limit",
		},
		&requestflag.String{
			Name:      "status",
			Usage:     "Filter by job status",
			QueryPath: "status",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:     "tool-type",
			Usage:    "Identifier for the extraction tool used to run a job.",
			Required: true,
			BodyPath: "toolType",
		},
		&requestflag.String{
			Name:     "advanced-query",
			Usage:    "Raw advanced search query appended as-is (tweet_search_extractor)",
			BodyPath: "advancedQuery",
		},
		&requestflag.String{
			Name:     "any-words",
			Usage:    "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "anyWords",
		},
		&requestflag.String{
			Name:     "bio-contains",
			Usage:    "Bio terms separated by commas or lines.",
			BodyPath: "bioContains",
		},
		&requestflag.Bool{
			Name:     "blue-verified-only",
			Usage:    "Return only Blue-verified Tweet authors.",
			Default:  false,
			BodyPath: "blueVerifiedOnly",
		},
		&requestflag.String{
			Name:     "bounding-box",
			Usage:    "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8 (tweet_search_extractor)",
			BodyPath: "boundingBox",
		},
		&requestflag.String{
			Name:     "card-name",
			Usage:    "Match the Tweet card name.",
			BodyPath: "cardName",
		},
		&requestflag.String{
			Name:     "cashtags",
			Usage:    "Cashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "cashtags",
		},
		&requestflag.String{
			Name:     "collection-strategy",
			Usage:    "Reply collection strategy.",
			Default:  "auto",
			BodyPath: "collectionStrategy",
		},
		&requestflag.String{
			Name:     "conversation-id",
			Usage:    "Conversation ID filter (tweet_search_extractor)",
			BodyPath: "conversationId",
		},
		&requestflag.Bool{
			Name:     "dedupe-across-targets",
			Usage:    "Merge duplicate results across collection targets.",
			Default:  true,
			BodyPath: "dedupeAcrossTargets",
		},
		&requestflag.String{
			Name:     "dedupe-mode",
			Usage:    "Keep target duplicates, first rows, or merged overlap.",
			BodyPath: "dedupeMode",
		},
		&requestflag.String{
			Name:     "exact-phrase",
			Usage:    "Exact phrase to match (tweet_search_extractor)",
			BodyPath: "exactPhrase",
		},
		&requestflag.Bool{
			Name:     "exclude-original-author",
			Usage:    "Exclude replies from the source author.",
			Default:  false,
			BodyPath: "excludeOriginalAuthor",
		},
		&requestflag.String{
			Name:     "exclude-source",
			Usage:    "Exclude a source application.",
			BodyPath: "excludeSource",
		},
		&requestflag.String{
			Name:     "exclude-words",
			Usage:    "Words or quoted phrases to exclude. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "excludeWords",
		},
		&requestflag.String{
			Name:     "from-user",
			Usage:    "Filter by author username (tweet_search_extractor)",
			BodyPath: "fromUser",
		},
		&requestflag.String{
			Name:     "geocode",
			Usage:    "Match latitude, longitude, and radius.",
			BodyPath: "geocode",
		},
		&requestflag.String{
			Name:     "hashtags",
			Usage:    "Hashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "hashtags",
		},
		&requestflag.Bool{
			Name:     "has-location",
			Usage:    "Require a profile location.",
			Default:  false,
			BodyPath: "hasLocation",
		},
		&requestflag.Bool{
			Name:     "has-media-only",
			Usage:    "Return only replies with media.",
			Default:  false,
			BodyPath: "hasMediaOnly",
		},
		&requestflag.Bool{
			Name:     "has-website",
			Usage:    "Require a profile website.",
			Default:  false,
			BodyPath: "hasWebsite",
		},
		&requestflag.Bool{
			Name:     "include-original-post",
			Usage:    "Include the source post in reply results.",
			Default:  false,
			BodyPath: "includeOriginalPost",
		},
		&requestflag.Bool{
			Name:     "include-search-terms",
			Usage:    "Add matching search terms to collection metadata.",
			Default:  false,
			BodyPath: "includeSearchTerms",
		},
		&requestflag.Bool{
			Name:     "include-target-metadata",
			Usage:    "Add source target metadata to each result.",
			Default:  true,
			BodyPath: "includeTargetMetadata",
		},
		&requestflag.String{
			Name:     "in-reply-to-tweet-id",
			Usage:    "Only replies to this tweet ID (tweet_search_extractor)",
			BodyPath: "inReplyToTweetId",
		},
		&requestflag.String{
			Name:     "language",
			Usage:    "Language code filter (tweet_search_extractor)",
			BodyPath: "language",
		},
		&requestflag.String{
			Name:     "list-id",
			Usage:    "Search within a list ID (tweet_search_extractor)",
			BodyPath: "listId",
		},
		&requestflag.String{
			Name:     "location-contains",
			Usage:    "Required profile location text.",
			BodyPath: "locationContains",
		},
		&requestflag.Int{
			Name:     "max-depth",
			Usage:    "Maximum nested reply depth.",
			BodyPath: "maxDepth",
		},
		&requestflag.Int{
			Name:     "max-followers",
			Usage:    "Maximum follower count for profile results.",
			BodyPath: "maxFollowers",
		},
		&requestflag.Int{
			Name:     "max-following",
			Usage:    "Maximum following count for profile results.",
			BodyPath: "maxFollowing",
		},
		&requestflag.String{
			Name:     "max-id",
			Usage:    "Return Tweets older than this Tweet ID.",
			BodyPath: "maxId",
		},
		&requestflag.Int{
			Name:     "max-items-per-target",
			Usage:    "Maximum results collected for each target.",
			BodyPath: "maxItemsPerTarget",
		},
		&requestflag.Int{
			Name:     "max-likes",
			Usage:    "Maximum Tweet like count.",
			BodyPath: "maxLikes",
		},
		&requestflag.Int{
			Name:     "max-pages-per-target",
			Usage:    "Reply pages collected for each target.",
			BodyPath: "maxPagesPerTarget",
		},
		&requestflag.Int{
			Name:     "max-posts",
			Usage:    "Maximum post count for profile results.",
			BodyPath: "maxPosts",
		},
		&requestflag.Int{
			Name:     "max-quotes",
			Usage:    "Maximum Tweet quote count.",
			BodyPath: "maxQuotes",
		},
		&requestflag.Int{
			Name:     "max-replies",
			Usage:    "Maximum Tweet reply count.",
			BodyPath: "maxReplies",
		},
		&requestflag.Int{
			Name:     "max-retweets",
			Usage:    "Maximum Tweet repost count.",
			BodyPath: "maxRetweets",
		},
		&requestflag.String{
			Name:     "media-type",
			Usage:    "Media type filter (tweet_search_extractor)",
			BodyPath: "mediaType",
		},
		&requestflag.String{
			Name:     "mentioning",
			Usage:    "Filter tweets mentioning a username (tweet_search_extractor)",
			BodyPath: "mentioning",
		},
		&requestflag.Int{
			Name:     "min-account-age-days",
			Usage:    "Minimum profile age in days.",
			BodyPath: "minAccountAgeDays",
		},
		&requestflag.Int{
			Name:     "min-bookmarks",
			Usage:    "Minimum Tweet bookmark count.",
			BodyPath: "minBookmarks",
		},
		&requestflag.Int{
			Name:     "min-faves",
			Usage:    "Minimum likes threshold (tweet_search_extractor)",
			BodyPath: "minFaves",
		},
		&requestflag.Int{
			Name:     "min-followers",
			Usage:    "Minimum follower count for profile results.",
			BodyPath: "minFollowers",
		},
		&requestflag.Int{
			Name:     "min-following",
			Usage:    "Minimum following count for profile results.",
			BodyPath: "minFollowing",
		},
		&requestflag.Int{
			Name:     "min-posts",
			Usage:    "Minimum post count for profile results.",
			BodyPath: "minPosts",
		},
		&requestflag.Int{
			Name:     "min-quotes",
			Usage:    "Minimum quote count threshold (tweet_search_extractor)",
			BodyPath: "minQuotes",
		},
		&requestflag.Int{
			Name:     "min-replies",
			Usage:    "Minimum replies threshold (tweet_search_extractor)",
			BodyPath: "minReplies",
		},
		&requestflag.Int{
			Name:     "min-retweets",
			Usage:    "Minimum retweets threshold (tweet_search_extractor)",
			BodyPath: "minRetweets",
		},
		&requestflag.Int{
			Name:     "min-views",
			Usage:    "Minimum Tweet view count.",
			BodyPath: "minViews",
		},
		&requestflag.Bool{
			Name:     "native-retweets",
			Usage:    "Only return native reposts.",
			Default:  false,
			BodyPath: "nativeRetweets",
		},
		&requestflag.String{
			Name:     "near",
			Usage:    "Match a place name.",
			BodyPath: "near",
		},
		&requestflag.Bool{
			Name:     "news",
			Usage:    "Only return news results.",
			Default:  false,
			BodyPath: "news",
		},
		&requestflag.Bool{
			Name:     "overlap-mode",
			Usage:    "Shortcut for dedupeMode=merge.",
			Default:  false,
			BodyPath: "overlapMode",
		},
		&requestflag.String{
			Name:     "place",
			Usage:    "Search within a place ID (tweet_search_extractor)",
			BodyPath: "place",
		},
		&requestflag.String{
			Name:     "place-country",
			Usage:    "Search within a country code (tweet_search_extractor)",
			BodyPath: "placeCountry",
		},
		&requestflag.String{
			Name:     "point-radius",
			Usage:    "Geo point radius, e.g. -73.99 40.73 25mi (tweet_search_extractor)",
			BodyPath: "pointRadius",
		},
		&requestflag.String{
			Name:     "query-type",
			Usage:    "Search ranking applied to every query.",
			Default:  "Latest",
			BodyPath: "queryType",
		},
		&requestflag.String{
			Name:     "quotes",
			Usage:    "Quote mode (tweet_search_extractor)",
			BodyPath: "quotes",
		},
		&requestflag.String{
			Name:     "quotes-of-tweet-id",
			Usage:    "Only quotes of this tweet ID (tweet_search_extractor)",
			BodyPath: "quotesOfTweetId",
		},
		&requestflag.ObjectList{
			Name:     "relation-target",
			Usage:    "Profile relations processed within one job.",
			BodyPath: "relationTargets",
		},
		&requestflag.String{
			Name:     "replies",
			Usage:    "Reply mode (tweet_search_extractor)",
			BodyPath: "replies",
		},
		&requestflag.Int{
			Name:     "results-limit",
			Usage:    "Maximum number of results to extract. When set, the extraction stops after reaching this limit.",
			BodyPath: "resultsLimit",
		},
		&requestflag.String{
			Name:     "retweets",
			Usage:    "Retweet mode (tweet_search_extractor)",
			BodyPath: "retweets",
		},
		&requestflag.String{
			Name:     "retweets-of-tweet-id",
			Usage:    "Only retweets of this tweet ID (tweet_search_extractor)",
			BodyPath: "retweetsOfTweetId",
		},
		&requestflag.Bool{
			Name:     "safe",
			Usage:    "Enable the safe-search filter.",
			Default:  false,
			BodyPath: "safe",
		},
		&requestflag.String{
			Name:     "scope",
			Usage:    "Reply depth scope.",
			Default:  "all",
			BodyPath: "scope",
		},
		&requestflag.StringList{
			Name:     "search-query",
			Usage:    "Search queries processed as one collection job.",
			BodyPath: "searchQueries",
		},
		&requestflag.String{
			Name:     "since-date",
			Usage:    "Start date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "sinceDate",
		},
		&requestflag.String{
			Name:     "since-id",
			Usage:    "Return Tweets newer than this Tweet ID.",
			BodyPath: "sinceId",
		},
		&requestflag.JSON{
			Name:     "since-time",
			Usage:    "Reply start time as ISO 8601 or Unix seconds.",
			BodyPath: "sinceTime",
		},
		&requestflag.String{
			Name:     "sort",
			Usage:    "Reply result order.",
			Default:  "relevance",
			BodyPath: "sort",
		},
		&requestflag.String{
			Name:     "source",
			Usage:    "Match the source application.",
			BodyPath: "source",
		},
		&requestflag.String{
			Name:     "start-cursor",
			Usage:    "Resume one reply target from this cursor.",
			BodyPath: "startCursor",
		},
		&requestflag.String{
			Name:     "target-community-id",
			Usage:    "Required for community_post_extractor & community_search.",
			BodyPath: "targetCommunityId",
		},
		&requestflag.String{
			Name:     "target-list-id",
			Usage:    "Required for list_follower_explorer, list_member_extractor & list_post_extractor.",
			BodyPath: "targetListId",
		},
		&requestflag.JSONList{
			Name:     "target",
			Usage:    "Mixed targets auto-routed within one job.",
			BodyPath: "targets",
		},
		&requestflag.String{
			Name:     "target-space-id",
			Usage:    "Required for space_explorer.",
			BodyPath: "targetSpaceId",
		},
		&requestflag.String{
			Name:     "target-tweet-id",
			BodyPath: "targetTweetId",
		},
		&requestflag.String{
			Name:     "target-username",
			BodyPath: "targetUsername",
		},
		&requestflag.String{
			Name:     "to-user",
			Usage:    "Filter replies sent to a username (tweet_search_extractor)",
			BodyPath: "toUser",
		},
		&requestflag.String{
			Name:     "until-date",
			Usage:    "End date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "untilDate",
		},
		&requestflag.JSON{
			Name:     "until-time",
			Usage:    "Reply end time as ISO 8601 or Unix seconds.",
			BodyPath: "untilTime",
		},
		&requestflag.String{
			Name:     "url",
			Usage:    "URL substring or domain filter (tweet_search_extractor)",
			BodyPath: "url",
		},
		&requestflag.String{
			Name:     "username-contains",
			Usage:    "Required username text.",
			BodyPath: "usernameContains",
		},
		&requestflag.Bool{
			Name:     "verified-only",
			Usage:    "Only verified authors (tweet_search_extractor)",
			BodyPath: "verifiedOnly",
		},
		&requestflag.String{
			Name:     "verified-type",
			Usage:    "Exact profile verification type.",
			BodyPath: "verifiedType",
		},
		&requestflag.String{
			Name:     "within",
			Usage:    "Set the radius for the near filter.",
			BodyPath: "within",
		},
		&requestflag.String{
			Name:     "within-time",
			Usage:    "Match Tweets inside a recent time window.",
			BodyPath: "withinTime",
		},
	},
	Action:          handleExtractionsEstimateCost,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"relation-target": {
		&requestflag.InnerString{
			Name:       "relation-target.relation",
			Usage:      `Choose "community_members", "followers", "following", "list_followers", "list_members" or "verified_followers".`,
			InnerField: "relation",
		},
		&requestflag.InnerString{
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "format",
			Usage:     "Export file format",
			Required:  true,
			QueryPath: "format",
		},
		&requestflag.Bool{
			Name:      "has-description",
			Usage:     "Require a non-empty description.",
			QueryPath: "hasDescription",
		},
		&requestflag.Bool{
			Name:      "has-location",
			Usage:     "Require a non-empty location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Bool{
			Name:      "has-media",
			Usage:     "Require media.",
			QueryPath: "hasMedia",
		},
		&requestflag.String{
			Name:      "lang",
			Usage:     "Filter by language code.",
			QueryPath: "lang",
		},
		&requestflag.Int{
			Name:      "max-followers",
			Usage:     "Maximum follower count.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Int{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Int{
			Name:      "max-posts",
			Usage:     "Maximum post count.",
			QueryPath: "maxPosts",
		},
		&requestflag.Int{
			Name:      "min-followers",
			Usage:     "Minimum follower count.",
			QueryPath: "minFollowers",
		},
		&requestflag.Int{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Int{
			Name:      "min-likes",
			Usage:     "Minimum like count.",
			QueryPath: "minLikes",
		},
		&requestflag.Int{
			Name:      "min-posts",
			Usage:     "Minimum post count.",
			QueryPath: "minPosts",
		},
		&requestflag.Int{
			Name:      "min-replies",
			Usage:     "Minimum reply count.",
			QueryPath: "minReplies",
		},
		&requestflag.Int{
			Name:      "min-retweets",
			Usage:     "Minimum repost count.",
			QueryPath: "minRetweets",
		},
		&requestflag.Int{
			Name:      "min-views",
			Usage:     "Minimum view count.",
			QueryPath: "minViews",
		},
		&requestflag.String{
			Name:      "search",
			Usage:     "Search exported result text.",
			QueryPath: "search",
		},
		&requestflag.String{
			Name:      "since-date",
			Usage:     "Include results on or after this date.",
			QueryPath: "sinceDate",
		},
		&requestflag.String{
			Name:      "until-date",
			Usage:     "Include results on or before this date.",
			QueryPath: "untilDate",
		},
		&requestflag.Bool{
			Name:      "verified",
			Usage:     "Filter by verified status.",
			QueryPath: "verified",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:     "tool-type",
			Usage:    "Identifier for the extraction tool used to run a job.",
			Required: true,
			BodyPath: "toolType",
		},
		&requestflag.Bool{
			Name:      "dry-run",
			Usage:     "Estimate cost without creating an extraction.",
			Default:   false,
			QueryPath: "dry_run",
		},
		&requestflag.String{
			Name:     "advanced-query",
			Usage:    "Raw advanced search query appended as-is (tweet_search_extractor)",
			BodyPath: "advancedQuery",
		},
		&requestflag.String{
			Name:     "any-words",
			Usage:    "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "anyWords",
		},
		&requestflag.String{
			Name:     "bio-contains",
			Usage:    "Bio terms separated by commas or lines.",
			BodyPath: "bioContains",
		},
		&requestflag.Bool{
			Name:     "blue-verified-only",
			Usage:    "Return only Blue-verified Tweet authors.",
			Default:  false,
			BodyPath: "blueVerifiedOnly",
		},
		&requestflag.String{
			Name:     "bounding-box",
			Usage:    "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8 (tweet_search_extractor)",
			BodyPath: "boundingBox",
		},
		&requestflag.String{
			Name:     "card-name",
			Usage:    "Match the Tweet card name.",
			BodyPath: "cardName",
		},
		&requestflag.String{
			Name:     "cashtags",
			Usage:    "Cashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "cashtags",
		},
		&requestflag.String{
			Name:     "collection-strategy",
			Usage:    "Reply collection strategy.",
			Default:  "auto",
			BodyPath: "collectionStrategy",
		},
		&requestflag.String{
			Name:     "conversation-id",
			Usage:    "Conversation ID filter (tweet_search_extractor)",
			BodyPath: "conversationId",
		},
		&requestflag.Bool{
			Name:     "dedupe-across-targets",
			Usage:    "Merge duplicate results across collection targets.",
			Default:  true,
			BodyPath: "dedupeAcrossTargets",
		},
		&requestflag.String{
			Name:     "dedupe-mode",
			Usage:    "Keep target duplicates, first rows, or merged overlap.",
			BodyPath: "dedupeMode",
		},
		&requestflag.String{
			Name:     "exact-phrase",
			Usage:    "Exact phrase to match (tweet_search_extractor)",
			BodyPath: "exactPhrase",
		},
		&requestflag.Bool{
			Name:     "exclude-original-author",
			Usage:    "Exclude replies from the source author.",
			Default:  false,
			BodyPath: "excludeOriginalAuthor",
		},
		&requestflag.String{
			Name:     "exclude-source",
			Usage:    "Exclude a source application.",
			BodyPath: "excludeSource",
		},
		&requestflag.String{
			Name:     "exclude-words",
			Usage:    "Words or quoted phrases to exclude. Separate with spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "excludeWords",
		},
		&requestflag.String{
			Name:     "from-user",
			Usage:    "Filter by author username (tweet_search_extractor)",
			BodyPath: "fromUser",
		},
		&requestflag.String{
			Name:     "geocode",
			Usage:    "Match latitude, longitude, and radius.",
			BodyPath: "geocode",
		},
		&requestflag.String{
			Name:     "hashtags",
			Usage:    "Hashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "hashtags",
		},
		&requestflag.Bool{
			Name:     "has-location",
			Usage:    "Require a profile location.",
			Default:  false,
			BodyPath: "hasLocation",
		},
		&requestflag.Bool{
			Name:     "has-media-only",
			Usage:    "Return only replies with media.",
			Default:  false,
			BodyPath: "hasMediaOnly",
		},
		&requestflag.Bool{
			Name:     "has-website",
			Usage:    "Require a profile website.",
			Default:  false,
			BodyPath: "hasWebsite",
		},
		&requestflag.Bool{
			Name:     "include-original-post",
			Usage:    "Include the source post in reply results.",
			Default:  false,
			BodyPath: "includeOriginalPost",
		},
		&requestflag.Bool{
			Name:     "include-search-terms",
			Usage:    "Add matching search terms to collection metadata.",
			Default:  false,
			BodyPath: "includeSearchTerms",
		},
		&requestflag.Bool{
			Name:     "include-target-metadata",
			Usage:    "Add source target metadata to each result.",
			Default:  true,
			BodyPath: "includeTargetMetadata",
		},
		&requestflag.String{
			Name:     "in-reply-to-tweet-id",
			Usage:    "Only replies to this tweet ID (tweet_search_extractor)",
			BodyPath: "inReplyToTweetId",
		},
		&requestflag.String{
			Name:     "language",
			Usage:    "Language code filter (tweet_search_extractor)",
			BodyPath: "language",
		},
		&requestflag.String{
			Name:     "list-id",
			Usage:    "Search within a list ID (tweet_search_extractor)",
			BodyPath: "listId",
		},
		&requestflag.String{
			Name:     "location-contains",
			Usage:    "Required profile location text.",
			BodyPath: "locationContains",
		},
		&requestflag.Int{
			Name:     "max-depth",
			Usage:    "Maximum nested reply depth.",
			BodyPath: "maxDepth",
		},
		&requestflag.Int{
			Name:     "max-followers",
			Usage:    "Maximum follower count for profile results.",
			BodyPath: "maxFollowers",
		},
		&requestflag.Int{
			Name:     "max-following",
			Usage:    "Maximum following count for profile results.",
			BodyPath: "maxFollowing",
		},
		&requestflag.String{
			Name:     "max-id",
			Usage:    "Return Tweets older than this Tweet ID.",
			BodyPath: "maxId",
		},
		&requestflag.Int{
			Name:     "max-items-per-target",
			Usage:    "Maximum results collected for each target.",
			BodyPath: "maxItemsPerTarget",
		},
		&requestflag.Int{
			Name:     "max-likes",
			Usage:    "Maximum Tweet like count.",
			BodyPath: "maxLikes",
		},
		&requestflag.Int{
			Name:     "max-pages-per-target",
			Usage:    "Reply pages collected for each target.",
			BodyPath: "maxPagesPerTarget",
		},
		&requestflag.Int{
			Name:     "max-posts",
			Usage:    "Maximum post count for profile results.",
			BodyPath: "maxPosts",
		},
		&requestflag.Int{
			Name:     "max-quotes",
			Usage:    "Maximum Tweet quote count.",
			BodyPath: "maxQuotes",
		},
		&requestflag.Int{
			Name:     "max-replies",
			Usage:    "Maximum Tweet reply count.",
			BodyPath: "maxReplies",
		},
		&requestflag.Int{
			Name:     "max-retweets",
			Usage:    "Maximum Tweet repost count.",
			BodyPath: "maxRetweets",
		},
		&requestflag.String{
			Name:     "media-type",
			Usage:    "Media type filter (tweet_search_extractor)",
			BodyPath: "mediaType",
		},
		&requestflag.String{
			Name:     "mentioning",
			Usage:    "Filter tweets mentioning a username (tweet_search_extractor)",
			BodyPath: "mentioning",
		},
		&requestflag.Int{
			Name:     "min-account-age-days",
			Usage:    "Minimum profile age in days.",
			BodyPath: "minAccountAgeDays",
		},
		&requestflag.Int{
			Name:     "min-bookmarks",
			Usage:    "Minimum Tweet bookmark count.",
			BodyPath: "minBookmarks",
		},
		&requestflag.Int{
			Name:     "min-faves",
			Usage:    "Minimum likes threshold (tweet_search_extractor)",
			BodyPath: "minFaves",
		},
		&requestflag.Int{
			Name:     "min-followers",
			Usage:    "Minimum follower count for profile results.",
			BodyPath: "minFollowers",
		},
		&requestflag.Int{
			Name:     "min-following",
			Usage:    "Minimum following count for profile results.",
			BodyPath: "minFollowing",
		},
		&requestflag.Int{
			Name:     "min-posts",
			Usage:    "Minimum post count for profile results.",
			BodyPath: "minPosts",
		},
		&requestflag.Int{
			Name:     "min-quotes",
			Usage:    "Minimum quote count threshold (tweet_search_extractor)",
			BodyPath: "minQuotes",
		},
		&requestflag.Int{
			Name:     "min-replies",
			Usage:    "Minimum replies threshold (tweet_search_extractor)",
			BodyPath: "minReplies",
		},
		&requestflag.Int{
			Name:     "min-retweets",
			Usage:    "Minimum retweets threshold (tweet_search_extractor)",
			BodyPath: "minRetweets",
		},
		&requestflag.Int{
			Name:     "min-views",
			Usage:    "Minimum Tweet view count.",
			BodyPath: "minViews",
		},
		&requestflag.Bool{
			Name:     "native-retweets",
			Usage:    "Only return native reposts.",
			Default:  false,
			BodyPath: "nativeRetweets",
		},
		&requestflag.String{
			Name:     "near",
			Usage:    "Match a place name.",
			BodyPath: "near",
		},
		&requestflag.Bool{
			Name:     "news",
			Usage:    "Only return news results.",
			Default:  false,
			BodyPath: "news",
		},
		&requestflag.Bool{
			Name:     "overlap-mode",
			Usage:    "Shortcut for dedupeMode=merge.",
			Default:  false,
			BodyPath: "overlapMode",
		},
		&requestflag.String{
			Name:     "place",
			Usage:    "Search within a place ID (tweet_search_extractor)",
			BodyPath: "place",
		},
		&requestflag.String{
			Name:     "place-country",
			Usage:    "Search within a country code (tweet_search_extractor)",
			BodyPath: "placeCountry",
		},
		&requestflag.String{
			Name:     "point-radius",
			Usage:    "Geo point radius, e.g. -73.99 40.73 25mi (tweet_search_extractor)",
			BodyPath: "pointRadius",
		},
		&requestflag.String{
			Name:     "query-type",
			Usage:    "Search ranking applied to every query.",
			Default:  "Latest",
			BodyPath: "queryType",
		},
		&requestflag.String{
			Name:     "quotes",
			Usage:    "Quote mode (tweet_search_extractor)",
			BodyPath: "quotes",
		},
		&requestflag.String{
			Name:     "quotes-of-tweet-id",
			Usage:    "Only quotes of this tweet ID (tweet_search_extractor)",
			BodyPath: "quotesOfTweetId",
		},
		&requestflag.ObjectList{
			Name:     "relation-target",
			Usage:    "Profile relations processed within one job.",
			BodyPath: "relationTargets",
		},
		&requestflag.String{
			Name:     "replies",
			Usage:    "Reply mode (tweet_search_extractor)",
			BodyPath: "replies",
		},
		&requestflag.Int{
			Name:     "results-limit",
			Usage:    "Maximum number of results to extract. When set, the extraction stops after reaching this limit.",
			BodyPath: "resultsLimit",
		},
		&requestflag.String{
			Name:     "retweets",
			Usage:    "Retweet mode (tweet_search_extractor)",
			BodyPath: "retweets",
		},
		&requestflag.String{
			Name:     "retweets-of-tweet-id",
			Usage:    "Only retweets of this tweet ID (tweet_search_extractor)",
			BodyPath: "retweetsOfTweetId",
		},
		&requestflag.Bool{
			Name:     "safe",
			Usage:    "Enable the safe-search filter.",
			Default:  false,
			BodyPath: "safe",
		},
		&requestflag.String{
			Name:     "scope",
			Usage:    "Reply depth scope.",
			Default:  "all",
			BodyPath: "scope",
		},
		&requestflag.StringList{
			Name:     "search-query",
			Usage:    "Search queries processed as one collection job.",
			BodyPath: "searchQueries",
		},
		&requestflag.String{
			Name:     "since-date",
			Usage:    "Start date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "sinceDate",
		},
		&requestflag.String{
			Name:     "since-id",
			Usage:    "Return Tweets newer than this Tweet ID.",
			BodyPath: "sinceId",
		},
		&requestflag.JSON{
			Name:     "since-time",
			Usage:    "Reply start time as ISO 8601 or Unix seconds.",
			BodyPath: "sinceTime",
		},
		&requestflag.String{
			Name:     "sort",
			Usage:    "Reply result order.",
			Default:  "relevance",
			BodyPath: "sort",
		},
		&requestflag.String{
			Name:     "source",
			Usage:    "Match the source application.",
			BodyPath: "source",
		},
		&requestflag.String{
			Name:     "start-cursor",
			Usage:    "Resume one reply target from this cursor.",
			BodyPath: "startCursor",
		},
		&requestflag.String{
			Name:     "target-community-id",
			Usage:    "Required for community_post_extractor & community_search.",
			BodyPath: "targetCommunityId",
		},
		&requestflag.String{
			Name:     "target-list-id",
			Usage:    "Required for list_follower_explorer, list_member_extractor & list_post_extractor.",
			BodyPath: "targetListId",
		},
		&requestflag.JSONList{
			Name:     "target",
			Usage:    "Mixed targets auto-routed within one job.",
			BodyPath: "targets",
		},
		&requestflag.String{
			Name:     "target-space-id",
			Usage:    "Required for space_explorer.",
			BodyPath: "targetSpaceId",
		},
		&requestflag.String{
			Name:     "target-tweet-id",
			BodyPath: "targetTweetId",
		},
		&requestflag.String{
			Name:     "target-username",
			BodyPath: "targetUsername",
		},
		&requestflag.String{
			Name:     "to-user",
			Usage:    "Filter replies sent to a username (tweet_search_extractor)",
			BodyPath: "toUser",
		},
		&requestflag.String{
			Name:     "until-date",
			Usage:    "End date YYYY-MM-DD (tweet_search_extractor)",
			BodyPath: "untilDate",
		},
		&requestflag.JSON{
			Name:     "until-time",
			Usage:    "Reply end time as ISO 8601 or Unix seconds.",
			BodyPath: "untilTime",
		},
		&requestflag.String{
			Name:     "url",
			Usage:    "URL substring or domain filter (tweet_search_extractor)",
			BodyPath: "url",
		},
		&requestflag.String{
			Name:     "username-contains",
			Usage:    "Required username text.",
			BodyPath: "usernameContains",
		},
		&requestflag.Bool{
			Name:     "verified-only",
			Usage:    "Only verified authors (tweet_search_extractor)",
			BodyPath: "verifiedOnly",
		},
		&requestflag.String{
			Name:     "verified-type",
			Usage:    "Exact profile verification type.",
			BodyPath: "verifiedType",
		},
		&requestflag.String{
			Name:     "within",
			Usage:    "Set the radius for the near filter.",
			BodyPath: "within",
		},
		&requestflag.String{
			Name:     "within-time",
			Usage:    "Match Tweets inside a recent time window.",
			BodyPath: "withinTime",
		},
	},
	Action:          handleExtractionsRun,
	HideHelpCommand: true,
}, map[string][]requestflag.HasOuterFlag{
	"relation-target": {
		&requestflag.InnerString{
			Name:       "relation-target.relation",
			Usage:      `Choose "community_members", "followers", "following", "list_followers", "list_members" or "verified_followers".`,
			InnerField: "relation",
		},
		&requestflag.InnerString{
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ExtractionGetParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Extractions.Get(
		ctx,
		cmd.String("id"),
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ExtractionExportResultsParams{}

	response, err := client.Extractions.ExportResults(
		ctx,
		cmd.String("id"),
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
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
