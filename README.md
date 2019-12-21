# gosolar

[![GoDoc](https://godoc.org/github.com/stobias123/gosolar?status.png)](http://godoc.org/github.com/stobias123/gosolar) [![Go Report Card](https://goreportcard.com/badge/github.com/stobias123/gosolar)](https://goreportcard.com/report/github.com/stobias123/gosolar)

GoSolar is a SolarWinds client library written in Go. It allows you
to submit queries to the SolarWinds Information Service (SWIS) and
do various other things.

## About

**stobias123/gosolar** is a wrapper around REST calls to the SWIS and makes
working with a SolarWinds install a little easier.

## Overview

GoSolar has the following generic methods:

* **Read** - read a SolarWinds object with all its properties.
* **Query** - query information via SWQL.
* **Create** - create new entities (nodes, pollers, etc.).
* **Delete** - delete an entity using its URI.
* **Invoke** - run verbs found in the SolarWinds API.

GoSolar has the following query wrappers for ease of use:

* **QueryOne** - returns a single `interface{}` from the query.
* **QueryRow** - returns a `[]byte` representing the single row.
* **QueryColumn** - returns a `[]interface{}` from the query.

GoSolar has the following convenience methods:

* Custom Properties
  * **SetCustomProperty** - set a custom property on a single entity.
  * **SetCustomProperties** - set custom properties on a single entity.
  * **BulkSetCustomProperties** - set a custom property on a series of entities.
  * **CreateCustomProperty** - create a custom property.
* Network Configuration Manager (NCM)
  * **RemoveNCMNodes** - remove nodes from NCM monitoring.
* Inventory Management
  * **BulkDelete** - delete multiple URIs in one request.
* Universal Device Poller (UnDP)
  * **GetAssignments** - get all the current UnDP assignments.
  * **AddNodePoller** - add a UnDP poller to a node.
  * **AddInterfacePoller** - add a UnDP poller to an interface.

## Installation

Install via **go get**:

```shell
go get -u github.com/stobias123/gosolar
```

## Documentation

See [http://godoc.org/github.com/stobias123/gosolar](http://godoc.org/github.com/stobias123/gosolar) or your local go doc
server for full documentation, as well as the examples.

```shell
cd $GOPATH
godoc -http=:6060 &
$preferred_browser http://localhost:6060/pkg &
```

## Usage

Basic usage can be found below but more specific examples are in the examples folder:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/stobias123/gosolar"
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
```

## Bugs

Please create an [issue](https://github.com/stobias123/gosolar/issues) on
GitHub with details about the bug and steps to reproduce it.
