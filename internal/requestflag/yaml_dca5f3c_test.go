// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type yamlBody struct {
	input string
	want  string
}

func assertYAMLBodies(t *testing.T, bodies map[string]yamlBody) {
	t.Helper()
	for name, body := range bodies {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := YAMLToJSON([]byte(body.input))
			require.NoError(t, err)
			assert.Equal(t, body.want, string(got))
		})
	}
}

func TestYAMLToJSONSendsTheBodiesDca5f3cSent(t *testing.T) {
	t.Parallel()

	assertYAMLBodies(t, map[string]yamlBody{
		"key anchor":                      {"&k key: 1\nv: *k", `{"key":1,"v":"key"}`},
		"alias key":                       {"a: &k key\n*k : 2", `{"a":"key","key":2}`},
		"nested key anchor":               {"m:\n  &k key: 1\n  v: *k", `{"m":{"key":1,"v":"key"}}`},
		"int key anchor":                  {"&k 5: x\nv: *k", `{"5":"x","v":5}`},
		"explicit key anchor":             {"? &k key\n: 1\nv: *k", `{"key":1,"v":"key"}`},
		"key anchor merge value":          {"&k key: {a: 1}\nv: *k", `{"key":{"a":1},"v":"key"}`},
		"merge empty list":                {"<<: []", `{}`},
		"merge empty map":                 {"<<: {}", `{}`},
		"str over date anchor":            {"a: &a 2001-12-14\nb: !!str *a", `{"a":"2001-12-14","b":"2001-12-14"}`},
		"str over float anchor":           {"a: &a 1.10\nb: !!str *a", `{"a":1.1,"b":"1.1"}`},
		"str over bool anchor":            {"a: &a true\nb: !!str *a", `{"a":true,"b":"true"}`},
		"str over null anchor":            {"a: &a ~\nb: !!str *a", `{"a":null,"b":""}`},
		"str over text anchor":            {"a: &a hi\nb: !!str *a", `{"a":"hi","b":"hi"}`},
		"str over int tag anchor":         {"a: &a !!int 5\nb: !!str *a", `{"a":5,"b":"5"}`},
		"str over str tag anchor":         {"a: &a !!str 5\nb: !!str *a", `{"a":"5","b":"5"}`},
		"str over date":                   {"b: !!str 2001-12-14", `{"b":"2001-12-14"}`},
		"int over str tag anchor":         {"a: &a !!str 5\nb: !!int *a", `{"a":"5","b":5}`},
		"int over map anchor":             {"a: &a {x: 1}\nb: !!int *a", `{"a":{"x":1},"b":0}`},
		"float over map anchor":           {"a: &a {x: 1}\nb: !!float *a", `{"a":{"x":1},"b":0}`},
		"null over map anchor":            {"a: &a {x: 1}\nb: !!null *a", `{"a":{"x":1},"b":null}`},
		"timestamp over map anchor":       {"a: &a {x: 1}\nb: !!timestamp *a", `{"a":{"x":1},"b":"0001-01-01T00:00:00Z"}`},
		"custom over map anchor":          {"a: &a {x: 1}\nb: !custom *a", `{"a":{"x":1},"b":{"x":1}}`},
		"binary bad padding":              {"v: !!binary aGk", `{"v":""}`},
		"binary bad text":                 {"v: !!binary '!!!'", `{"v":""}`},
		"timestamp":                       {"v: !!timestamp 2001-12-14", `{"v":"2001-12-14T00:00:00Z"}`},
		"timestamp alias":                 {"v: &a !!timestamp 2001-12-14\nw: *a", `{"v":"2001-12-14T00:00:00Z","w":"2001-12-14T00:00:00Z"}`},
		"date":                            {"v: 2001-12-14", `{"v":"2001-12-14"}`},
		"timestamp zone":                  {"v: !!timestamp 2001-12-14T21:59:43.10-05:00", `{"v":"2001-12-14T21:59:43.1-05:00"}`},
		"int key":                         {"1: a", `{"1":"a"}`},
		"bool key":                        {"true: a", `{"true":"a"}`},
		"float key":                       {"1.5: a", `{"1.5":"a"}`},
		"anchor on map":                   {"a: &a {x: 1}\nb: *a", `{"a":{"x":1},"b":{"x":1}}`},
		"anchor on seq":                   {"a: &a [1, x]\nb: *a", `{"a":[1,"x"],"b":[1,"x"]}`},
		"tag over anchor":                 {"x: !!int &a 5\ny: *a", `{"x":5,"y":5}`},
		"anchor over tag":                 {"x: &a !!int \"5\"\ny: *a", `{"x":5,"y":5}`},
		"redefined anchor":                {"a: &a 1\nb: &a 2\nc: *a", `{"a":1,"b":2,"c":2}`},
		"merge override":                  {"base: &base {a: 1, b: 2}\nchild:\n  <<: *base\n  b: 3", `{"base":{"a":1,"b":2},"child":{"a":1,"b":3}}`},
		"int over seq anchor":             {"a: &a [1]\nb: !!int *a", `{"a":[1],"b":0}`},
		"int over timestamp anchor":       {"a: &a !!timestamp 2001-12-14\nb: !!int *a", `{"a":"2001-12-14T00:00:00Z","b":0}`},
		"int over bool anchor":            {"a: &a true\nb: !!int *a", `{"a":true,"b":0}`},
		"int over null anchor":            {"a: &a ~\nb: !!int *a", `{"a":null,"b":0}`},
		"int over float anchor":           {"a: &a 2.5\nb: !!int *a", `{"a":2.5,"b":0}`},
		"int over text anchor":            {"a: &a x\nb: !!int *a", `{"a":"x","b":0}`},
		"int over quoted int anchor":      {"a: &a \"7\"\nb: !!int *a", `{"a":"7","b":7}`},
		"str over float tag anchor":       {"a: &a !!float 0.0000001\nb: !!str *a", `{"a":1e-7,"b":"1e-07"}`},
		"str over big float tag anchor":   {"a: &a !!float 1e21\nb: !!str *a", `{"a":1e+21,"b":"1e+21"}`},
		"str over bool tag anchor":        {"a: &a !!bool yes\nb: !!str *a", `{"a":true,"b":"true"}`},
		"str over null tag anchor":        {"a: &a !!null x\nb: !!str *a", `{"a":null,"b":""}`},
		"str over custom tag anchor":      {"a: &a !custom hi\nb: !!str *a", `{"a":"hi","b":"hi"}`},
		"str over int tag float anchor":   {"a: &a !!int 3.5\nb: !!str *a", `{"a":0,"b":"0"}`},
		"str over float tag text anchor":  {"a: &a !!float \"2.5\"\nb: !!str *a", `{"a":2.5,"b":"2.5"}`},
		"key anchor self value":           {"&k key: *k", `{"key":"key"}`},
		"alias key to int":                {"a: &k 5\n*k : x", `{"5":"x","a":5}`},
		"tagged key":                      {"!!str 5: x", `{"5":"x"}`},
		"timestamp over text anchor":      {"a: &a 2001-12-14\nb: !!timestamp *a", `{"a":"2001-12-14","b":"2001-12-14T00:00:00Z"}`},
		"bool over text anchor":           {"a: &a yes\nb: !!bool *a", `{"a":"yes","b":true}`},
		"float over text anchor":          {"a: &a \"2.5\"\nb: !!float *a", `{"a":"2.5","b":2.5}`},
		"float tag key":                   {"!!float 1: y", `{"1":"y"}`},
		"str over small float tag anchor": {"a: &a !!float 0.00001\nb: !!str *a", `{"a":0.00001,"b":"1e-05"}`},
		"float over float tag anchor":     {"a: &a !!float 2.5\nb: !!float *a", `{"a":2.5,"b":2.5}`},
		"timestamp over timestamp anchor": {"a: &a !!timestamp 2001-12-14\nb: !!timestamp *a", `{"a":"2001-12-14T00:00:00Z","b":"2001-12-14T00:00:00Z"}`},
		"str over float tag 1e6":          {"b: !!str !!float 1000000", `{"b":"1e+06"}`},
		"float tag key 1e6":               {"!!float 1000000: y", `{"1e+06":"y"}`},
		"float over int tag anchor":       {"a: &a !!int 5\nb: !!float *a", `{"a":5,"b":5}`},
		"float over str tag anchor":       {"a: &a !!str \"2.5\"\nb: !!float *a", `{"a":"2.5","b":2.5}`},
		"bool over str tag anchor":        {"a: &a !!str yes\nb: !!bool *a", `{"a":"yes","b":true}`},
		"timestamp over str tag anchor":   {"a: &a !!str 2001-12-14\nb: !!timestamp *a", `{"a":"2001-12-14","b":"2001-12-14T00:00:00Z"}`},
		"float over custom anchor":        {"a: &a !custom 2.5\nb: !!float *a", `{"a":"2.5","b":2.5}`},
		"str over custom float":           {"b: !!str !custom 1.10", `{"b":"1.10"}`},
		"float over str tag":              {"b: !!float !!str 2.5", `{"b":2.5}`},
		"int over str tag":                {"b: !!int !!str 7", `{"b":7}`},
		"str over float tag 1e6 anchor":   {"a: &a !!float 1000000\nb: !!str *a", `{"a":1000000,"b":"1e+06"}`},
		"custom float key":                {"!custom 1.10: y", `{"1.10":"y"}`},
		"timestamp over str tag":          {"b: !!timestamp !!str 2001-12-14", `{"b":"2001-12-14T00:00:00Z"}`},
		"int over float tag anchor":       {"a: &a !!float 7\nb: !!int *a", `{"a":7,"b":7}`},
		"bool over bool tag anchor":       {"a: &a !!bool yes\nb: !!bool *a", `{"a":true,"b":true}`},
		"float over timestamp anchor":     {"a: &a !!timestamp 2001-12-14\nb: !!float *a", `{"a":"2001-12-14T00:00:00Z","b":0}`},
		"binary over custom int anchor":   {"a: &a !custom 5\nb: !!binary *a", `{"a":"5","b":""}`},
		"str over inf":                    {"b: !!str !!float .inf", `{"b":"+Inf"}`},
		"directive int":                   {"%TAG !! tag:example.com,2000:\n---\na: !!int 5", `{"a":"5"}`},
		"directive bool":                  {"%TAG !! tag:example.com,2000:\n---\na: !!bool yes", `{"a":"yes"}`},
		"directive custom":                {"%TAG ! tag:example.com,2000:\n---\na: !foo 5", `{"a":"5"}`},
		"custom alias redefined":          {"x: &x 1\nm: &m [*x]\ny: &x 2\nc: !custom *m", `{"c":[1],"m":[1],"x":1,"y":2}`},
		"tagged alias redefined":          {"x: &x 1\nt: &t !!str *x\ny: &x 2\nb: !!str *t", `{"b":"1","t":"1","x":1,"y":2}`},
		"tagged alias inner anchor":       {"t: &t !!str &x 1\ny: &x 2\nb: !!str *t\nd: *x", `{"b":"1","d":2,"t":"1","y":2}`},
		"timestamp over failing bool":     {"b: !!timestamp !!bool x", `{"b":"0001-01-01T00:00:00Z"}`},
	})
}

