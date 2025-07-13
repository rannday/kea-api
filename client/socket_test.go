package client

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// startSocketServer starts a TCP or UNIX socket server that responds with the given payload.
// If skipDecode is true, the request is not read from the client.
func startSocketServer(t *testing.T, network, address string, response any, skipDecode bool) net.Listener {
	t.Helper()

	if network == "unix" {
		_ = os.Remove(address) // Remove stale socket
	}

	l, err := net.Listen(network, address)
	if err != nil {
		t.Fatalf("failed to start %s listener: %v", network, err)
	}

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Logf("server accept error: %v", err)
			return
		}
		defer conn.Close()

		if !skipDecode {
			var req CommandRequest
			if err := json.NewDecoder(conn).Decode(&req); err != nil {
				t.Logf("server decode error: %v", err)
				return
			}
		}

		if response != nil {
			if err := json.NewEncoder(conn).Encode(response); err != nil {
				t.Logf("server encode error: %v", err)
			}
		}
	}()

	return l
}

// startGarbageServer starts a socket server that responds with invalid JSON for testing decode errors.
func startGarbageServer(t *testing.T, network, address string) net.Listener {
	t.Helper()

	if network == "unix" {
		_ = os.Remove(address)
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

		// Read request (optional)
		_, _ = json.NewDecoder(conn).Token()

		// Send invalid JSON
		conn.Write([]byte("{bad json"))
	}()

	return l
}

// TempSocketPath returns a temp file path for use as a UNIX socket.
func TempSocketPath(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, name+".sock")
}

// TestSocketTransport_TCP verifies TCP socket communication.
func TestSocketTransport_TCP(t *testing.T) {
	response := []CommandResponse{{
		Result:    ResultSuccess,
		Arguments: json.RawMessage(`{}`),
	}}

	l := startSocketServer(t, "tcp", "127.0.0.1:12345", response, false)
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

	l := startSocketServer(t, "unix", socketPath, response, false)
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

// TestSocketTransport_EncodeError ensures encoding failures return an appropriate error.
func TestSocketTransport_EncodeError(t *testing.T) {
	l := startSocketServer(t, "tcp", "127.0.0.1:12348", nil, true)
	defer l.Close()

	tr := NewSocketTransport("tcp", "127.0.0.1:12348", 1*time.Second)
	c := NewClient(tr)

	badArgs := map[string]interface{}{
		"bad": make(chan int),
	}

	err := c.Call(CommandRequest{
		Command:   "status-get",
		Arguments: badArgs,
	}, &[]CommandResponse{})

	if err == nil || err.Error() != "encode request: json: unsupported type: chan int" {
		t.Errorf("expected encode error, got: %v", err)
	}
}

// TestSocketTransport_DecodeError ensures decoding failures return an appropriate error.
func TestSocketTransport_DecodeError(t *testing.T) {
	l := startGarbageServer(t, "tcp", "127.0.0.1:12349")
	defer l.Close()

	tr := NewSocketTransport("tcp", "127.0.0.1:12349", 1*time.Second)
	c := NewClient(tr)

	var out []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || !contains(err.Error(), "decode response") {
		t.Errorf("expected decode error, got: %v", err)
	}
}

// TestSocketTransport_EmptyResponse checks that an empty response triggers an error.
func TestSocketTransport_EmptyResponse(t *testing.T) {
	response := []CommandResponse{} // empty

	l := startSocketServer(t, "tcp", "127.0.0.1:12347", response, false)
	defer l.Close()

	tr := NewSocketTransport("tcp", "127.0.0.1:12347", 2*time.Second)
	c := NewClient(tr)

	var out []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &out)
	if err == nil || err.Error() != "empty response from Kea" {
		t.Errorf("expected empty response error, got: %v", err)
	}
}

// TestSocketTransport_CommandResponseHandling ensures ResultError logic is triggered
// when the response contains a non-zero result code.
func TestSocketTransport_CommandResponseHandling(t *testing.T) {
	response := []CommandResponse{
		{
			Result: ResultGeneralFailure,
			Text:   "custom failure",
		},
	}

	l := startSocketServer(t, "tcp", "127.0.0.1:12351", response, false)
	defer l.Close()

	tr := NewSocketTransport("tcp", "127.0.0.1:12351", 1*time.Second)
	c := NewClient(tr)

	var out []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || err.Error() != "general error: custom failure" {
		t.Errorf("expected custom failure error, got: %v", err)
	}
}
