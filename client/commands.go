package client

import (
	"encoding/json"
	"fmt"
)

// CallCommand sends a command request and returns the first successful response.
// It validates the response result and ensures the response is non-empty.
func CallCommand(c *Client, cmd string, service Service) (CommandResponse, error) {
	req := CommandRequest{
		Command: cmd,
	}
	if service != "" {
		req.Service = []Service{service}
	}

	var res []CommandResponse
	if err := c.Call(req, &res); err != nil {
		return CommandResponse{}, fmt.Errorf("%s failed: %w", cmd, err)
	}
	if len(res) == 0 {
		return CommandResponse{}, fmt.Errorf("%s returned empty response", cmd)
	}
	if res[0].Result != ResultSuccess {
		return CommandResponse{}, res[0].Result.ResultError(res[0].Text)
	}

	return res[0], nil
}

// CallAndDecode is a shared helper to send a command and unmarshal the result into T.
func CallAndDecode[T any](c *Client, cmd string, service Service) (string, T, error) {
	var zero T

	resp, err := CallCommand(c, cmd, service)
	if err != nil {
		return "", zero, err
	}

	if err := json.Unmarshal(resp.Arguments, &zero); err != nil {
		return resp.Text, zero, fmt.Errorf("decode %s arguments: %w", cmd, err)
	}

	return resp.Text, zero, nil
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

// BuildReport fetches compilation metadata from a given service.
func BuildReport(c *Client, service Service) (string, error) {
	resp, err := CallCommand(c, "build-report", service)
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

// ListCommands gets a list of commands for a given service.
func ListCommands(c *Client, service Service) ([]string, error) {
	_, cmds, err := CallAndDecode[[]string](c, "list-commands", service)
	return cmds, err
}

// StatusGet is a generic helper for unmarshaling a status-get command.
func StatusGet[T any](c *Client, service Service) (T, error) {
	_, val, err := CallAndDecode[T](c, "status-get", service)
	return val, err
}

// VersionGet is a generic helper for unmarshaling a version-get command.
func VersionGet[T any](c *Client, service Service) (string, T, error) {
	return CallAndDecode[T](c, "version-get", service)
}
