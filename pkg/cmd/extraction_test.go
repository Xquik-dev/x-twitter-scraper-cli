// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestExtractionsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"extractions", "retrieve",
			"--id", "id",
			"--after", "after",
			"--limit", "1",
		)
	})
}

func TestExtractionsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"extractions", "list",
			"--after", "after",
			"--limit", "1",
			"--status", "running",
			"--tool-type", "follower_explorer",
		)
	})
}

func TestExtractionsEstimateCost(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"extractions", "estimate-cost",
			"--tool-type", "follower_explorer",
			"--advanced-query", "min_faves:100",
			"--exact-phrase", "artificial intelligence",
			"--exclude-words", "spam",
			"--search-query", "AI trends 2025",
			"--target-community-id", "1500000000000000000",
			"--target-list-id", "1234567890",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-username", "elonmusk",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: follower_explorer\n" +
			"advancedQuery: min_faves:100\n" +
			"exactPhrase: artificial intelligence\n" +
			"excludeWords: spam\n" +
			"searchQuery: AI trends 2025\n" +
			"targetCommunityId: '1500000000000000000'\n" +
			"targetListId: '1234567890'\n" +
			"targetSpaceId: 1vOGwMdBqpwGB\n" +
			"targetTweetId: '1234567890'\n" +
			"targetUsername: elonmusk\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"extractions", "estimate-cost",
		)
	})
}

func TestExtractionsExportResults(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"extractions", "export-results",
			"--id", "id",
			"--format", "csv",
			"--output", "/dev/null",
		)
	})
}

func TestExtractionsRun(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"extractions", "run",
			"--tool-type", "follower_explorer",
			"--advanced-query", "min_faves:100",
			"--exact-phrase", "artificial intelligence",
			"--exclude-words", "spam",
			"--search-query", "AI trends 2025",
			"--target-community-id", "1500000000000000000",
			"--target-list-id", "1234567890",
			"--target-space-id", "1vOGwMdBqpwGB",
			"--target-tweet-id", "1234567890",
			"--target-username", "elonmusk",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: follower_explorer\n" +
			"advancedQuery: min_faves:100\n" +
			"exactPhrase: artificial intelligence\n" +
			"excludeWords: spam\n" +
			"searchQuery: AI trends 2025\n" +
			"targetCommunityId: '1500000000000000000'\n" +
			"targetListId: '1234567890'\n" +
			"targetSpaceId: 1vOGwMdBqpwGB\n" +
			"targetTweetId: '1234567890'\n" +
			"targetUsername: elonmusk\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"extractions", "run",
		)
	})
}
