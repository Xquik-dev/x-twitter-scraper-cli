// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"strings"
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXMediaDownload(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"x:media", "download",
			"--tweet-id", "1234567890",
			"--tweet-id", "1234567891",
			"--tweet-input", "https://x.com/elonmusk/status/1234567890",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"tweetIds:\n" +
			"  - '1234567890'\n" +
			"  - '1234567891'\n" +
			"tweetInput: https://x.com/elonmusk/status/1234567890\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"x:media", "download",
		)
	})
}

func TestXMediaUpload(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"x:media", "upload",
			"--account", "@elonmusk",
			"--file", mocktest.TestFile(t, "Example data"),
			"--is-long-video=true",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		testFile := mocktest.TestFile(t, "Example data")
		// Test piping YAML data over stdin
		pipeDataStr := "" +
			"account: '@elonmusk'\n" +
			"file: Example data\n" +
			"is_long_video: true\n"
		pipeDataStr = strings.ReplaceAll(pipeDataStr, "Example data", testFile)
		pipeData := []byte(pipeDataStr)
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"x:media", "upload",
		)
	})
}
