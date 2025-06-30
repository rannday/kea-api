package client

import (
	"encoding/json"
	"fmt"
)

// ResultCode represents the result code returned by Kea commands.
type ResultCode int

// Kea result codes.
const (
	ResultSuccess        ResultCode = 0
	ResultGeneralFailure ResultCode = 1
	ResultUnsupported    ResultCode = 2
	ResultNotFound       ResultCode = 3
	ResultConflict       ResultCode = 4
)

// ResultError converts a ResultCode to an error.
func (r ResultCode) ResultError(text string) error {
	switch r {
	case ResultSuccess:
		return nil
	case ResultGeneralFailure:
		return fmt.Errorf("general error: %s", text)
	case ResultUnsupported:
		return fmt.Errorf("unsupported command: %s", text)
	case ResultNotFound:
		return fmt.Errorf("resource not found: %s", text)
	case ResultConflict:
		return fmt.Errorf("conflict: %s", text)
	default:
		return fmt.Errorf("unknown result code %d: %s", r, text)
	}
}

// Service represents a Kea service name.
type Service string

// Known services in Kea.
const (
	ServiceDHCP4 Service = "dhcp4"
	ServiceDHCP6 Service = "dhcp6"
)

// CommandRequest represents a Kea API command.
type CommandRequest struct {
	Command   string                 `json:"command"`
	Service   []Service              `json:"service,omitempty"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CommandResponse represents a response from Kea.
type CommandResponse struct {
	Result    ResultCode      `json:"result"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
	Text      string          `json:"text,omitempty"`
}

// AsMap attempts to decode Arguments into a map.
func (r *CommandResponse) AsMap() (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := json.Unmarshal(r.Arguments, &out); err != nil {
		return nil, fmt.Errorf("arguments is not a JSON object: %w", err)
	}
	return out, nil
}

// AsList attempts to unmarshal Arguments into a list of strings.
func (r *CommandResponse) AsList() ([]string, error) {
	var out []string
	if err := json.Unmarshal(r.Arguments, &out); err != nil {
		return nil, fmt.Errorf("arguments is not a list of strings: %w", err)
	}
	return out, nil
}
