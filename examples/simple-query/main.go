package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mrxinu/gosolar"
)

func main() {
	hostname := "localhost"
	username := "admin"
	password := ""

	// NewClient creates a client that will handle the connection to SolarWinds
	// along with the timeout and HTTP conversation.
	client := gosolar.NewClient(hostname, username, password, true)

	// run the query without any parameters by passing nil as the 2nd parameter
	res, err := client.Query("SELECT Caption, IPAddress FROM Orion.Nodes", nil)
	if err != nil {
		log.Fatal(err)
	}

	// build a structure to unmarshal the results into
	var nodes []struct {
		Caption   string `json:"caption"`
		IPAddress string `json:"ipaddress"`
	}

	// let unmarshal do the work of unpacking the JSON
	if err := json.Unmarshal(res, &nodes); err != nil {
		log.Fatal(err)
	}

	// iterate over the resulting slice of node structures
	for _, n := range nodes {
		fmt.Printf("Working with node [%s] on IP address [%s]...\n", n.Caption, n.IPAddress)
	}
}
