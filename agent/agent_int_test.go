//go:build integration

package agent

import (
	"testing"

	"github.com/rannday/kea-api/internal/testenv"
)

func TestIntegration_BuildReport(t *testing.T) {
	t.Parallel()

	c := testenv.NewIntegrationClient()

	got, err := BuildReport(c)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}

	if got == "" {
		t.Error("BuildReport() returned empty string")
	}
}

func TestIntegration_StatusGet(t *testing.T) {
	t.Parallel()

	c := testenv.NewIntegrationClient()

	got, err := StatusGet(c)
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

func TestIntegration_ListCommands(t *testing.T) {
	t.Parallel()

	c := testenv.NewIntegrationClient()

	got, err := ListCommands(c)
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
