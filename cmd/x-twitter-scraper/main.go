// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/Xquik-dev/x-twitter-scraper-cli/pkg/cmd"
	"github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

func main() {
	exitCode := run(cmd.Command, os.Args, os.LookupEnv, os.Stderr)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func run(
	app *cli.Command,
	args []string,
	lookupEnv func(string) (string, bool),
	stderr io.Writer,
) int {
	if slices.Contains(args, "__complete") {
		prepareForAutocomplete(app)
	}

	if baseURL, ok := lookupEnv("X_TWITTER_SCRAPER_BASE_URL"); ok {
		if err := cmd.ValidateBaseURL(baseURL, "X_TWITTER_SCRAPER_BASE_URL"); err != nil {
			fmt.Fprintf(stderr, "%s\n", err.Error())
			return 1
		}
	}

	previousExitErrHandler := app.ExitErrHandler
	app.ExitErrHandler = func(context.Context, *cli.Command, error) {}
	defer func() {
		app.ExitErrHandler = previousExitErrHandler
	}()

	if err := app.Run(context.Background(), args); err != nil {
		exitCode := 1

		// Check if error has a custom exit code
		if exitErr, ok := err.(cli.ExitCoder); ok {
			exitCode = exitErr.ExitCode()
		}

		var apierr *xtwitterscraper.Error
		if errors.As(err, &apierr) {
			fmt.Fprintf(stderr, "%s %q: %d %s\n", apierr.Request.Method, apierr.Request.URL, apierr.Response.StatusCode, http.StatusText(apierr.Response.StatusCode))
			format := app.String("format-error")
			json := gjson.Parse(apierr.RawJSON())
			showErr := cmd.ShowJSON(json, cmd.ShowJSONOpts{
				ExplicitFormat: app.IsSet("format-error"),
				Format:         format,
				Title:          "Error",
				Transform:      app.String("transform-error"),
			})
			if showErr != nil {
				// Just print the original error:
				fmt.Fprintf(stderr, "%s\n", err.Error())
			}
		} else {
			if cmd.CommandErrorBuffer.Len() > 0 {
				_, _ = stderr.Write(cmd.CommandErrorBuffer.Bytes())
			} else {
				fmt.Fprintf(stderr, "%s\n", err.Error())
			}
		}
		return exitCode
	}
	return 0
}

func prepareForAutocomplete(cmd *cli.Command) {
	// urfave/cli does not handle flag completions and will print an error if we inspect a command with invalid flags.
	// This skips that sort of validation
	cmd.SkipFlagParsing = true
	for _, child := range cmd.Commands {
		prepareForAutocomplete(child)
	}
}
