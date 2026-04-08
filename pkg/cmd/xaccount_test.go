// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXAccountsCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "create",
			"--email", "user@example.com",
			"--password", "s3cur3Pa$$w0rd",
			"--username", "elonmusk",
			"--proxy-country", "US",
			"--totp-secret", "JBSWY3DPEHPK3PXP",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"email: user@example.com\n" +
			"password: s3cur3Pa$$w0rd\n" +
			"username: elonmusk\n" +
			"proxy_country: US\n" +
			"totp_secret: JBSWY3DPEHPK3PXP\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "create",
		)
	})
}

func TestXAccountsRetrieve(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "retrieve",
			"--id", "id",
		)
	})
}

func TestXAccountsList(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "list",
		)
	})
}

func TestXAccountsDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "delete",
			"--id", "id",
		)
	})
}

func TestXAccountsReauth(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "reauth",
			"--id", "id",
			"--password", "password_value",
			"--totp-secret", "totp_secret_value",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"password: password_value\n" +
			"totp_secret: totp_secret_value\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "reauth",
			"--id", "id",
		)
	})
}
