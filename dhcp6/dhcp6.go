package dhcp6

import (
	"github.com/rannday/isc-kea/client"
)

// StatusGet fetches the DHCPv6 server status using the generic helper.
func StatusGet(c *client.Client) (DHCP6Status, error) {
	return client.StatusGet[DHCP6Status](c, client.ServiceDHCP6)
}

// ListCommands fetches the list of commands for DHCPv6.
func ListCommands(c *client.Client) ([]string, error) {
	return client.ListCommands(c, client.ServiceDHCP6)
}

// VersionGet fetches the DHCPv6 server version using the shared helper.
func VersionGet(c *client.Client) (string, DHCP6Version, error) {
	return client.VersionGet[DHCP6Version](c, client.ServiceDHCP6)
}
