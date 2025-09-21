package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mrxinu/gosolar"
)

func main() {
	// Configuration from environment
	hostname := getEnvOrDefault("SOLARWINDS_HOST", "localhost")
	username := getEnvOrDefault("SOLARWINDS_USERNAME", "admin")
	password := os.Getenv("SOLARWINDS_PASSWORD")
	nodeID := getEnvIntOrDefault("NODE_ID", 1)
	siteName := getEnvOrDefault("SITE_NAME", "Serenity Valley")

	if password == "" {
		log.Fatal("SOLARWINDS_PASSWORD environment variable is required")
	}

	// Create modern client configuration
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the URI for the specified node
	res, err := client.QueryOneContext(ctx, "SELECT URI FROM Orion.Nodes WHERE NodeID = @nodeID", map[string]int{"nodeID": nodeID})
	if err != nil {
		log.Fatalf("Failed to query node URI: %v", err)
	}

	if res == nil {
		log.Fatalf("Node with ID %d not found", nodeID)
	}

	uri, ok := res.(string)
	if !ok {
		log.Fatalf("Expected string URI, got %T", res)
	}

	log.Printf("Setting custom property Site_Name='%s' on node %d (URI: %s)", siteName, nodeID, uri)

	// Set the custom property with context
	if err := client.SetCustomPropertyContext(ctx, uri, "Site_Name", siteName); err != nil {
		log.Fatalf("Failed to set custom property: %v", err)
	}

	log.Println("Custom property updated successfully!")
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, _ := time.ParseDuration(value); intVal > 0 {
			return int(intVal.Seconds())
		}
	}
	return defaultValue
}
