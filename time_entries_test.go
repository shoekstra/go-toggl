package toggl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testClient creates a test client backed by an httptest.Server running handler.
// The server is closed automatically when the test finishes.
func testClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := NewClient("test-token", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

const timeEntryJSON = `{
	"id": 123,
	"description": "Meeting",
	"project_id": 456,
	"task_id": null,
	"client_id": null,
	"workspace_id": 1,
	"user_id": 99,
	"billable": false,
	"tag_ids": [1, 2],
	"tags": ["client", "meeting"],
	"start": "2024-01-15T09:00:00Z",
	"stop": "2024-01-15T10:00:00Z",
	"duration": 3600,
	"created_with": "go-toggl",
	"at": "2024-01-15T10:00:00Z"
}`

// assertBody reads and unmarshals the request body into a map, calling t.Fatal on failure.
func assertBody(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read request body: %v", err)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	return body
}

func TestTimeEntriesService_ListTimeEntries(t *testing.T) {
	wantPath := "/api/v9/me/time_entries"

	tests := []struct {
		name       string
		opts       *ListTimeEntriesOptions
		statusCode int
		response   string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   "[" + timeEntryJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name: "success with options",
			opts: &ListTimeEntriesOptions{
				StartDate: String("2024-01-01"),
				EndDate:   String("2024-01-31"),
			},
			statusCode: http.StatusOK,
			response:   "[" + timeEntryJSON + "," + timeEntryJSON + "]",
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "empty result",
			opts:       nil,
			statusCode: http.StatusOK,
			response:   "[]",
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantCount:  0,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantCount:  0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entries, _, err := client.TimeEntries.ListTimeEntries(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ListTimeEntries() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(entries) != tt.wantCount {
				t.Errorf("ListTimeEntries() count = %d, want %d", len(entries), tt.wantCount)
			}
		})
	}
}

func TestTimeEntriesService_ListTimeEntries_QueryParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("start_date"); got != "2024-01-01" {
			t.Errorf("start_date = %q, want %q", got, "2024-01-01")
		}
		if got := q.Get("end_date"); got != "2024-01-31" {
			t.Errorf("end_date = %q, want %q", got, "2024-01-31")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	})

	client := testClient(t, handler)
	_, _, err := client.TimeEntries.ListTimeEntries(context.Background(), &ListTimeEntriesOptions{
		StartDate: String("2024-01-01"),
		EndDate:   String("2024-01-31"),
	})
	if err != nil {
		t.Fatalf("ListTimeEntries() error = %v", err)
	}
}

func TestTimeEntriesService_GetTimeEntry(t *testing.T) {
	tests := []struct {
		name       string
		entryID    int
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			entryID:    123,
			statusCode: http.StatusOK,
			response:   timeEntryJSON,
			wantID:     123,
			wantErr:    false,
		},
		{
			name:       "not found",
			entryID:    999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantID:     0,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			entryID:    123,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantID:     0,
			wantErr:    true,
		},
		{
			name:       "server error",
			entryID:    123,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantID:     0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/me/time_entries/%d", tt.entryID)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entry, _, err := client.TimeEntries.GetTimeEntry(context.Background(), tt.entryID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetTimeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && entry.ID != tt.wantID {
				t.Errorf("GetTimeEntry() ID = %d, want %d", entry.ID, tt.wantID)
			}
		})
	}
}

func TestTimeEntriesService_GetRunningTimeEntry(t *testing.T) {
	wantPath := "/api/v9/me/time_entries/current"

	tests := []struct {
		name       string
		statusCode int
		response   string
		wantID     int
		wantNil    bool
		wantErr    bool
	}{
		{
			name:       "running entry exists",
			statusCode: http.StatusOK,
			response:   timeEntryJSON,
			wantID:     123,
		},
		{
			name:       "no running entry — API returns null/200",
			statusCode: http.StatusOK,
			response:   "null",
			wantNil:    true,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entry, _, err := client.TimeEntries.GetRunningTimeEntry(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetRunningTimeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && entry != nil {
				t.Errorf("GetRunningTimeEntry() = %+v, want nil", entry)
			}
			if !tt.wantErr && !tt.wantNil && entry.ID != tt.wantID {
				t.Errorf("GetRunningTimeEntry() ID = %d, want %d", entry.ID, tt.wantID)
			}
		})
	}
}

