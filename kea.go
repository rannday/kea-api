package kea

// Package kea provides a unified, high-level entrypoint for interacting with ISC Kea services.

import (
	"github.com/rannday/kea-api/agent"
	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/dhcp4"
	"github.com/rannday/kea-api/dhcp6"
)

// Shared Client type
type Client = client.Client

// Top-level constructors (unified)
var (
	NewHTTP   = client.NewHTTP
	NewSocket = client.NewSocket
)

// Control Agent API
var (
	StatusGet    = agent.StatusGet
	ListCommands = agent.ListCommands
)

// DHCPv4 API
var (
	StatusGet4    = dhcp4.StatusGet
	ListCommands4 = dhcp4.ListCommands
)

// DHCPv6 API
var (
	StatusGet6    = dhcp6.StatusGet
	ListCommands6 = dhcp6.ListCommands
)
