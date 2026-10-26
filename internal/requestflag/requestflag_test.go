// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func run(t *testing.T, flags []cli.Flag, args ...string) (*cli.Command, error) {
	t.Helper()
	command := &cli.Command{
		Name:   "probe",
		Flags:  flags,
		Action: func(context.Context, *cli.Command) error { return nil },
	}
	err := command.Run(context.Background(), append([]string{"probe"}, args...))
	return command, err
}

func contents(t *testing.T, command *cli.Command) RequestContents {
	t.Helper()
	extracted, err := ExtractRequestContents(command)
	require.NoError(t, err)
	return extracted
}

func bodyJSON(t *testing.T, command *cli.Command) string {
	t.Helper()
	encoded, err := json.Marshal(contents(t, command).Body)
	require.NoError(t, err)
	return string(encoded)
}

func queryText(t *testing.T, command *cli.Command) string {
	t.Helper()
	return apiquery.Encode(contents(t, command).Queries, apiquery.ArrayQueryFormatComma).Encode()
}

func headerText(t *testing.T, command *cli.Command) string {
	t.Helper()
	return apiquery.Encode(contents(t, command).Headers, apiquery.ArrayQueryFormatRepeat).Encode()
}

func displayed(t *testing.T, flag interface {
	PreParse() error
	GetValue() string
}) string {
	t.Helper()
	require.NoError(t, flag.PreParse())
	return flag.GetValue()
}

func TestScalarFlagsParseTheirValues(t *testing.T) {
	t.Parallel()

	command, err := run(t, []cli.Flag{
		&String{Name: "text", BodyPath: "text"},
		&Int{Name: "count", BodyPath: "count"},
		&Bool{Name: "safe", BodyPath: "safe"},
		&Bool{Name: "quiet", BodyPath: "quiet"},
		&StringList{Name: "tag", BodyPath: "tags"},
	}, "--text", "hello", "--count", "0x10", "--safe=false", "--quiet", "--tag", "a,b", "--tag", "c")
	require.NoError(t, err)
	assert.Equal(t, `{"count":16,"quiet":true,"safe":false,"tags":["a,b","c"],"text":"hello"}`, bodyJSON(t, command))
	assert.Equal(t, "hello", command.String("text"))
}

func TestInvalidScalarsFail(t *testing.T) {
	t.Parallel()

	_, err := run(t, []cli.Flag{&Int{Name: "count"}}, "--count", "many")
	require.ErrorContains(t, err, "invalid syntax")
	_, err = run(t, []cli.Flag{&Bool{Name: "safe"}}, "--safe=maybe")
	require.Error(t, err)
}

func TestValidatorRunsOnDefaultAndValue(t *testing.T) {
	t.Parallel()

	tooBig := errors.New("count is too big")
	validator := func(count int64) error {
		if count > 10 {
			return tooBig
		}
		return nil
	}
	_, err := run(t, []cli.Flag{&Int{Name: "count", Validator: validator}}, "--count", "11")
	require.ErrorContains(t, err, tooBig.Error())
	_, err = run(t, []cli.Flag{&Int{Name: "count", Default: 20, Validator: validator}})
	require.ErrorContains(t, err, tooBig.Error())
	_, err = run(t, []cli.Flag{&Int{Name: "count", Validator: validator}}, "--count", "3")
	require.NoError(t, err)
}

func TestListFlagReplacesItsDefault(t *testing.T) {
	t.Parallel()

	command, err := run(t, []cli.Flag{
		&StringList{Name: "tag", BodyPath: "tags", Default: []string{"default"}},
	}, "--tag", "b", "--tag", "c")
	require.NoError(t, err)
	assert.Equal(t, `{"tags":["b","c"]}`, bodyJSON(t, command))
}

