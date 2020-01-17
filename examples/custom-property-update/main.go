package main

import (
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

	// get the URI for the first node
	res, err := client.QueryOne("SELECT URI FROM Orion.Nodes WHERE NodeID = @nodeID", map[string]int{"nodeID": 1})
	if err != nil {
		log.Fatal(err)
	}

	uri := res.(string) // cast to a string from interface{}

	// set the Site_Name property on that node (the custom property name is case insensitive)
	if err := client.SetCustomProperty(uri, "Site_Name", "Serenity Valley"); err != nil {
		log.Fatal(err)
	}
}
