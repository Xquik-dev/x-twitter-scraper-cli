// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
)

func TestComposeCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"compose", "create",
			"--step", "compose",
			"--additional-context", "additionalContext",
			"--call-to-action", "callToAction",
			"--draft", "draft",
			"--goal", "engagement",
			"--has-link=true",
			"--has-media=true",
			"--media-type", "photo",
			"--style-username", "styleUsername",
			"--tone", "tone",
			"--topic", "topic",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"step: compose\n" +
			"additionalContext: additionalContext\n" +
			"callToAction: callToAction\n" +
			"draft: draft\n" +
			"goal: engagement\n" +
			"hasLink: true\n" +
			"hasMedia: true\n" +
			"mediaType: photo\n" +
			"styleUsername: styleUsername\n" +
			"tone: tone\n" +
			"topic: topic\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"compose", "create",
		)
	})
}
