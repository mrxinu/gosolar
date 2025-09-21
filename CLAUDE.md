# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoSolar is a modernized Go client library for SolarWinds Information Service (SWIS) REST API. It provides a comprehensive, production-ready wrapper around SWIS REST calls with modern Go practices including context support, structured error handling, and strong typing.

## Architecture

The codebase follows modern Go patterns with the following structure:

- `gosolar.go` - Core client with context support, retries, and structured errors
- `config.go` - Configuration management with validation and defaults
- `errors.go` - Structured error types with proper error handling
- `types.go` - Strongly typed structures for common SolarWinds entities
- `customproperties.go` - Custom property management with validation
- `custompollers.go` - Universal Device Poller (UnDP) operations
- `ncm.go` - Network Configuration Manager (NCM) operations
- `*_test.go` - Comprehensive unit tests with mocking
- `examples/` - Modernized examples with environment configuration

### Core Client Structure

The `Client` uses modern configuration patterns:
- Configuration-based initialization with `Config` struct
- Context support for all operations (cancellation, timeouts)
- Structured error handling with typed errors
- Built-in retry logic with exponential backoff
- Structured logging with `slog`
- Connection pooling and timeout management

### Key Methods

**All methods have both legacy and context-aware versions:**

**Core Operations:**
- `QueryContext(ctx, swql, params)` - Execute SWQL with context
- `QueryOneContext/QueryRowContext/QueryColumnContext` - Typed query helpers
- `ReadContext(ctx, uri)` - Read entities with context
- `CreateContext(ctx, entityType, properties)` - Create with context
- `DeleteContext(ctx, uri)` - Delete with context
- `InvokeContext(ctx, entity, verb, args)` - Execute verbs with context

**Specialized Operations:**
- Custom Properties: `SetCustomPropertyContext`, `BulkSetCustomPropertyContext`, `CreateCustomPropertyContext`
- Bulk operations: `BulkDeleteContext`
- Strongly typed helpers: `UnmarshalQueryResult[T]()`

## Development Commands

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

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

Run tests: `go test ./...`

## Modern Usage Patterns

### Basic Client Creation
```go
config := gosolar.DefaultConfig()
config.Host = "solarwinds.example.com"
config.Username = "admin"
config.Password = os.Getenv("SOLARWINDS_PASSWORD")
config.Timeout = 30 * time.Second

client, err := gosolar.NewClient(config)
if err != nil {
    log.Fatal(err)
}
```

### Context-Aware Operations
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.QueryContext(ctx, "SELECT * FROM Orion.Nodes", nil)
```

### Strongly Typed Results
```go
result, err := client.QueryContext(ctx, "SELECT NodeID, Caption FROM Orion.Nodes", nil)
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

## Environment Variables

Examples support environment-based configuration:
- `SOLARWINDS_HOST` - SolarWinds server hostname
- `SOLARWINDS_USERNAME` - Authentication username
- `SOLARWINDS_PASSWORD` - Authentication password (required)
- `VENDOR_FILTER` - Vendor filter for parameterized queries
- `STATUS_FILTER` - Status filter for parameterized queries

## Backward Compatibility

Legacy methods are preserved with deprecation warnings:
- `NewClientLegacy()` - Original constructor
- Non-context methods delegate to context versions

New code should use the modern patterns with configuration and context support.