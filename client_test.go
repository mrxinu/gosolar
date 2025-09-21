package gosolar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Host:     "example.com",
				Username: "admin",
				Password: "password",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: &Config{
				Username: "admin",
				Password: "password",
				Timeout:  30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "missing username",
			config: &Config{
				Host:     "example.com",
				Password: "password",
				Timeout:  30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			config: &Config{
				Host:     "example.com",
				Username: "admin",
				Password: "password",
				Timeout:  -1 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.config.Host, client.config.Host)
			}
		})
	}
}

func TestNewClientLegacy(t *testing.T) {
	client, err := NewClientLegacy("example.com", "admin", "password", true)
	require.NoError(t, err)
	require.NotNil(t, client)

	assert.Equal(t, "example.com", client.config.Host)
	assert.Equal(t, "admin", client.config.Username)
	assert.Equal(t, "password", client.config.Password)
	assert.True(t, client.config.InsecureSkipVerify)
}

func TestClient_QueryContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/SolarWinds/InformationService/v3/Json/Query", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "admin", username)
		assert.Equal(t, "password", password)

		response := map[string]interface{}{
			"results": []map[string]interface{}{
				{"NodeID": 1, "Caption": "Node1"},
				{"NodeID": 2, "Caption": "Node2"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := DefaultConfig()
	config.Host = server.URL[7:] // Remove "http://" prefix
	config.Username = "admin"
	config.Password = "password"
	config.InsecureSkipVerify = true

	client, err := NewClient(config)
	require.NoError(t, err)

	// Override the base URL to point to our test server
	client.baseURL.Scheme = "http"
	client.baseURL.Host = server.URL[7:]

	ctx := context.Background()
	result, err := client.QueryContext(ctx, "SELECT NodeID, Caption FROM Orion.Nodes", nil)

	require.NoError(t, err)
	assert.NotEmpty(t, result)

	var nodes []map[string]interface{}
	err = json.Unmarshal(result, &nodes)
	require.NoError(t, err)
	assert.Len(t, nodes, 2)
	assert.Equal(t, float64(1), nodes[0]["NodeID"])
	assert.Equal(t, "Node1", nodes[0]["Caption"])
}

func TestClient_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Authentication failed"))
	}))
	defer server.Close()

	config := DefaultConfig()
	config.Host = server.URL[7:]
	config.Username = "wrong"
	config.Password = "credentials"
	config.InsecureSkipVerify = true
	config.MaxRetries = 0 // Disable retries for this test

	client, err := NewClient(config)
	require.NoError(t, err)

	client.baseURL.Scheme = "http"
	client.baseURL.Host = server.URL[7:]

	ctx := context.Background()
	_, err = client.QueryContext(ctx, "SELECT * FROM Orion.Nodes", nil)

	require.Error(t, err)
	swErr, ok := err.(*Error)
	require.True(t, ok)
	assert.Equal(t, ErrorTypeAuthentication, swErr.Type)
	assert.Equal(t, 401, swErr.StatusCode)
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // Longer than context timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := DefaultConfig()
	config.Host = server.URL[7:]
	config.Username = "admin"
	config.Password = "password"
	config.Timeout = 10 * time.Second // Long client timeout
	config.MaxRetries = 0             // Disable retries

	client, err := NewClient(config)
	require.NoError(t, err)

	client.baseURL.Scheme = "http"
	client.baseURL.Host = server.URL[7:]

	// Use a short context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = client.QueryContext(ctx, "SELECT * FROM Orion.Nodes", nil)

	require.Error(t, err)
	// Should be a context/network error
	swErr, ok := err.(*Error)
	require.True(t, ok)
	assert.Equal(t, ErrorTypeNetwork, swErr.Type)
}

func TestClient_Retries(t *testing.T) {
	// Test that retries work for network errors
	// Use a non-existent host to trigger network errors
	config := DefaultConfig()
	config.Host = "nonexistent.invalid"
	config.Username = "admin"
	config.Password = "password"
	config.MaxRetries = 2
	config.RetryDelay = 10 * time.Millisecond
	config.Timeout = 100 * time.Millisecond

	client, err := NewClient(config)
	require.NoError(t, err)

	ctx := context.Background()
	_, err = client.QueryContext(ctx, "SELECT * FROM Orion.Nodes", nil)

	// Should fail with network error after retries
	require.Error(t, err)
	swErr, ok := err.(*Error)
	require.True(t, ok)
	assert.Equal(t, ErrorTypeNetwork, swErr.Type)
	assert.Contains(t, err.Error(), "request failed after retries")
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, time.Second, config.RetryDelay)
	assert.False(t, config.InsecureSkipVerify)
	assert.Equal(t, "gosolar/2.0", config.UserAgent)
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Host:       "example.com",
				Username:   "admin",
				Password:   "password",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
			},
			wantErr: false,
		},
		{
			name: "negative max retries",
			config: &Config{
				Host:       "example.com",
				Username:   "admin",
				Password:   "password",
				Timeout:    30 * time.Second,
				MaxRetries: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
