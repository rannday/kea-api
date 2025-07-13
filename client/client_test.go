package client

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// StartSocketServer starts a TCP or UNIX socket server that responds with the given payload.
func startSocketServer(t *testing.T, network, address string, response any) net.Listener {
	t.Helper()

	if network == "unix" {
		_ = os.Remove(address) // Ensure no stale socket
	}

	l, err := net.Listen(network, address)
	if err != nil {
		t.Fatalf("failed to start %s listener: %v", network, err)
	}

	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		var req CommandRequest
		if err := json.NewDecoder(conn).Decode(&req); err != nil {
			return
		}

		_ = json.NewEncoder(conn).Encode(response)
	}()

	return l
}

// TempSocketPath returns a temp file path for use as a UNIX socket.
func TempSocketPath(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, name+".sock")
}

// TestClientCall_RejectsFailureResult tests that an error is returned for non-success result.
func TestClientCall_RejectsFailureResult(t *testing.T) {
	response := []CommandResponse{{
		Result:    ResultUnsupported,
		Text:      "unsupported command: status-xyz",
		Arguments: json.RawMessage(`{}`),
	}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := NewHTTP(server.URL)

	var res []CommandResponse
	err := c.Call(CommandRequest{Command: "status-xyz"}, &res)
	if err == nil {
		t.Fatal("expected error due to unsupported command, got nil")
	}

	want := "unsupported command: unsupported command: status-xyz"
	if got := err.Error(); got != want {
		t.Errorf("unexpected error: %v", got)
	}
}

// TestSocketTransport_TCP verifies TCP socket communication.
func TestSocketTransport_TCP(t *testing.T) {
	response := []CommandResponse{{
		Result:    ResultSuccess,
		Arguments: json.RawMessage(`{}`),
	}}

	l := startSocketServer(t, "tcp", "127.0.0.1:12345", response)
	defer l.Close()

	tr := NewSocketTransport("tcp", "127.0.0.1:12345", 2*time.Second)
	c := NewClient(tr)

	var out []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &out)
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}
	if len(out) != 1 || out[0].Result != ResultSuccess {
		t.Errorf("unexpected result: %+v", out)
	}
}

// TestSocketTransport_UNIX verifies UNIX socket communication.
func TestSocketTransport_UNIX(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("UNIX socket test skipped on Windows")
	}

	socketPath := "/tmp/kea-test.sock"
	_ = os.Remove(socketPath)

	response := []CommandResponse{{
		Result:    ResultSuccess,
		Arguments: json.RawMessage(`{}`),
	}}

	l := startSocketServer(t, "unix", socketPath, response)
	defer func() {
		l.Close()
		_ = os.Remove(socketPath)
	}()

	tr := NewSocketTransport("unix", socketPath, 2*time.Second)
	c := NewClient(tr)

	var out []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &out)
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}
	if len(out) != 1 || out[0].Result != ResultSuccess {
		t.Errorf("unexpected result: %+v", out)
	}
}

// TestSocketTransport_ConnectError verifies connection error handling.
func TestSocketTransport_ConnectError(t *testing.T) {
	tr := NewSocketTransport("tcp", "127.0.0.1:65000", 1*time.Second)
	c := NewClient(tr)

	var out []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &out)
	if err == nil {
		t.Fatal("expected error on connect, got nil")
	}
}
