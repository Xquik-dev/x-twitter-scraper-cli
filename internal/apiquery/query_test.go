// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package apiquery

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		param  Param
		format ArrayQueryFormat
		want   string
	}{
		"scalar":                {Scalar("hello world"), ArrayQueryFormatComma, "query=hello+world"},
		"scalar repeat":         {Scalar("42"), ArrayQueryFormatRepeat, "query=42"},
		"empty scalar":          {Scalar(""), ArrayQueryFormatComma, "query="},
		"list comma":            {List([]string{"10", "20", "30"}), ArrayQueryFormatComma, "query=10%2C20%2C30"},
		"list repeat":           {List([]string{"10", "20", "30"}), ArrayQueryFormatRepeat, "query=10&query=20&query=30"},
		"empty list comma":      {List(nil), ArrayQueryFormatComma, "query="},
		"empty list repeat":     {List(nil), ArrayQueryFormatRepeat, ""},
		"comma inside a scalar": {Scalar("a,b"), ArrayQueryFormatRepeat, "query=a%2Cb"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := Encode(map[string]Param{"query": test.param}, test.format)
			assert.Equal(t, test.want, got.Encode())
		})
	}
}

func TestParamMap(t *testing.T) {
	t.Parallel()

	upper, err := List([]string{"a", "b"}).Map(func(value string) (string, error) {
		return strings.ToUpper(value), nil
	})
	require.NoError(t, err)
	assert.Equal(t, "query=A%2CB", Encode(map[string]Param{"query": upper}, ArrayQueryFormatComma).Encode())

	failure := errors.New("transform failed")
	_, err = Scalar("a").Map(func(string) (string, error) { return "", failure })
	require.ErrorIs(t, err, failure)
}
