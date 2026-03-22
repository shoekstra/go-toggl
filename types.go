package toggl

import "time"

// Common types

// Workspace represents a Toggl workspace.
type Workspace struct {
	ID                          int       `json:"id"`
	Name                        string    `json:"name"`
	OrganizationID              int       `json:"organization_id"`
	ActiveProjectCount          int       `json:"active_project_count"`
	At                          time.Time `json:"at"`
	Premium                     bool      `json:"premium"`
	BusinessWs                  bool      `json:"business_ws"`
	DefaultCurrency             string    `json:"default_currency"`
	DefaultHourlyRate           *float64  `json:"default_hourly_rate"`
	OnlyAdminsMayCreateProjects bool      `json:"only_admins_may_create_projects"`
	OnlyAdminsMayCreateTags     bool      `json:"only_admins_may_create_tags"`
	OnlyAdminsSeeDashboard      bool      `json:"only_admins_see_team_dashboard"`
	ProjectsBillableByDefault   bool      `json:"projects_billable_by_default"`
	ProjectsEnforceBillable     bool      `json:"projects_enforce_billable"`
	ProjectsPrivateByDefault    bool      `json:"projects_private_by_default"`
	ReportsCollapse             bool      `json:"reports_collapse"`
	Rounding                    int       `json:"rounding"`
	RoundingMinutes             int       `json:"rounding_minutes"`
	LogoURL                     *string   `json:"logo_url"`
	IcalEnabled                 bool      `json:"ical_enabled"`
	IcalURL                     *string   `json:"ical_url"`
	Role                        string    `json:"role"`
	Permissions                 []string  `json:"permissions"`
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
	Description *string    `json:"description"`
	ProjectID   *int       `json:"project_id"`
	TaskID      *int       `json:"task_id"`
	ClientID    *int       `json:"client_id"`
	WorkspaceID int        `json:"workspace_id"`
	UserID      int        `json:"user_id"`
	Billable    bool       `json:"billable"`
	TagIDs      []int      `json:"tag_ids"`
	Tags        []string   `json:"tags"`
	Start       time.Time  `json:"start"`
	Stop        *time.Time `json:"stop"`
	Duration    int        `json:"duration"`
	CreatedWith string     `json:"created_with"`
	At          time.Time  `json:"at"`
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
