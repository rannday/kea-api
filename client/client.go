package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents a client for the ISC Kea Control Agent API over HTTP.
type Client struct {
	endpoint   string
	httpClient *http.Client
	auth       AuthProvider
}

// Option defines a functional option for configuring a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithAuth sets the client's AuthProvider (e.g. BasicAuth).
func WithAuth(auth AuthProvider) Option {
	return func(c *Client) {
		c.auth = auth
	}
}

// NewHTTPClient creates a new Client using the given endpoint and optional configuration.
func NewHTTPClient(endpoint string, opts ...Option) *Client {
	client := &Client{
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// Call sends a CommandRequest and decodes the JSON response into `out`.
// It expects `out` to be *[]CommandResponse and checks each response for result == 0.
func (c *Client) Call(req CommandRequest, out interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if c.auth != nil {
		if err := c.auth.Apply(httpReq); err != nil {
			return fmt.Errorf("apply auth: %w", err)
		}
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("kea error: status=%d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	// Check result codes if we got a []CommandResponse
	if responses, ok := out.(*[]CommandResponse); ok {
		if len(*responses) == 0 {
			return fmt.Errorf("empty response from Kea")
		}
		for _, r := range *responses {
			if err := r.Result.ResultError(r.Text); err != nil {
				return err
			}
		}
	}

	return nil
}
