package dhcp6

import (
	"github.com/rannday/kea-api/client"
)

// BuildReport fetches the build configuration report of the DHCPv6 server.
func BuildReport(c *client.Client) (string, error) {
	return client.BuildReport(c, client.Services.DHCP6)
}

// ConfigGet fetches the config from the DHCPv6 service.
func ConfigGet(c *client.Client) (Dhcp6Config, error) {
	return client.ConfigGet[Dhcp6Config](c, "dhcp6")
}

// StatusGet fetches the DHCPv6 server status using the generic helper.
func StatusGet(c *client.Client) (DHCP6Status, error) {
	return client.StatusGet[DHCP6Status](c, client.Services.DHCP6)
}

// ListCommands fetches the list of commands for DHCPv6.
func ListCommands(c *client.Client) ([]string, error) {
	return client.ListCommands(c, client.Services.DHCP6)
}

// VersionGet fetches the DHCPv6 server version using the shared helper.
func VersionGet(c *client.Client) (string, DHCP6Version, error) {
	return client.VersionGet[DHCP6Version](c, client.Services.DHCP6)
}
