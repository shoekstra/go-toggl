package toggl

import (
	"strings"
	"testing"
)

func TestErrorResponse_Error(t *testing.T) {
	err := &ErrorResponse{StatusCode: 404, Message: "resource not found"}
	want := "toggl: 404: resource not found"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestParseErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "message key",
			body: `{"message": "Resource can not be found", "code": 404}`,
			want: "Resource can not be found",
		},
		{
			name: "error key",
			body: `{"error": "unauthorized"}`,
			want: "unauthorized",
		},
		{
			name: "message takes precedence over error",
			body: `{"message": "detail here", "error": "short"}`,
			want: "detail here",
		},
		{
			name: "plain string body",
			body: `Not Found`,
			want: "Not Found",
		},
		{
			name: "empty body",
			body: "",
			want: "empty response",
		},
		{
			name: "whitespace-only body",
			body: "   \n\t  ",
			want: "empty response",
		},
		{
			name: "whitespace-only message field",
			body: `{"message": "   "}`,
			want: `{"message": "   "}`,
		},
		{
			name: "whitespace-only error field",
			body: `{"error": "\t\n"}`,
			want: `{"error": "\t\n"}`,
		},
		{
			name: "empty JSON object",
			body: `{}`,
			want: "{}",
		},
		{
			name: "body truncated at 200 characters",
			body: strings.Repeat("x", 300),
			want: strings.Repeat("x", 200) + "...",
		},
		{
			name: "HTML error page truncated",
			body: "<html><body>" + strings.Repeat("e", 300) + "</body></html>",
			want: "<html><body>" + strings.Repeat("e", 188) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseErrorMessage([]byte(tt.body))
			if got != tt.want {
				t.Errorf("parseErrorMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}
