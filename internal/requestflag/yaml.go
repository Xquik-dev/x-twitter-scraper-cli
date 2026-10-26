// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

type yamlDocument struct {
	anchors map[string]json.RawMessage
}

func YAMLToJSON(data []byte) (json.RawMessage, error) {
	file, err := parser.ParseBytes(data, 0)
	if err != nil {
		return nil, err
	}
	var body ast.Node
	if len(file.Docs) > 0 {
		body = file.Docs[0].Body
	}
	document := yamlDocument{anchors: map[string]json.RawMessage{}}
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
		if node.Start.Value == "!!str" {
			return encodeScalar[string](node.Value)
		}
		return d.convert(node.Value)
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

func (d yamlDocument) mapping(pairs []*ast.MappingValueNode) (json.RawMessage, error) {
	fields := map[string]json.RawMessage{}
	var merged []map[string]json.RawMessage
	for _, pair := range pairs {
		value, err := d.convert(pair.Value)
		if err != nil {
			return nil, err
		}
		if pair.Key.IsMergeKey() {
			sources, err := mergeSources(value)
			if err != nil {
				return nil, err
			}
			merged = append(merged, sources...)
			continue
		}
		var key string
		if err := yaml.NodeToValue(pair.Key, &key); err != nil {
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

func mergeSources(raw json.RawMessage) ([]map[string]json.RawMessage, error) {
	if fields, ok := JSONObject(raw); ok {
		return []map[string]json.RawMessage{fields}, nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
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
