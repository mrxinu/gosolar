# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoSolar is a Go client library for the SolarWinds Information Service (SWIS) REST API. It provides a production-ready wrapper around SWIS REST calls with modern Go practices including context-first APIs, structured error handling, generics, and strong typing.

## Architecture

The codebase follows modern Go patterns with the following structure:

- `gosolar.go` - Core client with context-first API, retries with exponential backoff, and structured errors
- `config.go` - Configuration management with validation and defaults
- `errors.go` - Structured error types with proper error handling
- `types.go` - Strongly typed structures for common SolarWinds entities
- `customproperties.go` - Custom property management with validation
- `custompollers.go` - Universal Device Poller (UnDP) operations
- `ncm.go` - Network Configuration Manager (NCM) operations
- `*_test.go` - Comprehensive unit tests with mocking
- `examples/` - Examples with environment configuration

### Core Client Structure

The `Client` uses modern configuration patterns:
- Configuration-based initialization with `Config` struct
- Context as first parameter on all methods (no legacy wrappers)
- Structured error handling with typed errors
- Built-in retry logic with exponential backoff
- Structured logging with `slog`
- Configurable port (default: 17774, the current SWIS REST API port)
- Connection pooling and timeout management

### Key Methods

**All methods take `context.Context` as the first parameter:**

**Core Operations:**
- `Query(ctx, swql, params)` - Execute SWQL query
- `QueryOne/QueryRow/QueryColumn` - Typed query helpers
- `Read(ctx, uri)` - Read entities
- `Create(ctx, entityType, properties)` - Create entities
- `Update(ctx, uri, body)` - Update entities
- `Delete(ctx, uri)` - Delete entities
- `Invoke(ctx, entity, verb, args)` - Execute verbs

**Specialized Operations:**
- Custom Properties: `SetCustomProperty`, `BulkSetCustomProperty`, `CreateCustomProperty`
- Custom Pollers: `GetAssignments`, `AddNodePoller`, `AddInterfacePoller`
- NCM: `RemoveNCMNodes`
- Bulk operations: `BulkDelete`
- Strongly typed helpers: `UnmarshalQueryResult[T]()`

## Development Commands

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage and race detection
go test -v -race -cover ./...

# Run specific test
go test -run TestNewClient

# Build examples
go build -o bin/simple-query examples/simple-query/main.go

# Format code
go fmt ./...

# Lint (if golangci-lint is installed)
golangci-lint run

# Generate documentation
go doc
godoc -http=:6060
```

## Testing

Comprehensive test suite with:
- Unit tests for all major components
- HTTP mocking for integration testing
- Error condition testing
- Context cancellation testing
- Configuration validation testing
- Port configuration testing

Run tests: `go test ./...`

## Usage Patterns

### Basic Client Creation
```go
config := gosolar.DefaultConfig()
config.Host = "solarwinds.example.com"
config.Username = "admin"
config.Password = os.Getenv("SOLARWINDS_PASSWORD")

client, err := gosolar.NewClient(config)
if err != nil {
    log.Fatal(err)
}
```

### Query with Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.Query(ctx, "SELECT * FROM Orion.Nodes", nil)
```

### Strongly Typed Results
```go
result, err := client.Query(ctx, "SELECT NodeID, Caption FROM Orion.Nodes", nil)
if err != nil {
    return err
}

queryResult, err := gosolar.UnmarshalQueryResult[gosolar.CommonNode](result)
if err != nil {
    return err
}

for _, node := range queryResult.Results {
    fmt.Printf("Node: %s\n", node.Caption)
}
```

### Error Handling
```go
if err != nil {
    var swErr *gosolar.Error
    if errors.As(err, &swErr) {
        switch swErr.Type {
        case gosolar.ErrorTypeAuthentication:
            // Handle auth errors
        case gosolar.ErrorTypeNetwork:
            // Handle network errors
        }
    }
}
```

### Custom Port (legacy SolarWinds installations)
```go
config := gosolar.DefaultConfig()
config.Port = 17778 // Legacy port, default is 17774
```

## Environment Variables

Examples support environment-based configuration:
- `SOLARWINDS_HOST` - SolarWinds server hostname
- `SOLARWINDS_USERNAME` - Authentication username
- `SOLARWINDS_PASSWORD` - Authentication password (required)
- `VENDOR_FILTER` - Vendor filter for parameterized queries
- `STATUS_FILTER` - Status filter for parameterized queries
