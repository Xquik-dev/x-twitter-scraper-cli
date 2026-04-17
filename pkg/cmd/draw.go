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

var drawsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get draw details",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
	},
	Action:          handleDrawsRetrieve,
	HideHelpCommand: true,
}

var drawsList = cli.Command{
	Name:    "list",
	Usage:   "List draws",
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
	},
	Action:          handleDrawsList,
	HideHelpCommand: true,
}

var drawsExport = cli.Command{
	Name:    "export",
	Usage:   "Export draw data",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "id",
			Required: true,
		},
		&requestflag.Flag[string]{
			Name:      "format",
			Usage:     "Export output format",
			Default:   "csv",
			QueryPath: "format",
		},
		&requestflag.Flag[string]{
			Name:      "type",
			Usage:     "Export winners or all entries",
			Default:   "winners",
			QueryPath: "type",
		},
		&requestflag.Flag[string]{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "The file where the response contents will be stored. Use the value '-' to force output to stdout.",
		},
	},
	Action:          handleDrawsExport,
	HideHelpCommand: true,
}

var drawsRun = cli.Command{
	Name:    "run",
	Usage:   "Run giveaway draw",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "tweet-url",
			Required: true,
			BodyPath: "tweetUrl",
		},
		&requestflag.Flag[int64]{
			Name:     "backup-count",
			BodyPath: "backupCount",
		},
		&requestflag.Flag[int64]{
			Name:     "filter-account-age-days",
			BodyPath: "filterAccountAgeDays",
		},
		&requestflag.Flag[string]{
			Name:     "filter-language",
			BodyPath: "filterLanguage",
		},
		&requestflag.Flag[int64]{
			Name:     "filter-min-followers",
			BodyPath: "filterMinFollowers",
		},
		&requestflag.Flag[string]{
			Name:     "must-follow-username",
			BodyPath: "mustFollowUsername",
		},
		&requestflag.Flag[bool]{
			Name:     "must-retweet",
			BodyPath: "mustRetweet",
		},
		&requestflag.Flag[[]string]{
			Name:     "required-hashtag",
			BodyPath: "requiredHashtags",
		},
		&requestflag.Flag[[]string]{
			Name:     "required-keyword",
			BodyPath: "requiredKeywords",
		},
		&requestflag.Flag[[]string]{
			Name:     "required-mention",
			BodyPath: "requiredMentions",
		},
		&requestflag.Flag[bool]{
			Name:     "unique-authors-only",
			BodyPath: "uniqueAuthorsOnly",
		},
		&requestflag.Flag[int64]{
			Name:     "winner-count",
			Default:  1,
			BodyPath: "winnerCount",
		},
	},
	Action:          handleDrawsRun,
	HideHelpCommand: true,
}

func handleDrawsRetrieve(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.Draws.Get(ctx, cmd.Value("id").(string), options...)
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
		Title:          "draws retrieve",
		Transform:      transform,
	})
}

func handleDrawsList(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.DrawListParams{}

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
	_, err = client.Draws.List(ctx, params, options...)
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
		Title:          "draws list",
		Transform:      transform,
	})
}

func handleDrawsExport(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.DrawExportParams{}

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

	response, err := client.Draws.Export(
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

func handleDrawsRun(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.DrawRunParams{}

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
	_, err = client.Draws.Run(ctx, params, options...)
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
		Title:          "draws run",
		Transform:      transform,
	})
}
