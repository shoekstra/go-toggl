package toggl

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

const workspaceJSON = `{
	"id": 1,
	"name": "My Workspace",
	"organization_id": 10,
	"active_project_count": 3,
	"at": "2024-01-15T10:00:00Z",
	"premium": false,
	"business_ws": false,
	"default_currency": "USD",
	"default_hourly_rate": 0,
	"only_admins_may_create_projects": false,
	"only_admins_may_create_tags": false,
	"only_admins_see_team_dashboard": false,
	"projects_billable_by_default": true,
	"projects_enforce_billable": false,
	"projects_private_by_default": true,
	"reports_collapse": true,
	"rounding": 0,
	"rounding_minutes": 0,
	"logo_url": "",
	"ical_enabled": false,
	"ical_url": "",
	"role": "admin",
	"permissions": ["read", "write"]
}`

func TestWorkspacesService_ListWorkspaces(t *testing.T) {
	wantPath := "/api/v9/me/workspaces"

	tests := []struct {
		name       string
		statusCode int
		response   string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response:   "[" + workspaceJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "multiple workspaces",
			statusCode: http.StatusOK,
			response:   "[" + workspaceJSON + "," + workspaceJSON + "]",
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "empty result",
			statusCode: http.StatusOK,
			response:   "[]",
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantCount:  0,
			wantErr:    true,
		},
		{
			name:       "server error",
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
			workspaces, _, err := client.Workspaces.ListWorkspaces(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("ListWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(workspaces) != tt.wantCount {
				t.Errorf("ListWorkspaces() count = %d, want %d", len(workspaces), tt.wantCount)
			}
		})
	}
}

func TestWorkspacesService_ListWorkspaces_Fields(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[" + workspaceJSON + "]"))
	})

	client := testClient(t, handler)
	workspaces, _, err := client.Workspaces.ListWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("ListWorkspaces() error = %v", err)
	}
	if len(workspaces) != 1 {
		t.Fatalf("ListWorkspaces() count = %d, want 1", len(workspaces))
	}

	ws := workspaces[0]
	if ws.ID != 1 {
		t.Errorf("ID = %d, want 1", ws.ID)
	}
	if ws.Name != "My Workspace" {
		t.Errorf("Name = %q, want %q", ws.Name, "My Workspace")
	}
	if ws.OrganizationID != 10 {
		t.Errorf("OrganizationID = %d, want 10", ws.OrganizationID)
	}
	if ws.ActiveProjectCount != 3 {
		t.Errorf("ActiveProjectCount = %d, want 3", ws.ActiveProjectCount)
	}
	if ws.DefaultCurrency != "USD" {
		t.Errorf("DefaultCurrency = %q, want %q", ws.DefaultCurrency, "USD")
	}
	if ws.Role != "admin" {
		t.Errorf("Role = %q, want %q", ws.Role, "admin")
	}
}

func TestWorkspacesService_GetWorkspace(t *testing.T) {
	tests := []struct {
		name        string
		workspaceID int
		statusCode  int
		response    string
		wantID      int
		wantErr     bool
	}{
		{
			name:        "success",
			workspaceID: 1,
			statusCode:  http.StatusOK,
			response:    workspaceJSON,
			wantID:      1,
			wantErr:     false,
		},
		{
			name:        "not found",
			workspaceID: 999,
			statusCode:  http.StatusNotFound,
			response:    `{"error": "not found"}`,
			wantID:      0,
			wantErr:     true,
		},
		{
			name:        "unauthorized",
			workspaceID: 1,
			statusCode:  http.StatusUnauthorized,
			response:    `{"error": "unauthorized"}`,
			wantID:      0,
			wantErr:     true,
		},
		{
			name:        "server error",
			workspaceID: 1,
			statusCode:  http.StatusInternalServerError,
			response:    `{"error": "internal server error"}`,
			wantID:      0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/%d", tt.workspaceID)
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
			workspace, _, err := client.Workspaces.GetWorkspace(context.Background(), tt.workspaceID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && workspace.ID != tt.wantID {
				t.Errorf("GetWorkspace() ID = %d, want %d", workspace.ID, tt.wantID)
			}
		})
	}
}

func TestWorkspacesService_UpdateWorkspace(t *testing.T) {
	tests := []struct {
		name        string
		workspaceID int
		opts        *UpdateWorkspaceOptions
		statusCode  int
		response    string
		wantID      int
		wantErr     bool
	}{
		{
			name:        "success",
			workspaceID: 1,
			opts:        &UpdateWorkspaceOptions{Name: String("Renamed Workspace")},
			statusCode:  http.StatusOK,
			response:    workspaceJSON,
			wantID:      1,
			wantErr:     false,
		},
		{
			name:        "success with all options",
			workspaceID: 1,
			opts: &UpdateWorkspaceOptions{
				Name:                        String("Renamed"),
				Admins:                      []int{1, 2},
				DefaultCurrency:             String("EUR"),
				DefaultHourlyRate:           Float64(100.0),
				OnlyAdminsMayCreateProjects: Bool(true),
				OnlyAdminsMayCreateTags:     Bool(true),
				OnlyAdminsSeeDashboard:      Bool(true),
				ProjectsBillableByDefault:   Bool(false),
				ProjectsEnforceBillable:     Bool(true),
				ProjectsPrivateByDefault:    Bool(false),
				LimitPublicProjectData:      Bool(true),
				RateChangeMode:              String("start-today"),
				ReportsCollapse:             Bool(false),
				Rounding:                    Int(1),
				RoundingMinutes:             Int(15),
			},
			statusCode: http.StatusOK,
			response:   workspaceJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:        "nil options",
			workspaceID: 1,
			opts:        nil,
			wantErr:     true,
		},
		{
			name:        "not found",
			workspaceID: 999,
			opts:        &UpdateWorkspaceOptions{Name: String("x")},
			statusCode:  http.StatusNotFound,
			response:    `{"error": "not found"}`,
			wantErr:     true,
		},
		{
			name:        "unauthorized",
			workspaceID: 1,
			opts:        &UpdateWorkspaceOptions{Name: String("x")},
			statusCode:  http.StatusUnauthorized,
			response:    `{"error": "unauthorized"}`,
			wantErr:     true,
		},
		{
			name:        "server error",
			workspaceID: 1,
			opts:        &UpdateWorkspaceOptions{Name: String("x")},
			statusCode:  http.StatusInternalServerError,
			response:    `{"error": "internal server error"}`,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/%d", tt.workspaceID)
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
				if tt.opts.Name != nil {
					if got, ok := body["name"].(string); !ok || got != *tt.opts.Name {
						t.Errorf("name = %v, want %q", body["name"], *tt.opts.Name)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			workspace, _, err := client.Workspaces.UpdateWorkspace(context.Background(), tt.workspaceID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && workspace.ID != tt.wantID {
				t.Errorf("UpdateWorkspace() ID = %d, want %d", workspace.ID, tt.wantID)
			}
		})
	}
}