func TestYAMLToJSONSendsTheListedFixes(t *testing.T) {
	t.Parallel()

	assertYAMLBodies(t, map[string]yamlBody{
		"merge tagged map":               {"<<: !!map {a: 1}", `{"a":1}`},
		"merge list tagged map":          {"<<: [!!map {a: 1}]", `{"a":1}`},
		"merge seq tag":                  {"<<: !!seq [{a: 1}]", `{"a":1}`},
		"str over timestamp anchor":      {"a: &a !!timestamp 2001-12-14\nb: !!str *a", `{"a":"2001-12-14T00:00:00Z","b":"2001-12-14T00:00:00Z"}`},
		"str over binary anchor":         {"a: &a !!binary aGk=\nb: !!str *a", `{"a":"aGk=","b":"aGk="}`},
		"str over binary tag":            {"b: !!str !!binary aGk=", `{"b":"aGk="}`},
		"binary over text anchor":        {"a: &a aGk=\nb: !!binary *a", `{"a":"aGk=","b":"aGk="}`},
		"binary":                         {"v: !!binary aGk=", `{"v":"aGk="}`},
		"binary alias":                   {"v: &a !!binary aGk=\nw: *a", `{"v":"aGk=","w":"aGk="}`},
		"binary block":                   {"v: !!binary |\n  aGVs\n  bG8=\n", `{"v":"aGVsbG8="}`},
		"null key":                       {"~: 1", `{"":1}`},
		"merge list order":               {"one: &one {a: 1}\ntwo: &two {a: 2, b: 2}\nchild:\n  <<: [*one, *two]", `{"child":{"a":1,"b":2},"one":{"a":1},"two":{"a":2,"b":2}}`},
		"merge before key":               {"child:\n  b: 3\n  <<: {b: 2, c: 4}", `{"child":{"b":3,"c":4}}`},
		"str over timestamp zone anchor": {"a: &a !!timestamp 2001-12-14T21:59:43.10-05:00\nb: !!str *a", `{"a":"2001-12-14T21:59:43.1-05:00","b":"2001-12-14T21:59:43.1-05:00"}`},
		"merge alias tagged map":         {"m: &m !!map {a: 1}\nc:\n  <<: *m", `{"c":{"a":1},"m":{"a":1}}`},
		"binary over block text anchor":  {"a: &a |\n  aGk=\nb: !!binary *a", `{"a":"aGk=\n","b":"aGk="}`},
		"timestamp tag key":              {"!!timestamp 2001-12-14: y", `{"2001-12-14T00:00:00Z":"y"}`},
		"binary tag key":                 {"!!binary aGk=: y", `{"aGk=":"y"}`},
		"binary over str tag anchor":     {"a: &a !!str aGk=\nb: !!binary *a", `{"a":"aGk=","b":"aGk="}`},
		"custom merge":                   {"<<: !custom {a: 1}", `{"a":1}`},
		"binary over custom anchor":      {"a: &a !custom aGk=\nb: !!binary *a", `{"a":"aGk=","b":"aGk="}`},
		"binary over str tag":            {"b: !!binary !!str aGk=", `{"b":"aGk="}`},
		"int over binary":                {"b: !!int !!binary '1234'", `{"b":1234}`},
		"bool over binary":               {"b: !!bool !!binary 'true'", `{"b":true}`},
		"merge alias redefined":          {"x: &x 1\nm: &m {a: *x}\ny: &x 2\nc: {<<: *m}", `{"c":{"a":1},"m":{"a":1},"x":1,"y":2}`},
	})
}

