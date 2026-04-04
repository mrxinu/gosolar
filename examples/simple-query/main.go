package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mrxinu/gosolar"
)

// Node represents a SolarWinds network node.
type Node struct {
	Caption   string `json:"caption"`
	IPAddress string `json:"ipaddress"`
}

func main() {
	hostname := getEnvOrDefault("SOLARWINDS_HOST", "localhost")
	username := getEnvOrDefault("SOLARWINDS_USERNAME", "admin")
	password := os.Getenv("SOLARWINDS_PASSWORD")

	if password == "" {
		log.Fatal("SOLARWINDS_PASSWORD environment variable is required")
	}

	config := gosolar.DefaultConfig()
	config.Host = hostname
	config.Username = username
	config.Password = password
	config.InsecureSkipVerify = true // Only for demo - use proper certificates in production

	client, err := gosolar.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Query(ctx, "SELECT Caption, IPAddress FROM Orion.Nodes", nil)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	var nodes []Node
	if err := json.Unmarshal(res, &nodes); err != nil {
		log.Fatalf("Failed to unmarshal results: %v", err)
	}

	fmt.Printf("Found %d nodes:\n", len(nodes))
	for i, node := range nodes {
		fmt.Printf("%d. %s (%s)\n", i+1, node.Caption, node.IPAddress)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
