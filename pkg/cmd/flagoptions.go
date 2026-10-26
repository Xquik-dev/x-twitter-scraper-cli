// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/debugmiddleware"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"

	"github.com/urfave/cli/v3"
)

type BodyContentType int

const (
	EmptyBody BodyContentType = iota
	ApplicationJSON
)

type onceStdinReader struct {
	stdinReader   io.Reader
	failureReason string
}

func (o *onceStdinReader) readAll() ([]byte, error) {
	if o.failureReason != "" {
		return nil, fmt.Errorf("cannot read from stdin: %s", o.failureReason)
	}
	if o.stdinReader == nil {
		return nil, fmt.Errorf("stdin has already been read by another parameter; it can only be read once")
	}
	r := o.stdinReader
	o.stdinReader = nil
	return io.ReadAll(r)
}

func isStdinPath(s string) bool {
	switch s {
	case "-", "/dev/fd/0", "/dev/stdin":
		return true
	}
	return false
}

func embedJSON(raw json.RawMessage, stdin *onceStdinReader) (json.RawMessage, error) {
	switch jsonKind(raw) {
	case '"':
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			return nil, err
		}
		embedded, err := embedText(text, stdin)
		if err != nil || embedded == text {
			return raw, err
		}
		return json.Marshal(embedded)
	case '[':
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, err
		}
		for i, item := range items {
			embedded, err := embedJSON(item, stdin)
			if err != nil {
				return nil, err
			}
			items[i] = embedded
		}
		return json.Marshal(items)
	case '{':
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fields); err != nil {
			return nil, err
		}
		embedded, err := embedFields(fields, stdin)
		if err != nil {
			return nil, err
		}
		return json.Marshal(embedded)
	default:
		return raw, nil
	}
}

func jsonKind(raw json.RawMessage) byte {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return 0
	}
	return trimmed[0]
}

func embedFields(fields map[string]json.RawMessage, stdin *onceStdinReader) (map[string]json.RawMessage, error) {
	embedded := make(map[string]json.RawMessage, len(fields))
	for key, value := range fields {
		replaced, err := embedJSON(value, stdin)
		if err != nil {
			return nil, err
		}
		embedded[key] = replaced
	}
	return embedded, nil
}

func embedParams(params map[string]apiquery.Param, stdin *onceStdinReader) (map[string]apiquery.Param, error) {
	embedded := make(map[string]apiquery.Param, len(params))
	for key, param := range params {
		replaced, err := param.Map(func(value string) (string, error) { return embedText(value, stdin) })
		if err != nil {
			return nil, err
		}
		embedded[key] = replaced
	}
	return embedded, nil
}

