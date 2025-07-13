package ddns

import "github.com/rannday/kea-api/types"

// DdnsConfig is the typed response from config-get on the d2 (DDNS) service.
type DdnsConfig struct {
	DhcpDdns DhcpDdnsBlock `json:"DhcpDdns"` // DDNS service configuration
	Hash     string        `json:"hash"`     // Configuration hash
}

// DhcpDdnsBlock contains all Kea DDNS service configuration options.
type DhcpDdnsBlock struct {
	ControlSocket    types.SocketConfig   `json:"control-socket"`     // Control channel socket
	DNSServerTimeout int                  `json:"dns-server-timeout"` // Timeout in ms for communicating with DNS servers
	ForwardDDNS      DdnsDirectionConfig  `json:"forward-ddns"`       // Forward DDNS configuration
	ReverseDDNS      DdnsDirectionConfig  `json:"reverse-ddns"`       // Reverse DDNS configuration
	HooksLibraries   []types.HookLibrary  `json:"hooks-libraries"`    // Loaded DDNS hook libraries
	IPAddress        string               `json:"ip-address"`         // IP the DDNS server listens on
	Loggers          []types.LoggerConfig `json:"loggers"`            // Logging configuration
	NCRFormat        string               `json:"ncr-format"`         // NCR format (e.g. "JSON")
	NCRProtocol      string               `json:"ncr-protocol"`       // Protocol for NCRs (e.g. "UDP")
	Port             int                  `json:"port"`               // Listening port
	TSIGKeys         []TsigKey            `json:"tsig-keys"`          // List of TSIG key configurations
}

// DdnsDirectionConfig holds DDNS configuration for forward or reverse zones.
type DdnsDirectionConfig struct {
	Domains []interface{} `json:"ddns-domains"` // List of domain blocks (optional)
}

// TsigKey defines a TSIG key used for secure DNS updates.
type TsigKey struct {
	Algorithm  string `json:"algorithm"`             // TSIG algorithm (e.g. "hmac-sha256")
	DigestBits int    `json:"digest-bits,omitempty"` // Optional: number of digest bits
	KeyName    string `json:"key-name"`              // Identifier for the key
	Secret     string `json:"secret"`                // Base64-encoded secret
}
