package gosolar

import (
	"fmt"
	"net/http"
)

// ErrorType represents the category of error that occurred
type ErrorType string

const (
	ErrorTypeNetwork        ErrorType = "network"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypePermission     ErrorType = "permission"
	ErrorTypeSWQL           ErrorType = "swql"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeValidation     ErrorType = "validation"
	ErrorTypeInternal       ErrorType = "internal"
)

// Error represents a structured error from the SolarWinds API
type Error struct {
	Type       ErrorType `json:"type"`
	Operation  string    `json:"operation"`
	Endpoint   string    `json:"endpoint,omitempty"`
	StatusCode int       `json:"status_code,omitempty"`
	Message    string    `json:"message"`
	Cause      error     `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("gosolar %s error in %s (HTTP %d): %s", e.Type, e.Operation, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("gosolar %s error in %s: %s", e.Type, e.Operation, e.Message)
}

// Unwrap implements the errors.Unwrap interface
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is implements the errors.Is interface for error type comparison
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Type == t.Type
	}
	return false
}

// NewError creates a new structured error
func NewError(errType ErrorType, operation, message string) *Error {
	return &Error{
		Type:      errType,
		Operation: operation,
		Message:   message,
	}
}

// NewHTTPError creates a new error from an HTTP response
func NewHTTPError(operation, endpoint string, resp *http.Response, message string) *Error {
	var errType ErrorType
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		errType = ErrorTypeAuthentication
	case http.StatusForbidden:
		errType = ErrorTypePermission
	case http.StatusNotFound:
		errType = ErrorTypeNotFound
	case http.StatusBadRequest:
		errType = ErrorTypeValidation
	default:
		if resp.StatusCode >= 500 {
			errType = ErrorTypeInternal
		} else {
			errType = ErrorTypeNetwork
		}
	}

	return &Error{
		Type:       errType,
		Operation:  operation,
		Endpoint:   endpoint,
		StatusCode: resp.StatusCode,
		Message:    message,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, errType ErrorType, operation, message string) *Error {
	return &Error{
		Type:      errType,
		Operation: operation,
		Message:   message,
		Cause:     err,
	}
}
