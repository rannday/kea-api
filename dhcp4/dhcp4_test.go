package dhcp4

import (
	"reflect"
	"testing"

	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/internal/testenv"
	"github.com/rannday/kea-api/types"
)

// TestBuildReport tests the BuildReport function for the DHCPv4 service.
func TestBuildReport(t *testing.T) {
	t.Parallel()

	want := `Kea source configure results:
Package:
  Name: kea-dhcp4
  Version: 2.6.3
`

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "build-report", client.Services.DHCP4),
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

// TestConfigGet tests the ConfigGet function for the Dhcp4Config type.
func TestConfigGet(t *testing.T) {
	t.Parallel()

	want := Dhcp4Config{
		Dhcp4: Dhcp4Block{
			Allocator:            "iterative",
			Authoritative:        true,
			BootFileName:         "pxelinux.0",
			ControlSocket:        types.SocketConfig{SocketName: "/run/kea/kea4-ctrl.sock", SocketType: "unix"},
			EchoClientID:         true,
			IPReservationsUnique: true,
			LeaseDatabase:        types.DatabaseConfig{Type: "memfile"},
			Loggers: []types.LoggerConfig{{
				Name:       "kea-dhcp4",
				Severity:   "INFO",
				DebugLevel: 0,
				OutputOptions: []types.LogOutputOption{{
					Output:  "stdout",
					Pattern: "%-5p %m\n",
					Flush:   true,
				}},
			}},
			InterfacesConfig: types.InterfacesConfig{
				Interfaces: []string{"eth0"},
				ReDetect:   false,
			},
			ServerTag:      "default",
			Subnet4:        []interface{}{},
			OptionData:     []interface{}{},
			OptionDef:      []interface{}{},
			SharedNetworks: []interface{}{},
			HostsDatabases: []types.DatabaseConfig{},
		},
		Hash: "abcdef1234567890",
	}

	client := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "config-get", client.Services.DHCP4),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := ConfigGet(client)
	if err != nil {
		t.Fatalf("ConfigGet() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ConfigGet() = %+v, want %+v", got, want)
	}
}

// TestStatusGet tests the StatusGet function for the CtrlDHCP4 type.
func TestStatusGet(t *testing.T) {
	t.Parallel()

	want := DHCP4Status{
		PID:                   12345,
		Uptime:                100,
		Reload:                2,
		ThreadPoolSize:        4,
		MultiThreadingEnabled: true,
		PacketQueueSize:       16,
		PacketQueueStatistics: []float64{0.1, 0.2, 0.3},
		Sockets:               map[string]interface{}{"eth0": "listening"},
		DHCPState:             types.DHCPState{GloballyDisabled: false},
	}

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "status-get", client.Services.DHCP4),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := StatusGet(mockClient)
	if err != nil {
		t.Fatalf("StatusGet() error = %v", err)
	}

	if got.PID != want.PID || got.Reload != want.Reload || got.Uptime != want.Uptime {
		t.Errorf("StatusGet() = %+v, want %+v", got, want)
	}
}

// TestListCommands tests the ListCommands function for the CtrlDHCP4 type.
func TestListCommands(t *testing.T) {
	t.Parallel()

	want := []string{
		"build-report", "config-backend-pull", "config-get", "config-hash-get", "config-reload",
		"config-set", "config-test", "config-write", "dhcp-disable", "dhcp-enable",
		"lease4-add", "lease4-del", "lease4-get", "lease4-get-all", "lease4-get-by-client-id",
		"lease4-get-by-hostname", "lease4-get-by-hw-address", "lease4-get-page",
		"lease4-resend-ddns", "lease4-update", "lease4-wipe", "lease4-write", "lease6-add",
		"lease6-bulk-apply", "lease6-del", "lease6-get", "lease6-get-all", "lease6-get-by-duid",
		"lease6-get-by-hostname", "lease6-get-page", "lease6-resend-ddns", "lease6-update",
		"lease6-wipe", "lease6-write", "leases-reclaim", "list-commands", "server-tag-get",
		"shutdown", "stat-lease4-get", "stat-lease6-get", "statistic-get", "statistic-get-all",
		"statistic-remove", "statistic-remove-all", "statistic-reset", "statistic-reset-all",
		"statistic-sample-age-set", "statistic-sample-age-set-all", "statistic-sample-count-set",
		"statistic-sample-count-set-all", "status-get", "version-get",
	}

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "list-commands", client.Services.DHCP4),
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

// TestVersionGet checks the VersionGet function for DHCP4 returns both text and typed extended version info.
func TestVersionGet(t *testing.T) {
	t.Parallel()

	wantText := "2.6.3"
	wantArgs := DHCP4Version{
		Extended: `2.6.3 (isc20250522135511 deb)
premium: yes (isc20250522135511 deb)
linked with:
- log4cplus 2.0.8
- OpenSSL 3.0.16 11 Feb 2025
backends:
- MySQL backend 22.2, library 3.3.14
- PostgreSQL backend 22.2, library 150013
- Memfile backend 3.0`,
	}

	mockClient := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "version-get", client.Services.DHCP4),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Text:      wantText,
			Arguments: testenv.MustEncodeRawJSON(t, wantArgs),
		}},
	)

	gotText, gotArgs, err := VersionGet(mockClient)
	if err != nil {
		t.Fatalf("VersionGet() error = %v", err)
	}

	if gotText != wantText {
		t.Errorf("VersionGet() text = %q, want %q", gotText, wantText)
	}
	if gotArgs.Extended != wantArgs.Extended {
		t.Errorf("VersionGet() extended = %q, want %q", gotArgs.Extended, wantArgs.Extended)
	}
}
