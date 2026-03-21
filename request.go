package toggl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response wraps the standard http.Response with additional context.
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Response   *http.Response
}

// newRequest creates a new HTTP request with common headers.
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.SetBasicAuth(c.token, basicAuthPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// do executes an HTTP request and returns the response.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       respBody,
		Response:   resp,
	}

	// Check for HTTP errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return response, &ErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
			Response:   resp,
		}
	}

	// Decode response if provided
	if v != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, v); err != nil {
			return response, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return response, nil
}

// get performs a GET request.
func (c *Client) get(ctx context.Context, path string, v interface{}) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, v)
}

// post performs a POST request.
func (c *Client) post(ctx context.Context, path string, body, v interface{}) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, v)
}

// put performs a PUT request.
func (c *Client) put(ctx context.Context, path string, body, v interface{}) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, v)
}

// delete performs a DELETE request.
func (c *Client) delete(ctx context.Context, path string) (*Response, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req, nil)
}
