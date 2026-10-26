// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

type InnerFlag[T any, C codec[T]] struct {
	Name        string
	DefaultText string
	Usage       string
	Aliases     []string
	Validator   func(T) error

	OuterFlag   cli.Flag
	InnerField  string
	DataAliases []string
}

type InnerString = InnerFlag[string, textCodec]

var (
	_ HasOuterFlag                    = (*InnerString)(nil)
	_ cli.DocGenerationMultiValueFlag = (*InnerString)(nil)
)

type HasOuterFlag interface {
	cli.Flag
	SetOuterFlag(cli.Flag)
	GetOuterFlag() cli.Flag
	GetInnerField() string
	GetDataAliases() []string
}

type fieldSetter interface {
	setInnerField(field string, value json.RawMessage) error
}

func (f *InnerFlag[T, C]) GetDataAliases() []string {
	return f.DataAliases
}

func (f *InnerFlag[T, C]) GetInnerField() string {
	return f.InnerField
}

func (f *InnerFlag[T, C]) SetOuterFlag(flag cli.Flag) {
	f.OuterFlag = flag
}

func (f *InnerFlag[T, C]) GetOuterFlag() cli.Flag {
	return f.OuterFlag
}

func (f *InnerFlag[T, C]) PreParse() error {
	return nil
}

func (f *InnerFlag[T, C]) PostParse() error {
	return nil
}

func (f *InnerFlag[T, C]) Set(name string, input string) error {
	var c C
	var empty T
	parsed, err := c.parse(empty, input)
	if err != nil {
		return err
	}
	if f.Validator != nil {
		if err := f.Validator(parsed); err != nil {
			return err
		}
	}
	outer, ok := f.OuterFlag.(fieldSetter)
	if !ok {
		return fmt.Errorf("cannot set an inner field on %v", f.OuterFlag)
	}
	encoded, err := json.Marshal(parsed)
	if err != nil {
		return err
	}
	return outer.setInnerField(f.InnerField, encoded)
}

func (f *InnerFlag[T, C]) Get() any {
	var empty T
	return empty
}

func (f *InnerFlag[T, C]) String() string {
	return cli.FlagStringer(f)
}

func (f *InnerFlag[T, C]) IsSet() bool {
	return false
}

func (f *InnerFlag[T, C]) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

func (f *InnerFlag[T, C]) TakesValue() bool {
	var c C
	return !c.traits().isBool
}

func (f *InnerFlag[T, C]) GetUsage() string {
	return f.Usage
}

func (f *InnerFlag[T, C]) GetValue() string {
	return ""
}

func (f *InnerFlag[T, C]) GetDefaultText() string {
	return f.DefaultText
}

func (f *InnerFlag[T, C]) GetEnvVars() []string {
	return nil
}

func (f *InnerFlag[T, C]) IsDefaultVisible() bool {
	return false
}

func (f *InnerFlag[T, C]) TypeName() string {
	var c C
	return c.traits().typeName
}

func (f *InnerFlag[T, C]) IsMultiValueFlag() bool {
	return false
}

func (f *InnerFlag[T, C]) IsBoolFlag() bool {
	var c C
	return c.traits().isBool
}

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
