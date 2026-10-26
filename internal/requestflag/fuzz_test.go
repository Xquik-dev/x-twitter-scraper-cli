// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"bytes"
	"encoding/json"
	"testing"
)

func FuzzJSONFlagValue(f *testing.F) {
	for _, seed := range []string{`{a: 1, b: [x, "y"]}`, `[1, 2]`, `null`, `@file.txt`, `abc_def-1.2=3`, `{a: 1`, `key: value`, ``, `'quoted'`, `1.5e3`} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		raw, err := parseJSON(input)
		if err != nil {
			return
		}
		var compact bytes.Buffer
		if err := json.Compact(&compact, raw); err != nil {
			t.Fatalf("parseJSON(%q) returned invalid JSON %q: %v", input, raw, err)
		}
		if !bytes.Equal(compact.Bytes(), raw) {
			t.Fatalf("parseJSON(%q) returned %q, which is not compact", input, raw)
		}
	})
}
