package client

import (
	"encoding/json"
	"fmt"
)

// CallCommand sends a command to one or more services and returns validated responses.
// If only Services.Agent is specified or no services are passed, the request is sent to the control-agent itself.
func CallCommand(c *Client, cmd string, services ...Service) ([]CommandResponse, error) {
	req := CommandRequest{Command: cmd}

	// Only set the "service" field if not targeting the control agent directly
	if len(services) > 0 && !(len(services) == 1 && services[0] == Services.Agent) {
		req.Service = services
	}

	var res []CommandResponse
	if err := c.Call(req, &res); err != nil {
		return nil, fmt.Errorf("%s failed: %w", cmd, err)
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("%s returned empty response", cmd)
	}

	for _, r := range res {
		if r.Result != ResultSuccess {
			return nil, r.Result.ResultError(r.Text)
		}
	}

	return res, nil
}

// CallAndExtractText sends a command and returns only the .Text from the first successful response.
func CallAndExtractText(c *Client, cmd string, service Service) (string, error) {
  responses, err := CallCommand(c, cmd, service)
  if err != nil {
    return "", err
  }

  return responses[0].Text, nil
}

// CallAndDecode sends a command and decodes all successful responses into a slice of T.
func CallAndDecode[T any](c *Client, cmd string, services ...Service) ([]T, error) {
	responses, err := CallCommand(c, cmd, services...)
	if err != nil {
		return nil, err
	}

	var out []T
	for _, r := range responses {
		var decoded T
		if err := json.Unmarshal(r.Arguments, &decoded); err != nil {
			return nil, fmt.Errorf("decode %s arguments: %w", cmd, err)
		}
		out = append(out, decoded)
	}

	return out, nil
}

// DecodeFirst sends a command and decodes the first response into T.
func DecodeFirst[T any](c *Client, cmd string, service Service) (T, error) {
	vals, err := CallAndDecode[T](c, cmd, service)
	if err != nil {
		var zero T
		return zero, err
	}
	return vals[0], nil
}

// DecodeFirstWithText sends a command and decodes the first response into T, also returning the text field.
func DecodeFirstWithText[T any](c *Client, cmd string, service Service) (string, T, error) {
	responses, err := CallCommand(c, cmd, service)
	if err != nil {
		var zero T
		return "", zero, err
	}

	text := responses[0].Text
	var decoded T
	if err := json.Unmarshal(responses[0].Arguments, &decoded); err != nil {
		return text, decoded, fmt.Errorf("%s: decode arguments: %w", cmd, err)
	}
	return text, decoded, nil
}
