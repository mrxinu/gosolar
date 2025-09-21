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

// Node represents a SolarWinds network node
type Node struct {
	Caption   string `json:"caption"`
	IPAddress string `json:"ipaddress"`
}

func main() {
	// Get configuration from environment variables or use defaults
	hostname := getEnvOrDefault("SOLARWINDS_HOST", "localhost")
	username := getEnvOrDefault("SOLARWINDS_USERNAME", "admin")
	password := getEnvOrDefault("SOLARWINDS_PASSWORD", "")

	if password == "" {
		log.Fatal("SOLARWINDS_PASSWORD environment variable is required")
	}

	// Create client configuration with modern settings
	config := gosolar.DefaultConfig()
	config.Host = hostname
	config.Username = username
	config.Password = password
	config.InsecureSkipVerify = true // Only for demo - use proper certificates in production
	config.Timeout = 30 * time.Second

	// Create the client
	client, err := gosolar.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a context with timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute the query with context
	res, err := client.QueryContext(ctx, "SELECT Caption, IPAddress FROM Orion.Nodes", nil)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	// Unmarshal the results into our struct
	var nodes []Node
	if err := json.Unmarshal(res, &nodes); err != nil {
		log.Fatalf("Failed to unmarshal results: %v", err)
	}

	// Display the results
	fmt.Printf("Found %d nodes:\n", len(nodes))
	for i, node := range nodes {
		fmt.Printf("%d. %s (%s)\n", i+1, node.Caption, node.IPAddress)
	}
}

// getEnvOrDefault returns the environment variable value or a default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