func TestJSONFlagsReadYAML(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input string
		want  string
	}{
		"object":            {`{a: 1, b: [x, "y"]}`, `{"a":1,"b":["x","y"]}`},
		"list":              {`[1, 2]`, `[1,2]`},
		"number":            {`1700000000`, `1700000000`},
		"quoted text":       {`'2019-12-27'`, `"2019-12-27"`},
		"null":              {`null`, `null`},
		"file reference":    {`@file.txt`, `"@file.txt"`},
		"identifier":        {`abc_def-1.2=3`, `"abc_def-1.2=3"`},
		"text with a colon": {`key: value`, `{"key":"value"}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			command, err := run(t, []cli.Flag{&JSON{Name: "value", BodyPath: "value"}}, "--value", test.input)
			require.NoError(t, err)
			assert.Equal(t, `{"value":`+test.want+`}`, bodyJSON(t, command))
		})
	}
}

func TestJSONFlagRejectsInvalidYAML(t *testing.T) {
	t.Parallel()

	_, err := run(t, []cli.Flag{&JSON{Name: "value"}}, "--value", "{a: 1")
	require.ErrorContains(t, err, "failed to parse as YAML")
	_, err = run(t, []cli.Flag{&Object{Name: "value"}}, "--value", "[1]")
	require.ErrorContains(t, err, "to an object")
	_, err = run(t, []cli.Flag{&ObjectList{Name: "value"}}, "--value", "{a: 1")
	require.ErrorContains(t, err, "failed to parse as YAML")
	_, err = run(t, []cli.Flag{&JSONList{Name: "value"}}, "--value", "{a: 1")
	require.ErrorContains(t, err, "failed to parse as YAML")
}

func TestJSONListAndObjectFlags(t *testing.T) {
	t.Parallel()

	command, err := run(t, []cli.Flag{
		&JSONList{Name: "target", BodyPath: "targets"},
		&Object{Name: "meta", BodyPath: "meta"},
		&ObjectList{Name: "row", BodyPath: "rows"},
	}, "--target", "{kind: tweet, value: '1'}", "--target", "NASA", "--meta", "{a: 1}", "--meta", "{b: 2}",
		"--row", "{a: 1}", "--row", "null")
	require.NoError(t, err)
	assert.Equal(t, `{"meta":{"b":2},"rows":[{"a":1},null],"targets":[{"kind":"tweet","value":"1"},"NASA"]}`, bodyJSON(t, command))
}

func objectListCommand() cli.Command {
	return WithInnerFlags(cli.Command{
		Name: "probe",
		Flags: []cli.Flag{
			&ObjectList{Name: "relation-target", BodyPath: "relationTargets"},
			&Object{Name: "parent", BodyPath: "parent"},
		},
		Action: func(context.Context, *cli.Command) error { return nil },
	}, map[string][]HasOuterFlag{
		"relation-target": {
			&InnerString{Name: "relation-target.relation", InnerField: "relation"},
			&InnerString{Name: "relation-target.value", InnerField: "value"},
		},
		"parent": {
			&InnerString{Name: "parent.name", InnerField: "name"},
		},
	})
}

func TestInnerFlagsFillTheTrailingObject(t *testing.T) {
	t.Parallel()

	command := objectListCommand()
	require.NoError(t, CheckInnerFlags(command))
	require.NoError(t, command.Run(context.Background(), []string{
		"probe",
		"--relation-target.relation", "followers",
		"--relation-target.value", "NASA",
		"--relation-target.relation", "following",
		"--parent", "{id: 7}",
		"--parent.name", "root",
	}))
	assert.Equal(t,
		`{"parent":{"id":7,"name":"root"},"relationTargets":[{"relation":"followers","value":"NASA"},{"relation":"following"}]}`,
		bodyJSON(t, &command))
}

func TestInnerFlagAppendsToAnObjectList(t *testing.T) {
	t.Parallel()

	command := objectListCommand()
	require.NoError(t, command.Run(context.Background(), []string{
		"probe",
		"--relation-target", "{relation: followers, value: NASA}",
		"--relation-target.relation", "following",
		"--relation-target.value", "SpaceX",
	}))
	assert.Equal(t,
		`{"relationTargets":[{"relation":"followers","value":"NASA"},{"relation":"following","value":"SpaceX"}]}`,
		bodyJSON(t, &command))
}

func TestInnerFlagValidatorAndOuterChecks(t *testing.T) {
	t.Parallel()

	rejected := errors.New("rejected")
	outer := &ObjectList{Name: "row", BodyPath: "rows"}
	inner := &InnerString{Name: "row.name", InnerField: "name", OuterFlag: outer, Validator: func(string) error { return rejected }}
	require.ErrorIs(t, inner.Set("row.name", "x"), rejected)

	scalar := &String{Name: "text"}
	require.ErrorContains(t, (&InnerString{Name: "text.name", InnerField: "name", OuterFlag: scalar}).Set("text.name", "x"), "not an object")
	require.ErrorContains(t, (&InnerString{Name: "orphan.name", InnerField: "name"}).Set("orphan.name", "x"), "cannot set an inner field")

	require.ErrorContains(t, CheckInnerFlags(cli.Command{Flags: []cli.Flag{&InnerString{Name: "orphan.name"}}}), "missing an outer flag")
	require.ErrorContains(t, CheckInnerFlags(cli.Command{Flags: []cli.Flag{&InnerString{Name: "other.name", OuterFlag: outer}}}), "must start with")
	assert.Panics(t, func() {
		WithInnerFlags(cli.Command{}, map[string][]HasOuterFlag{"missing": {&InnerString{Name: "missing.name"}}})
	})
	unchanged := cli.Command{Name: "same"}
	assert.Equal(t, "same", WithInnerFlags(unchanged, nil).Name)
}

func TestRequestLocations(t *testing.T) {
	t.Parallel()

	command, err := run(t, []cli.Flag{
		&String{Name: "query", QueryPath: "q"},
		&StringList{Name: "field", QueryPath: "fields"},
		&Int{Name: "limit", QueryPath: "limit"},
		&Bool{Name: "safe", QueryPath: "safe"},
		&String{Name: "key", HeaderPath: "Idempotency-Key"},
		&StringList{Name: "trace", HeaderPath: "X-Trace"},
		&String{Name: "id", PathParam: "id"},
		&String{Name: "step", BodyPath: "step", Default: "compose", Const: true},
		&String{Name: "unset", BodyPath: "unset"},
	}, "--query", "from:nasa", "--field", "a", "--field", "b", "--limit", "5", "--safe",
		"--key", "k1", "--trace", "t1", "--trace", "t2", "--id", "42")
	require.NoError(t, err)
	assert.Equal(t, "fields=a%2Cb&limit=5&q=from%3Anasa&safe=true", queryText(t, command))
	assert.Equal(t, "Idempotency-Key=k1&X-Trace=t1&X-Trace=t2", headerText(t, command))
	assert.Equal(t, `{"step":"compose"}`, bodyJSON(t, command))
	assert.Equal(t, "42", command.String("id"))
}

func TestApplyStdinDataToFlags(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		flags  []cli.Flag
		args   []string
		data   string
		query  string
		header string
	}{
		"query path":             {[]cli.Flag{&String{Name: "account-id", QueryPath: "account_id"}}, nil, `{"account_id":"acct_123"}`, "account_id=acct_123", ""},
		"header path":            {[]cli.Flag{&String{Name: "key", HeaderPath: "Idempotency-Key"}}, nil, `{"Idempotency-Key":"key-xyz"}`, "", "Idempotency-Key=key-xyz"},
		"body path untouched":    {[]cli.Flag{&String{Name: "name", BodyPath: "name"}}, nil, `{"name":"x"}`, "", ""},
		"flag set on the line":   {[]cli.Flag{&String{Name: "account-id", QueryPath: "account_id"}}, []string{"--account-id", "cli"}, `{"account_id":"piped"}`, "account_id=cli", ""},
		"integer":                {[]cli.Flag{&Int{Name: "limit", QueryPath: "limit"}}, nil, `{"limit":50}`, "limit=50", ""},
		"boolean":                {[]cli.Flag{&Bool{Name: "safe", QueryPath: "safe"}}, nil, `{"safe":true}`, "safe=true", ""},
		"object as text":         {[]cli.Flag{&String{Name: "filter", QueryPath: "filter"}}, nil, `{"filter":{"a":[1, 2]}}`, "filter=%7B%22a%22%3A%5B1%2C2%5D%7D", ""},
		"data alias":             {[]cli.Flag{&String{Name: "account-id", QueryPath: "account_id", DataAliases: []string{"accountId"}}}, nil, `{"accountId":"acct_456"}`, "account_id=acct_456", ""},
		"canonical before alias": {[]cli.Flag{&String{Name: "account-id", QueryPath: "account_id", DataAliases: []string{"accountId"}}}, nil, `{"accountId":"alias","account_id":"canonical"}`, "account_id=canonical", ""},
		"no matching key":        {[]cli.Flag{&String{Name: "account-id", QueryPath: "account_id"}}, nil, `{"other":"x"}`, "", ""},
		"no path":                {[]cli.Flag{&String{Name: "loose"}}, nil, `{"loose":"x"}`, "", ""},
		"empty data":             {[]cli.Flag{&String{Name: "account-id", QueryPath: "account_id"}}, nil, `{}`, "", ""},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			command, err := run(t, test.flags, test.args...)
			require.NoError(t, err)
			data, ok := JSONObject(json.RawMessage(test.data))
			require.True(t, ok)
			require.NoError(t, ApplyStdinDataToFlags(command, data))
			assert.Equal(t, test.query, queryText(t, command))
			assert.Equal(t, test.header, headerText(t, command))
		})
	}
}

func TestApplyStdinDataToInnerFlags(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		data string
		want string
	}{
		"nested field":      {`{"parent":{"name":"San Francisco"}}`, `{"parent":{"name":"San Francisco"}}`},
		"nested data alias": {`{"parent":{"label":"Oakland"}}`, `{"parent":{"name":"Oakland"}}`},
		"no nested object":  {`{"parent":"flat"}`, `{}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			outer := &Object{Name: "parent", BodyPath: "parent"}
			inner := &InnerString{Name: "parent.name", InnerField: "name", DataAliases: []string{"label"}, OuterFlag: outer}
			command := &cli.Command{Flags: []cli.Flag{outer, inner}}
			data, ok := JSONObject(json.RawMessage(test.data))
			require.True(t, ok)
			require.NoError(t, ApplyStdinDataToFlags(command, data))
			assert.Equal(t, test.want, bodyJSON(t, command))
		})
	}

	outer := &Object{Name: "parent"}
	inner := &InnerString{Name: "parent.name", InnerField: "name", OuterFlag: outer}
	command := &cli.Command{Flags: []cli.Flag{outer, inner}}
	data, ok := JSONObject(json.RawMessage(`{"parent":{"name":"x"}}`))
	require.True(t, ok)
	require.NoError(t, ApplyStdinDataToFlags(command, data))
	assert.False(t, outer.IsSet())
}

