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

var xMediaDownload = cli.Command{
	Name:    "download",
	Usage:   "Download images and videos from tweets",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:     "tweet-id",
			Usage:    "Numeric tweet ID alias for tweetInput",
			BodyPath: "tweetId",
		},
		&requestflag.String{
			Name:     "tweet-input",
			Usage:    "Tweet URL or ID (single tweet)",
			BodyPath: "tweetInput",
		},
		&requestflag.String{
			Name:     "tweet-url",
			Usage:    "Tweet URL alias for tweetInput",
			BodyPath: "tweetUrl",
		},
	},
	Action:          handleXMediaDownload,
	HideHelpCommand: true,
}

var xMediaUpload = cli.Command{
	Name:    "upload",
	Usage:   "Upload media",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:     "account",
			Usage:    "X account (@username or ID) uploading media from URL",
			Required: true,
			BodyPath: "account",
		},
		&requestflag.String{
			Name:     "url",
			Usage:    "HTTPS URL to download and upload as media",
			Required: true,
			BodyPath: "url",
		},
		&requestflag.String{
			Name:       "idempotency-key",
			Required:   true,
			HeaderPath: "Idempotency-Key",
		},
	},
	Action:          handleXMediaUpload,
	HideHelpCommand: true,
}

func handleXMediaDownload(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XMediaDownloadParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Media.Download(ctx, params, options...)
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
		Title:          "x:media download",
		Transform:      transform,
	})
}

func handleXMediaUpload(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, ApplicationJSON)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XMediaUploadParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Media.Upload(ctx, params, options...)
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
		Title:          "x:media upload",
		Transform:      transform,
	})
}
