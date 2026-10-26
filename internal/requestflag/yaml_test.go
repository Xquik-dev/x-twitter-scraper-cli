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
		"whole float":       {`{"limit": 20.0}`, `{"limit":20}`},
		"escape character":  {`{"text": "\u001b[31m"}`, `{"text":"\u001b[31m"}`},
		"control character": {"text: \"\\x01\"", `{"text":"\u0001"}`},
		"sorted keys":       {`{b: 1, a: 2}`, `{"a":2,"b":1}`},
		"scalars":           {`[1, 1.5, true, null, x, ~]`, `[1,1.5,true,null,"x",null]`},
		"large integer":     {`big: 18446744073709551615`, `{"big":18446744073709551615}`},
		"string tag":        {`!!str 123`, `"123"`},
		"block literal":     {"text: |\n  line\n", `{"text":"line\n"}`},
		"anchor and merge":  {"base: &base {a: 1, b: 2}\nchild:\n  <<: *base\n  b: 3\n", `{"base":{"a":1,"b":2},"child":{"a":1,"b":3}}`},
		"merge list":        {"one: &one {a: 1}\ntwo: &two {a: 2, b: 2}\nchild:\n  <<: [*one, *two]\n", `{"child":{"a":1,"b":2},"one":{"a":1},"two":{"a":2,"b":2}}`},
		"empty document":    {``, `null`},
		"comment only":      {"# nothing\n", `null`},
		"nested sequence":   {"- [a, {b: c}]\n", `[["a",{"b":"c"}]]`},
		"html in text":      {`text: "<b>&</b>"`, `{"text":"\u003cb\u003e\u0026\u003c/b\u003e"}`},
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

	for _, input := range []string{"{a: 1", "key: *missing", "<<: 1", "value: .inf"} {
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
