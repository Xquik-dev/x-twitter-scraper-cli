// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestYAMLToJSONEncodesWithEncodingJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input string
		want  string
	}{
		"whole float":        {`{"limit": 20.0}`, `{"limit":20}`},
		"escape character":   {`{"text": "\u001b[31m"}`, `{"text":"\u001b[31m"}`},
		"control character":  {"text: \"\\x01\"", `{"text":"\u0001"}`},
		"sorted keys":        {`{b: 1, a: 2}`, `{"a":2,"b":1}`},
		"scalars":            {`[1, 1.5, true, null, x, ~]`, `[1,1.5,true,null,"x",null]`},
		"large integer":      {`big: 18446744073709551615`, `{"big":18446744073709551615}`},
		"string tag":         {`!!str 123`, `"123"`},
		"string tag float":   {`!!str 1.10`, `"1.1"`},
		"string tag null":    {`!!str null`, `""`},
		"integer tag":        {`!!int "3"`, `3`},
		"integer tag float":  {`!!int 3.5`, `0`},
		"float tag":          {`!!float "2.5"`, `2.5`},
		"float tag text":     {`!!float abc`, `0`},
		"bool tag":           {`!!bool yes`, `true`},
		"null tag":           {`!!null x`, `null`},
		"timestamp tag":      {`!!timestamp 2001-12-14`, `"2001-12-14T00:00:00Z"`},
		"timestamp tag zero": {`!!timestamp 5`, `"0001-01-01T00:00:00Z"`},
		"binary tag":         {"!!binary |\n  aGVs\n  bG8=\n", `"aGVsbG8="`},
		"binary tag value":   {`v: !!binary aGk=`, `{"v":"aGk="}`},
		"directive document": {"%YAML 1.2\n---\ntext: hi\n", `{"text":"hi"}`},
		"tag over anchor":    {"x: !!int &a 5\ny: *a\n", `{"x":5,"y":5}`},
		"tag over alias":     {"x: &a 5\ny: !!str *a\n", `{"x":5,"y":"5"}`},
		"str over timestamp": {"a: &a !!timestamp 2001-12-14\nb: !!str *a\n", `{"a":"2001-12-14T00:00:00Z","b":"2001-12-14T00:00:00Z"}`},
		"str over binary":    {"a: &a !!binary aGk=\nb: !!str *a\n", `{"a":"aGk=","b":"aGk="}`},
		"int over map alias": {"a: &a {x: 1}\nb: !!int *a\n", `{"a":{"x":1},"b":0}`},
		"str over float tag": {"a: &a !!float 0.0000001\nb: !!str *a\n", `{"a":1e-7,"b":"1e-07"}`},
		"str over bool tag":  {"a: &a !!bool yes\nb: !!str *a\n", `{"a":true,"b":"true"}`},
		"str over null tag":  {"a: &a !!null x\nb: !!str *a\n", `{"a":null,"b":""}`},
		"str over int tag":   {"a: &a !!int 5\nb: !!str *a\n", `{"a":5,"b":"5"}`},
		"timestamp tag key":  {"!!timestamp 2001-12-14: y\n", `{"2001-12-14T00:00:00Z":"y"}`},
		"binary tag key":     {"!!binary aGk=: y\n", `{"aGk=":"y"}`},
		"self key alias":     {"&k key: *k\n", `{"key":"key"}`},
		"alias int key":      {"a: &k 5\n*k : x\n", `{"5":"x","a":5}`},
		"key anchor":         {"&k key: 1\nv: *k\n", `{"key":1,"v":"key"}`},
		"integer key anchor": {"&k 5: x\nv: *k\n", `{"5":"x","v":5}`},
		"explicit key":       {"? &k key\n: 1\nv: *k\n", `{"key":1,"v":"key"}`},
		"alias key":          {"a: &k key\n*k : 2\n", `{"a":"key","key":2}`},
		"tagged merge":       {"<<: !!map {a: 1}\n", `{"a":1}`},
		"tagged merge list":  {"<<: !!seq [{a: 1}]\n", `{"a":1}`},
		"other tag":          {`!custom [1]`, `[1]`},
		"block literal":      {"text: |\n  line\n", `{"text":"line\n"}`},
		"anchor and merge":   {"base: &base {a: 1, b: 2}\nchild:\n  <<: *base\n  b: 3\n", `{"base":{"a":1,"b":2},"child":{"a":1,"b":3}}`},
		"merge list":         {"one: &one {a: 1}\ntwo: &two {a: 2, b: 2}\nchild:\n  <<: [*one, *two]\n", `{"child":{"a":1,"b":2},"one":{"a":1},"two":{"a":2,"b":2}}`},
		"empty document":     {``, `null`},
		"comment only":       {"# nothing\n", `null`},
		"nested sequence":    {"- [a, {b: c}]\n", `[["a",{"b":"c"}]]`},
		"html in text":       {`text: "<b>&</b>"`, `{"text":"\u003cb\u003e\u0026\u003c/b\u003e"}`},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := YAMLToJSON([]byte(test.input))
			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestYAMLToJSONRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	for _, input := range []string{
		"{a: 1", "key: *missing", "<<: 1", "value: .inf", "!!bool maybe", "!!binary [1]",
		"<<: ~", "<<: !!null x", "n: &n ~\nc:\n  <<: *n", "a: &a {x: 1}\nb: !!str *a", "a: &a [1]\nb: !!str *a",
		"a: &a !custom {x: 1}\nb: !!str *a", "a: &k {x: 1}\n*k : y",
	} {
		_, err := YAMLToJSON([]byte(input))
		require.Error(t, err, input)
	}
}

func TestPipedNumbersAndTextReachTheirFlags(t *testing.T) {
	t.Parallel()

	command, err := run(t, []cli.Flag{
		&Int{Name: "limit", QueryPath: "limit"},
		&String{Name: "text", QueryPath: "text"},
	})
	require.NoError(t, err)
	piped, err := YAMLToJSON([]byte(`{"limit": 20.0, "text": "\u001b[31m"}`))
	require.NoError(t, err)
	data, ok := JSONObject(piped)
	require.True(t, ok)
	require.NoError(t, ApplyStdinDataToFlags(command, data))
	assert.Equal(t, "limit=20&text=%1B%5B31m", queryText(t, command))

	encoded, err := json.Marshal(data)
	require.NoError(t, err)
	assert.JSONEq(t, `{"limit":20,"text":"\u001b[31m"}`, string(encoded))
}
