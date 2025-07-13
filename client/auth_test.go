package client

import (
	"net/http"
	"os"
	"testing"
)

// TestBasicAuth_Apply sets the Authorization header correctly.
func TestBasicAuth_Apply(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	auth := &BasicAuth{
		Username: "admin",
		Password: "secret",
	}

	if err := auth.Apply(req); err != nil {
		t.Fatalf("Apply() returned error: %v", err)
	}

	want := "Basic YWRtaW46c2VjcmV0" // base64("admin:secret")
	got := req.Header.Get("Authorization")
	if got != want {
		t.Errorf("Authorization header = %q, want %q", got, want)
	}
}

// TestBasicAuth_Empty ensures empty credentials still set header.
func TestBasicAuth_EmptyFails(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost", nil)

	auth := &BasicAuth{
		Username: "",
		Password: "",
	}
	err := auth.Apply(req)
	if err == nil {
		t.Fatal("expected error for empty credentials, got nil")
	}

	want := "auth.Apply: empty credentials"
	if err.Error() != want {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}

// TestBasicAuth_NilRequest should error if request is nil.
func TestBasicAuth_NilRequest(t *testing.T) {
	auth := &BasicAuth{Username: "foo", Password: "bar"}
	err := auth.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
}

// TestTLSAuth_ConfigureClient_Success verifies that valid TLS files configure the client correctly.
func TestTLSAuth_ConfigureClient_Success(t *testing.T) {
	auth := &TLSAuth{
		CertFile: "testdata/client.crt",
		KeyFile:  "testdata/client.key",
		CAFile:   "testdata/ca.crt",
	}
	client := &http.Client{}

	err := auth.ConfigureClient(client)
	if err != nil {
		t.Fatalf("ConfigureClient failed: %v", err)
	}

	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.Transport)
	}
	if tr.TLSClientConfig == nil {
		t.Fatal("TLSClientConfig was nil")
	}
	if len(tr.TLSClientConfig.Certificates) != 1 {
		t.Errorf("expected 1 certificate, got %d", len(tr.TLSClientConfig.Certificates))
	}
}

// TestTLSAuth_ConfigureClient_Fails checks that missing TLS files cause configuration to fail.
func TestTLSAuth_ConfigureClient_Fails(t *testing.T) {
	tests := []struct {
		name    string
		auth    *TLSAuth
		wantErr string
	}{
		{
			name: "missing cert file",
			auth: &TLSAuth{
				CertFile: "testdata/missing.crt",
				KeyFile:  "testdata/client.key",
				CAFile:   "testdata/ca.crt",
			},
			wantErr: "no such file or directory",
		},
		{
			name: "missing key file",
			auth: &TLSAuth{
				CertFile: "testdata/client.crt",
				KeyFile:  "testdata/missing.key",
				CAFile:   "testdata/ca.crt",
			},
			wantErr: "no such file or directory",
		},
		{
			name: "missing CA file",
			auth: &TLSAuth{
				CertFile: "testdata/client.crt",
				KeyFile:  "testdata/client.key",
				CAFile:   "testdata/missing_ca.crt",
			},
			wantErr: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			err := tt.auth.ConfigureClient(client)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !os.IsNotExist(err) {
				t.Errorf("expected file-not-found error, got: %v", err)
			}
		})
	}
}
