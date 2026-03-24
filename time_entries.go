package toggl

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// TimeEntriesService handles operations related to time entries.
type TimeEntriesService struct {
	client *Client
}

// ListTimeEntriesOptions specifies the optional parameters to
// TimeEntriesService.ListTimeEntries.
type ListTimeEntriesOptions struct {
	// Since filters entries modified since this UNIX timestamp (includes deleted).
	Since *int64
	// Before filters entries with start time before this date (YYYY-MM-DD or RFC3339).
	Before *string
	// StartDate filters entries with start time from this date (YYYY-MM-DD or RFC3339).
	// Use together with EndDate.
	StartDate *string
	// EndDate filters entries with start time until this date (YYYY-MM-DD or RFC3339).
	// Use together with StartDate.
	EndDate *string
}

// CreateTimeEntryOptions specifies the parameters to
// TimeEntriesService.CreateTimeEntry.
type CreateTimeEntryOptions struct {
	// Start is required. The start time of the entry in UTC.
	Start time.Time
	// Description is an optional time entry description.
	Description *string
	// ProjectID is an optional project association.
	ProjectID *int
	// TaskID is an optional task association.
	TaskID *int
	// Billable marks the entry as billable.
	Billable *bool
	// Tags are names of tags to apply. New tags are created automatically.
	Tags []string
	// TagIDs are IDs of tags to apply.
	TagIDs []int
	// Stop is the optional stop time. If omitted the entry runs until stopped.
	Stop *time.Time
	// Duration in seconds. Defaults to -1 (running entry) when not set.
	// If Stop and Duration are both provided they must satisfy start+duration==stop.
	Duration *int
	// CreatedWith identifies the application creating the entry.
	// Defaults to "go-toggl" when empty.
	CreatedWith string
}

// UpdateTimeEntryOptions specifies the optional parameters to
// TimeEntriesService.UpdateTimeEntry.
type UpdateTimeEntryOptions struct {
	// Description updates the time entry description.
	Description *string
	// ProjectID updates the project association.
	ProjectID *int
	// TaskID updates the task association.
	TaskID *int
	// Billable updates the billable flag.
	Billable *bool
	// Tags are tag names to add or remove (see TagAction).
	Tags []string
	// TagIDs are tag IDs to add or remove (see TagAction).
	TagIDs []int
	// TagAction is "add" or "delete". Controls how Tags/TagIDs are applied.
	TagAction *string
	// Start updates the start time.
	Start *time.Time
	// Stop updates the stop time.
	Stop *time.Time
	// Duration updates the duration in seconds.
	Duration *int
}

// ListTimeEntries lists the current user's time entries with optional filters.
//
// API: GET /api/v9/me/time_entries
//
// See: https://engineering.toggl.com/docs/api/time_entries#get-timeentries
func (s *TimeEntriesService) ListTimeEntries(ctx context.Context, opts *ListTimeEntriesOptions) ([]*TimeEntry, *Response, error) {
	path := "/api/v9/me/time_entries"

	if opts != nil {
		params := url.Values{}
		if opts.Since != nil {
			params.Set("since", strconv.FormatInt(*opts.Since, 10))
		}
		if opts.Before != nil {
			params.Set("before", *opts.Before)
		}
		if opts.StartDate != nil {
			params.Set("start_date", *opts.StartDate)
		}
		if opts.EndDate != nil {
			params.Set("end_date", *opts.EndDate)
		}
		if q := params.Encode(); q != "" {
			path += "?" + q
		}
	}

	var entries []*TimeEntry
	resp, err := s.client.get(ctx, path, &entries)
	if err != nil {
		return nil, resp, err
	}

	return entries, resp, nil
}

// GetTimeEntry gets a single time entry by ID.
//
// API: GET /api/v9/me/time_entries/{time_entry_id}
//
// See: https://engineering.toggl.com/docs/api/time_entries#get-get-a-time-entry-by-id
func (s *TimeEntriesService) GetTimeEntry(ctx context.Context, entryID int) (*TimeEntry, *Response, error) {
	path := fmt.Sprintf("/api/v9/me/time_entries/%d", entryID)

	entry := new(TimeEntry)
	resp, err := s.client.get(ctx, path, entry)
	if err != nil {
		return nil, resp, err
	}

	return entry, resp, nil
}

// GetRunningTimeEntry returns the currently running time entry, if any.
// Returns a 404 error when no entry is running.
//
// API: GET /api/v9/me/time_entries/current
//
// See: https://engineering.toggl.com/docs/api/time_entries#get-get-current-time-entry
func (s *TimeEntriesService) GetRunningTimeEntry(ctx context.Context) (*TimeEntry, *Response, error) {
	path := "/api/v9/me/time_entries/current"

	entry := new(TimeEntry)
	resp, err := s.client.get(ctx, path, entry)
	if err != nil {
		return nil, resp, err
	}

	return entry, resp, nil
}

