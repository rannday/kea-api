//go:build integration

package dhcp6

import (
	"testing"

	"github.com/rannday/kea-api/internal/testenv"
)

// TestIntegration_BuildReportDHCP6 verifies BuildReport returns non-empty build info.
func TestIntegration_BuildReportDHCP6(t *testing.T) {
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

// TestIntegration_StatusGetDHCP6 checks that StatusGet returns valid PID and uptime.
func TestIntegration_StatusGetDHCP6(t *testing.T) {
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

// TestIntegration_ListCommandsDHCP6 ensures ListCommands returns known DHCPv6 commands.
func TestIntegration_ListCommandsDHCP6(t *testing.T) {
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

// TestIntegration_ConfigGetDHCP6 verifies ConfigGet returns a populated config and hash.
func TestIntegration_ConfigGetDHCP6(t *testing.T) {
	t.Parallel()

	c := testenv.NewIntegrationClient()

	got, err := ConfigGet(c)
	if err != nil {
		t.Fatalf("ConfigGet() error = %v", err)
	}

	if got.Hash == "" {
		t.Error("ConfigGet() returned empty config hash")
	}
	if got.Dhcp6.ServerTag == "" && got.Dhcp6.Allocator == "" {
		t.Errorf("ConfigGet() returned unexpected Dhcp6 block: %+v", got.Dhcp6)
	}
}

// TestIntegration_VersionGetDHCP6 checks VersionGet returns version text and extended details.
func TestIntegration_VersionGetDHCP6(t *testing.T) {
	t.Parallel()

	c := testenv.NewIntegrationClient()

	gotText, gotVersion, err := VersionGet(c)
	if err != nil {
		t.Fatalf("VersionGet() error = %v", err)
	}

	if gotText == "" {
		t.Error("VersionGet() returned empty version text")
	}
	if gotVersion.Extended == "" {
		t.Errorf("VersionGet() returned empty extended version: %+v", gotVersion)
	}
}
