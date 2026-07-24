// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package apiform

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type describedReader struct {
	*strings.Reader
}

func (describedReader) Filename() string    { return `file"name.txt` }
func (describedReader) ContentType() string { return "text/plain" }

type namedReader struct {
	*strings.Reader
}

func (namedReader) Name() string { return "/tmp/named.txt" }

func encodeForCoverage(t *testing.T, value any, format FormFormat) (string, error) {
	t.Helper()
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	require.NoError(t, writer.SetBoundary("coverage"))
	err := MarshalWithSettings(value, writer, format)
	closeErr := writer.Close()
	if err == nil {
		err = closeErr
	}
	return buffer.String(), err
}

func TestMarshalPrimitiveAndPointerVariants(t *testing.T) {
	t.Parallel()

	value := "pointer"
	var nilPointer *string
	output, err := encodeForCoverage(t, map[string]any{
		"false":   false,
		"float32": float32(1.25),
		"int8":    int8(-8),
		"pointer": &value,
		"nil":     nilPointer,
		"uint8":   uint8(8),
	}, FormatRepeat)
	require.NoError(t, err)
	assert.Contains(t, output, "pointer")
	assert.Contains(t, output, "false")

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	require.NoError(t, Marshal(map[string]any{"value": "default"}, writer))
	require.NoError(t, writer.Close())
	assert.Contains(t, buffer.String(), "default")

	_, err = encodeForCoverage(t, map[string]any{"invalid": struct{}{}}, FormatRepeat)
	require.Error(t, err)
	require.NoError(t, Marshal(nil, multipart.NewWriter(io.Discard)))
}

func TestMarshalArrayFormatsAndFailures(t *testing.T) {
	t.Parallel()

	for _, format := range []FormFormat{
		FormatRepeat,
		FormatBrackets,
		FormatIndicesDots,
		FormatIndicesBrackets,
		FormatComma,
	} {
		_, err := encodeForCoverage(t, []any{"value", int64(2), uint64(3), 4.5, true, nil}, format)
		require.NoError(t, err)
	}

	_, err := encodeForCoverage(t, []any{map[string]any{"nested": true}}, FormatComma)
	require.Error(t, err)
	_, err = encodeForCoverage(t, []string{"value"}, FormFormat(99))
	require.Error(t, err)
	_, err = encodeForCoverage(t, map[int]string{1: "value"}, FormatRepeat)
	require.Error(t, err)
}

func TestMarshalReadersAndEscaping(t *testing.T) {
	t.Parallel()

	output, err := encodeForCoverage(t, map[string]any{
		`key"name`: describedReader{Reader: strings.NewReader("contents")},
		"named":    namedReader{Reader: strings.NewReader("named contents")},
		"plain":    strings.NewReader("plain contents"),
	}, FormatRepeat)
	require.NoError(t, err)
	assert.Contains(t, output, `file\"name.txt`)
	assert.Contains(t, output, "named.txt")
	assert.Contains(t, output, "anonymous_file")
	assert.Equal(t, `a\\b\"c`, escapeQuotes(`a\b"c`))
}
