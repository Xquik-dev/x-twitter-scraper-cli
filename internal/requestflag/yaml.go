// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

type yamlDocument struct {
	anchors     map[string]json.RawMessage
	anchorNodes map[string]ast.Node
}

func YAMLToJSON(data []byte) (json.RawMessage, error) {
	file, err := parser.ParseBytes(data, 0)
	if err != nil {
		return nil, err
	}
	var body ast.Node
	for _, document := range file.Docs {
		if _, isDirective := document.Body.(*ast.DirectiveNode); !isDirective {
			body = document.Body
			break
		}
	}
	document := yamlDocument{anchors: map[string]json.RawMessage{}, anchorNodes: map[string]ast.Node{}}
	return document.convert(body)
}

func (d yamlDocument) convert(node ast.Node) (json.RawMessage, error) {
	switch node := node.(type) {
	case nil, *ast.NullNode, *ast.CommentGroupNode:
		return json.RawMessage("null"), nil
	case *ast.AnchorNode:
		raw, err := d.convert(node.Value)
		if err == nil {
			d.anchors[node.Name.GetToken().Value] = raw
			d.anchorNodes[node.Name.GetToken().Value] = node.Value
		}
		return raw, err
	case *ast.AliasNode:
		name := node.Value.GetToken().Value
		raw, found := d.anchors[name]
		if !found {
			return nil, fmt.Errorf("alias %q has no anchor", name)
		}
		return raw, nil
	case *ast.TagNode:
		return d.tagged(node)
	case *ast.MappingNode:
		return d.mapping(node.Values)
	case *ast.MappingValueNode:
		return d.mapping([]*ast.MappingValueNode{node})
	case *ast.SequenceNode:
		items := make([]json.RawMessage, len(node.Values))
		for i, value := range node.Values {
			item, err := d.convert(value)
			if err != nil {
				return nil, err
			}
			items[i] = item
		}
		return json.Marshal(items)
	case *ast.IntegerNode:
		if raw, err := encodeScalar[int64](node); err == nil {
			return raw, nil
		}
		if raw, err := encodeScalar[uint64](node); err == nil {
			return raw, nil
		}
		return encodeScalar[float64](node)
	case *ast.FloatNode, *ast.InfinityNode, *ast.NanNode:
		return encodeScalar[float64](node)
	case *ast.BoolNode:
		return encodeScalar[bool](node)
	default:
		return encodeScalar[string](node)
	}
}

var timestampFormats = []string{
	"2006-1-2T15:4:5.999999999Z07:00",
	"2006-1-2t15:4:5.999999999Z07:00",
	"2006-1-2 15:4:5.999999999",
	"2006-1-2",
}

func (d yamlDocument) resolve(node ast.Node) (ast.Node, error) {
	for {
		switch reference := node.(type) {
		case *ast.AnchorNode:
			if _, err := d.convert(reference); err != nil {
				return nil, err
			}
			node = reference.Value
		case *ast.AliasNode:
			name := reference.Value.GetToken().Value
			target, found := d.anchorNodes[name]
			if !found {
				return nil, fmt.Errorf("alias %q has no anchor", name)
			}
			node = target
		case *ast.MappingKeyNode:
			node = reference.Value
		default:
			return node, nil
		}
	}
}

type tagValue interface {
	printed() string
	encoded() (json.RawMessage, error)
}

type textValue string

type nullValue struct{}

type intValue int

type floatValue float64

type boolValue bool

type timeValue time.Time

type bytesValue []byte

func (v textValue) printed() string                   { return string(v) }
func (v textValue) encoded() (json.RawMessage, error) { return json.Marshal(string(v)) }

func (nullValue) printed() string                   { return "" }
func (nullValue) encoded() (json.RawMessage, error) { return json.RawMessage("null"), nil }

func (v intValue) printed() string                   { return strconv.Itoa(int(v)) }
func (v intValue) encoded() (json.RawMessage, error) { return json.Marshal(int(v)) }

