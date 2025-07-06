package client

/* Shared Kea API commands across ctrl-agent, dhcp4, dhcp6, and ddns
build-report               config-get                 config-hash-get
config-reload              config-set                 config-test
config-write               list-commands              shutdown
status-get                 version-get

// Shared by dhcp4, dhcp6, and ddns (not ctrl-agent)
statistic-get              statistic-get-all          statistic-reset
statistic-reset-all
*/

// ListCommands gets a list of commands for a given service.
func ListCommands(c *Client, service Service) ([]string, error) {
	_, cmds, err := DoCommand[[]string](c, "list-commands", service)
	return cmds, err
}

// StatusGet is a generic helper for unmarshaling a status-get command.
func StatusGet[T any](c *Client, service Service) (T, error) {
	_, val, err := DoCommand[T](c, "status-get", service)
	return val, err
}

// VersionGet is a generic helper for unmarshaling a version-get command.
func VersionGet[T any](c *Client, service Service) (string, T, error) {
	return DoCommand[T](c, "version-get", service)
}
