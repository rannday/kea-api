package types

// DHCPState captures the service's enabled/disabled status.
type DHCPState struct {
	GloballyDisabled        bool     `json:"globally-disabled"`
	DisabledByUser          bool     `json:"disabled-by-user"`
	DisabledByRemoteCommand []string `json:"disabled-by-remote-command"`
	DisabledByLocalCommand  []string `json:"disabled-by-local-command"`
	DisabledByDBConnection  []string `json:"disabled-by-db-connection"`
}

// ListCommandsResponse is the typed structure for a `list-commands` reply.
type ListCommandsResponse struct {
	Arguments []string `json:"arguments"`
}
