package client

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Transport defines how to send a CommandRequest and receive a response.
type Transport interface {
	Call(req CommandRequest, out interface{}) error
}

// Client routes Kea API calls to the underlying Transport.
type Client struct {
	transport Transport
}

// NewClient returns a new Kea client using the provided transport.
func NewClient(t Transport) *Client {
	return &Client{transport: t}
}

// NewHTTP returns a Kea client using HTTP transport.
func NewHTTP(endpoint string, opts ...HTTPOption) *Client {
	t := NewHTTPTransport(endpoint, opts...)
	return NewClient(t)
}

// timeoutOrDefault returns the first timeout provided or a default of 5 seconds.
// It returns an error if the timeout is zero or negative.
// If the timeout is unusually long (>60s), a warning is printed to stderr.
func timeoutOrDefault(t []time.Duration) (time.Duration, error) {
	const maxRecommended = 60 * time.Second
	const defaultTimeout = 5 * time.Second
	if len(t) == 0 {
		return defaultTimeout, nil
	}
	timeout := t[0]
	if timeout <= 0 {
		return 0, errors.New("timeout must be greater than zero")
	}
	if timeout > maxRecommended {
		fmt.Fprintln(os.Stderr, "kea-api [warning] unusually long socket timeout:", timeout)
	}
	return timeout, nil
}

// NewSocket returns a Kea client using a UNIX or TCP socket.
func NewSocket(network, address string, timeout ...time.Duration) (*Client, error) {
	tout, err := timeoutOrDefault(timeout)
	if err != nil {
		return nil, err
	}
	t := NewSocketTransport(network, address, tout)
	return NewClient(t), nil
}

// Call sends the request using the chosen transport.
func (c *Client) Call(req CommandRequest, out interface{}) error {
	return c.transport.Call(req, out)
}
