// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"

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
		&requestflag.String{
			Name:     "account",
			Usage:    "X account (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.String{
			Name:       "idempotency-key",
			Required:   true,
			HeaderPath: "Idempotency-Key",
		},
		&requestflag.String{
			Name:     "community-id",
			BodyPath: "community_id",
		},
		&requestflag.Bool{
			Name:     "is-note-tweet",
			BodyPath: "is_note_tweet",
		},
		&requestflag.StringList{
			Name:     "media",
			Usage:    "Array of public media URLs to attach. Supports up to 4 images or exactly 1 MP4 video up to 100 MB. Each URL must be publicly reachable. Attached media adds 2 credits per started MB across all files.",
			BodyPath: "media",
		},
		&requestflag.String{
			Name:     "reply-to-tweet-id",
			BodyPath: "reply_to_tweet_id",
		},
		&requestflag.String{
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
		&requestflag.String{
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
		&requestflag.String{
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:     "account",
			Usage:    "X account identifier (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for favoriters",
			QueryPath: "cursor",
		},
		&requestflag.Bool{
			Name:      "has-location",
			Usage:     "Only return profiles with a location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Bool{
			Name:      "has-website",
			Usage:     "Only return profiles with a website.",
			QueryPath: "hasWebsite",
		},
		&requestflag.String{
			Name:      "location-contains",
			Usage:     "Match a location substring, ignoring case.",
			QueryPath: "locationContains",
		},
		&requestflag.Int{
			Name:      "max-followers",
			Usage:     "Maximum follower count. Missing counts pass this maximum.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Int{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Int{
			Name:      "max-statuses",
			Usage:     "Maximum post count. maxPosts is also accepted.",
			QueryPath: "maxStatuses",
		},
		&requestflag.Int{
			Name:      "min-account-age-days",
			Usage:     "Minimum account age in whole days.",
			QueryPath: "minAccountAgeDays",
		},
		&requestflag.Int{
			Name:      "min-followers",
			Usage:     "Minimum follower count. Filtering happens before billing.",
			QueryPath: "minFollowers",
		},
		&requestflag.Int{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Int{
			Name:      "min-statuses",
			Usage:     "Minimum post count. minPosts is also accepted.",
			QueryPath: "minStatuses",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Maximum user profiles requested from this page (20-200, default 200). Source, filters, or credits can return fewer profiles. Keep requesting next_cursor while has_next_page is true. Deprecated aliases remain accepted.\n",
			Default:   200,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "username-contains",
			Usage:     "Match a username substring, ignoring case.",
			QueryPath: "usernameContains",
		},
		&requestflag.Bool{
			Name:      "verified-only",
			Usage:     "Only return verified profiles.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Bool{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.String{
			Name:      "card-name",
			Usage:     "Match the Tweet card name.",
			QueryPath: "cardName",
		},
		&requestflag.String{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.String{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for quote tweets",
			QueryPath: "cursor",
		},
		&requestflag.String{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.String{
			Name:      "exclude-source",
			Usage:     "Exclude a source application.",
			QueryPath: "excludeSource",
		},
		&requestflag.String{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.String{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.String{
			Name:      "geocode",
			Usage:     "Match latitude, longitude, and radius.",
			QueryPath: "geocode",
		},
		&requestflag.String{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.Bool{
			Name:      "include-replies",
			Usage:     "Include reply quotes (default false)",
			QueryPath: "includeReplies",
		},
		&requestflag.String{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.String{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Int{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.String{
			Name:      "max-id",
			Usage:     "Return Tweets older than this Tweet ID.",
			QueryPath: "maxId",
		},
		&requestflag.Int{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Int{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Int{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.String{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.String{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Int{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Int{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Int{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Int{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Int{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Int{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.Bool{
			Name:      "native-retweets",
			Usage:     "Only return native reposts.",
			QueryPath: "nativeRetweets",
		},
		&requestflag.String{
			Name:      "near",
			Usage:     "Match a place name.",
			QueryPath: "near",
		},
		&requestflag.Bool{
			Name:      "news",
			Usage:     "Only return news results.",
			QueryPath: "news",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Maximum page items (1-100, default 20). Source, filters, or credits can reduce results. Continue while has_next_page is true. Deprecated limit and count aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.String{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.String{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.String{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.String{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.Bool{
			Name:      "safe",
			Usage:     "Enable the safe-search filter.",
			QueryPath: "safe",
		},
		&requestflag.String{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.String{
			Name:      "since-id",
			Usage:     "Return Tweets newer than this Tweet ID.",
			QueryPath: "sinceId",
		},
		&requestflag.String{
			Name:      "since-time",
			Usage:     "Unix timestamp - return quotes posted after this time",
			QueryPath: "sinceTime",
		},
		&requestflag.String{
			Name:      "source",
			Usage:     "Match the source application.",
			QueryPath: "source",
		},
		&requestflag.String{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.String{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.String{
			Name:      "until-time",
			Usage:     "Unix timestamp - return quotes posted before this time",
			QueryPath: "untilTime",
		},
		&requestflag.String{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Bool{
			Name:      "verified-only",
			Usage:     "Only return tweets from verified authors.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.String{
			Name:      "within",
			Usage:     "Set the radius for the near filter.",
			QueryPath: "within",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Bool{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.String{
			Name:      "card-name",
			Usage:     "Match the Tweet card name.",
			QueryPath: "cardName",
		},
		&requestflag.String{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.String{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Cursor from the previous response. Xquik cursors resume automatic coverage. Existing unprefixed cursors keep legacy standard behavior.\n",
			QueryPath: "cursor",
		},
		&requestflag.String{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.Bool{
			Name:      "exclude-original-author",
			Usage:     "Exclude replies written by the source-post author.",
			Default:   false,
			QueryPath: "excludeOriginalAuthor",
		},
		&requestflag.String{
			Name:      "exclude-source",
			Usage:     "Exclude a source application.",
			QueryPath: "excludeSource",
		},
		&requestflag.String{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.String{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.String{
			Name:      "geocode",
			Usage:     "Match latitude, longitude, and radius.",
			QueryPath: "geocode",
		},
		&requestflag.String{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.Bool{
			Name:      "has-media-only",
			Usage:     "Only return replies containing media.",
			Default:   false,
			QueryPath: "hasMediaOnly",
		},
		&requestflag.Bool{
			Name:      "include-original-post",
			Usage:     "Include the source post and count it toward limit.",
			Default:   false,
			QueryPath: "includeOriginalPost",
		},
		&requestflag.String{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.String{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Int{
			Name:      "limit",
			Usage:     "With mode=complete, maximum combined direct and nested reply rows (1-25000, default 25000). Automatic pages accept 1-300. Standard pages accept 1-100. Prefer pageSize outside complete mode.\n",
			QueryPath: "limit",
		},
		&requestflag.Int{
			Name:      "max-depth",
			Usage:     "Maximum reply depth from the source post.",
			QueryPath: "maxDepth",
		},
		&requestflag.Int{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.String{
			Name:      "max-id",
			Usage:     "Return Tweets older than this Tweet ID.",
			QueryPath: "maxId",
		},
		&requestflag.Int{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Int{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Int{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.String{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.String{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Int{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Int{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Int{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Int{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Int{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Int{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.String{
			Name:      "mode",
			Usage:     "Optional advanced override. Omit mode for automatic maximum direct reply coverage with pagination. Standard keeps legacy pagination. Complete returns direct and nested replies with diagnostics, scope, depth, sorting, and original-post controls.\n",
			QueryPath: "mode",
		},
		&requestflag.Bool{
			Name:      "native-retweets",
			Usage:     "Only return native reposts.",
			QueryPath: "nativeRetweets",
		},
		&requestflag.String{
			Name:      "near",
			Usage:     "Match a place name.",
			QueryPath: "near",
		},
		&requestflag.Bool{
			Name:      "news",
			Usage:     "Only return news results.",
			QueryPath: "news",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Automatic pages accept 1-300 Tweets. Standard pages keep 1-100. Default 20. Continue while has_next_page is true. Deprecated aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.String{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.String{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.String{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.String{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.Bool{
			Name:      "safe",
			Usage:     "Enable the safe-search filter.",
			QueryPath: "safe",
		},
		&requestflag.String{
			Name:      "scope",
			Usage:     "Select all replies, direct replies, or nested replies.",
			Default:   "all",
			QueryPath: "scope",
		},
		&requestflag.String{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.String{
			Name:      "since-id",
			Usage:     "Return Tweets newer than this Tweet ID.",
			QueryPath: "sinceId",
		},
		&requestflag.String{
			Name:      "since-time",
			Usage:     "Unix timestamp - return replies posted after this time",
			QueryPath: "sinceTime",
		},
		&requestflag.String{
			Name:      "sort",
			Usage:     "Sort the selected replies before applying limit.",
			Default:   "relevance",
			QueryPath: "sort",
		},
		&requestflag.String{
			Name:      "source",
			Usage:     "Match the source application.",
			QueryPath: "source",
		},
		&requestflag.String{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.String{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.String{
			Name:      "until-time",
			Usage:     "Unix timestamp - return replies posted before this time",
			QueryPath: "untilTime",
		},
		&requestflag.String{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Bool{
			Name:      "verified-only",
			Usage:     "Only return tweets from verified authors.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.String{
			Name:      "within",
			Usage:     "Set the radius for the near filter.",
			QueryPath: "within",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for retweeters",
			QueryPath: "cursor",
		},
		&requestflag.Bool{
			Name:      "has-location",
			Usage:     "Only return profiles with a location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Bool{
			Name:      "has-website",
			Usage:     "Only return profiles with a website.",
			QueryPath: "hasWebsite",
		},
		&requestflag.String{
			Name:      "location-contains",
			Usage:     "Match a location substring, ignoring case.",
			QueryPath: "locationContains",
		},
		&requestflag.Int{
			Name:      "max-followers",
			Usage:     "Maximum follower count. Missing counts pass this maximum.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Int{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Int{
			Name:      "max-statuses",
			Usage:     "Maximum post count. maxPosts is also accepted.",
			QueryPath: "maxStatuses",
		},
		&requestflag.Int{
			Name:      "min-account-age-days",
			Usage:     "Minimum account age in whole days.",
			QueryPath: "minAccountAgeDays",
		},
		&requestflag.Int{
			Name:      "min-followers",
			Usage:     "Minimum follower count. Filtering happens before billing.",
			QueryPath: "minFollowers",
		},
		&requestflag.Int{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Int{
			Name:      "min-statuses",
			Usage:     "Minimum post count. minPosts is also accepted.",
			QueryPath: "minStatuses",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Maximum user profiles requested from this page (20-200, default 200). Source, filters, or credits can return fewer profiles. Keep requesting next_cursor while has_next_page is true. Deprecated aliases remain accepted.\n",
			Default:   200,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "username-contains",
			Usage:     "Match a username substring, ignoring case.",
			QueryPath: "usernameContains",
		},
		&requestflag.Bool{
			Name:      "verified-only",
			Usage:     "Only return verified profiles.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.String{
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Bool{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.String{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.String{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for thread tweets",
			QueryPath: "cursor",
		},
		&requestflag.String{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.String{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.String{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.String{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.String{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.String{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Int{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.Int{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Int{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Int{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.String{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.String{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Int{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Int{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Int{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Int{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Int{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Int{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Maximum page items (1-100, default 20). Source, filters, or credits can reduce results. Continue while has_next_page is true. Deprecated limit and count aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.String{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.String{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.String{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.String{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.String{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.String{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.String{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.String{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Bool{
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
		&requestflag.String{
			Name:      "q",
			Usage:     "Query, Tweet ID, or status URL. Valid inline bounds apply per page.",
			Required:  true,
			QueryPath: "q",
		},
		&requestflag.String{
			Name:      "advanced-query",
			Usage:     "Raw advanced search query appended as-is.",
			QueryPath: "advancedQuery",
		},
		&requestflag.String{
			Name:      "any-words",
			Usage:     "Words or quoted phrases where any one can match. Separate with spaces, commas, or lines.",
			QueryPath: "anyWords",
		},
		&requestflag.Bool{
			Name:      "blue-verified-only",
			Usage:     "Only return tweets from Blue-verified authors.",
			QueryPath: "blueVerifiedOnly",
		},
		&requestflag.String{
			Name:      "bounding-box",
			Usage:     "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8.",
			QueryPath: "boundingBox",
		},
		&requestflag.String{
			Name:      "card-name",
			Usage:     "Match the Tweet card name.",
			QueryPath: "cardName",
		},
		&requestflag.String{
			Name:      "cashtags",
			Usage:     "Cashtags separated by spaces, commas, or lines.",
			QueryPath: "cashtags",
		},
		&requestflag.String{
			Name:      "conversation-id",
			Usage:     "Conversation ID filter.",
			QueryPath: "conversationId",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Cursor from the previous response. Xquik cursors resume automatic coverage. Existing unprefixed cursors keep legacy standard behavior.\n",
			QueryPath: "cursor",
		},
		&requestflag.String{
			Name:      "exact-phrase",
			Usage:     "Exact phrase to match.",
			QueryPath: "exactPhrase",
		},
		&requestflag.String{
			Name:      "exclude-source",
			Usage:     "Exclude a source application.",
			QueryPath: "excludeSource",
		},
		&requestflag.String{
			Name:      "exclude-words",
			Usage:     "Words or quoted phrases to exclude. Separate with spaces, commas, or lines.",
			QueryPath: "excludeWords",
		},
		&requestflag.String{
			Name:      "from-user",
			Usage:     "Filter by author username.",
			QueryPath: "fromUser",
		},
		&requestflag.String{
			Name:      "geocode",
			Usage:     "Match latitude, longitude, and radius.",
			QueryPath: "geocode",
		},
		&requestflag.String{
			Name:      "hashtags",
			Usage:     "Hashtags separated by spaces, commas, or lines.",
			QueryPath: "hashtags",
		},
		&requestflag.String{
			Name:      "in-reply-to-tweet-id",
			Usage:     "Only replies to this tweet ID.",
			QueryPath: "inReplyToTweetId",
		},
		&requestflag.String{
			Name:      "language",
			Usage:     "Language code filter, e.g. en or tr.",
			QueryPath: "language",
		},
		&requestflag.Int{
			Name:      "limit",
			Usage:     "Result upper bound. Omit it for the existing 20-row page size. Explicit coverage defaults to 2000 and allows 10000. For paid requests, remaining credits can reduce results. Zero affordable results returns 402.\n",
			QueryPath: "limit",
		},
		&requestflag.String{
			Name:      "list-id",
			Usage:     "Search within a list ID.",
			QueryPath: "listId",
		},
		&requestflag.Int{
			Name:      "max-faves",
			Usage:     "Maximum likes threshold. maxLikes is also accepted.",
			QueryPath: "maxFaves",
		},
		&requestflag.String{
			Name:      "max-id",
			Usage:     "Return Tweets older than this Tweet ID.",
			QueryPath: "maxId",
		},
		&requestflag.Int{
			Name:      "max-quotes",
			Usage:     "Maximum quotes threshold.",
			QueryPath: "maxQuotes",
		},
		&requestflag.Int{
			Name:      "max-replies",
			Usage:     "Maximum replies threshold.",
			QueryPath: "maxReplies",
		},
		&requestflag.Int{
			Name:      "max-retweets",
			Usage:     "Maximum retweets threshold.",
			QueryPath: "maxRetweets",
		},
		&requestflag.String{
			Name:      "media-type",
			Usage:     "Filter by media type.",
			QueryPath: "mediaType",
		},
		&requestflag.String{
			Name:      "mentioning",
			Usage:     "Filter tweets mentioning a username.",
			QueryPath: "mentioning",
		},
		&requestflag.Int{
			Name:      "min-bookmarks",
			Usage:     "Minimum bookmark count threshold.",
			QueryPath: "minBookmarks",
		},
		&requestflag.Int{
			Name:      "min-faves",
			Usage:     "Minimum likes threshold. minLikes is also accepted.",
			QueryPath: "minFaves",
		},
		&requestflag.Int{
			Name:      "min-quotes",
			Usage:     "Minimum quote count threshold.",
			QueryPath: "minQuotes",
		},
		&requestflag.Int{
			Name:      "min-replies",
			Usage:     "Minimum replies threshold.",
			QueryPath: "minReplies",
		},
		&requestflag.Int{
			Name:      "min-retweets",
			Usage:     "Minimum retweets threshold.",
			QueryPath: "minRetweets",
		},
		&requestflag.Int{
			Name:      "min-views",
			Usage:     "Minimum view count threshold.",
			QueryPath: "minViews",
		},
		&requestflag.String{
			Name:      "mode",
			Usage:     "Omit mode for resumable maximum coverage. Standard keeps legacy pagination. Coverage returns diagnostics once and rejects cursors.\n",
			QueryPath: "mode",
		},
		&requestflag.Bool{
			Name:      "native-retweets",
			Usage:     "Only return native reposts.",
			QueryPath: "nativeRetweets",
		},
		&requestflag.String{
			Name:      "near",
			Usage:     "Match a place name.",
			QueryPath: "near",
		},
		&requestflag.Bool{
			Name:      "news",
			Usage:     "Only return news results.",
			QueryPath: "news",
		},
		&requestflag.String{
			Name:      "place",
			Usage:     "Search within a place ID.",
			QueryPath: "place",
		},
		&requestflag.String{
			Name:      "place-country",
			Usage:     "Search within a country code.",
			QueryPath: "placeCountry",
		},
		&requestflag.String{
			Name:      "point-radius",
			Usage:     "Geo point radius, e.g. -73.99 40.73 25mi.",
			QueryPath: "pointRadius",
		},
		&requestflag.String{
			Name:      "query-type",
			Usage:     "Sort order - Latest (chronological) or Top (engagement-ranked)",
			Default:   "Latest",
			QueryPath: "queryType",
		},
		&requestflag.String{
			Name:      "quotes",
			Usage:     "Quote mode.",
			QueryPath: "quotes",
		},
		&requestflag.String{
			Name:      "quotes-of-tweet-id",
			Usage:     "Only quotes of this tweet ID.",
			QueryPath: "quotesOfTweetId",
		},
		&requestflag.String{
			Name:      "replies",
			Usage:     "Reply mode.",
			QueryPath: "replies",
		},
		&requestflag.String{
			Name:      "retweets",
			Usage:     "Retweet mode.",
			QueryPath: "retweets",
		},
		&requestflag.String{
			Name:      "retweets-of-tweet-id",
			Usage:     "Only retweets of this tweet ID.",
			QueryPath: "retweetsOfTweetId",
		},
		&requestflag.Bool{
			Name:      "safe",
			Usage:     "Enable the safe-search filter.",
			QueryPath: "safe",
		},
		&requestflag.String{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.String{
			Name:      "since-id",
			Usage:     "Return Tweets newer than this Tweet ID.",
			QueryPath: "sinceId",
		},
		&requestflag.String{
			Name:      "since-time",
			Usage:     "Inclusive ISO bound.",
			QueryPath: "sinceTime",
		},
		&requestflag.String{
			Name:      "source",
			Usage:     "Match the source application.",
			QueryPath: "source",
		},
		&requestflag.String{
			Name:      "to-user",
			Usage:     "Filter replies sent to a username.",
			QueryPath: "toUser",
		},
		&requestflag.String{
			Name:      "until-date",
			Usage:     "End date in YYYY-MM-DD format.",
			QueryPath: "untilDate",
		},
		&requestflag.String{
			Name:      "until-time",
			Usage:     "Exclusive ISO bound.",
			QueryPath: "untilTime",
		},
		&requestflag.String{
			Name:      "url",
			Usage:     "URL substring or domain filter.",
			QueryPath: "url",
		},
		&requestflag.Bool{
			Name:      "verified-only",
			Usage:     "Only return tweets from verified authors.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.String{
			Name:      "within",
			Usage:     "Set the radius for the near filter.",
			QueryPath: "within",
		},
		&requestflag.String{
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.Get(ctx, cmd.String("id"), options...)
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XTweetDeleteParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.Delete(
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XTweetGetFavoritersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetFavoriters(
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XTweetGetQuotesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetQuotes(
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XTweetGetRepliesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetReplies(
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XTweetGetRetweetersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetRetweeters(
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XTweetGetThreadParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.GetThread(
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
		Title:          "x:tweets get-thread",
		Transform:      transform,
	})
}

func handleXTweetsSearch(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
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
