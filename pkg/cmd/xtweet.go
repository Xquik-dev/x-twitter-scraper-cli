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
			Name:      "cursor",
			Usage:     "Pagination cursor for favoriters",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum user profiles requested from this page (20-200, default 200). The response can contain fewer profiles because the source returned fewer or remaining credits cover fewer results. Keep requesting next_cursor while has_next_page is true. The deprecated limit and count aliases remain accepted.\n",
			Default:   200,
			QueryPath: "pageSize",
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
			Name:      "min-faves",
			Usage:     "Minimum likes threshold.",
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
			Name:      "page-size",
			Usage:     "Maximum items requested from this page (1-100, default 20). The response can contain fewer items because the source returned fewer, filters removed items, or remaining credits cover fewer results. Keep requesting next_cursor while has_next_page is true, even when a page is empty. The deprecated limit and count aliases remain accepted.\n",
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
			Name:      "since-time",
			Usage:     "Unix timestamp - return quotes posted after this time",
			QueryPath: "sinceTime",
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
	},
	Action:          handleXTweetsGetQuotes,
	HideHelpCommand: true,
}

var xTweetsGetReplies = cli.Command{
	Name:    "get-replies",
	Usage:   "Returns visible replies. For an unfiltered first page, Xquik compares a terminal\npage with the post's reported reply count. If the page is visibly incomplete,\nthe endpoint returns 424 `replies_incomplete` instead of presenting partial\ncoverage as complete. Use tweet search with a `conversation_id:{id}` query as\nthe broader fallback.",
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
			Usage:     "Pagination cursor for tweet replies",
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
			Name:      "min-faves",
			Usage:     "Minimum likes threshold.",
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
			Name:      "page-size",
			Usage:     "Maximum items requested from this page (1-100, default 20). The response can contain fewer items because the source returned fewer, filters removed items, or remaining credits cover fewer results. Keep requesting next_cursor while has_next_page is true, even when a page is empty. The deprecated limit and count aliases remain accepted.\n",
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
			Name:      "since-time",
			Usage:     "Unix timestamp - return replies posted after this time",
			QueryPath: "sinceTime",
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
			Name:      "cursor",
			Usage:     "Pagination cursor for retweeters",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum user profiles requested from this page (20-200, default 200). The response can contain fewer profiles because the source returned fewer or remaining credits cover fewer results. Keep requesting next_cursor while has_next_page is true. The deprecated limit and count aliases remain accepted.\n",
			Default:   200,
			QueryPath: "pageSize",
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
			Name:      "cursor",
			Usage:     "Pagination cursor for thread tweets",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Maximum items requested from this page (1-100, default 20). The response can contain fewer items because the source returned fewer, filters removed items, or remaining credits cover fewer results. Keep requesting next_cursor while has_next_page is true, even when a page is empty. The deprecated limit and count aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
	},
	Action:          handleXTweetsGetThread,
	HideHelpCommand: true,
}

var xTweetsSearch = cli.Command{
	Name:    "search",
	Usage:   "Search tweets by query, Tweet ID, X status URL, or account date window",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "q",
			Usage:     "Search query (keywords,",
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
		&requestflag.Flag[string]{
			Name:      "bounding-box",
			Usage:     "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8.",
			QueryPath: "boundingBox",
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
			Usage:     "Pagination cursor from previous response",
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
			Name:      "limit",
			Usage:     "Max tweets to return (server paginates internally). Omit for single page (~20). This is an upper bound for paid authenticated calls: remaining credits can reduce the returned page size, and zero affordable results returns 402 insufficient_credits.\n",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "list-id",
			Usage:     "Search within a list ID.",
			QueryPath: "listId",
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
			Name:      "min-faves",
			Usage:     "Minimum likes threshold.",
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
		&requestflag.Flag[any]{
			Name:      "since-date",
			Usage:     "Start date in YYYY-MM-DD format.",
			QueryPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:      "since-time",
			Usage:     "ISO 8601 timestamp - only return tweets after this time",
			QueryPath: "sinceTime",
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
			Usage:     "ISO 8601 timestamp - only return tweets before this time",
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
