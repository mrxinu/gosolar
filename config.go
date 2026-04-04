package gosolar

import (
	"log/slog"
	"time"
)

// DefaultPort is the default SWIS REST API port (17774 as of SolarWinds 2024.2).
const DefaultPort = 17774

// Config holds configuration options for the SolarWinds client.
type Config struct {
	// Host is the SolarWinds server hostname or IP.
	Host string

	// Port is the SWIS REST API port (default: 17774).
	Port int

	// Username for authentication.
	Username string

	// Password for authentication.
	Password string

	// InsecureSkipVerify controls whether TLS certificate verification is skipped.
	InsecureSkipVerify bool

	// Timeout for HTTP requests (default: 30s).
	Timeout time.Duration

	// MaxIdleConns controls the maximum number of idle connections per host.
	MaxIdleConns int

	// MaxRetries for failed requests (default: 3).
	MaxRetries int

	// RetryBaseDelay is the base delay for exponential backoff between retries (default: 1s).
	RetryBaseDelay time.Duration

	// Logger for structured logging (optional, defaults to slog.Default).
	Logger *slog.Logger

	// UserAgent for HTTP requests.
	UserAgent string
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Port:               DefaultPort,
		Timeout:            30 * time.Second,
		MaxIdleConns:       10,
		MaxRetries:         3,
		RetryBaseDelay:     time.Second,
		InsecureSkipVerify: false,
		UserAgent:          "gosolar/3.0",
	}
}

// Validate checks if the configuration is valid.
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
	if c.Port <= 0 || c.Port > 65535 {
		return NewError(ErrorTypeValidation, "config", "port must be between 1 and 65535")
	}
	return nil
}
