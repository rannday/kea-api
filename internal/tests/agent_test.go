//go:build integration

package tests

import (
	"testing"

	"github.com/rannday/kea-api/agent"
)

// TestBuildReport checks the actual response from the control agent for build-report.
func TestBuildReport(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := agent.BuildReport(client)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}

	if got == "" {
		t.Error("BuildReport() returned empty string")
	}
}

// TestStatusGet checks that the control agent returns valid status data.
func TestStatusGet(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := agent.StatusGet(client)
	if err != nil {
		t.Fatalf("StatusGet() error = %v", err)
	}

	if got.PID == 0 {
		t.Errorf("StatusGet() returned zero PID: %+v", got)
	}
	if got.Uptime == 0 {
		t.Errorf("StatusGet() returned zero Uptime: %+v", got)
	}
}

// TestListCommands ensures the control agent returns a full command list.
func TestListCommands(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := agent.ListCommands(client)
	if err != nil {
		t.Fatalf("ListCommands() error = %v", err)
	}

	if len(got) < 3 {
		t.Errorf("ListCommands() returned too few commands: %d", len(got))
	}

	want := "status-get"
	found := false
	for _, cmd := range got {
		if cmd == want {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ListCommands() missing expected command %q", want)
	}
}
