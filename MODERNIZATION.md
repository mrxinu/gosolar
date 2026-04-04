# GoSolar Modernization Summary

This document summarizes the modernization changes made to the GoSolar library.

## v3 Changes (Current)

### Context-First API
All methods now take `context.Context` as the first parameter. The `*Context` suffix and legacy non-context wrappers have been removed.

**Before (v2):**
```go
result, err := client.Query("SELECT * FROM Orion.Nodes", nil)
// or
result, err := client.QueryContext(ctx, "SELECT * FROM Orion.Nodes", nil)
```

**After (v3):**
```go
result, err := client.Query(ctx, "SELECT * FROM Orion.Nodes", nil)
```

### Port Update
The default SWIS REST API port changed from 17778 to 17774, matching SolarWinds 2024.2+. The port is now configurable via `Config.Port` for legacy installations.

```go
config := gosolar.DefaultConfig()
config.Port = 17778 // For pre-2024.2 SolarWinds installations
```

### Exponential Backoff
Retry logic now uses exponential backoff instead of a fixed delay. `Config.RetryDelay` has been renamed to `Config.RetryBaseDelay`.

### Removed Legacy Methods
- `NewClientLegacy()` has been removed
- Non-context method wrappers (`Query()`, `Read()`, etc. without context) have been removed
- All methods now require `context.Context` as the first parameter

### Go 1.22+ Required
- Minimum Go version bumped from 1.21 to 1.22
- Uses `range` over integers (Go 1.22 feature)
- Uses `any` type alias throughout (replacing `interface{}`)

### Consistent Structured Errors
All modules now use the structured error system (`*Error` types). Previously, `custompollers.go` and `ncm.go` used plain `fmt.Errorf`.

### NewHTTPError Simplified
`NewHTTPError` now takes an `int` status code instead of `*http.Response`, avoiding unnecessary coupling.

### Retry Body Reuse Fix
The retry logic now properly recreates the request body for each attempt, fixing a bug where the body reader was consumed on the first attempt.

### Bug Fixes
- Fixed `getEnvIntOrDefault` in custom-property-update example (was using `time.ParseDuration` instead of `strconv.Atoi`)

## v2 Changes (Previous)

### Dependency Management
- Added `go.mod` for proper module management
- Testing dependencies (testify)

### Structured Error Handling
- `errors.go` with typed error system
- Error types: Network, Authentication, Permission, SWQL, NotFound, Validation, Internal
- Proper error wrapping and unwrapping support
- HTTP status code mapping to error types

### Configuration Management
- `config.go` with comprehensive configuration
- Configuration validation with sensible defaults
- Timeout, retry, and connection pool settings

### Strong Typing
- `types.go` with common SolarWinds entities
- Generic query result helpers (`UnmarshalQueryResult[T]`)
- Constants for status values and alert severities
- Type-safe custom property management

### Comprehensive Testing
- `client_test.go`, `errors_test.go`
- HTTP mocking for integration testing
- Error condition, context cancellation, and configuration validation testing

## Migration Guide (v2 to v3)

1. **Add context to all calls**: Every client method now requires `context.Context` as the first parameter
2. **Remove `Context` suffix**: `QueryContext` becomes `Query`, `ReadContext` becomes `Read`, etc.
3. **Replace `NewClientLegacy`**: Use `NewClient` with a `Config` struct
4. **Update retry config**: Rename `RetryDelay` to `RetryBaseDelay`
5. **Check port**: Default port is now 17774. Set `config.Port = 17778` for older SolarWinds
6. **Replace `interface{}`**: Use `any` (optional, both work)
7. **Update Go version**: Requires Go 1.22+
