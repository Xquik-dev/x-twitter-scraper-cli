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

var radarRetrieveTrendingTopics = cli.Command{
	Name:    "retrieve-trending-topics",
	Usage:   "Get trending topics from curated sources",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "after",
			Usage:     "Cursor for pagination (from prior response nextCursor).",
			QueryPath: "after",
		},
		&requestflag.Flag[string]{
			Name:      "category",
			Usage:     "Filter by category.",
			QueryPath: "category",
		},
		&requestflag.Flag[int64]{
			Name:      "hours",
			Usage:     "Lookback window in hours (1-168, default 24).",
			Default:   24,
			QueryPath: "hours",
		},
		&requestflag.Flag[int64]{
			Name:      "limit",
			Usage:     "Number of items to return (1-100, default 50).",
			Default:   50,
			QueryPath: "limit",
		},
		&requestflag.Flag[string]{
			Name:      "region",
			Usage:     "Region filter (us, global, etc.)",
			QueryPath: "region",
		},
		&requestflag.Flag[string]{
			Name:      "source",
			Usage:     "Source filter. One of: github, google_trends, hacker_news, polymarket, reddit, trustmrr, wikipedia",
			QueryPath: "source",
		},
	},
	Action:          handleRadarRetrieveTrendingTopics,
	HideHelpCommand: true,
}

func handleRadarRetrieveTrendingTopics(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.RadarGetTrendingTopicsParams{}

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
	_, err = client.Radar.GetTrendingTopics(ctx, params, options...)
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
		Title:          "radar retrieve-trending-topics",
		Transform:      transform,
	})
}
