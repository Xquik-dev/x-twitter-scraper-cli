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
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiform"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/apiquery"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/debugmiddleware"
	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/requestflag"
	"github.com/Xquik-dev/x-twitter-scraper-go/option"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"
)

type BodyContentType int

const (
	EmptyBody BodyContentType = iota
	MultipartFormEncoded
	ApplicationJSON
	ApplicationOctetStream
)

type FileEmbedStyle int

const (
	// EmbedText loads files inline. It base64-encodes binary data.
	EmbedText FileEmbedStyle = iota

	// EmbedIOReader streams binary request bodies through io.Reader.
	EmbedIOReader
)

// onceStdinReader gives 1 request parameter exclusive access to stdin.
// A failure reason blocks access before the first read.
type onceStdinReader struct {
	stdinReader   io.Reader
	failureReason string
}

func (o *onceStdinReader) read() (io.Reader, error) {
	if o.failureReason != "" {
		return nil, fmt.Errorf("cannot read from stdin: %s", o.failureReason)
	}
	if o.stdinReader == nil {
		return nil, fmt.Errorf("stdin has already been read by another parameter; it can only be read once")
	}
	r := o.stdinReader
	o.stdinReader = nil
	return r, nil
}

func (o *onceStdinReader) readAll() ([]byte, error) {
	r, err := o.read()
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

func isStdinPath(s string) bool {
	switch s {
	case "-", "/dev/fd/0", "/dev/stdin":
		return true
	}
	return false
}

func embedFiles(obj any, embedStyle FileEmbedStyle, stdin *onceStdinReader) (any, error) {
	if obj == nil {
		return obj, nil
	}
	v := reflect.ValueOf(obj)
	result, err := embedFilesValue(v, embedStyle, stdin)
	if err != nil {
		return nil, err
	}
	return result.Interface(), nil
}

// embedFilesValue replaces file references with their contents.
func embedFilesValue(v reflect.Value, embedStyle FileEmbedStyle, stdin *onceStdinReader) (reflect.Value, error) {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		if v.Len() == 0 {
			return v, nil
		}
		result := reflect.MakeMap(reflect.TypeOf(map[string]any{}))

		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()
			newVal, err := embedFilesValue(val, embedStyle, stdin)
			if err != nil {
				return reflect.Value{}, err
			}
			result.SetMapIndex(key, newVal)
		}
		return result, nil

	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return v, nil
		}
		result := reflect.MakeSlice(reflect.TypeOf([]any{}), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			newVal, err := embedFilesValue(v.Index(i), embedStyle, stdin)
			if err != nil {
				return reflect.Value{}, err
			}
			result.Index(i).Set(newVal)
		}
		return result, nil

	case reflect.String:
		// FilePathValue never needs an @ prefix.
		if v.Type() == reflect.TypeOf(FilePathValue("")) {
			s := v.String()
			if s == "" {
				return v, nil
			}
			if embedStyle == EmbedIOReader {
				if isStdinPath(s) {
					r, err := stdin.read()
					if err != nil {
						return v, err
					}
					return reflect.ValueOf(io.NopCloser(r)), nil
				}
				upload, err := openFileUpload(s)
				if err != nil {
					return v, err
				}
				return reflect.ValueOf(upload), nil
			}
			if isStdinPath(s) {
				content, err := stdin.readAll()
				if err != nil {
					return v, err
				}
				return reflect.ValueOf(string(content)), nil
			}
			content, err := os.ReadFile(s)
			if err != nil {
				return v, err
			}
			return reflect.ValueOf(string(content)), nil
		}

		s := v.String()
		if literal, ok := strings.CutPrefix(s, "\\@"); ok {
			// Preserve escaped @ literals.
			return reflect.ValueOf("@" + literal), nil
		}

		if embedStyle == EmbedText {
			if filename, ok := strings.CutPrefix(s, "@data://"); ok {
				// @data:// always uses base64.
				if isStdinPath(filename) {
					content, err := stdin.readAll()
					if err != nil {
						return v, err
					}
					return reflect.ValueOf(base64.StdEncoding.EncodeToString(content)), nil
				}
				content, err := os.ReadFile(filename)
				if err != nil {
					return v, err
				}
				return reflect.ValueOf(base64.StdEncoding.EncodeToString(content)), nil
			} else if filename, ok := strings.CutPrefix(s, "@file://"); ok {
				// @file:// always uses text.
				if isStdinPath(filename) {
					content, err := stdin.readAll()
					if err != nil {
						return v, err
					}
					return reflect.ValueOf(string(content)), nil
				}
				content, err := os.ReadFile(filename)
				if err != nil {
					return v, err
				}
				return reflect.ValueOf(string(content)), nil
			} else if filename, ok := strings.CutPrefix(s, "@"); ok {
				if isStdinPath(filename) {
					content, err := stdin.readAll()
					if err != nil {
						return v, err
					}
					if isUTF8TextFile(content) {
						return reflect.ValueOf(string(content)), nil
					}
					return reflect.ValueOf(base64.StdEncoding.EncodeToString(content)), nil
				}
				content, err := os.ReadFile(filename)
				if err != nil {
					// Treat @username as text, but file-like values as paths.
					probablyFile := strings.Contains(filename, ".") || strings.Contains(filename, "/")
					if probablyFile {
						return v, err
					}
					return v, nil
				}
				if isUTF8TextFile(content) {
					return reflect.ValueOf(string(content)), nil
				}
				return reflect.ValueOf(base64.StdEncoding.EncodeToString(content)), nil
			}
		} else {
			if filename, ok := strings.CutPrefix(s, "@"); ok {
				// Prefixes share upload behavior; @username remains text.
				expectsFile := true
				if withoutPrefix, ok := strings.CutPrefix(filename, "data://"); ok {
					filename = withoutPrefix
				} else if withoutPrefix, ok := strings.CutPrefix(filename, "file://"); ok {
					filename = withoutPrefix
				} else {
					expectsFile = strings.Contains(filename, ".") || strings.Contains(filename, "/")
				}

				if isStdinPath(filename) {
					r, err := stdin.read()
					if err != nil {
						return v, err
					}
					return reflect.ValueOf(io.NopCloser(r)), nil
				}

				upload, err := openFileUpload(filename)
				if err != nil {
					if !expectsFile {
						return v, nil
					}
					return v, err
				}
				return reflect.ValueOf(upload), nil
			}
		}
		return v, nil

	default:
		return v, nil
	}
}

