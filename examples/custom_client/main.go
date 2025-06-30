package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rannday/isc-kea/agent"
	"github.com/rannday/isc-kea/utils"
)

func main() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       50,
			IdleConnTimeout:    90 * time.Second,
			DisableKeepAlives:  false,
			DisableCompression: false,
		},
	}

	client := agent.NewHTTPClient(
		"http://192.168.66.2:8000/",
		agent.WithHTTPClient(httpClient),
		agent.WithAuth(&agent.BasicAuth{
			Username: "kea-api",
			Password: "kea",
		}),
	)

	status, err := agent.StatusGet(client)
	if err != nil {
		log.Fatalf("Failed to get agent status: %v", err)
	}

	fmt.Println("Control Agent Status:")
	fmt.Printf("  PID:    %d\n", status.PID)
	fmt.Printf("  Uptime: %s\n", utils.HumanDuration(status.Uptime))
	fmt.Printf("  Reload: %s\n", utils.HumanDuration(status.Reload))
}
