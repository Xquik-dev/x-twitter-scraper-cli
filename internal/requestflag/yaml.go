// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package requestflag

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
)

// YAMLToJSON follows the decoder that goccy/go-yaml runs for an empty interface target, with each value typed.
// https://github.com/goccy/go-yaml/blob/v1.19.2/decode.go#L345
func YAMLToJSON(data []byte) (json.RawMessage, error) {
	file, err := parser.ParseBytes(data, 0)
	if err != nil {
		return nil, err
	}
	d := &yamlDecoder{nodes: map[string]ast.Node{}, values: map[string]yamlValue{}}
	var bodies []ast.Node
	for _, document := range file.Docs {
		value, err := d.value(nil, document.Body)
		if err != nil {
			return nil, err
		}
		if _, isNull := value.(yamlNull); !isNull || (document.Body != nil && document.Body.Type() == ast.NullType) {
			bodies = append(bodies, document.Body)
		}
	}
	if len(bodies) == 0 {
		return json.RawMessage("null"), nil
	}
	d.values = map[string]yamlValue{}
	d.depth = 1
	var anchors enclosingAnchors
	body := bodies[0]
	if anchor, isAnchor := body.(*ast.AnchorNode); isAnchor {
		d.depth++
		anchors = anchors.with(anchor.Name.GetToken().Value)
		body = anchor.Value
	}
	// goccy's decodeValue reads an interface target's node a second time, and that read sees the anchor values of the first.
	// https://github.com/goccy/go-yaml/blob/v1.19.2/decode.go#L1037
	var value yamlValue
	for range 2 {
		if value, err = d.value(anchors, body); err != nil {
			return nil, err
		}
	}
	return value.appendJSON(nil)
}

type yamlDecoder struct {
	nodes  map[string]ast.Node
	values map[string]yamlValue
	depth  int
}

type enclosingAnchors map[string]bool

func (s enclosingAnchors) with(name string) enclosingAnchors {
	inner := enclosingAnchors{name: true}
	maps.Copy(inner, s)
	return inner
}

const maxYAMLDepth = 10000

func (d *yamlDecoder) enter() error {
	d.depth++
	if d.depth > maxYAMLDepth {
		return yaml.ErrExceededMaxDepth
	}
	return nil
}

func (d *yamlDecoder) leave() { d.depth-- }

func (d *yamlDecoder) value(anchors enclosingAnchors, node ast.Node) (yamlValue, error) {
	defer d.leave()
	if err := d.enter(); err != nil {
		return nil, err
	}
	switch node := node.(type) {
	case *ast.StringNode:
		return yamlText(node.Value), nil
	case *ast.LiteralNode:
		return yamlText(node.Value.Value), nil
	case *ast.IntegerNode:
		return integer(node), nil
	case *ast.FloatNode:
		return yamlFloat(node.Value), nil
	case *ast.InfinityNode:
		return yamlFloat(node.Value), nil
	case *ast.NanNode:
		return yamlFloat(math.NaN()), nil
	case *ast.BoolNode:
		return yamlBool(node.Value), nil
	case *ast.TagNode:
		return d.tagged(anchors, node)
	case *ast.AnchorNode:
		name := node.Name.GetToken().Value
		d.nodes[name] = nil
		value, err := d.value(anchors.with(name), node.Value)
		if err != nil {
			delete(d.nodes, name)
			return nil, err
		}
		d.nodes[name] = node.Value
		d.values[name] = value
		return value, nil
	case *ast.AliasNode:
		name := node.Value.String()
		if anchors[name] {
			return yamlNull{}, nil
		}
		if value, found := d.values[name]; found {
			return value, nil
		}
		if target, found := d.nodes[name]; found {
			return d.value(anchors, target)
		}
		return nil, yamlSyntaxError{fmt.Sprintf("could not find alias %q", name), node.Value.GetToken()}
	case *ast.MappingKeyNode:
		return d.value(anchors, node.Value)
	case *ast.MappingValueNode:
		fields := yamlMap{}
		if err := d.setPair(anchors, node, fields); err != nil {
			return nil, err
		}
		return fields, nil
	case *ast.MappingNode:
		fields := yamlMap{}
		for _, pair := range node.Values {
			if err := d.setField(anchors, pair, fields); err != nil {
				return nil, err
			}
		}
		return fields, nil
	case *ast.SequenceNode:
		items := make(yamlList, len(node.Values))
		for i, item := range node.Values {
			value, err := d.value(anchors, item)
			if err != nil {
				return nil, err
			}
			items[i] = value
		}
		return items, nil
	}
	return yamlNull{}, nil
}

func integer(node *ast.IntegerNode) yamlValue {
	if whole, ok := node.Value.(uint64); ok {
		return yamlUint(whole)
	}
	if whole, ok := node.Value.(int64); ok {
		return yamlInt(whole)
	}
	return yamlNull{}
}

