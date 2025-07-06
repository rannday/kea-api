package testenv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rannday/kea-api/client"
)

// MustEncodeRawJSON marshals a value into json.RawMessage.
func MustEncodeRawJSON(t *testing.T, v interface{}) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	return json.RawMessage(b)
}

// NewTestClient returns a mock *client.Client that expects a specific request
// (validated by the test function) and returns the given responses.
func NewTestClient(
	t *testing.T,
	validate func(t *testing.T, req client.CommandRequest),
	responses []client.CommandResponse,
) *client.Client {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req client.CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		validate(t, req)

		if err := json.NewEncoder(w).Encode(responses); err != nil {
			t.Errorf("failed to encode mock response: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	})

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return client.NewHTTP(server.URL)
}

// ExpectCommand returns a function that can be used to validate a CommandRequest.
func ExpectCommand(t *testing.T, wantCommand string, wantService ...client.Service) func(*testing.T, client.CommandRequest) {
	return func(t *testing.T, req client.CommandRequest) {
		t.Helper()
		if req.Command != wantCommand {
			t.Errorf("unexpected command: got %q, want %q", req.Command, wantCommand)
		}
		if len(wantService) == 0 {
			if req.Service != nil {
				t.Errorf("expected no service, got: %v", req.Service)
			}
		} else {
			if len(req.Service) != len(wantService) || req.Service[0] != wantService[0] {
				t.Errorf("unexpected service: got %v, want %v", req.Service, wantService)
			}
		}
	}
}
