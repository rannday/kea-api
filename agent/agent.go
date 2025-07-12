package agent

import "github.com/rannday/kea-api/client"

/*
 * Control Agent API
 * Supported by kea-ctrl-agent daemon:
 * build-report, config-get, config-hash-get, config-reload,
 * config-set, config-test, config-write, list-commands, shutdown,
 * status-get, version-get.
 */

// BuildReport fetches the build configuration report of the control-agent.
func BuildReport(c *client.Client) (string, error) {
	return client.BuildReport(c, "")
}

// ListCommands fetches the list of commands for the control-agent.
func ListCommands(c *client.Client) ([]string, error) {
	return client.ListCommands(c, "")
}

// StatusGet fetches the status of the control-agent.
func StatusGet(c *client.Client) (CtrlAgentStatus, error) {
	return client.StatusGet[CtrlAgentStatus](c, "")
}

// VersionGet fetches the version of the control-agent.
func VersionGet(c *client.Client) (string, CtrlAgentVersion, error) {
	return client.VersionGet[CtrlAgentVersion](c, "")
}