// isUTF8TextFile reports whether MIME sniffing finds valid UTF-8 text.
func isUTF8TextFile(content []byte) bool {
	// DetectContentType follows https://mimesniff.spec.whatwg.org/.
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

func flagOptions(
	cmd *cli.Command,
	nestedFormat apiquery.NestedQueryFormat,
	arrayFormat apiquery.ArrayQueryFormat,
	bodyType BodyContentType,

	// ignoreStdin reserves stdin for a binary parameter named "-".
	ignoreStdin bool,
) ([]option.RequestOption, error) {
	var options []option.RequestOption
	if cmd.Bool("debug") {
		options = append(options, option.WithMiddleware(debugmiddleware.NewRequestLogger().Middleware()))
	}

	requestContents := requestflag.ExtractRequestContents(cmd)

	// Map YAML aliases to canonical API fields.
	if bodyMap, ok := requestContents.Body.(map[string]any); ok {
		applyDataAliases(cmd, bodyMap)
	}

	stdinConsumedByPipe := false
	if bodyType != ApplicationOctetStream && !ignoreStdin && isInputPiped() {
		pipeData, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}

		if len(pipeData) > 0 {
			stdinConsumedByPipe = true
			var bodyData any
			if err := yaml.Unmarshal(pipeData, &bodyData); err != nil {
				return nil, fmt.Errorf("piped data is invalid YAML or JSON. Fix the input:\n%w", err)
			}
			if bodyMap, ok := bodyData.(map[string]any); ok {
				applyDataAliases(cmd, bodyMap)
				if err := requestflag.ApplyStdinDataToFlags(cmd, bodyMap); err != nil {
					return nil, err
				}
				requestContents = requestflag.ExtractRequestContents(cmd)
				// Remove path, query, and header keys before merging the body.
				for _, flag := range cmd.Flags {
					inReq, ok := flag.(requestflag.InRequest)
					if !ok || !flag.IsSet() {
						continue
					}
					if inReq.GetQueryPath() != "" || inReq.GetHeaderPath() != "" || inReq.GetPathParam() != "" {
						delete(bodyMap, inReq.GetQueryPath())
						delete(bodyMap, inReq.GetHeaderPath())
						delete(bodyMap, inReq.GetPathParam())
						for _, alias := range inReq.GetDataAliases() {
							delete(bodyMap, alias)
						}
					}
				}
				if bodyType != EmptyBody {
					if flagMap, ok := requestContents.Body.(map[string]any); ok {
						maps.Copy(bodyMap, flagMap)
						requestContents.Body = bodyMap
					} else {
						bodyData = requestContents.Body
					}
				}
			} else if bodyType != EmptyBody {
				if flagMap, ok := requestContents.Body.(map[string]any); ok && len(flagMap) > 0 {
					return nil, fmt.Errorf("request body must be a map when using flags. Fix the input: %v", bodyData)
				} else {
					requestContents.Body = bodyData
				}
			}
		}
	}

	if missingFlags := requestflag.GetMissingRequiredFlags(cmd, requestContents.Body); len(missingFlags) > 0 {
		if len(missingFlags) == 1 {
			return nil, fmt.Errorf("required flag %q is missing. Run '%s --help' for usage", missingFlags[0].Names()[0], cmd.FullName())
		} else {
			names := []string{}
			for _, flag := range missingFlags {
				names = append(names, flag.Names()[0])
			}
			return nil, fmt.Errorf("required flags %q are missing. Run '%s --help' for usage", strings.Join(names, ", "), cmd.FullName())
		}
	}

	// Wrap binary paths so file expansion never needs an @ prefix.
	wrapFileInputValues(cmd, &requestContents)

	// Reserve stdin for one input source.
	var stdinReader onceStdinReader
	if ignoreStdin {
		stdinReader = onceStdinReader{failureReason: "stdin is already being used for the request body"}
	} else if stdinConsumedByPipe {
		stdinReader = onceStdinReader{failureReason: "stdin was already consumed by piped YAML/JSON input"}
	} else {
		stdinReader = onceStdinReader{stdinReader: os.Stdin}
	}

	// Expand file references across the request.
	embedStyle := EmbedText
	if bodyType == ApplicationOctetStream || bodyType == MultipartFormEncoded {
		embedStyle = EmbedIOReader
	}

	if embedded, err := embedFiles(requestContents.Body, embedStyle, &stdinReader); err != nil {
		return nil, err
	} else {
		requestContents.Body = embedded
	}

	if headersWithFiles, err := embedFiles(requestContents.Headers, EmbedText, &stdinReader); err != nil {
		return nil, err
	} else {
		requestContents.Headers = headersWithFiles.(map[string]any)
	}
	if queriesWithFiles, err := embedFiles(requestContents.Queries, EmbedText, &stdinReader); err != nil {
		return nil, err
	} else {
		requestContents.Queries = queriesWithFiles.(map[string]any)
	}

	querySettings := apiquery.QuerySettings{
		NestedFormat: nestedFormat,
		ArrayFormat:  arrayFormat,
	}

	if values, err := apiquery.MarshalWithSettings(requestContents.Queries, querySettings); err != nil {
		return nil, err
	} else {
		for k, vs := range values {
			if len(vs) == 0 {
				options = append(options, option.WithQueryDel(k))
			} else {
				options = append(options, option.WithQuery(k, vs[0]))
				for _, v := range vs[1:] {
					options = append(options, option.WithQueryAdd(k, v))
				}
			}
		}
	}

	headerSettings := apiquery.QuerySettings{
		NestedFormat: apiquery.NestedQueryFormatDots,
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
	}
	if values, err := apiquery.MarshalWithSettings(requestContents.Headers, headerSettings); err != nil {
		return nil, err
	} else {
		for k, vs := range values {
			if len(vs) == 0 {
				options = append(options, option.WithHeaderDel(k))
			} else {
				options = append(options, option.WithHeader(k, vs[0]))
				for _, v := range vs[1:] {
					options = append(options, option.WithHeaderAdd(k, v))
				}
			}
		}
	}

	switch bodyType {
	case EmptyBody:
		break
	case MultipartFormEncoded:
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		bodyMap, ok := requestContents.Body.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("form request body must be a map. Fix the input: %v", requestContents.Body)
		}
		encodingFormat := apiform.FormatComma
		if err := apiform.MarshalWithSettings(bodyMap, writer, encodingFormat); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		options = append(options, option.WithRequestBody(writer.FormDataContentType(), buf))

	case ApplicationJSON:
		bodyBytes, err := json.Marshal(requestContents.Body)
		if err != nil {
			return nil, err
		}
		options = append(options, option.WithRequestBody("application/json", bodyBytes))

	case ApplicationOctetStream:
		// Body-root flags already set the request body.
		for _, flag := range cmd.Flags {
			if toSend, ok := flag.(requestflag.InRequest); ok && toSend.IsBodyRoot() {
				return options, nil
			}
		}
		if bodyBytes, ok := requestContents.Body.([]byte); ok {
			options = append(options, option.WithRequestBody("application/octet-stream", bodyBytes))
		} else if bodyStr, ok := requestContents.Body.(string); ok {
			options = append(options, option.WithRequestBody("application/octet-stream", []byte(bodyStr)))
		} else {
			return nil, fmt.Errorf("octet-stream body is unsupported. Pass bytes or text: %v", requestContents.Body)
		}

	default:
		panic("invalid body content type")
	}

	return options, nil
}

