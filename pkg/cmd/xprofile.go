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

var xProfileUpdate = cli.Command{
	Name:    "update",
	Usage:   "Update X profile",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "account",
			Usage:    "X account (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.Flag[string]{
			Name:     "description",
			Usage:    "Bio description",
			BodyPath: "description",
		},
		&requestflag.Flag[string]{
			Name:     "location",
			BodyPath: "location",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Display name",
			BodyPath: "name",
		},
		&requestflag.Flag[string]{
			Name:     "url",
			Usage:    "Website URL",
			BodyPath: "url",
		},
	},
	Action:          handleXProfileUpdate,
	HideHelpCommand: true,
}

var xProfileUpdateAvatar = cli.Command{
	Name:    "update-avatar",
	Usage:   "Update profile avatar",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "account",
			Usage:    "X account (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.Flag[string]{
			Name:      "file",
			Usage:     "Avatar image (max 716KB)",
			Required:  true,
			BodyPath:  "file",
			FileInput: true,
		},
	},
	Action:          handleXProfileUpdateAvatar,
	HideHelpCommand: true,
}

var xProfileUpdateBanner = cli.Command{
	Name:    "update-banner",
	Usage:   "Update profile banner",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "account",
			Usage:    "X account (@username or account ID)",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.Flag[string]{
			Name:      "file",
			Usage:     "Banner image (max 2MB)",
			Required:  true,
			BodyPath:  "file",
			FileInput: true,
		},
	},
	Action:          handleXProfileUpdateBanner,
	HideHelpCommand: true,
}

func handleXProfileUpdate(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XProfileUpdateParams{}

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
	_, err = client.X.Profile.Update(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:profile update", obj, format, transform)
}

func handleXProfileUpdateAvatar(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XProfileUpdateAvatarParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		MultipartFormEncoded,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Profile.UpdateAvatar(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:profile update-avatar", obj, format, transform)
}

func handleXProfileUpdateBanner(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	params := xtwitterscraper.XProfileUpdateBannerParams{}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatComma,
		MultipartFormEncoded,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Profile.UpdateBanner(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(os.Stdout, "x:profile update-banner", obj, format, transform)
}
