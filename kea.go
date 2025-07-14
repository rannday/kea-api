package kea

// Package kea provides a unified, high-level entrypoint for interacting with ISC Kea services.

import (
	"github.com/rannday/kea-api/agent"
	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/dhcp4"
	"github.com/rannday/kea-api/dhcp6"
)

type Client = client.Client

// Client
var (
	NewHTTP   = client.NewHTTP
	NewSocket = client.NewSocket
)

// Control Agent
var (
	BuildReport  = agent.BuildReport
	ConfigGet    = agent.ConfigGet
	ListCommands = agent.ListCommands
	StatusGet    = agent.StatusGet
	VersionGet   = agent.VersionGet
)

// DHCPv4
var (
	BuildReport4  = dhcp4.BuildReport
	ConfigGet4    = dhcp4.ConfigGet
	ListCommands4 = dhcp4.ListCommands
	StatusGet4    = dhcp4.StatusGet
	VersionGet4   = dhcp4.VersionGet
)

// DHCPv6
var (
	BuildReport6  = dhcp6.BuildReport
	ConfigGet6    = dhcp6.ConfigGet
	ListCommands6 = dhcp6.ListCommands
	StatusGet6    = dhcp6.StatusGet
	VersionGet6   = dhcp6.VersionGet
)

// DDNS
/*var (
	
)*/