// CreateTimeEntry creates a new time entry in the given workspace.
//
// API: POST /api/v9/workspaces/{workspace_id}/time_entries
//
// See: https://engineering.toggl.com/docs/api/time_entries#post-timeentries
func (s *TimeEntriesService) CreateTimeEntry(ctx context.Context, workspaceID int, opts *CreateTimeEntryOptions) (*TimeEntry, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}
	if opts.Start.IsZero() {
		return nil, nil, fmt.Errorf("start is required")
	}

	// Resolve duration: explicit > derived from Stop > default running (-1).
	duration := -1
	if opts.Stop != nil && opts.Duration == nil {
		if !opts.Stop.After(opts.Start) {
			return nil, nil, fmt.Errorf("stop must be after start")
		}
		duration = int(opts.Stop.Sub(opts.Start).Seconds())
	} else if opts.Duration != nil {
		duration = *opts.Duration
	}

	createdWith := opts.CreatedWith
	if createdWith == "" {
		createdWith = "go-toggl"
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/time_entries", workspaceID)

	body := map[string]interface{}{
		"workspace_id": workspaceID,
		"start":        opts.Start.UTC().Format(time.RFC3339),
		"created_with": createdWith,
		"duration":     duration,
	}

	if opts.Description != nil {
		body["description"] = *opts.Description
	}
	if opts.ProjectID != nil {
		body["project_id"] = *opts.ProjectID
	}
	if opts.TaskID != nil {
		body["task_id"] = *opts.TaskID
	}
	if opts.Billable != nil {
		body["billable"] = *opts.Billable
	}
	if len(opts.Tags) > 0 {
		body["tags"] = opts.Tags
	}
	if len(opts.TagIDs) > 0 {
		body["tag_ids"] = opts.TagIDs
	}
	if opts.Stop != nil {
		body["stop"] = opts.Stop.UTC().Format(time.RFC3339)
	}

	entry := new(TimeEntry)
	resp, err := s.client.post(ctx, path, body, entry)
	if err != nil {
		return nil, resp, err
	}

	return entry, resp, nil
}

// StartTimeEntry creates a new running time entry in the given workspace.
// It is a convenience wrapper around CreateTimeEntry that always sets
// duration=-1 and clears any stop time, ensuring the entry is running.
//
// API: POST /api/v9/workspaces/{workspace_id}/time_entries
//
// See: https://engineering.toggl.com/docs/api/time_entries#post-timeentries
func (s *TimeEntriesService) StartTimeEntry(ctx context.Context, workspaceID int, opts *CreateTimeEntryOptions) (*TimeEntry, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}

	runOpts := *opts
	runOpts.Duration = Int(-1)
	runOpts.Stop = nil

	return s.CreateTimeEntry(ctx, workspaceID, &runOpts)
}

// UpdateTimeEntry updates an existing time entry.
//
// API: PUT /api/v9/workspaces/{workspace_id}/time_entries/{time_entry_id}
//
// See: https://engineering.toggl.com/docs/api/time_entries#put-timeentries
func (s *TimeEntriesService) UpdateTimeEntry(ctx context.Context, workspaceID, entryID int, opts *UpdateTimeEntryOptions) (*TimeEntry, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/time_entries/%d", workspaceID, entryID)

	body := map[string]interface{}{
		"workspace_id": workspaceID,
	}

	if opts.Description != nil {
		body["description"] = *opts.Description
	}
	if opts.ProjectID != nil {
		body["project_id"] = *opts.ProjectID
	}
	if opts.TaskID != nil {
		body["task_id"] = *opts.TaskID
	}
	if opts.Billable != nil {
		body["billable"] = *opts.Billable
	}
	if len(opts.Tags) > 0 {
		body["tags"] = opts.Tags
	}
	if len(opts.TagIDs) > 0 {
		body["tag_ids"] = opts.TagIDs
	}
	if opts.TagAction != nil {
		body["tag_action"] = *opts.TagAction
	}
	if opts.Start != nil {
		body["start"] = opts.Start.UTC().Format(time.RFC3339)
	}
	if opts.Stop != nil {
		body["stop"] = opts.Stop.UTC().Format(time.RFC3339)
	}
	if opts.Duration != nil {
		body["duration"] = *opts.Duration
	}

	entry := new(TimeEntry)
	resp, err := s.client.put(ctx, path, body, entry)
	if err != nil {
		return nil, resp, err
	}

	return entry, resp, nil
}

// DeleteTimeEntry deletes a time entry.
//
// API: DELETE /api/v9/workspaces/{workspace_id}/time_entries/{time_entry_id}
//
// See: https://engineering.toggl.com/docs/api/time_entries#delete-timeentries
func (s *TimeEntriesService) DeleteTimeEntry(ctx context.Context, workspaceID, entryID int) (*Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/time_entries/%d", workspaceID, entryID)
	return s.client.delete(ctx, path)
}

// StopTimeEntry stops a running time entry.
//
// API: PATCH /api/v9/workspaces/{workspace_id}/time_entries/{time_entry_id}/stop
//
// See: https://engineering.toggl.com/docs/api/time_entries#patch-stop-timeentry
func (s *TimeEntriesService) StopTimeEntry(ctx context.Context, workspaceID, entryID int) (*TimeEntry, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/time_entries/%d/stop", workspaceID, entryID)

	entry := new(TimeEntry)
	resp, err := s.client.patch(ctx, path, nil, entry)
	if err != nil {
		return nil, resp, err
	}

	return entry, resp, nil
}
