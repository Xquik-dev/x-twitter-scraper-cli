// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestComposeCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"compose", "create",
			"--step", "compose",
			"--additional-context", "https://x.com/elonmusk/status/1234567890",
			"--call-to-action", "Follow for more",
			"--draft", "AI is changing everything. Here's why.",
			"--goal", "engagement",
			"--has-link=false",
			"--has-media=false",
			"--media-type", "none",
			"--style-username", "elonmusk",
			"--tone", "professional",
			"--topic", "AI trends in 2025",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"step: compose\n" +
			"additionalContext: https://x.com/elonmusk/status/1234567890\n" +
			"callToAction: Follow for more\n" +
			"draft: AI is changing everything. Here's why.\n" +
			"goal: engagement\n" +
			"hasLink: false\n" +
			"hasMedia: false\n" +
			"mediaType: none\n" +
			"styleUsername: elonmusk\n" +
			"tone: professional\n" +
			"topic: AI trends in 2025\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"compose", "create",
		)
	})
}
