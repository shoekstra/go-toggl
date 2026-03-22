package toggl

import (
	"net/http"
	"time"
)

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client for API requests.
// A nil value is ignored and the default client is retained.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithTimeout sets the timeout for API requests. It can be used alongside
// WithHTTPClient in any order; the timeout is always applied after all options
// are processed. Passing zero explicitly disables the timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
		c.timeoutSet = true
	}
}
