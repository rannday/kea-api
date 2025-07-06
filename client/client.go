package client

import "time"

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

// NewSocket returns a Kea client using a UNIX or TCP socket.
func NewSocket(network, address string, timeout time.Duration) *Client {
	t := NewSocketTransport(network, address, timeout)
	return NewClient(t)
}

// Call sends the request using the chosen transport.
func (c *Client) Call(req CommandRequest, out interface{}) error {
	return c.transport.Call(req, out)
}
