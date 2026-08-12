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

var xUsersRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get user profile with follower counts and verification",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
	},
	Action:          handleXUsersRetrieve,
	HideHelpCommand: true,
}

var xUsersRemoveFollower = cli.Command{
	Name:    "remove-follower",
	Usage:   "Remove follower",
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
	Action:          handleXUsersRemoveFollower,
	HideHelpCommand: true,
}

var xUsersRetrieveBatch = cli.Command{
	Name:    "retrieve-batch",
	Usage:   "Look up multiple users by IDs in one call",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "ids",
			Usage:     "Comma-separated numeric user IDs (1-100 values). Duplicate IDs are ignored while preserving first-seen order.",
			Required:  true,
			QueryPath: "ids",
		},
	},
	Action:          handleXUsersRetrieveBatch,
	HideHelpCommand: true,
}

var xUsersRetrieveFollowers = cli.Command{
	Name:    "retrieve-followers",
	Usage:   "List followers of a user",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "after",
			Usage:     "Legacy cursor alias. Prefer cursor.",
			QueryPath: "after",
		},
		&requestflag.Flag[string]{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Cursor from the previous response. Xquik cursors resume automatic coverage. Existing unprefixed cursors keep legacy standard behavior.\n",
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
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Legacy page-size alias outside explicit coverage mode. Coverage accepts 1-10000. Prefer pageSize.\n",
			QueryPath: "limit",
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
		&requestflag.Flag[string]{
			Name:      "mode",
			Usage:     "Omit mode for resumable maximum coverage. Standard keeps legacy pagination. Coverage returns diagnostics once and rejects cursors.\n",
			QueryPath: "mode",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum user profiles: automatic 300; standard 200. Sources return fewer profiles. Continue with has_next_page.\n",
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
	Action:          handleXUsersRetrieveFollowers,
	HideHelpCommand: true,
}

var xUsersRetrieveFollowersYouKnow = cli.Command{
	Name:    "retrieve-followers-you-know",
	Usage:   "List mutual followers between you and a user",
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
			Usage:     "Pagination cursor for followers-you-know",
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
	Action:          handleXUsersRetrieveFollowersYouKnow,
	HideHelpCommand: true,
}

var xUsersRetrieveFollowing = cli.Command{
	Name:    "retrieve-following",
	Usage:   "List accounts a user follows",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "after",
			Usage:     "Deprecated following cursor alias. Prefer cursor.",
			QueryPath: "after",
		},
		&requestflag.Flag[string]{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Cursor from the previous response. Xquik cursors resume automatic coverage. Existing unprefixed cursors keep legacy standard behavior.\n",
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
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Legacy page-size alias outside explicit coverage mode. Coverage accepts 1-10000. Prefer pageSize.\n",
			QueryPath: "limit",
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
		&requestflag.Flag[string]{
			Name:      "mode",
			Usage:     "Omit mode for resumable maximum coverage. Standard keeps legacy pagination. Coverage returns diagnostics once and rejects cursors.\n",
			QueryPath: "mode",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum user profiles: automatic 300; standard 200. Sources return fewer profiles. Continue with has_next_page.\n",
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
	Action:          handleXUsersRetrieveFollowing,
	HideHelpCommand: true,
}

var xUsersRetrieveLikes = cli.Command{
	Name:    "retrieve-likes",
	Usage:   "List tweets liked by a user",
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
			Usage:     "Pagination cursor for liked tweets",
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
	Action:          handleXUsersRetrieveLikes,
	HideHelpCommand: true,
}

var xUsersRetrieveMedia = cli.Command{
	Name:    "retrieve-media",
	Usage:   "List media tweets posted by a user",
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
			Usage:     "Pagination cursor for media tweets",
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
	Action:          handleXUsersRetrieveMedia,
	HideHelpCommand: true,
}

var xUsersRetrieveMentions = cli.Command{
	Name:    "retrieve-mentions",
	Usage:   "List tweets mentioning a user",
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
			Usage:     "Pagination cursor for mentions",
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
			Usage:     "Unix timestamp - return mentions after this time",
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
			Usage:     "Unix timestamp - return mentions before this time",
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
	Action:          handleXUsersRetrieveMentions,
	HideHelpCommand: true,
}

var xUsersRetrieveReplies = cli.Command{
	Name:    "retrieve-replies",
	Usage:   "Returns target-authored posts and replies. Omit mode for automatic maximum\ncoverage. Pass next_cursor unchanged. Unprefixed cursors stay legacy. Excludes\nother-author context.",
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
			Name:      "include-parent-tweet",
			Usage:     "Include each reply's parent tweet.",
			Default:   false,
			QueryPath: "includeParentTweet",
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
	Action:          handleXUsersRetrieveReplies,
	HideHelpCommand: true,
}

var xUsersRetrieveSearch = cli.Command{
	Name:    "retrieve-search",
	Usage:   "Search users by name or username",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "q",
			Usage:     "User search query",
			Required:  true,
			QueryPath: "q",
		},
		&requestflag.Flag[string]{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for user search",
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
	Action:          handleXUsersRetrieveSearch,
	HideHelpCommand: true,
}

var xUsersRetrieveTweets = cli.Command{
	Name:    "retrieve-tweets",
	Usage:   "Omit mode for automatic maximum coverage. Pass next_cursor unchanged. Unprefixed\ncursors use legacy pagination. Shape and billing stay the same.",
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
			Name:      "include-parent-tweet",
			Usage:     "Include parent tweet for replies",
			Default:   false,
			QueryPath: "includeParentTweet",
		},
		&requestflag.Flag[bool]{
			Name:      "include-replies",
			Usage:     "Include reply tweets",
			Default:   false,
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
	Action:          handleXUsersRetrieveTweets,
	HideHelpCommand: true,
}

var xUsersRetrieveVerifiedFollowers = cli.Command{
	Name:    "retrieve-verified-followers",
	Usage:   "List verified followers of a user",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:      "after",
			Usage:     "Legacy cursor alias. Prefer cursor.",
			QueryPath: "after",
		},
		&requestflag.Flag[string]{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Cursor from the previous response. Xquik cursors resume automatic coverage. Existing unprefixed cursors keep legacy standard behavior.\n",
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
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Legacy page-size alias outside explicit coverage mode. Coverage accepts 1-10000. Prefer pageSize.\n",
			QueryPath: "limit",
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
		&requestflag.Flag[string]{
			Name:      "mode",
			Usage:     "Omit mode for resumable maximum coverage. Standard keeps legacy pagination. Coverage returns diagnostics once and rejects cursors.\n",
			QueryPath: "mode",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum user profiles: automatic 300; standard 200. Sources return fewer profiles. Continue with has_next_page.\n",
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
	Action:          handleXUsersRetrieveVerifiedFollowers,
	HideHelpCommand: true,
}

func handleXUsersRetrieve(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.X.Users.Get(ctx, cmd.Value("id").(string), options...)
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
		Title:          "x:users retrieve",
		Transform:      transform,
	})
}

func handleXUsersRemoveFollower(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserRemoveFollowerParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.RemoveFollower(
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
		Title:          "x:users remove-follower",
		Transform:      transform,
	})
}