func embedText(s string, stdin *onceStdinReader) (string, error) {
	if literal, ok := strings.CutPrefix(s, "\\@"); ok {
		return "@" + literal, nil
	}
	if filename, ok := strings.CutPrefix(s, "@data://"); ok {
		content, err := readInput(filename, stdin)
		if err != nil {
			return s, err
		}
		return base64.StdEncoding.EncodeToString(content), nil
	}
	if filename, ok := strings.CutPrefix(s, "@file://"); ok {
		content, err := readInput(filename, stdin)
		if err != nil {
			return s, err
		}
		return string(content), nil
	}
	filename, ok := strings.CutPrefix(s, "@")
	if !ok {
		return s, nil
	}
	content, err := readInput(filename, stdin)
	if err != nil {
		probablyFile := isStdinPath(filename) || strings.Contains(filename, ".") || strings.Contains(filename, "/")
		if probablyFile {
			return s, err
		}
		return s, nil
	}
	if isUTF8TextFile(content) {
		return string(content), nil
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

func readInput(filename string, stdin *onceStdinReader) ([]byte, error) {
	if isStdinPath(filename) {
		return stdin.readAll()
	}
	return os.ReadFile(filename)
}

func isUTF8TextFile(content []byte) bool {
	textTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-javascript",
		"application/ecmascript",
		"application/x-ecmascript",
	}

	contentType := http.DetectContentType(content)
	for _, prefix := range textTypes {
		if strings.HasPrefix(contentType, prefix) {
			return utf8.Valid(content)
		}
	}
	return false
}

func flagOptions(cmd *cli.Command, bodyType BodyContentType) ([]option.RequestOption, error) {
	var options []option.RequestOption
	if cmd.Bool("debug") {
		options = append(options, option.WithMiddleware(debugmiddleware.NewRequestLogger().Middleware()))
	}

	requestContents, err := requestflag.ExtractRequestContents(cmd)
	if err != nil {
		return nil, err
	}
	if err := applyDataAliases(cmd, requestContents.Body); err != nil {
		return nil, err
	}

	var rootBody json.RawMessage
	stdinReader := onceStdinReader{stdinReader: os.Stdin}
	if isInputPiped() {
		pipeData, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		if len(pipeData) > 0 {
			stdinReader = onceStdinReader{failureReason: "stdin was already consumed by piped YAML/JSON input"}
			requestContents, rootBody, err = applyPipedData(cmd, pipeData, bodyType, requestContents)
			if err != nil {
				return nil, err
			}
		}
	}

	if missingFlags := requestflag.GetMissingRequiredFlags(cmd, requestContents.Body); len(missingFlags) > 0 {
		if len(missingFlags) == 1 {
			return nil, fmt.Errorf("required flag %q is missing. Run '%s --help' for usage", missingFlags[0].Names()[0], cmd.FullName())
		}
		names := []string{}
		for _, flag := range missingFlags {
			names = append(names, flag.Names()[0])
		}
		return nil, fmt.Errorf("required flags %q are missing. Run '%s --help' for usage", strings.Join(names, ", "), cmd.FullName())
	}

	var bodyBytes []byte
	if rootBody != nil {
		bodyBytes, err = embedJSON(rootBody, &stdinReader)
	} else {
		var fields map[string]json.RawMessage
		fields, err = embedFields(requestContents.Body, &stdinReader)
		if err == nil {
			bodyBytes, err = json.Marshal(fields)
		}
	}
	if err != nil {
		return nil, err
	}
	headers, err := embedParams(requestContents.Headers, &stdinReader)
	if err != nil {
		return nil, err
	}
	queries, err := embedParams(requestContents.Queries, &stdinReader)
	if err != nil {
		return nil, err
	}

	for key, values := range apiquery.Encode(queries, apiquery.ArrayQueryFormatComma) {
		options = append(options, option.WithQuery(key, values[0]))
		for _, value := range values[1:] {
			options = append(options, option.WithQueryAdd(key, value))
		}
	}
	for key, values := range apiquery.Encode(headers, apiquery.ArrayQueryFormatRepeat) {
		options = append(options, option.WithHeader(key, values[0]))
		for _, value := range values[1:] {
			options = append(options, option.WithHeaderAdd(key, value))
		}
	}

	if bodyType == ApplicationJSON {
		options = append(options, option.WithRequestBody("application/json", bodyBytes))
	}
	return options, nil
}

func applyPipedData(
	cmd *cli.Command,
	pipeData []byte,
	bodyType BodyContentType,
	contents requestflag.RequestContents,
) (requestflag.RequestContents, json.RawMessage, error) {
	bodyData, err := requestflag.YAMLToJSON(pipeData)
	if err != nil {
		return contents, nil, fmt.Errorf("piped data is invalid YAML or JSON. Fix the input:\n%w", err)
	}
	bodyMap, isObject := requestflag.JSONObject(bodyData)
	if !isObject {
		if bodyType == EmptyBody {
			return contents, nil, nil
		}
		if len(contents.Body) > 0 {
			return contents, nil, fmt.Errorf("request body must be a map when using flags. Fix the input: %s", bodyData)
		}
		return contents, bodyData, nil
	}

	if err := applyDataAliases(cmd, bodyMap); err != nil {
		return contents, nil, err
	}
	if err := requestflag.ApplyStdinDataToFlags(cmd, bodyMap); err != nil {
		return contents, nil, err
	}
	contents, err = requestflag.ExtractRequestContents(cmd)
	if err != nil {
		return contents, nil, err
	}
	for _, flag := range cmd.Flags {
		inRequest, ok := flag.(requestflag.InRequest)
		if !ok || !flag.IsSet() {
			continue
		}
		if inRequest.GetQueryPath() != "" || inRequest.GetHeaderPath() != "" || inRequest.GetPathParam() != "" {
			delete(bodyMap, inRequest.GetQueryPath())
			delete(bodyMap, inRequest.GetHeaderPath())
			delete(bodyMap, inRequest.GetPathParam())
			for _, alias := range inRequest.GetDataAliases() {
				delete(bodyMap, alias)
			}
		}
	}
	if bodyType != EmptyBody {
		maps.Copy(bodyMap, contents.Body)
		contents.Body = bodyMap
	}
	return contents, nil, nil
}

func applyDataAliases(cmd *cli.Command, bodyMap map[string]json.RawMessage) error {
	for _, flag := range cmd.Flags {
		if inner, ok := flag.(requestflag.HasOuterFlag); ok {
			outer, ok := inner.GetOuterFlag().(requestflag.InRequest)
			if !ok || inner.GetInnerField() == "" {
				continue
			}
			nested, ok := requestflag.JSONObject(bodyMap[outer.GetBodyPath()])
			if !ok || !rewriteAliases(nested, inner.GetInnerField(), inner.GetDataAliases()) {
				continue
			}
			encoded, err := json.Marshal(nested)
			if err != nil {
				return err
			}
			bodyMap[outer.GetBodyPath()] = encoded
			continue
		}
		if inRequest, ok := flag.(requestflag.InRequest); ok && inRequest.GetBodyPath() != "" {
			rewriteAliases(bodyMap, inRequest.GetBodyPath(), inRequest.GetDataAliases())
		}
	}
	return nil
}

func rewriteAliases(m map[string]json.RawMessage, canonical string, aliases []string) bool {
	rewritten := false
	for _, alias := range aliases {
		if alias == "" || alias == canonical {
			continue
		}
		if val, exists := m[alias]; exists {
			m[canonical] = val
			delete(m, alias)
			rewritten = true
		}
	}
	return rewritten
}
