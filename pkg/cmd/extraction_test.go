// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestExtractionsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
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
			"--bearer-token", "string",
			"extractions", "list",
			"--after", "after",
			"--limit", "1",
			"--status", "running",
			"--tool-type", "article_extractor",
		)
	})
}

func TestExtractionsEstimateCost(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "estimate-cost",
			"--tool-type", "article_extractor",
			"--advanced-query", "advancedQuery",
			"--exact-phrase", "exactPhrase",
			"--exclude-words", "excludeWords",
			"--search-query", "searchQuery",
			"--target-community-id", "targetCommunityId",
			"--target-list-id", "targetListId",
			"--target-space-id", "targetSpaceId",
			"--target-tweet-id", "targetTweetId",
			"--target-username", "targetUsername",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: article_extractor\n" +
			"advancedQuery: advancedQuery\n" +
			"exactPhrase: exactPhrase\n" +
			"excludeWords: excludeWords\n" +
			"searchQuery: searchQuery\n" +
			"targetCommunityId: targetCommunityId\n" +
			"targetListId: targetListId\n" +
			"targetSpaceId: targetSpaceId\n" +
			"targetTweetId: targetTweetId\n" +
			"targetUsername: targetUsername\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
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
			"--bearer-token", "string",
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
			"--bearer-token", "string",
			"extractions", "run",
			"--tool-type", "article_extractor",
			"--advanced-query", "advancedQuery",
			"--exact-phrase", "exactPhrase",
			"--exclude-words", "excludeWords",
			"--search-query", "searchQuery",
			"--target-community-id", "targetCommunityId",
			"--target-list-id", "targetListId",
			"--target-space-id", "targetSpaceId",
			"--target-tweet-id", "targetTweetId",
			"--target-username", "targetUsername",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"toolType: article_extractor\n" +
			"advancedQuery: advancedQuery\n" +
			"exactPhrase: exactPhrase\n" +
			"excludeWords: excludeWords\n" +
			"searchQuery: searchQuery\n" +
			"targetCommunityId: targetCommunityId\n" +
			"targetListId: targetListId\n" +
			"targetSpaceId: targetSpaceId\n" +
			"targetTweetId: targetTweetId\n" +
			"targetUsername: targetUsername\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"extractions", "run",
		)
	})
}