func (v floatValue) printed() string                   { return fmt.Sprint(float64(v)) }
func (v floatValue) encoded() (json.RawMessage, error) { return json.Marshal(float64(v)) }

func (v boolValue) printed() string                   { return strconv.FormatBool(bool(v)) }
func (v boolValue) encoded() (json.RawMessage, error) { return json.Marshal(bool(v)) }

func (v timeValue) printed() string                   { return time.Time(v).Format(time.RFC3339Nano) }
func (v timeValue) encoded() (json.RawMessage, error) { return json.Marshal(time.Time(v)) }

func (v bytesValue) printed() string                   { return base64.StdEncoding.EncodeToString(v) }
func (v bytesValue) encoded() (json.RawMessage, error) { return json.Marshal([]byte(v)) }

func (d yamlDocument) tagged(tag *ast.TagNode) (json.RawMessage, error) {
	value, isScalar, err := d.readTag(tag)
	if err != nil {
		return nil, err
	}
	if !isScalar {
		return d.convert(tag.Value)
	}
	return value.encoded()
}

func (d yamlDocument) readTag(tag *ast.TagNode) (tagValue, bool, error) {
	name := tag.Start.Value
	if tag.Directive != nil {
		name = "!!str"
	}
	var read func(ast.Node) (tagValue, error)
	switch name {
	case "!!str":
		read = d.readText
	case "!!null":
		read = readNull
	case "!!int":
		read = d.readInt
	case "!!float":
		read = d.readFloat
	case "!!bool":
		read = d.readBool
	case "!!timestamp":
		read = d.readTimestamp
	case "!!binary":
		read = d.readBinary
	default:
		return nil, false, nil
	}
	node, err := d.resolve(tag.Value)
	if err != nil {
		return nil, true, err
	}
	value, err := read(node)
	return value, true, err
}

func (d yamlDocument) readText(node ast.Node) (tagValue, error) {
	text, err := d.text(node)
	return textValue(text), err
}

func readNull(ast.Node) (tagValue, error) {
	return nullValue{}, nil
}

func (d yamlDocument) readInt(node ast.Node) (tagValue, error) {
	text, _ := d.text(node)
	whole, _ := strconv.Atoi(text)
	return intValue(whole), nil
}

func (d yamlDocument) readFloat(node ast.Node) (tagValue, error) {
	switch node.(type) {
	case *ast.IntegerNode, *ast.FloatNode, *ast.InfinityNode, *ast.NanNode:
		var number float64
		err := yaml.NodeToValue(node, &number)
		return floatValue(number), err
	case *ast.StringNode, *ast.TagNode:
		text, _ := d.text(node)
		number, _ := strconv.ParseFloat(text, 64)
		return floatValue(number), nil
	default:
		return floatValue(0), nil
	}
}

func (d yamlDocument) readBool(node ast.Node) (tagValue, error) {
	text, err := d.text(node)
	if err != nil {
		return nil, err
	}
	switch lowered := strings.ToLower(text); lowered {
	case "yes", "no":
		return boolValue(lowered == "yes"), nil
	default:
		truth, err := strconv.ParseBool(lowered)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to boolean", text)
		}
		return boolValue(truth), nil
	}
}

func (d yamlDocument) readTimestamp(node ast.Node) (tagValue, error) {
	switch node.(type) {
	case *ast.StringNode, *ast.TagNode:
		text, _ := d.text(node)
		for _, format := range timestampFormats {
			if moment, err := time.Parse(format, text); err == nil {
				return timeValue(moment), nil
			}
		}
	}
	return timeValue(time.Time{}), nil
}

func (d yamlDocument) readBinary(node ast.Node) (tagValue, error) {
	if !d.holdsText(node) {
		return nil, fmt.Errorf("cannot convert %s to binary", node.String())
	}
	text, err := d.text(node)
	if err != nil {
		return nil, err
	}
	decoded, _ := base64.StdEncoding.DecodeString(text)
	return bytesValue(decoded), nil
}

