package toggl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ErrorResponse represents an error returned by the Toggl API.
type ErrorResponse struct {
	// StatusCode is the HTTP status code returned by the server.
	StatusCode int
	// Message is a human-readable description of the error, parsed from the
	// response body where possible.
	Message string
	// Response is the raw HTTP response.
	Response *http.Response
}

// Error implements the error interface.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("toggl: %d: %s", e.StatusCode, e.Message)
}

// parseErrorMessage extracts a human-readable message from a Toggl API error
// body. It tries common JSON shapes ("message" and "error" keys) and falls
// back to the raw body, truncated to 200 characters to avoid embedding large
// HTML error pages in error strings.
func parseErrorMessage(body []byte) string {
	var payload struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		if s := strings.TrimSpace(payload.Message); s != "" {
			return s
		}
		if s := strings.TrimSpace(payload.Error); s != "" {
			return s
		}
	}

	const maxLen = 200
	s := strings.TrimSpace(string(body))
	if s == "" {
		return "empty response"
	}
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}
