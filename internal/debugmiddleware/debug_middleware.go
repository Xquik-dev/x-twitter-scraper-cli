// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package debugmiddleware logs each HTTP request without its body or
// credentials, and each response without its cookies.
package debugmiddleware

import (
	"log"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
)

type (
	Middleware     = func(*http.Request, MiddlewareNext) (*http.Response, error)
	MiddlewareNext = func(*http.Request) (*http.Response, error)
)

const redactedPlaceholder = "<REDACTED>"

var sensitiveHeaders = []string{
	"api-key",
	"x-api-key",
	"cookie",
	"set-cookie",
}

type RequestLogger struct {
	logger           *log.Logger
	sensitiveHeaders []string
}

type redactedRequest struct {
	request *http.Request
}

func NewRequestLogger() *RequestLogger {
	return &RequestLogger{
		logger:           log.Default(),
		sensitiveHeaders: sensitiveHeaders,
	}
}

func (m *RequestLogger) Middleware() Middleware {
	return func(req *http.Request, mn MiddlewareNext) (*http.Response, error) {
		if reqBytes, err := httputil.DumpRequest(m.redact(req).request, false); err == nil {
			m.logger.Printf("HTTP Request:\n%s\n", reqBytes)
		}

		resp, err := mn(req)
		if err != nil {
			return resp, err
		}

		shown := *resp
		shown.Header = resp.Header.Clone()
		shown.Header.Del("Set-Cookie")
		respBytes, dumpErr := httputil.DumpResponse(&shown, true)
		resp.Body = shown.Body
		if dumpErr == nil {
			m.logger.Printf("HTTP Response:\n%s\n", respBytes)
		}

		return resp, err
	}
}

func (m *RequestLogger) redact(req *http.Request) redactedRequest {
	redactedHeaders := req.Header.Clone()

	if values := redactedHeaders.Values("Authorization"); len(values) > 0 {
		redactedHeaders.Del("Authorization")

		for _, value := range values {
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
		return redactedRequest{request: req}
	}

	redacted := req.Clone(req.Context())
	redacted.Header = redactedHeaders
	return redactedRequest{request: redacted}
}
