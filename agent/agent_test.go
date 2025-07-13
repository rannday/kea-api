package agent

import (
	"reflect"
	"testing"

	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/internal/testenv"
	"github.com/rannday/kea-api/types"
)

// TestBuildReport tests the BuildReport function for the CtrlAgent type.
func TestBuildReport(t *testing.T) {
	t.Parallel()

	want := `Kea source configure results:
Package:
  Name: kea
  Version: 2.6.3
`

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "build-report", client.Services.Agent),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Text:      want,
			Arguments: testenv.MustEncodeRawJSON(t, map[string]any{}),
		}},
	)

	got, err := BuildReport(mockClient)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}

	if got != want {
		t.Errorf("BuildReport() = %q, want %q", got, want)
	}
}

// TestConfigGet tests the ConfigGet function for the CtrlAgentConfig type.
func TestConfigGet(t *testing.T) {
	t.Parallel()

	want := CtrlAgentConfig{
		ControlAgent: CtrlAgentBlock{
			Authentication: AuthConfig{
				Clients: []AuthClient{{
					User:         "kea-api",
					PasswordFile: "kea-api-password",
				}},
				Directory: "/etc/kea",
				Realm:     "Kea Control Agent",
				Type:      "basic",
			},
			ControlSockets: map[string]types.SocketConfig{
				"dhcp4": {SocketName: "kea4-ctrl-socket", SocketType: "unix"},
				"dhcp6": {SocketName: "kea6-ctrl-socket", SocketType: "unix"},
				"d2":    {SocketName: "kea-ddns-ctrl-socket", SocketType: "unix"},
			},
			HooksLibraries: []types.HookLibrary{},
			HTTPHost:       "0.0.0.0",
			HTTPPort:       8000,
			Loggers: []types.LoggerConfig{{
				Name:       "kea-ctrl-agent",
				Severity:   "INFO",
				DebugLevel: 0,
				OutputOptions: []types.LogOutputOption{{
					Output:  "stdout",
					Pattern: "%-5p %m\n",
					Flush:   true,
				}},
			}},
		},
		Hash: "07FC3D1A1717B92D5A3D7E2E4900BCAF9916C8DDA3DB013390C49FBD5D035CB6",
	}

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "config-get", client.Services.Agent),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := ConfigGet[CtrlAgentConfig](mockClient)
	if err != nil {
		t.Fatalf("ConfigGet() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ConfigGet() = %+v, want %+v", got, want)
	}
}

// TestStatusGet tests the StatusGet function for the CtrlAgentStatus type.
func TestStatusGet(t *testing.T) {
	t.Parallel()

	want := CtrlAgentStatus{
		PID:    14013,
		Uptime: 123,
		Reload: 456,
	}

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "status-get", client.Services.Agent),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := StatusGet(mockClient)
	if err != nil {
		t.Fatalf("StatusGet() error = %v", err)
	}
	if got != want {
		t.Errorf("StatusGet() = %+v, want %+v", got, want)
	}
}

// TestListCommands tests the ListCommands function for the CtrlAgent type.
func TestListCommands(t *testing.T) {
	t.Parallel()

	want := []string{
		"build-report", "config-get", "config-hash-get", "config-reload", "config-set",
		"config-test", "config-write", "list-commands", "shutdown", "status-get", "version-get",
	}

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "list-commands", client.Services.Agent),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := ListCommands(mockClient)
	if err != nil {
		t.Fatalf("ListCommands() error = %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("ListCommands() len = %d, want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("ListCommands()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// TestVersionGet tests the VersionGet function for the CtrlAgent type.
func TestVersionGet(t *testing.T) {
	t.Parallel()

	wantText := "2.6.3"
	want := CtrlAgentVersion{
		Extended: "2.6.3 (isc20250522135511 deb)\npremium: yes (isc20250522135511 deb)\nlinked with:\n- log4cplus 2.0.8\n- OpenSSL 3.0.16 11 Feb 2025",
	}

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "version-get", client.Services.Agent),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Text:      wantText,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	gotText, gotVersion, err := VersionGet(mockClient)
	if err != nil {
		t.Fatalf("VersionGet() error = %v", err)
	}

	if gotText != wantText {
		t.Errorf("VersionGet() text = %q, want %q", gotText, wantText)
	}

	if gotVersion != want {
		t.Errorf("VersionGet() = %+v, want %+v", gotVersion, want)
	}
}
