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

var extractionsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get extraction results",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "after",
			Usage:     "Cursor for keyset pagination",
			QueryPath: "after",
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
			Name:      "after",
			Usage:     "Cursor for keyset pagination",
			QueryPath: "after",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Maximum number of items to return (1-100, default 50)",
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
			Name:     "exact-phrase",
			Usage:    "Exact phrase filter for search estimation",
			BodyPath: "exactPhrase",
		},
		&requestflag.Flag[string]{
			Name:     "exclude-words",
			Usage:    "Words excluded from estimated search results",
			BodyPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:     "search-query",
			BodyPath: "searchQuery",
		},
		&requestflag.Flag[string]{
			Name:     "target-community-id",
			BodyPath: "targetCommunityId",
		},
		&requestflag.Flag[string]{
			Name:     "target-list-id",
			BodyPath: "targetListId",
		},
		&requestflag.Flag[string]{
			Name:     "target-space-id",
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
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "format",
			Usage:     "Export file format",
			Default:   "csv",
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
			Name:     "exact-phrase",
			Usage:    "Exact phrase to match (tweet_search_extractor)",
			BodyPath: "exactPhrase",
		},
		&requestflag.Flag[string]{
			Name:     "exclude-words",
			Usage:    "Words to exclude from results (tweet_search_extractor)",
			BodyPath: "excludeWords",
		},
		&requestflag.Flag[string]{
			Name:     "search-query",
			BodyPath: "searchQuery",
		},
		&requestflag.Flag[string]{
			Name:     "target-community-id",
			BodyPath: "targetCommunityId",
		},
		&requestflag.Flag[string]{
			Name:     "target-list-id",
			BodyPath: "targetListId",
		},
		&requestflag.Flag[string]{
			Name:     "target-space-id",
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

	params := xtwitterscraper.ExtractionGetParams{}

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

	params := xtwitterscraper.ExtractionListParams{}

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

	params := xtwitterscraper.ExtractionEstimateCostParams{}

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

	params := xtwitterscraper.ExtractionExportResultsParams{}

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

	params := xtwitterscraper.ExtractionRunParams{}

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
		Title:          "extractions run",
		Transform:      transform,
	})
}
