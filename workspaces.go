package toggl

import (
	"context"
	"fmt"
)

// WorkspacesService handles operations related to workspaces.
type WorkspacesService struct {
	client *Client
}

// UpdateWorkspaceOptions specifies the optional parameters to
// WorkspacesService.UpdateWorkspace.
type UpdateWorkspaceOptions struct {
	// Name updates the workspace name.
	Name *string
	// Admins replaces the list of workspace admin user IDs.
	Admins []int
	// DefaultCurrency updates the default currency (premium feature).
	DefaultCurrency *string
	// DefaultHourlyRate updates the default hourly rate (premium feature).
	DefaultHourlyRate *float64
	// OnlyAdminsMayCreateProjects restricts project creation to admins.
	OnlyAdminsMayCreateProjects *bool
	// OnlyAdminsMayCreateTags restricts tag creation to admins.
	OnlyAdminsMayCreateTags *bool
	// OnlyAdminsSeeDashboard restricts the team dashboard to admins.
	OnlyAdminsSeeDashboard *bool
	// ProjectsBillableByDefault sets new projects as billable by default (premium feature).
	ProjectsBillableByDefault *bool
	// ProjectsEnforceBillable enforces the billable setting when tracking time to projects.
	ProjectsEnforceBillable *bool
	// ProjectsPrivateByDefault sets new projects as private by default.
	ProjectsPrivateByDefault *bool
	// LimitPublicProjectData limits public project data in reports to admins.
	LimitPublicProjectData *bool
	// RateChangeMode controls how rate changes are applied (premium feature).
	// Valid values: "start-today", "override-current", "override-all".
	RateChangeMode *string
	// ReportsCollapse sets whether reports are collapsed by default.
	ReportsCollapse *bool
	// Rounding sets the default rounding (premium feature).
	// Values: 0 = nearest, 1 = round up, -1 = round down.
	Rounding *int
	// RoundingMinutes sets the default rounding interval in minutes (premium feature).
	RoundingMinutes *int
}

// ListWorkspaces lists all workspaces for the current user.
//
// API: GET /api/v9/me/workspaces
//
// See: https://engineering.toggl.com/docs/api/me#get-workspaces
func (s *WorkspacesService) ListWorkspaces(ctx context.Context) ([]*Workspace, *Response, error) {
	path := "/api/v9/me/workspaces"

	var workspaces []*Workspace
	resp, err := s.client.get(ctx, path, &workspaces)
	if err != nil {
		return nil, resp, err
	}

	return workspaces, resp, nil
}

// GetWorkspace gets a single workspace by ID.
//
// API: GET /api/v9/workspaces/{workspace_id}
//
// See: https://engineering.toggl.com/docs/api/workspaces#get-get-single-workspace
func (s *WorkspacesService) GetWorkspace(ctx context.Context, workspaceID int) (*Workspace, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d", workspaceID)

	workspace := new(Workspace)
	resp, err := s.client.get(ctx, path, workspace)
	if err != nil {
		return nil, resp, err
	}

	return workspace, resp, nil
}

// UpdateWorkspace updates a workspace.
//
// API: PUT /api/v9/workspaces/{workspace_id}
//
// See: https://engineering.toggl.com/docs/api/workspaces#put-update-workspace
func (s *WorkspacesService) UpdateWorkspace(ctx context.Context, workspaceID int, opts *UpdateWorkspaceOptions) (*Workspace, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d", workspaceID)

	body := make(map[string]interface{})

	if opts.Name != nil {
		body["name"] = *opts.Name
	}
	if len(opts.Admins) > 0 {
		body["admins"] = opts.Admins
	}
	if opts.DefaultCurrency != nil {
		body["default_currency"] = *opts.DefaultCurrency
	}
	if opts.DefaultHourlyRate != nil {
		body["default_hourly_rate"] = *opts.DefaultHourlyRate
	}
	if opts.OnlyAdminsMayCreateProjects != nil {
		body["only_admins_may_create_projects"] = *opts.OnlyAdminsMayCreateProjects
	}
	if opts.OnlyAdminsMayCreateTags != nil {
		body["only_admins_may_create_tags"] = *opts.OnlyAdminsMayCreateTags
	}
	if opts.OnlyAdminsSeeDashboard != nil {
		body["only_admins_see_team_dashboard"] = *opts.OnlyAdminsSeeDashboard
	}
	if opts.ProjectsBillableByDefault != nil {
		body["projects_billable_by_default"] = *opts.ProjectsBillableByDefault
	}
	if opts.ProjectsEnforceBillable != nil {
		body["projects_enforce_billable"] = *opts.ProjectsEnforceBillable
	}
	if opts.ProjectsPrivateByDefault != nil {
		body["projects_private_by_default"] = *opts.ProjectsPrivateByDefault
	}
	if opts.LimitPublicProjectData != nil {
		body["limit_public_project_data"] = *opts.LimitPublicProjectData
	}
	if opts.RateChangeMode != nil {
		body["rate_change_mode"] = *opts.RateChangeMode
	}
	if opts.ReportsCollapse != nil {
		body["reports_collapse"] = *opts.ReportsCollapse
	}
	if opts.Rounding != nil {
		body["rounding"] = *opts.Rounding
	}
	if opts.RoundingMinutes != nil {
		body["rounding_minutes"] = *opts.RoundingMinutes
	}

	workspace := new(Workspace)
	resp, err := s.client.put(ctx, path, body, workspace)
	if err != nil {
		return nil, resp, err
	}

	return workspace, resp, nil
}
