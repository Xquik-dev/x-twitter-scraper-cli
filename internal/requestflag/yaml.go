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

func (d yamlDocument) tagged(tag *ast.TagNode) (json.RawMessage, error) {
	node, err := d.resolve(tag.Value)
	if err != nil {
		return nil, err
	}
	switch tag.Start.Value {
	case "!!str":
		return d.taggedText(node)
	case "!!null":
		return json.RawMessage("null"), nil
	case "!!int":
		text, _ := scalarText(node)
		whole, _ := strconv.Atoi(text)
		return json.Marshal(whole)
	case "!!float":
		return taggedFloat(node)
	case "!!bool":
		text, err := scalarText(node)
		if err != nil {
			return nil, err
		}
		switch lowered := strings.ToLower(text); lowered {
		case "yes", "no":
			return json.Marshal(lowered == "yes")
		default:
			truth, err := strconv.ParseBool(lowered)
			if err != nil {
				return nil, fmt.Errorf("cannot convert %q to boolean", text)
			}
			return json.Marshal(truth)
		}
	case "!!timestamp":
		var moment time.Time
		if _, isText := node.(*ast.StringNode); isText {
			text, err := scalarText(node)
			if err != nil {
				return nil, err
			}
			for _, format := range timestampFormats {
				if parsed, err := time.Parse(format, text); err == nil {
					moment = parsed
					break
				}
			}
		}
		return json.Marshal(moment)
	case "!!binary":
		switch node.(type) {
		case *ast.StringNode, *ast.LiteralNode:
		default:
			return nil, fmt.Errorf("cannot convert %s to binary", node.String())
		}
		text, err := scalarText(node)
		if err != nil {
			return nil, err
		}
		decoded, _ := base64.StdEncoding.DecodeString(text)
		return json.Marshal(decoded)
	default:
		return d.convert(node)
	}
}

func (d yamlDocument) taggedText(node ast.Node) (json.RawMessage, error) {
	text, err := d.text(node)
	if err != nil {
		return nil, err
	}
	return json.Marshal(text)
}

func (d yamlDocument) text(node ast.Node) (string, error) {
	tag, isTagged := node.(*ast.TagNode)
	if !isTagged {
		return scalarText(node)
	}
	raw, err := d.tagged(tag)
	if err != nil {
		return "", err
	}
	switch raw[0] {
	case 'n':
		return "", nil
	case 't', 'f':
		return string(raw), nil
	case '"':
		var text string
		err := json.Unmarshal(raw, &text)
		return text, err
	case '{', '[':
		return "", fmt.Errorf("cannot convert %s to string", raw)
	}
	var whole int64
	if err := json.Unmarshal(raw, &whole); err == nil {
		return strconv.FormatInt(whole, 10), nil
	}
	var number float64
	err = json.Unmarshal(raw, &number)
	return fmt.Sprint(number), err
}

func taggedFloat(node ast.Node) (json.RawMessage, error) {
	switch node.(type) {
	case *ast.IntegerNode, *ast.FloatNode, *ast.InfinityNode, *ast.NanNode:
		return encodeScalar[float64](node)
	case *ast.StringNode:
		text, err := scalarText(node)
		if err != nil {
			return nil, err
		}
		number, _ := strconv.ParseFloat(text, 64)
		return json.Marshal(number)
	default:
		return json.Marshal(0)
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
