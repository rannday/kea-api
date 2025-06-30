package main

import (
	"fmt"
	"log"

	"github.com/rannday/isc-kea/agent"
	"github.com/rannday/isc-kea/utils"
)

func main() {
	client := agent.NewHTTPClient("http://192.168.66.2:8000/",
		agent.WithAuth(&agent.BasicAuth{
			Username: "kea-api",
			Password: "kea",
		}),
	)

	status, err := agent.StatusGet(client)
	if err != nil {
		log.Fatalf("Failed to get agent status: %v", err)
	}

	fmt.Printf("Control Agent Status:\n")
	fmt.Printf("  PID:    %d\n", status.PID)
	fmt.Printf("  Uptime: %s\n", utils.HumanDuration(status.Uptime))
	fmt.Printf("  Reload: %s\n", utils.HumanDuration(status.Reload))
}
