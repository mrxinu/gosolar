package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mrxinu/gosolar"
)

// Node represents a SolarWinds network node with vendor information.
type Node struct {
	Caption   string `json:"caption"`
	IPAddress string `json:"ipaddress"`
}

func main() {
	hostname := getEnvOrDefault("SOLARWINDS_HOST", "localhost")
	username := getEnvOrDefault("SOLARWINDS_USERNAME", "admin")
	password := os.Getenv("SOLARWINDS_PASSWORD")
	vendor := getEnvOrDefault("VENDOR_FILTER", "Cisco")
	status := getEnvIntOrDefault("STATUS_FILTER", 1)

	if password == "" {
		log.Fatal("SOLARWINDS_PASSWORD environment variable is required")
	}

	config := gosolar.DefaultConfig()
	config.Host = hostname
	config.Username = username
	config.Password = password
	config.InsecureSkipVerify = true // Only for demo

	client, err := gosolar.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	query := `
		SELECT
			Caption,
			IPAddress
		FROM Orion.Nodes
		WHERE Vendor = @vendor
		  AND Status = @status
		ORDER BY Caption
	`

	parameters := map[string]any{
		"vendor": vendor,
		"status": status,
	}

	fmt.Printf("Searching for %s nodes with status %d...\n", vendor, status)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, err := client.Query(ctx, query, parameters)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	var nodes []Node
	if err := json.Unmarshal(res, &nodes); err != nil {
		log.Fatalf("Failed to unmarshal results: %v", err)
	}

	if len(nodes) == 0 {
		fmt.Printf("No %s nodes found with status %d\n", vendor, status)
		return
	}

	fmt.Printf("Found %d matching nodes:\n", len(nodes))
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

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
