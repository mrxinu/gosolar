package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mrxinu/gosolar"
)

func main() {
	hostname := getEnvOrDefault("SOLARWINDS_HOST", "localhost")
	username := getEnvOrDefault("SOLARWINDS_USERNAME", "admin")
	password := os.Getenv("SOLARWINDS_PASSWORD")
	nodeID := getEnvIntOrDefault("NODE_ID", 1)
	siteName := getEnvOrDefault("SITE_NAME", "Serenity Valley")

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.QueryOne(ctx, "SELECT URI FROM Orion.Nodes WHERE NodeID = @nodeID", map[string]int{"nodeID": nodeID})
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

	if err := client.SetCustomProperty(ctx, uri, "Site_Name", siteName); err != nil {
		log.Fatalf("Failed to set custom property: %v", err)
	}

	log.Println("Custom property updated successfully!")
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
