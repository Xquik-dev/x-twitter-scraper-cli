// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package requestflag defines command-line flags that set the query, header,
// path and body fields of a request, on top of
// https://pkg.go.dev/github.com/urfave/cli/v3.
package requestflag

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/urfave/cli/v3"
)

type Flag[T any, C codec[T]] struct {
	Name        string
	Category    string
	DefaultText string
	HideDefault bool
	Usage       string
	Sources     cli.ValueSourceChain
	Required    bool
	Hidden      bool
	Default     T
	Aliases     []string
	Validator   func(T) error

	QueryPath  string
	HeaderPath string
	BodyPath   string
	PathParam  string

	Const bool

	DataAliases []string

	count      int
	hasBeenSet bool
	applied    bool
	hasValue   bool
	value      T
}

type (
	String     = Flag[string, textCodec]
	Int        = Flag[int64, intCodec]
	Bool       = Flag[bool, boolCodec]
	StringList = Flag[[]string, textListCodec]
	JSON       = Flag[json.RawMessage, jsonCodec]
	JSONList   = Flag[[]json.RawMessage, jsonListCodec]
	Object     = Flag[map[string]json.RawMessage, objectCodec]
	ObjectList = Flag[[]map[string]json.RawMessage, objectListCodec]
)

var (
	_ cli.CategorizableFlag           = (*String)(nil)
	_ cli.RequiredFlag                = (*String)(nil)
	_ cli.DocGenerationMultiValueFlag = (*String)(nil)
	_ cli.Countable                   = (*String)(nil)
	_ cli.LocalFlag                   = (*String)(nil)
	_ InRequest                       = (*String)(nil)
	_ fieldSetter                     = (*ObjectList)(nil)
)

type InRequest interface {
	GetQueryPath() string
	GetHeaderPath() string
	GetBodyPath() string
	GetPathParam() string
	GetDataAliases() []string
	requestParam() apiquery.Param
	requestJSON() (json.RawMessage, error)
}

type RequiredFlagOrStdin interface {
	IsRequiredAsFlagOrStdin() bool
}

type RequestContents struct {
	Queries map[string]apiquery.Param
	Headers map[string]apiquery.Param
	Body    map[string]json.RawMessage
}

func (f *Flag[T, C]) GetQueryPath() string     { return f.QueryPath }
func (f *Flag[T, C]) GetHeaderPath() string    { return f.HeaderPath }
func (f *Flag[T, C]) GetBodyPath() string      { return f.BodyPath }
func (f *Flag[T, C]) GetPathParam() string     { return f.PathParam }
func (f *Flag[T, C]) GetDataAliases() []string { return f.DataAliases }

func (f *Flag[T, C]) current() T {
	if f.hasValue {
		return f.value
	}
	return f.Default
}

func (f *Flag[T, C]) requestParam() apiquery.Param {
	var c C
	return c.param(f.current())
}

func (f *Flag[T, C]) requestJSON() (json.RawMessage, error) {
	return json.Marshal(f.current())
}

func (f *Flag[T, C]) setInnerField(field string, value json.RawMessage) error {
	var c C
	updated, err := c.setField(f.current(), field, value)
	if err != nil {
		return fmt.Errorf("cannot set an inner field on flag %s: %w", f.Name, err)
	}
	f.value = updated
	f.hasValue = true
	f.hasBeenSet = true
	return nil
}

func (f *Flag[T, C]) PreParse() error {
	if f.Validator != nil {
		if err := f.Validator(f.Default); err != nil {
			return err
		}
	}
	f.value = f.Default
	f.hasValue = true
	f.applied = true
	return nil
}

func (f *Flag[T, C]) PostParse() error {
	if f.hasBeenSet {
		return nil
	}
	input, source, found := f.Sources.LookupWithSource()
	if !found {
		return nil
	}
	var c C
	traits := c.traits()
	if input != "" || traits.acceptsEmpty {
		if err := f.Set(f.Name, input); err != nil {
			return fmt.Errorf("could not parse %q as %s value from %s for flag %s: %s", input, traits.typeName, source, f.Name, err)
		}
	} else if traits.isBool {
		if err := f.Set(f.Name, "false"); err != nil {
			return err
		}
	}
	f.hasBeenSet = true
	return nil
}

func (f *Flag[T, C]) Set(name string, input string) error {
	if !f.applied {
		if err := f.PreParse(); err != nil {
			return err
		}
	}
	f.count++
	current := f.value
	if f.count == 1 {
		var empty T
		current = empty
	}
	var c C
	parsed, err := c.parse(current, input)
	if err != nil {
		return err
	}
	f.value = parsed
	f.hasValue = true
	f.hasBeenSet = true
	if f.Validator != nil {
		return f.Validator(f.value)
	}
	return nil
}

func (f *Flag[T, C]) Get() any {
	return f.current()
}

func (f *Flag[T, C]) String() string {
	return cli.FlagStringer(f)
}

func (f *Flag[T, C]) IsSet() bool {
	return f.hasBeenSet || f.Const
}

func (f *Flag[T, C]) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

func (f *Flag[T, C]) IsVisible() bool {
	return !f.Hidden
}

func (f *Flag[T, C]) GetCategory() string {
	return f.Category
}

func (f *Flag[T, C]) SetCategory(category string) {
	f.Category = category
}

