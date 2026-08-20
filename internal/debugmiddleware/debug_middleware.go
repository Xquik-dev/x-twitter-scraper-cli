// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package debugmiddleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
)

// These aliases keep redaction tests independent from generated code.
type (
	Middleware     = func(*http.Request, MiddlewareNext) (*http.Response, error)
	MiddlewareNext = func(*http.Request) (*http.Response, error)
)

const redactedPlaceholder = "<REDACTED>"

// Authorization needs separate redaction, so this list excludes it.
var sensitiveHeaders = []string{
	"api-key",
	"x-api-key",
	"cookie",
	"set-cookie",
}

// RequestLogger is a middleware that logs HTTP requests and responses.
type RequestLogger struct {
	logger           interface{ Printf(string, ...any) } // field for testability; usually log.Default()
	sensitiveHeaders []string                            // field for testability; usually sensitiveHeaders
}

// NewRequestLogger returns a logger with default redaction.
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{
		logger:           log.Default(),
		sensitiveHeaders: sensitiveHeaders,
	}
}

func (m *RequestLogger) Middleware() Middleware {
	return func(req *http.Request, mn MiddlewareNext) (*http.Response, error) {
		redacted, err := m.redactRequest(req)
		if err != nil {
			return nil, err
		}
		if reqBytes, err := httputil.DumpRequest(redacted, true); err == nil {
			m.logger.Printf("HTTP Request:\n%s\n", reqBytes)
		}

		resp, err := mn(req)
		if err != nil {
			return resp, err
		}

		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			m.logger.Printf("HTTP Response:\n%s\n", respBytes)
		}

		return resp, err
	}
}

// redactRequest clones requests only when a sensitive header needs redaction.
func (m *RequestLogger) redactRequest(req *http.Request) (*http.Request, error) {
	redactedHeaders := req.Header.Clone()

	// Redact every value when a header appears more than once.
	if values := redactedHeaders.Values("Authorization"); len(values) > 0 {
		redactedHeaders.Del("Authorization")

		for _, value := range values {
			// Keep the authorization scheme for debugging.
			if authKind, _, ok := strings.Cut(value, " "); ok {
				redactedHeaders.Add("Authorization", authKind+" "+redactedPlaceholder)
			} else {
				redactedHeaders.Add("Authorization", redactedPlaceholder)
			}
		}
	}

	for _, header := range m.sensitiveHeaders {
		values := redactedHeaders.Values(header)
		if len(values) == 0 {
			continue
		}

		redactedHeaders.Del(header)

		for range values {
			redactedHeaders.Add(header, redactedPlaceholder)
		}
	}

	if reflect.DeepEqual(req.Header, redactedHeaders) {
		return req, nil
	}

	redacted := req.Clone(req.Context())
	redacted.Header = redactedHeaders
	var err error
	redacted.Body, req.Body, err = cloneBody(req.Body)
	return redacted, err
}

// cloneBody returns 2 independent readers, following net/http/httputil.drainBody.
func cloneBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
