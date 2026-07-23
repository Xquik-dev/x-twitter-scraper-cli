package requestflag

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestFlagMetadataAndTypeHelpers(t *testing.T) {
	t.Parallel()

	flag := &Flag[string]{
		Name:        "value",
		Aliases:     []string{"v"},
		Category:    "input",
		Default:     "default",
		DefaultText: "default text",
		FileInput:   true,
		HideDefault: true,
		Hidden:      true,
		Required:    true,
		Sources:     cli.EnvVars("XQUIK_COVERAGE_VALUE"),
		Usage:       "A value.",
	}

	assert.Equal(t, []string{"value", "v"}, flag.Names())
	assert.False(t, flag.IsVisible())
	assert.Equal(t, "input", flag.GetCategory())
	flag.SetCategory("request")
	assert.Equal(t, "request", flag.GetCategory())
	assert.True(t, flag.IsRequired())
	assert.True(t, flag.IsRequiredAsFlagOrStdin())
	assert.True(t, flag.TakesValue())
	assert.Equal(t, "A value.", flag.GetUsage())
	assert.Empty(t, flag.GetValue())
	assert.Equal(t, "default text", flag.GetDefaultText())
	assert.Equal(t, []string{"XQUIK_COVERAGE_VALUE"}, flag.GetEnvVars())
	assert.False(t, flag.IsDefaultVisible())
	assert.Equal(t, "string", flag.TypeName())
	assert.False(t, flag.IsMultiValueFlag())
	assert.False(t, flag.IsBoolFlag())
	assert.Zero(t, flag.Count())
	assert.True(t, flag.IsLocal())
	assert.True(t, flag.IsFileInput())
	assert.Contains(t, flag.String(), "value")

	require.NoError(t, flag.PreParse())
	assert.Equal(t, "default", flag.GetValue())
	require.NoError(t, flag.Set(flag.Name, "changed"))
	assert.Equal(t, 1, flag.Count())

	constFlag := &Flag[string]{Const: true, Required: true}
	assert.True(t, constFlag.IsSet())
	assert.False(t, constFlag.IsRequired())
	assert.False(t, constFlag.IsRequiredAsFlagOrStdin())

	requestFlag := &Flag[string]{BodyPath: "value", Required: true}
	assert.False(t, requestFlag.IsRequired())
	assert.True(t, requestFlag.IsRequiredAsFlagOrStdin())

	boolFlag := &Flag[bool]{Default: false}
	assert.False(t, boolFlag.TakesValue())
	assert.True(t, boolFlag.IsBoolFlag())
	assert.True(t, (&cliValue[bool]{value: true}).IsBoolFlag())

	mapFlag := &Flag[map[string]any]{}
	assert.True(t, mapFlag.IsMultiValueFlag())
	assert.Equal(t, "string=any", mapFlag.TypeName())
	assert.Equal(t, "int", (&Flag[*int64]{}).TypeName())
	assert.Empty(t, (&Flag[any]{}).TypeName())
	assert.Equal(t, "value", *Ptr("value"))
}

func TestFlagPostParseReadsEmptyAndInvalidEnvironmentValues(t *testing.T) {
	t.Setenv("XQUIK_COVERAGE_STRING", "")
	t.Setenv("XQUIK_COVERAGE_BOOL", "")
	t.Setenv("XQUIK_COVERAGE_INT", "invalid")
	t.Setenv("XQUIK_COVERAGE_POINTER", "")

	stringFlag := &Flag[string]{
		Name:    "string",
		Default: "fallback",
		Sources: cli.EnvVars("XQUIK_COVERAGE_STRING"),
	}
	require.NoError(t, stringFlag.PreParse())
	require.NoError(t, stringFlag.PostParse())
	assert.Equal(t, "", stringFlag.Get())

	boolFlag := &Flag[bool]{
		Name:    "bool",
		Default: true,
		Sources: cli.EnvVars("XQUIK_COVERAGE_BOOL"),
	}
	require.NoError(t, boolFlag.PreParse())
	require.NoError(t, boolFlag.PostParse())
	assert.Equal(t, false, boolFlag.Get())

	pointerFlag := &Flag[*string]{
		Name:    "pointer",
		Sources: cli.EnvVars("XQUIK_COVERAGE_POINTER"),
	}
	require.NoError(t, pointerFlag.PreParse())
	require.NoError(t, pointerFlag.PostParse())
	require.NotNil(t, pointerFlag.Get())
	assert.Equal(t, "", *pointerFlag.Get().(*string))

	intFlag := &Flag[int64]{
		Name:    "int",
		Sources: cli.EnvVars("XQUIK_COVERAGE_INT"),
	}
	require.NoError(t, intFlag.PreParse())
	assert.ErrorContains(t, intFlag.PostParse(), "could not parse")

	alreadySet := &Flag[string]{
		Name:       "set",
		hasBeenSet: true,
		Sources:    cli.EnvVars("XQUIK_COVERAGE_STRING"),
	}
	require.NoError(t, alreadySet.PostParse())
}

