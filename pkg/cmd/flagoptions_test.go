// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsUTF8TextFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		content  []byte
		expected bool
	}{
		{[]byte("Hello, world!"), true},
		{[]byte(`{"key": "value"}`), true},
		{[]byte(`<?xml version="1.0"?><root/>`), true},
		{[]byte(`function test() {}`), true},
		{[]byte{0xFF, 0xD8, 0xFF, 0xE0}, false},
		{[]byte{0x00, 0x01, 0xFF, 0xFE}, false},
		{[]byte("Hello \xFF\xFE"), false},
		{[]byte("Hello ☺️"), true},
		{[]byte{}, true},
	}

	for _, tt := range tests {
		require.Equal(t, tt.expected, isUTF8TextFile(tt.content))
	}
}

func TestEmbedJSON(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	configContent := "host=localhost\nport=8080"
	templateContent := "<html><body>Hello</body></html>"
	dataContent := `{"key": "value"}`

	writeTestFile(t, tmpDir, "config.txt", configContent)
	writeTestFile(t, tmpDir, "template.html", templateContent)
	writeTestFile(t, tmpDir, "data.json", dataContent)
	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	writeTestFile(t, tmpDir, "image.jpg", string(jpegHeader))

	text := func(value string) json.RawMessage {
		encoded, err := json.Marshal(value)
		require.NoError(t, err)
		return encoded
	}
	file := func(prefix, name string) json.RawMessage {
		return text(prefix + filepath.Join(tmpDir, name))
	}
	object := func(fields map[string]json.RawMessage) json.RawMessage {
		encoded, err := json.Marshal(fields)
		require.NoError(t, err)
		return encoded
	}
	list := func(items ...json.RawMessage) json.RawMessage {
		encoded, err := json.Marshal(items)
		require.NoError(t, err)
		return encoded
	}
	number := json.RawMessage(`42`)
	truth := json.RawMessage(`true`)

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "object with file references",
			input: string(object(map[string]json.RawMessage{"config": file("@", "config.txt"), "template": file("@file://", "template.html"), "count": number})),
			want:  string(object(map[string]json.RawMessage{"config": text(configContent), "template": text(templateContent), "count": number})),
		},
		{
			name:  "list with file references",
			input: string(list(file("@", "config.txt"), number, truth, file("@file://", "data.json"))),
			want:  string(list(text(configContent), number, truth, text(dataContent))),
		},
		{
			name:  "nested structures",
			input: string(object(map[string]json.RawMessage{"outer": object(map[string]json.RawMessage{"inner": list(file("@", "config.txt"), object(map[string]json.RawMessage{"data": file("@", "data.json")}))})})),
			want:  string(object(map[string]json.RawMessage{"outer": object(map[string]json.RawMessage{"inner": list(text(configContent), object(map[string]json.RawMessage{"data": text(dataContent)}))})})),
		},
		{
			name:  "base64 encoding",
			input: string(object(map[string]json.RawMessage{"encoded": file("@data://", "config.txt"), "image": file("@", "image.jpg")})),
			want:  string(object(map[string]json.RawMessage{"encoded": text(base64.StdEncoding.EncodeToString([]byte(configContent))), "image": text(base64.StdEncoding.EncodeToString(jpegHeader))})),
		},
		{
			name:    "missing file with @ prefix",
			input:   `{"missing":"@file.txt"}`,
			wantErr: true,
		},
		{
			name:  "non-file-like thing with @ prefix",
			input: `{"username":"@user","favorite_symbol":"@"}`,
			want:  `{"username":"@user","favorite_symbol":"@"}`,
		},
		{
			name:    "missing file with @file:// prefix",
			input:   `{"missing":"@file:///nonexistent/file.txt"}`,
			wantErr: true,
		},
		{
			name:    "missing file with @data:// prefix",
			input:   `{"missing":"@data:///nonexistent/file.txt"}`,
			wantErr: true,
		},
		{
			name:  "escaping",
			input: `{"simple":"\\@file.txt","file":"\\@file://file.txt","data":"\\@data://file.txt","keep_escape":"user\\@example.com"}`,
			want:  `{"simple":"@file.txt","file":"@file://file.txt","data":"@data://file.txt","keep_escape":"user\\@example.com"}`,
		},
		{
			name:  "primitive types",
			input: `{"int":123,"float":45.67,"bool":true,"null":null,"string":"no prefix","email":"user@example.com"}`,
			want:  `{"int":123,"float":45.67,"bool":true,"null":null,"string":"no prefix","email":"user@example.com"}`,
		},
		{
			name:  "bare scalars",
			input: `[1,2,3]`,
			want:  `[1,2,3]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := embedJSON(json.RawMessage(tt.input), &onceStdinReader{})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestEmbedJSONRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	for _, input := range []string{`"open`, `[1,`, `{"a":`} {
		_, err := embedJSON(json.RawMessage(input), &onceStdinReader{})
		require.Error(t, err, input)
	}
	_, err := embedJSON(json.RawMessage(`[{"a":"@missing.txt"}]`), &onceStdinReader{})
	require.Error(t, err)
	_, err = embedJSON(json.RawMessage(`{"a":{"b":"@missing.txt"}}`), &onceStdinReader{})
	require.Error(t, err)
}

