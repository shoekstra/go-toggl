package toggl

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ProjectsService handles operations related to projects.
type ProjectsService struct {
	client *Client
}

// ListProjectsOptions specifies the optional parameters to
// ProjectsService.ListProjects.
type ListProjectsOptions struct {
	// Active filters by project status. Valid values: "true", "false", "both".
	// Defaults to "true" when not set.
	Active *string
	// Billable filters by billable status (premium feature).
	Billable *bool
	// Name filters projects by name.
	Name *string
	// Since filters projects modified since this UNIX timestamp.
	Since *int64
	// Page is the page number for pagination.
	Page *int
	// PerPage is the number of results per page.
	PerPage *int
}

// CreateProjectOptions specifies the parameters to
// ProjectsService.CreateProject.
type CreateProjectOptions struct {
	// Name is required. The project name.
	Name string
	// Active sets the project as active or archived.
	Active *bool
	// Billable marks the project as billable (premium feature).
	Billable *bool
	// ClientID associates the project with a client.
	ClientID *int
	// Color is the project color in hex format (e.g. "#06aaf5").
	Color *string
	// Currency sets the project currency (premium feature).
	Currency *string
	// EndDate is the project end date in YYYY-MM-DD format.
	EndDate *string
	// EstimatedHours is the estimated number of hours for the project.
	EstimatedHours *int
	// FixedFee is the fixed fee for the project (premium feature).
	FixedFee *float64
	// IsPrivate controls whether the project is private.
	IsPrivate *bool
	// Rate is the hourly rate for the project (premium feature).
	Rate *float64
	// StartDate is the project start date in YYYY-MM-DD format.
	StartDate *string
	// Template marks the project as a template.
	Template *bool
	// TemplateID creates the project from a template.
	TemplateID *int
}

// UpdateProjectOptions specifies the optional parameters to
// ProjectsService.UpdateProject.
type UpdateProjectOptions struct {
	// Name updates the project name.
	Name *string
	// Active sets the project as active or archived.
	Active *bool
	// Billable marks the project as billable (premium feature).
	Billable *bool
	// ClientID updates the client association.
	ClientID *int
	// Color updates the project color in hex format.
	Color *string
	// Currency updates the project currency (premium feature).
	Currency *string
	// EndDate updates the project end date in YYYY-MM-DD format.
	EndDate *string
	// EstimatedHours updates the estimated number of hours.
	EstimatedHours *int
	// FixedFee updates the fixed fee (premium feature).
	FixedFee *float64
	// IsPrivate controls whether the project is private.
	IsPrivate *bool
	// Rate updates the hourly rate (premium feature).
	Rate *float64
	// StartDate updates the project start date in YYYY-MM-DD format.
	StartDate *string
	// Template marks the project as a template.
	Template *bool
}

