// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/goccy/go-yaml"
)

type codec[T any] interface {
	parse(current T, input string) (T, error)
	param(value T) apiquery.Param
	display(value T) string
	setField(current T, field string, value json.RawMessage) (T, error)
	traits() traits
}

type traits struct {
	typeName     string
	isBool       bool
	multiValue   bool
	acceptsEmpty bool
}

var errNotAnObject = errors.New("the flag value is not an object")

type textCodec struct{}

func (textCodec) parse(_ string, input string) (string, error) { return input, nil }
func (textCodec) param(value string) apiquery.Param            { return apiquery.Scalar(value) }
func (textCodec) display(value string) string                  { return value }
func (textCodec) setField(current string, _ string, _ json.RawMessage) (string, error) {
	return current, errNotAnObject
}
func (textCodec) traits() traits { return traits{typeName: "string", acceptsEmpty: true} }

type intCodec struct{}

func (intCodec) parse(_ int64, input string) (int64, error) { return strconv.ParseInt(input, 0, 64) }
func (intCodec) param(value int64) apiquery.Param {
	return apiquery.Scalar(strconv.FormatInt(value, 10))
}
func (intCodec) display(value int64) string { return strconv.FormatInt(value, 10) }
func (intCodec) setField(current int64, _ string, _ json.RawMessage) (int64, error) {
	return current, errNotAnObject
}
func (intCodec) traits() traits { return traits{typeName: "int"} }

type boolCodec struct{}

func (boolCodec) parse(_ bool, input string) (bool, error) { return strconv.ParseBool(input) }
func (boolCodec) param(value bool) apiquery.Param {
	return apiquery.Scalar(strconv.FormatBool(value))
}
func (boolCodec) display(value bool) string { return strconv.FormatBool(value) }
func (boolCodec) setField(current bool, _ string, _ json.RawMessage) (bool, error) {
	return current, errNotAnObject
}
func (boolCodec) traits() traits { return traits{typeName: "boolean", isBool: true} }

type textListCodec struct{}

func (textListCodec) parse(current []string, input string) ([]string, error) {
	return append(current, input), nil
}
func (textListCodec) param(value []string) apiquery.Param { return apiquery.List(value) }
func (textListCodec) display(value []string) string       { return fmt.Sprint(value) }
func (textListCodec) setField(current []string, _ string, _ json.RawMessage) ([]string, error) {
	return current, errNotAnObject
}
func (textListCodec) traits() traits { return traits{typeName: "string", multiValue: true} }

type jsonCodec struct{}

func (jsonCodec) parse(_ json.RawMessage, input string) (json.RawMessage, error) {
	return parseJSON(input)
}
func (jsonCodec) param(value json.RawMessage) apiquery.Param {
	return apiquery.Scalar(jsonText(value))
}
func (jsonCodec) display(value json.RawMessage) string { return string(value) }
func (jsonCodec) setField(current json.RawMessage, _ string, _ json.RawMessage) (json.RawMessage, error) {
	return current, errNotAnObject
}
func (jsonCodec) traits() traits { return traits{} }

type jsonListCodec struct{}

func (jsonListCodec) parse(current []json.RawMessage, input string) ([]json.RawMessage, error) {
	parsed, err := parseJSON(input)
	if err != nil {
		return current, err
	}
	return append(current, parsed), nil
}
func (jsonListCodec) param(value []json.RawMessage) apiquery.Param {
	texts := make([]string, len(value))
	for i, item := range value {
		texts[i] = jsonText(item)
	}
	return apiquery.List(texts)
}
func (jsonListCodec) display(value []json.RawMessage) string {
	if len(value) == 0 {
		return ""
	}
	return encodeJSON(value)
}
func (jsonListCodec) setField(current []json.RawMessage, _ string, _ json.RawMessage) ([]json.RawMessage, error) {
	return current, errNotAnObject
}
func (jsonListCodec) traits() traits { return traits{typeName: "any", multiValue: true} }

type objectCodec struct{}

func (objectCodec) parse(_ map[string]json.RawMessage, input string) (map[string]json.RawMessage, error) {
	return parseObject(input)
}
func (objectCodec) param(value map[string]json.RawMessage) apiquery.Param {
	return apiquery.Scalar(encodeJSON(value))
}
func (objectCodec) display(value map[string]json.RawMessage) string {
	if len(value) == 0 {
		return ""
	}
	return encodeJSON(value)
}
func (objectCodec) setField(current map[string]json.RawMessage, field string, value json.RawMessage) (map[string]json.RawMessage, error) {
	if current == nil {
		return map[string]json.RawMessage{field: value}, nil
	}
	current[field] = value
	return current, nil
}
func (objectCodec) traits() traits { return traits{typeName: "string=any", multiValue: true} }

type objectListCodec struct{}

func (objectListCodec) parse(current []map[string]json.RawMessage, input string) ([]map[string]json.RawMessage, error) {
	parsed, err := parseObject(input)
	if err != nil {
		return current, err
	}
	return append(current, parsed), nil
}
func (objectListCodec) param(value []map[string]json.RawMessage) apiquery.Param {
	texts := make([]string, len(value))
	for i, item := range value {
		texts[i] = encodeJSON(item)
	}
	return apiquery.List(texts)
}
func (objectListCodec) display(value []map[string]json.RawMessage) string {
	if len(value) == 0 {
		return ""
	}
	return encodeJSON(value)
}
func (objectListCodec) setField(current []map[string]json.RawMessage, field string, value json.RawMessage) ([]map[string]json.RawMessage, error) {
	if last := len(current) - 1; last >= 0 && current[last] != nil {
		if _, taken := current[last][field]; !taken {
			current[last][field] = value
			return current, nil
		}
	}
	return append(current, map[string]json.RawMessage{field: value}), nil
}
func (objectListCodec) traits() traits { return traits{typeName: "any", multiValue: true} }

func parseJSON(input string) (json.RawMessage, error) {
	if strings.HasPrefix(input, "@") {
		return json.Marshal(input)
	}
	converted, err := YAMLToJSON([]byte(input))
	if err != nil {
		if allowAsLiteralString(input) {
			return json.Marshal(input)
		}
		return nil, fmt.Errorf("failed to parse as YAML: %w", err)
	}
	return converted, nil
}

func YAMLToJSON(data []byte) (json.RawMessage, error) {
	converted, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	return compactJSON(converted)
}

func parseObject(input string) (map[string]json.RawMessage, error) {
	raw, err := parseJSON(input)
	if err != nil {
		return nil, err
	}
	if string(raw) == "null" {
		return nil, nil
	}
	fields, ok := JSONObject(raw)
	if !ok {
		return nil, fmt.Errorf("cannot convert %q to an object", input)
	}
	return fields, nil
}

func compactJSON(raw []byte) (json.RawMessage, error) {
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		return nil, err
	}
	return compact.Bytes(), nil
}

func jsonText(raw json.RawMessage) string {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	return string(raw)
}

func encodeJSON[T any](value T) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(encoded)
}

func allowAsLiteralString(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) &&
			c != '_' && c != '-' && c != '.' && c != '=' {
			return false
		}
	}
	return true
}
