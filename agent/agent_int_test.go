//go:build integration

package agent

import (
	"testing"

	"github.com/rannday/kea-api/internal/testenv"
)

// TestIntegration_BuildReport verifies BuildReport returns a non-empty build info string.
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

// TestIntegration_StatusGet checks that StatusGet returns valid PID and uptime.
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

// TestIntegration_ListCommands ensures ListCommands returns a usable list of command names.
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

// TestIntegration_ConfigGet verifies ConfigGet returns a populated configuration and hash.
func TestIntegration_ConfigGet(t *testing.T) {
  t.Parallel()
  c := testenv.NewIntegrationClient()

  got, err := ConfigGet(c)
  if err != nil {
    t.Fatalf("ConfigGet() error = %v", err)
  }

  if got.ControlAgent.HTTPPort == 0 {
    t.Errorf("ConfigGet() returned empty HTTPPort: %+v", got.ControlAgent)
  }
  if got.Hash == "" {
    t.Error("ConfigGet() returned empty config hash")
  }
}

// TestIntegration_VersionGet checks that VersionGet returns both version text and extended info.
func TestIntegration_VersionGet(t *testing.T) {
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
