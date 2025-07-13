package dhcp4

import (
	"github.com/rannday/kea-api/client"
)

/*
 * Commands supported by kea-dhcp4 daemon: build-report, cache-clear, cache-flush, cache-get, cache-get-by-id,
 * cache-insert, cache-load, cache-remove, cache-size, cache-write, class-add, class-del, class-get, class-list,
 * class-update, config-backend-pull, config-get, config-hash-get, config-reload, config-set, config-test,
 * config-write, dhcp-disable, dhcp-enable, extended-info4-upgrade, ha-continue, ha-heartbeat, ha-maintenance-cancel,
 * ha-maintenance-notify, ha-maintenance-start, ha-reset, ha-scopes, ha-sync, ha-sync-complete-notify, lease4-add,
 * lease4-del, lease4-get, lease4-get-all, lease4-get-by-client-id, lease4-get-by-hostname, lease4-get-by-hw-address,
 * lease4-get-page, lease4-resend-ddns, lease4-update, lease4-wipe, lease4-write, leases-reclaim, list-commands,
 * network4-add, network4-del, network4-get, network4-list, network4-subnet-add, network4-subnet-del, perfmon-control,
 * perfmon-get-all-durations, remote-class4-del, remote-class4-get, remote-class4-get-all, remote-class4-set,
 * remote-global-parameter4-del, remote-global-parameter4-get, remote-global-parameter4-get-all, remote-global-parameter4-set,
 * remote-network4-del, remote-network4-get, remote-network4-list, remote-network4-set, remote-option-def4-del,
 * remote-option-def4-get, remote-option-def4-get-all, remote-option-def4-set, remote-option4-global-del, remote-option4-global-get,
 * remote-option4-global-get-all, remote-option4-global-set, remote-option4-network-del, remote-option4-network-set,
 * remote-option4-pool-del, remote-option4-pool-set, remote-option4-subnet-del, remote-option4-subnet-set, remote-server4-del,
 * remote-server4-get, remote-server4-get-all, remote-server4-set, remote-subnet4-del-by-id, remote-subnet4-del-by-prefix,
 * remote-subnet4-get-by-id, remote-subnet4-get-by-prefix, remote-subnet4-list, remote-subnet4-set, reservation-add,
 * reservation-del, reservation-get, reservation-get-all, reservation-get-by-address, reservation-get-by-hostname,
 * reservation-get-by-id, reservation-get-page, reservation-update, server-tag-get, shutdown, stat-lease4-get, statistic-get,
 * statistic-get-all, statistic-remove, statistic-remove-all, statistic-reset, statistic-reset-all, statistic-sample-age-set,
 * statistic-sample-age-set-all, statistic-sample-count-set, statistic-sample-count-set-all, status-get, subnet4-add, subnet4-del,
 * subnet4-delta-add, subnet4-delta-del, subnet4-get, subnet4-list, subnet4-select-test, subnet4-update, subnet4o6-select-test, version-get.
 */

// BuildReport fetches the build configuration report of the DHCPv4 server.
func BuildReport(c *client.Client) (string, error) {
	return client.BuildReport(c, client.Services.DHCP4)
}

// ConfigGet fetches the DHCPv4 configuration using the generic helper.
func ConfigGet(c *client.Client) (Dhcp4Config, error) {
	return client.ConfigGet[Dhcp4Config](c, client.Services.DHCP4)
}

// StatusGet fetches the DHCPv4 server status using the generic helper.
func StatusGet(c *client.Client) (DHCP4Status, error) {
	return client.StatusGet[DHCP4Status](c, client.Services.DHCP4)
}

// ListCommands fetches the list of commands for DHCPv4.
func ListCommands(c *client.Client) ([]string, error) {
	return client.ListCommands(c, client.Services.DHCP4)
}

// VersionGet fetches the DHCPv4 server version using the shared helper.
func VersionGet(c *client.Client) (string, DHCP4Version, error) {
	return client.VersionGet[DHCP4Version](c, client.Services.DHCP4)
}
