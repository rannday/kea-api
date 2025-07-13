package ddns

import "github.com/rannday/kea-api/client"

/*
 * Commands supported by kea-dhcp-ddns daemon: build-report, config-get, config-hash-get, config-reload, config-set, config-test,
 * config-write, gss-tsig-get, gss-tsig-get-all, gss-tsig-key-del, gss-tsig-key-expire, gss-tsig-key-get, gss-tsig-list, gss-tsig-purge,
 * gss-tsig-purge-all, gss-tsig-rekey, gss-tsig-rekey-all, list-commands, shutdown, statistic-get, statistic-get-all, statistic-reset,
 * statistic-reset-all, status-get, version-get.
 */

// ConfigGet retrieves the configuration from the Kea DDNS (d2) service.
func ConfigGet(c *client.Client) (DdnsConfig, error) {
	return client.ConfigGet[DdnsConfig](c, "d2")
}
