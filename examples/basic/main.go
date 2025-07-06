package main

import (
	"fmt"
	"log"

	"github.com/rannday/kea-api/agent"
	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/internal/utils"
)

func main() {
	c := client.NewHTTP("http://192.168.66.2:8000/",
		client.WithAuth(&client.BasicAuth{
			Username: "kea-api",
			Password: "kea",
		}),
	)

	status, err := agent.StatusGet(c)
	if err != nil {
		log.Fatalf("Failed to get agent status: %v", err)
	}

	fmt.Println("Control Agent Status:")
	fmt.Printf("  PID:    %d\n", status.PID)
	fmt.Printf("  Uptime: %s\n", utils.HumanDuration(status.Uptime))
	fmt.Printf("  Reload: %s\n", utils.HumanDuration(status.Reload))
}
