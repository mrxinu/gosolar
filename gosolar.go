package gosolar

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

// Client represents a SolarWinds SWIS API client
type Client struct {
	config *Config
	baseURL *url.URL
	httpClient *http.Client
	logger *slog.Logger
}

// NewClient creates a new SolarWinds client with the provided configuration
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(fmt.Sprintf("https://%s:17778/SolarWinds/InformationService/v3/Json/", config.Host))
	if err != nil {
		return nil, WrapError(err, ErrorTypeValidation, "new_client", "invalid host URL")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
		MaxIdleConnsPerHost: config.MaxIdleConns,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	logger := config.Logger
	if logger == nil {
		logger = slog.Default()
	}

	return &Client{
		config:     config,
		baseURL:    baseURL,
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// NewClientLegacy creates a client using the legacy constructor signature for backward compatibility
// Deprecated: Use NewClient with Config instead
func NewClientLegacy(host, user, pass string, ignoreSSL bool) (*Client, error) {
	config := DefaultConfig()
	config.Host = host
	config.Username = user
	config.Password = pass
	config.InsecureSkipVerify = ignoreSSL
	return NewClient(config)
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, WrapError(err, ErrorTypeValidation, "request", "failed to marshal request body")
		}
		reqBody = &buf
	}

	endpointURL, err := c.baseURL.Parse(endpoint)
	if err != nil {
		return nil, WrapError(err, ErrorTypeValidation, "request", "invalid endpoint")
	}

	req, err := http.NewRequestWithContext(ctx, method, endpointURL.String(), reqBody)
	if err != nil {
		return nil, WrapError(err, ErrorTypeNetwork, "request", "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.SetBasicAuth(c.config.Username, c.config.Password)

	c.logger.DebugContext(ctx, "making request", "method", method, "endpoint", endpoint)

	var resp *http.Response
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			c.logger.DebugContext(ctx, "retrying request", "attempt", attempt)
			time.Sleep(c.config.RetryDelay)
		}

		resp, err = c.httpClient.Do(req)
		if err == nil {
			break
		}
		lastErr = err
	}

	if err != nil {
		return nil, WrapError(lastErr, ErrorTypeNetwork, "request", "request failed after retries")
	}
	defer resp.Body.Close()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, WrapError(err, ErrorTypeNetwork, "response", "failed to read response body")
	}

	if resp.StatusCode >= 400 {
		return nil, NewHTTPError("request", endpoint, resp, string(output))
	}

	c.logger.DebugContext(ctx, "request completed", "status", resp.StatusCode)
	return output, nil
}

func (c *Client) post(endpoint string, body interface{}) ([]byte, error) {
	return c.PostContext(context.Background(), endpoint, body)
}

// PostContext performs a POST request with context
func (c *Client) PostContext(ctx context.Context, endpoint string, body interface{}) ([]byte, error) {
	return c.doRequest(ctx, "POST", endpoint, body)
}

func (c *Client) get(endpoint string) ([]byte, error) {
	return c.GetContext(context.Background(), endpoint)
}

// GetContext performs a GET request with context
func (c *Client) GetContext(ctx context.Context, endpoint string) ([]byte, error) {
	return c.doRequest(ctx, "GET", endpoint, nil)
}

// Query executes a SWQL query and returns the results
func (c *Client) Query(query string, parameters interface{}) ([]byte, error) {
	return c.QueryContext(context.Background(), query, parameters)
}

// QueryContext executes a SWQL query with context
func (c *Client) QueryContext(ctx context.Context, query string, parameters interface{}) ([]byte, error) {
	req := struct {
		Query      string      `json:"query"`
		Parameters interface{} `json:"parameters"`
	}{
		Query:      query,
		Parameters: parameters,
	}

	result, err := c.PostContext(ctx, "Query", &req)
	if err != nil {
		// Preserve the original error type if it's already a structured error
		if swErr, ok := err.(*Error); ok {
			return nil, swErr
		}
		return nil, WrapError(err, ErrorTypeSWQL, "query", "SWQL query failed")
	}

	sr := struct {
		Result *json.RawMessage `json:"results"`
	}{}

	if err := json.Unmarshal(result, &sr); err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "query", "failed to parse query response")
	}

	if sr.Result == nil {
		return []byte("[]"), nil
	}

	return []byte(*sr.Result), nil
}

