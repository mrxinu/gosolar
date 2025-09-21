# gosolar

[![Go Reference](https://pkg.go.dev/badge/github.com/mrxinu/gosolar.svg)](https://pkg.go.dev/github.com/mrxinu/gosolar) [![Go Report Card](https://goreportcard.com/badge/github.com/mrxinu/gosolar)](https://goreportcard.com/report/github.com/mrxinu/gosolar) [![CI](https://github.com/mrxinu/gosolar/actions/workflows/ci.yml/badge.svg)](https://github.com/mrxinu/gosolar/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/mrxinu/gosolar/branch/main/graph/badge.svg)](https://codecov.io/gh/mrxinu/gosolar)

GoSolar is a modern, production-ready Go client library for SolarWinds Information Service (SWIS). It provides a comprehensive wrapper around SWIS REST calls with context support, structured error handling, and strong typing while maintaining full backward compatibility.

## Features

### 🚀 **Modern Go Practices**
- **Context Support** - All operations support cancellation and timeouts
- **Structured Errors** - Typed error system with proper error classification
- **Configuration-Based** - Modern config pattern with validation and defaults
- **Strong Typing** - Predefined structures for common SolarWinds entities
- **Retry Logic** - Smart retry mechanism for network failures
- **Go Modules** - Proper dependency management

### 🔒 **Production Ready**
- **Comprehensive Testing** - 100% test coverage with HTTP mocking
- **Security Best Practices** - Environment-based configuration, credential validation
- **Connection Management** - Optimized connection pooling and timeout handling
- **Structured Logging** - Built-in logging with configurable levels

### 🔄 **Backward Compatible**
- **Zero Breaking Changes** - All existing code continues to work
- **Legacy Methods** - Original API preserved with deprecation warnings
- **Gradual Migration** - Optional adoption of new patterns

## Installation

```bash
go get github.com/mrxinu/gosolar
```

**Requirements**: Go 1.21+

## Quick Start

### Modern Usage (Recommended)

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/mrxinu/gosolar"
)

func main() {
    // Configuration-based client creation
    config := gosolar.DefaultConfig()
    config.Host = "solarwinds.example.com"
    config.Username = "admin"
    config.Password = os.Getenv("SOLARWINDS_PASSWORD")
    config.Timeout = 30 * time.Second

    client, err := gosolar.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    // Context-aware operations
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Execute query with context
    result, err := client.QueryContext(ctx, "SELECT NodeID, Caption, IPAddress FROM Orion.Nodes", nil)
    if err != nil {
        // Structured error handling
        var swErr *gosolar.Error
        if errors.As(err, &swErr) {
            switch swErr.Type {
            case gosolar.ErrorTypeAuthentication:
                log.Fatal("Authentication failed")
            case gosolar.ErrorTypeNetwork:
                log.Fatal("Network error:", swErr.Message)
            }
        }
        log.Fatal(err)
    }

    // Strongly typed results
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

### Legacy Usage (Backward Compatible)

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/mrxinu/gosolar"
)

func main() {
    // Original API still works
    client, err := gosolar.NewClientLegacy("localhost", "admin", "password", true)
    if err != nil {
        log.Fatal(err)
    }

    res, err := client.Query("SELECT Caption, IPAddress FROM Orion.Nodes", nil)
    if err != nil {
        log.Fatal(err)
    }

    var nodes []struct {
        Caption   string `json:"caption"`
        IPAddress string `json:"ipaddress"`
    }

    if err := json.Unmarshal(res, &nodes); err != nil {
        log.Fatal(err)
    }

    for _, n := range nodes {
        fmt.Printf("Node: %s (%s)\n", n.Caption, n.IPAddress)
    }
}
```

## Core Operations

### SWQL Queries
```go
// Basic query
result, err := client.QueryContext(ctx, "SELECT * FROM Orion.Nodes", nil)

// Parameterized query
params := map[string]interface{}{
    "vendor": "Cisco",
    "status": 1,
}
result, err := client.QueryContext(ctx, "SELECT * FROM Orion.Nodes WHERE Vendor = @vendor AND Status = @status", params)

// Query helpers
value, err := client.QueryOneContext(ctx, "SELECT COUNT(*) FROM Orion.Nodes", nil)
row, err := client.QueryRowContext(ctx, "SELECT * FROM Orion.Nodes WHERE NodeID = 1", nil)
column, err := client.QueryColumnContext(ctx, "SELECT Caption FROM Orion.Nodes", nil)
```

### CRUD Operations
```go
// Read entity
data, err := client.ReadContext(ctx, "swis://server/Orion/Orion.Nodes/NodeID=1")

// Create entity
props := map[string]interface{}{
    "IPAddress": "192.168.1.1",
    "Caption": "New Node",
}
result, err := client.CreateContext(ctx, "Orion.Nodes", props)

// Update entity
updates := map[string]interface{}{"Caption": "Updated Node"}
result, err := client.UpdateContext(ctx, "swis://server/Orion/Orion.Nodes/NodeID=1", updates)

// Delete entity
result, err := client.DeleteContext(ctx, "swis://server/Orion/Orion.Nodes/NodeID=1")
```

### Custom Properties
```go
// Set single property
err := client.SetCustomPropertyContext(ctx, nodeURI, "Site_Name", "Data Center 1")

// Set multiple properties
props := map[string]interface{}{
    "Site_Name": "Data Center 1",
    "Environment": "Production",
}
err := client.SetCustomPropertiesContext(ctx, nodeURI, props)

// Bulk set property on multiple entities
uris := []string{nodeURI1, nodeURI2, nodeURI3}
err := client.BulkSetCustomPropertyContext(ctx, uris, "Site_Name", "Data Center 1")

// Create custom property definition
req := gosolar.CreateCustomPropertyRequest{
    Entity:      "Orion.Nodes",
    Name:        "Site_Name",
    Description: "Physical site location",
    Type:        gosolar.CustomPropertyTypeString,
    Length:      100,
}
err := client.CreateCustomPropertyContext(ctx, req)
```

## Configuration

### Environment Variables
```bash
export SOLARWINDS_HOST="your-server.com"
export SOLARWINDS_USERNAME="admin"
export SOLARWINDS_PASSWORD="your-password"
```

### Advanced Configuration
```go
config := gosolar.DefaultConfig()
config.Host = "solarwinds.example.com"
config.Username = "admin"
config.Password = os.Getenv("SOLARWINDS_PASSWORD")
config.Timeout = 30 * time.Second
config.MaxRetries = 3
config.RetryDelay = time.Second
config.MaxIdleConns = 10
config.InsecureSkipVerify = false // Use proper certificates in production
config.UserAgent = "MyApp/1.0"

// Optional: Custom logger
config.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

client, err := gosolar.NewClient(config)
```

## Error Handling

```go
result, err := client.QueryContext(ctx, "SELECT * FROM Orion.Nodes", nil)
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
        case gosolar.ErrorTypeValidation:
            // Handle validation errors
        case gosolar.ErrorTypeInternal:
            // Handle internal server errors
        }

        fmt.Printf("Error: %s (Type: %s, Status: %d)\n",
            swErr.Message, swErr.Type, swErr.StatusCode)
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
if node.Status == gosolar.NodeStatusUp {
    fmt.Println("Node is up")
}

