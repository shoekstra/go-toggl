package toggl

import (
	"fmt"
	"net/http"
	"time"
)

const (
	defaultBaseURL    = "https://api.track.toggl.com"
	defaultTimeout    = 30 * time.Second
	userAgent         = "go-toggl"
	basicAuthPassword = "api_token"
)

// Client is the Toggl Track API client.
type Client struct {
	baseURL     string
	token       string
	httpClient  *http.Client
	timeout     time.Duration
	timeoutSet  bool

	// Services
	TimeEntries *TimeEntriesService
	Projects    *ProjectsService
	Tags        *TagsService
	Clients     *ClientsService
	Workspaces  *WorkspacesService
	Reports     *ReportsService
}

// NewClient creates a new Toggl API client with the given token.
//
// Use the functional options (WithBaseURL, WithHTTPClient, WithTimeout) to
// customise the client. When WithTimeout is used it takes effect regardless of
// whether WithHTTPClient appears before or after it. When only WithHTTPClient
// is provided, the custom client's existing Timeout is preserved; if the custom
// client has no timeout set, defaultTimeout is applied.
func NewClient(token string, opts ...ClientOption) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	c := &Client{
		baseURL:    defaultBaseURL,
		token:      token,
		httpClient: &http.Client{},
		// timeout is zero until WithTimeout sets it explicitly.
	}

	// Apply options.
	for _, opt := range opts {
		opt(c)
	}

	// Resolve the final timeout:
	//   - WithTimeout was called: always use that value (even zero).
	//   - WithHTTPClient was called without WithTimeout: preserve the custom
	//     client's timeout if it has one, otherwise apply defaultTimeout.
	if c.timeoutSet {
		c.httpClient.Timeout = c.timeout
	} else if c.httpClient.Timeout == 0 {
		c.httpClient.Timeout = defaultTimeout
	}

	// Initialize services.
	c.TimeEntries = &TimeEntriesService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Tags = &TagsService{client: c}
	c.Clients = &ClientsService{client: c}
	c.Workspaces = &WorkspacesService{client: c}
	c.Reports = &ReportsService{client: c}

	return c, nil
}
