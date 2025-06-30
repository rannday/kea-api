package kea

// Package kea provides a unified, high-level entrypoint for interacting with ISC Kea services

import (
	"github.com/rannday/isc-kea/agent"
	"github.com/rannday/isc-kea/client"
	"github.com/rannday/isc-kea/dhcp4"
	"github.com/rannday/isc-kea/dhcp6"
)

type Client = client.Client

var (
	NewHTTPClient  = client.NewHTTPClient
	WithHTTPClient = client.WithHTTPClient
	WithAuth       = client.WithAuth
)

// Control Agent API
var (
	StatusGet = agent.StatusGet
)

// DHCP4 API
var (
	StatusGet4 = dhcp4.StatusGet
)

// DHCP6 API
var (
	StatusGet6 = dhcp6.StatusGet
)
