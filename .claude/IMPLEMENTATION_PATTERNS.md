# go-toggl Implementation Patterns

## Service Method Template

Use this template for all service methods:

### List Method

```go
// ListTimeEntriesOptions contains optional parameters for ListTimeEntries.
type ListTimeEntriesOptions struct {
    StartDate *time.Time
    EndDate   *time.Time
    Billable  *bool
    ProjectID *int
    Limit     *int
    Page      *int
}

// ListTimeEntries lists all time entries with optional filters.
//
// API: GET /me/time_entries
func (s *TimeEntriesService) ListTimeEntries(ctx context.Context, opts *ListTimeEntriesOptions) ([]*TimeEntry, *Response, error) {
    path := "/me/time_entries"

    // Build query parameters if options provided
    if opts != nil {
        // Use url.Values for query string building
    }

    var entries []*TimeEntry
    resp, err := s.client.get(ctx, path, &entries)
    if err != nil {
        return nil, resp, err
    }

    return entries, resp, nil
}
```

### Get (Single Resource) Method

```go
// GetTimeEntry gets a single time entry by ID.
//
// API: GET /me/time_entries/{id}
func (s *TimeEntriesService) GetTimeEntry(ctx context.Context, id int) (*TimeEntry, *Response, error) {
    path := fmt.Sprintf("/me/time_entries/%d", id)

    entry := new(TimeEntry)
    resp, err := s.client.get(ctx, path, entry)
    if err != nil {
        return nil, resp, err
    }

    return entry, resp, nil
}
```

### Create Method

```go
// CreateTimeEntryOptions contains options for creating a time entry.
type CreateTimeEntryOptions struct {
    Name        string
    Description *string
    ProjectID   *int
    Tags        []string
    Billable    *bool
    Start       time.Time
    Stop        *time.Time
    Duration    *int
}

// CreateTimeEntry creates a new time entry.
//
// API: POST /workspaces/{workspace_id}/time_entries
func (s *TimeEntriesService) CreateTimeEntry(ctx context.Context, workspaceID int, opts *CreateTimeEntryOptions) (*TimeEntry, *Response, error) {
    if opts == nil {
        return nil, nil, fmt.Errorf("options required")
    }

    path := fmt.Sprintf("/workspaces/%d/time_entries", workspaceID)

    body := map[string]interface{}{
        "name":  opts.Name,
        "start": opts.Start,
    }

    if opts.Description != nil {
        body["description"] = *opts.Description
    }
    if opts.ProjectID != nil {
        body["project_id"] = *opts.ProjectID
    }
    // ... add other optional fields

    entry := new(TimeEntry)
    resp, err := s.client.post(ctx, path, body, entry)
    if err != nil {
        return nil, resp, err
    }

    return entry, resp, nil
}
```

### Update Method

```go
// UpdateTimeEntryOptions contains options for updating a time entry.
type UpdateTimeEntryOptions struct {
    Name        *string
    Description *string
    ProjectID   *int
    Tags        []string
    Billable    *bool
    Stop        *time.Time
}

// UpdateTimeEntry updates an existing time entry.
//
// API: PUT /workspaces/{workspace_id}/time_entries/{id}
func (s *TimeEntriesService) UpdateTimeEntry(ctx context.Context, workspaceID, id int, opts *UpdateTimeEntryOptions) (*TimeEntry, *Response, error) {
    if opts == nil {
        return nil, nil, fmt.Errorf("options required")
    }

    path := fmt.Sprintf("/workspaces/%d/time_entries/%d", workspaceID, id)

    body := make(map[string]interface{})

    if opts.Name != nil {
        body["name"] = *opts.Name
    }
    if opts.Description != nil {
        body["description"] = *opts.Description
    }
    // ... add other optional fields

    entry := new(TimeEntry)
    resp, err := s.client.put(ctx, path, body, entry)
    if err != nil {
        return nil, resp, err
    }

    return entry, resp, nil
}
```

### Delete Method