func TestApplyStdinDataReportsBadValues(t *testing.T) {
	t.Parallel()

	command := &cli.Command{Flags: []cli.Flag{&Int{Name: "limit", QueryPath: "limit"}}}
	require.ErrorContains(t, ApplyStdinDataToFlags(command, map[string]json.RawMessage{"limit": json.RawMessage(`"many"`)}), `cannot set flag "limit"`)
	require.ErrorContains(t, ApplyStdinDataToFlags(command, map[string]json.RawMessage{"limit": json.RawMessage(`{`)}), `cannot format piped value`)
}

func TestRequiredFlags(t *testing.T) {
	t.Parallel()

	bodyFlag := &String{Name: "text", Required: true, BodyPath: "text"}
	plainFlag := &String{Name: "mode", Required: true}
	constFlag := &String{Name: "step", Required: true, Const: true, BodyPath: "step"}
	assert.False(t, bodyFlag.IsRequired())
	assert.True(t, bodyFlag.IsRequiredAsFlagOrStdin())
	assert.True(t, plainFlag.IsRequired())
	assert.False(t, constFlag.IsRequired())
	assert.False(t, constFlag.IsRequiredAsFlagOrStdin())

	command := &cli.Command{Flags: []cli.Flag{bodyFlag, plainFlag, constFlag, &String{Name: "free"}}}
	names := func(flags []cli.Flag) []string {
		result := []string{}
		for _, flag := range flags {
			result = append(result, flag.Names()[0])
		}
		return result
	}
	assert.Equal(t, []string{"text", "mode"}, names(GetMissingRequiredFlags(command, nil)))
	assert.Equal(t, []string{"mode"}, names(GetMissingRequiredFlags(command, map[string]json.RawMessage{"text": json.RawMessage(`"x"`)})))
}

