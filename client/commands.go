package client

import (
	"encoding/json"
	"fmt"
)

// DoCommand is a shared helper to send a command and unmarshal the result into T.
func DoCommand[T any](c *Client, cmd string, service Service) (string, T, error) {
	req := CommandRequest{
		Command: cmd,
	}
	if service != "" {
		req.Service = []Service{service}
	}

	var res []CommandResponse
	var zero T

	if err := c.Call(req, &res); err != nil {
		return "", zero, fmt.Errorf("%s failed: %w", cmd, err)
	}
	if len(res) == 0 {
		return "", zero, fmt.Errorf("%s returned empty response", cmd)
	}
	if res[0].Result != ResultSuccess {
		return "", zero, res[0].Result.ResultError(res[0].Text)
	}

	if err := json.Unmarshal(res[0].Arguments, &zero); err != nil {
		return res[0].Text, zero, fmt.Errorf("decode %s arguments: %w", cmd, err)
	}

	return res[0].Text, zero, nil
}

/* Shared Kea API commands across ctrl-agent, dhcp4, dhcp6, and ddns
build-report               config-get                 config-hash-get
config-reload              config-set                 config-test
config-write               list-commands              shutdown
status-get                 version-get

// Shared by dhcp4, dhcp6, and ddns (not ctrl-agent)
statistic-get              statistic-get-all          statistic-reset
statistic-reset-all
*/

// ListCommands gets a list of commands for a given service.
func ListCommands(c *Client, service Service) ([]string, error) {
	_, cmds, err := DoCommand[[]string](c, "list-commands", service)
	return cmds, err
}

// StatusGet is a generic helper for unmarshaling a status-get command.
func StatusGet[T any](c *Client, service Service) (T, error) {
	_, val, err := DoCommand[T](c, "status-get", service)
	return val, err
}

// VersionGet is a generic helper for unmarshaling a version-get command.
func VersionGet[T any](c *Client, service Service) (string, T, error) {
	return DoCommand[T](c, "version-get", service)
}
