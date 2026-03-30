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
			"--bearer-token", "string",
			"integrations", "create",
			"--config", "{chatId: chatId}",
			"--event-type", "tweet.new",
			"--name", "name",
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
			"--bearer-token", "string",
			"integrations", "create",
			"--config.chat-id", "chatId",
			"--event-type", "tweet.new",
			"--name", "name",
			"--type", "telegram",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"config:\n" +
			"  chatId: chatId\n" +
			"eventTypes:\n" +
			"  - tweet.new\n" +
			"name: name\n" +
			"type: telegram\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
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
			"--bearer-token", "string",
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
			"--bearer-token", "string",
			"integrations", "update",
			"--id", "id",
			"--event-type", "tweet.new",
			"--filters", "{foo: bar}",
			"--is-active=true",
			"--message-template", "{foo: bar}",
			"--name", "name",
			"--scope-all-monitors=true",
			"--silent-push=true",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"eventTypes:\n" +
			"  - tweet.new\n" +
			"filters:\n" +
			"  foo: bar\n" +
			"isActive: true\n" +
			"messageTemplate:\n" +
			"  foo: bar\n" +
			"name: name\n" +
			"scopeAllMonitors: true\n" +
			"silentPush: true\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
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
			"--bearer-token", "string",
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
			"--bearer-token", "string",
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
			"--bearer-token", "string",
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
			"--bearer-token", "string",
			"integrations", "send-test",
			"--id", "id",
		)
	})
}
