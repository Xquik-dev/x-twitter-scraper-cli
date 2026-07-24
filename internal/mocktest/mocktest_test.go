// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package mocktest

import (
	"context"
	"errors"
	"net/http"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandHelperScenarios(t *testing.T) {
	RegisterRunner(func(_ *testing.T, ctx context.Context, args []string, pipeData []byte) error {
		if string(pipeData) == "{" {
			return errors.New("invalid piped data")
		}
		if slices.Contains(args, "unexpected") {
			return errors.New("unexpected argument")
		}

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, args[2], nil)
		if err != nil {
			return err
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode >= http.StatusBadRequest {
			return errors.New("request failed")
		}
		return nil
	})

	TestRunMockTestWithFlags(t, "resource", "action")
	TestRunMockTestWithPipeAndFlags(t, []byte("{}"), "resource", "action")
}

func TestFileCreatesFixture(t *testing.T) {
	path := TestFile(t, "contents")
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "contents", string(contents))
}

func TestNetworkGuardRejectsOtherServers(t *testing.T) {
	restore := restrictNetworkToMockServer()
	t.Cleanup(restore)

	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://example.invalid",
		nil,
	)
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.ErrorContains(t, err, "blocked test network connection")
	require.Nil(t, response)
}
