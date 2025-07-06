package client

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// SocketTransport connects to Kea directly via UNIX or TCP sockets.
type SocketTransport struct {
	network string // "unix" or "tcp"
	address string
	timeout time.Duration
}

// NewSocketTransport creates a socket-based Transport.
func NewSocketTransport(network, address string, timeout time.Duration) *SocketTransport {
	return &SocketTransport{
		network: network,
		address: address,
		timeout: timeout,
	}
}

// Call implements the Transport interface for sockets.
func (s *SocketTransport) Call(req CommandRequest, out interface{}) error {
	conn, err := net.DialTimeout(s.network, s.address, s.timeout)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	// Optional: set read/write deadlines for additional protection
	_ = conn.SetDeadline(time.Now().Add(s.timeout))

	enc := json.NewEncoder(conn)
	if err := enc.Encode(&req); err != nil {
		return fmt.Errorf("encode request: %w", err)
	}

	dec := json.NewDecoder(conn)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
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
