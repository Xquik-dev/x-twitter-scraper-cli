// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/mocktest"
	"github.com/stainless-sdks/x-twitter-scraper-cli/internal/requestflag"
)

func TestIntegrationsCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "create",
			"--config", "{chatId: '-1001234567890'}",
			"--event-type", "tweet.new",
			"--event-type", "follower.gained",
			"--name", "My Telegram Bot",
			"--type", "telegram",
		)
	})

	t.Run("inner flags", func(t *testing.T) {
		// Check that inner flags have been set up correctly
		requestflag.CheckInnerFlags(integrationsCreate)

		// Alternative argument passing style using inner flags
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "create",
			"--config.chat-id", "-1001234567890",
			"--event-type", "tweet.new",
			"--event-type", "follower.gained",
			"--name", "My Telegram Bot",
			"--type", "telegram",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"config:\n" +
			"  chatId: '-1001234567890'\n" +
			"eventTypes:\n" +
			"  - tweet.new\n" +
			"  - follower.gained\n" +
			"name: My Telegram Bot\n" +
			"type: telegram\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"integrations", "create",
		)
	})
}

func TestIntegrationsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "retrieve",
			"--id", "id",
		)
	})
}

func TestIntegrationsUpdate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "update",
			"--id", "id",
			"--event-type", "tweet.new",
			"--filters", "{}",
			"--is-active=true",
			"--message-template", "{}",
			"--name", "My Telegram Bot",
			"--scope-all-monitors=true",
			"--silent-push=false",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"eventTypes:\n" +
			"  - tweet.new\n" +
			"filters: {}\n" +
			"isActive: true\n" +
			"messageTemplate: {}\n" +
			"name: My Telegram Bot\n" +
			"scopeAllMonitors: true\n" +
			"silentPush: false\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"integrations", "update",
			"--id", "id",
		)
	})
}

func TestIntegrationsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "list",
		)
	})
}

func TestIntegrationsDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "delete",
			"--id", "id",
		)
	})
}

func TestIntegrationsListDeliveries(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "list-deliveries",
			"--id", "id",
			"--limit", "1",
		)
	})
}

func TestIntegrationsSendTest(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"integrations", "send-test",
			"--id", "id",
		)
	})
}
