// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var xAccountsCreate = cli.Command{
	Name:    "create",
	Usage:   "Connect X account",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:     "email",
			Usage:    "Account email",
			Required: true,
			BodyPath: "email",
		},
		&requestflag.String{
			Name:     "password",
			Usage:    "Account password",
			Required: true,
			BodyPath: "password",
		},
		&requestflag.String{
			Name:     "totp-secret",
			Usage:    "Authenticator App TOTP secret required for durable login",
			Required: true,
			BodyPath: "totp_secret",
		},
		&requestflag.String{
			Name:     "username",
			Usage:    "X username",
			Required: true,
			BodyPath: "username",
		},
	},
	Action:          handleXAccountsCreate,
	HideHelpCommand: true,
}

var xAccountsRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "Get X account details",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
	},
	Action:          handleXAccountsRetrieve,
	HideHelpCommand: true,
}

var xAccountsList = cli.Command{
	Name:            "list",
	Usage:           "List connected X accounts",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleXAccountsList,
	HideHelpCommand: true,
}

var xAccountsDelete = cli.Command{
	Name:    "delete",
	Usage:   "Disconnect X account",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
	},
	Action:          handleXAccountsDelete,
	HideHelpCommand: true,
}

var xAccountsBulkRetry = cli.Command{
	Name:            "bulk-retry",
	Usage:           "Clears loginFailedAt and loginFailureReason for all accounts with transient or\nautomated failure reasons, making them eligible for retry on next use.",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleXAccountsBulkRetry,
	HideHelpCommand: true,
}

var xAccountsReauth = cli.Command{
	Name:    "reauth",
	Usage:   "Re-authenticate X account",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:     "password",
			Usage:    "Updated account password",
			Required: true,
			BodyPath: "password",
		},
		&requestflag.String{
			Name:     "email",
			Usage:    "Email for the X account (updates stored email)",
			BodyPath: "email",
		},
		&requestflag.String{
			Name:     "totp-secret",
			Usage:    "Replacement Authenticator App TOTP secret. Omit it to reuse the saved secret.",
			BodyPath: "totp_secret",
		},
	},
	Action:          handleXAccountsReauth,
	HideHelpCommand: true,
}

func handleXAccountsCreate(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XAccountNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Accounts.New(ctx, params, options...)
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
		Title:          "x:accounts create",
		Transform:      transform,
	})
}

func handleXAccountsRetrieve(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.X.Accounts.Get(ctx, cmd.String("id"), options...)
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
		Title:          "x:accounts retrieve",
		Transform:      transform,
	})
}

func handleXAccountsList(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Accounts.List(ctx, options...)
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
		Title:          "x:accounts list",
		Transform:      transform,
	})
}

func handleXAccountsDelete(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.X.Accounts.Delete(ctx, cmd.String("id"), options...)
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
		Title:          "x:accounts delete",
		Transform:      transform,
	})
}

func handleXAccountsBulkRetry(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Accounts.BulkRetry(ctx, options...)
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
		Title:          "x:accounts bulk-retry",
		Transform:      transform,
	})
}

func handleXAccountsReauth(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XAccountReauthParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Accounts.Reauth(
		ctx,
		cmd.String("id"),
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
		Title:          "x:accounts reauth",
		Transform:      transform,
	})
}
