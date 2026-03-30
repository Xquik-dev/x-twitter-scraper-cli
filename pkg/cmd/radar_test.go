// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestRadarRetrieveTrendingTopics(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"radar", "retrieve-trending-topics",
			"--category", "category",
			"--count", "0",
			"--hours", "0",
			"--region", "region",
			"--source", "source",
		)
	})
}
