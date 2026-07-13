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

var xDmRetrieveHistory = cli.Command{
	Name:    "retrieve-history",
	Usage:   "Get DM conversation history",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "user-id",
			Required:  true,
			PathParam: "userId",
		},
		&requestflag.Flag[string]{
			Name:      "account",
			Usage:     "X handle (without the `@` prefix) of the connected X account used to read the conversation. The account must be a participant in the conversation.\n",
			Required:  true,
			QueryPath: "account",
		},
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Pagination cursor for DM history",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "max-id",
			Usage:     "Legacy pagination cursor (backward compat)",
			QueryPath: "maxId",
		},
	},
	Action:          handleXDmRetrieveHistory,
	HideHelpCommand: true,
}

var xDmSend = cli.Command{
	Name:    "send",
	Usage:   "Send direct message",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "user-id",
			Required:  true,
			PathParam: "userId",
		},
		&requestflag.Flag[string]{
			Name:     "account",
			Usage:    "X account (@username or ID) sending the DM",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.Flag[string]{
			Name:     "text",
			Required: true,
			BodyPath: "text",
		},
		&requestflag.Flag[[]string]{
			Name:     "media-id",
			Usage:    "Optional array containing exactly 1 uploaded media ID.",
			BodyPath: "media_ids",
		},
	},
	Action:          handleXDmSend,
	HideHelpCommand: true,
}

func handleXDmRetrieveHistory(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("user-id") && len(unusedArgs) > 0 {
		cmd.Set("user-id", unusedArgs[0])
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

	params := xtwitterscraper.XDmGetHistoryParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Dm.GetHistory(
		ctx,
		cmd.Value("user-id").(string),
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
		Title:          "x:dm retrieve-history",
		Transform:      transform,
	})
}

func handleXDmSend(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("user-id") && len(unusedArgs) > 0 {
		cmd.Set("user-id", unusedArgs[0])
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

	params := xtwitterscraper.XDmSendParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Dm.Send(
		ctx,
		cmd.Value("user-id").(string),
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
		Title:          "x:dm send",
		Transform:      transform,
	})
}