func (d *yamlDecoder) tagged(anchors enclosingAnchors, tag *ast.TagNode) (yamlValue, error) {
	if tag.Directive != nil {
		return d.scalar(anchors, tag.Value, toText)
	}
	switch token.ReservedTagKeyword(tag.Start.Value) {
	case token.NullTag:
		return yamlNull{}, nil
	case token.TimestampTag:
		return d.timestamp(anchors, tag.Value), nil
	case token.StringTag:
		return d.scalar(anchors, tag.Value, toText)
	case token.IntegerTag:
		return d.scalar(anchors, tag.Value, toInt)
	case token.FloatTag:
		return d.scalar(anchors, tag.Value, toFloat)
	case token.BinaryTag:
		return d.scalar(anchors, tag.Value, toBytes)
	case token.BooleanTag:
		return d.scalar(anchors, tag.Value, toBool)
	}
	return d.value(anchors, tag.Value)
}

func (d *yamlDecoder) scalar(anchors enclosingAnchors, node ast.Node, convert func(yamlValue) (yamlValue, error)) (yamlValue, error) {
	value, err := d.value(anchors, node)
	if err != nil {
		return nil, err
	}
	converted, err := convert(value)
	if err != nil {
		return nil, yamlSyntaxError{err.Error(), node.GetToken()}
	}
	return converted, nil
}

func toText(value yamlValue) (yamlValue, error) {
	if _, isNull := value.(yamlNull); isNull {
		return yamlText(""), nil
	}
	return yamlText(value.sprint()), nil
}

func toInt(value yamlValue) (yamlValue, error) {
	whole, _ := strconv.Atoi(value.sprint())
	return yamlInt(whole), nil
}

func toFloat(value yamlValue) (yamlValue, error) {
	switch value := value.(type) {
	case yamlInt:
		return yamlFloat(value), nil
	case yamlUint:
		return yamlFloat(value), nil
	case yamlFloat:
		return value, nil
	case yamlText:
		number, _ := strconv.ParseFloat(string(value), 64)
		return yamlFloat(number), nil
	}
	return yamlFloat(0), nil
}

func toBytes(value yamlValue) (yamlValue, error) {
	text, isText := value.(yamlText)
	if !isText {
		return nil, fmt.Errorf("cannot convert %q to string", value.sprint())
	}
	decoded, _ := base64.StdEncoding.DecodeString(string(text))
	return yamlBytes(decoded), nil
}

func toBool(value yamlValue) (yamlValue, error) {
	lowered := strings.ToLower(value.sprint())
	if truth, err := strconv.ParseBool(lowered); err == nil {
		return yamlBool(truth), nil
	}
	switch lowered {
	case "yes":
		return yamlBool(true), nil
	case "no":
		return yamlBool(false), nil
	}
	return nil, fmt.Errorf("cannot convert %q to boolean", value.sprint())
}

var timestampFormats = []string{
	"2006-1-2T15:4:5.999999999Z07:00",
	"2006-1-2t15:4:5.999999999Z07:00",
	"2006-1-2 15:4:5.999999999",
	"2006-1-2",
}

func (d *yamlDecoder) timestamp(anchors enclosingAnchors, node ast.Node) yamlTime {
	value, err := d.value(anchors, node)
	if err != nil {
		return yamlTime{}
	}
	switch value := value.(type) {
	case yamlTime:
		return value
	case yamlText:
		for _, format := range timestampFormats {
			if moment, err := time.Parse(format, string(value)); err == nil {
				return yamlTime(moment)
			}
		}
	}
	return yamlTime{}
}

func (d *yamlDecoder) setField(anchors enclosingAnchors, pair *ast.MappingValueNode, fields yamlMap) error {
	defer d.leave()
	if err := d.enter(); err != nil {
		return err
	}
	return d.setPair(anchors, pair, fields)
}

func (d *yamlDecoder) setPair(anchors enclosingAnchors, pair *ast.MappingValueNode, fields yamlMap) error {
	if pair.Key.IsMergeKey() {
		source, err := d.mapNode(pair.Value, true)
		if err != nil {
			return err
		}
		for merged := source.MapRange(); merged.Next(); {
			if err := d.setField(anchors, merged.KeyValue(), fields); err != nil {
				return err
			}
		}
		return nil
	}
	key, err := d.value(anchors, pair.Key)
	if err != nil {
		return err
	}
	value, err := d.value(anchors, pair.Value)
	if err != nil {
		return err
	}
	if _, isNull := key.(yamlNull); isNull {
		fields["null"] = value
	} else {
		fields[key.sprint()] = value
	}
	return nil
}

