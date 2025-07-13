//go:build integration

package client

import (
	"os"
	"testing"
)

// Only defined locally so we don't import testenv and cause a cycle.
func newIntegrationClient() *Client {
	url := os.Getenv("KEA_API_URL")
	if url == "" {
		url = "http://localhost:8000"
	}

	auth := &BasicAuth{
		Username: "kea",
		Password: "kea",
	}

	return NewHTTP(url, WithAuth(auth))
}

// TestBasicAuthLive checks that BasicAuth works with Kea's control agent.
func TestBasicAuthLive(t *testing.T) {
	c := newIntegrationClient()

	var out []CommandResponse
	err := c.Call(CommandRequest{Command: "status-get"}, &out)
	if err != nil {
		t.Fatalf("Call failed with BasicAuth: %v", err)
	}

	if len(out) == 0 || out[0].Result != ResultSuccess {
		t.Errorf("unexpected response: %+v", out)
	}
}
