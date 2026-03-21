package toggl

import "time"

// Common types

// Workspace represents a Toggl workspace.
type Workspace struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Premium       bool       `json:"premium"`
	Admin         bool       `json:"admin"`
	SuspendDate   *time.Time `json:"suspend_date"`
	SuspendReason *string    `json:"suspend_reason"`
}

// Project represents a Toggl project.
type Project struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Color          *string   `json:"color"`
	Billable       *bool     `json:"billable"`
	Active         bool      `json:"active"`
	Public         bool      `json:"public"`
	Template       bool      `json:"template"`
	TemplateID     *int      `json:"template_id"`
	AutoEstimates  *bool     `json:"auto_estimates"`
	EstimatedHours *int      `json:"estimated_hours"`
	Rate           *float64  `json:"rate"`
	Currency       *string   `json:"currency"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Tag represents a Toggl tag.
type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Workspace int       `json:"workspace_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TogglClient represents a Toggl client (customer/organization).
type TogglClient struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Workspace int       `json:"workspace_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TimeEntry represents a time entry in Toggl.
type TimeEntry struct {
	ID          int        `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Project     *int       `json:"project_id"`
	Client      *int       `json:"client_id"`
	Tags        []string   `json:"tags"`
	Billable    *bool      `json:"billable"`
	Start       time.Time  `json:"start"`
	Stop        *time.Time `json:"stop"`
	Duration    int        `json:"duration"`
	Workspace   int        `json:"workspace_id"`
	User        int        `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Helper functions for pointer creation

// String returns a pointer to the provided string value.
func String(s string) *string {
	return &s
}

// Int returns a pointer to the provided int value.
func Int(i int) *int {
	return &i
}

// Bool returns a pointer to the provided bool value.
func Bool(b bool) *bool {
	return &b
}

// Float64 returns a pointer to the provided float64 value.
func Float64(f float64) *float64 {
	return &f
}

// Time returns a pointer to the provided time.Time value.
func Time(t time.Time) *time.Time {
	return &t
}
