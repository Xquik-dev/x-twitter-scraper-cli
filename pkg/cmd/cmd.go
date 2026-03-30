// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/autocomplete"
	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/requestflag"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var (
	Command            *cli.Command
	CommandErrorBuffer bytes.Buffer
)

func init() {
	Command = &cli.Command{
		Name:      "x-twitter-scraper",
		Usage:     "CLI for the x-twitter-scraper API",
		Suggest:   true,
		Version:   Version,
		ErrWriter: &CommandErrorBuffer,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug logging",
			},
			&cli.StringFlag{
				Name:        "base-url",
				DefaultText: "url",
				Usage:       "Override the base URL for API requests",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "The format for displaying response data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "format-error",
				Usage: "The format for displaying error data (one of: " + strings.Join(OutputFormats, ", ") + ")",
				Value: "auto",
				Validator: func(format string) error {
					if !slices.Contains(OutputFormats, strings.ToLower(format)) {
						return fmt.Errorf("format must be one of: %s", strings.Join(OutputFormats, ", "))
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "transform",
				Usage: "The GJSON transformation for data output.",
			},
			&cli.StringFlag{
				Name:  "transform-error",
				Usage: "The GJSON transformation for errors.",
			},
			&requestflag.Flag[string]{
				Name:    "api-key",
				Sources: cli.EnvVars("X_TWITTER_SCRAPER_API_KEY"),
			},
			&requestflag.Flag[string]{
				Name:    "bearer-token",
				Usage:   "OAuth 2.1 access token",
				Sources: cli.EnvVars("X_TWITTER_SCRAPER_BEARER_TOKEN"),
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "account",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&accountRetrieve,
					&accountSetXUsername,
					&accountUpdateLocale,
				},
			},
			{
				Name:     "api-keys",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&apiKeysCreate,
					&apiKeysList,
					&apiKeysRevoke,
				},
			},
			{
				Name:     "subscribe",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&subscribeCreate,
				},
			},
			{
				Name:     "compose",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&composeCreate,
				},
			},
			{
				Name:     "drafts",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&draftsCreate,
					&draftsRetrieve,
					&draftsList,
					&draftsDelete,
				},
			},
			{
				Name:     "styles",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&stylesRetrieve,
					&stylesUpdate,
					&stylesList,
					&stylesDelete,
					&stylesAnalyze,
					&stylesCompare,
					&stylesGetPerformance,
				},
			},
			{
				Name:     "radar",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&radarRetrieveTrendingTopics,
				},
			},
			{
				Name:     "monitors",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&monitorsCreate,
					&monitorsRetrieve,
					&monitorsUpdate,
					&monitorsList,
					&monitorsDeactivate,
				},
			},
			{
				Name:     "events",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&eventsRetrieve,
					&eventsList,
				},
			},
			{
				Name:     "extractions",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&extractionsRetrieve,
					&extractionsList,
					&extractionsEstimateCost,
					&extractionsExportResults,
					&extractionsRun,
				},
			},
			{
				Name:     "draws",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&drawsRetrieve,
					&drawsList,
					&drawsExport,
					&drawsRun,
				},
			},
			{
				Name:     "webhooks",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&webhooksCreate,
					&webhooksUpdate,
					&webhooksList,
					&webhooksDeactivate,
					&webhooksListDeliveries,
					&webhooksTest,
				},
			},
			{
				Name:     "integrations",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&integrationsCreate,
					&integrationsRetrieve,
					&integrationsUpdate,
					&integrationsList,
					&integrationsDelete,
					&integrationsListDeliveries,
					&integrationsSendTest,
				},
			},
			{
				Name:     "x",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xGetArticle,
					&xGetHomeTimeline,
					&xGetNotifications,
					&xGetTrends,
				},
			},
			{
				Name:     "x:tweets",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xTweetsCreate,
					&xTweetsRetrieve,
					&xTweetsList,
					&xTweetsDelete,
					&xTweetsGetFavoriters,
					&xTweetsGetQuotes,
					&xTweetsGetReplies,
					&xTweetsGetRetweeters,
					&xTweetsGetThread,
					&xTweetsSearch,
				},
			},
			{
				Name:     "x:tweets:like",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xTweetsLikeCreate,
					&xTweetsLikeDelete,
				},
			},
			{
				Name:     "x:tweets:retweet",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xTweetsRetweetCreate,
					&xTweetsRetweetDelete,
				},
			},
			{
				Name:     "x:users",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xUsersRetrieve,
					&xUsersRetrieveBatch,
					&xUsersRetrieveFollowers,
					&xUsersRetrieveFollowersYouKnow,
					&xUsersRetrieveFollowing,
					&xUsersRetrieveLikes,
					&xUsersRetrieveMedia,
					&xUsersRetrieveMentions,
					&xUsersRetrieveSearch,
					&xUsersRetrieveTweets,
					&xUsersRetrieveVerifiedFollowers,
				},
			},
			{
				Name:     "x:users:follow",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xUsersFollowCreate,
					&xUsersFollowDeleteAll,
				},
			},
			{
				Name:     "x:followers",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xFollowersCheck,
				},
			},
			{
				Name:     "x:dm",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xDmRetrieveHistory,
					&xDmSend,
				},
			},
			{
				Name:     "x:media",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xMediaDownload,
					&xMediaUpload,
				},
			},
			{
				Name:     "x:profile",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xProfileUpdate,
					&xProfileUpdateAvatar,
					&xProfileUpdateBanner,
				},
			},
			{
				Name:     "x:communities",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xCommunitiesCreate,
					&xCommunitiesDelete,
					&xCommunitiesRetrieveInfo,
					&xCommunitiesRetrieveMembers,
					&xCommunitiesRetrieveModerators,
					&xCommunitiesRetrieveSearch,
				},
			},
			{
				Name:     "x:communities:join",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xCommunitiesJoinCreate,
					&xCommunitiesJoinDeleteAll,
				},
			},
			{
				Name:     "x:communities:tweets",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xCommunitiesTweetsList,
				},
			},
			{
				Name:     "x:accounts",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xAccountsCreate,
					&xAccountsRetrieve,
					&xAccountsList,
					&xAccountsDelete,
					&xAccountsReauth,
				},
			},
			{
				Name:     "x:bookmarks",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xBookmarksList,
					&xBookmarksRetrieveFolders,
				},
			},
			{
				Name:     "x:lists",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&xListsRetrieveFollowers,
					&xListsRetrieveMembers,
					&xListsRetrieveTweets,
				},
			},
			{
				Name:     "trends",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&trendsList,
				},
			},
			{
				Name:     "support:tickets",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&supportTicketsCreate,
					&supportTicketsRetrieve,
					&supportTicketsUpdate,
					&supportTicketsList,
					&supportTicketsReply,
				},
			},
			{
				Name:     "credits",
				Category: "API RESOURCE",
				Suggest:  true,
				Commands: []*cli.Command{
					&creditsRetrieveBalance,
					&creditsTopupBalance,
				},
			},
			{
				Name:            "@manpages",
				Usage:           "Generate documentation for 'man'",
				UsageText:       "x-twitter-scraper @manpages [-o x-twitter-scraper.1] [--gzip]",
				Hidden:          true,
				Action:          generateManpages,
				HideHelpCommand: true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "write manpages to the given folder",
						Value:   "man",
					},
					&cli.BoolFlag{
						Name:    "gzip",
						Aliases: []string{"z"},
						Usage:   "output gzipped manpage files to .gz",
						Value:   true,
					},
					&cli.BoolFlag{
						Name:    "text",
						Aliases: []string{"z"},
						Usage:   "output uncompressed text files",
						Value:   false,
					},
				},
			},
			{
				Name:            "__complete",
				Hidden:          true,
				HideHelpCommand: true,
				Action:          autocomplete.ExecuteShellCompletion,
			},
			{
				Name:            "@completion",
				Hidden:          true,
				HideHelpCommand: true,
				Action:          autocomplete.OutputCompletionScript,
			},
		},
		HideHelpCommand: true,
	}
}

func generateManpages(ctx context.Context, c *cli.Command) error {
	manpage, err := docs.ToManWithSection(Command, 1)
	if err != nil {
		return err
	}
	dir := c.String("output")
	err = os.MkdirAll(filepath.Join(dir, "man1"), 0755)
	if err != nil {
		// handle error
	}
	if c.Bool("text") {
		file, err := os.Create(filepath.Join(dir, "man1", "x-twitter-scraper.1"))
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.WriteString(manpage); err != nil {
			return err
		}
	}
	if c.Bool("gzip") {
		file, err := os.Create(filepath.Join(dir, "man1", "x-twitter-scraper.1.gz"))
		if err != nil {
			return err
		}
		defer file.Close()
		gzWriter := gzip.NewWriter(file)
		defer gzWriter.Close()
		_, err = gzWriter.Write([]byte(manpage))
		if err != nil {
			return err
		}
	}
	fmt.Printf("Wrote manpages to %s\n", dir)
	return nil
}
