package agent

import "github.com/rannday/kea-api/client"

/*
 * Control Agent API
 * Supported by kea-ctrl-agent daemon:
 * build-report, config-get, config-hash-get, config-reload,
 * config-set, config-test, config-write, list-commands, shutdown,
 * status-get, version-get.
 */

// ListCommands fetches the list of commands for the control-agent.
func ListCommands(c *client.Client) ([]string, error) {
	return client.ListCommands(c, "")
}

func StatusGet(c *client.Client) (CtrlAgentStatus, error) {
	return client.StatusGet[CtrlAgentStatus](c, "")
}

func VersionGet(c *client.Client) (string, CtrlAgentVersion, error) {
	return client.VersionGet[CtrlAgentVersion](c, "")
}
