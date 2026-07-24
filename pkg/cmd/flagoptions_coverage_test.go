// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestRewriteAliases(t *testing.T) {
	t.Parallel()

	body := map[string]any{
		"alias":     "value",
		"canonical": "original",
		"same":      "unchanged",
	}
	rewriteAliases(body, "canonical", []string{"", "canonical", "missing", "alias"})

	assert.Equal(t, "value", body["canonical"])
	assert.Equal(t, "unchanged", body["same"])
	assert.NotContains(t, body, "alias")
}

func TestApplyDataAliases(t *testing.T) {
	t.Parallel()

	topLevel := &requestflag.Flag[string]{
		Name:        "top",
		BodyPath:    "canonical",
		DataAliases: []string{"alias"},
	}
	outer := &requestflag.Flag[map[string]any]{
		Name:     "outer",
		BodyPath: "parent",
	}
	inner := &requestflag.InnerFlag[string]{
		Name:        "outer.inner",
		OuterFlag:   outer,
		InnerField:  "inner_canonical",
		DataAliases: []string{"inner_alias"},
	}
	ignoredInner := &requestflag.InnerFlag[string]{
		Name:       "ignored.inner",
		OuterFlag:  &cli.StringFlag{Name: "ignored"},
		InnerField: "inner",
	}
	command := &cli.Command{Flags: []cli.Flag{
		topLevel,
		outer,
		inner,
		ignoredInner,
		&requestflag.Flag[string]{Name: "no-body"},
	}}
	body := map[string]any{
		"alias": "top value",
		"parent": map[string]any{
			"inner_alias": "inner value",
		},
	}

	applyDataAliases(command, body)

	assert.Equal(t, "top value", body["canonical"])
	assert.NotContains(t, body, "alias")
	nested := body["parent"].(map[string]any)
	assert.Equal(t, "inner value", nested["inner_canonical"])
	assert.NotContains(t, nested, "inner_alias")
}

func TestWrapFileInputValueVariants(t *testing.T) {
	t.Parallel()

	value, changed := wrapFileInputValue("")
	assert.False(t, changed)
	assert.Equal(t, "", value)

	value, changed = wrapFileInputValue("file.txt")
	assert.True(t, changed)
	assert.Equal(t, FilePathValue("file.txt"), value)

	value, changed = wrapFileInputValue([]string{"one.txt", "two.txt"})
	assert.True(t, changed)
	assert.Equal(t, []any{
		FilePathValue("one.txt"),
		FilePathValue("two.txt"),
	}, value)

	value, changed = wrapFileInputValue([]any{"one.txt", 2})
	assert.True(t, changed)
	assert.Equal(t, []any{FilePathValue("one.txt"), 2}, value)

	value, changed = wrapFileInputValue(42)
	assert.False(t, changed)
	assert.Equal(t, 42, value)
}

func TestWrapFileInputValuesByRequestLocation(t *testing.T) {
	t.Parallel()

	bodyFlag := &requestflag.Flag[string]{
		Name:      "body",
		BodyPath:  "body_file",
		FileInput: true,
	}
	queryFlag := &requestflag.Flag[string]{
		Name:      "query",
		QueryPath: "query_file",
		FileInput: true,
	}
	headerFlag := &requestflag.Flag[string]{
		Name:       "header",
		HeaderPath: "X-File",
		FileInput:  true,
	}
	bodyRoot := &requestflag.Flag[string]{
		Name:      "root",
		BodyRoot:  true,
		FileInput: true,
	}
	nonFile := &requestflag.Flag[string]{
		Name:      "plain",
		BodyPath:  "plain",
		FileInput: false,
	}
	pipedFile := &requestflag.Flag[string]{
		Name:      "piped",
		BodyPath:  "piped_file",
		FileInput: true,
	}
	require.NoError(t, bodyFlag.Set(bodyFlag.Name, "explicit-body.txt"))
	require.NoError(t, queryFlag.Set(queryFlag.Name, "query.txt"))
	require.NoError(t, headerFlag.Set(headerFlag.Name, "header.txt"))
	require.NoError(t, bodyRoot.Set(bodyRoot.Name, "root.txt"))
	require.NoError(t, nonFile.Set(nonFile.Name, "plain.txt"))

	body := map[string]any{
		"body_file":  "piped-body.txt",
		"piped_file": "piped-only.txt",
		"plain":      "plain.txt",
	}
	contents := &requestflag.RequestContents{
		Body:    body,
		Headers: map[string]any{"X-File": headerFlag.Get()},
		Queries: map[string]any{"query_file": queryFlag.Get()},
	}
	command := &cli.Command{Flags: []cli.Flag{
		bodyFlag,
		queryFlag,
		headerFlag,
		bodyRoot,
		nonFile,
		pipedFile,
		&cli.StringFlag{Name: "standard"},
	}}

	wrapFileInputValues(command, contents)

	assert.Equal(t, FilePathValue("explicit-body.txt"), body["body_file"])
	assert.Equal(t, FilePathValue("piped-only.txt"), body["piped_file"])
	assert.Equal(t, "plain.txt", body["plain"])
	assert.Equal(t, FilePathValue("query.txt"), contents.Queries["query_file"])
	assert.Equal(t, FilePathValue("header.txt"), contents.Headers["X-File"])
}
