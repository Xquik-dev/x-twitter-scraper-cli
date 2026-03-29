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

var botPlatformLinksCreate = cli.Command{
	Name:    "create",
	Usage:   "Link a platform user to an Xquik account",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "platform",
			Usage:    `Allowed values: "telegram".`,
			Required: true,
			BodyPath: "platform",
		},
		&requestflag.Flag[string]{
			Name:     "platform-user-id",
			Required: true,
			BodyPath: "platformUserId",
		},
	},
	Action:          handleBotPlatformLinksCreate,
	HideHelpCommand: true,
}

var botPlatformLinksDelete = cli.Command{
	Name:    "delete",
	Usage:   "Unlink a platform user from an Xquik account",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "platform",
			Usage:    `Allowed values: "telegram".`,
			Required: true,
			BodyPath: "platform",
		},
		&requestflag.Flag[string]{
			Name:     "platform-user-id",
			Required: true,
			BodyPath: "platformUserId",
		},
	},
	Action:          handleBotPlatformLinksDelete,
	HideHelpCommand: true,
}

var botPlatformLinksLookup = cli.Command{
	Name:    "lookup",
	Usage:   "Look up an Xquik user by platform identity",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "platform",
			Required:  true,
			QueryPath: "platform",
		},
		&requestflag.Flag[string]{
			Name:      "platform-user-id",
			Required:  true,
			QueryPath: "platformUserId",
		},
	},
	Action:          handleBotPlatformLinksLookup,
	HideHelpCommand: true,
}

func handleBotPlatformLinksCreate(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.BotPlatformLinkNewParams{}

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
	_, err = client.Bot.PlatformLinks.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "bot:platform-links create", obj, format, transform)
}

func handleBotPlatformLinksDelete(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.BotPlatformLinkDeleteParams{}

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
	_, err = client.Bot.PlatformLinks.Delete(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "bot:platform-links delete", obj, format, transform)
}

func handleBotPlatformLinksLookup(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.BotPlatformLinkLookupParams{}

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
	_, err = client.Bot.PlatformLinks.Lookup(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "bot:platform-links lookup", obj, format, transform)
}
