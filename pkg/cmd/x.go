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

var xGetArticle = cli.Command{
	Name:    "get-article",
	Usage:   "Retrieve the full content of an X Article (long-form post) by tweet ID.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "tweet-id",
			Required: true,
		},
	},
	Action:          handleXGetArticle,
	HideHelpCommand: true,
}

var xGetHomeTimeline = cli.Command{
	Name:    "get-home-timeline",
	Usage:   "Get home timeline",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for timeline",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "seen-tweet-ids",
			Usage:     "Comma-separated tweet IDs to exclude from results",
			QueryPath: "seenTweetIds",
		},
	},
	Action:          handleXGetHomeTimeline,
	HideHelpCommand: true,
}

var xGetNotifications = cli.Command{
	Name:    "get-notifications",
	Usage:   "Get notifications",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for notifications",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "type",
			Usage:     "Notification type filter",
			Default:   "All",
			QueryPath: "type",
		},
	},
	Action:          handleXGetNotifications,
	HideHelpCommand: true,
}

var xGetTrends = cli.Command{
	Name:    "get-trends",
	Usage:   "Get trending hashtags & topics from X by region",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:      "count",
			Usage:     "Number of trending topics to return (1-50, default 30)",
			Default:   30,
			QueryPath: "count",
		},
		&requestflag.Flag[int64]{
			Name:      "woeid",
			Usage:     "Region WOEID (1=Worldwide, 23424977=US, 23424975=UK, 23424969=Turkey)",
			Default:   1,
			QueryPath: "woeid",
		},
	},
	Action:          handleXGetTrends,
	HideHelpCommand: true,
}

func handleXGetArticle(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.X.GetArticle(ctx, cmd.Value("tweet-id").(string), options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, os.Stderr, "x get-article", obj, format, explicitFormat, transform)
}

func handleXGetHomeTimeline(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XGetHomeTimelineParams{}

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
	_, err = client.X.GetHomeTimeline(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, os.Stderr, "x get-home-timeline", obj, format, explicitFormat, transform)
}

func handleXGetNotifications(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XGetNotificationsParams{}

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
	_, err = client.X.GetNotifications(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, os.Stderr, "x get-notifications", obj, format, explicitFormat, transform)
}

func handleXGetTrends(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XGetTrendsParams{}

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
	_, err = client.X.GetTrends(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, os.Stderr, "x get-trends", obj, format, explicitFormat, transform)
}
