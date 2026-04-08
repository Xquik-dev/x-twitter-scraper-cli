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

var xCommunitiesTweetsList = cli.Command{
	Name:    "list",
	Usage:   "Search tweets across all communities",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "q",
			Usage:     "Search query for cross-community tweets",
			Required:  true,
			QueryPath: "q",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for cross-community results",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "query-type",
			Usage:     "Sort order for cross-community results (Latest or Top)",
			QueryPath: "queryType",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleXCommunitiesTweetsList,
	HideHelpCommand: true,
}

var xCommunitiesTweetsListByCommunity = cli.Command{
	Name:    "list-by-community",
	Usage:   "Get community tweets",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for community tweets",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleXCommunitiesTweetsListByCommunity,
	HideHelpCommand: true,
}

func handleXCommunitiesTweetsList(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XCommunityTweetListParams{}

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

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.X.Communities.Tweets.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "x:communities:tweets list", obj, format, transform)
	} else {
		iter := client.X.Communities.Tweets.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(os.Stdout, "x:communities:tweets list", iter, format, transform, maxItems)
	}
}

func handleXCommunitiesTweetsListByCommunity(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XCommunityTweetListByCommunityParams{}

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

	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.X.Communities.Tweets.ListByCommunity(
			ctx,
			cmd.Value("id").(string),
			params,
			options...,
		)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(os.Stdout, "x:communities:tweets list-by-community", obj, format, transform)
	} else {
		iter := client.X.Communities.Tweets.ListByCommunityAutoPaging(
			ctx,
			cmd.Value("id").(string),
			params,
			options...,
		)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(os.Stdout, "x:communities:tweets list-by-community", iter, format, transform, maxItems)
	}
}
