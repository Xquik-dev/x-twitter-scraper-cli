// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

type explorerTerminal struct {
	output bytes.Buffer
	input  io.Writer
	steps  []explorerStep
	offset int
}

type explorerStep struct {
	waitFor string
	key     string
}

func (terminal *explorerTerminal) String() string { return terminal.output.String() }

func (terminal *explorerTerminal) Write(data []byte) (int, error) {
	written, err := terminal.output.Write(data)
	if err == nil && len(terminal.steps) > 0 && strings.Contains(terminal.String()[terminal.offset:], terminal.steps[0].waitFor) {
		key := terminal.steps[0].key
		terminal.steps = terminal.steps[1:]
		terminal.offset = terminal.output.Len()
		_, err = io.WriteString(terminal.input, key)
	}
	return written, err
}

func TestCLIExplorerPrintsSelectedValue(t *testing.T) {
	binary := buildCLI(t)
	for _, scenario := range []struct {
		name     string
		response string
		steps    []explorerStep
		want     string
	}{
		{"object", `{"id":"account-1","locale":"en"}`, []explorerStep{{"account-1", "p"}}, "account-1"},
		{"nested object", `{"profile":{"id":"account-1"}}`, []explorerStep{{"profile", "l"}, {"→", "p"}}, "account-1"},
		{"raw mode", `{"id":"account-1"}`, []explorerStep{{"account-1", "r"}, {"(JSON)", "p"}}, "account-1"},
		{"expanded text", `{"description":"line one\nline two"}`, []explorerStep{{"description", "l"}, {"→", "p"}}, `"line one\nline two"`},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			runExplorerTerminal(t, binary, scenario.response, scenario.steps, scenario.want)
		})
	}
}

func runExplorerTerminal(t *testing.T, binary, response string, steps []explorerStep, want string) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	quote := func(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'" }
	arguments := []string{binary, "--base-url", server.URL, "--api-key", "test-api-key", "--format", "explore", "account", "retrieve"}
	for index, argument := range arguments {
		arguments[index] = quote(argument)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, "script", "--quiet", "--return", "--command", "stty rows 24 cols 100; exec "+strings.Join(arguments, " "), "/dev/null")
	command.Cancel = func() error { return command.Process.Signal(syscall.SIGTERM) }
	command.WaitDelay = 5 * time.Second
	command.Env = append(cliEnvironment(t), "TERM=xterm-256color", "NO_COLOR=1")
	input, err := command.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	output := &explorerTerminal{input: input, steps: steps}
	var stderr bytes.Buffer
	command.Stdout, command.Stderr = output, &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("terminal explorer: %v\nstdout: %s\nstderr: %s", err, output.String(), stderr.String())
	}
	if len(output.steps) != 0 || !strings.HasSuffix(strings.TrimSpace(output.String()), want) {
		t.Fatalf("print selected value did not return %q: %q", want, output.String())
	}
}
