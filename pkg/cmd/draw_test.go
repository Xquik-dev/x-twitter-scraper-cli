// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestDrawsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"draws", "retrieve",
			"--id", "id",
		)
	})
}

func TestDrawsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"draws", "list",
			"--after", "after",
			"--limit", "1",
		)
	})
}

func TestDrawsExport(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"draws", "export",
			"--id", "id",
			"--format", "csv",
			"--type", "winners",
			"--output", "/dev/null",
		)
	})
}

func TestDrawsRun(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"draws", "run",
			"--tweet-url", "https://example.com",
			"--backup-count", "0",
			"--filter-account-age-days", "0",
			"--filter-language", "filterLanguage",
			"--filter-min-followers", "0",
			"--must-follow-username", "mustFollowUsername",
			"--must-retweet=true",
			"--required-hashtag", "string",
			"--required-keyword", "string",
			"--required-mention", "string",
			"--unique-authors-only=true",
			"--winner-count", "0",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"tweetUrl: https://example.com\n" +
			"backupCount: 0\n" +
			"filterAccountAgeDays: 0\n" +
			"filterLanguage: filterLanguage\n" +
			"filterMinFollowers: 0\n" +
			"mustFollowUsername: mustFollowUsername\n" +
			"mustRetweet: true\n" +
			"requiredHashtags:\n" +
			"  - string\n" +
			"requiredKeywords:\n" +
			"  - string\n" +
			"requiredMentions:\n" +
			"  - string\n" +
			"uniqueAuthorsOnly: true\n" +
			"winnerCount: 0\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"draws", "run",
		)
	})
}