func (f *Flag[T, C]) IsRequired() bool {
	if f.Const {
		return false
	}
	if f.BodyPath != "" || f.PathParam != "" || f.QueryPath != "" || f.HeaderPath != "" {
		return false
	}
	return f.Required
}

func (f *Flag[T, C]) IsRequiredAsFlagOrStdin() bool {
	return !f.Const && f.Required
}

func (f *Flag[T, C]) TakesValue() bool {
	var c C
	return !c.traits().isBool
}

func (f *Flag[T, C]) GetUsage() string {
	return f.Usage
}

func (f *Flag[T, C]) GetValue() string {
	if !f.hasValue {
		return ""
	}
	var c C
	return c.display(f.value)
}

func (f *Flag[T, C]) GetDefaultText() string {
	return f.DefaultText
}

func (f *Flag[T, C]) GetEnvVars() []string {
	return f.Sources.EnvKeys()
}

func (f *Flag[T, C]) IsDefaultVisible() bool {
	return !f.HideDefault
}

func (f *Flag[T, C]) TypeName() string {
	var c C
	return c.traits().typeName
}

func (f *Flag[T, C]) IsMultiValueFlag() bool {
	var c C
	return c.traits().multiValue
}

func (f *Flag[T, C]) IsBoolFlag() bool {
	var c C
	return c.traits().isBool
}

func (f *Flag[T, C]) Count() int {
	return f.count
}

func (f *Flag[T, C]) IsLocal() bool {
	return true
}

func ApplyStdinDataToFlags(cmd *cli.Command, data map[string]json.RawMessage) error {
	for _, flag := range cmd.Flags {
		if flag.IsSet() {
			continue
		}
		raw, found := stdinValue(flag, data)
		if !found {
			continue
		}
		input, err := formatForFlagSet(raw)
		if err != nil {
			return fmt.Errorf("cannot format piped value for flag %q: %w", flag.Names()[0], err)
		}
		if err := flag.Set(flag.Names()[0], input); err != nil {
			return fmt.Errorf("cannot set flag %q from piped data: %w", flag.Names()[0], err)
		}
	}
	return nil
}

func stdinValue(flag cli.Flag, data map[string]json.RawMessage) (json.RawMessage, bool) {
	if inner, ok := flag.(HasOuterFlag); ok {
		outer, ok := inner.GetOuterFlag().(InRequest)
		if !ok || outer.GetBodyPath() == "" {
			return nil, false
		}
		nested, ok := JSONObject(data[outer.GetBodyPath()])
		if !ok {
			return nil, false
		}
		return lookup(nested, inner.GetInnerField(), inner.GetDataAliases())
	}
	inRequest, ok := flag.(InRequest)
	if !ok {
		return nil, false
	}
	for _, path := range []string{inRequest.GetQueryPath(), inRequest.GetHeaderPath(), inRequest.GetPathParam()} {
		if path == "" {
			continue
		}
		if raw, found := lookup(data, path, inRequest.GetDataAliases()); found {
			return raw, true
		}
	}
	return nil, false
}

func lookup(data map[string]json.RawMessage, key string, aliases []string) (json.RawMessage, bool) {
	for _, name := range append([]string{key}, aliases...) {
		if name == "" {
			continue
		}
		if raw, found := data[name]; found {
			return raw, true
		}
	}
	return nil, false
}

func formatForFlagSet(raw json.RawMessage) (string, error) {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text, nil
	}
	compact, err := compactJSON(raw)
	if err != nil {
		return "", err
	}
	return string(compact), nil
}

func ExtractRequestContents(cmd *cli.Command) (RequestContents, error) {
	contents := RequestContents{
		Queries: map[string]apiquery.Param{},
		Headers: map[string]apiquery.Param{},
		Body:    map[string]json.RawMessage{},
	}
	for _, flag := range cmd.Flags {
		inRequest, ok := flag.(InRequest)
		if !ok || !flag.IsSet() {
			continue
		}
		if path := inRequest.GetQueryPath(); path != "" {
			contents.Queries[path] = inRequest.requestParam()
		}
		if path := inRequest.GetHeaderPath(); path != "" {
			contents.Headers[path] = inRequest.requestParam()
		}
		if path := inRequest.GetBodyPath(); path != "" {
			raw, err := inRequest.requestJSON()
			if err != nil {
				return RequestContents{}, fmt.Errorf("cannot encode flag %q: %w", flag.Names()[0], err)
			}
			contents.Body[path] = raw
		}
	}
	return contents, nil
}

func GetMissingRequiredFlags(cmd *cli.Command, body map[string]json.RawMessage) []cli.Flag {
	missing := []cli.Flag{}
	for _, flag := range cmd.Flags {
		if flag.IsSet() {
			continue
		}
		if required, ok := flag.(cli.RequiredFlag); ok && required.IsRequired() {
			missing = append(missing, flag)
			continue
		}
		if required, ok := flag.(RequiredFlagOrStdin); !ok || !required.IsRequiredAsFlagOrStdin() {
			continue
		}
		if inRequest, ok := flag.(InRequest); ok && inRequest.GetBodyPath() != "" {
			if _, found := body[inRequest.GetBodyPath()]; found {
				continue
			}
		}
		missing = append(missing, flag)
	}
	return missing
}

func JSONObject(raw json.RawMessage) (map[string]json.RawMessage, bool) {
	if !strings.HasPrefix(strings.TrimSpace(string(raw)), "{") {
		return nil, false
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, false
	}
	return fields, true
}
