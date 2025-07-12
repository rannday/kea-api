//go:build integration

package tests

import (
	"testing"

	"github.com/rannday/kea-api/dhcp4"
)

func TestBuildReportDHCP4(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := dhcp4.BuildReport(client)
	if err != nil {
		t.Fatalf("dhcp4.BuildReport() error = %v", err)
	}

	if got == "" {
		t.Error("dhcp4.BuildReport() returned empty string")
	}
}

func TestStatusGetDHCP4(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := dhcp4.StatusGet(client)
	if err != nil {
		t.Fatalf("dhcp4.StatusGet() error = %v", err)
	}

	if got.PID == 0 {
		t.Errorf("dhcp4.StatusGet() returned zero PID: %+v", got)
	}
	if got.Uptime == 0 {
		t.Errorf("dhcp4.StatusGet() returned zero Uptime: %+v", got)
	}
}

func TestListCommandsDHCP4(t *testing.T) {
	t.Parallel()

	client := NewClient()

	got, err := dhcp4.ListCommands(client)
	if err != nil {
		t.Fatalf("dhcp4.ListCommands() error = %v", err)
	}

	if len(got) < 3 {
		t.Errorf("dhcp4.ListCommands() returned too few commands: %d", len(got))
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
		t.Errorf("dhcp4.ListCommands() missing expected command %q", want)
	}
}