func (d *yamlDecoder) mapNode(node ast.Node, isMerge bool) (ast.MapNode, error) {
	defer d.leave()
	if err := d.enter(); err != nil {
		return nil, err
	}
	switch node := node.(type) {
	case ast.MapNode:
		return node, nil
	case *ast.AnchorNode:
		d.nodes[node.Name.GetToken().Value] = node.Value
		return d.mapNode(node.Value, isMerge)
	case *ast.AliasNode:
		name := node.Value.GetToken().Value
		target := d.nodes[name]
		if target == nil {
			return nil, fmt.Errorf("cannot find anchor by alias name %s", name)
		}
		return d.mapNode(target, isMerge)
	case *ast.SequenceNode:
		if isMerge {
			sources := make([]ast.MapNode, len(node.Values))
			for i, item := range node.Values {
				source, err := d.mapNode(item, false)
				if err != nil {
					return nil, err
				}
				sources[i] = source
			}
			return ast.SequenceMergeValue(sources...), nil
		}
	}
	return nil, yamlSyntaxError{fmt.Sprintf("%s was used where %s is expected", node.Type().YAMLName(), ast.MappingType.YAMLName()), node.GetToken()}
}

type yamlSyntaxError struct {
	message string
	token   *token.Token
}

func (e yamlSyntaxError) Error() string {
	var text printer.Printer
	position := fmt.Sprintf("[%d:%d] ", e.token.Position.Line, e.token.Position.Column)
	return text.PrintErrorMessage(position+e.message, false) + "\n" + text.PrintErrorToken(e.token, false)
}

type yamlValue interface {
	sprint() string
	appendJSON(to []byte) ([]byte, error)
}

type yamlNull struct{}

type yamlText string

type yamlInt int64

type yamlUint uint64

type yamlFloat float64

type yamlBool bool

type yamlTime time.Time

type yamlBytes []byte

type yamlMap map[string]yamlValue

type yamlList []yamlValue

func (yamlNull) sprint() string    { return "<nil>" }
func (v yamlText) sprint() string  { return string(v) }
func (v yamlInt) sprint() string   { return strconv.FormatInt(int64(v), 10) }
func (v yamlUint) sprint() string  { return strconv.FormatUint(uint64(v), 10) }
func (v yamlFloat) sprint() string { return strconv.FormatFloat(float64(v), 'g', -1, 64) }
func (v yamlBool) sprint() string  { return strconv.FormatBool(bool(v)) }
func (v yamlTime) sprint() string  { return time.Time(v).String() }
func (v yamlBytes) sprint() string { return fmt.Sprint([]byte(v)) }

func (v yamlMap) sprint() string {
	keys := slices.Sorted(maps.Keys(v))
	fields := make([]string, len(keys))
	for i, key := range keys {
		fields[i] = key + ":" + v[key].sprint()
	}
	return "map[" + strings.Join(fields, " ") + "]"
}

func (v yamlList) sprint() string {
	items := make([]string, len(v))
	for i, item := range v {
		items[i] = item.sprint()
	}
	return "[" + strings.Join(items, " ") + "]"
}

func (yamlNull) appendJSON(to []byte) ([]byte, error)   { return append(to, "null"...), nil }
func (v yamlText) appendJSON(to []byte) ([]byte, error) { return appendMarshaled(to, string(v)) }
func (v yamlInt) appendJSON(to []byte) ([]byte, error) {
	return strconv.AppendInt(to, int64(v), 10), nil
}
func (v yamlUint) appendJSON(to []byte) ([]byte, error) {
	return strconv.AppendUint(to, uint64(v), 10), nil
}
func (v yamlFloat) appendJSON(to []byte) ([]byte, error) {
	return appendMarshaled(to, float64(v))
}
func (v yamlBool) appendJSON(to []byte) ([]byte, error)  { return strconv.AppendBool(to, bool(v)), nil }
func (v yamlTime) appendJSON(to []byte) ([]byte, error)  { return appendMarshaled(to, time.Time(v)) }
func (v yamlBytes) appendJSON(to []byte) ([]byte, error) { return appendMarshaled(to, []byte(v)) }

func (v yamlMap) appendJSON(to []byte) ([]byte, error) {
	to = append(to, '{')
	keys := slices.Sorted(maps.Keys(v))
	for i, key := range keys {
		if i > 0 {
			to = append(to, ',')
		}
		var err error
		if to, err = appendMarshaled(to, key); err != nil {
			return nil, err
		}
		if to, err = v[key].appendJSON(append(to, ':')); err != nil {
			return nil, err
		}
	}
	return append(to, '}'), nil
}

func (v yamlList) appendJSON(to []byte) ([]byte, error) {
	to = append(to, '[')
	for i, item := range v {
		if i > 0 {
			to = append(to, ',')
		}
		var err error
		if to, err = item.appendJSON(to); err != nil {
			return nil, err
		}
	}
	return append(to, ']'), nil
}

func appendMarshaled[T string | float64 | time.Time | []byte](to []byte, value T) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return append(to, encoded...), nil
}
