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

var xAccountConnectionChallengesSubmit = cli.Command{
	Name:    "submit",
	Usage:   "Submit X account email verification code",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.Flag[string]{
			Name:     "email-code",
			Usage:    "Code sent to the account email.",
			Required: true,
			BodyPath: "email_code",
		},
	},
	Action:          handleXAccountConnectionChallengesSubmit,
	HideHelpCommand: true,
}

func handleXAccountConnectionChallengesSubmit(ctx context.Context, cmd *cli.Command) error {
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
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XAccountConnectionChallengeSubmitParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.AccountConnectionChallenges.Submit(
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
		Title:          "x:account-connection-challenges submit",
		Transform:      transform,
	})
}
