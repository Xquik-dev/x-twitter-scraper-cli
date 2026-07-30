// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXAccountsCreate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "create",
			"--email", "account@example.invalid",
			"--password", "<ACCOUNT_PASSWORD>",
			"--totp-secret", "<TOTP_SECRET>",
			"--username", "your_x_username",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"email: account@example.invalid\n" +
			"password: <ACCOUNT_PASSWORD>\n" +
			"totp_secret: <TOTP_SECRET>\n" +
			"username: your_x_username\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "create",
		)
	})
}

func TestXAccountsRetrieve(t *testing.T) {
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

func TestXAccountsBulkRetry(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "bulk-retry",
		)
	})
}

func TestXAccountsReauth(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "reauth",
			"--id", "id",
			"--password", "<ACCOUNT_PASSWORD>",
			"--email", "account@example.invalid",
			"--totp-secret", "<TOTP_SECRET>",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"password: <ACCOUNT_PASSWORD>\n" +
			"email: account@example.invalid\n" +
			"totp_secret: <TOTP_SECRET>\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:accounts", "reauth",
			"--id", "id",
		)
	})
}