func (d yamlDocument) text(node ast.Node) (string, error) {
	tag, isTagged := node.(*ast.TagNode)
	if !isTagged {
		return scalarText(node)
	}
	value, isScalar, err := d.readTag(tag)
	if err != nil {
		return "", err
	}
	if isScalar {
		return value.printed(), nil
	}
	inner, err := d.resolve(tag.Value)
	if err != nil {
		return "", err
	}
	return d.text(inner)
}

func (d yamlDocument) holdsText(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.StringNode, *ast.LiteralNode:
		return true
	case *ast.TagNode:
		value, isScalar, err := d.readTag(node)
		if isScalar {
			_, isText := value.(textValue)
			return err == nil && isText
		}
		inner, err := d.resolve(node.Value)
		return err == nil && d.holdsText(inner)
	default:
		return false
	}
}

func scalarText(node ast.Node) (string, error) {
	switch node.(type) {
	case nil, *ast.NullNode:
		return "", nil
	case *ast.IntegerNode:
		var whole int64
		if err := yaml.NodeToValue(node, &whole); err == nil {
			return strconv.FormatInt(whole, 10), nil
		}
		var natural uint64
		if err := yaml.NodeToValue(node, &natural); err == nil {
			return strconv.FormatUint(natural, 10), nil
		}
		var number float64
		err := yaml.NodeToValue(node, &number)
		return fmt.Sprint(number), err
	case *ast.FloatNode, *ast.InfinityNode, *ast.NanNode:
		var number float64
		err := yaml.NodeToValue(node, &number)
		return fmt.Sprint(number), err
	case *ast.BoolNode:
		var truth bool
		err := yaml.NodeToValue(node, &truth)
		return strconv.FormatBool(truth), err
	default:
		var text string
		err := yaml.NodeToValue(node, &text)
		return text, err
	}
}

func (d yamlDocument) mapping(pairs []*ast.MappingValueNode) (json.RawMessage, error) {
	fields := map[string]json.RawMessage{}
	var merged []map[string]json.RawMessage
	for _, pair := range pairs {
		if pair.Key.IsMergeKey() {
			sources, err := d.mergeSources(pair.Value)
			if err != nil {
				return nil, err
			}
			merged = append(merged, sources...)
			continue
		}
		key, err := d.keyText(pair.Key)
		if err != nil {
			return nil, err
		}
		value, err := d.convert(pair.Value)
		if err != nil {
			return nil, err
		}
		fields[key] = value
	}
	for _, source := range merged {
		for key, value := range source {
			if _, taken := fields[key]; !taken {
				fields[key] = value
			}
		}
	}
	return json.Marshal(fields)
}

func (d yamlDocument) keyText(key ast.MapKeyNode) (string, error) {
	node, err := d.resolve(key)
	if err != nil {
		return "", err
	}
	if _, isTagged := node.(*ast.TagNode); isTagged {
		return d.text(node)
	}
	var text string
	err = yaml.NodeToValue(node, &text)
	return text, err
}

func (d yamlDocument) mergeSources(node ast.Node) ([]map[string]json.RawMessage, error) {
	raw, err := d.convert(node)
	if err != nil {
		return nil, err
	}
	if fields, ok := JSONObject(raw); ok {
		return []map[string]json.RawMessage{fields}, nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil || items == nil {
		return nil, fmt.Errorf("a merge key needs a mapping or a list of mappings")
	}
	sources := make([]map[string]json.RawMessage, 0, len(items))
	for _, item := range items {
		fields, ok := JSONObject(item)
		if !ok {
			return nil, fmt.Errorf("a merge key needs a mapping or a list of mappings")
		}
		sources = append(sources, fields)
	}
	return sources, nil
}

func encodeScalar[T string | int64 | uint64 | float64 | bool](node ast.Node) (json.RawMessage, error) {
	var value T
	if err := yaml.NodeToValue(node, &value); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}
