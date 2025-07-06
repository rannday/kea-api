package dhcp6

import (
	"testing"

	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/internal/testenv"
	"github.com/rannday/kea-api/types"
)

// TestStatusGet tests the StatusGet function for the CtrlDHCP6 type.
func TestStatusGet(t *testing.T) {
	t.Parallel()

	want := DHCP6Status{
		PID:                   67890,
		Uptime:                200,
		Reload:                3,
		ThreadPoolSize:        8,
		MultiThreadingEnabled: false,
		PacketQueueSize:       32,
		PacketQueueStatistics: []float64{1.1, 1.2},
		Sockets:               map[string]interface{}{"eth1": "bound"},
		DHCPState:             types.DHCPState{DisabledByUser: true},
		ExtendedInfoTables:    true,
	}

	client := testenv.NewTestClient(t,
		testenv.ExpectCommand(t, "status-get", client.ServiceDHCP6),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := StatusGet(client)
	if err != nil {
		t.Fatalf("StatusGet() error = %v", err)
	}

	if got.PID != want.PID || got.Reload != want.Reload || !got.ExtendedInfoTables {
		t.Errorf("StatusGet() = %+v, want %+v", got, want)
	}
}

// TestListCommands tests the ListCommands function for the CtrlDHCP6 type.
func TestListCommands(t *testing.T) {
	t.Parallel()

	want := []string{
		"build-report", "config-backend-pull", "config-get", "config-hash-get", "config-reload",
		"config-set", "config-test", "config-write", "dhcp-disable", "dhcp-enable",
		"leases-reclaim", "list-commands", "server-tag-get", "shutdown", "statistic-get",
		"statistic-get-all", "statistic-remove", "statistic-remove-all", "statistic-reset",
		"statistic-reset-all", "statistic-sample-age-set", "statistic-sample-age-set-all",
		"statistic-sample-count-set", "statistic-sample-count-set-all", "status-get", "version-get",
	}

	client := testenv.NewTestClient(t,
		testenv.ExpectCommand(t, "list-commands", client.ServiceDHCP6),
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
