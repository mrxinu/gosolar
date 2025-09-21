package gosolar

import (
	"log/slog"
	"time"
)

// Config holds configuration options for the SolarWinds client
type Config struct {
	// Host is the SolarWinds server hostname or IP
	Host string

	// Username for authentication
	Username string

	// Password for authentication
	Password string

	// InsecureSkipVerify controls whether SSL certificate verification is skipped
	InsecureSkipVerify bool

	// Timeout for HTTP requests (default: 30s)
	Timeout time.Duration

	// MaxIdleConns controls the maximum number of idle connections per host
	MaxIdleConns int

	// MaxRetries for failed requests (default: 3)
	MaxRetries int

	// RetryDelay between retry attempts (default: 1s)
	RetryDelay time.Duration

	// Logger for structured logging (optional)
	Logger *slog.Logger

	// UserAgent for HTTP requests
	UserAgent string
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Timeout:            30 * time.Second,
		MaxIdleConns:       10,
		MaxRetries:         3,
		RetryDelay:         time.Second,
		InsecureSkipVerify: false,
		UserAgent:          "gosolar/2.0",
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Host == "" {
		return NewError(ErrorTypeValidation, "config", "host is required")
	}
	if c.Username == "" {
		return NewError(ErrorTypeValidation, "config", "username is required")
	}
	if c.Password == "" {
		return NewError(ErrorTypeValidation, "config", "password is required")
	}
	if c.Timeout <= 0 {
		return NewError(ErrorTypeValidation, "config", "timeout must be positive")
	}
	if c.MaxRetries < 0 {
		return NewError(ErrorTypeValidation, "config", "max retries cannot be negative")
	}
	return nil
}