func TestTimeEntriesService_CreateTimeEntry(t *testing.T) {
	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		workspaceID int
		opts        *CreateTimeEntryOptions
		statusCode  int
		response    string
		wantID      int
		wantErr     bool
	}{
		{
			name:        "success",
			workspaceID: 1,
			opts: &CreateTimeEntryOptions{
				Start:       start,
				Description: String("Meeting"),
				ProjectID:   Int(456),
			},
			statusCode: http.StatusOK,
			response:   timeEntryJSON,
			wantID:     123,
			wantErr:    false,
		},
		{
			name:        "success with all options",
			workspaceID: 1,
			opts: &CreateTimeEntryOptions{
				Start:       start,
				Description: String("Deep work"),
				ProjectID:   Int(456),
				TaskID:      Int(789),
				Billable:    Bool(true),
				Tags:        []string{"focus"},
				TagIDs:      []int{1, 2},
				Duration:    Int(7200),
				CreatedWith: "my-app",
			},
			statusCode: http.StatusOK,
			response:   timeEntryJSON,
			wantID:     123,
			wantErr:    false,
		},
		{
			name:        "nil options",
			workspaceID: 1,
			opts:        nil,
			wantErr:     true,
		},
		{
			name:        "zero start time",
			workspaceID: 1,
			opts:        &CreateTimeEntryOptions{},
			wantErr:     true,
		},
		{
			name:        "stop before start",
			workspaceID: 1,
			opts: &CreateTimeEntryOptions{
				Start: start,
				Stop:  func() *time.Time { t := start.Add(-time.Hour); return &t }(),
			},
			wantErr: true,
		},
		{
			name:        "stop equals start",
			workspaceID: 1,
			opts: &CreateTimeEntryOptions{
				Start: start,
				Stop:  &start,
			},
			wantErr: true,
		},
		{
			name:        "stop derives duration",
			workspaceID: 1,
			opts: &CreateTimeEntryOptions{
				Start: start,
				Stop:  func() *time.Time { t := start.Add(time.Hour); return &t }(),
			},
			statusCode: http.StatusOK,
			response:   timeEntryJSON,
			wantID:     123,
			wantErr:    false,
		},
		{
			name:        "unauthorized",
			workspaceID: 1,
			opts:        &CreateTimeEntryOptions{Start: start},
			statusCode:  http.StatusUnauthorized,
			response:    `{"error": "unauthorized"}`,
			wantID:      0,
			wantErr:     true,
		},
		{
			name:        "server error",
			workspaceID: 1,
			opts:        &CreateTimeEntryOptions{Start: start},
			statusCode:  http.StatusInternalServerError,
			response:    `{"error": "internal server error"}`,
			wantID:      0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/%d/time_entries", tt.workspaceID)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				body := assertBody(t, r)
				if wid, ok := body["workspace_id"].(float64); !ok || int(wid) != tt.workspaceID {
					t.Errorf("workspace_id = %v, want %d", body["workspace_id"], tt.workspaceID)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				dur, ok := body["duration"].(float64)
				if !ok {
					t.Errorf("body missing required field duration")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// "stop derives duration" case: Stop-Start == 1h → 3600s.
				if tt.opts != nil && tt.opts.Stop != nil && tt.opts.Duration == nil {
					if int(dur) != 3600 {
						t.Errorf("derived duration = %d, want 3600", int(dur))
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entry, _, err := client.TimeEntries.CreateTimeEntry(context.Background(), tt.workspaceID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateTimeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && entry.ID != tt.wantID {
				t.Errorf("CreateTimeEntry() ID = %d, want %d", entry.ID, tt.wantID)
			}
		})
	}
}

func TestTimeEntriesService_StartTimeEntry(t *testing.T) {
	start := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	wantPath := "/api/v9/workspaces/1/time_entries"

	t.Run("success", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if r.URL.Path != wantPath {
				t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			body := assertBody(t, r)
			if dur, ok := body["duration"].(float64); !ok || int(dur) != -1 {
				t.Errorf("duration = %v, want -1", body["duration"])
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if _, hasStop := body["stop"]; hasStop {
				t.Errorf("stop field must not be set for a running entry")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(timeEntryJSON))
		})

		client := testClient(t, handler)
		entry, _, err := client.TimeEntries.StartTimeEntry(context.Background(), 1, &CreateTimeEntryOptions{
			Start:       start,
			Description: String("Focus session"),
		})
		if err != nil {
			t.Fatalf("StartTimeEntry() error = %v", err)
		}
		if entry.ID != 123 {
			t.Errorf("StartTimeEntry() ID = %d, want 123", entry.ID)
		}
	})

	t.Run("nil options", func(t *testing.T) {
		client := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		_, _, err := client.TimeEntries.StartTimeEntry(context.Background(), 1, nil)
		if err == nil {
			t.Error("StartTimeEntry() expected error for nil options, got nil")
		}
	})

	t.Run("does not mutate caller opts", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(timeEntryJSON))
		})

		client := testClient(t, handler)
		stop := start.Add(time.Hour)
		opts := &CreateTimeEntryOptions{
			Start:    start,
			Stop:     &stop,
			Duration: Int(3600),
		}
		_, _, _ = client.TimeEntries.StartTimeEntry(context.Background(), 1, opts)

		// Caller's opts must not be modified.
		if opts.Duration == nil || *opts.Duration != 3600 {
			t.Error("StartTimeEntry() mutated caller's Duration")
		}
		if opts.Stop == nil {
			t.Error("StartTimeEntry() mutated caller's Stop")
		}
	})
}

