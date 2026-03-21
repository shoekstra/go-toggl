package toggl

import (
	"fmt"
	"net/http"
	"time"
)

const (
	defaultBaseURL     = "https://api.track.toggl.com"
	defaultTimeout     = 30 * time.Second
	userAgent          = "go-toggl/1.0.0"
	basicAuthPassword  = "api_token"
)

// Client is the Toggl Track API client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client

	// Services
	TimeEntries *TimeEntriesService
	Projects    *ProjectsService
	Tags        *TagsService
	Clients     *ClientsService
	Workspaces  *WorkspacesService
	Reports     *ReportsService
}

// NewClient creates a new Toggl API client with the given token.
func NewClient(token string, opts ...ClientOption) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	c := &Client{
		baseURL: defaultBaseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.TimeEntries = &TimeEntriesService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Tags = &TagsService{client: c}
	c.Clients = &ClientsService{client: c}
	c.Workspaces = &WorkspacesService{client: c}
	c.Reports = &ReportsService{client: c}

	return c, nil
}

// BaseURL returns the base URL of the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// SetBaseURL sets the base URL of the client.
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// HTTPClient returns the HTTP client used by the Toggl client.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// SetHTTPClient sets the HTTP client used by the Toggl client.
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}
