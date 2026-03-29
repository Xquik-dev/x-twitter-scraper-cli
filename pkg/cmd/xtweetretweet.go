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

var xTweetsRetweetCreate = cli.Command{
	Name:    "create",
	Usage:   "Retweet",
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
	Action:          handleXTweetsRetweetCreate,
	HideHelpCommand: true,
}

var xTweetsRetweetDelete = cli.Command{
	Name:    "delete",
	Usage:   "Unretweet",
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
	Action:          handleXTweetsRetweetDelete,
	HideHelpCommand: true,
}

func handleXTweetsRetweetCreate(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("tweet-id") && len(unusedArgs) > 0 {
		cmd.Set("tweet-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XTweetRetweetNewParams{}

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
	_, err = client.X.Tweets.Retweet.New(
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
	return ShowJSON(os.Stdout, "x:tweets:retweet create", obj, format, transform)
}

func handleXTweetsRetweetDelete(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("tweet-id") && len(unusedArgs) > 0 {
		cmd.Set("tweet-id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XTweetRetweetDeleteParams{}

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
	_, err = client.X.Tweets.Retweet.Delete(
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
	return ShowJSON(os.Stdout, "x:tweets:retweet delete", obj, format, transform)
}
