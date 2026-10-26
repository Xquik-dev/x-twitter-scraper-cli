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

var xListsRetrieveFollowers = cli.Command{
	Name:    "retrieve-followers",
	Usage:   "List followers of an X List",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for list followers",
			QueryPath: "cursor",
		},
		&requestflag.Bool{
			Name:      "has-location",
			Usage:     "Only return profiles with a location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Bool{
			Name:      "has-website",
			Usage:     "Only return profiles with a website.",
			QueryPath: "hasWebsite",
		},
		&requestflag.String{
			Name:      "location-contains",
			Usage:     "Match a location substring, ignoring case.",
			QueryPath: "locationContains",
		},
		&requestflag.Int{
			Name:      "max-followers",
			Usage:     "Maximum follower count. Missing counts pass this maximum.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Int{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Int{
			Name:      "max-statuses",
			Usage:     "Maximum post count. maxPosts is also accepted.",
			QueryPath: "maxStatuses",
		},
		&requestflag.Int{
			Name:      "min-account-age-days",
			Usage:     "Minimum account age in whole days.",
			QueryPath: "minAccountAgeDays",
		},
		&requestflag.Int{
			Name:      "min-followers",
			Usage:     "Minimum follower count. Filtering happens before billing.",
			QueryPath: "minFollowers",
		},
		&requestflag.Int{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Int{
			Name:      "min-statuses",
			Usage:     "Minimum post count. minPosts is also accepted.",
			QueryPath: "minStatuses",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Maximum user profiles requested from this page (20-200, default 200). Source, filters, or credits can return fewer profiles. Keep requesting next_cursor while has_next_page is true. Deprecated aliases remain accepted.\n",
			Default:   200,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "username-contains",
			Usage:     "Match a username substring, ignoring case.",
			QueryPath: "usernameContains",
		},
		&requestflag.Bool{
			Name:      "verified-only",
			Usage:     "Only return verified profiles.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.String{
			Name:      "verified-type",
			Usage:     "Match the verification type exactly, ignoring case.",
			QueryPath: "verifiedType",
		},
	},
	Action:          handleXListsRetrieveFollowers,
	HideHelpCommand: true,
}

var xListsRetrieveMembers = cli.Command{
	Name:    "retrieve-members",
	Usage:   "List members of an X List",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "bio-contains",
			Usage:     "Match any comma-separated or line-separated bio term, ignoring case.\n",
			QueryPath: "bioContains",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for list members",
			QueryPath: "cursor",
		},
		&requestflag.Bool{
			Name:      "has-location",
			Usage:     "Only return profiles with a location.",
			QueryPath: "hasLocation",
		},
		&requestflag.Bool{
			Name:      "has-website",
			Usage:     "Only return profiles with a website.",
			QueryPath: "hasWebsite",
		},
		&requestflag.String{
			Name:      "location-contains",
			Usage:     "Match a location substring, ignoring case.",
			QueryPath: "locationContains",
		},
		&requestflag.Int{
			Name:      "max-followers",
			Usage:     "Maximum follower count. Missing counts pass this maximum.",
			QueryPath: "maxFollowers",
		},
		&requestflag.Int{
			Name:      "max-following",
			Usage:     "Maximum following count.",
			QueryPath: "maxFollowing",
		},
		&requestflag.Int{
			Name:      "max-statuses",
			Usage:     "Maximum post count. maxPosts is also accepted.",
			QueryPath: "maxStatuses",
		},
		&requestflag.Int{
			Name:      "min-account-age-days",
			Usage:     "Minimum account age in whole days.",
			QueryPath: "minAccountAgeDays",
		},
		&requestflag.Int{
			Name:      "min-followers",
			Usage:     "Minimum follower count. Filtering happens before billing.",
			QueryPath: "minFollowers",
		},
		&requestflag.Int{
			Name:      "min-following",
			Usage:     "Minimum following count.",
			QueryPath: "minFollowing",
		},
		&requestflag.Int{
			Name:      "min-statuses",
			Usage:     "Minimum post count. minPosts is also accepted.",
			QueryPath: "minStatuses",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Members per page (20-200, default 20)",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "username-contains",
			Usage:     "Match a username substring, ignoring case.",
			QueryPath: "usernameContains",
		},
		&requestflag.Bool{
			Name:      "verified-only",
			Usage:     "Only return verified profiles.",
			QueryPath: "verifiedOnly",
		},
		&requestflag.String{
			Name:      "verified-type",
			Usage:     "Match the verification type exactly, ignoring case.",
			QueryPath: "verifiedType",
		},
	},
	Action:          handleXListsRetrieveMembers,
	HideHelpCommand: true,
}

var xListsRetrieveTweets = cli.Command{
	Name:    "retrieve-tweets",
	Usage:   "List tweets from an X List",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.String{
			Name:      "id",
			Required:  true,
			PathParam: "id",
		},
		&requestflag.String{
			Name:      "cursor",
			Usage:     "Pagination cursor for list tweets",
			QueryPath: "cursor",
		},
		&requestflag.Bool{
			Name:      "include-replies",
			Usage:     "Include replies (default false)",
			QueryPath: "includeReplies",
		},
		&requestflag.Int{
			Name:      "page-size",
			Usage:     "Maximum page items (1-100, default 20). Source, filters, or credits can reduce results. Continue while has_next_page is true. Deprecated limit and count aliases remain accepted.\n",
			Default:   20,
			QueryPath: "pageSize",
		},
		&requestflag.String{
			Name:      "since-time",
			Usage:     "Unix timestamp - filter after",
			QueryPath: "sinceTime",
		},
		&requestflag.String{
			Name:      "until-time",
			Usage:     "Unix timestamp - filter before",
			QueryPath: "untilTime",
		},
	},
	Action:          handleXListsRetrieveTweets,
	HideHelpCommand: true,
}

func handleXListsRetrieveFollowers(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XListGetFollowersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Lists.GetFollowers(
		ctx,
		cmd.String("id"),
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
		Title:          "x:lists retrieve-followers",
		Transform:      transform,
	})
}

func handleXListsRetrieveMembers(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XListGetMembersParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Lists.GetMembers(
		ctx,
		cmd.String("id"),
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
		Title:          "x:lists retrieve-members",
		Transform:      transform,
	})
}

func handleXListsRetrieveTweets(ctx context.Context, cmd *cli.Command) error {
	client := xtwitterscraper.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("id") && len(unusedArgs) > 0 {
		cmd.Set("id", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return unexpectedArgumentsError(unusedArgs)
	}

	options, err := flagOptions(cmd, EmptyBody)
	if err != nil {
		return err
	}

	params := xtwitterscraper.XListGetTweetsParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.X.Lists.GetTweets(
		ctx,
		cmd.String("id"),
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
		Title:          "x:lists retrieve-tweets",
		Transform:      transform,
	})
}
