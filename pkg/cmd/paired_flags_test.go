// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const pairedFlagsArgsVariable = "XQUIK_CLI_PAIRED_FLAGS_ARGS"

func init() {
	encoded, ok := os.LookupEnv(pairedFlagsArgsVariable)
	if !ok {
		return
	}
	var args []string
	if err := json.Unmarshal([]byte(encoded), &args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	command := newCommand()
	command.Writer = io.Discard
	command.ErrWriter = io.Discard
	if err := command.Run(context.Background(), args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func runCLIProcess(t *testing.T, stdin string, args ...string) (string, int) {
	t.Helper()
	encoded, err := json.Marshal(append([]string{"x-twitter-scraper"}, args...))
	require.NoError(t, err)
	child := exec.Command(os.Args[0], "-test.run=^$")
	child.Env = append(os.Environ(), pairedFlagsArgsVariable+"="+string(encoded))
	child.Stdin = strings.NewReader(stdin)
	output, err := child.CombinedOutput()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return string(output), exitErr.ExitCode()
	}
	require.NoError(t, err)
	return string(output), 0
}

type recordedRequest struct {
	method string
	path   string
	body   string
}

func runAgainstRecorder(t *testing.T, stdin string, args ...string) recordedRequest {
	t.Helper()
	var mu sync.Mutex
	var requests []recordedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		requests = append(requests, recordedRequest{method: r.Method, path: r.URL.Path, body: string(body)})
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	defer server.Close()

	output, code := runCLIProcess(t, stdin, append([]string{"--base-url", server.URL, "--api-key", "test-key"}, args...)...)
	require.Equal(t, 0, code, "the CLI run failed with this output\n%s", output)

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, requests, 1)
	return requests[0]
}

func TestPairedFlagsKeepTheFieldsThatRelease0133Sent(t *testing.T) {
	extractionFlags := []string{
		"--tool-type", "tweet_search_extractor",
		"--search-query", "alpha",
		"--search-query", "beta",
		"--target-community-id", "community-1",
		"--target-list-id", "list-1",
		"--target-tweet-id", "tweet-1",
		"--target-username", "alice",
	}
	extractionBody := `{"searchQueries":["alpha","beta"],"targetCommunityId":"community-1",` +
		`"targetListId":"list-1","targetTweetId":"tweet-1","targetUsername":"alice",` +
		`"toolType":"tweet_search_extractor"}`

	tests := []struct {
		name string
		args []string
		path string
		body string
	}{
		{
			name: "extractions estimate-cost",
			args: append([]string{"extractions", "estimate-cost"}, extractionFlags...),
			path: "/extractions/estimate",
			body: extractionBody,
		},
		{
			name: "extractions run",
			args: append([]string{"extractions", "run"}, extractionFlags...),
			path: "/extractions",
			body: extractionBody,
		},
		{
			name: "x:media download",
			args: []string{"x:media", "download", "--tweet-id", "123"},
			path: "/x/media/download",
			body: `{"tweetId":"123"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := runAgainstRecorder(t, "", test.args...)
			require.Equal(t, http.MethodPost, request.method)
			require.Equal(t, test.path, request.path)
			require.JSONEq(t, test.body, request.body)
		})
	}
}

func TestPipedInputReachesTheOtherFieldOfEachPair(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stdin string
		body  string
	}{
		{
			name: "extractions estimate-cost",
			args: []string{"extractions", "estimate-cost", "--tool-type", "community_search"},
			stdin: "searchQuery: startup\n" +
				"targetCommunityIds: [community-1, community-2]\n" +
				"targetListIds: [list-1]\n" +
				"targetTweetIds: [tweet-1]\n" +
				"targetUsernames: [alice, bob]\n",
			body: `{"searchQuery":"startup","targetCommunityIds":["community-1","community-2"],` +
				`"targetListIds":["list-1"],"targetTweetIds":["tweet-1"],` +
				`"targetUsernames":["alice","bob"],"toolType":"community_search"}`,
		},
		{
			name:  "x:media download",
			args:  []string{"x:media", "download"},
			stdin: "tweetIds: [\"1\", \"2\"]\n",
			body:  `{"tweetIds":["1","2"]}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := runAgainstRecorder(t, test.stdin, test.args...)
			require.JSONEq(t, test.body, request.body)
		})
	}
}

func TestManpagesShortFlagZStillSetsGzip(t *testing.T) {
	output := t.TempDir()
	result, code := runCLIProcess(t, "", "@manpages", "--output", output, "--text", "-z=false")
	require.Equal(t, 0, code, "the CLI run failed with this output\n%s", result)

	entries, err := os.ReadDir(filepath.Join(output, "man1"))
	require.NoError(t, err)
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	require.Equal(t, []string{"x-twitter-scraper.1"}, names)
}
