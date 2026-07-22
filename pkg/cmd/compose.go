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

var composeCreate = cli.Command{
	Name:    "create",
	Usage:   "Run one step of Xquik's three-step writing workflow. Compose returns questions\nand editorial rules. Refine returns goal-specific guidance. Score applies\ndeterministic text checks. It does not predict reach or expose X ranking\nweights.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "step",
			Usage:    `Allowed values: "compose".`,
			Default:  "compose",
			Const:    true,
			BodyPath: "step",
		},
		&requestflag.Flag[string]{
			Name:     "topic",
			Usage:    "Subject for the post.",
			BodyPath: "topic",
		},
		&requestflag.Flag[string]{
			Name:     "goal",
			Usage:    "Editorial goal used to order the rules and questions.",
			Default:  "engagement",
			BodyPath: "goal",
		},
		&requestflag.Flag[string]{
			Name:     "style-username",
			Usage:    "Username from a style analysis saved to this account.",
			BodyPath: "styleUsername",
		},
		&requestflag.Flag[string]{
			Name:     "tone",
			Usage:    "Requested writing tone.",
			BodyPath: "tone",
		},
		&requestflag.Flag[string]{
			Name:     "additional-context",
			Usage:    "Audience, constraints, sources, or other writing context.",
			BodyPath: "additionalContext",
		},
		&requestflag.Flag[string]{
			Name:     "call-to-action",
			Usage:    "Specific action the draft should request.",
			BodyPath: "callToAction",
		},
		&requestflag.Flag[string]{
			Name:     "media-type",
			Usage:    "Planned media type.",
			BodyPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:     "draft",
			Usage:    "Full post text for deterministic editorial checks.",
			BodyPath: "draft",
		},
		&requestflag.Flag[bool]{
			Name:     "has-link",
			Usage:    "True when a separate link card is attached.",
			Default:  false,
			BodyPath: "hasLink",
		},
		&requestflag.Flag[bool]{
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
