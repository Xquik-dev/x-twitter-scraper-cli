// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestRadarRetrieveTrendingTopics(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"radar", "retrieve-trending-topics",
			"--after", "after",
			"--category", "general",
			"--hours", "1",
			"--limit", "1",
			"--region", "region",
			"--source", "github",
		)
	})
}
