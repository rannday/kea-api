package agent

// CtrlAgentStatus represents status-get when no service is specified.
type CtrlAgentStatus struct {
	PID    int `json:"pid"`
	Uptime int `json:"uptime"`
	Reload int `json:"reload"`
}

// CtrlAgentVersion is the response type for version-get on the control-agent.
type CtrlAgentVersion struct {
	Extended string `json:"extended"`
}
