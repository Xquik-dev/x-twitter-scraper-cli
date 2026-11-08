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
		"alias in a list inside its own anchor":                   {"x: &a {k: [*a]}", `{"x":{"k":[null]}}`},
		"alias inside its own anchor":                             {"&a [*a]", `[null]`},
		"alias key":                                               {"a: &k key\n*k : 2", `{"a":"key","key":2}`},
		"alias key to int":                                        {"a: &k 5\n*k : x", `{"5":"x","a":5}`},
		"alias key to map":                                        {"a: &k {x: 1}\n*k : y", `{"a":{"x":1},"map[x:1]":"y"}`},
		"alias read again after a later anchor":                   {"<<: &c {q: 1}\nx: *c\ny: &c 5", `{"q":1,"x":5,"y":5}`},
		"alias to a merge anchor read again after a later anchor": {"m: {<<: &c {q: 1}}\nx: *c\ny: &c 5", `{"m":{"q":1},"x":5,"y":5}`},
		"alias to an anchor after its list":                       {"a: &a [&b 1]\nc: *b", `{"a":[1],"c":1}`},
		"alias to an anchor in a merge":                           {"x: &a {<<: &m {q: 1}, r: *m}", `{"x":{"q":1,"r":{"q":1}}}`},
		"aliases inside 2 open anchors":                           {"x: &a [&b [*a, *b]]", `{"x":[[null,null]]}`},
		"anchor on map":                                           {"a: &a {x: 1}\nb: *a", `{"a":{"x":1},"b":{"x":1}}`},
		"anchor on seq":                                           {"a: &a [1, x]\nb: *a", `{"a":[1,"x"],"b":[1,"x"]}`},
		"anchor over tag":                                         {"x: &a !!int \"5\"\ny: *a", `{"x":5,"y":5}`},
		"binary bad padding":                                      {"v: !!binary aGk", `{"v":""}`},
		"binary bad text":                                         {"v: !!binary '!!!'", `{"v":""}`},
		"binary over custom int anchor":                           {"a: &a !custom 5\nb: !!binary *a", `{"a":"5","b":""}`},
		"binary tag key":                                          {"!!binary aGk=: y", `{"[104 105]":"y"}`},
		"binary tag over a string alias key":                      {"a: &a !!str 1\n!!binary *a : v", `{"[]":"v","a":"1"}`},
		"binary tag over a tagged alias chain":                    {"x: &x 1\na: &a !!str *x\ny: &x 9\nb: !!binary *a\nk: {!!binary *a : v}", `{"a":"1","b":"","k":{"[]":"v"},"x":1,"y":9}`},
		"bool key":                                                {"true: a", `{"true":"a"}`},
		"bool over bool tag anchor":                               {"a: &a !!bool yes\nb: !!bool *a", `{"a":true,"b":true}`},
		"bool over str tag anchor":                                {"a: &a !!str yes\nb: !!bool *a", `{"a":"yes","b":true}`},
		"bool over text anchor":                                   {"a: &a yes\nb: !!bool *a", `{"a":"yes","b":true}`},
		"bool tag over a custom tag alias":                        {"x: &x 1\na: &a !custom *x\ny: &x 9\nb: !!bool *a\nk: {!!bool *a : v}", `{"a":1,"b":true,"k":{"true":"v"},"x":1,"y":9}`},
		"bool tag over a custom tag alias of true":                {"x: &x true\na: &a !custom *x\ny: &x 9\nb: !!bool *a\nk: {!!bool *a : v}", `{"a":true,"b":true,"k":{"true":"v"},"x":true,"y":9}`},
		"custom alias redefined":                                  {"x: &x 1\nm: &m [*x]\ny: &x 2\nc: !custom *m", `{"c":[1],"m":[1],"x":1,"y":2}`},
		"custom float key":                                        {"!custom 1.10: y", `{"1.10":"y"}`},
		"custom over map anchor":                                  {"a: &a {x: 1}\nb: !custom *a", `{"a":{"x":1},"b":{"x":1}}`},
		"custom tag alias after an anchor redefinition":           {"x: &x 1\na: &a !custom *x\ny: &x 9\nb: *a\nk: {*a : v}", `{"a":1,"b":1,"k":{"1":"v"},"x":1,"y":9}`},
		"date":                                 {"v: 2001-12-14", `{"v":"2001-12-14"}`},
		"directive bool":                       {"%TAG !! tag:example.com,2000:\n---\na: !!bool yes", `{"a":"yes"}`},
		"directive custom":                     {"%TAG ! tag:example.com,2000:\n---\na: !foo 5", `{"a":"5"}`},
		"directive int":                        {"%TAG !! tag:example.com,2000:\n---\na: !!int 5", `{"a":"5"}`},
		"explicit key anchor":                  {"? &k key\n: 1\nv: *k", `{"key":1,"v":"key"}`},
		"float key":                            {"1.5: a", `{"1.5":"a"}`},
		"float over custom anchor":             {"a: &a !custom 2.5\nb: !!float *a", `{"a":"2.5","b":2.5}`},
		"float over float tag anchor":          {"a: &a !!float 2.5\nb: !!float *a", `{"a":2.5,"b":2.5}`},
		"float over int tag anchor":            {"a: &a !!int 5\nb: !!float *a", `{"a":5,"b":5}`},
		"float over map anchor":                {"a: &a {x: 1}\nb: !!float *a", `{"a":{"x":1},"b":0}`},
		"float over str tag":                   {"b: !!float !!str 2.5", `{"b":2.5}`},
		"float over str tag anchor":            {"a: &a !!str \"2.5\"\nb: !!float *a", `{"a":"2.5","b":2.5}`},
		"float over text anchor":               {"a: &a \"2.5\"\nb: !!float *a", `{"a":"2.5","b":2.5}`},
		"float over timestamp anchor":          {"a: &a !!timestamp 2001-12-14\nb: !!float *a", `{"a":"2001-12-14T00:00:00Z","b":0}`},
		"float tag key":                        {"!!float 1: y", `{"1":"y"}`},
		"float tag key 1e6":                    {"!!float 1000000: y", `{"1e+06":"y"}`},
		"int key":                              {"1: a", `{"1":"a"}`},
		"int key anchor":                       {"&k 5: x\nv: *k", `{"5":"x","v":5}`},
		"int over binary":                      {"b: !!int !!binary '1234'", `{"b":0}`},
		"int over bool anchor":                 {"a: &a true\nb: !!int *a", `{"a":true,"b":0}`},
		"int over float anchor":                {"a: &a 2.5\nb: !!int *a", `{"a":2.5,"b":0}`},
		"int over float tag anchor":            {"a: &a !!float 7\nb: !!int *a", `{"a":7,"b":7}`},
		"int over map anchor":                  {"a: &a {x: 1}\nb: !!int *a", `{"a":{"x":1},"b":0}`},
		"int over null anchor":                 {"a: &a ~\nb: !!int *a", `{"a":null,"b":0}`},
		"int over quoted int anchor":           {"a: &a \"7\"\nb: !!int *a", `{"a":"7","b":7}`},
		"int over seq anchor":                  {"a: &a [1]\nb: !!int *a", `{"a":[1],"b":0}`},
		"int over str tag":                     {"b: !!int !!str 7", `{"b":7}`},
		"int over str tag anchor":              {"a: &a !!str 5\nb: !!int *a", `{"a":"5","b":5}`},
		"int over text anchor":                 {"a: &a x\nb: !!int *a", `{"a":"x","b":0}`},
		"int over timestamp anchor":            {"a: &a !!timestamp 2001-12-14\nb: !!int *a", `{"a":"2001-12-14T00:00:00Z","b":0}`},
		"int tag over an alias of a map tag":   {"a: &a !!map {x: 1}\nb: !!int *a", `{"a":{"x":1},"b":0}`},
		"int tag over an integer past 64 bits": {"!!int 99999999999999999999", `9223372036854775807`},
		"integer past 64 bits":                 {"a: 99999999999999999999", `{"a":"99999999999999999999"}`},
		"key anchor":                           {"&k key: 1\nv: *k", `{"key":1,"v":"key"}`},
		"key anchor merge value":               {"&k key: {a: 1}\nv: *k", `{"key":{"a":1},"v":"key"}`},
		"key anchor self value":                {"&k key: *k", `{"key":"key"}`},
		"merge after an anchor redefinition":   {"x: &x 1\nm: &m {q: *x}\ny: &x 9\nc: {<<: *m}", `{"c":{"q":9},"m":{"q":1},"x":1,"y":9}`},
		"merge after an explicit key":          {"m: &m {a: 1}\nc: {a: 2, <<: *m}", `{"c":{"a":1},"m":{"a":1}}`},
		"merge alias redefined":                {"x: &x 1\nm: &m {a: *x}\ny: &x 2\nc: {<<: *m}", `{"c":{"a":2},"m":{"a":1},"x":1,"y":2}`},
		"merge before key":                     {"child:\n  b: 3\n  <<: {b: 2, c: 4}", `{"child":{"b":2,"c":4}}`},
		"merge empty list":                     {"<<: []", `{}`},
		"merge empty map":                      {"<<: {}", `{}`},
		"merge list order":                     {"one: &one {a: 1}\ntwo: &two {a: 2, b: 2}\nchild:\n  <<: [*one, *two]", `{"child":{"a":2,"b":2},"one":{"a":1},"two":{"a":2,"b":2}}`},
		"merge list takes the later value":     {"m: &m {a: 1}\nn: &n {a: 2, b: 2}\nc: {<<: [*m, *n]}", `{"c":{"a":2,"b":2},"m":{"a":1},"n":{"a":2,"b":2}}`},
		"merge override":                       {"base: &base {a: 1, b: 2}\nchild:\n  <<: *base\n  b: 3", `{"base":{"a":1,"b":2},"child":{"a":1,"b":3}}`},
		"nested key anchor":                    {"m:\n  &k key: 1\n  v: *k", `{"m":{"key":1,"v":"key"}}`},
		"null document before a document":      {"--- !!null\n---\na: 1", `{"a":1}`},
		"null key":                             {"~: 1", `{"null":1}`},
		"null over map anchor":                 {"a: &a {x: 1}\nb: !!null *a", `{"a":{"x":1},"b":null}`},
		"null tag over a tagged alias key":     {"a: &a !!str 1\n!!null *a : v", `{"a":"1","null":"v"}`},
		"null tag over an alias key":           {"a: &a 1\n!!null *a : v", `{"a":1,"null":"v"}`},
		"null word key":                        {"null: v", `{"null":"v"}`},
		"redefined anchor":                     {"a: &a 1\nb: &a 2\nc: *a", `{"a":1,"b":2,"c":2}`},
		"str over big float tag anchor":        {"a: &a !!float 1e21\nb: !!str *a", `{"a":1e+21,"b":"1e+21"}`},
		"str over binary tag":                  {"b: !!str !!binary aGk=", `{"b":"[104 105]"}`},
		"str over bool anchor":                 {"a: &a true\nb: !!str *a", `{"a":true,"b":"true"}`},
		"str over bool tag anchor":             {"a: &a !!bool yes\nb: !!str *a", `{"a":true,"b":"true"}`},
		"str over custom float":                {"b: !!str !custom 1.10", `{"b":"1.10"}`},
		"str over custom map anchor":           {"a: &a !custom {x: 1}\nb: !!str *a", `{"a":{"x":1},"b":"map[x:1]"}`},
		"str over custom tag anchor":           {"a: &a !custom hi\nb: !!str *a", `{"a":"hi","b":"hi"}`},
		"str over date":                        {"b: !!str 2001-12-14", `{"b":"2001-12-14"}`},
		"str over date anchor":                 {"a: &a 2001-12-14\nb: !!str *a", `{"a":"2001-12-14","b":"2001-12-14"}`},
		"str over float anchor":                {"a: &a 1.10\nb: !!str *a", `{"a":1.1,"b":"1.1"}`},
		"str over float tag 1e6":               {"b: !!str !!float 1000000", `{"b":"1e+06"}`},
		"str over float tag 1e6 anchor":        {"a: &a !!float 1000000\nb: !!str *a", `{"a":1000000,"b":"1e+06"}`},
		"str over float tag anchor":            {"a: &a !!float 0.0000001\nb: !!str *a", `{"a":1e-7,"b":"1e-07"}`},
		"str over float tag text anchor":       {"a: &a !!float \"2.5\"\nb: !!str *a", `{"a":2.5,"b":"2.5"}`},
		"str over inf":                         {"b: !!str !!float .inf", `{"b":"+Inf"}`},
		"str over int tag anchor":              {"a: &a !!int 5\nb: !!str *a", `{"a":5,"b":"5"}`},
		"str over int tag float anchor":        {"a: &a !!int 3.5\nb: !!str *a", `{"a":0,"b":"0"}`},
		"str over map anchor":                  {"a: &a {x: 1}\nb: !!str *a", `{"a":{"x":1},"b":"map[x:1]"}`},
		"str over null anchor":                 {"a: &a ~\nb: !!str *a", `{"a":null,"b":""}`},
		"str over null tag anchor":             {"a: &a !!null x\nb: !!str *a", `{"a":null,"b":""}`},
		"str over seq anchor":                  {"a: &a [1, 2]\nb: !!str *a", `{"a":[1,2],"b":"[1 2]"}`},
		"str over small float tag anchor":      {"a: &a !!float 0.00001\nb: !!str *a", `{"a":0.00001,"b":"1e-05"}`},
		"str over str tag anchor":              {"a: &a !!str 5\nb: !!str *a", `{"a":"5","b":"5"}`},
		"str over text anchor":                 {"a: &a hi\nb: !!str *a", `{"a":"hi","b":"hi"}`},
		"str over timestamp anchor":            {"a: &a !!timestamp 2001-12-14\nb: !!str *a", `{"a":"2001-12-14T00:00:00Z","b":"2001-12-14 00:00:00 +0000 UTC"}`},
		"str over timestamp zone anchor":       {"a: &a !!timestamp 2001-12-14T21:59:43.10-05:00\nb: !!str *a", `{"a":"2001-12-14T21:59:43.1-05:00","b":"2001-12-14 21:59:43.1 -0500 -0500"}`},
		"string tag over a custom tag alias":   {"x: &x 1\na: &a !custom *x\ny: &x 9\nb: !!str *a\nk: {!!str *a : v}", `{"a":1,"b":"1","k":{"1":"v"},"x":1,"y":9}`},
		"string tag over a custom tag alias after a redefinition": {"x: &x 1\na: &a !custom *x\ny: &x 9\nb: !!str *a", `{"a":1,"b":"1","x":1,"y":9}`},
		"string tag over an anchored map":                         {"a: !!str &x {x: 1}", `{"a":"map[x:1]"}`},
		"tag over anchor":                                         {"x: !!int &a 5\ny: *a", `{"x":5,"y":5}`},
		"tagged alias inner anchor":                               {"t: &t !!str &x 1\ny: &x 2\nb: !!str *t\nd: *x", `{"b":"1","d":2,"t":"1","y":2}`},
		"tagged alias redefined":                                  {"x: &x 1\nt: &t !!str *x\ny: &x 2\nb: !!str *t", `{"b":"1","t":"1","x":1,"y":2}`},
		"tagged key":                                              {"!!str 5: x", `{"5":"x"}`},
		"timestamp":                                               {"v: !!timestamp 2001-12-14", `{"v":"2001-12-14T00:00:00Z"}`},
		"timestamp alias":                                         {"v: &a !!timestamp 2001-12-14\nw: *a", `{"v":"2001-12-14T00:00:00Z","w":"2001-12-14T00:00:00Z"}`},
		"timestamp over failing bool":                             {"b: !!timestamp !!bool x", `{"b":"0001-01-01T00:00:00Z"}`},
		"timestamp over map anchor":                               {"a: &a {x: 1}\nb: !!timestamp *a", `{"a":{"x":1},"b":"0001-01-01T00:00:00Z"}`},
		"timestamp over str tag":                                  {"b: !!timestamp !!str 2001-12-14", `{"b":"2001-12-14T00:00:00Z"}`},
		"timestamp over str tag anchor":                           {"a: &a !!str 2001-12-14\nb: !!timestamp *a", `{"a":"2001-12-14","b":"2001-12-14T00:00:00Z"}`},
		"timestamp over text anchor":                              {"a: &a 2001-12-14\nb: !!timestamp *a", `{"a":"2001-12-14","b":"2001-12-14T00:00:00Z"}`},
		"timestamp over timestamp anchor":                         {"a: &a !!timestamp 2001-12-14\nb: !!timestamp *a", `{"a":"2001-12-14T00:00:00Z","b":"2001-12-14T00:00:00Z"}`},
		"timestamp tag key":                                       {"!!timestamp 2001-12-14: y", `{"2001-12-14 00:00:00 +0000 UTC":"y"}`},
		"timestamp tag over a string alias key":                   {"a: &a !!str 1\n!!timestamp *a : v", `{"0001-01-01 00:00:00 +0000 UTC":"v","a":"1"}`},
		"timestamp tag over a string alias of a list":             {"x: &x [1]\na: &a !!str *x\ny: &x 9\nb: !!timestamp *a\nk: {!!timestamp *a : v}", `{"a":"[1]","b":"0001-01-01T00:00:00Z","k":{"0001-01-01 00:00:00 +0000 UTC":"v"},"x":[1],"y":9}`},
		"timestamp tag over a string alias of a map":              {"x: &x {x: 1}\na: &a !!str *x\ny: &x 9\nb: !!timestamp *a\nk: {!!timestamp *a : v}", `{"a":"map[x:1]","b":"0001-01-01T00:00:00Z","k":{"0001-01-01 00:00:00 +0000 UTC":"v"},"x":{"x":1},"y":9}`},
		"timestamp tag over an alias key":                         {"a: &a 1\n!!timestamp *a : v", `{"0001-01-01 00:00:00 +0000 UTC":"v","a":1}`},
		"timestamp zone":                                          {"v: !!timestamp 2001-12-14T21:59:43.10-05:00", `{"v":"2001-12-14T21:59:43.1-05:00"}`},
	})
}

