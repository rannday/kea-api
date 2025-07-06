package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rannday/kea-api/agent"
	"github.com/rannday/kea-api/client"
	"github.com/rannday/kea-api/internal/utils"
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

	transport := client.NewHTTPTransport(
		"http://192.168.66.2:8000/",
		client.WithHTTPClient(httpClient),
		client.WithAuth(&client.BasicAuth{
			Username: "kea-api",
			Password: "kea",
		}),
	)

	c := client.NewClient(transport)

	status, err := agent.StatusGet(c)
	if err != nil {
		log.Fatalf("Failed to get agent status: %v", err)
	}

	fmt.Println("Control Agent Status:")
	fmt.Printf("  PID:    %d\n", status.PID)
	fmt.Printf("  Uptime: %s\n", utils.HumanDuration(status.Uptime))
	fmt.Printf("  Reload: %s\n", utils.HumanDuration(status.Reload))
}
