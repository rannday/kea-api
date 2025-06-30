package agent

// CtrlAgentStatus represents status-get when no service is specified.
type CtrlAgentStatus struct {
	PID    int `json:"pid"`
	Uptime int `json:"uptime"`
	Reload int `json:"reload"`
}
