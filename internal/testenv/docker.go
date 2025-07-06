package testenv

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// StartKeaContainer starts the Docker container for integration testing.
func StartKeaContainer() error {
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"--name", "kea-int-test",
		"-p", "8000:8000",
		"kea-custom:latest",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start kea container: %v\n%s", err, string(output))
	}
	return waitForKea()
}

// StopKeaContainer stops the Docker container.
func StopKeaContainer() {
	_ = exec.Command("docker", "stop", "kea-int-test").Run()
}

// waitForKea pings the Kea control agent until it's ready.
func waitForKea() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for Kea to become ready")
		default:
			resp := exec.Command("curl", "-s", "http://localhost:8000")
			if out, err := resp.Output(); err == nil && len(out) > 0 {
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}
