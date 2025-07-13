package dhcp4

import "github.com/rannday/kea-api/types"

// DHCP4Status represents the response from "status-get" on kea-dhcp4.
type DHCP4Status struct {
	PID                   int                    `json:"pid"`
	Uptime                int                    `json:"uptime"`
	Reload                int                    `json:"reload"`
	ThreadPoolSize        int                    `json:"thread-pool-size"`
	MultiThreadingEnabled bool                   `json:"multi-threading-enabled"`
	PacketQueueSize       int                    `json:"packet-queue-size"`
	PacketQueueStatistics []float64              `json:"packet-queue-statistics"`
	Sockets               map[string]interface{} `json:"sockets"`
	DHCPState             types.DHCPState        `json:"dhcp-state"`
}

// DHCP4Version is the response from "version-get" on kea-dhcp4.
type DHCP4Version struct {
	Extended string `json:"extended"`
}

// Dhcp4Config is the typed response from "config-get" on the dhcp4 service.
type Dhcp4Config struct {
	Dhcp4 Dhcp4Block `json:"Dhcp4"`
	Hash  string     `json:"hash"`
}

// Dhcp4Block holds the core DHCPv4 configuration parameters.
type Dhcp4Block struct {
	Allocator                  string                     `json:"allocator"`
	Authoritative              bool                       `json:"authoritative"`
	BootFileName               string                     `json:"boot-file-name"`
	CalculateTeeTimes          bool                       `json:"calculate-tee-times"`
	ConfigControl              ConfigControl              `json:"config-control"`
	ControlSocket              types.SocketConfig         `json:"control-socket"`
	DDNSConflictMode           string                     `json:"ddns-conflict-resolution-mode"`
	DDNSGeneratedPrefix        string                     `json:"ddns-generated-prefix"`
	DDNSOverrideClientUpdate   bool                       `json:"ddns-override-client-update"`
	DDNSOverrideNoUpdate       bool                       `json:"ddns-override-no-update"`
	DDNSQualifyingSuffix       string                     `json:"ddns-qualifying-suffix"`
	DDNSReplaceClientName      string                     `json:"ddns-replace-client-name"`
	DDNSSendUpdates            bool                       `json:"ddns-send-updates"`
	DDNSUpdateOnRenew          bool                       `json:"ddns-update-on-renew"`
	DeclineProbationPeriod     int                        `json:"decline-probation-period"`
	DhcpDDNS                   DhcpDDNSConfig             `json:"dhcp-ddns"`
	DhcpQueueControl           DhcpQueueControl           `json:"dhcp-queue-control"`
	Dhcp4o6Port                int                        `json:"dhcp4o6-port"`
	EchoClientID               bool                       `json:"echo-client-id"`
	EarlyGlobalResLookup       bool                       `json:"early-global-reservations-lookup"`
	ExpiredLeasesProcessing    ExpiredLeasesProcessing    `json:"expired-leases-processing"`
	HooksLibraries             []types.HookLibrary        `json:"hooks-libraries"`
	HostReservationIdentifiers []string                   `json:"host-reservation-identifiers"`
	HostnameCharReplacement    string                     `json:"hostname-char-replacement"`
	HostnameCharSet            string                     `json:"hostname-char-set"`
	HostsDatabases             []types.DatabaseConfig     `json:"hosts-databases"`
	InterfacesConfig           types.InterfacesConfig     `json:"interfaces-config"`
	IPReservationsUnique       bool                       `json:"ip-reservations-unique"`
	LeaseDatabase              types.DatabaseConfig       `json:"lease-database"`
	Loggers                    []types.LoggerConfig       `json:"loggers"`
	MatchClientID              bool                       `json:"match-client-id"`
	MultiThreading             types.MultiThreadingConfig `json:"multi-threading"`
	NextServer                 string                     `json:"next-server"`
	OptionData                 []interface{}              `json:"option-data"`
	OptionDef                  []interface{}              `json:"option-def"`
	ParkedPacketLimit          int                        `json:"parked-packet-limit"`
	ReservationsGlobal         bool                       `json:"reservations-global"`
	ReservationsInSubnet       bool                       `json:"reservations-in-subnet"`
	ReservationsLookupFirst    bool                       `json:"reservations-lookup-first"`
	ReservationsOutOfPool      bool                       `json:"reservations-out-of-pool"`
	SanityChecks               types.SanityChecks         `json:"sanity-checks"`
	ServerHostname             string                     `json:"server-hostname"`
	ServerTag                  string                     `json:"server-tag"`
	SharedNetworks             []interface{}              `json:"shared-networks"`
	StashAgentOptions          bool                       `json:"stash-agent-options"`
	StatisticSampleAge         int                        `json:"statistic-default-sample-age"`
	StatisticSampleCount       int                        `json:"statistic-default-sample-count"`
	StoreExtendedInfo          bool                       `json:"store-extended-info"`
	Subnet4                    []interface{}              `json:"subnet4"`
	T1Percent                  float64                    `json:"t1-percent"`
	T2Percent                  float64                    `json:"t2-percent"`
	ValidLifetime              int                        `json:"valid-lifetime"`
}

// ConfigControl defines external config database references.
type ConfigControl struct {
	ConfigDatabases []types.DatabaseConfig `json:"config-databases"`
}

// DhcpDDNSConfig specifies how to send DDNS updates.
type DhcpDDNSConfig struct {
	EnableUpdates bool   `json:"enable-updates"`
	MaxQueueSize  int    `json:"max-queue-size"`
	NCRFormat     string `json:"ncr-format"`
	NCRProtocol   string `json:"ncr-protocol"`
	SenderIP      string `json:"sender-ip"`
	SenderPort    int    `json:"sender-port"`
	ServerIP      string `json:"server-ip"`
	ServerPort    int    `json:"server-port"`
}

// DhcpQueueControl governs packet queueing for async processing.
type DhcpQueueControl struct {
	Capacity    int    `json:"capacity"`
	EnableQueue bool   `json:"enable-queue"`
	QueueType   string `json:"queue-type"`
}

// ExpiredLeasesProcessing configures lease reclamation behavior.
type ExpiredLeasesProcessing struct {
	FlushReclaimedTimerWaitTime int `json:"flush-reclaimed-timer-wait-time"`
	HoldReclaimedTime           int `json:"hold-reclaimed-time"`
	MaxReclaimLeases            int `json:"max-reclaim-leases"`
	MaxReclaimTime              int `json:"max-reclaim-time"`
	ReclaimTimerWaitTime        int `json:"reclaim-timer-wait-time"`
	UnwarnedReclaimCycles       int `json:"unwarned-reclaim-cycles"`
}
