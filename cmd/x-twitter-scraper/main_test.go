// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/pkg/cmd"
	"github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func unsetEnvironment(string) (string, bool) {
	return "", false
}

func TestRunSuccessAndAutocomplete(t *testing.T) {
	child := &cli.Command{Name: "child"}
	app := &cli.Command{
		Name:     "test",
		Commands: []*cli.Command{child},
		Action: func(context.Context, *cli.Command) error {
			return nil
		},
	}
	var stderr bytes.Buffer

	assert.Zero(t, run(app, []string{"test", "__complete"}, unsetEnvironment, &stderr))
	assert.True(t, app.SkipFlagParsing)
	assert.True(t, child.SkipFlagParsing)
	assert.Empty(t, stderr.String())
}

func TestRunRejectsInvalidBaseURL(t *testing.T) {
	lookupEnv := func(key string) (string, bool) {
		require.Equal(t, "X_TWITTER_SCRAPER_BASE_URL", key)
		return "://invalid", true
	}
	var stderr bytes.Buffer

	assert.Equal(t, 1, run(&cli.Command{Name: "test"}, []string{"test"}, lookupEnv, &stderr))
	assert.Contains(t, stderr.String(), "X_TWITTER_SCRAPER_BASE_URL")
}

func TestRunReportsGenericAndExitErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		exitCode int
	}{
		{name: "generic", err: errors.New("failed"), exitCode: 1},
		{name: "custom exit", err: cli.Exit("failed", 7), exitCode: 7},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := &cli.Command{
				Name: "test",
				Action: func(context.Context, *cli.Command) error {
					return test.err
				},
			}
			var stderr bytes.Buffer

			assert.Equal(t, test.exitCode, run(app, []string{"test"}, unsetEnvironment, &stderr))
			assert.Contains(t, stderr.String(), "failed")
		})
	}
}

func TestRunUsesCommandErrorBuffer(t *testing.T) {
	cmd.CommandErrorBuffer.Reset()
	t.Cleanup(cmd.CommandErrorBuffer.Reset)
	cmd.CommandErrorBuffer.WriteString("buffered error\n")

	app := &cli.Command{
		Name: "test",
		Action: func(context.Context, *cli.Command) error {
			return errors.New("fallback error")
		},
	}
	var stderr bytes.Buffer

	assert.Equal(t, 1, run(app, []string{"test"}, unsetEnvironment, &stderr))
	assert.Equal(t, "buffered error\n", stderr.String())
}

func TestRunReportsAPIError(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "https://example.com/resource", nil)
	require.NoError(t, err)
	apiError := &xtwitterscraper.Error{
		Request: request,
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
	}
	require.NoError(t, apiError.UnmarshalJSON([]byte(`{"error":"invalid"}`)))

	app := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "format-error", Value: "invalid"},
			&cli.StringFlag{Name: "transform-error"},
		},
		Action: func(context.Context, *cli.Command) error {
			return apiError
		},
	}
	var stderr bytes.Buffer

	assert.Equal(t, 1, run(app, []string{"test"}, unsetEnvironment, &stderr))
	assert.Contains(t, stderr.String(), `POST "https://example.com/resource": 400 Bad Request`)
	assert.Contains(t, stderr.String(), "invalid")
}
