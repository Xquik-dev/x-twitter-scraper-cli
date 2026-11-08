// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestRewriteAliases(t *testing.T) {
	t.Parallel()

	body := map[string]json.RawMessage{
		"alias":     json.RawMessage(`"value"`),
		"canonical": json.RawMessage(`"original"`),
		"same":      json.RawMessage(`"unchanged"`),
	}
	assert.True(t, rewriteAliases(body, "canonical", []string{"", "canonical", "missing", "alias"}))
	assert.False(t, rewriteAliases(body, "canonical", []string{"missing"}))

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.Equal(t, `{"canonical":"value","same":"unchanged"}`, string(encoded))
}

func TestApplyDataAliases(t *testing.T) {
	t.Parallel()

	topLevel := &requestflag.String{
		Name:        "top",
		BodyPath:    "canonical",
		DataAliases: []string{"alias"},
	}
	outer := &requestflag.Object{
		Name:     "outer",
		BodyPath: "parent",
	}
	inner := &requestflag.InnerString{
		Name:        "outer.inner",
		OuterFlag:   outer,
		InnerField:  "inner_canonical",
		DataAliases: []string{"inner_alias"},
	}
	ignoredInner := &requestflag.InnerString{
		Name:       "ignored.inner",
		OuterFlag:  &cli.StringFlag{Name: "ignored"},
		InnerField: "inner",
	}
	command := &cli.Command{Flags: []cli.Flag{
		topLevel,
		outer,
		inner,
		ignoredInner,
		&requestflag.String{Name: "no-body"},
	}}
	body, ok := requestflag.JSONObject(json.RawMessage(`{"alias":"top value","parent":{"inner_alias":"inner value","kept":1}}`))
	require.True(t, ok)

	require.NoError(t, applyDataAliases(command, body))

	encoded, err := json.Marshal(body)
	require.NoError(t, err)
	assert.Equal(t, `{"canonical":"top value","parent":{"inner_canonical":"inner value","kept":1}}`, string(encoded))
}
