// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var drawsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get draw details",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
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
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Previous nextCursor.",
			QueryPath: "cursor",
		},
		&requestflag.Int{
			Name:      "limit",
			Usage:     "Maximum number of items to return (1-100, default 50). For paid per-result endpoints, the returned count may be lower when remaining credits cannot cover the requested page. If zero paid results are affordable, the endpoint returns 402 insufficient_credits.\n",
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
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "format",
			Usage:     "Export output format. PDF entry exports include up to 10,000 rows. Other entry formats include up to 100,000 rows.",
			Required:  true,
			QueryPath: "format",
		},
		&requestflag.String{
			Name:      "type",
			Usage:     "Export winners or all entries",
			Default:   "winners",
			QueryPath: "type",
		},
		&requestflag.String{
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
	Usage:   "Runs a giveaway draw from a source tweet. The draw first checks the minimum\ncredits needed to inspect the source tweet and at least one candidate. Remaining\ncredits cap how many replies and retweeters can be inspected before filters and\nwinner selection run.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:     "tweet-url",
			Required: true,
			BodyPath: "tweetUrl",
		},
		&requestflag.Int{
			Name:     "backup-count",
			BodyPath: "backupCount",
		},
		&requestflag.Int{
			Name:     "filter-account-age-days",
			BodyPath: "filterAccountAgeDays",
		},
		&requestflag.String{
			Name:     "filter-language",
			BodyPath: "filterLanguage",
		},
		&requestflag.Int{
			Name:     "filter-min-followers",
			BodyPath: "filterMinFollowers",
		},
		&requestflag.String{
			Name:     "must-follow-username",
			BodyPath: "mustFollowUsername",
		},
		&requestflag.Bool{
			Name:     "must-retweet",
			BodyPath: "mustRetweet",
		},
		&requestflag.StringList{
			Name:     "required-hashtag",
			BodyPath: "requiredHashtags",
		},
		&requestflag.StringList{
			Name:     "required-keyword",
			BodyPath: "requiredKeywords",
		},
		&requestflag.StringList{
			Name:     "required-mention",
			BodyPath: "requiredMentions",
		},
		&requestflag.Bool{
			Name:     "unique-authors-only",
			BodyPath: "uniqueAuthorsOnly",
		},
		&requestflag.Int{
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Draws.Get(ctx, cmd.String("id"), options...)
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.DrawListParams{}

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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.DrawExportParams{}

	response, err := client.Draws.Export(
		ctx,
		cmd.String("id"),
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
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
	if err != nil {
		return err
	}

	params := xtwitterscraper.DrawRunParams{}

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
