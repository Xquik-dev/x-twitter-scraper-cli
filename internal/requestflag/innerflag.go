// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/urfave/cli/v3"
)

// InnerFlag sets nested values. For example, --foo.baz sets foo.baz.
type InnerFlag[
	T []any | []map[string]any | []DateTimeValue | []DateValue | []TimeValue | []string |
		[]float64 | []int64 | []bool | any | map[string]any | DateTimeValue | DateValue | TimeValue |
		string | float64 | int64 | bool |
		*string | *float64 | *int64 | *bool | *DateTimeValue | *DateValue | *TimeValue,
] struct {
	Name        string
	DefaultText string
	Usage       string
	Aliases     []string
	Validator   func(T) error

	OuterFlag   cli.Flag
	InnerField  string
	DataAliases []string

	// OuterIsArrayOfObjects selects a list seed for untyped nullable schemas.
	OuterIsArrayOfObjects bool
}

// GetDataAliases returns the aliases recognized when parsing inner field keys from piped or flag YAML.
func (f *InnerFlag[T]) GetDataAliases() []string {
	return f.DataAliases
}

// GetInnerField returns the canonical API field name.
func (f *InnerFlag[T]) GetInnerField() string {
	return f.InnerField
}

type HasOuterFlag interface {
	cli.Flag
	SetOuterFlag(cli.Flag)
	GetOuterFlag() cli.Flag
	GetInnerField() string
	GetDataAliases() []string
}

func (f *InnerFlag[T]) SetOuterFlag(flag cli.Flag) {
	f.OuterFlag = flag
}

func (f *InnerFlag[T]) GetOuterFlag() cli.Flag {
	return f.OuterFlag
}

var _ cli.Flag = (*InnerFlag[any])(nil)

func (f *InnerFlag[T]) PreParse() error {
	return nil
}

func (f *InnerFlag[T]) PostParse() error {
	return nil
}

func (f *InnerFlag[T]) Set(name string, rawVal string) error {
	if parsedValue, err := parseCLIArg[T](rawVal); err != nil {
		return err
	} else {
		if f.Validator != nil {
			if err := f.Validator(parsedValue); err != nil {
				return err
			}
		}

		if seeder, ok := f.OuterFlag.(InnerFieldSeeder); ok {
			seeder.SeedInnerCollection(f.OuterIsArrayOfObjects)
		}

		if settableInnerField, ok := f.OuterFlag.(SettableInnerField); ok {
			settableInnerField.SetInnerField(f.InnerField, parsedValue)
		} else {
			return fmt.Errorf("cannot set an inner field on %v", f.OuterFlag)
		}
		return nil
	}
}

func (f *InnerFlag[T]) Get() any {
	var zeroValue T
	return zeroValue
}

func (f *InnerFlag[T]) String() string {
	return cli.FlagStringer(f)
}

func (f *InnerFlag[T]) IsSet() bool {
	return false
}

func (f *InnerFlag[T]) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

var _ cli.DocGenerationFlag = (*InnerFlag[any])(nil)

func (f *InnerFlag[T]) TakesValue() bool {
	var t T
	return reflect.TypeOf(t) == nil || reflect.TypeOf(t).Kind() != reflect.Bool
}

func (f *InnerFlag[T]) GetUsage() string {
	return f.Usage
}

func (f *InnerFlag[T]) GetValue() string {
	return ""
}

func (f *InnerFlag[T]) GetDefaultText() string {
	return f.DefaultText
}

func (f *InnerFlag[T]) GetEnvVars() []string {
	return nil
}

func (f *InnerFlag[T]) IsDefaultVisible() bool {
	return false
}

func (f *InnerFlag[T]) TypeName() string {
	var zeroValue T
	ty := reflect.TypeOf(zeroValue)
	if ty == nil {
		return ""
	}
	if ty.Kind() == reflect.Pointer {
		ty = ty.Elem()
	}

	getTypeName := func(t reflect.Type) string {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "int"
		case reflect.Float32, reflect.Float64:
			return "float"
		case reflect.Bool:
			return "boolean"
		case reflect.String:
			switch t.Name() {
			case "DateTimeValue":
				return "datetime"
			case "DateValue":
				return "date"
			case "TimeValue":
				return "time"
			default:
				return "string"
			}
		default:
			if t.Name() == "" {
				return "any"
			}
			return strings.ToLower(t.Name())
		}
	}

	switch ty.Kind() {
	case reflect.Slice:
		elemType := ty.Elem()
		return getTypeName(elemType)
	case reflect.Map:
		keyType := ty.Key()
		valueType := ty.Elem()
		return fmt.Sprintf("%s=%s", getTypeName(keyType), getTypeName(valueType))
	default:
		return getTypeName(ty)
	}
}

var _ cli.DocGenerationMultiValueFlag = (*InnerFlag[any])(nil)

func (f *InnerFlag[T]) IsMultiValueFlag() bool {
	return false
}

func (f *InnerFlag[T]) IsBoolFlag() bool {
	var zeroValue T
	_, isBool := any(zeroValue).(bool)
	return isBool
}

// WithInnerFlags attaches nested flags to their outer flags.
func WithInnerFlags(cmd cli.Command, innerFlagMap map[string][]HasOuterFlag) cli.Command {
	if len(innerFlagMap) == 0 {
		return cmd
	}

	unusedInnerFlagKeys := make(map[string]struct{})
	for name := range innerFlagMap {
		unusedInnerFlagKeys[name] = struct{}{}
	}

	updatedFlags := make([]cli.Flag, 0, len(cmd.Flags))
	for _, flag := range cmd.Flags {
		updatedFlags = append(updatedFlags, flag)
		for _, name := range flag.Names() {
			innerFlags, hasInnerFlags := innerFlagMap[name]
			if !hasInnerFlags {
				continue
			}

			delete(unusedInnerFlagKeys, name)

			for _, innerFlag := range innerFlags {
				innerFlag.SetOuterFlag(flag)
				updatedFlags = append(updatedFlags, innerFlag)
			}
		}
	}

	// A missing outer flag indicates invalid generated wiring.
	if len(unusedInnerFlagKeys) > 0 {
		unusedKeys := make([]string, 0, len(unusedInnerFlagKeys))
		for key := range unusedInnerFlagKeys {
			unusedKeys = append(unusedKeys, key)
		}
		panic(fmt.Sprintf("inner flags are missing outer flags: %v", unusedKeys))
	}

	result := cmd
	result.Flags = updatedFlags
	return result
}

// CheckInnerFlags validates outer references and --foo.baz prefixes.
func CheckInnerFlags(cmd cli.Command) error {
	var errors []string
	for _, flag := range cmd.Flags {
		if innerFlag, ok := flag.(HasOuterFlag); ok {
			outerFlag := innerFlag.GetOuterFlag()
			if outerFlag == nil {
				errors = append(errors, fmt.Sprintf("inner flag %s is missing an outer flag", flag.Names()))
				continue
			}

			innerFlagName := flag.Names()[0]
			valid := false
			for _, outerName := range outerFlag.Names() {
				if strings.HasPrefix(innerFlagName, outerName+".") {
					valid = true
					break
				}
			}

			if !valid {
				errors = append(errors, fmt.Sprintf("inner flag %s must start with one of its outer flag's names followed by a dot", innerFlagName))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "; "))
	}
	return nil
}
