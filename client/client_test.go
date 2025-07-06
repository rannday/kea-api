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

// TestBasicAuthApply tests the Apply method of BasicAuth.
func TestBasicAuthApply(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	auth := &BasicAuth{
		Username: "user",
		Password: "pass",
	}

	if err := auth.Apply(req); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	want := "Basic dXNlcjpwYXNz" // base64("user:pass")
	if got := req.Header.Get("Authorization"); got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
}

// TestTLSAuthConfigureClientFailsGracefully ensures TLSAuth fails on missing cert files.
func TestTLSAuthConfigureClientFailsGracefully(t *testing.T) {
	auth := &TLSAuth{
		CertFile: "testdata/missing.crt",
		KeyFile:  "testdata/missing.key",
		CAFile:   "testdata/missing_ca.crt",
	}
	hc := &http.Client{}

	err := auth.ConfigureClient(hc)
	if err == nil {
		t.Fatal("expected error due to missing files, got nil")
	}
}

// TestWithAuthOption checks that WithAuth applies the Authorization header.
func TestWithAuthOption(t *testing.T) {
	auth := &BasicAuth{Username: "u", Password: "p"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := "Basic dTpw"
		got := r.Header.Get("Authorization")
		if got != want {
			t.Errorf("Authorization header = %q, want %q", got, want)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]CommandResponse{{
			Result:    ResultSuccess,
			Arguments: json.RawMessage(`{}`),
		}})
	}))
	defer server.Close()

	c := NewHTTP(server.URL, WithAuth(auth))
	var res []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &res)
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}
}

// TestClientCall_UsesAuthAndDecodes tests full request/response with auth and decoding.
func TestClientCall_UsesAuthAndDecodes(t *testing.T) {
	response := []CommandResponse{{
		Result:    ResultSuccess,
		Arguments: json.RawMessage(`{}`),
	}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Basic Zm9vOmJhcg==" {
			t.Errorf("unexpected auth header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	auth := &BasicAuth{Username: "foo", Password: "bar"}
	c := NewHTTP(server.URL, WithAuth(auth))

	var res []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &res)
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	if len(res) != 1 || res[0].Result != 0 {
		t.Errorf("unexpected response: %+v", res)
	}
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