func TestEmbedStdin(t *testing.T) {
	t.Parallel()

	t.Run("AtDash", func(t *testing.T) {
		t.Parallel()

		stdin := &onceStdinReader{stdinReader: strings.NewReader("piped content")}
		got, err := embedJSON(json.RawMessage(`{"data":"@-"}`), stdin)
		require.NoError(t, err)
		require.JSONEq(t, `{"data":"piped content"}`, string(got))
	})

	t.Run("AtDevStdinAsBase64", func(t *testing.T) {
		t.Parallel()

		stdin := &onceStdinReader{stdinReader: strings.NewReader("piped content")}
		got, err := embedText("@data:///dev/stdin", stdin)
		require.NoError(t, err)
		require.Equal(t, base64.StdEncoding.EncodeToString([]byte("piped content")), got)
	})

	t.Run("StdinReadsOnce", func(t *testing.T) {
		t.Parallel()

		stdin := &onceStdinReader{stdinReader: strings.NewReader("stdin content")}
		_, err := embedJSON(json.RawMessage(`["@-","@file://-"]`), stdin)
		require.ErrorContains(t, err, "already been read")
	})

	t.Run("UnavailableStdin", func(t *testing.T) {
		t.Parallel()

		stdin := &onceStdinReader{failureReason: "stdin was already consumed by piped YAML/JSON input"}
		_, err := embedText("@-", stdin)
		require.ErrorContains(t, err, "cannot read from stdin")
	})

	t.Run("RealFile", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		writeTestFile(t, tmpDir, "test.txt", "file content")
		stdin := &onceStdinReader{stdinReader: strings.NewReader("unused stdin")}
		got, err := embedText("@"+filepath.Join(tmpDir, "test.txt"), stdin)
		require.NoError(t, err)
		require.Equal(t, "file content", got)
	})
}

func TestEmbedParams(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeTestFile(t, tmpDir, "token.txt", "secret")
	params := map[string]apiquery.Param{
		"token": apiquery.Scalar("@" + filepath.Join(tmpDir, "token.txt")),
		"tags":  apiquery.List([]string{"\\@a", "b"}),
	}
	embedded, err := embedParams(params, &onceStdinReader{})
	require.NoError(t, err)
	assert.Equal(t, "tags=%40a%2Cb&token=secret", apiquery.Encode(embedded, apiquery.ArrayQueryFormatComma).Encode())

	_, err = embedParams(map[string]apiquery.Param{"missing": apiquery.Scalar("@missing.txt")}, &onceStdinReader{})
	require.Error(t, err)
}

func writeTestFile(t *testing.T, dir, filename, content string) {
	t.Helper()

	path := filepath.Join(dir, filename)

	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err, "failed to write test file %s", path)
}
