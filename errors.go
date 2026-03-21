package toggl

import "net/http"

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	StatusCode int
	Message    string
	Response   *http.Response
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	return e.Message
}
