//go:build integration

package dhcp4

import (
	"testing"

	"github.com/rannday/kea-api/internal/testenv"
)

// TestIntegration_BuildReportDHCP4 verifies BuildReport returns non-empty build details.
func TestIntegration_BuildReportDHCP4(t *testing.T) {
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

// TestIntegration_StatusGetDHCP4 checks that StatusGet returns valid PID and uptime.
func TestIntegration_StatusGetDHCP4(t *testing.T) {
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

// TestIntegration_ListCommandsDHCP4 ensures ListCommands returns known DHCPv4 commands.
func TestIntegration_ListCommandsDHCP4(t *testing.T) {
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

// TestIntegration_ConfigGetDHCP4 verifies ConfigGet returns a valid DHCPv4 config and hash.
func TestIntegration_ConfigGetDHCP4(t *testing.T) {
	t.Parallel()

	c := testenv.NewIntegrationClient()

	got, err := ConfigGet(c)
	if err != nil {
		t.Fatalf("ConfigGet() error = %v", err)
	}

	if got.Hash == "" {
		t.Error("ConfigGet() returned empty config hash")
	}
	if got.Dhcp4.Allocator == "" && got.Dhcp4.ServerTag == "" {
		t.Errorf("ConfigGet() returned unexpected Dhcp4 block: %+v", got.Dhcp4)
	}
}

// TestIntegration_VersionGetDHCP4 checks VersionGet returns both short and extended version info.
func TestIntegration_VersionGetDHCP4(t *testing.T) {
	t.Parallel()

	c := testenv.NewIntegrationClient()

	gotText, gotVersion, err := VersionGet(c)
	if err != nil {
		t.Fatalf("VersionGet() error = %v", err)
	}

	if gotText == "" {
		t.Error("VersionGet() returned empty text")
	}
	if gotVersion.Extended == "" {
		t.Errorf("VersionGet() returned empty extended version: %+v", gotVersion)
	}
}
