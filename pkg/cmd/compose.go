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

var composeCreate = cli.Command{
	Name:    "create",
	Usage:   "Run one step of Xquik's three-step writing workflow. Compose returns questions,\neditorial rules, and source-specific Radar recommendations. Refine returns\ngoal-specific guidance. Score applies deterministic text checks. It does not\npredict reach or expose X ranking weights.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:     "step",
			Usage:    `Use "compose".`,
			Default:  "compose",
			Const:    true,
			BodyPath: "step",
		},
		&requestflag.String{
			Name:     "topic",
			Usage:    "Subject for the post.",
			BodyPath: "topic",
		},
		&requestflag.String{
			Name:     "goal",
			Usage:    "Editorial goal used to order the rules and questions.",
			Default:  "engagement",
			BodyPath: "goal",
		},
		&requestflag.String{
			Name:     "style-username",
			Usage:    "Username from a style analysis saved to this account.",
			BodyPath: "styleUsername",
		},
		&requestflag.String{
			Name:     "tone",
			Usage:    "Requested writing tone.",
			BodyPath: "tone",
		},
		&requestflag.String{
			Name:     "additional-context",
			Usage:    "Audience, constraints, sources, or other writing context.",
			BodyPath: "additionalContext",
		},
		&requestflag.String{
			Name:     "call-to-action",
			Usage:    "Specific action the draft should request.",
			BodyPath: "callToAction",
		},
		&requestflag.String{
			Name:     "media-type",
			Usage:    "Planned media type.",
			BodyPath: "mediaType",
		},
		&requestflag.String{
			Name:     "draft",
			Usage:    "Full post text for deterministic editorial checks.",
			BodyPath: "draft",
		},
		&requestflag.Bool{
			Name:     "has-link",
			Usage:    "True when a separate link card is attached.",
			Default:  false,
			BodyPath: "hasLink",
		},
		&requestflag.Bool{
			Name:     "has-media",
			Usage:    "Accepted for backward compatibility. Text checks ignore this field.\n",
			BodyPath: "hasMedia",
		},
	},
	Action:          handleComposeCreate,
	HideHelpCommand: true,
}

func handleComposeCreate(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
	if err != nil {
		return err
	}

	params := xtwitterscraper.ComposeNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Compose.New(ctx, params, options...)
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
		Title:          "compose create",
		Transform:      transform,
	})
}
