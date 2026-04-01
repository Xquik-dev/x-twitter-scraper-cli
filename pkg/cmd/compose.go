// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var composeCreate = cli.Command{
	Name:    "create",
	Usage:   "Compose, refine, or score a tweet",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "step",
			Usage:    "Workflow step",
			Required: true,
			BodyPath: "step",
		},
		&requestflag.Flag[string]{
			Name:     "additional-context",
			Usage:    "Extra context or URLs (refine)",
			BodyPath: "additionalContext",
		},
		&requestflag.Flag[string]{
			Name:     "call-to-action",
			Usage:    "Desired call to action (refine)",
			BodyPath: "callToAction",
		},
		&requestflag.Flag[string]{
			Name:     "draft",
			Usage:    "Tweet draft text to evaluate (score)",
			BodyPath: "draft",
		},
		&requestflag.Flag[string]{
			Name:     "goal",
			Usage:    "Optimization goal",
			BodyPath: "goal",
		},
		&requestflag.Flag[bool]{
			Name:     "has-link",
			Usage:    "Whether a link is attached (score)",
			BodyPath: "hasLink",
		},
		&requestflag.Flag[bool]{
			Name:     "has-media",
			Usage:    "Whether media is attached (score)",
			BodyPath: "hasMedia",
		},
		&requestflag.Flag[string]{
			Name:     "media-type",
			Usage:    "Media type (refine)",
			BodyPath: "mediaType",
		},
		&requestflag.Flag[string]{
			Name:     "style-username",
			Usage:    "Cached style username for voice matching (compose)",
			BodyPath: "styleUsername",
		},
		&requestflag.Flag[string]{
			Name:     "tone",
			Usage:    "Desired tone (refine)",
			BodyPath: "tone",
		},
		&requestflag.Flag[string]{
			Name:     "topic",
			Usage:    "Tweet topic (compose, refine)",
			BodyPath: "topic",
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

	params := xtwitterscraper.ComposeNewParams{}

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
	_, err = client.Compose.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "compose create", obj, format, transform)
}