// ListProjects lists all projects in the given workspace.
//
// API: GET /api/v9/workspaces/{workspace_id}/projects
//
// See: https://engineering.toggl.com/docs/api/projects#get-workspaceprojects
func (s *ProjectsService) ListProjects(ctx context.Context, workspaceID int, opts *ListProjectsOptions) ([]*Project, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/projects", workspaceID)

	if opts != nil {
		params := url.Values{}
		if opts.Active != nil {
			params.Set("active", *opts.Active)
		}
		if opts.Billable != nil {
			params.Set("billable", strconv.FormatBool(*opts.Billable))
		}
		if opts.Name != nil {
			params.Set("name", *opts.Name)
		}
		if opts.Since != nil {
			params.Set("since", strconv.FormatInt(*opts.Since, 10))
		}
		if opts.Page != nil {
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if q := params.Encode(); q != "" {
			path += "?" + q
		}
	}

	var projects []*Project
	resp, err := s.client.get(ctx, path, &projects)
	if err != nil {
		return nil, resp, err
	}

	return projects, resp, nil
}

// GetProject gets a single project by ID.
//
// API: GET /api/v9/workspaces/{workspace_id}/projects/{project_id}
//
// See: https://engineering.toggl.com/docs/api/projects#get-workspaceproject
func (s *ProjectsService) GetProject(ctx context.Context, workspaceID, projectID int) (*Project, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/projects/%d", workspaceID, projectID)

	project := new(Project)
	resp, err := s.client.get(ctx, path, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}

// CreateProject creates a new project in the given workspace.
//
// API: POST /api/v9/workspaces/{workspace_id}/projects
//
// See: https://engineering.toggl.com/docs/api/projects#post-workspaceprojects
func (s *ProjectsService) CreateProject(ctx context.Context, workspaceID int, opts *CreateProjectOptions) (*Project, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}
	if opts.Name == "" {
		return nil, nil, fmt.Errorf("Name is required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/projects", workspaceID)

	body := map[string]interface{}{
		"name": opts.Name,
	}

	if opts.Active != nil {
		body["active"] = *opts.Active
	}
	if opts.Billable != nil {
		body["billable"] = *opts.Billable
	}
	if opts.ClientID != nil {
		body["client_id"] = *opts.ClientID
	}
	if opts.Color != nil {
		body["color"] = *opts.Color
	}
	if opts.Currency != nil {
		body["currency"] = *opts.Currency
	}
	if opts.EndDate != nil {
		body["end_date"] = *opts.EndDate
	}
	if opts.EstimatedHours != nil {
		body["estimated_hours"] = *opts.EstimatedHours
	}
	if opts.FixedFee != nil {
		body["fixed_fee"] = *opts.FixedFee
	}
	if opts.IsPrivate != nil {
		body["is_private"] = *opts.IsPrivate
	}
	if opts.Rate != nil {
		body["rate"] = *opts.Rate
	}
	if opts.StartDate != nil {
		body["start_date"] = *opts.StartDate
	}
	if opts.Template != nil {
		body["template"] = *opts.Template
	}
	if opts.TemplateID != nil {
		body["template_id"] = *opts.TemplateID
	}

	project := new(Project)
	resp, err := s.client.post(ctx, path, body, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}

// UpdateProject updates an existing project.
//
// API: PUT /api/v9/workspaces/{workspace_id}/projects/{project_id}
//
// See: https://engineering.toggl.com/docs/api/projects#put-workspaceproject
func (s *ProjectsService) UpdateProject(ctx context.Context, workspaceID, projectID int, opts *UpdateProjectOptions) (*Project, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/projects/%d", workspaceID, projectID)

	body := make(map[string]interface{})

	if opts.Name != nil {
		body["name"] = *opts.Name
	}
	if opts.Active != nil {
		body["active"] = *opts.Active
	}
	if opts.Billable != nil {
		body["billable"] = *opts.Billable
	}
	if opts.ClientID != nil {
		body["client_id"] = *opts.ClientID
	}
	if opts.Color != nil {
		body["color"] = *opts.Color
	}
	if opts.Currency != nil {
		body["currency"] = *opts.Currency
	}
	if opts.EndDate != nil {
		body["end_date"] = *opts.EndDate
	}
	if opts.EstimatedHours != nil {
		body["estimated_hours"] = *opts.EstimatedHours
	}
	if opts.FixedFee != nil {
		body["fixed_fee"] = *opts.FixedFee
	}
	if opts.IsPrivate != nil {
		body["is_private"] = *opts.IsPrivate
	}
	if opts.Rate != nil {
		body["rate"] = *opts.Rate
	}
	if opts.StartDate != nil {
		body["start_date"] = *opts.StartDate
	}
	if opts.Template != nil {
		body["template"] = *opts.Template
	}

	project := new(Project)
	resp, err := s.client.put(ctx, path, body, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, nil
}

// DeleteProject deletes a project from the given workspace.
//
// The optional teDeletionMode controls what happens to time entries linked
// to this project: "delete" removes them, "unassign" removes the project
// association but keeps the entries. When not set the server default applies.
//
// API: DELETE /api/v9/workspaces/{workspace_id}/projects/{project_id}
//
// See: https://engineering.toggl.com/docs/api/projects#delete-workspaceproject
func (s *ProjectsService) DeleteProject(ctx context.Context, workspaceID, projectID int, teDeletionMode *string) (*Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/projects/%d", workspaceID, projectID)

	if teDeletionMode != nil {
		params := url.Values{}
		params.Set("teDeletionMode", *teDeletionMode)
		path += "?" + params.Encode()
	}

	return s.client.delete(ctx, path)
}
