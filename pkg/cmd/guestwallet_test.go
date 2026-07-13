// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestGuestWalletsCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"guest-wallets", "create",
			"--amount-minor", "1000",
			"--currency", "usd",
			"--idempotency-key", "e1cb97D8-dDF3-4AaA-ad0a-49E4A0d1CfAa",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"amount_minor: 1000\n" +
			"currency: usd\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"guest-wallets", "create",
			"--idempotency-key", "e1cb97D8-dDF3-4AaA-ad0a-49E4A0d1CfAa",
		)
	})
}

func TestGuestWalletsRetrieveStatus(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"guest-wallets", "retrieve-status",
		)
	})
}

func TestGuestWalletsTopup(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"guest-wallets", "topup",
			"--amount-minor", "1000",
			"--currency", "usd",
			"--idempotency-key", "e1cb97D8-dDF3-4AaA-ad0a-49E4A0d1CfAa",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"amount_minor: 1000\n" +
			"currency: usd\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"--cookie-session", "string",
			"guest-wallets", "topup",
			"--idempotency-key", "e1cb97D8-dDF3-4AaA-ad0a-49E4A0d1CfAa",
		)
	})
}
