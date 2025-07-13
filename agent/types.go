package agent

import "github.com/rannday/kea-api/types"

// CtrlAgentStatus represents the response from "status-get" when no service is specified.
type CtrlAgentStatus struct {
	PID    int `json:"pid"`    // Process ID of the control-agent
	Uptime int `json:"uptime"` // Uptime in seconds
	Reload int `json:"reload"` // Number of times configuration has been reloaded
}

// CtrlAgentVersion is the response from "version-get" on the control-agent.
type CtrlAgentVersion struct {
	Extended string `json:"extended"` // Full version string with build metadata
}

// CtrlAgentConfig is the response from "config-get" on the control-agent.
// It contains configuration parameters and control-socket definitions.
type CtrlAgentConfig struct {
	ControlAgent CtrlAgentBlock `json:"Control-agent"` // Main configuration block
	Hash         string         `json:"hash"`          // SHA256 hash of the current config
}

// CtrlAgentBlock describes the Control-agent settings block.
type CtrlAgentBlock struct {
	Authentication AuthConfig                    `json:"authentication"`  // Authentication settings
	ControlSockets map[string]types.SocketConfig `json:"control-sockets"` // Per-service socket paths
	HooksLibraries []types.HookLibrary           `json:"hooks-libraries"` // Installed hooks (usually empty)
	HTTPHost       string                        `json:"http-host"`       // HTTP bind address
	HTTPPort       int                           `json:"http-port"`       // HTTP bind port
	Loggers        []types.LoggerConfig          `json:"loggers"`         // Logging configuration
}

// AuthConfig defines HTTP Basic Authentication parameters.
type AuthConfig struct {
	Clients   []AuthClient `json:"clients"`   // List of valid users
	Directory string       `json:"directory"` // Directory where password file is stored
	Realm     string       `json:"realm"`     // HTTP auth realm name
	Type      string       `json:"type"`      // Authentication type (e.g. "basic")
}

// AuthClient represents a single user credential entry.
type AuthClient struct {
	PasswordFile string `json:"password-file"` // Filename of password hash file
	User         string `json:"user"`          // Username
}
