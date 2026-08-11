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

var creditsRedirectTopupCheckout = cli.Command{
	Name:    "redirect-topup-checkout",
	Usage:   "Redirect to an active top-up payment page",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "session-id",
			Usage:     "Billing session ID returned by the top-up billing flow.",
			Required:  true,
			QueryPath: "session_id",
		},
	},
	Action:          handleCreditsRedirectTopupCheckout,
	HideHelpCommand: true,
}

var creditsRetrieveBalance = cli.Command{
	Name:            "retrieve-balance",
	Usage:           "Get credits balance",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleCreditsRetrieveBalance,
	HideHelpCommand: true,
}

var creditsRetrieveTopupStatus = cli.Command{
	Name:    "retrieve-topup-status",
	Usage:   "Get top-up billing status",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "session-id",
			Usage:     "Top-up session ID to inspect.",
			Required:  true,
			QueryPath: "session_id",
		},
	},
	Action:          handleCreditsRetrieveTopupStatus,
	HideHelpCommand: true,
}

var creditsTopupBalance = cli.Command{
	Name:    "topup-balance",
	Usage:   "Create a hosted checkout only after the user confirms. The request never\ncompletes payment or adds credits.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[int64]{
			Name:     "dollars",
			Usage:    "Amount to top up in US dollars. Minimum 10.",
			Required: true,
			BodyPath: "dollars",
		},
		&requestflag.Flag[string]{
			Name:     "locale",
			Usage:    "Optional checkout locale. Defaults to en.",
			BodyPath: "locale",
		},
	},
	Action:          handleCreditsTopupBalance,
	HideHelpCommand: true,
}

func handleCreditsRedirectTopupCheckout(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.CreditRedirectTopupCheckoutParams{}

	return client.Credits.RedirectTopupCheckout(ctx, params, options...)
}

func handleCreditsRetrieveBalance(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.Credits.GetBalance(ctx, options...)
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
		Title:          "credits retrieve-balance",
		Transform:      transform,
	})
}

func handleCreditsRetrieveTopupStatus(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.CreditGetTopupStatusParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Credits.GetTopupStatus(ctx, params, options...)
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
		Title:          "credits retrieve-topup-status",
		Transform:      transform,
	})
}

func handleCreditsTopupBalance(ctx context.Context, cmd *cli.Command) error {
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

	params := xtwitterscraper.CreditTopupBalanceParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Credits.TopupBalance(ctx, params, options...)
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
		Title:          "credits topup-balance",
		Transform:      transform,
	})
}
