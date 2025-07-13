package types

// SocketConfig defines the control socket location and type.
type SocketConfig struct {
	SocketName string `json:"socket-name"`
	SocketType string `json:"socket-type"`
}

// LoggerConfig represents a logger instance for outputting diagnostic information.
type LoggerConfig struct {
	DebugLevel    int               `json:"debuglevel"`
	Name          string            `json:"name"`
	OutputOptions []LogOutputOption `json:"output-options"`
	Severity      string            `json:"severity"`
}

// LogOutputOption controls destination, format, and flushing for logger output.
type LogOutputOption struct {
	Flush   bool   `json:"flush"`
	Output  string `json:"output"`
	Pattern string `json:"pattern"`
}

// HookLibrary represents a dynamically loaded Kea hook module.
type HookLibrary struct {
	Library string `json:"library"`
}

// InterfacesConfig lists interfaces to bind and whether to auto-detect them.
type InterfacesConfig struct {
	Interfaces []string `json:"interfaces"`
	ReDetect   bool     `json:"re-detect"`
}

// SanityChecks configures validation levels for leases and configuration.
type SanityChecks struct {
	ExtendedInfoChecks string `json:"extended-info-checks"`
	LeaseChecks        string `json:"lease-checks"`
}

// DatabaseConfig represents a generic database connection config (used in config-databases or hosts-databases).
type DatabaseConfig struct {
	Host     string `json:"host"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Type     string `json:"type"` // e.g. "mysql", "postgresql"
	User     string `json:"user"`
}

// LeaseDatabaseConfig represents the lease-database block, which may include additional memfile-specific fields.
type LeaseDatabaseConfig struct {
	Host        string `json:"host,omitempty"`
	Name        string `json:"name,omitempty"`
	Password    string `json:"password,omitempty"`
	Port        int    `json:"port,omitempty"`
	Type        string `json:"type"` // e.g. "mysql", "memfile"
	User        string `json:"user,omitempty"`
	LFCInterval int    `json:"lfc-interval,omitempty"` // Only used by memfile backends
}

// MultiThreadingConfig controls concurrency behavior of the DHCP server.
type MultiThreadingConfig struct {
	EnableMultiThreading bool `json:"enable-multi-threading"`
	PacketQueueSize      int  `json:"packet-queue-size"`
	ThreadPoolSize       int  `json:"thread-pool-size"`
}

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
