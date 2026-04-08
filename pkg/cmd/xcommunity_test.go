// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestXCommunitiesCreate(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "create",
			"--account", "@elonmusk",
			"--name", "Example Name",
			"--description", "A community for Tesla enthusiasts",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"name: Example Name\n" +
			"description: A community for Tesla enthusiasts\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "create",
		)
	})
}

func TestXCommunitiesDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "delete",
			"--id", "id",
			"--account", "@elonmusk",
			"--community-name", "Tesla Fans",
		)
	})

	t.Run("piping data", func(t *testing.T) {
		// Test piping YAML data over stdin
		pipeData := []byte("" +
			"account: '@elonmusk'\n" +
			"community_name: Tesla Fans\n")
		mocktest.TestRunMockTestWithPipeAndFlags(
			t, pipeData,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "delete",
			"--id", "id",
		)
	})
}

func TestXCommunitiesRetrieveInfo(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "retrieve-info",
			"--id", "id",
		)
	})
}

func TestXCommunitiesRetrieveMembers(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "retrieve-members",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXCommunitiesRetrieveModerators(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "retrieve-moderators",
			"--id", "id",
			"--cursor", "cursor",
		)
	})
}

func TestXCommunitiesRetrieveSearch(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"x:communities", "retrieve-search",
			"--q", "q",
			"--cursor", "cursor",
			"--query-type", "queryType",
		)
	})
}
