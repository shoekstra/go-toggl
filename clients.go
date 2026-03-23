package toggl

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ClientsService handles operations related to clients.
type ClientsService struct {
	client *Client
}

// ListClientsOptions specifies the optional parameters to
// ClientsService.ListClients.
type ListClientsOptions struct {
	// Status filters clients by status. Valid values: "active", "archived", "both".
	// Defaults to "active" when not set.
	Status *string
	// Name filters clients by name (case-insensitive).
	Name *string
	// Page is the page number for pagination.
	Page *int
	// PerPage is the number of results per page.
	PerPage *int
}

// CreateClientOptions specifies the parameters to
// ClientsService.CreateClient.
type CreateClientOptions struct {
	// Name is required. The client name.
	Name string
	// Notes is an optional free-text note.
	Notes *string
	// ExternalReference is an optional reference to an external system.
	ExternalReference *string
}

// UpdateClientOptions specifies the optional parameters to
// ClientsService.UpdateClient.
type UpdateClientOptions struct {
	// Name updates the client name.
	Name *string
	// Notes updates the free-text note.
	Notes *string
	// ExternalReference updates the external system reference.
	ExternalReference *string
}

// RestoreClientOptions specifies the optional parameters to
// ClientsService.RestoreClient.
type RestoreClientOptions struct {
	// RestoreAllProjects restores all projects linked to the client.
	RestoreAllProjects *bool
	// Projects is a list of specific project IDs to restore alongside the client.
	Projects []int
}

// ListClients lists all clients in the given workspace.
//
// API: GET /api/v9/workspaces/{workspace_id}/clients
//
// See: https://engineering.toggl.com/docs/api/clients#get-list-clients
func (s *ClientsService) ListClients(ctx context.Context, workspaceID int, opts *ListClientsOptions) ([]*TogglClient, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/clients", workspaceID)

	if opts != nil {
		params := url.Values{}
		if opts.Status != nil {
			switch *opts.Status {
			case "active", "archived", "both":
				params.Set("status", *opts.Status)
			default:
				return nil, nil, fmt.Errorf("invalid status %q: must be \"active\", \"archived\", or \"both\"", *opts.Status)
			}
		}
		if opts.Name != nil {
			params.Set("name", *opts.Name)
		}
		if opts.Page != nil {
			if *opts.Page < 1 {
				return nil, nil, fmt.Errorf("invalid page %d: must be greater than 0", *opts.Page)
			}
			params.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			if *opts.PerPage < 1 {
				return nil, nil, fmt.Errorf("invalid per_page %d: must be greater than 0", *opts.PerPage)
			}
			params.Set("per_page", strconv.Itoa(*opts.PerPage))
		}
		if q := params.Encode(); q != "" {
			path += "?" + q
		}
	}

	var clients []*TogglClient
	resp, err := s.client.get(ctx, path, &clients)
	if err != nil {
		return nil, resp, err
	}

	return clients, resp, nil
}

// GetClient gets a single client by ID.
//
// API: GET /api/v9/workspaces/{workspace_id}/clients/{client_id}
//
// See: https://engineering.toggl.com/docs/api/clients#get-load-client-by-id
func (s *ClientsService) GetClient(ctx context.Context, workspaceID, clientID int) (*TogglClient, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/clients/%d", workspaceID, clientID)

	c := new(TogglClient)
	resp, err := s.client.get(ctx, path, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// CreateClient creates a new client in the given workspace.
//
// API: POST /api/v9/workspaces/{workspace_id}/clients
//
// See: https://engineering.toggl.com/docs/api/clients#post-create-client
func (s *ClientsService) CreateClient(ctx context.Context, workspaceID int, opts *CreateClientOptions) (*TogglClient, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}
	if opts.Name == "" {
		return nil, nil, fmt.Errorf("Name is required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/clients", workspaceID)

	body := map[string]interface{}{
		"name": opts.Name,
	}

	if opts.Notes != nil {
		body["notes"] = *opts.Notes
	}
	if opts.ExternalReference != nil {
		body["external_reference"] = *opts.ExternalReference
	}

	c := new(TogglClient)
	resp, err := s.client.post(ctx, path, body, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// UpdateClient updates an existing client.
//
// API: PUT /api/v9/workspaces/{workspace_id}/clients/{client_id}
//
// See: https://engineering.toggl.com/docs/api/clients#put-change-client
func (s *ClientsService) UpdateClient(ctx context.Context, workspaceID, clientID int, opts *UpdateClientOptions) (*TogglClient, *Response, error) {
	if opts == nil {
		return nil, nil, fmt.Errorf("options required")
	}

	path := fmt.Sprintf("/api/v9/workspaces/%d/clients/%d", workspaceID, clientID)

	body := make(map[string]interface{})

	if opts.Name != nil {
		body["name"] = *opts.Name
	}
	if opts.Notes != nil {
		body["notes"] = *opts.Notes
	}
	if opts.ExternalReference != nil {
		body["external_reference"] = *opts.ExternalReference
	}

	if len(body) == 0 {
		return nil, nil, fmt.Errorf("at least one field must be set")
	}

	c := new(TogglClient)
	resp, err := s.client.put(ctx, path, body, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}

// DeleteClient deletes a client from the given workspace.
//
// API: DELETE /api/v9/workspaces/{workspace_id}/clients/{client_id}
//
// See: https://engineering.toggl.com/docs/api/clients#delete-delete-client
func (s *ClientsService) DeleteClient(ctx context.Context, workspaceID, clientID int) (*Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/clients/%d", workspaceID, clientID)
	return s.client.delete(ctx, path)
}

// ArchiveResponse is the response from ClientsService.ArchiveClient.
// Items contains the IDs of projects archived alongside the client.
type ArchiveResponse struct {
	Items []int `json:"items"`
}

// ArchiveClient archives a client and its related projects. This is a
// premium workspace feature.
//
// API: POST /api/v9/workspaces/{workspace_id}/clients/{client_id}/archive
//
// See: https://engineering.toggl.com/docs/api/clients#post-archive-client
func (s *ClientsService) ArchiveClient(ctx context.Context, workspaceID, clientID int) (*ArchiveResponse, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/clients/%d/archive", workspaceID, clientID)

	ar := new(ArchiveResponse)
	resp, err := s.client.post(ctx, path, nil, ar)
	if err != nil {
		return nil, resp, err
	}

	return ar, resp, nil
}

// RestoreClient restores an archived client and optionally its projects.
//
// API: POST /api/v9/workspaces/{workspace_id}/clients/{client_id}/restore
//
// See: https://engineering.toggl.com/docs/api/clients#post-restore-client
func (s *ClientsService) RestoreClient(ctx context.Context, workspaceID, clientID int, opts *RestoreClientOptions) (*TogglClient, *Response, error) {
	path := fmt.Sprintf("/api/v9/workspaces/%d/clients/%d/restore", workspaceID, clientID)

	var body interface{}
	if opts != nil {
		m := make(map[string]interface{})
		if opts.RestoreAllProjects != nil {
			m["restore_all_projects"] = *opts.RestoreAllProjects
		}
		if len(opts.Projects) > 0 {
			m["projects"] = opts.Projects
		}
		if len(m) > 0 {
			body = m
		}
	}

	c := new(TogglClient)
	resp, err := s.client.post(ctx, path, body, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, nil
}
