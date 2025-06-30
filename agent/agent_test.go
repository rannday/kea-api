package agent

import (
	"testing"

	"github.com/rannday/isc-kea/client"
	"github.com/rannday/isc-kea/utils"
)

// TestStatusGet tests the StatusGet function for the CtrlAgentStatus type.
func TestStatusGet(t *testing.T) {
	t.Parallel()

	want := CtrlAgentStatus{
		PID:    14013,
		Uptime: 123,
		Reload: 456,
	}

	client := utils.NewTestClient(t,
		utils.ExpectCommand(t, "status-get"),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: utils.MustEncodeRawJSON(t, want),
		}},
	)

	got, err := StatusGet(client)
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

	client := utils.NewTestClient(t,
		utils.ExpectCommand(t, "list-commands"),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: utils.MustEncodeRawJSON(t, want),
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
