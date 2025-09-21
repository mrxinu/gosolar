package gosolar

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "with status code",
			err: &Error{
				Type:       ErrorTypeAuthentication,
				Operation:  "login",
				StatusCode: 401,
				Message:    "invalid credentials",
			},
			want: "gosolar authentication error in login (HTTP 401): invalid credentials",
		},
		{
			name: "without status code",
			err: &Error{
				Type:      ErrorTypeNetwork,
				Operation: "connect",
				Message:   "connection refused",
			},
			want: "gosolar network error in connect: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := &Error{
		Type:      ErrorTypeInternal,
		Operation: "test",
		Message:   "wrapped error",
		Cause:     originalErr,
	}

	assert.Equal(t, originalErr, errors.Unwrap(wrappedErr))
}

func TestError_Is(t *testing.T) {
	authErr := &Error{Type: ErrorTypeAuthentication}
	networkErr := &Error{Type: ErrorTypeNetwork}
	otherErr := errors.New("other error")

	assert.True(t, errors.Is(authErr, &Error{Type: ErrorTypeAuthentication}))
	assert.False(t, errors.Is(authErr, &Error{Type: ErrorTypeNetwork}))
	assert.False(t, errors.Is(authErr, otherErr))
	assert.False(t, errors.Is(networkErr, authErr))
}

func TestNewHTTPError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedType   ErrorType
		expectedStatus int
	}{
		{
			name:           "unauthorized",
			statusCode:     http.StatusUnauthorized,
			expectedType:   ErrorTypeAuthentication,
			expectedStatus: 401,
		},
		{
			name:           "forbidden",
			statusCode:     http.StatusForbidden,
			expectedType:   ErrorTypePermission,
			expectedStatus: 403,
		},
		{
			name:           "not found",
			statusCode:     http.StatusNotFound,
			expectedType:   ErrorTypeNotFound,
			expectedStatus: 404,
		},
		{
			name:           "bad request",
			statusCode:     http.StatusBadRequest,
			expectedType:   ErrorTypeValidation,
			expectedStatus: 400,
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			expectedType:   ErrorTypeInternal,
			expectedStatus: 500,
		},
		{
			name:           "other 4xx error",
			statusCode:     http.StatusConflict,
			expectedType:   ErrorTypeNetwork,
			expectedStatus: 409,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.WriteHeader(tt.statusCode)
			resp := recorder.Result()

			err := NewHTTPError("test_operation", "/test/endpoint", resp, "test message")

			assert.Equal(t, tt.expectedType, err.Type)
			assert.Equal(t, "test_operation", err.Operation)
			assert.Equal(t, "/test/endpoint", err.Endpoint)
			assert.Equal(t, tt.expectedStatus, err.StatusCode)
			assert.Equal(t, "test message", err.Message)
		})
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(originalErr, ErrorTypeNetwork, "test_operation", "wrapped message")

	assert.Equal(t, ErrorTypeNetwork, wrappedErr.Type)
	assert.Equal(t, "test_operation", wrappedErr.Operation)
	assert.Equal(t, "wrapped message", wrappedErr.Message)
	assert.Equal(t, originalErr, wrappedErr.Cause)
	assert.Equal(t, originalErr, errors.Unwrap(wrappedErr))
}

func TestNewError(t *testing.T) {
	err := NewError(ErrorTypeSWQL, "query", "invalid SWQL syntax")

	assert.Equal(t, ErrorTypeSWQL, err.Type)
	assert.Equal(t, "query", err.Operation)
	assert.Equal(t, "invalid SWQL syntax", err.Message)
	assert.Nil(t, err.Cause)
	assert.Equal(t, 0, err.StatusCode)
}