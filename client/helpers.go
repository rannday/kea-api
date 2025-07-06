package client

import (
	"encoding/json"
	"fmt"
)

// doCommand is a shared helper to send a command and unmarshal the result into T.
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
