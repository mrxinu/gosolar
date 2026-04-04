# gosolar

[![Go Reference](https://pkg.go.dev/badge/github.com/mrxinu/gosolar.svg)](https://pkg.go.dev/github.com/mrxinu/gosolar) [![Go Report Card](https://goreportcard.com/badge/github.com/mrxinu/gosolar)](https://goreportcard.com/report/github.com/mrxinu/gosolar) [![CI](https://github.com/mrxinu/gosolar/actions/workflows/ci.yml/badge.svg)](https://github.com/mrxinu/gosolar/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/mrxinu/gosolar/branch/main/graph/badge.svg)](https://codecov.io/gh/mrxinu/gosolar)

GoSolar is a Go client library for the SolarWinds Information Service (SWIS) REST API. It provides a production-ready wrapper around SWIS REST calls with context-first APIs, structured error handling, generics, and strong typing.

## Features

- **Context-First API** - All operations take `context.Context` as the first parameter
- **Structured Errors** - Typed error system with proper error classification
- **Configuration-Based** - Config pattern with validation and sensible defaults
- **Strong Typing** - Predefined structures for common SolarWinds entities with generics
- **Exponential Backoff** - Retry logic with exponential backoff for transient failures
- **Structured Logging** - Built-in `slog` integration
- **Configurable Port** - Default port 17774 (current SWIS), configurable for legacy (17778)

## Installation

```bash
go get github.com/mrxinu/gosolar
```

**Requirements**: Go 1.22+

## Quick Start

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/mrxinu/gosolar"
)

func main() {
    config := gosolar.DefaultConfig()
    config.Host = "solarwinds.example.com"
    config.Username = "admin"
    config.Password = os.Getenv("SOLARWINDS_PASSWORD")

    client, err := gosolar.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := client.Query(ctx, "SELECT NodeID, Caption, IPAddress FROM Orion.Nodes", nil)
    if err != nil {
        var swErr *gosolar.Error
        if errors.As(err, &swErr) {
            log.Fatalf("%s error: %s", swErr.Type, swErr.Message)
        }
        log.Fatal(err)
    }

    queryResult, err := gosolar.UnmarshalQueryResult[gosolar.CommonNode](result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d nodes:\n", queryResult.Count)
    for _, node := range queryResult.Results {
        fmt.Printf("- %s (%s)\n", node.Caption, node.IPAddress)
    }
}
```

## Core Operations

### SWQL Queries
```go
// Basic query
result, err := client.Query(ctx, "SELECT * FROM Orion.Nodes", nil)

// Parameterized query
params := map[string]any{
    "vendor": "Cisco",
    "status": 1,
}
result, err := client.Query(ctx, "SELECT * FROM Orion.Nodes WHERE Vendor = @vendor AND Status = @status", params)

// Query helpers
value, err := client.QueryOne(ctx, "SELECT COUNT(*) FROM Orion.Nodes", nil)
row, err := client.QueryRow(ctx, "SELECT * FROM Orion.Nodes WHERE NodeID = 1", nil)
column, err := client.QueryColumn(ctx, "SELECT Caption FROM Orion.Nodes", nil)
```

### CRUD Operations
```go
// Read entity
data, err := client.Read(ctx, "swis://server/Orion/Orion.Nodes/NodeID=1")

// Create entity
props := map[string]any{
    "IPAddress": "192.168.1.1",
    "Caption":   "New Node",
}
result, err := client.Create(ctx, "Orion.Nodes", props)

// Update entity
updates := map[string]any{"Caption": "Updated Node"}
result, err := client.Update(ctx, "swis://server/Orion/Orion.Nodes/NodeID=1", updates)

// Delete entity
result, err := client.Delete(ctx, "swis://server/Orion/Orion.Nodes/NodeID=1")
```

### Custom Properties
```go
// Set single property
err := client.SetCustomProperty(ctx, nodeURI, "Site_Name", "Data Center 1")

// Set multiple properties
props := map[string]any{
    "Site_Name":   "Data Center 1",
    "Environment": "Production",
}
err := client.SetCustomProperties(ctx, nodeURI, props)

// Bulk set property on multiple entities
uris := []string{nodeURI1, nodeURI2, nodeURI3}
err := client.BulkSetCustomProperty(ctx, uris, "Site_Name", "Data Center 1")

// Create custom property definition
req := gosolar.CreateCustomPropertyRequest{
    Entity:      "Orion.Nodes",
    Name:        "Site_Name",
    Description: "Physical site location",
    Type:        gosolar.CustomPropertyTypeString,
    Length:      100,
}
err := client.CreateCustomProperty(ctx, req)
```

## Configuration

```go
config := gosolar.DefaultConfig()
config.Host = "solarwinds.example.com"
config.Port = 17774              // Default; use 17778 for legacy installations
config.Username = "admin"
config.Password = os.Getenv("SOLARWINDS_PASSWORD")
config.Timeout = 30 * time.Second
config.MaxRetries = 3
config.RetryBaseDelay = time.Second
config.MaxIdleConns = 10
config.InsecureSkipVerify = false // Use proper certificates in production
config.UserAgent = "MyApp/1.0"
config.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

client, err := gosolar.NewClient(config)
```

## Error Handling

```go
result, err := client.Query(ctx, "SELECT * FROM Orion.Nodes", nil)
if err != nil {
    var swErr *gosolar.Error
    if errors.As(err, &swErr) {
        switch swErr.Type {
        case gosolar.ErrorTypeAuthentication:
            // Handle authentication errors
        case gosolar.ErrorTypePermission:
            // Handle permission errors
        case gosolar.ErrorTypeNetwork:
            // Handle network errors
        case gosolar.ErrorTypeSWQL:
            // Handle SWQL syntax errors
        case gosolar.ErrorTypeNotFound:
            // Handle not found errors
        }
    }
}
```

## Predefined Types

```go
// Common SolarWinds entities
var nodes []gosolar.CommonNode
var interfaces []gosolar.Interface
var alerts []gosolar.Alert
var volumes []gosolar.Volume
var applications []gosolar.Application

// Status constants
if node.Status == gosolar.NodeStatusUp { ... }
if alert.Severity == gosolar.AlertSeverityCritical { ... }
```

## Testing

```bash
go test ./...                          # Run all tests
go test -v -race -cover ./...          # Verbose with race detection and coverage
go test -run TestNewClient             # Run specific test
```

## Examples

See the [examples](examples/) directory:
- [simple-query](examples/simple-query/) - Basic SWQL queries
- [simple-query-with-parameters](examples/simple-query-with-parameters/) - Parameterized queries
- [simple-query-with-slice](examples/simple-query-with-slice/) - Array parameters
- [custom-property-update](examples/custom-property-update/) - Custom property management

## Migration from v2

See [MODERNIZATION.md](MODERNIZATION.md) for the full migration guide.

Key changes in v3:
- All methods now take `context.Context` as the first parameter (no more `*Context` suffix)
- Removed `NewClientLegacy()` and non-context method wrappers
- Default port changed from 17778 to 17774 (SolarWinds 2024.2+)
- `RetryDelay` renamed to `RetryBaseDelay` (now uses exponential backoff)
- `interface{}` replaced with `any` throughout
- Uses `any` type alias (`interface{}` still works but `any` is idiomatic)

## License

See [LICENSE.md](LICENSE.md) for license information.

## Bugs and Issues

Please create an [issue](https://github.com/mrxinu/gosolar/issues) on GitHub with details about bugs and steps to reproduce them.