func handleXUsersRetrieveBatch(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetBatchParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetBatch(ctx, params, options...)
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
		Title:          "x:users retrieve-batch",
		Transform:      transform,
	})
}

func handleXUsersRetrieveFollowers(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetFollowersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetFollowers(
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
		Title:          "x:users retrieve-followers",
		Transform:      transform,
	})
}

func handleXUsersRetrieveFollowersYouKnow(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetFollowersYouKnowParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetFollowersYouKnow(
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
		Title:          "x:users retrieve-followers-you-know",
		Transform:      transform,
	})
}

func handleXUsersRetrieveFollowing(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetFollowingParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetFollowing(
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
		Title:          "x:users retrieve-following",
		Transform:      transform,
	})
}

func handleXUsersRetrieveLikes(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetLikesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetLikes(
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
		Title:          "x:users retrieve-likes",
		Transform:      transform,
	})
}

func handleXUsersRetrieveMedia(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetMediaParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetMedia(
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
		Title:          "x:users retrieve-media",
		Transform:      transform,
	})
}

func handleXUsersRetrieveMentions(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetMentionsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetMentions(
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
		Title:          "x:users retrieve-mentions",
		Transform:      transform,
	})
}

func handleXUsersRetrieveReplies(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetRepliesParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetReplies(
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
		Title:          "x:users retrieve-replies",
		Transform:      transform,
	})
}

func handleXUsersRetrieveSearch(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetSearchParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetSearch(ctx, params, options...)
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
		Title:          "x:users retrieve-search",
		Transform:      transform,
	})
}

func handleXUsersRetrieveTweets(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetTweetsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetTweets(
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
		Title:          "x:users retrieve-tweets",
		Transform:      transform,
	})
}

func handleXUsersRetrieveVerifiedFollowers(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.XUserGetVerifiedFollowersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Users.GetVerifiedFollowers(
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
		Title:          "x:users retrieve-verified-followers",
		Transform:      transform,
	})
}
