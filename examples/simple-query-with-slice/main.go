package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mrxinu/gosolar"
)

type Node struct {
	Caption   string `json:"caption"`
	IPAddress string `json:"ipaddress"`
}

func main() {
	// Configuration from environment
	hostname := getEnvOrDefault("SOLARWINDS_HOST", "localhost")
	username := getEnvOrDefault("SOLARWINDS_USERNAME", "admin")
	password := os.Getenv("SOLARWINDS_PASSWORD")
	vendor := getEnvOrDefault("VENDOR_FILTER", "Cisco")
	statusesStr := getEnvOrDefault("STATUS_FILTERS", "1,2,3")

	if password == "" {
		log.Fatal("SOLARWINDS_PASSWORD environment variable is required")
	}

	// Parse statuses from comma-separated string
	statuses, err := parseIntSlice(statusesStr)
	if err != nil {
		log.Fatalf("Invalid STATUS_FILTERS format: %v", err)
	}

	// Create modern client
	config := gosolar.DefaultConfig()
	config.Host = hostname
	config.Username = username
	config.Password = password
	config.InsecureSkipVerify = true // Only for demo
	config.Timeout = 30 * time.Second

	client, err := gosolar.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// SWQL query with slice parameter
	query := `
		SELECT
			Caption,
			IPAddress
		FROM Orion.Nodes
		WHERE Vendor = @vendor
		  AND Status IN @statuses
		ORDER BY Caption
	`

	// Parameters including slice for IN clause
	parameters := map[string]interface{}{
		"vendor":   vendor,
		"statuses": statuses,
	}

	fmt.Printf("Searching for %s nodes with statuses %v...\n", vendor, statuses)

	// Execute query with context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, err := client.QueryContext(ctx, query, parameters)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	// Unmarshal results using strongly typed struct
	var nodes []Node
	if err := json.Unmarshal(res, &nodes); err != nil {
		log.Fatalf("Failed to unmarshal results: %v", err)
	}

	// Display results
	if len(nodes) == 0 {
		fmt.Printf("No %s nodes found with statuses %v\n", vendor, statuses)
		return
	}

	fmt.Printf("Found %d matching nodes:\n", len(nodes))
	for i, node := range nodes {
		fmt.Printf("%d. %s (%s)\n", i+1, node.Caption, node.IPAddress)
	}
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseIntSlice(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}

	return result, nil
}
