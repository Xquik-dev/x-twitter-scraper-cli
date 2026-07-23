package mocktest

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

type Runner func(*testing.T, context.Context, []string, []byte) error

var runCLI Runner
var commandTestMutex sync.Mutex
var responseStatus atomic.Int64

var mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	_, _ = io.Copy(io.Discard, r.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(responseStatus.Load()))
	_, _ = w.Write([]byte("{}"))
}))

func restrictNetworkToMockServer() func() {
	dialer := &net.Dialer{}
	allowedAddress := mockServer.Listener.Addr().String()
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			if address != allowedAddress {
				return nil, fmt.Errorf("blocked test network connection to %s", address)
			}
			return dialer.DialContext(ctx, network, address)
		},
	}

	previousTransport := http.DefaultTransport
	http.DefaultTransport = transport
	return func() {
		http.DefaultTransport = previousTransport
		transport.CloseIdleConnections()
	}
}

func RegisterRunner(runner Runner) {
	runCLI = runner
}

// TestRunMockTestWithFlags runs a test against a mock server with the provided
// CLI args and ensures it succeeds
func TestRunMockTestWithFlags(t *testing.T, args ...string) {
	t.Helper()
	TestRunMockTestWithPipeAndFlags(t, nil, args...)
}

// TestRunMockTestWithPipeAndFlags runs a test against a mock server with the provided
// data piped over stdin and CLI args and ensures it succeeds
func TestRunMockTestWithPipeAndFlags(t *testing.T, pipeData []byte, args ...string) {
	t.Helper()
	require.NotNil(t, runCLI, "Register a CLI test runner before running command tests")
	commandTestMutex.Lock()
	defer commandTestMutex.Unlock()
	defer restrictNetworkToMockServer()()

	responseStatus.Store(http.StatusOK)
	commandArgs := append(
		[]string{"x-twitter-scraper", "--base-url", mockServer.URL},
		args...,
	)
	require.NoError(t, runCLI(t, context.Background(), commandArgs, pipeData))

	unexpectedArgs := append(append([]string{}, commandArgs...), "unexpected")
	require.Error(t, runCLI(t, context.Background(), unexpectedArgs, pipeData))

	require.Error(t, runCLI(t, context.Background(), commandArgs, []byte("{")))

	responseStatus.Store(http.StatusBadRequest)
	require.Error(t, runCLI(t, context.Background(), commandArgs, pipeData))
}

func TestFile(t *testing.T, contents string) string {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "file.txt")
	require.NoError(t, os.WriteFile(filename, []byte(contents), 0644))
	return filename
}
