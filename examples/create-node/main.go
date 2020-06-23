package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mrxinu/gosolar"
)

type Node struct {
	IPAddress     string
	EngineID      int
	ObjectSubType string
	SNMPVersion   int
	Community     string
	DNS           string
	NodeName      string
	SysName       string
}

func main() {
	hostname := "localhost"
	username := "admin"
	password := ""

	// NewClient creates a client that will handle the connection to SolarWinds
	// along with the timeout and HTTP conversation.
	client := gosolar.NewClient(hostname, username, password, true)

	node := Node{
		IPAddress:     "10.10.10.10",
		EngineID:      1,
		ObjectSubType: "SNMP",
		Community:     "SNMP community string here",
		SNMPVersion:   2,
		DNS:           "device.name.here",
		NodeName:      "device.name.here",
		SysName:       "device.name.here",
	}

	// Node creation
	res, err := client.Create("Orion.Nodes", node)

	if err != nil {
		fmt.Println(err)
	}

	// In case you want to get the NodeID
	re := regexp.MustCompile(`NodeID=\d+`)
	match := re.FindAllString(string(res), -1)
	if len(match) == 0 {
		fmt.Println("No NodeID found in create Node response")
	}

	fmt.Println(strconv.Atoi(strings.Replace(match[0], "NodeID=", "", -1)))
}
