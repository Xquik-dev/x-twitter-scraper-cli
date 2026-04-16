// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/apiquery"
	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/requestflag"
	"github.com/stainless-sdks/x-twitter-scraper-go"
	"github.com/stainless-sdks/x-twitter-scraper-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var xDmRetrieveHistory = cli.Command{
	Name:    "retrieve-history",
	Usage:   "Get DM conversation history",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "user-id",
			Required: true,
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
			Name:     "user-id",
			Required: true,
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
			BodyPath: "media_ids",
		},
		&requestflag.Flag[string]{
			Name:     "reply-to-message-id",
			BodyPath: "reply_to_message_id",
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

	params := xtwitterscraper.XDmGetHistoryParams{}

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

	params := xtwitterscraper.XDmSendParams{}

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
		Title:          "x:dm send",
		Transform:      transform,
	})
}
