// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

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

var guestWalletsCreate = cli.Command{
	Name:    "create",
	Usage:   "Create a one-use hosted checkout after the user confirms $10-$250 USD. The\nrequest creates no charge. It returns a paid-read API key without an Xquik\naccount. Replays return the same key.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:     "amount-minor",
			Usage:    "USD cents accepted for this checkout.",
			Required: true,
			BodyPath: "amount_minor",
		},
		&requestflag.Flag[string]{
			Name:     "currency",
			Usage:    `Allowed values: "usd".`,
			Default:  "usd",
			Const:    true,
			BodyPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:       "idempotency-key",
			Required:   true,
			HeaderPath: "Idempotency-Key",
		},
	},
	Action:          handleGuestWalletsCreate,
	HideHelpCommand: true,
}

var guestWalletsRetrieveStatus = cli.Command{
	Name:            "retrieve-status",
	Usage:           "Poll after payment. Use usable to decide whether paid reads can run. An active\nwallet can remain usable while a top-up is pending. A new wallet becomes usable\nonly after payment is verified. Send the guest key as Authorization: Bearer.",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleGuestWalletsRetrieveStatus,
	HideHelpCommand: true,
}

var guestWalletsTopup = cli.Command{
	Name:    "topup",
	Usage:   "Create a one-use hosted checkout after the user confirms a $10-$250 USD amount\nfor an existing paid-read guest key. The key remains the same. This request\ncreates no charge and never redirects through Xquik.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:     "amount-minor",
			Usage:    "USD cents accepted for this checkout.",
			Required: true,
			BodyPath: "amount_minor",
		},
		&requestflag.Flag[string]{
			Name:     "currency",
			Usage:    `Allowed values: "usd".`,
			Default:  "usd",
			Const:    true,
			BodyPath: "currency",
		},
		&requestflag.Flag[string]{
			Name:       "idempotency-key",
			Required:   true,
			HeaderPath: "Idempotency-Key",
		},
	},
	Action:          handleGuestWalletsTopup,
	HideHelpCommand: true,
}

func handleGuestWalletsCreate(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.GuestWalletNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.GuestWallets.New(ctx, params, options...)
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
		Title:          "guest-wallets create",
		Transform:      transform,
	})
}

func handleGuestWalletsRetrieveStatus(ctx context.Context, cmd *cli.Command) error {
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

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.GuestWallets.GetStatus(ctx, options...)
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
		Title:          "guest-wallets retrieve-status",
		Transform:      transform,
	})
}

func handleGuestWalletsTopup(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.GuestWalletTopupParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.GuestWallets.Topup(ctx, params, options...)
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
		Title:          "guest-wallets topup",
		Transform:      transform,
	})
}
