package toggl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPagination_Headers(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    Pagination
		wantErr bool
	}{
		{
			name: "all page-based headers present",
			headers: map[string]string{
				"X-Page":      "3",
				"X-Pages":     "10",
				"X-Page-Size": "50",
			},
			want: Pagination{CurrentPage: 3, TotalPages: 10, PageSize: 50},
		},
		{
			name: "cursor-based headers present",
			headers: map[string]string{
				"X-Next-ID":         "9876",
				"X-Next-Row-Number": "150",
			},
			want: Pagination{NextID: 9876, NextRowNumber: 150},
		},
		{
			name:    "no pagination headers",
			headers: map[string]string{},
			want:    Pagination{},
		},
		{
			name:    "malformed X-Page returns error",
			headers: map[string]string{"X-Page": "bad"},
			wantErr: true,
		},
		{
			name:    "malformed X-Pages returns error",
			headers: map[string]string{"X-Pages": "bad"},
			wantErr: true,
		},
		{
			name:    "malformed X-Page-Size returns error",
			headers: map[string]string{"X-Page-Size": "bad"},
			wantErr: true,
		},
		{
			name:    "malformed X-Next-ID returns error",
			headers: map[string]string{"X-Next-ID": "bad"},
			wantErr: true,
		},
		{
			name:    "malformed X-Next-Row-Number returns error",
			headers: map[string]string{"X-Next-Row-Number": "bad"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for k, v := range tt.headers {
					w.Header().Set(k, v)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{}"))
			})

			server := httptest.NewServer(handler)
			defer server.Close()

			client, err := NewClient("test-token", WithBaseURL(server.URL))
			if err != nil {
				t.Fatalf("NewClient: %v", err)
			}

			resp, err := client.get(context.Background(), "/", nil)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error for malformed header, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("get: %v", err)
			}

			got := resp.Pagination
			if got != tt.want {
				t.Errorf("Pagination = %+v, want %+v", got, tt.want)
			}
		})
	}
}
