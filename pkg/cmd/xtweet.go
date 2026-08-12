// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var xTweetsCreate = cli.Command{
	Name:    "create",
	Usage:   "Create tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "account",
			Usage:    "X account (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.Flag[string]{
			Name:       "idempotency-key",
			Required:   true,
			HeaderPath: "Idempotency-Key",
		},
		&requestflag.Flag[string]{
			Name:     "community-id",
			BodyPath: "community_id",
		},
		&requestflag.Flag[bool]{
			Name:     "is-note-tweet",
			BodyPath: "is_note_tweet",
		},
		&requestflag.Flag[[]string]{
			Name:     "media",
			Usage:    "Array of public media URLs to attach. Supports up to 4 images or exactly 1 MP4 video up to 100 MB. Each URL must be publicly reachable. Attached media adds 2 credits per started MB across all files.",
			BodyPath: "media",
		},
		&requestflag.Flag[string]{
			Name:     "reply-to-tweet-id",
			BodyPath: "reply_to_tweet_id",
		},
		&requestflag.Flag[string]{
			Name:     "text",
			Usage:    "Tweet text (optional when media is provided)",
			BodyPath: "text",
		},
	},
	Action:          handleXTweetsCreate,
	HideHelpCommand: true,
}

var xTweetsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get tweet with full text, author, metrics and media",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
	},
	Action:          handleXTweetsRetrieve,
	HideHelpCommand: true,
}

var xTweetsList = cli.Command{
	Name:    "list",
	Usage:   "Get multiple tweets by IDs",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "ids",
			Usage:     "Comma-separated tweet IDs (max 100)",
			Required:  true,
			QueryPath: "ids",
		},
	},
	Action:          handleXTweetsList,
	HideHelpCommand: true,
}

var xTweetsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Delete tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:     "account",
			Usage:    "X account identifier (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.Flag[string]{
			Name:       "idempotency-key",
			Required:   true,
			HeaderPath: "Idempotency-Key",
		},
	},
	Action:          handleXTweetsDelete,
	HideHelpCommand: true,
}

var xTweetsGetFavoriters = cli.Command{
	Name:    "get-favoriters",
	Usage:   "Returns liker profiles that X makes visible for the post. X can withhold liker\nidentities even when the post reports likes. In that case this endpoint returns\n424 `favoriters_unavailable` instead of a misleading empty success.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for favoriters",
			QueryPath: "cursor",
		},
		&requestflag.Flag[bool]{
			Name:      "has-location",
			Usage:     "Only return profiles with a location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Flag[bool]{
			Name:      "has-website",
			Usage:     "Only return profiles with a website.",
			QueryPath: "hasWebsite",
		},
		&requestflag.Flag[string]{
			Name:      "location-contains",
			Usage:     "Match a location substring, ignoring case.",
			QueryPath: "locationContains",
		},
		&requestflag.Flag[int64]{
			Name:      "max-followers",
			Usage:     "Maximum follower count. Missing counts pass this maximum.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Flag[int64]{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Flag[int64]{
			Name:      "max-statuses",
			Usage:     "Maximum post count. maxPosts is also accepted.",
			QueryPath: "maxStatuses",
		},
		&requestflag.Flag[int64]{
			Name:      "min-account-age-days",
			Usage:     "Minimum account age in whole days.",
			QueryPath: "minAccountAgeDays",
		},
		&requestflag.Flag[int64]{
			Name:      "min-followers",
			Usage:     "Minimum follower count. Filtering happens before billing.",
			QueryPath: "minFollowers",
		},
		&requestflag.Flag[int64]{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Flag[int64]{
			Name:      "min-statuses",
			Usage:     "Minimum post count. minPosts is also accepted.",
			QueryPath: "minStatuses",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum user profiles requested from this page (20-200, default 200). Source, filters, or credits can return fewer profiles. Keep requesting next_cursor while has_next_page is true. Deprecated aliases remain accepted.\n",
			Default:   200,
			QueryPath: "pageSize",
		},
		&requestflag.Flag[string]{
			Name:      "username-contains",
			Usage:     "Match a username substring, ignoring case.",
			QueryPath: "usernameContains",
		},
		&requestflag.Flag[bool]{
			Name:      "verified-only",
			Usage:     "Only return verified profiles.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "verified-type",
			Usage:     "Match the verification type exactly, ignoring case.",
			QueryPath: "verifiedType",
		},
	},
	Action:          handleXTweetsGetFavoriters,
	HideHelpCommand: true,
}

