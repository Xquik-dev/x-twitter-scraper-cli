package apiquery

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalCoversPrimitiveWidthsAndPointers(t *testing.T) {
	t.Parallel()

	value := "pointer"
	var nilValue *string
	values, err := Marshal(map[string]any{
		"false":       false,
		"float32":     float32(1.25),
		"int8":        int8(-8),
		"int16":       int16(-16),
		"nil-pointer": nilValue,
		"pointer":     &value,
		"uint8":       uint8(8),
		"uint16":      uint16(16),
	})
	require.NoError(t, err)
	assert.Equal(t, "false", values.Get("false"))
	assert.Equal(t, "1.25", values.Get("float32"))
	assert.Equal(t, "-8", values.Get("int8"))
	assert.Equal(t, "-16", values.Get("int16"))
	assert.Equal(t, "", values.Get("nil-pointer"))
	assert.Equal(t, "pointer", values.Get("pointer"))
	assert.Equal(t, "8", values.Get("uint8"))
	assert.Equal(t, "16", values.Get("uint16"))

	empty, err := Marshal(nil)
	require.NoError(t, err)
	assert.Nil(t, empty)
}

func TestEncoderRejectsUnknownArrayFormat(t *testing.T) {
	t.Parallel()

	queryEncoder := encoder{
		settings: QuerySettings{ArrayFormat: ArrayQueryFormat(99)},
	}
	assert.PanicsWithValue(t, "Unknown ArrayFormat value: 99", func() {
		_, _ = queryEncoder.Encode("values", reflect.ValueOf([]string{"value"}))
	})
}

func TestEncoderIgnoresUnsupportedPrimitive(t *testing.T) {
	t.Parallel()

	queryEncoder := encoder{}
	pairs, err := queryEncoder.Encode("value", reflect.ValueOf(struct{}{}))
	require.NoError(t, err)
	assert.Empty(t, pairs)
}
