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

var xBookmarksList = cli.Command{
	Name:    "list",
	Usage:   "Get bookmarked tweets",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for bookmarks",
			QueryPath: "cursor",
		},
		&requestflag.String{
			Name:      "folder-id",
			Usage:     "Optional bookmark folder ID",
			QueryPath: "folderId",
		},
	},
	Action:          handleXBookmarksList,
	HideHelpCommand: true,
}

var xBookmarksRetrieveFolders = cli.Command{
	Name:            "retrieve-folders",
	Usage:           "Get bookmark folders",
	Suggest:         true,
	Flags:           []cli.Flag{},
	Action:          handleXBookmarksRetrieveFolders,
	HideHelpCommand: true,
}

func handleXBookmarksList(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XBookmarkListParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Bookmarks.List(ctx, params, options...)
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
		Title:          "x:bookmarks list",
		Transform:      transform,
	})
}

func handleXBookmarksRetrieveFolders(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Bookmarks.GetFolders(ctx, options...)
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
		Title:          "x:bookmarks retrieve-folders",
		Transform:      transform,
	})
}