func TestYAMLToJSONSendsBinaryAsBase64(t *testing.T) {
	t.Parallel()

	assertYAMLBodies(t, map[string]yamlBody{
		"binary":                                   {"v: !!binary aGk=", `{"v":"aGk="}`},
		"binary alias":                             {"v: &a !!binary aGk=\nw: *a", `{"v":"aGk=","w":"aGk="}`},
		"binary block":                             {"v: !!binary |\n  aGVs\n  bG8=\n", `{"v":"aGVsbG8="}`},
		"binary over block text anchor":            {"a: &a |\n  aGk=\nb: !!binary *a", `{"a":"aGk=\n","b":"aGk="}`},
		"binary over custom anchor":                {"a: &a !custom aGk=\nb: !!binary *a", `{"a":"aGk=","b":"aGk="}`},
		"binary over str tag":                      {"b: !!binary !!str aGk=", `{"b":"aGk="}`},
		"binary over str tag anchor":               {"a: &a !!str aGk=\nb: !!binary *a", `{"a":"aGk=","b":"aGk="}`},
		"binary over text anchor":                  {"a: &a aGk=\nb: !!binary *a", `{"a":"aGk=","b":"aGk="}`},
		"binary tag in a list over a string alias": {"a: &a !!str true\ns: [!!binary *a]", `{"a":"true","s":["true"]}`},
		"binary tag over a string alias":           {"a: &a !!str true\nb: !!binary *a", `{"a":"true","b":"true"}`},
		"int tag over a binary alias":              {"a: &a !!binary 2001-12-14\nb: !!int *a", `{"a":"2001","b":0}`},
		"int tag over a binary alias key":          {"a: &a !!binary 2001-12-14\n!!int *a : v", `{"0":"v","a":"2001"}`},
		"str over binary anchor":                   {"a: &a !!binary aGk=\nb: !!str *a", `{"a":"aGk=","b":"[104 105]"}`},
	})
}

