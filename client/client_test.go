package client

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// TestTimeoutOrDefault_EmptySlice checks that the default timeout is returned when none is provided.
func TestTimeoutOrDefault_EmptySlice(t *testing.T) {
	got, err := timeoutOrDefault(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := 5 * time.Second
	if got != want {
		t.Errorf("expected %v, got %v", want, got)
	}
}

// TestTimeoutOrDefault_ValidShort checks that a valid short timeout is returned as-is.
func TestTimeoutOrDefault_ValidShort(t *testing.T) {
	got, err := timeoutOrDefault([]time.Duration{2 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2*time.Second {
		t.Errorf("expected 2s, got %v", got)
	}
}

// TestTimeoutOrDefault_ZeroDuration verifies that zero timeout returns an error.
func TestTimeoutOrDefault_ZeroDuration(t *testing.T) {
	_, err := timeoutOrDefault([]time.Duration{0})
	if err == nil {
		t.Fatal("expected error for zero timeout, got nil")
	}
	if !strings.Contains(err.Error(), "greater than zero") {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestTimeoutOrDefault_LongTimeoutWarning captures and checks for a warning on long timeout values.
func TestTimeoutOrDefault_LongTimeoutWarning(t *testing.T) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	_, err := timeoutOrDefault([]time.Duration{61 * time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stderr = oldStderr

	output, _ := io.ReadAll(r)
	got := string(output)

	if !strings.Contains(got, "unusually long socket timeout") {
		t.Errorf("expected warning, got: %q", got)
	}
}

// TestNewSocket_Valid checks that NewSocket creates a client successfully with valid inputs.
func TestNewSocket_Valid(t *testing.T) {
	c, err := NewSocket("tcp", "127.0.0.1:1234", 1*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

// TestNewSocket_InvalidTimeout verifies that NewSocket returns an error for zero timeout.
func TestNewSocket_InvalidTimeout(t *testing.T) {
	c, err := NewSocket("unix", "/tmp/fake.sock", 0)
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
	if c != nil {
		t.Errorf("expected nil client, got %+v", c)
	}
}

// TestResultError checks all branches of ResultError including the default fallback.
func TestResultError(t *testing.T) {
	tests := []struct {
		code     ResultCode
		text     string
		wantErr  bool
		wantText string
	}{
		{ResultSuccess, "ok", false, ""},
		{ResultGeneralFailure, "failed", true, "general error: failed"},
		{ResultUnsupported, "bad cmd", true, "unsupported command: bad cmd"},
		{ResultNotFound, "nope", true, "resource not found: nope"},
		{ResultConflict, "duplicate", true, "conflict: duplicate"},
		{ResultCode(99), "???", true, "unknown result code 99: ???"},
	}

	for _, tt := range tests {
		err := tt.code.ResultError(tt.text)
		if tt.wantErr {
			if err == nil {
				t.Errorf("expected error for code %d, got nil", tt.code)
			} else if err.Error() != tt.wantText {
				t.Errorf("unexpected error: got %q, want %q", err.Error(), tt.wantText)
			}
		} else {
			if err != nil {
				t.Errorf("expected nil for code %d, got: %v", tt.code, err)
			}
		}
	}
}
