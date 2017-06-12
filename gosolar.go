package gosolar

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	ioutil "io/ioutil"
	"log"
	"net/http"
)

// Client structure for the SolarWinds (SWIS) connection.
type Client struct {
	// connection parameters
	URL      string
	Username string
	Password string

	// internal state
	http *http.Client
}

// NewClient creates a new reference to the Client struct.
func NewClient(host, user, pass string, ignoreSSL bool) *Client {
	return &Client{
		URL:      fmt.Sprintf("https://%s:17778/SolarWinds/InformationService/v3/Json/", host),
		Username: user,
		Password: pass,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: ignoreSSL,
				},
				MaxIdleConnsPerHost: 4, // DEFAULT: 2
			},
		},
	}
}

func (c *Client) post(endpoint string, body interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", c.URL+endpoint, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to submit query: %v", err)
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		err2 := res.Body.Close()
		if err2 != nil {
			log.Fatalf("failed to close result body: %v", err2)
		}
	}()

	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("query failed - status code %d: %v", res.StatusCode, err)
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("swis failure message [status: %d]:\n%s",
			res.StatusCode, string(output))
	}

	return output, nil
}

// Query retrieves a result from the SolarWinds API.
func (c *Client) Query(query string, parameters map[string]string) ([]byte, error) {
	req := struct {
		Query      string            `json:"query"`
		Parameters map[string]string `json:"parameters"`
	}{
		Query:      query,
		Parameters: parameters,
	}

	result, err := c.post("Query", &req)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	type swisResult struct {
		Result *json.RawMessage `json:"results"`
	}

	// unpack the result and return it
	var sr swisResult
	err = json.Unmarshal(result, &sr)
	if err != nil {
		return nil, err
	}

	return []byte(*sr.Result), nil
}