var xTweetsGetQuotes = cli.Command{
	Name:    "get-quotes",
	Usage:   "List quote tweets of a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Flag[bool]{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "card-name",
			Usage:     "Match the Tweet card name.",
			QueryPath: "cardName",
		},
		&requestflag.Flag[string]{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for quote tweets",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.Flag[string]{
			Name:      "exclude-source",
			Usage:     "Exclude a source application.",
			QueryPath: "excludeSource",
		},
		&requestflag.Flag[string]{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.Flag[string]{
			Name:      "geocode",
			Usage:     "Match latitude, longitude, and radius.",
			QueryPath: "geocode",
		},
		&requestflag.Flag[string]{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.Flag[bool]{
			Name:      "include-replies",
			Usage:     "Include reply quotes (default false)",
			QueryPath: "includeReplies",
		},
		&requestflag.Flag[string]{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Flag[int64]{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.Flag[string]{
			Name:      "max-id",
			Usage:     "Return Tweets older than this Tweet ID.",
			QueryPath: "maxId",
		},
		&requestflag.Flag[int64]{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.Flag[string]{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Flag[int64]{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Flag[int64]{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Flag[int64]{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Flag[int64]{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.Flag[bool]{
			Name:      "native-retweets",
			Usage:     "Only return native reposts.",
			QueryPath: "nativeRetweets",
		},
		&requestflag.Flag[string]{
			Name:      "near",
			Usage:     "Match a place name.",
			QueryPath: "near",
		},
		&requestflag.Flag[bool]{
			Name:      "news",
			Usage:     "Only return news results.",
			QueryPath: "news",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum page items (1-100, default 20). Source, filters, or credits can reduce results. Continue while has_next_page is true. Deprecated limit and count aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.Flag[string]{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.Flag[string]{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.Flag[string]{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.Flag[bool]{
			Name:      "safe",
			Usage:     "Enable the safe-search filter.",
			QueryPath: "safe",
		},
		&requestflag.Flag[any]{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:      "since-id",
			Usage:     "Return Tweets newer than this Tweet ID.",
			QueryPath: "sinceId",
		},
		&requestflag.Flag[string]{
			Name:      "since-time",
			Usage:     "Unix timestamp - return quotes posted after this time",
			QueryPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:      "source",
			Usage:     "Match the source application.",
			QueryPath: "source",
		},
		&requestflag.Flag[string]{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.Flag[any]{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.Flag[string]{
			Name:      "until-time",
			Usage:     "Unix timestamp - return quotes posted before this time",
			QueryPath: "untilTime",
		},
		&requestflag.Flag[string]{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Flag[bool]{
			Name:      "verified-only",
			Usage:     "Only return tweets from verified authors.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "within",
			Usage:     "Set the radius for the near filter.",
			QueryPath: "within",
		},
		&requestflag.Flag[string]{
			Name:      "within-time",
			Usage:     "Match Tweets inside a recent time window.",
			QueryPath: "withinTime",
		},
	},
	Action:          handleXTweetsGetQuotes,
	HideHelpCommand: true,
}

var xTweetsGetReplies = cli.Command{
	Name:    "get-replies",
	Usage:   "Returns direct replies. Omit mode for automatic maximum coverage with resumable\npagination. Complete mode returns nested replies, diagnostics, and 424 when\ndirect coverage stays below 80%.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Flag[bool]{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "card-name",
			Usage:     "Match the Tweet card name.",
			QueryPath: "cardName",
		},
		&requestflag.Flag[string]{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Cursor from the previous response. Xquik cursors resume automatic coverage. Existing unprefixed cursors keep legacy standard behavior.\n",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.Flag[bool]{
			Name:      "exclude-original-author",
			Usage:     "Exclude replies written by the source-post author.",
			Default:   false,
			QueryPath: "excludeOriginalAuthor",
		},
		&requestflag.Flag[string]{
			Name:      "exclude-source",
			Usage:     "Exclude a source application.",
			QueryPath: "excludeSource",
		},
		&requestflag.Flag[string]{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.Flag[string]{
			Name:      "geocode",
			Usage:     "Match latitude, longitude, and radius.",
			QueryPath: "geocode",
		},
		&requestflag.Flag[string]{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.Flag[bool]{
			Name:      "has-media-only",
			Usage:     "Only return replies containing media.",
			Default:   false,
			QueryPath: "hasMediaOnly",
		},
		&requestflag.Flag[bool]{
			Name:      "include-original-post",
			Usage:     "Include the source post and count it toward limit.",
			Default:   false,
			QueryPath: "includeOriginalPost",
		},
		&requestflag.Flag[string]{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "With mode=complete, maximum combined direct and nested reply rows (1-25000, default 25000). Automatic pages accept 1-300. Standard pages accept 1-100. Prefer pageSize outside complete mode.\n",
			QueryPath: "limit",
		},
		&requestflag.Flag[int64]{
			Name:      "max-depth",
			Usage:     "Maximum reply depth from the source post.",
			QueryPath: "maxDepth",
		},
		&requestflag.Flag[int64]{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.Flag[string]{
			Name:      "max-id",
			Usage:     "Return Tweets older than this Tweet ID.",
			QueryPath: "maxId",
		},
		&requestflag.Flag[int64]{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.Flag[string]{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Flag[int64]{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Flag[int64]{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Flag[int64]{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Flag[int64]{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.Flag[string]{
			Name:      "mode",
			Usage:     "Optional advanced override. Omit mode for automatic maximum direct reply coverage with pagination. Standard keeps legacy pagination. Complete returns direct and nested replies with diagnostics, scope, depth, sorting, and original-post controls.\n",
			QueryPath: "mode",
		},
		&requestflag.Flag[bool]{
			Name:      "native-retweets",
			Usage:     "Only return native reposts.",
			QueryPath: "nativeRetweets",
		},
		&requestflag.Flag[string]{
			Name:      "near",
			Usage:     "Match a place name.",
			QueryPath: "near",
		},
		&requestflag.Flag[bool]{
			Name:      "news",
			Usage:     "Only return news results.",
			QueryPath: "news",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Automatic pages accept 1-300 Tweets. Standard pages keep 1-100. Default 20. Continue while has_next_page is true. Deprecated aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.Flag[string]{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.Flag[string]{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.Flag[string]{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.Flag[bool]{
			Name:      "safe",
			Usage:     "Enable the safe-search filter.",
			QueryPath: "safe",
		},
		&requestflag.Flag[string]{
			Name:      "scope",
			Usage:     "Select all replies, direct replies, or nested replies.",
			Default:   "all",
			QueryPath: "scope",
		},
		&requestflag.Flag[any]{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:      "since-id",
			Usage:     "Return Tweets newer than this Tweet ID.",
			QueryPath: "sinceId",
		},
		&requestflag.Flag[string]{
			Name:      "since-time",
			Usage:     "Unix timestamp - return replies posted after this time",
			QueryPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:      "sort",
			Usage:     "Sort the selected replies before applying limit.",
			Default:   "relevance",
			QueryPath: "sort",
		},
		&requestflag.Flag[string]{
			Name:      "source",
			Usage:     "Match the source application.",
			QueryPath: "source",
		},
		&requestflag.Flag[string]{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.Flag[any]{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.Flag[string]{
			Name:      "until-time",
			Usage:     "Unix timestamp - return replies posted before this time",
			QueryPath: "untilTime",
		},
		&requestflag.Flag[string]{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Flag[bool]{
			Name:      "verified-only",
			Usage:     "Only return tweets from verified authors.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "within",
			Usage:     "Set the radius for the near filter.",
			QueryPath: "within",
		},
		&requestflag.Flag[string]{
			Name:      "within-time",
			Usage:     "Match Tweets inside a recent time window.",
			QueryPath: "withinTime",
		},
	},
	Action:          handleXTweetsGetReplies,
	HideHelpCommand: true,
}

var xTweetsGetRetweeters = cli.Command{
	Name:    "get-retweeters",
	Usage:   "List users who retweeted a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for retweeters",
			QueryPath: "cursor",
		},
		&requestflag.Flag[bool]{
			Name:      "has-location",
			Usage:     "Only return profiles with a location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Flag[bool]{
			Name:      "has-website",
			Usage:     "Only return profiles with a website.",
			QueryPath: "hasWebsite",
		},
		&requestflag.Flag[string]{
			Name:      "location-contains",
			Usage:     "Match a location substring, ignoring case.",
			QueryPath: "locationContains",
		},
		&requestflag.Flag[int64]{
			Name:      "max-followers",
			Usage:     "Maximum follower count. Missing counts pass this maximum.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Flag[int64]{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Flag[int64]{
			Name:      "max-statuses",
			Usage:     "Maximum post count. maxPosts is also accepted.",
			QueryPath: "maxStatuses",
		},
		&requestflag.Flag[int64]{
			Name:      "min-account-age-days",
			Usage:     "Minimum account age in whole days.",
			QueryPath: "minAccountAgeDays",
		},
		&requestflag.Flag[int64]{
			Name:      "min-followers",
			Usage:     "Minimum follower count. Filtering happens before billing.",
			QueryPath: "minFollowers",
		},
		&requestflag.Flag[int64]{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Flag[int64]{
			Name:      "min-statuses",
			Usage:     "Minimum post count. minPosts is also accepted.",
			QueryPath: "minStatuses",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum user profiles requested from this page (20-200, default 200). Source, filters, or credits can return fewer profiles. Keep requesting next_cursor while has_next_page is true. Deprecated aliases remain accepted.\n",
			Default:   200,
			QueryPath: "pageSize",
		},
		&requestflag.Flag[string]{
			Name:      "username-contains",
			Usage:     "Match a username substring, ignoring case.",
			QueryPath: "usernameContains",
		},
		&requestflag.Flag[bool]{
			Name:      "verified-only",
			Usage:     "Only return verified profiles.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "verified-type",
			Usage:     "Match the verification type exactly, ignoring case.",
			QueryPath: "verifiedType",
		},
	},
	Action:          handleXTweetsGetRetweeters,
	HideHelpCommand: true,
}

var xTweetsGetThread = cli.Command{
	Name:    "get-thread",
	Usage:   "Get full conversation thread for a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Flag[bool]{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for thread tweets",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.Flag[string]{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.Flag[string]{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.Flag[string]{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Flag[int64]{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.Flag[int64]{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.Flag[string]{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Flag[int64]{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Flag[int64]{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Flag[int64]{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Flag[int64]{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum page items (1-100, default 20). Source, filters, or credits can reduce results. Continue while has_next_page is true. Deprecated limit and count aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.Flag[string]{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.Flag[string]{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.Flag[string]{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.Flag[any]{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.Flag[any]{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.Flag[string]{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Flag[bool]{
			Name:      "verified-only",
			Usage:     "Only return tweets from verified authors.",
			QueryPath: "verifiedOnly",
		},
	},
	Action:          handleXTweetsGetThread,
	HideHelpCommand: true,
}

var xTweetsSearch = cli.Command{
	Name:    "search",
	Usage:   "No-mode search maximizes coverage. New cursorless `Latest` sessions return rows\nnewest-first across cursor pages. Existing cursors preserve their established\nordering.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "q",
			Usage:     "Query, Tweet ID, or status URL. Valid inline bounds apply per page.",
			Required:  true,
			QueryPath: "q",
		},
		&requestflag.Flag[string]{
			Name:      "advanced-query",
			Usage:     "Raw advanced search query appended as-is.",
			QueryPath: "advancedQuery",
		},
		&requestflag.Flag[string]{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Flag[bool]{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "bounding-box",
			Usage:     "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8.",
			QueryPath: "boundingBox",
		},
		&requestflag.Flag[string]{
			Name:      "card-name",
			Usage:     "Match the Tweet card name.",
			QueryPath: "cardName",
		},
		&requestflag.Flag[string]{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Cursor from the previous response. Xquik cursors resume automatic coverage. Existing unprefixed cursors keep legacy standard behavior.\n",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.Flag[string]{
			Name:      "exclude-source",
			Usage:     "Exclude a source application.",
			QueryPath: "excludeSource",
		},
		&requestflag.Flag[string]{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.Flag[string]{
			Name:      "geocode",
			Usage:     "Match latitude, longitude, and radius.",
			QueryPath: "geocode",
		},
		&requestflag.Flag[string]{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.Flag[string]{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Result upper bound. Omit it for the existing 20-row page size. Explicit coverage defaults to 2000 and allows 10000. For paid requests, remaining credits can reduce results. Zero affordable results returns 402.\n",
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "list-id",
			Usage:     "Search within a list ID.",
			QueryPath: "listId",
		},
		&requestflag.Flag[int64]{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.Flag[string]{
			Name:      "max-id",
			Usage:     "Return Tweets older than this Tweet ID.",
			QueryPath: "maxId",
		},
		&requestflag.Flag[int64]{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.Flag[string]{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Flag[int64]{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Flag[int64]{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Flag[int64]{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Flag[int64]{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Flag[int64]{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.Flag[string]{
			Name:      "mode",
			Usage:     "Omit mode for resumable maximum coverage. Standard keeps legacy pagination. Coverage returns diagnostics once and rejects cursors.\n",
			QueryPath: "mode",
		},
		&requestflag.Flag[bool]{
			Name:      "native-retweets",
			Usage:     "Only return native reposts.",
			QueryPath: "nativeRetweets",
		},
		&requestflag.Flag[string]{
			Name:      "near",
			Usage:     "Match a place name.",
			QueryPath: "near",
		},
		&requestflag.Flag[bool]{
			Name:      "news",
			Usage:     "Only return news results.",
			QueryPath: "news",
		},
		&requestflag.Flag[string]{
			Name:      "place",
			Usage:     "Search within a place ID.",
			QueryPath: "place",
		},
		&requestflag.Flag[string]{
			Name:      "place-country",
			Usage:     "Search within a country code.",
			QueryPath: "placeCountry",
		},
		&requestflag.Flag[string]{
			Name:      "point-radius",
			Usage:     "Geo point radius, e.g. -73.99 40.73 25mi.",
			QueryPath: "pointRadius",
		},
		&requestflag.Flag[string]{
			Name:      "query-type",
			Usage:     "Sort order - Latest (chronological) or Top (engagement-ranked)",
			Default:   "Latest",
			QueryPath: "queryType",
		},
		&requestflag.Flag[string]{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.Flag[string]{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.Flag[string]{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.Flag[string]{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.Flag[bool]{
			Name:      "safe",
			Usage:     "Enable the safe-search filter.",
			QueryPath: "safe",
		},
		&requestflag.Flag[any]{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:      "since-id",
			Usage:     "Return Tweets newer than this Tweet ID.",
			QueryPath: "sinceId",
		},
		&requestflag.Flag[string]{
			Name:      "since-time",
			Usage:     "Inclusive ISO bound.",
			QueryPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:      "source",
			Usage:     "Match the source application.",
			QueryPath: "source",
		},
		&requestflag.Flag[string]{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.Flag[any]{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.Flag[string]{
			Name:      "until-time",
			Usage:     "Exclusive ISO bound.",
			QueryPath: "untilTime",
		},
		&requestflag.Flag[string]{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Flag[bool]{
			Name:      "verified-only",
			Usage:     "Only return tweets from verified authors.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.Flag[string]{
			Name:      "within",
			Usage:     "Set the radius for the near filter.",
			QueryPath: "within",
		},
		&requestflag.Flag[string]{
			Name:      "within-time",
			Usage:     "Match Tweets inside a recent time window.",
			QueryPath: "withinTime",
		},
	},
	Action:          handleXTweetsSearch,
	HideHelpCommand: true,
}

func handleXTweetsCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.New(ctx, params, options...)
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
		Title:          "x:tweets create",
		Transform:      transform,
	})
}

func handleXTweetsRetrieve(ctx context.Context, cmd *cli.Command) error {
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.Get(ctx, cmd.Value("id").(string), options...)
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
		Title:          "x:tweets retrieve",
		Transform:      transform,
	})
}

func handleXTweetsList(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetListParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.List(ctx, params, options...)
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
		Title:          "x:tweets list",
		Transform:      transform,
	})
}

func handleXTweetsDelete(ctx context.Context, cmd *cli.Command) error {
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
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XTweetDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.Delete(
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
		Title:          "x:tweets delete",
		Transform:      transform,
	})
}

func handleXTweetsGetFavoriters(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetGetFavoritersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetFavoriters(
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
		Title:          "x:tweets get-favoriters",
		Transform:      transform,
	})
}

func handleXTweetsGetQuotes(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetGetQuotesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetQuotes(
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
		Title:          "x:tweets get-quotes",
		Transform:      transform,
	})
}

func handleXTweetsGetReplies(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetGetRepliesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetReplies(
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
		Title:          "x:tweets get-replies",
		Transform:      transform,
	})
}

func handleXTweetsGetRetweeters(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetGetRetweetersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetRetweeters(
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
		Title:          "x:tweets get-retweeters",
		Transform:      transform,
	})
}

func handleXTweetsGetThread(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetGetThreadParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetThread(
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
		Title:          "x:tweets get-thread",
		Transform:      transform,
	})
}

func handleXTweetsSearch(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XTweetSearchParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.Search(ctx, params, options...)
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
		Title:          "x:tweets search",
		Transform:      transform,
	})
}
