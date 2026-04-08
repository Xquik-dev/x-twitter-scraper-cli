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
			"--tweet-url", "https://x.com/elonmusk/status/1234567890",
			"--backup-count", "2",
			"--filter-account-age-days", "30",
			"--filter-language", "en",
			"--filter-min-followers", "50",
			"--must-follow-username", "elonmusk",
			"--must-retweet=true",
			"--required-hashtag", "#giveaway",
			"--required-keyword", "entered",
			"--required-mention", "@elonmusk",
			"--unique-authors-only=true",
			"--winner-count", "3",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"tweetUrl: https://x.com/elonmusk/status/1234567890\n" +
			"backupCount: 2\n" +
			"filterAccountAgeDays: 30\n" +
			"filterLanguage: en\n" +
			"filterMinFollowers: 50\n" +
			"mustFollowUsername: elonmusk\n" +
			"mustRetweet: true\n" +
			"requiredHashtags:\n" +
			"  - '#giveaway'\n" +
			"requiredKeywords:\n" +
			"  - entered\n" +
			"requiredMentions:\n" +
			"  - '@elonmusk'\n" +
			"uniqueAuthorsOnly: true\n" +
			"winnerCount: 3\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"draws", "run",
		)
	})
}