// FilePathValue marks a path for expansion without an @ prefix.
type FilePathValue string

// fileUpload carries multipart filename and content type metadata.
type fileUpload struct {
	io.Reader
	filename    string
	contentType string
}

func (f fileUpload) Filename() string    { return f.filename }
func (f fileUpload) ContentType() string { return f.contentType }
func (f fileUpload) Close() error {
	if c, ok := f.Reader.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

// openFileUpload opens a path and derives its multipart metadata.
func openFileUpload(path string) (fileUpload, error) {
	file, err := os.Open(path)
	if err != nil {
		return fileUpload{}, err
	}
	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return fileUpload{
		Reader:      file,
		filename:    filepath.Base(path),
		contentType: contentType,
	}, nil
}

// applyDataAliases maps top-level and nested aliases to canonical API fields.
func applyDataAliases(cmd *cli.Command, bodyMap map[string]any) {
	for _, flag := range cmd.Flags {
		if inner, ok := flag.(requestflag.HasOuterFlag); ok {
			outer, outerOk := inner.GetOuterFlag().(requestflag.InRequest)
			if !outerOk {
				continue
			}
			if nested, ok := bodyMap[outer.GetBodyPath()].(map[string]any); ok && inner.GetInnerField() != "" {
				rewriteAliases(nested, inner.GetInnerField(), inner.GetDataAliases())
			}
			continue
		}
		if inReq, ok := flag.(requestflag.InRequest); ok && inReq.GetBodyPath() != "" {
			rewriteAliases(bodyMap, inReq.GetBodyPath(), inReq.GetDataAliases())
		}
	}
}

// rewriteAliases replaces alias keys with one canonical key.
func rewriteAliases(m map[string]any, canonical string, aliases []string) {
	for _, alias := range aliases {
		if alias == "" || alias == canonical {
			continue
		}
		if val, exists := m[alias]; exists {
			m[canonical] = val
			delete(m, alias)
		}
	}
}

// wrapFileInputValues marks binary paths from flags and piped data.
func wrapFileInputValues(cmd *cli.Command, contents *requestflag.RequestContents) {
	bodyMap, _ := contents.Body.(map[string]any)

	for _, flag := range cmd.Flags {
		inReq, ok := flag.(requestflag.InRequest)
		if !ok || !inReq.IsFileInput() || inReq.IsBodyRoot() {
			continue
		}

		if flag.IsSet() {
			if wrapped, changed := wrapFileInputValue(flag.Get()); changed {
				if bodyPath := inReq.GetBodyPath(); bodyPath != "" {
					if bodyMap != nil {
						bodyMap[bodyPath] = wrapped
					}
				} else if queryPath := inReq.GetQueryPath(); queryPath != "" {
					contents.Queries[queryPath] = wrapped
				} else if headerPath := inReq.GetHeaderPath(); headerPath != "" {
					contents.Headers[headerPath] = wrapped
				}
			}
		}

		if bodyPath := inReq.GetBodyPath(); bodyPath != "" && bodyMap != nil {
			if value, exists := bodyMap[bodyPath]; exists {
				if wrapped, changed := wrapFileInputValue(value); changed {
					bodyMap[bodyPath] = wrapped
				}
			}
		}
	}
}

func wrapFileInputValue(value any) (any, bool) {
	switch v := value.(type) {
	case string:
		if v == "" {
			return value, false
		}
		return FilePathValue(v), true

	case []string:
		result := make([]any, len(v))
		for i, s := range v {
			result[i] = FilePathValue(s)
		}
		return result, true

	case []any:
		result := make([]any, len(v))
		for i, elem := range v {
			if s, ok := elem.(string); ok {
				result[i] = FilePathValue(s)
			} else {
				result[i] = elem
			}
		}
		return result, true

	default:
		return value, false
	}
}
