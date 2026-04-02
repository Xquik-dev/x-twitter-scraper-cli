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

var stylesList = cli.Command{
	Name:            "list",
	Usage:           "List cached style profiles",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleStylesList,
	HideHelpCommand: true,
}

var stylesAnalyze = cli.Command{
	Name:    "analyze",
	Usage:   "Analyze writing style from recent tweets",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "username",
			Usage:    "X username to analyze",
			Required: true,
			BodyPath: "username",
		},
	},
	Action:          handleStylesAnalyze,
	HideHelpCommand: true,
}

var stylesCompare = cli.Command{
	Name:    "compare",
	Usage:   "Compare two style profiles",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "username1",
			Usage:     "First username to compare",
			Required:  true,
			QueryPath: "username1",
		},
		&requestflag.Flag[string]{
			Name:      "username2",
			Usage:     "Second username to compare",
			Required:  true,
			QueryPath: "username2",
		},
	},
	Action:          handleStylesCompare,
	HideHelpCommand: true,
}

func handleStylesList(ctx context.Context, cmd *cli.Command) error {
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
	_, err = client.Styles.List(ctx, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "styles list", obj, format, transform)
}

func handleStylesAnalyze(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.StyleAnalyzeParams{}

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
	_, err = client.Styles.Analyze(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "styles analyze", obj, format, transform)
}

func handleStylesCompare(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.StyleCompareParams{}

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
	_, err = client.Styles.Compare(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "styles compare", obj, format, transform)
}
