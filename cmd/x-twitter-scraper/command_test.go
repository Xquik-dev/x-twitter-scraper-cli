// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func buildCLI(t *testing.T) string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "x-twitter-scraper")
	build := exec.Command("go", "build", "-cover", "-coverpkg=./...", "-o", binary, "./cmd/x-twitter-scraper")
	build.Dir = "../.."
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	return binary
}

func invokeCLI(t *testing.T, binary string, args ...string) (string, string, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, binary, args...)
	command.Env = cliEnvironment(t)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if ctx.Err() != nil {
		t.Fatalf("CLI timed out: %v", ctx.Err())
	}
	exit := 0
	if err != nil {
		failure, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatal(err)
		}
		exit = failure.ExitCode()
	}
	return stdout.String(), stderr.String(), exit
}

func cliEnvironment(t *testing.T) []string {
	t.Helper()
	var environment []string
	for _, variable := range os.Environ() {
		if !strings.HasPrefix(variable, "X_TWITTER_SCRAPER_") && !strings.HasPrefix(variable, "GOCOVERDIR=") {
			environment = append(environment, variable)
		}
	}
	counters := os.Getenv("XQUIK_GO_COVERAGE_DIR")
	if counters == "" {
		counters = t.TempDir()
	}
	return append(environment, "GOCOVERDIR="+counters)
}

func TestCLIAccountReadPreservesJSONTransform(t *testing.T) {
	binary := buildCLI(t)
	var mutex sync.Mutex
	var requests []*http.Request
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mutex.Lock()
		requests = append(requests, request.Clone(request.Context()))
		mutex.Unlock()
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"id":"account-1","locale":"en"}`))
	}))
	t.Cleanup(server.Close)
	stdout, stderr, exit := invokeCLI(t, binary, "--base-url", server.URL, "--api-key", "test-api-key", "--format", "json", "--transform", "id", "account", "retrieve")
	if exit != 0 || stderr != "" || stdout != "\"account-1\"\n" {
		t.Fatalf("account read returned exit %d, stdout %q, stderr %q", exit, stdout, stderr)
	}
	mutex.Lock()
	defer mutex.Unlock()
	if len(requests) != 1 {
		t.Fatalf("account read made %d HTTP requests, want 1", len(requests))
	}
	request := requests[0]
	if request.Method != "GET" || request.URL.Path != "/account" || request.Header.Get("X-API-Key") != "test-api-key" {
		t.Fatalf("request = %s %s, API key %q", request.Method, request.URL.Path, request.Header.Get("X-API-Key"))
	}
}

func TestCLIAPIErrorUsesStderr(t *testing.T) {
	binary := buildCLI(t)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(`{"code":"AUTH_INVALID","message":"API key is invalid."}`))
	}))
	t.Cleanup(server.Close)
	stdout, stderr, exit := invokeCLI(t, binary, "--base-url", server.URL, "--api-key", "test-api-key", "--format-error", "json", "account", "retrieve")
	if exit != 1 || stdout != "" || !strings.Contains(stderr, "401 Unauthorized") || !strings.Contains(stderr, `"code": "AUTH_INVALID"`) {
		t.Fatalf("API error returned exit %d, stdout %q, stderr %q", exit, stdout, stderr)
	}
}
