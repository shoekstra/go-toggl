package toggl

import (
	"encoding/json"
	"time"
)

// Common types

// Me represents the authenticated user's profile.
type Me struct {
	// ID is the user's unique identifier.
	ID int `json:"id"`
	// Email is the user's email address.
	Email string `json:"email"`
	// Fullname is the user's display name.
	Fullname string `json:"fullname"`
	// Timezone is the user's configured IANA timezone (e.g. "Europe/Amsterdam").
	Timezone string `json:"timezone"`
	// DefaultWorkspaceID is the ID of the user's default workspace.
	DefaultWorkspaceID int `json:"default_workspace_id"`
	// BeginningOfWeek is the first day of the week: 0 = Sunday, 1 = Monday, … 6 = Saturday.
	BeginningOfWeek int `json:"beginning_of_week"`
	// ImageURL is the URL of the user's profile image.
	ImageURL string `json:"image_url"`
	// CountryID is the user's country (may be absent).
	CountryID *int `json:"country_id"`
	// HasPassword indicates whether the account has a password set.
	HasPassword bool `json:"has_password"`
	// OpenIDEnabled indicates whether OpenID login is enabled for this account.
	OpenIDEnabled bool `json:"openid_enabled"`
	// APIToken is the user's API token (identical to the token used to authenticate this request).
	APIToken string `json:"api_token"`
	// At is the timestamp of the most recent change to the user record.
	At time.Time `json:"at"`
	// CreatedAt is when the account was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the account was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

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
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	WorkspaceID      int       `json:"workspace_id"`
	ClientID         *int      `json:"client_id"`
	Active           bool      `json:"active"`
	IsPrivate        bool      `json:"is_private"`
	Billable         *bool     `json:"billable"`
	AutoEstimates    *bool     `json:"auto_estimates"`
	EstimatedHours   *int      `json:"estimated_hours"`
	EstimatedSeconds *int      `json:"estimated_seconds"`
	ActualHours      *int      `json:"actual_hours"`
	ActualSeconds    *int      `json:"actual_seconds"`
	Color            string    `json:"color"`
	Rate             float64   `json:"rate"`
	RateLastUpdated  *string   `json:"rate_last_updated"`
	Currency         *string   `json:"currency"`
	Template         *bool     `json:"template"`
	TemplateID       *int      `json:"template_id"`
	FixedFee         float64   `json:"fixed_fee"`
	StartDate        *string   `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	At               time.Time `json:"at"`
	CreatedAt        time.Time `json:"created_at"`
}

// Tag represents a Toggl tag.
type Tag struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	WorkspaceID        int        `json:"workspace_id"` // WorkspaceID is the workspace this tag belongs to (json: "workspace_id").
	Workspace          int        `json:"-"`            // Workspace is an alias for WorkspaceID retained for backwards compatibility.
	CreatorID          int        `json:"creator_id"`
	At                 time.Time  `json:"at"`
	DeletedAt          *time.Time `json:"deleted_at"`
	IntegrationExtID   *string    `json:"integration_ext_id"`
	IntegrationExtType *string    `json:"integration_ext_type"`
	Permissions        []string   `json:"permissions"`
}

// UnmarshalJSON implements json.Unmarshaler for Tag, accepting both
// "workspace_id" (current API) and "workspace" (legacy) and keeping
// both WorkspaceID and Workspace in sync.
func (t *Tag) UnmarshalJSON(data []byte) error {
	type tagAlias Tag
	aux := struct {
		*tagAlias
		WorkspaceAlt int `json:"workspace"`
	}{tagAlias: (*tagAlias)(t)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if t.WorkspaceID == 0 && aux.WorkspaceAlt != 0 {
		t.WorkspaceID = aux.WorkspaceAlt
	}
	t.Workspace = t.WorkspaceID
	return nil
}

// WorkspaceClient represents a Toggl client (customer/organization) within a workspace.
type WorkspaceClient struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	WorkspaceID       int       `json:"wid"`
	Notes             *string   `json:"notes"`
	Archived          bool      `json:"archived"`
	At                time.Time `json:"at"`
	CreatorID         int       `json:"creator_id"`
	ExternalReference *string   `json:"external_reference"`
	Permissions       []string  `json:"permissions"`
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

// Reports API response types

// SummaryReportData is the response from ReportsService.SummaryReport.
type SummaryReportData struct {
	Groups []SummaryGroup `json:"groups"`
}

// SummaryGroup is a top-level grouping in a summary report.
type SummaryGroup struct {
	ID        *int              `json:"id"`
	Seconds   int               `json:"seconds"`
	Rates     []BillableRate    `json:"rates"`
	SubGroups []SummarySubGroup `json:"sub_groups"`
}

// SummarySubGroup is a sub-grouping within a SummaryGroup.
type SummarySubGroup struct {
	ID      *int           `json:"id"`
	Seconds int            `json:"seconds"`
	Rates   []BillableRate `json:"rates"`
	Title   *string        `json:"title"`
}

// DetailedTimeEntry is a single time entry row in a detailed report.
type DetailedTimeEntry struct {
	Billable            bool     `json:"billable"`
	BillableAmountCents *int     `json:"billable_amount_in_cents"`
	ClientName          *string  `json:"client_name"`
	Currency            *string  `json:"currency"`
	Description         *string  `json:"description"`
	HourlyRateCents     *int     `json:"hourly_rate_in_cents"`
	ProjectColor        *string  `json:"project_color"`
	ProjectHex          *string  `json:"project_hex"`
	ProjectID           *int     `json:"project_id"`
	ProjectName         *string  `json:"project_name"`
	RowNumber           int      `json:"row_number"`
	TagIDs              []int    `json:"tag_ids"`
	TagNames            []string `json:"tag_names"`
	UserID              *int     `json:"user_id"`
	Username            *string  `json:"username"`
}

// WeeklyReportEntry is a single row in a weekly report.
type WeeklyReportEntry struct {
	ProjectID *int    `json:"project_id"`
	ClientID  *int    `json:"client_id"`
	UserID    *int    `json:"user_id"`
	Title     *string `json:"title"`
	Seconds   []int   `json:"seconds"`
}

// TotalsReport is the response from ReportsService.DetailedReportTotals.
type TotalsReport struct {
	BillableAmountCents int            `json:"billable_amount_in_cents"`
	LabourCostCents     int            `json:"labour_cost_in_cents"`
	Seconds             int            `json:"seconds"`
	TrackedDays         int            `json:"tracked_days"`
	Resolution          string         `json:"resolution"`
	Graph               []TotalsGraph  `json:"graph"`
	Rates               []BillableRate `json:"rates"`
}

// TotalsGraph is a single data point in the TotalsReport graph.
type TotalsGraph struct {
	BillableAmountCents int `json:"billable_amount_in_cents"`
	LabourCostCents     int `json:"labour_cost_in_cents"`
	Seconds             int `json:"seconds"`
}

// BillableRate represents an hourly rate with associated billable seconds.
type BillableRate struct {
	BillableSeconds int    `json:"billable_seconds"`
	Currency        string `json:"currency"`
	HourlyRateCents int    `json:"hourly_rate_in_cents"`
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
