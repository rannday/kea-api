package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"unsafe"
)

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

// TestTLSAuthConfigureClientFailsGracefully tests that TLSAuth.ConfigureClient fails gracefully when cert files are missing.
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

// TestWithAuthOption confirms that WithAuth sets the AuthProvider.
func TestWithAuthOption(t *testing.T) {
	auth := &BasicAuth{
		Username: "u",
		Password: "p",
	}
	c := NewHTTPClient("http://localhost:1234", WithAuth(auth))

	internalAuth := cAuth(c)
	if internalAuth == nil {
		t.Fatal("expected AuthProvider to be set")
	}

	req, _ := http.NewRequest("GET", "http://fake", nil)
	if err := internalAuth.Apply(req); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if got := req.Header.Get("Authorization"); got == "" {
		t.Error("Authorization header was not set by applied auth")
	}
}

// TestWithHTTPClientOption checks that a custom HTTP client is applied.
func TestWithHTTPClientOption(t *testing.T) {
	custom := &http.Client{}
	c := NewHTTPClient("http://x", WithHTTPClient(custom))
	if cHTTP(c) != custom {
		t.Error("custom http.Client not applied")
	}
}

// TestClientCall_UsesAuthAndDecodes tests that Call uses the auth provider and decodes the response.
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
	c := NewHTTPClient(server.URL, WithAuth(auth))

	var res []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &res)
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	if len(res) != 1 || res[0].Result != 0 {
		t.Errorf("unexpected response: %+v", res)
	}
}

// TestClientCall_RejectsFailureResult tests that Call rejects responses with non-success results.
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

	c := NewHTTPClient(server.URL)

	var res []CommandResponse
	err := c.Call(CommandRequest{Command: "status-xyz"}, &res)
	if err == nil {
		t.Fatal("expected error due to unsupported command, got nil")
	}

	if got := err.Error(); got != "unsupported command: unsupported command: status-xyz" {
		t.Errorf("unexpected error: %v", got)
	}
}

// --- Helpers ---

// cAuth returns the AuthProvider from the client using reflection.
func cAuth(c *Client) AuthProvider {
	return reflectGet[AuthProvider](c, "auth")
}

// cHTTP returns the *http.Client from the client using reflection.
func cHTTP(c *Client) *http.Client {
	return reflectGet[*http.Client](c, "httpClient")
}

// reflectGet is a helper to get an unexported field from Client.
func reflectGet[T any](c *Client, field string) T {
	v := reflect.ValueOf(c).Elem().FieldByName(field)

	// Create a pointer to the field using unsafe
	ptr := unsafe.Pointer(v.UnsafeAddr())
	val := *(*T)(ptr)

	return val
}
