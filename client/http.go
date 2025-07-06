package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPTransport sends commands to Kea over HTTP.
type HTTPTransport struct {
	endpoint   string
	httpClient *http.Client
	auth       AuthProvider
}

// HTTPOption configures an HTTPTransport.
type HTTPOption func(*HTTPTransport)

func WithHTTPClient(hc *http.Client) HTTPOption {
	return func(t *HTTPTransport) {
		t.httpClient = hc
	}
}

func WithAuth(auth AuthProvider) HTTPOption {
	return func(t *HTTPTransport) {
		t.auth = auth
	}
}

// NewHTTPTransport returns a configured HTTPTransport.
func NewHTTPTransport(endpoint string, opts ...HTTPOption) *HTTPTransport {
	t := &HTTPTransport{
		endpoint:   endpoint,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Call implements the Transport interface for HTTP.
func (t *HTTPTransport) Call(req CommandRequest, out interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", t.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	if t.auth != nil {
		if err := t.auth.Apply(httpReq); err != nil {
			return fmt.Errorf("apply auth: %w", err)
		}
	}

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("kea error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

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
