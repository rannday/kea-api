package dhcp4

import (
	"testing"

	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/internal/testenv"
	"github.com/rannday/kea-api/types"
)

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

	client := testenv.NewTestClient(t,
		testenv.ExpectCommand(t, "status-get", client.ServiceDHCP4),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := StatusGet(client)
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

	client := testenv.NewTestClient(t,
		testenv.ExpectCommand(t, "list-commands", client.ServiceDHCP4),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := ListCommands(client)
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
