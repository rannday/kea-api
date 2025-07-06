//go:build integration

package tests

import (
	"os"
	"testing"

	"github.com/rannday/kea-api/client"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func keaURL() string {
	if url := os.Getenv("KEA_API_URL"); url != "" {
		return url
	}
	return "http://localhost:8000"
}

// TestBasicAuthLive checks that BasicAuth works with Kea's control agent
func TestBasicAuthLive(t *testing.T) {
	auth := &client.BasicAuth{
		Username: "kea",
		Password: "kea", // matches contents of kea-api-password
	}
	c := client.NewHTTP(keaURL(), client.WithAuth(auth))

	var out []client.CommandResponse
	err := c.Call(client.CommandRequest{Command: "status-get"}, &out)
	if err != nil {
		t.Fatalf("Call failed with BasicAuth: %v", err)
	}

	if len(out) == 0 || out[0].Result != client.ResultSuccess {
		t.Errorf("unexpected response: %+v", out)
	}
}
