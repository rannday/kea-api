package dhcp6

import "github.com/rannday/kea-api/types"

// DHCP6Status represents the result of status-get for dhcp6.
type DHCP6Status struct {
	PID                   int                    `json:"pid"`
	Uptime                int                    `json:"uptime"`
	Reload                int                    `json:"reload"`
	ThreadPoolSize        int                    `json:"thread-pool-size"`
	MultiThreadingEnabled bool                   `json:"multi-threading-enabled"`
	PacketQueueSize       int                    `json:"packet-queue-size"`
	PacketQueueStatistics []float64              `json:"packet-queue-statistics"`
	Sockets               map[string]interface{} `json:"sockets"`
	DHCPState             types.DHCPState        `json:"dhcp-state"`
	ExtendedInfoTables    bool                   `json:"extended-info-tables"`
}

// DHCP6Version is the response type for version-get on kea-dhcp6.
type DHCP6Version struct {
	Extended string `json:"extended"`
}

// Dhcp6Config is the typed response from config-get on the dhcp6 service.
type Dhcp6Config struct {
	Dhcp6 Dhcp6Block `json:"Dhcp6"`
	Hash  string     `json:"hash"`
}

// Dhcp6Block contains all DHCPv6 configuration parameters as returned by config-get.
type Dhcp6Block struct {
	Allocator                  string                     `json:"allocator"`
	CalculateTeeTimes          bool                       `json:"calculate-tee-times"`
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
	EarlyGlobalResLookup       bool                       `json:"early-global-reservations-lookup"`
	ExpiredLeasesProcessing    ExpiredLeasesProcessing    `json:"expired-leases-processing"`
	HooksLibraries             []types.HookLibrary        `json:"hooks-libraries"`
	HostReservationIdentifiers []string                   `json:"host-reservation-identifiers"`
	HostnameCharReplacement    string                     `json:"hostname-char-replacement"`
	HostnameCharSet            string                     `json:"hostname-char-set"`
	HostsDatabases             []types.DatabaseConfig     `json:"hosts-databases"` // ← added
	InterfacesConfig           types.InterfacesConfig     `json:"interfaces-config"`
	IPReservationsUnique       bool                       `json:"ip-reservations-unique"`
	LeaseDatabase              types.LeaseDatabaseConfig  `json:"lease-database"`
	Loggers                    []types.LoggerConfig       `json:"loggers"`
	MacSources                 []string                   `json:"mac-sources"`
	MultiThreading             types.MultiThreadingConfig `json:"multi-threading"`
	OptionData                 []interface{}              `json:"option-data"`
	OptionDef                  []interface{}              `json:"option-def"`
	ParkedPacketLimit          int                        `json:"parked-packet-limit"`
	PDAllocator                string                     `json:"pd-allocator"`
	PreferredLifetime          int                        `json:"preferred-lifetime"`
	RebindTimer                int                        `json:"rebind-timer"`
	RelaySuppliedOptions       []string                   `json:"relay-supplied-options"`
	RenewTimer                 int                        `json:"renew-timer"`
	ReservationsGlobal         bool                       `json:"reservations-global"`
	ReservationsInSubnet       bool                       `json:"reservations-in-subnet"`
	ReservationsLookupFirst    bool                       `json:"reservations-lookup-first"`
	ReservationsOutOfPool      bool                       `json:"reservations-out-of-pool"`
	SanityChecks               types.SanityChecks         `json:"sanity-checks"`
	ServerID                   ServerID                   `json:"server-id"`
	ServerTag                  string                     `json:"server-tag"`
	SharedNetworks             []interface{}              `json:"shared-networks"`
	StatisticSampleAge         int                        `json:"statistic-default-sample-age"`
	StatisticSampleCount       int                        `json:"statistic-default-sample-count"`
	StoreExtendedInfo          bool                       `json:"store-extended-info"`
	Subnet6                    []interface{}              `json:"subnet6"`
	T1Percent                  float64                    `json:"t1-percent"`
	T2Percent                  float64                    `json:"t2-percent"`
	ValidLifetime              int                        `json:"valid-lifetime"`
}

// DhcpDDNSConfig defines how and where to send DDNS updates.
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

// DhcpQueueControl controls queue behavior for asynchronous packet processing.
type DhcpQueueControl struct {
	Capacity    int    `json:"capacity"`
	EnableQueue bool   `json:"enable-queue"`
	QueueType   string `json:"queue-type"`
}

// ExpiredLeasesProcessing defines how expired or released leases are reclaimed.
type ExpiredLeasesProcessing struct {
	FlushReclaimedTimerWaitTime int `json:"flush-reclaimed-timer-wait-time"`
	HoldReclaimedTime           int `json:"hold-reclaimed-time"`
	MaxReclaimLeases            int `json:"max-reclaim-leases"`
	MaxReclaimTime              int `json:"max-reclaim-time"`
	ReclaimTimerWaitTime        int `json:"reclaim-timer-wait-time"`
	UnwarnedReclaimCycles       int `json:"unwarned-reclaim-cycles"`
}

// ServerID identifies the DHCPv6 server in transactions.
type ServerID struct {
	EnterpriseID int    `json:"enterprise-id"`
	HType        int    `json:"htype"`
	Identifier   string `json:"identifier"`
	Persist      bool   `json:"persist"`
	Time         int    `json:"time"`
	Type         string `json:"type"`
}
