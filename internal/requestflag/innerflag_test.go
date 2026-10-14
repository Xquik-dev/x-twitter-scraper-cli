// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestInnerFlagSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		flagType  string
		inputVal  string
		expected  any
		expectErr bool
	}{
		{"string", "string", "hello", "hello", false},
		{"int64", "int64", "42", int64(42), false},
		{"float64", "float64", "3.14", float64(3.14), false},
		{"bool", "bool", "true", true, false},
		{"invalid int", "int64", "not-a-number", nil, true},
		{"invalid float", "float64", "not-a-float", nil, true},
		{"invalid bool", "bool", "not-a-bool", nil, true},
		{"yaml map", "map", "key: value", map[string]any{"key": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			outerFlag := &Flag[map[string]any]{
				Name: "test-flag",
			}

			var innerFlag cli.Flag
			switch tt.flagType {
			case "string":
				innerFlag = &InnerFlag[string]{
					Name:       "test-flag.test-field",
					OuterFlag:  outerFlag,
					InnerField: "test_field",
				}
			case "int64":
				innerFlag = &InnerFlag[int64]{
					Name:       "test-flag.test-field",
					OuterFlag:  outerFlag,
					InnerField: "test_field",
				}
			case "float64":
				innerFlag = &InnerFlag[float64]{
					Name:       "test-flag.test-field",
					OuterFlag:  outerFlag,
					InnerField: "test_field",
				}
			case "bool":
				innerFlag = &InnerFlag[bool]{
					Name:       "test-flag.test-field",
					OuterFlag:  outerFlag,
					InnerField: "test_field",
				}
			case "map":
				innerFlag = &InnerFlag[map[string]any]{
					Name:       "test-flag.test-field",
					OuterFlag:  outerFlag,
					InnerField: "test_field",
				}
			}

			err := innerFlag.Set(innerFlag.Names()[0], tt.inputVal)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			actual, ok := outerFlag.Get().(map[string]any)["test_field"]
			assert.True(t, ok, "Field 'test_field' should exist in the map")
			assert.Equal(t, tt.expected, actual, "Expected %v (%T), got %v (%T)", tt.expected, tt.expected, actual, actual)
		})
	}
}

func TestInnerFlagValidator(t *testing.T) {
	t.Parallel()

	outerFlag := &Flag[map[string]any]{Name: "test-flag"}

	innerFlag := &InnerFlag[int64]{
		Name:       "test-flag.test-field",
		OuterFlag:  outerFlag,
		InnerField: "test_field",
		Validator: func(val int64) error {
			if val < 0 {
				return cli.Exit("Value must be non-negative", 1)
			}
			return nil
		},
	}

	err := innerFlag.Set(innerFlag.Name, "42")
	assert.NoError(t, err, "Expected no error for valid value, got: %v", err)

	err = innerFlag.Set(innerFlag.Name, "-5")
	assert.Error(t, err, "Expected error for invalid value, got none")
}

func TestWithInnerFlags(t *testing.T) {
	t.Parallel()

	outerFlag := &Flag[map[string]any]{Name: "outer"}
	innerFlag := &InnerFlag[string]{
		Name:       "outer.baz",
		InnerField: "baz",
	}

	cmd := WithInnerFlags(cli.Command{
		Name:  "test-command",
		Flags: []cli.Flag{outerFlag},
	}, map[string][]HasOuterFlag{
		"outer": {innerFlag},
	})

	assert.Len(t, cmd.Flags, 2, "Expected 2 flags, got %d", len(cmd.Flags))
	assert.Equal(t, outerFlag, cmd.Flags[0], "First flag should be outerFlag")
	assert.Equal(t, innerFlag, cmd.Flags[1], "Second flag should be innerFlag")
	assert.Same(t, outerFlag, innerFlag.OuterFlag, "innerFlag.OuterFlag should point to outerFlag")
}

func TestInnerFlagTypeNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		flag     cli.DocGenerationFlag
		expected string
	}{
		{"string", &InnerFlag[string]{}, "string"},
		{"int64", &InnerFlag[int64]{}, "int"},
		{"float64", &InnerFlag[float64]{}, "float"},
		{"bool", &InnerFlag[bool]{}, "boolean"},
		{"string slice", &InnerFlag[[]string]{}, "string"},
		{"date", &InnerFlag[DateValue]{}, "date"},
		{"datetime", &InnerFlag[DateTimeValue]{}, "datetime"},
		{"time", &InnerFlag[TimeValue]{}, "time"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			typeName := tt.flag.TypeName()
			assert.Equal(t, tt.expected, typeName, "Expected type name %q, got %q", tt.expected, typeName)
		})
	}
}