func TestEnvironmentSources(t *testing.T) {
	t.Setenv("REQUESTFLAG_TEXT", "")
	t.Setenv("REQUESTFLAG_SAFE", "")
	t.Setenv("REQUESTFLAG_LIMIT", "")
	t.Setenv("REQUESTFLAG_BAD", "many")
	t.Setenv("REQUESTFLAG_KEY", "secret")

	command, err := run(t, []cli.Flag{
		&String{Name: "text", BodyPath: "text", Sources: cli.EnvVars("REQUESTFLAG_TEXT")},
		&Bool{Name: "safe", BodyPath: "safe", Default: true, Sources: cli.EnvVars("REQUESTFLAG_SAFE")},
		&Int{Name: "limit", BodyPath: "limit", Default: 20, Sources: cli.EnvVars("REQUESTFLAG_LIMIT")},
		&String{Name: "key", Sources: cli.EnvVars("REQUESTFLAG_KEY")},
		&String{Name: "cli", BodyPath: "cli", Sources: cli.EnvVars("REQUESTFLAG_KEY")},
	}, "--cli", "line")
	require.NoError(t, err)
	assert.Equal(t, `{"cli":"line","limit":20,"safe":false,"text":""}`, bodyJSON(t, command))
	assert.Equal(t, "secret", command.String("key"))

	_, err = run(t, []cli.Flag{&Int{Name: "bad", Sources: cli.EnvVars("REQUESTFLAG_BAD")}})
	require.ErrorContains(t, err, `could not parse "many" as int value`)
}

