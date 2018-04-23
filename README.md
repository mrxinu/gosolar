# gosolar

[![GoDoc](https://godoc.org/github.com/mrxinu/gosolar?status.png)](http://godoc.org/github.com/mrxinu/gosolar) [![Go Report Card](https://goreportcard.com/badge/github.com/mrxinu/gosolar)](https://goreportcard.com/report/github.com/mrxinu/gosolar)

GoSolar is a SolarWinds client library written in Go. It allows you
to submit queries to the SolarWinds Information Service (SWIS) and
do various other things.

## About

**mrxinu/gosolar** is a wrapper around REST calls to the SWIS and makes
working with a SolarWinds install a little easier.

## Overview

GoSolar has the following generic methods:

* **Read** - read a SolarWinds object with all its properties.
* **Query** - query information via SWQL.
* **Create** - create new entities (nodes, pollers, etc.).
* **Delete** - delete an entity using its URI.
* **Invoke** - run verbs found in the SolarWinds API.

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

GoSolar is under development, so I would not start using this straight
away. Initially it's here so I can write some utilities without having
to rework the client code every time I do.

## Installation

Install via **go get**:

```shell
go get -u github.com/mrxinu/gosolar
```

## Documentation

See [http://godoc.org/github.com/mrxinu/gosolar](http://godoc.org/github.com/mrxinu/gosolar) or your local go doc
server for full documentation, as well as the examples.

```shell
cd $GOPATH
godoc -http=:6060 &
$preferred_browser http://localhost:6060/pkg &
```

## Testing

```
make test
```

The `test` make target will test the entire `gosolar` package.

## Usage

TBD

## Bugs

Please create an [issue](https://github.com/mrxinu/gosolar/issues) on
GitHub with details about the bug and steps to reproduce it.