// QueryOne executes a query and returns a single value
func (c *Client) QueryOne(query string, parameters interface{}) (interface{}, error) {
	return c.QueryOneContext(context.Background(), query, parameters)
}

// QueryOneContext executes a query with context and returns a single value
func (c *Client) QueryOneContext(ctx context.Context, query string, parameters interface{}) (interface{}, error) {
	res, err := c.QueryRowContext(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(res, &m); err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "query_one", "failed to unmarshal result")
	}

	for _, v := range m {
		return v, nil
	}

	return nil, nil
}

// QueryRow executes a query and returns a single row
func (c *Client) QueryRow(query string, parameters interface{}) ([]byte, error) {
	return c.QueryRowContext(context.Background(), query, parameters)
}

// QueryRowContext executes a query with context and returns a single row
func (c *Client) QueryRowContext(ctx context.Context, query string, parameters interface{}) ([]byte, error) {
	res, err := c.QueryContext(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if len(res) < 2 {
		return []byte("{}"), nil
	}

	return res[1 : len(res)-1], nil
}

// QueryColumn executes a query and returns values from a single column
func (c *Client) QueryColumn(query string, parameters interface{}) ([]interface{}, error) {
	return c.QueryColumnContext(context.Background(), query, parameters)
}

// QueryColumnContext executes a query with context and returns values from a single column
func (c *Client) QueryColumnContext(ctx context.Context, query string, parameters interface{}) ([]interface{}, error) {
	res, err := c.QueryContext(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(res, &rows); err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "query_column", "failed to unmarshal result")
	}

	var values []interface{}
	for _, m := range rows {
		for _, v := range m {
			values = append(values, v)
			break // Only take the first column
		}
	}

	return values, nil
}

// Create creates a new entity in SolarWinds
func (c *Client) Create(entity, body interface{}) ([]byte, error) {
	return c.CreateContext(context.Background(), entity, body)
}

// CreateContext creates a new entity with context
func (c *Client) CreateContext(ctx context.Context, entity, body interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("Create/%s", entity)
	return c.PostContext(ctx, endpoint, body)
}

// Read retrieves an entity by URI
func (c *Client) Read(uri string) ([]byte, error) {
	return c.ReadContext(context.Background(), uri)
}

// ReadContext retrieves an entity by URI with context
func (c *Client) ReadContext(ctx context.Context, uri string) ([]byte, error) {
	return c.GetContext(ctx, uri)
}

// Invoke executes a SolarWinds verb on an entity
func (c *Client) Invoke(entity, verb string, body interface{}) ([]byte, error) {
	return c.InvokeContext(context.Background(), entity, verb, body)
}

// InvokeContext executes a SolarWinds verb with context
func (c *Client) InvokeContext(ctx context.Context, entity, verb string, body interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("Invoke/%s/%s", entity, verb)
	return c.PostContext(ctx, endpoint, body)
}

// BulkDelete deletes multiple entities by URI
func (c *Client) BulkDelete(uris []string) ([]byte, error) {
	return c.BulkDeleteContext(context.Background(), uris)
}

// BulkDeleteContext deletes multiple entities with context
func (c *Client) BulkDeleteContext(ctx context.Context, uris []string) ([]byte, error) {
	req := map[string][]string{
		"uris": uris,
	}
	return c.PostContext(ctx, "BulkDelete", req)
}

// Delete removes an entity by URI
func (c *Client) Delete(uri string) ([]byte, error) {
	return c.DeleteContext(context.Background(), uri)
}

// DeleteContext removes an entity by URI with context
func (c *Client) DeleteContext(ctx context.Context, uri string) ([]byte, error) {
	return c.doRequest(ctx, "DELETE", uri, nil)
}

// Update modifies an existing entity
func (c *Client) Update(uri string, body map[string]interface{}) ([]byte, error) {
	return c.UpdateContext(context.Background(), uri, body)
}

// UpdateContext modifies an existing entity with context
func (c *Client) UpdateContext(ctx context.Context, uri string, body map[string]interface{}) ([]byte, error) {
	return c.PostContext(ctx, uri, body)
}
