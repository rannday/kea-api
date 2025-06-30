package dhcp4

import (
	"github.com/rannday/isc-kea/client"
)

type Client = client.Client

// StatusGet fetches the DHCPv4 server status using the generic helper.
func StatusGet(c *client.Client) (DHCP4Status, error) {
	return client.StatusGet[DHCP4Status](c, client.ServiceDHCP4)
}

// ListCommands fetches the list of commands for DHCPv4.
func ListCommands(c *client.Client) ([]string, error) {
	return client.ListCommands(c, client.ServiceDHCP4)
}