func TestHelpTraits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		flag       cli.DocGenerationMultiValueFlag
		typeName   string
		takesValue bool
		multiValue bool
	}{
		{&String{}, "string", true, false},
		{&Int{}, "int", true, false},
		{&Bool{}, "boolean", false, false},
		{&StringList{}, "string", true, true},
		{&JSON{}, "", true, false},
		{&JSONList{}, "any", true, true},
		{&Object{}, "string=any", true, true},
		{&ObjectList{}, "any", true, true},
		{&InnerString{}, "string", true, false},
	}
	for _, test := range tests {
		assert.Equal(t, test.typeName, test.flag.TypeName())
		assert.Equal(t, test.takesValue, test.flag.TakesValue())
		assert.Equal(t, test.multiValue, test.flag.IsMultiValueFlag())
	}
	assert.True(t, (&Bool{}).IsBoolFlag())
	assert.False(t, (&InnerString{}).IsBoolFlag())
}

func TestFlagMetadata(t *testing.T) {
	t.Parallel()

	flag := &Int{
		Name:        "limit",
		Aliases:     []string{"l"},
		Category:    "paging",
		DefaultText: "twenty",
		HideDefault: true,
		Usage:       "How many to return.",
		Hidden:      true,
		Default:     20,
		Sources:     cli.EnvVars("REQUESTFLAG_LIMIT_HELP"),
		QueryPath:   "limit",
		DataAliases: []string{"max"},
	}
	assert.Equal(t, "", flag.GetValue())
	require.NoError(t, flag.PreParse())
	assert.Equal(t, "20", flag.GetValue())
	assert.Equal(t, []string{"limit", "l"}, flag.Names())
	assert.Equal(t, "paging", flag.GetCategory())
	flag.SetCategory("other")
	assert.Equal(t, "other", flag.GetCategory())
	assert.Equal(t, "twenty", flag.GetDefaultText())
	assert.False(t, flag.IsDefaultVisible())
	assert.False(t, flag.IsVisible())
	assert.Equal(t, "How many to return.", flag.GetUsage())
	assert.Equal(t, []string{"REQUESTFLAG_LIMIT_HELP"}, flag.GetEnvVars())
	assert.Equal(t, "limit", flag.GetQueryPath())
	assert.Equal(t, []string{"max"}, flag.GetDataAliases())
	assert.True(t, flag.IsLocal())
	assert.Equal(t, 0, flag.Count())
	require.NoError(t, flag.Set("limit", "3"))
	assert.Equal(t, 1, flag.Count())
	assert.Contains(t, flag.String(), "--limit")

	assert.Equal(t, "[a b]", displayed(t, &StringList{Default: []string{"a", "b"}, Name: "tag"}))
	assert.Equal(t, `{"a":1}`, displayed(t, &Object{Default: map[string]json.RawMessage{"a": json.RawMessage(`1`)}, Name: "meta"}))
	assert.Equal(t, `[{"a":1}]`, displayed(t, &ObjectList{Default: []map[string]json.RawMessage{{"a": json.RawMessage(`1`)}}, Name: "rows"}))
	assert.Equal(t, `["x"]`, displayed(t, &JSONList{Default: []json.RawMessage{json.RawMessage(`"x"`)}, Name: "targets"}))
	assert.Equal(t, `{"a":1}`, displayed(t, &JSON{Default: json.RawMessage(`{"a":1}`), Name: "value"}))
	assert.Equal(t, "true", displayed(t, &Bool{Default: true, Name: "safe"}))

	inner := &InnerString{Name: "row.name", Aliases: []string{"n"}, Usage: "Name.", DefaultText: "none", DataAliases: []string{"label"}}
	require.NoError(t, inner.PreParse())
	require.NoError(t, inner.PostParse())
	assert.False(t, inner.IsSet())
	assert.Equal(t, []string{"row.name", "n"}, inner.Names())
	assert.Equal(t, "Name.", inner.GetUsage())
	assert.Equal(t, "none", inner.GetDefaultText())
	assert.Equal(t, "", inner.GetValue())
	assert.Nil(t, inner.GetEnvVars())
	assert.False(t, inner.IsDefaultVisible())
	assert.Contains(t, inner.String(), "--row.name")
}
