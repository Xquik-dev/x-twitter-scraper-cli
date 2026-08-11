// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXListsRetrieveFollowers(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:lists", "retrieve-followers",
			"--id", "id",
			"--bio-contains", "bioContains",
			"--cursor", "cursor",
			"--has-location=true",
			"--has-website=true",
			"--location-contains", "locationContains",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-statuses", "0",
			"--min-account-age-days", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-statuses", "0",
			"--page-size", "20",
			"--username-contains", "usernameContains",
			"--verified-only=true",
			"--verified-type", "verifiedType",
		)
	})
}

func TestXListsRetrieveMembers(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:lists", "retrieve-members",
			"--id", "id",
			"--bio-contains", "bioContains",
			"--cursor", "cursor",
			"--has-location=true",
			"--has-website=true",
			"--location-contains", "locationContains",
			"--max-followers", "0",
			"--max-following", "0",
			"--max-statuses", "0",
			"--min-account-age-days", "0",
			"--min-followers", "0",
			"--min-following", "0",
			"--min-statuses", "0",
			"--page-size", "20",
			"--username-contains", "usernameContains",
			"--verified-only=true",
			"--verified-type", "verifiedType",
		)
	})
}

func TestXListsRetrieveTweets(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:lists", "retrieve-tweets",
			"--id", "id",
			"--cursor", "cursor",
			"--include-replies=true",
			"--page-size", "1",
			"--since-time", "sinceTime",
			"--until-time", "untilTime",
		)
	})
}