func TestFlagValidationAndRequiredInputHelpers(t *testing.T) {
	t.Parallel()

	invalidDefault := &Flag[string]{
		Default: "invalid",
		Validator: func(value string) error {
			return errors.New("invalid default")
		},
	}
	assert.EqualError(t, invalidDefault.PreParse(), "invalid default")

	required := &Flag[string]{Name: "required", Required: true}
	bodyRoot := &Flag[any]{Name: "root", BodyRoot: true, Required: true}
	bodyField := &Flag[string]{
		Name:     "body-field",
		BodyPath: "body_field",
		Required: true,
	}
	command := &cli.Command{Flags: []cli.Flag{required, bodyRoot, bodyField}}

	missing := GetMissingRequiredFlags(command, nil)
	assert.ElementsMatch(t, []cli.Flag{required, bodyRoot, bodyField}, missing)

	missing = GetMissingRequiredFlags(command, map[string]any{"body_field": "set"})
	assert.Equal(t, []cli.Flag{required}, missing)

	missing = GetMissingRequiredFlags(command, "body")
	assert.ElementsMatch(t, []cli.Flag{required, bodyField}, missing)

	require.NoError(t, required.Set(required.Name, "set"))
	missing = GetMissingRequiredFlags(command, "body")
	assert.Equal(t, []cli.Flag{bodyField}, missing)
}

func TestFormatForFlagSet(t *testing.T) {
	t.Parallel()

	value, err := formatForFlagSet("plain")
	require.NoError(t, err)
	assert.Equal(t, "plain", value)

	value, err = formatForFlagSet(true)
	require.NoError(t, err)
	assert.Equal(t, "true", value)

	value, err = formatForFlagSet(map[string]any{"key": "value"})
	require.NoError(t, err)
	assert.JSONEq(t, `{"key":"value"}`, value)

	_, err = formatForFlagSet(make(chan int))
	assert.ErrorContains(t, err, "cannot format value")
}

func TestInnerFlagMetadataAndValidation(t *testing.T) {
	t.Parallel()

	outer := &Flag[map[string]any]{Name: "outer", Aliases: []string{"o"}}
	inner := &InnerFlag[string]{
		Name:        "outer.value",
		Aliases:     []string{"o.value"},
		DefaultText: "default",
		InnerField:  "value",
		OuterFlag:   outer,
		Usage:       "Nested value.",
	}

	require.NoError(t, inner.PreParse())
	require.NoError(t, inner.PostParse())
	assert.Empty(t, inner.Get())
	assert.Contains(t, inner.String(), "outer.value")
	assert.False(t, inner.IsSet())
	assert.True(t, inner.TakesValue())
	assert.Equal(t, "Nested value.", inner.GetUsage())
	assert.Empty(t, inner.GetValue())
	assert.Equal(t, "default", inner.GetDefaultText())
	assert.Nil(t, inner.GetEnvVars())
	assert.False(t, inner.IsDefaultVisible())
	assert.Equal(t, "string", inner.TypeName())
	assert.False(t, inner.IsMultiValueFlag())
	assert.False(t, inner.IsBoolFlag())

	boolInner := &InnerFlag[bool]{}
	assert.False(t, boolInner.TakesValue())
	assert.True(t, boolInner.IsBoolFlag())
	assert.Empty(t, (&InnerFlag[any]{}).TypeName())
	assert.Equal(t, "string=any", (&InnerFlag[map[string]any]{}).TypeName())
	assert.Equal(t, "int", (&InnerFlag[*int64]{}).TypeName())

	unsupportedOuter := &InnerFlag[string]{
		InnerField: "value",
		OuterFlag:  &cli.StringFlag{Name: "outer"},
	}
	assert.ErrorContains(t, unsupportedOuter.Set("outer.value", "value"), "Cannot set inner field")
}

func TestInnerFlagAssemblyChecks(t *testing.T) {
	t.Parallel()

	empty := cli.Command{Name: "empty"}
	assert.Equal(t, empty, WithInnerFlags(empty, nil))

	outer := &Flag[map[string]any]{Name: "outer", Aliases: []string{"o"}}
	inner := &InnerFlag[string]{Name: "o.value", InnerField: "value"}
	command := WithInnerFlags(cli.Command{
		Name:  "valid",
		Flags: []cli.Flag{outer},
	}, map[string][]HasOuterFlag{"o": {inner}})
	require.NoError(t, CheckInnerFlags(command))

	missingOuter := cli.Command{
		Flags: []cli.Flag{&InnerFlag[string]{Name: "outer.value"}},
	}
	assert.ErrorContains(t, CheckInnerFlags(missingOuter), "missing an outer flag")

	badPrefix := cli.Command{
		Flags: []cli.Flag{&InnerFlag[string]{
			Name:      "wrong.value",
			OuterFlag: outer,
		}},
	}
	assert.ErrorContains(t, CheckInnerFlags(badPrefix), "must start")

	assert.Panics(t, func() {
		WithInnerFlags(cli.Command{Flags: []cli.Flag{outer}}, map[string][]HasOuterFlag{
			"missing": {inner},
		})
	})
}
