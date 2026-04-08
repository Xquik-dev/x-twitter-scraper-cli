// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestXDmRetrieveHistory(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:dm", "retrieve-history",
			"--user-id", "userId",
			"--cursor", "cursor",
			"--max-id", "maxId",
		)
	})
}

func TestXDmSend(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:dm", "send",
			"--user-id", "userId",
			"--account", "@elonmusk",
			"--text", "Example text content",
			"--media-id", "1234567890123456789",
			"--reply-to-message-id", "1234567890123456789",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"text: Example text content\n" +
			"media_ids:\n" +
			"  - '1234567890123456789'\n" +
			"reply_to_message_id: '1234567890123456789'\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:dm", "send",
			"--user-id", "userId",
		)
	})
}
