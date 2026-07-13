// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXMediaDownload(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"x:media", "download",
			"--tweet-id", "1234567890",
			"--tweet-id", "1234567890",
			"--tweet-id", "1234567891",
			"--tweet-input", "https://x.com/elonmusk/status/1234567890",
			"--tweet-url", "https://x.com/elonmusk/status/1234567890",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"tweetId: '1234567890'\n" +
			"tweetIds:\n" +
			"  - '1234567890'\n" +
			"  - '1234567891'\n" +
			"tweetInput: https://x.com/elonmusk/status/1234567890\n" +
			"tweetUrl: https://x.com/elonmusk/status/1234567890\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
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
			"--bearer-token", "string",
			"--cookie-session", "string",
			"x:media", "upload",
			"--account", "@elonmusk",
			"--url", "https://example.com/image.png",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"url: https://example.com/image.png\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"x:media", "upload",
		)
	})
}
