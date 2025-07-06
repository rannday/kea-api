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
