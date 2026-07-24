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
			Usage:     "Cursor for keyset pagination from prior response next_cursor",
			QueryPath: "cursor",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of results to return (1-1000, default 100)",
			Default:   100,
			QueryPath: "limit",
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
			Usage:     "Cursor for keyset pagination from prior response next_cursor",
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

var extractionsEstimateCost = cli.Command{
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
			Usage:    "Raw advanced query string appended to the estimate (tweet_search_extractor)",
			BodyPath: "advancedQuery",
		},
		&requestflag.Flag[string]{
			Name:     "any-words",
			Usage:    "Alternative words or quoted phrases for estimated results. Separate with spaces, commas, or lines.",
			BodyPath: "anyWords",
		},
		&requestflag.Flag[string]{
			Name:     "bounding-box",
			Usage:    "Geo bounding box used for estimation, e.g. -74.1 40.6 -73.9 40.8 (tweet_search_extractor)",
			BodyPath: "boundingBox",
		},
		&requestflag.Flag[string]{
			Name:     "cashtags",
			Usage:    "Cashtags applied to the estimate, separated by spaces, commas, or lines.",
			BodyPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:     "conversation-id",
			Usage:    "Conversation ID filter used for estimation (tweet_search_extractor)",
			BodyPath: "conversationId",
		},
		&requestflag.Flag[string]{
			Name:     "exact-phrase",
			Usage:    "Exact phrase filter for search estimation",
			BodyPath: "exactPhrase",
		},
		&requestflag.Flag[string]{
			Name:     "exclude-words",
			Usage:    "Words or quoted phrases excluded from estimated results. Separate with spaces, commas, or lines.",
			BodyPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:     "from-user",
			Usage:    "Estimate only tweets from this author username (tweet_search_extractor)",
			BodyPath: "fromUser",
		},
		&requestflag.Flag[string]{
			Name:     "hashtags",
			Usage:    "Hashtags applied to the estimate, separated by spaces, commas, or lines.",
			BodyPath: "hashtags",
		},
		&requestflag.Flag[string]{
			Name:     "in-reply-to-tweet-id",
			Usage:    "Estimate only replies to this tweet ID (tweet_search_extractor)",
			BodyPath: "inReplyToTweetId",
		},
		&requestflag.Flag[string]{
			Name:     "language",
			Usage:    "Language code used for estimate filtering (tweet_search_extractor)",
			BodyPath: "language",
		},
		&requestflag.Flag[string]{
			Name:     "list-id",
			Usage:    "Estimate search results within this list ID (tweet_search_extractor)",
			BodyPath: "listId",
		},
		&requestflag.Flag[string]{
			Name:     "media-type",
			Usage:    "Media type used for estimate filtering (tweet_search_extractor)",
			BodyPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:     "mentioning",
			Usage:    "Estimate tweets mentioning this username (tweet_search_extractor)",
			BodyPath: "mentioning",
		},
		&requestflag.Flag[int64]{
			Name:     "min-faves",
			Usage:    "Minimum likes threshold for estimated results (tweet_search_extractor)",
			BodyPath: "minFaves",
		},
		&requestflag.Flag[int64]{
			Name:     "min-quotes",
			Usage:    "Minimum quote count threshold for estimated results (tweet_search_extractor)",
			BodyPath: "minQuotes",
		},
		&requestflag.Flag[int64]{
			Name:     "min-replies",
			Usage:    "Minimum replies threshold for estimated results (tweet_search_extractor)",
			BodyPath: "minReplies",
		},
		&requestflag.Flag[int64]{
			Name:     "min-retweets",
			Usage:    "Minimum retweets threshold for estimated results (tweet_search_extractor)",
			BodyPath: "minRetweets",
		},
		&requestflag.Flag[string]{
			Name:     "place",
			Usage:    "Estimate search results within this place ID (tweet_search_extractor)",
			BodyPath: "place",
		},
		&requestflag.Flag[string]{
			Name:     "place-country",
			Usage:    "Estimate search results within this country code (tweet_search_extractor)",
			BodyPath: "placeCountry",
		},
		&requestflag.Flag[string]{
			Name:     "point-radius",
			Usage:    "Geo point radius used for estimation, e.g. -73.99 40.73 25mi (tweet_search_extractor)",
			BodyPath: "pointRadius",
		},
		&requestflag.Flag[string]{
			Name:     "quotes",
			Usage:    "Quote mode used for estimation (tweet_search_extractor)",
			BodyPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:     "quotes-of-tweet-id",
			Usage:    "Estimate only quotes of this tweet ID (tweet_search_extractor)",
			BodyPath: "quotesOfTweetId",
		},
		&requestflag.Flag[string]{
			Name:     "replies",
			Usage:    "Reply mode used for estimation (tweet_search_extractor)",
			BodyPath: "replies",
		},
		&requestflag.Flag[int64]{
			Name:     "results-limit",
			Usage:    "Maximum number of results to estimate. When set, the estimate caps projected results to this value.",
			BodyPath: "resultsLimit",
		},
		&requestflag.Flag[string]{
			Name:     "retweets",
			Usage:    "Retweet mode used for estimation (tweet_search_extractor)",
			BodyPath: "retweets",
		},
		&requestflag.Flag[string]{
			Name:     "retweets-of-tweet-id",
			Usage:    "Estimate only retweets of this tweet ID (tweet_search_extractor)",
			BodyPath: "retweetsOfTweetId",
		},
		&requestflag.Flag[string]{
			Name:     "search-query",
			Usage:    "Query used to price tweet_search_extractor or community_search.",
			BodyPath: "searchQuery",
		},
		&requestflag.Flag[any]{
			Name:     "since-date",
			Usage:    "Estimate start date in YYYY-MM-DD format (tweet_search_extractor)",
			BodyPath: "sinceDate",
		},
		&requestflag.Flag[string]{
			Name:     "target-community-id",
			Usage:    "Community ID used to price community_post_extractor or community_search.",
			BodyPath: "targetCommunityId",
		},
		&requestflag.Flag[string]{
			Name:     "target-list-id",
			Usage:    "List ID used to price list_follower_explorer, list_member_extractor, or list_post_extractor.",
			BodyPath: "targetListId",
		},
		&requestflag.Flag[string]{
			Name:     "target-space-id",
			Usage:    "Space ID used to price space_explorer.",
			BodyPath: "targetSpaceId",
		},
		&requestflag.Flag[string]{
			Name:     "target-tweet-id",
			BodyPath: "targetTweetId",
		},
		&requestflag.Flag[string]{
			Name:     "target-username",
			BodyPath: "targetUsername",
		},
		&requestflag.Flag[string]{
			Name:     "to-user",
			Usage:    "Estimate replies sent to this username (tweet_search_extractor)",
			BodyPath: "toUser",
		},
		&requestflag.Flag[any]{
			Name:     "until-date",
			Usage:    "Estimate end date in YYYY-MM-DD format (tweet_search_extractor)",
			BodyPath: "untilDate",
		},
		&requestflag.Flag[string]{
			Name:     "url",
			Usage:    "URL substring or domain filter used for estimation (tweet_search_extractor)",
			BodyPath: "url",
		},
		&requestflag.Flag[bool]{
			Name:     "verified-only",
			Usage:    "Estimate only verified authors (tweet_search_extractor)",
			BodyPath: "verifiedOnly",
		},
	},
	Action:          handleExtractionsEstimateCost,
	HideHelpCommand: true,
}

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
		&requestflag.Flag[string]{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The file where the response contents will be stored. Use the value '-' to force output to stdout.",
		},
	},
	Action:          handleExtractionsExportResults,
	HideHelpCommand: true,
}