```go
// DeleteTimeEntry deletes a time entry.
//
// API: DELETE /workspaces/{workspace_id}/time_entries/{id}
func (s *TimeEntriesService) DeleteTimeEntry(ctx context.Context, workspaceID, id int) (*Response, error) {
    path := fmt.Sprintf("/workspaces/%d/time_entries/%d", workspaceID, id)
    return s.client.delete(ctx, path)
}
```

## Test Template

```go
import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
)

// Helper to create test client with mocked server
func testClient(handler http.Handler) *Client {
    server := httptest.NewServer(handler)
    client, _ := NewClient("test-token", WithBaseURL(server.URL))
    return client
}

func TestTimeEntriesService_GetTimeEntry(t *testing.T) {
    tests := []struct {
        name       string
        id         int
        statusCode int
        response   string
        wantID     int
        wantErr    bool
    }{
        {
            name:       "success",
            id:         123,
            statusCode: 200,
            response:   `{"id":123,"name":"test","start":"2024-01-01T09:00:00Z","duration":3600,"workspace_id":1,"user_id":1,"created_at":"2024-01-01T09:00:00Z","updated_at":"2024-01-01T09:00:00Z"}`,
            wantID:     123,
            wantErr:    false,
        },
        {
            name:       "not found",
            id:         999,
            statusCode: 404,
            response:   `{"error":"not found"}`,
            wantID:     0,
            wantErr:    true,
        },
        {
            name:       "unauthorized",
            id:         123,
            statusCode: 401,
            response:   `{"error":"unauthorized"}`,
            wantID:     0,
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(tt.statusCode)
                w.Write([]byte(tt.response))
            })

            client := testClient(handler)
            entry, _, err := client.TimeEntries.GetTimeEntry(context.Background(), tt.id)

            if (err != nil) != tt.wantErr {
                t.Errorf("GetTimeEntry() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && entry.ID != tt.wantID {
                t.Errorf("GetTimeEntry() got ID %d, want %d", entry.ID, tt.wantID)
            }
        })
    }
}
```

## Type Definition Template

```go
// TimeEntry represents a time entry in Toggl.
type TimeEntry struct {
    ID          int        `json:"id"`
    Name        *string    `json:"name"`
    Description *string    `json:"description"`
    ProjectID   *int       `json:"project_id"`
    ClientID    *int       `json:"client_id"`
    Tags        []string   `json:"tags"`
    Billable    *bool      `json:"billable"`
    Start       time.Time  `json:"start"`
    Stop        *time.Time `json:"stop"`
    Duration    int        `json:"duration"`
    WorkspaceID int        `json:"workspace_id"`
    UserID      int        `json:"user_id"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

## Helper Functions

```go
// Pointer helpers for optional parameters

// String returns a pointer to a string.
func String(s string) *string {
    return &s
}

// Int returns a pointer to an int.
func Int(i int) *int {
    return &i
}

// Bool returns a pointer to a bool.
func Bool(b bool) *bool {
    return &b
}

// Time returns a pointer to a time.Time.
func Time(t time.Time) *time.Time {
    return &t
}
```

## Usage Pattern

```go
// Using options
entries, _, err := client.TimeEntries.ListTimeEntries(ctx, &toggl.ListTimeEntriesOptions{
    StartDate: toggl.Time(time.Now().AddDate(0, 0, -7)),
    Limit:     toggl.Int(50),
})

// Creating
entry, _, err := client.TimeEntries.CreateTimeEntry(ctx, workspaceID, &toggl.CreateTimeEntryOptions{
    Name:      "Client meeting",
    ProjectID: toggl.Int(projectID),
    Tags:      []string{"meeting", "client"},
    Start:     time.Now(),
})

// Updating
updated, _, err := client.TimeEntries.UpdateTimeEntry(ctx, workspaceID, entryID, &toggl.UpdateTimeEntryOptions{
    Name: toggl.String("Updated name"),
})

// Deleting
_, err := client.TimeEntries.DeleteTimeEntry(ctx, workspaceID, entryID)
```
