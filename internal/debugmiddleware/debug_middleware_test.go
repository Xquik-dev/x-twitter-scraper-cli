// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package debugmiddleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDebugMiddleware(t *testing.T) {
	t.Parallel()

	setup := func() (*RequestLogger, *bytes.Buffer) {
		var (
			logBuf     bytes.Buffer
			middleware = NewRequestLogger()
		)
		middleware.logger = log.New(&logBuf, "", 0)
		return middleware, &logBuf
	}

	t.Run("DoesNotRedactMostHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		const stainlessUserAgent = "Stainless"

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("User-Agent", stainlessUserAgent)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			require.Equal(t, stainlessUserAgent, req.Header.Get("User-Agent"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "User-Agent: "+stainlessUserAgent)
	})

	const secretToken = "secret-token"

	t.Run("RedactsAuthorizationHeader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("Authorization", secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			require.Equal(t, secretToken, req.Header.Get("Authorization"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "Authorization: "+redactedPlaceholder)
	})

	t.Run("RedactsOnlySecretInAuthorizationHeader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("Authorization", "Bearer "+secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "Authorization: Bearer "+redactedPlaceholder)
	})

	t.Run("RedactsMultipleAuthorizationHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Add("Authorization", secretToken+"1")
		req.Header.Add("Authorization", secretToken+"2")

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			require.Equal(t, []string{secretToken + "1", secretToken + "2"}, req.Header.Values("Authorization"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)

		if strings.Count(logBuf.String(), "Authorization: "+redactedPlaceholder) != 2 {
			t.Error("expected exactly two redacted placeholders in authorization headers")
		}
	})

	const customAPIKeyHeader = "X-My-Api-Key"

	t.Run("RedactsSensitiveHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set(customAPIKeyHeader, secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			require.Equal(t, secretToken, req.Header.Get(customAPIKeyHeader))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), customAPIKeyHeader+": "+redactedPlaceholder)
	})

	t.Run("RedactsMultipleSensitiveHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Add(customAPIKeyHeader, secretToken+"1")
		req.Header.Add(customAPIKeyHeader, secretToken+"2")

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			require.Equal(t, []string{secretToken + "1", secretToken + "2"}, req.Header.Values(customAPIKeyHeader))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Equal(t, 2, strings.Count(logBuf.String(), customAPIKeyHeader+": "+redactedPlaceholder))
	})

	t.Run("RedactsTheDefaultSensitiveHeaders", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()

		req := httptest.NewRequest("GET", "https://example.com", nil)
		req.Header.Set("Authorization", "Bearer "+secretToken+"-bearer")
		req.Header.Set("Cookie", "session="+secretToken+"-cookie")
		req.Header.Set("Set-Cookie", "session="+secretToken+"-set-cookie")
		req.Header.Set("X-Api-Key", secretToken+"-x-api-key")
		req.Header.Set("Api-Key", secretToken+"-api-key")

		_, err := middleware.Middleware()(req, func(*http.Request) (*http.Response, error) {
			return &http.Response{}, nil
		})
		require.NoError(t, err)

		logged := logBuf.String()
		require.NotContains(t, logged, secretToken)
		for _, header := range []string{"Authorization: Bearer ", "Cookie: ", "Set-Cookie: ", "X-Api-Key: ", "Api-Key: "} {
			require.Contains(t, logged, header+redactedPlaceholder)
		}
		require.Equal(t, "session="+secretToken+"-cookie", req.Header.Get("Cookie"))
	})

	t.Run("LeavesTheRequestBodyOutOfTheDump", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()
		const body = `{"email":"e","password":"` + secretToken + `-password","totpSecret":"` + secretToken + `-totp"}`
		req := httptest.NewRequest("POST", "https://example.com/x/accounts", strings.NewReader(body))

		_, err := middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			sent, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, body, string(sent))
			return &http.Response{Header: http.Header{}, Body: http.NoBody}, nil
		})
		require.NoError(t, err)
		require.Contains(t, logBuf.String(), "POST https://example.com/x/accounts")
		require.NotContains(t, logBuf.String(), secretToken)
	})

	t.Run("LeavesTheResponseSetCookieOutOfTheDump", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()
		req := httptest.NewRequest("GET", "https://example.com", nil)
		header := http.Header{}
		header.Set("Set-Cookie", "session="+secretToken)
		header.Set("Content-Type", "application/json")

		resp, err := middleware.Middleware()(req, func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     header,
				Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
			}, nil
		})
		require.NoError(t, err)
		require.NotContains(t, logBuf.String(), secretToken)
		require.Contains(t, logBuf.String(), `{"ok":true}`)
		require.Equal(t, "session="+secretToken, resp.Header.Get("Set-Cookie"))
		received, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, `{"ok":true}`, string(received))
	})

	t.Run("DoesNotConsumeRequestBodyWhenIoReader", func(t *testing.T) {
		t.Parallel()

		middleware, logBuf := setup()
		middleware.sensitiveHeaders = []string{customAPIKeyHeader}

		const bodyContent = "test request body content"
		bodyReader := strings.NewReader(bodyContent)

		req := httptest.NewRequest("POST", "https://example.com", bodyReader)
		req.Header.Set("Authorization", secretToken)

		var nextMiddlewareRan bool
		middleware.Middleware()(req, func(req *http.Request) (*http.Response, error) {
			nextMiddlewareRan = true

			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, bodyContent, string(body))

			require.Equal(t, secretToken, req.Header.Get("Authorization"))

			return &http.Response{}, nil
		})

		require.True(t, nextMiddlewareRan)
		require.Contains(t, logBuf.String(), "Authorization: "+redactedPlaceholder)
	})
}