func TestYAMLToJSONSendsTheValueOfAFinishedAnchorUnderAnOpenAnchor(t *testing.T) {
	t.Parallel()

	assertYAMLBodies(t, map[string]yamlBody{
		"alias to a finished anchor from a later anchor":      {"x: &a [&b 1, &c [*b]]", `{"x":[1,[1]]}`},
		"alias to a finished anchor in a list":                {"x: &a [&b 1, *b]", `{"x":[1,1]}`},
		"alias to a finished anchor in a map":                 {"base: &base\n  a: &a 1\n  b: *a", `{"base":{"a":1,"b":1}}`},
		"alias to a finished anchor under a top-level anchor": {"&a [&b 1, *b]", `[1,1]`},
		"alias to a finished nested anchor":                   {"x: &a {p: &b {q: &c 1}, r: *c}", `{"x":{"p":{"q":1},"r":1}}`},
	})
}

func TestYAMLToJSONRefusesWhatDca5f3cRefused(t *testing.T) {
	t.Parallel()

	for name, input := range map[string]string{
		"alias before its anchor":                   "a: [*c]\n<<: &c {q: 1}\nc: &c 2",
		"alias to a scalar anchor under a null tag": "x: !!null &x 1\ny: *x",
		"alias to an anchor under a null tag":       "a: &a !!null &x {x: 1}\ny: &x 9\nb: *a\nd: *x",
		"binary over binary anchor":                 "a: &a !!binary aGk=\nb: !!binary *a",
		"binary over int tag anchor":                "a: &a !!int 5\nb: !!binary *a",
		"binary over null tag anchor":               "a: &a !!null x\nb: !!binary *a",
		"binary over seq anchor":                    "a: &a [1]\nb: !!binary *a",
		"binary over timestamp anchor":              "a: &a !!timestamp 2001-12-14\nb: !!binary *a",
		"bool over binary":                          "b: !!bool !!binary 'true'",
		"bool over map anchor":                      "a: &a {x: 1}\nb: !!bool *a",
		"custom merge":                              "<<: !custom {a: 1}",
		"duplicate key":                             "a: 1\na: 2",
		"float inline map":                          "b: !!float {x: 1}",
		"float over failing bool":                   "b: !!float !!bool x",
		"float over inf alias":                      "a: &a !!float .inf\nb: !!str *a",
		"flow map key":                              "{x: 1}: y",
		"flow seq key":                              "[1, 2]: y",
		"int inline map":                            "b: !!int {x: 1}",
		"int over failing binary":                   "b: !!int !!binary !!int 5",
		"int over failing bool":                     "b: !!int !!bool x",
		"map key":                                   "? {x: 1}\n: y",
		"merge alias tagged map":                    "m: &m !!map {a: 1}\nc:\n  <<: *m",
		"merge list null":                           "<<: [~]",
		"merge list tagged map":                     "<<: [!!map {a: 1}]",
		"merge null":                                "<<: ~",
		"merge null alias":                          "n: &n ~\nc:\n  <<: *n",
		"merge null tag":                            "<<: !!null x",
		"merge of a map tag":                        "c: {<<: !!map {a: 1}}",
		"merge of an anchored map tag":              "x: &x 1\nm: &m !!map {q: *x}\ny: &x 9\nc: {<<: *m}",
		"merge seq tag":                             "<<: !!seq [{a: 1}]",
		"merge string":                              "<<: \"str\"",
		"merge tagged alias":                        "a: &a {x: 1}\nc:\n  <<: !!map *a",
		"merge tagged map":                          "<<: !!map {a: 1}",
		"merge through a custom tag":                "x: &x 1\nm: &m {q: *x}\ny: &x 9\nc: {<<: !custom *m}",
		"merge twice":                               "<<: {a: 1}\n<<: {b: 2}",
		"null tag over 2 anchors":                   "a: &a !!null &x {x: 1}\ny: &x 9\nd: *x",
		"seq key":                                   "? [1, 2]\n: y",
		"str over map":                              "b: !!str {x: 1}",
		"str over seq":                              "b: !!str [1]",
		"string tag over an alias to an anchor under a null tag": "a: &a !!null &x {x: 1}\ny: &x 9\nb: !!str *a\nd: *x",
	} {
		_, err := YAMLToJSON([]byte(input))
		require.Error(t, err, name)
	}
}