func TestInnerYamlHandling(t *testing.T) {
	t.Parallel()

	t.Run("Parse YAML to map", func(t *testing.T) {
		t.Parallel()

		outerFlag := &Flag[map[string]any]{Name: "outer"}
		innerFlag := &InnerFlag[map[string]any]{
			Name:       "outer.baz",
			OuterFlag:  outerFlag,
			InnerField: "baz",
		}

		err := innerFlag.Set(innerFlag.Name, "{name: test, value: 42}")
		assert.NoError(t, err)

		result, ok := outerFlag.Get().(map[string]any)
		assert.True(t, ok, "Expected map[string]any from outerFlag.Get()")
		yamlField, ok := result["baz"].(map[string]any)
		assert.True(t, ok, "Expected map[string]any, got %T", result["baz"])
		val := yamlField

		if ok {
			assert.Equal(t, map[string]any{"name": "test", "value": uint64(42)}, val)
		}
	})

	t.Run("Parse invalid YAML", func(t *testing.T) {
		t.Parallel()

		outerFlag := &Flag[map[string]any]{Name: "outer"}
		innerFlag := &InnerFlag[map[string]any]{
			Name:       "outer.baz",
			OuterFlag:  outerFlag,
			InnerField: "baz",
		}

		invalidYaml := `[not closed`
		err := innerFlag.Set(innerFlag.Name, invalidYaml)
		assert.Error(t, err)
	})

	t.Run("Set inner flags on map multiple times", func(t *testing.T) {
		t.Parallel()

		outerFlag := &Flag[map[string]any]{Name: "outer"}

		firstInnerFlag := &InnerFlag[string]{
			Name:       "outer.first-flag",
			OuterFlag:  outerFlag,
			InnerField: "first_field",
		}

		err := firstInnerFlag.Set(firstInnerFlag.Name, "first-value")
		assert.NoError(t, err)

		secondInnerFlag := &InnerFlag[int64]{
			Name:       "outer.second-flag",
			OuterFlag:  outerFlag,
			InnerField: "second_field",
		}

		err = secondInnerFlag.Set(secondInnerFlag.Name, "42")
		assert.NoError(t, err)

		result := outerFlag.Get().(map[string]any)
		assert.Equal(t, map[string]any{"first_field": "first-value", "second_field": int64(42)}, result)
	})

	t.Run("Set YAML and then inner flag", func(t *testing.T) {
		t.Parallel()

		outerFlag := &Flag[map[string]any]{Name: "outer"}

		err := outerFlag.Set(outerFlag.Name, `{existing: value, another: field}`)
		assert.NoError(t, err)

		innerFlag := &InnerFlag[string]{
			Name:       "outer.inner-flag",
			OuterFlag:  outerFlag,
			InnerField: "new_field",
		}

		err = innerFlag.Set(innerFlag.Name, "inner-value")
		assert.NoError(t, err)

		result := outerFlag.Get().(map[string]any)
		assert.Equal(t, map[string]any{
			"existing":  "value",
			"another":   "field",
			"new_field": "inner-value",
		}, result)
	})
}

func TestInnerFlagWithSliceType(t *testing.T) {
	t.Parallel()

	t.Run("Setting inner flags on slice of maps", func(t *testing.T) {
		t.Parallel()

		outerFlag := &Flag[[]map[string]any]{Name: "outer"}

		firstInnerFlag := &InnerFlag[string]{
			Name:       "outer.name-flag",
			OuterFlag:  outerFlag,
			InnerField: "name",
		}

		err := firstInnerFlag.Set(firstInnerFlag.Name, "item1")
		assert.NoError(t, err)

		secondInnerFlag := &InnerFlag[int64]{
			Name:       "outer.count-flag",
			OuterFlag:  outerFlag,
			InnerField: "count",
		}

		err = secondInnerFlag.Set(secondInnerFlag.Name, "42")
		assert.NoError(t, err)

		err = firstInnerFlag.Set(firstInnerFlag.Name, "item2")
		assert.NoError(t, err)

		result := outerFlag.Get().([]map[string]any)

		assert.Equal(t, []map[string]any{
			{"name": "item1", "count": int64(42)},
			{"name": "item2"},
		}, result)
		assert.Nil(t, result[1]["count"], "Second item should not have count field")
	})

	t.Run("Appending to existing slice", func(t *testing.T) {
		t.Parallel()

		outerFlag := &Flag[[]map[string]any]{Name: "outer"}
		err := outerFlag.Set(outerFlag.Name, `{name: initial}`)
		assert.NoError(t, err)

		modifyFlag := &InnerFlag[string]{
			Name:       "outer.value-flag",
			OuterFlag:  outerFlag,
			InnerField: "value",
		}

		err = modifyFlag.Set(modifyFlag.Name, "updated")
		assert.NoError(t, err)

		newItemFlag := &InnerFlag[string]{
			Name:       "outer.name-flag",
			OuterFlag:  outerFlag,
			InnerField: "name",
		}

		err = newItemFlag.Set(newItemFlag.Name, "second")
		assert.NoError(t, err)

		result := outerFlag.Get().([]map[string]any)
		assert.Equal(t, []map[string]any{
			{"name": "initial", "value": "updated"},
			{"name": "second"},
		}, result)
	})
}