func assertYAMLRefusals(t *testing.T, inputs map[string]string) {
	t.Helper()
	for name, input := range inputs {
		_, err := YAMLToJSON([]byte(input))
		require.Error(t, err, name)
	}
}

func TestYAMLToJSONRefusesWhatDca5f3cRefused(t *testing.T) {
	t.Parallel()

	assertYAMLRefusals(t, map[string]string{
		"merge null":                   "<<: ~",
		"merge list null":              "<<: [~]",
		"merge null alias":             "n: &n ~\nc:\n  <<: *n",
		"merge null tag":               "<<: !!null x",
		"merge tagged alias":           "a: &a {x: 1}\nc:\n  <<: !!map *a",
		"merge string":                 "<<: \"str\"",
		"merge twice":                  "<<: {a: 1}\n<<: {b: 2}",
		"str over map":                 "b: !!str {x: 1}",
		"str over seq":                 "b: !!str [1]",
		"bool over map anchor":         "a: &a {x: 1}\nb: !!bool *a",
		"binary over seq anchor":       "a: &a [1]\nb: !!binary *a",
		"duplicate key":                "a: 1\na: 2",
		"int inline map":               "b: !!int {x: 1}",
		"float inline map":             "b: !!float {x: 1}",
		"map key":                      "? {x: 1}\n: y",
		"seq key":                      "? [1, 2]\n: y",
		"flow seq key":                 "[1, 2]: y",
		"flow map key":                 "{x: 1}: y",
		"binary over binary anchor":    "a: &a !!binary aGk=\nb: !!binary *a",
		"binary over timestamp anchor": "a: &a !!timestamp 2001-12-14\nb: !!binary *a",
		"binary over null tag anchor":  "a: &a !!null x\nb: !!binary *a",
		"binary over int tag anchor":   "a: &a !!int 5\nb: !!binary *a",
		"float over inf alias":         "a: &a !!float .inf\nb: !!str *a",
		"int over failing bool":        "b: !!int !!bool x",
		"float over failing bool":      "b: !!float !!bool x",
		"int over failing binary":      "b: !!int !!binary !!int 5",
	})
}

func TestYAMLToJSONRefusesTheListedFixes(t *testing.T) {
	t.Parallel()

	assertYAMLRefusals(t, map[string]string{
		"str over map anchor":        "a: &a {x: 1}\nb: !!str *a",
		"str over seq anchor":        "a: &a [1, 2]\nb: !!str *a",
		"str over custom map anchor": "a: &a !custom {x: 1}\nb: !!str *a",
		"alias key to map":           "a: &k {x: 1}\n*k : y",
	})
}
