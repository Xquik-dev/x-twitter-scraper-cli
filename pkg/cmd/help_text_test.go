// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

var proseColon = regexp.MustCompile(`[A-Za-z0-9)\]'"]:(\s|$)`)

var allowedColons = []string{"Authorization: Bearer"}

func punctuationProblem(text string) string {
	for _, allowed := range allowedColons {
		text = strings.ReplaceAll(text, allowed, "")
	}
	if strings.Contains(text, ";") {
		return "a semicolon"
	}
	if proseColon.MatchString(text) {
		return "a colon"
	}
	return ""
}

func TestHelpTextUsesNoColonOrSemicolonInProse(t *testing.T) {
	var problems []string
	var walk func(command *cli.Command, path string)
	walk = func(command *cli.Command, path string) {
		texts := map[string]string{
			"usage":       command.Usage,
			"description": command.Description,
			"usage text":  command.UsageText,
		}
		for _, flag := range command.Flags {
			if documented, ok := flag.(cli.DocGenerationFlag); ok {
				texts["flag --"+flag.Names()[0]] = documented.GetUsage()
			}
		}
		for where, text := range texts {
			if problem := punctuationProblem(text); problem != "" {
				problems = append(problems, path+" "+where+" holds "+problem+" in "+strings.TrimSpace(text))
			}
		}
		for _, child := range command.Commands {
			walk(child, path+" "+child.Name)
		}
	}
	walk(newCommand(), "x-twitter-scraper")
	require.Empty(t, problems)
}

func runForError(t *testing.T, args ...string) error {
	t.Helper()
	command := newCommand()
	command.Writer = io.Discard
	command.ErrWriter = io.Discard
	return command.Run(context.Background(), append([]string{"x-twitter-scraper"}, args...))
}

func TestUserFacingErrorsSayWhatToDoWithoutAColon(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		message string
	}{
		{
			name:    "extra arguments",
			args:    []string{"--api-key", "test-key", "--base-url", "http://127.0.0.1:9", "account", "retrieve", "extra"},
			message: "Unexpected extra arguments [extra]. Remove them, then run the command again.",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := runForError(t, test.args...)
			require.EqualError(t, err, test.message)
		})
	}
}

func TestShellCompletionErrorsListTheShells(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		message string
	}{
		{
			name:    "missing shell",
			args:    []string{"@completion"},
			message: "Shell is missing. Choose bash, fish, pwsh or zsh.\n",
		},
		{
			name:    "unknown shell",
			args:    []string{"@completion", "tcsh"},
			message: "Shell \"tcsh\" is unsupported. Choose bash, fish, pwsh or zsh.\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, code := runCLIProcess(t, "", test.args...)
			require.Equal(t, 1, code)
			require.Equal(t, test.message, output)
		})
	}
}

func TestUnknownFormatsListTheChoices(t *testing.T) {
	for _, flag := range []string{"--format", "--format-error"} {
		t.Run(flag, func(t *testing.T) {
			err := runForError(t, flag, "csv", "account", "retrieve")
			require.Error(t, err)
			require.True(t, strings.HasSuffix(err.Error(), " Choose auto, explore, json, jsonl, pretty, raw or yaml."), err.Error())
		})
	}
}
