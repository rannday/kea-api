//go:build integration

package tests

import (
	"testing"

	"github.com/rannday/kea-api/dhcp6"
)

func TestBuildReportDHCP6(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := dhcp6.BuildReport(client)
	if err != nil {
		t.Fatalf("dhcp6.BuildReport() error = %v", err)
	}

	if got == "" {
		t.Error("dhcp6.BuildReport() returned empty string")
	}
}

func TestStatusGetDHCP6(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := dhcp6.StatusGet(client)
	if err != nil {
		t.Fatalf("dhcp6.StatusGet() error = %v", err)
	}

	if got.PID == 0 {
		t.Errorf("dhcp6.StatusGet() returned zero PID: %+v", got)
	}
	if got.Uptime == 0 {
		t.Errorf("dhcp6.StatusGet() returned zero Uptime: %+v", got)
	}
}

func TestListCommandsDHCP6(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := dhcp6.ListCommands(client)
	if err != nil {
		t.Fatalf("dhcp6.ListCommands() error = %v", err)
	}

	if len(got) < 3 {
		t.Errorf("dhcp6.ListCommands() returned too few commands: %d", len(got))
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
		t.Errorf("dhcp6.ListCommands() missing expected command %q", want)
	}
}
