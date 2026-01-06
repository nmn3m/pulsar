package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// TestClient wraps http.Client with test utilities
type TestClient struct {
	baseURL    string
	httpClient *http.Client
	t          *testing.T
	authToken  string
}

// NewTestClient creates a new test HTTP client
func NewTestClient(t *testing.T, baseURL string) *TestClient {
	return &TestClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		t:          t,
	}
}

// SetAuthToken sets the Bearer token for subsequent requests
func (c *TestClient) SetAuthToken(token string) {
	c.authToken = token
}

// ClearAuthToken removes the auth token
func (c *TestClient) ClearAuthToken() {
	c.authToken = ""
}

// Get performs a GET request
func (c *TestClient) Get(path string) *http.Response {
	return c.doRequest("GET", path, nil)
}

// GetWithQuery performs a GET request with query parameters
func (c *TestClient) GetWithQuery(path string, query map[string]string) *http.Response {
	if len(query) > 0 {
		path += "?"
		first := true
		for k, v := range query {
			if !first {
				path += "&"
			}
			path += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}
	return c.doRequest("GET", path, nil)
}

// Post performs a POST request
func (c *TestClient) Post(path string, body interface{}) *http.Response {
	return c.doRequest("POST", path, body)
}

// Patch performs a PATCH request
func (c *TestClient) Patch(path string, body interface{}) *http.Response {
	return c.doRequest("PATCH", path, body)
}

// Put performs a PUT request
func (c *TestClient) Put(path string, body interface{}) *http.Response {
	return c.doRequest("PUT", path, body)
}

// Delete performs a DELETE request
func (c *TestClient) Delete(path string) *http.Response {
	return c.doRequest("DELETE", path, nil)
}

// doRequest performs an HTTP request
func (c *TestClient) doRequest(method, path string, body interface{}) *http.Response {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			c.t.Fatalf("Failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		c.t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.t.Fatalf("Failed to execute request: %v", err)
	}

	return resp
}

// ParseJSON parses the response body as JSON
func (c *TestClient) ParseJSON(resp *http.Response, v interface{}) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.t.Fatalf("Failed to read response body: %v", err)
	}
	if err := json.Unmarshal(body, v); err != nil {
		c.t.Fatalf("Failed to parse JSON response: %v\nBody: %s", err, string(body))
	}
}

// ReadBody reads the response body as string
func (c *TestClient) ReadBody(resp *http.Response) string {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.t.Fatalf("Failed to read response body: %v", err)
	}
	return string(body)
}

// ExpectStatus asserts the response status code
func (c *TestClient) ExpectStatus(resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		c.t.Errorf("Expected status %d, got %d. Body: %s", expected, resp.StatusCode, string(body))
	}
}

// AssertStatus asserts the response status code and fails immediately if wrong
func (c *TestClient) AssertStatus(resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		c.t.Fatalf("Expected status %d, got %d. Body: %s", expected, resp.StatusCode, string(body))
	}
}
