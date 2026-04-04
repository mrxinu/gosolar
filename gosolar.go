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

// Client represents a SolarWinds SWIS API client.
type Client struct {
	config     *Config
	baseURL    *url.URL
	httpClient *http.Client
	logger     *slog.Logger
}

// NewClient creates a new SolarWinds client with the provided configuration.
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(fmt.Sprintf("https://%s:%d/SolarWinds/InformationService/v3/Json/", config.Host, config.Port))
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

func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any) ([]byte, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, WrapError(err, ErrorTypeValidation, "request", "failed to marshal request body")
		}
	}

	endpointURL, err := c.baseURL.Parse(endpoint)
	if err != nil {
		return nil, WrapError(err, ErrorTypeValidation, "request", "invalid endpoint")
	}

	c.logger.DebugContext(ctx, "making request", "method", method, "endpoint", endpoint)

	var lastErr error

	for attempt := range c.config.MaxRetries + 1 {
		if attempt > 0 {
			delay := c.config.RetryBaseDelay * time.Duration(1<<(attempt-1))
			c.logger.DebugContext(ctx, "retrying request", "attempt", attempt, "delay", delay)

			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, WrapError(ctx.Err(), ErrorTypeNetwork, "request", "request cancelled during retry backoff")
			case <-timer.C:
			}
		}

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, endpointURL.String(), reqBody)
		if err != nil {
			return nil, WrapError(err, ErrorTypeNetwork, "request", "failed to create request")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", c.config.UserAgent)
		req.SetBasicAuth(c.config.Username, c.config.Password)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		output, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, WrapError(err, ErrorTypeNetwork, "response", "failed to read response body")
		}

		if resp.StatusCode >= 400 {
			return nil, NewHTTPError("request", endpoint, resp.StatusCode, string(output))
		}

		c.logger.DebugContext(ctx, "request completed", "status", resp.StatusCode)
		return output, nil
	}

	return nil, WrapError(lastErr, ErrorTypeNetwork, "request", "request failed after retries")
}

// Post performs a POST request with context.
func (c *Client) Post(ctx context.Context, endpoint string, body any) ([]byte, error) {
	return c.doRequest(ctx, "POST", endpoint, body)
}

// Get performs a GET request with context.
func (c *Client) Get(ctx context.Context, endpoint string) ([]byte, error) {
	return c.doRequest(ctx, "GET", endpoint, nil)
}

// Query executes a SWQL query with context and returns the results.
func (c *Client) Query(ctx context.Context, query string, parameters any) ([]byte, error) {
	req := struct {
		Query      string `json:"query"`
		Parameters any    `json:"parameters"`
	}{
		Query:      query,
		Parameters: parameters,
	}

	result, err := c.Post(ctx, "Query", &req)
	if err != nil {
		var swErr *Error
		if ok := isError(err, &swErr); ok {
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

// QueryOne executes a query with context and returns a single value.
func (c *Client) QueryOne(ctx context.Context, query string, parameters any) (any, error) {
	res, err := c.QueryRow(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	m := make(map[string]any)
	if err := json.Unmarshal(res, &m); err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "query_one", "failed to unmarshal result")
	}

	for _, v := range m {
		return v, nil
	}

	return nil, nil
}

// QueryRow executes a query with context and returns a single row.
func (c *Client) QueryRow(ctx context.Context, query string, parameters any) ([]byte, error) {
	res, err := c.Query(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	if len(res) < 2 {
		return []byte("{}"), nil
	}

	return res[1 : len(res)-1], nil
}

// QueryColumn executes a query with context and returns values from a single column.
func (c *Client) QueryColumn(ctx context.Context, query string, parameters any) ([]any, error) {
	res, err := c.Query(ctx, query, parameters)
	if err != nil {
		return nil, err
	}

	var rows []map[string]any
	if err := json.Unmarshal(res, &rows); err != nil {
		return nil, WrapError(err, ErrorTypeInternal, "query_column", "failed to unmarshal result")
	}

	values := make([]any, 0, len(rows))
	for _, m := range rows {
		for _, v := range m {
			values = append(values, v)
			break
		}
	}

	return values, nil
}

// Create creates a new entity in SolarWinds with context.
func (c *Client) Create(ctx context.Context, entity string, body any) ([]byte, error) {
	endpoint := fmt.Sprintf("Create/%s", entity)
	return c.Post(ctx, endpoint, body)
}

// Read retrieves an entity by URI with context.
func (c *Client) Read(ctx context.Context, uri string) ([]byte, error) {
	return c.Get(ctx, uri)
}

// Invoke executes a SolarWinds verb on an entity with context.
func (c *Client) Invoke(ctx context.Context, entity, verb string, body any) ([]byte, error) {
	endpoint := fmt.Sprintf("Invoke/%s/%s", entity, verb)
	return c.Post(ctx, endpoint, body)
}

// BulkDelete deletes multiple entities by URI with context.
func (c *Client) BulkDelete(ctx context.Context, uris []string) ([]byte, error) {
	req := map[string][]string{
		"uris": uris,
	}
	return c.Post(ctx, "BulkDelete", req)
}

// Delete removes an entity by URI with context.
func (c *Client) Delete(ctx context.Context, uri string) ([]byte, error) {
	return c.doRequest(ctx, "DELETE", uri, nil)
}

// Update modifies an existing entity with context.
func (c *Client) Update(ctx context.Context, uri string, body map[string]any) ([]byte, error) {
	return c.Post(ctx, uri, body)
}

// isError checks if an error is a structured Error type.
func isError(err error, target **Error) bool {
	if swErr, ok := err.(*Error); ok {
		*target = swErr
		return true
	}
	return false
}