func TestTimeEntriesService_UpdateTimeEntry(t *testing.T) {
	tests := []struct {
		name        string
		workspaceID int
		entryID     int
		opts        *UpdateTimeEntryOptions
		statusCode  int
		response    string
		wantID      int
		wantErr     bool
	}{
		{
			name:        "success",
			workspaceID: 1,
			entryID:     123,
			opts: &UpdateTimeEntryOptions{
				Description: String("Updated meeting"),
				Billable:    Bool(true),
			},
			statusCode: http.StatusOK,
			response:   timeEntryJSON,
			wantID:     123,
			wantErr:    false,
		},
		{
			name:        "success with all options",
			workspaceID: 1,
			entryID:     123,
			opts: &UpdateTimeEntryOptions{
				Description: String("Updated"),
				ProjectID:   Int(789),
				TaskID:      Int(1),
				Billable:    Bool(false),
				Tags:        []string{"updated"},
				TagIDs:      []int{3},
				TagAction:   String("add"),
				Duration:    Int(7200),
			},
			statusCode: http.StatusOK,
			response:   timeEntryJSON,
			wantID:     123,
			wantErr:    false,
		},
		{
			name:        "nil options",
			workspaceID: 1,
			entryID:     123,
			opts:        nil,
			wantErr:     true,
		},
		{
			name:        "not found",
			workspaceID: 1,
			entryID:     999,
			opts:        &UpdateTimeEntryOptions{Description: String("x")},
			statusCode:  http.StatusNotFound,
			response:    `{"error": "not found"}`,
			wantErr:     true,
		},
		{
			name:        "unauthorized",
			workspaceID: 1,
			entryID:     123,
			opts:        &UpdateTimeEntryOptions{Description: String("x")},
			statusCode:  http.StatusUnauthorized,
			response:    `{"error": "unauthorized"}`,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/%d/time_entries/%d", tt.workspaceID, tt.entryID)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("expected PUT, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				body := assertBody(t, r)
				if wid, ok := body["workspace_id"].(float64); !ok || int(wid) != tt.workspaceID {
					t.Errorf("workspace_id = %v, want %d", body["workspace_id"], tt.workspaceID)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entry, _, err := client.TimeEntries.UpdateTimeEntry(context.Background(), tt.workspaceID, tt.entryID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateTimeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && entry.ID != tt.wantID {
				t.Errorf("UpdateTimeEntry() ID = %d, want %d", entry.ID, tt.wantID)
			}
		})
	}
}

func TestTimeEntriesService_DeleteTimeEntry(t *testing.T) {
	tests := []struct {
		name        string
		workspaceID int
		entryID     int
		statusCode  int
		response    string
		wantErr     bool
	}{
		{
			name:        "success",
			workspaceID: 1,
			entryID:     123,
			statusCode:  http.StatusOK,
			response:    `"OK"`,
			wantErr:     false,
		},
		{
			name:        "not found",
			workspaceID: 1,
			entryID:     999,
			statusCode:  http.StatusNotFound,
			response:    `{"error": "not found"}`,
			wantErr:     true,
		},
		{
			name:        "unauthorized",
			workspaceID: 1,
			entryID:     123,
			statusCode:  http.StatusUnauthorized,
			response:    `{"error": "unauthorized"}`,
			wantErr:     true,
		},
		{
			name:        "server error",
			workspaceID: 1,
			entryID:     123,
			statusCode:  http.StatusInternalServerError,
			response:    `{"error": "internal server error"}`,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/%d/time_entries/%d", tt.workspaceID, tt.entryID)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			_, err := client.TimeEntries.DeleteTimeEntry(context.Background(), tt.workspaceID, tt.entryID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteTimeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeEntriesService_StopTimeEntry(t *testing.T) {
	tests := []struct {
		name        string
		workspaceID int
		entryID     int
		statusCode  int
		response    string
		wantID      int
		wantErr     bool
	}{
		{
			name:        "success",
			workspaceID: 1,
			entryID:     123,
			statusCode:  http.StatusOK,
			response:    timeEntryJSON,
			wantID:      123,
			wantErr:     false,
		},
		{
			name:        "not found",
			workspaceID: 1,
			entryID:     999,
			statusCode:  http.StatusNotFound,
			response:    `{"error": "not found"}`,
			wantErr:     true,
		},
		{
			name:        "unauthorized",
			workspaceID: 1,
			entryID:     123,
			statusCode:  http.StatusUnauthorized,
			response:    `{"error": "unauthorized"}`,
			wantErr:     true,
		},
		{
			name:        "server error",
			workspaceID: 1,
			entryID:     123,
			statusCode:  http.StatusInternalServerError,
			response:    `{"error": "internal server error"}`,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/%d/time_entries/%d/stop", tt.workspaceID, tt.entryID)
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPatch {
					t.Errorf("expected PATCH, got %s", r.Method)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.URL.Path != wantPath {
					t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			entry, _, err := client.TimeEntries.StopTimeEntry(context.Background(), tt.workspaceID, tt.entryID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("StopTimeEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && entry.ID != tt.wantID {
				t.Errorf("StopTimeEntry() ID = %d, want %d", entry.ID, tt.wantID)
			}
		})
	}
}
