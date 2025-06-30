package client

import (
	"encoding/json"
	"fmt"
)

// StatusGet is a generic helper for unmarshaling a status-get command.
func StatusGet[T any](c *Client, service Service) (T, error) {
	req := CommandRequest{
		Command: "status-get",
	}
	if service != "" {
		req.Service = []Service{service}
	}

	var res []CommandResponse
	var zero T
	if err := c.Call(req, &res); err != nil {
		return zero, fmt.Errorf("status-get failed: %w", err)
	}
	if len(res) == 0 {
		return zero, fmt.Errorf("empty response")
	}

	var status T
	if err := json.Unmarshal(res[0].Arguments, &status); err != nil {
		return zero, fmt.Errorf("decode arguments: %w", err)
	}
	return status, nil
}

// ListCommands gets a list of commands for a given service.
func ListCommands(c *Client, service Service) ([]string, error) {
	req := CommandRequest{
		Command: "list-commands",
	}
	if service != "" {
		req.Service = []Service{service}
	}

	var res []CommandResponse
	if err := c.Call(req, &res); err != nil {
		return nil, fmt.Errorf("list-commands failed: %w", err)
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	var commands []string
	if err := json.Unmarshal(res[0].Arguments, &commands); err != nil {
		return nil, fmt.Errorf("decode arguments: %w", err)
	}
	return commands, nil
}
