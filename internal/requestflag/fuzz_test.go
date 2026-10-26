// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"unicode/utf8"
)

func FuzzJSONFlagValue(f *testing.F) {
	for _, seed := range []string{`{a: 1, b: [x, "y"]}`, `[1, 2]`, `null`, `@file.txt`, `abc_def-1.2=3`, `{a: 1`, `key: value`, ``, `'quoted'`, `1.5e3`, `{"limit": 20.0}`} {
		f.Add(seed, "\x1b[31m")
	}
	f.Add(`{"text": "\u0001"}`, "\x00\t\n\"\\<&>\u2028")
	f.Fuzz(func(t *testing.T, input string, text string) {
		if raw, err := parseJSON(input); err == nil {
			var compact bytes.Buffer
			if err := json.Compact(&compact, raw); err != nil {
				t.Fatalf("parseJSON(%q) returned invalid JSON %q: %v", input, raw, err)
			}
			if !bytes.Equal(compact.Bytes(), raw) {
				t.Fatalf("parseJSON(%q) returned %q, which is not compact", input, raw)
			}
		}

		if !utf8.ValidString(text) || strings.ContainsFunc(text, outsideYAMLPrintable) {
			return
		}
		document, err := json.Marshal(map[string]string{"text": text})
		if err != nil {
			t.Fatal(err)
		}
		raw, err := parseJSON(string(document))
		if err != nil {
			t.Fatalf("parseJSON(%s) failed: %v", document, err)
		}
		var decoded map[string]string
		if err := json.Unmarshal(raw, &decoded); err != nil || decoded["text"] != text {
			t.Fatalf("parseJSON(%s) returned %s", document, raw)
		}
	})
}

func outsideYAMLPrintable(r rune) bool {
	return r == 0x7f || (r >= 0x80 && r <= 0x9f) || r == 0xfeff || r == 0xfffe || r == 0xffff
}