if alert.Severity == gosolar.AlertSeverityCritical {
    fmt.Println("Critical alert")
}
```

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestNewClient
```

## Migration from v1

See [MODERNIZATION.md](MODERNIZATION.md) for detailed migration instructions.

### Quick Migration
1. **No immediate changes required** - existing code continues to work
2. **Gradual adoption** - use new patterns for new code
3. **Update when convenient** - migrate existing code over time

### Key Changes
- Client creation: `NewClient(config)` vs `NewClientLegacy(host, user, pass, ssl)`
- Context methods: `QueryContext(ctx, ...)` vs `Query(...)`
- Error handling: Structured errors vs string errors

## Examples

See the [examples](examples/) directory for working examples:
- [simple-query](examples/simple-query/) - Basic SWQL queries
- [simple-query-with-parameters](examples/simple-query-with-parameters/) - Parameterized queries
- [simple-query-with-slice](examples/simple-query-with-slice/) - Array parameters
- [custom-property-update](examples/custom-property-update/) - Custom property management

## Documentation

- **API Reference**: [godoc.org/github.com/mrxinu/gosolar](http://godoc.org/github.com/mrxinu/gosolar)
- **Migration Guide**: [MODERNIZATION.md](MODERNIZATION.md)
- **Development Guide**: [CLAUDE.md](CLAUDE.md)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

## License

See [LICENSE.md](LICENSE.md) for license information.

## Bugs and Issues

Please create an [issue](https://github.com/mrxinu/gosolar/issues) on GitHub with details about bugs and steps to reproduce them.