var extractionsRun = cli.Command{
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
			Name:     "bounding-box",
			Usage:    "Geo bounding box, e.g. -74.1 40.6 -73.9 40.8 (tweet_search_extractor)",
			BodyPath: "boundingBox",
		},
		&requestflag.Flag[string]{
			Name:     "cashtags",
			Usage:    "Cashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "cashtags",
		},
		&requestflag.Flag[string]{
			Name:     "conversation-id",
			Usage:    "Conversation ID filter (tweet_search_extractor)",
			BodyPath: "conversationId",
		},
		&requestflag.Flag[string]{
			Name:     "exact-phrase",
			Usage:    "Exact phrase to match (tweet_search_extractor)",
			BodyPath: "exactPhrase",
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
			Name:     "hashtags",
			Usage:    "Hashtags separated by spaces, commas, or lines. (tweet_search_extractor)",
			BodyPath: "hashtags",
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
			Name:     "min-faves",
			Usage:    "Minimum likes threshold (tweet_search_extractor)",
			BodyPath: "minFaves",
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
			Name:     "quotes",
			Usage:    "Quote mode (tweet_search_extractor)",
			BodyPath: "quotes",
		},
		&requestflag.Flag[string]{
			Name:     "quotes-of-tweet-id",
			Usage:    "Only quotes of this tweet ID (tweet_search_extractor)",
			BodyPath: "quotesOfTweetId",
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
			Name:     "target-community-id",
			Usage:    "Required for community_post_extractor & community_search.",
			BodyPath: "targetCommunityId",
		},
		&requestflag.Flag[string]{
			Name:     "target-list-id",
			Usage:    "Required for list_follower_explorer, list_member_extractor & list_post_extractor.",
			BodyPath: "targetListId",
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
		&requestflag.Flag[string]{
			Name:     "target-username",
			BodyPath: "targetUsername",
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
		&requestflag.Flag[string]{
			Name:     "url",
			Usage:    "URL substring or domain filter (tweet_search_extractor)",
			BodyPath: "url",
		},
		&requestflag.Flag[bool]{
			Name:     "verified-only",
			Usage:    "Only verified authors (tweet_search_extractor)",
			BodyPath: "verifiedOnly",
		},
	},
	Action:          handleExtractionsRun,
	HideHelpCommand: true,
}

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
