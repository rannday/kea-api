package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWithHTTPClientOption verifies that WithHTTPClient sets the custom HTTP client.
func TestWithHTTPClientOption(t *testing.T) {
	custom := &http.Client{}
	tp := NewHTTPTransport("http://localhost", WithHTTPClient(custom))

	if tp.httpClient != custom {
		t.Errorf("expected custom httpClient to be set")
	}
}

// TestHTTPTransport_ApplyAuthError checks that auth errors are wrapped and returned.
func TestHTTPTransport_ApplyAuthError(t *testing.T) {
	badAuth := AuthProviderFunc(func(r *http.Request) error {
		return errors.New("boom")
	})

	tp := NewHTTPTransport("http://localhost", WithAuth(badAuth))

	var out []CommandResponse
	err := tp.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || err.Error() != "apply auth: boom" {
		t.Errorf("expected auth error, got %v", err)
	}
}

// TestHTTPTransport_Success verifies full successful path including result decoding.
func TestHTTPTransport_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]CommandResponse{{
			Result:    ResultSuccess,
			Arguments: json.RawMessage(`{}`),
		}})
	}))
	defer srv.Close()

	tp := NewHTTPTransport(srv.URL)
	var out []CommandResponse
	err := tp.Call(CommandRequest{Command: "status-get"}, &out)

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(out) != 1 || out[0].Result != ResultSuccess {
		t.Errorf("unexpected result: %+v", out)
	}
}

// TestHTTPTransport_RequestFail checks that invalid URL returns error.
func TestHTTPTransport_RequestFail(t *testing.T) {
	tp := NewHTTPTransport("http://[::1]:namedport") // invalid URL format
	var out []CommandResponse
	err := tp.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || !contains(err.Error(), "create request") {
		t.Errorf("expected request creation error, got %v", err)
	}
}

// TestHTTPCall_RejectsFailureResult tests that an error is returned for non-success result.
func TestHTTPCall_RejectsFailureResult(t *testing.T) {
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

// TestHTTPTransport_MarshalFail ensures marshal errors are returned properly.
func TestHTTPTransport_MarshalFail(t *testing.T) {
	tp := NewHTTPTransport("http://localhost")
	var out []CommandResponse

	// Create a map with a value that can't be marshaled to JSON (e.g., a channel).
	args := map[string]interface{}{
		"bad": make(chan int),
	}

	req := CommandRequest{
		Command:   "status-get",
		Arguments: args,
	}

	err := tp.Call(req, &out)
	if err == nil || !contains(err.Error(), "marshal request") {
		t.Errorf("expected marshal error, got: %v", err)
	}
}

// TestHTTPTransport_SendFail checks HTTP client Do() failure.
func TestHTTPTransport_SendFail(t *testing.T) {
	// No server running on this port
	tp := NewHTTPTransport("http://127.0.0.1:65534")
	var out []CommandResponse
	err := tp.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || !contains(err.Error(), "send request") {
		t.Errorf("expected send error, got %v", err)
	}
}

// TestHTTPTransport_HTTPError checks HTTP 500 responses are caught.
func TestHTTPTransport_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "kaboom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	tp := NewHTTPTransport(srv.URL)
	var out []CommandResponse
	err := tp.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || !contains(err.Error(), "kea error: status=500") {
		t.Errorf("expected HTTP error, got %v", err)
	}
}

// TestHTTPTransport_BadJSONResponse checks that bad JSON is caught.
func TestHTTPTransport_BadJSONResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	tp := NewHTTPTransport(srv.URL)
	var out []CommandResponse
	err := tp.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || !contains(err.Error(), "unmarshal response") {
		t.Errorf("expected unmarshal error, got %v", err)
	}
}

// TestHTTPTransport_EmptyResponse checks for empty []CommandResponse error.
func TestHTTPTransport_EmptyResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode([]CommandResponse{})
	}))
	defer srv.Close()

	tp := NewHTTPTransport(srv.URL)
	var out []CommandResponse
	err := tp.Call(CommandRequest{Command: "status-get"}, &out)

	if err == nil || err.Error() != "empty response from Kea" {
		t.Errorf("expected empty response error, got %v", err)
	}
}

// helper for inline AuthProvider mocks
type AuthProviderFunc func(*http.Request) error

func (f AuthProviderFunc) Apply(r *http.Request) error { return f(r) }

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
