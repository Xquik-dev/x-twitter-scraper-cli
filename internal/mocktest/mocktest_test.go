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
