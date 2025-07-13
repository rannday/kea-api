package agent

import (
	"testing"

	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/internal/testenv"
)

// TestBuildReport tests the BuildReport function for the CtrlAgent type.
func TestBuildReport(t *testing.T) {
	t.Parallel()

	want := `Kea source configure results:
Package:
  Name: kea
  Version: 2.6.3
`

	client := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "build-report"),
		[]client.CommandResponse{{
			Result: client.ResultSuccess,
			Text:   want,
		}},
	)

	got, err := BuildReport(client)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}

	if got != want {
		t.Errorf("BuildReport() = %q, want %q", got, want)
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

	client := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "status-get"),
		[]client.CommandResponse{{
			Result:    client.ResultSuccess,
			Arguments: testenv.MustEncodeRawJSON(t, want),
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

	client := testenv.NewMockClient(t,
		testenv.ExpectCommand(t, "list-commands"),
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
