package agent

import (
	"net/http"

	"github.com/rannday/isc-kea/client"
)

type (
	Client       = client.Client
	AuthProvider = client.AuthProvider
	BasicAuth    = client.BasicAuth
	Option       = client.Option
)

// WithHTTPClient returns an Option that sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return client.WithHTTPClient(hc)
}

// WithAuth returns an Option that sets the client's AuthProvider.
func WithAuth(auth *BasicAuth) Option {
	return client.WithAuth(auth)
}

// NewHTTPClient creates a new Control Agent client using the provided endpoint and optional options.
func NewHTTPClient(endpoint string, opts ...Option) *client.Client {
	copts := make([]client.Option, len(opts))
	copy(copts, opts)
	return client.NewHTTPClient(endpoint, copts...)
}

// StatusGet fetches the control-agent status using the generic helper.
func StatusGet(c *client.Client) (CtrlAgentStatus, error) {
	return client.StatusGet[CtrlAgentStatus](c, "")
}

// ListCommands fetches the list of commands for the control-agent.
func ListCommands(c *client.Client) ([]string, error) {
	return client.ListCommands(c, "")
}
