// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/apiquery"
	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/requestflag"
	"github.com/stainless-sdks/x-twitter-scraper-go"
	"github.com/stainless-sdks/x-twitter-scraper-go/option"
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
			Name:     "text",
			Required: true,
			BodyPath: "text",
		},
		&requestflag.Flag[string]{
			Name:     "attachment-url",
			BodyPath: "attachment_url",
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
			Name:     "media-id",
			BodyPath: "media_ids",
		},
		&requestflag.Flag[string]{
			Name:     "reply-to-tweet-id",
			BodyPath: "reply_to_tweet_id",
		},
	},
	Action:          handleXTweetsCreate,
	HideHelpCommand: true,
}

var xTweetsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Look up tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "tweet-id",
			Required: true,
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
			Name:     "tweet-id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:     "account",
			Usage:    "X account (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
	},
	Action:          handleXTweetsDelete,
	HideHelpCommand: true,
}

var xTweetsGetFavoriters = cli.Command{
	Name:    "get-favoriters",
	Usage:   "Get users who liked a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor from previous response",
			QueryPath: "cursor",
		},
	},
	Action:          handleXTweetsGetFavoriters,
	HideHelpCommand: true,
}

var xTweetsGetQuotes = cli.Command{
	Name:    "get-quotes",
	Usage:   "Get quote tweets of a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor",
			QueryPath: "cursor",
		},
		&requestflag.Flag[bool]{
			Name:      "include-replies",
			Usage:     "Include replies (default false)",
			QueryPath: "includeReplies",
		},
		&requestflag.Flag[string]{
			Name:      "since-time",
			Usage:     "Unix timestamp - filter after",
			QueryPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:      "until-time",
			Usage:     "Unix timestamp - filter before",
			QueryPath: "untilTime",
		},
	},
	Action:          handleXTweetsGetQuotes,
	HideHelpCommand: true,
}

var xTweetsGetReplies = cli.Command{
	Name:    "get-replies",
	Usage:   "Get replies to a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "since-time",
			Usage:     "Unix timestamp - filter after",
			QueryPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:      "until-time",
			Usage:     "Unix timestamp - filter before",
			QueryPath: "untilTime",
		},
	},
	Action:          handleXTweetsGetReplies,
	HideHelpCommand: true,
}

var xTweetsGetRetweeters = cli.Command{
	Name:    "get-retweeters",
	Usage:   "Get users who retweeted a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor",
			QueryPath: "cursor",
		},
	},
	Action:          handleXTweetsGetRetweeters,
	HideHelpCommand: true,
}

var xTweetsGetThread = cli.Command{
	Name:    "get-thread",
	Usage:   "Get thread context for a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor",
			QueryPath: "cursor",
		},
	},
	Action:          handleXTweetsGetThread,
	HideHelpCommand: true,
}

var xTweetsSearch = cli.Command{
	Name:    "search",
	Usage:   "Search tweets",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "q",
			Usage:     "Search query (keywords,",
			Required:  true,
			QueryPath: "q",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor from previous response",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Deprecated — use cursor-based pagination instead",
			Default:   20,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "query-type",
			Usage:     "Sort order — Latest (chronological) or Top (engagement-ranked)",
			Default:   "Latest",
			QueryPath: "queryType",
		},
		&requestflag.Flag[string]{
			Name:      "since-time",
			Usage:     "ISO 8601 timestamp — only return tweets after this time",
			QueryPath: "sinceTime",
		},
		&requestflag.Flag[string]{
			Name:      "until-time",
			Usage:     "ISO 8601 timestamp — only return tweets before this time",
			QueryPath: "untilTime",
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

	params := xtwitterscraper.XTweetNewParams{}

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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets create", obj, format, transform)
}

func handleXTweetsRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("tweet-id") && len(unusedArgs) > 0 {
		cmd.Set("tweet-id", unusedArgs[0])
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
	_, err = client.X.Tweets.Get(ctx, cmd.Value("tweet-id").(string), options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets retrieve", obj, format, transform)
}

func handleXTweetsList(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XTweetListParams{}

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

	return client.X.Tweets.List(ctx, params, options...)
}

func handleXTweetsDelete(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("tweet-id") && len(unusedArgs) > 0 {
		cmd.Set("tweet-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XTweetDeleteParams{}

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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Tweets.Delete(
		ctx,
		cmd.Value("tweet-id").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets delete", obj, format, transform)
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

	params := xtwitterscraper.XTweetGetFavoritersParams{}

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
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets get-favoriters", obj, format, transform)
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

	params := xtwitterscraper.XTweetGetQuotesParams{}

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
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets get-quotes", obj, format, transform)
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

	params := xtwitterscraper.XTweetGetRepliesParams{}

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
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets get-replies", obj, format, transform)
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

	params := xtwitterscraper.XTweetGetRetweetersParams{}

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
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets get-retweeters", obj, format, transform)
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

	params := xtwitterscraper.XTweetGetThreadParams{}

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
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets get-thread", obj, format, transform)
}

func handleXTweetsSearch(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XTweetSearchParams{}

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
	_, err = client.X.Tweets.Search(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:tweets search", obj, format, transform)
}
