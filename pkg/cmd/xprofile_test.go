// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXProfileUpdate(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:profile", "update",
			"--account", "@elonmusk",
			"--idempotency-key", "Idempotency-Key",
			"--description", "description_value",
			"--location", "location_value",
			"--name", "Example Name",
			"--url", "https://xquik.com/example",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"description: description_value\n" +
			"location: location_value\n" +
			"name: Example Name\n" +
			"url: https://xquik.com/example\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:profile", "update",
			"--idempotency-key", "Idempotency-Key",
		)
	})
}

func TestXProfileUpdateAvatar(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:profile", "update-avatar",
			"--account", "@elonmusk",
			"--url", "https://example.com/avatar.png",
			"--idempotency-key", "Idempotency-Key",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"url: https://example.com/avatar.png\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:profile", "update-avatar",
			"--idempotency-key", "Idempotency-Key",
		)
	})
}

func TestXProfileUpdateBanner(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:profile", "update-banner",
			"--account", "@elonmusk",
			"--url", "https://example.com/banner.png",
			"--idempotency-key", "Idempotency-Key",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"url: https://example.com/banner.png\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:profile", "update-banner",
			"--idempotency-key", "Idempotency-Key",
		)
	})
}
