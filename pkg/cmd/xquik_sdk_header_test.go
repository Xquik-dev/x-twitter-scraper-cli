// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	xtwitterscraper "github.com/Xquik-dev/x-twitter-scraper-go"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestEveryRequestNamesTheCLIForXquik(t *testing.T) {
	headers := make(chan http.Header, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers <- r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	options := append(getDefaultRequestOptions(&cli.Command{Name: "retrieve"}),
		option.WithBaseURL(server.URL), option.WithAPIKey("My API Key"))
	_, err := xtwitterscraper.NewClient(options...).Account.Get(context.Background())
	require.NoError(t, err)

	sent := <-headers
	assert.Equal(t, "cli", sent.Get("X-Xquik-SDK"))
	assert.Equal(t, "cli", sent.Get("X-Stainless-Lang"))
}
