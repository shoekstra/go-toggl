package toggl

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

const projectJSON = `{
	"id": 1,
	"name": "My Project",
	"workspace_id": 100,
	"client_id": null,
	"active": true,
	"is_private": true,
	"billable": null,
	"auto_estimates": null,
	"estimated_hours": null,
	"estimated_seconds": null,
	"actual_hours": null,
	"actual_seconds": null,
	"color": "#06aaf5",
	"rate": 0,
	"rate_last_updated": null,
	"currency": null,
	"template": null,
	"template_id": null,
	"fixed_fee": 0,
	"start_date": null,
	"end_date": null,
	"at": "2024-01-15T10:00:00Z",
	"created_at": "2024-01-15T10:00:00Z"
}`

func TestProjectsService_ListProjects(t *testing.T) {
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
			response:   "[" + projectJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "multiple projects",
			statusCode: http.StatusOK,
			response:   "[" + projectJSON + "," + projectJSON + "]",
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
			wantPath := "/api/v9/workspaces/100/projects"
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
			projects, _, err := client.Projects.ListProjects(context.Background(), 100, nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ListProjects() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(projects) != tt.wantCount {
				t.Errorf("ListProjects() count = %d, want %d", len(projects), tt.wantCount)
			}
		})
	}
}

func TestProjectsService_ListProjects_QueryParams(t *testing.T) {
	tests := []struct {
		name        string
		opts        *ListProjectsOptions
		wantParams  map[string]string
		wantNoParam []string
	}{
		{
			name:        "nil opts — no query params",
			opts:        nil,
			wantNoParam: []string{"active", "billable", "name", "since", "page", "per_page"},
		},
		{
			name:       "active both",
			opts:       &ListProjectsOptions{Active: String("both")},
			wantParams: map[string]string{"active": "both"},
		},
		{
			name:       "active false",
			opts:       &ListProjectsOptions{Active: String("false")},
			wantParams: map[string]string{"active": "false"},
		},
		{
			name:       "billable true",
			opts:       &ListProjectsOptions{Billable: Bool(true)},
			wantParams: map[string]string{"billable": "true"},
		},
		{
			name:       "name filter",
			opts:       &ListProjectsOptions{Name: String("My Project")},
			wantParams: map[string]string{"name": "My Project"},
		},
		{
			name:       "since",
			opts:       &ListProjectsOptions{Since: func() *int64 { v := int64(1700000000); return &v }()},
			wantParams: map[string]string{"since": "1700000000"},
		},
		{
			name:       "page and per_page",
			opts:       &ListProjectsOptions{Page: Int(2), PerPage: Int(50)},
			wantParams: map[string]string{"page": "2", "per_page": "50"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				for k, want := range tt.wantParams {
					if got := q.Get(k); got != want {
						t.Errorf("param %s = %q, want %q", k, got, want)
					}
				}
				for _, k := range tt.wantNoParam {
					if q.Has(k) {
						t.Errorf("unexpected param %s = %q", k, q.Get(k))
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("[]"))
			})

			client := testClient(t, handler)
			client.Projects.ListProjects(context.Background(), 100, tt.opts) //nolint:errcheck
		})
	}
}

func TestProjectsService_ListProjects_Fields(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[" + projectJSON + "]"))
	})

	client := testClient(t, handler)
	projects, _, err := client.Projects.ListProjects(context.Background(), 100, nil)
	if err != nil {
		t.Fatalf("ListProjects() error = %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("ListProjects() count = %d, want 1", len(projects))
	}

	p := projects[0]
	if p.ID != 1 {
		t.Errorf("ID = %d, want 1", p.ID)
	}
	if p.Name != "My Project" {
		t.Errorf("Name = %q, want %q", p.Name, "My Project")
	}
	if p.WorkspaceID != 100 {
		t.Errorf("WorkspaceID = %d, want 100", p.WorkspaceID)
	}
	if p.ClientID != nil {
		t.Errorf("ClientID = %v, want nil", p.ClientID)
	}
	if !p.Active {
		t.Errorf("Active = false, want true")
	}
	if !p.IsPrivate {
		t.Errorf("IsPrivate = false, want true")
	}
	if p.Billable != nil {
		t.Errorf("Billable = %v, want nil", p.Billable)
	}
	if p.Color != "#06aaf5" {
		t.Errorf("Color = %q, want %q", p.Color, "#06aaf5")
	}
}

func TestProjectsService_GetProject(t *testing.T) {
	tests := []struct {
		name       string
		projectID  int
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			projectID:  1,
			statusCode: http.StatusOK,
			response:   projectJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:       "not found",
			projectID:  999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantID:     0,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			projectID:  1,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantID:     0,
			wantErr:    true,
		},
		{
			name:       "server error",
			projectID:  1,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantID:     0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/projects/%d", tt.projectID)
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
			project, _, err := client.Projects.GetProject(context.Background(), 100, tt.projectID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetProject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && project.ID != tt.wantID {
				t.Errorf("GetProject() ID = %d, want %d", project.ID, tt.wantID)
			}
		})
	}
}

func TestProjectsService_CreateProject(t *testing.T) {
	tests := []struct {
		name       string
		opts       *CreateProjectOptions
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			opts:       &CreateProjectOptions{Name: "My Project"},
			statusCode: http.StatusOK,
			response:   projectJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name: "success with all options",
			opts: &CreateProjectOptions{
				Name:           "Full Project",
				Active:         Bool(true),
				Billable:       Bool(true),
				ClientID:       Int(42),
				Color:          String("#ff0000"),
				Currency:       String("EUR"),
				EndDate:        String("2024-12-31"),
				EstimatedHours: Int(100),
				FixedFee:       Float64(500.0),
				IsPrivate:      Bool(false),
				Rate:           Float64(75.0),
				StartDate:      String("2024-01-01"),
				Template:       Bool(false),
				TemplateID:     Int(5),
			},
			statusCode: http.StatusOK,
			response:   projectJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "empty name",
			opts:    &CreateProjectOptions{Name: ""},
			wantErr: true,
		},
		{
			name:       "unauthorized",
			opts:       &CreateProjectOptions{Name: "x"},
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       &CreateProjectOptions{Name: "x"},
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := "/api/v9/workspaces/100/projects"
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
				if got, ok := body["name"].(string); !ok || got != tt.opts.Name {
					t.Errorf("name = %v, want %q", body["name"], tt.opts.Name)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			project, _, err := client.Projects.CreateProject(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && project.ID != tt.wantID {
				t.Errorf("CreateProject() ID = %d, want %d", project.ID, tt.wantID)
			}
		})
	}
}

func TestProjectsService_UpdateProject(t *testing.T) {
	tests := []struct {
		name       string
		projectID  int
		opts       *UpdateProjectOptions
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			projectID:  1,
			opts:       &UpdateProjectOptions{Name: String("Renamed Project")},
			statusCode: http.StatusOK,
			response:   projectJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:      "success with all options",
			projectID: 1,
			opts: &UpdateProjectOptions{
				Name:           String("Updated"),
				Active:         Bool(false),
				Billable:       Bool(true),
				ClientID:       Int(42),
				Color:          String("#ff0000"),
				Currency:       String("EUR"),
				EndDate:        String("2024-12-31"),
				EstimatedHours: Int(200),
				FixedFee:       Float64(1000.0),
				IsPrivate:      Bool(true),
				Rate:           Float64(100.0),
				StartDate:      String("2024-06-01"),
				Template:       Bool(false),
			},
			statusCode: http.StatusOK,
			response:   projectJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:      "nil options",
			projectID: 1,
			opts:      nil,
			wantErr:   true,
		},
		{
			name:       "not found",
			projectID:  999,
			opts:       &UpdateProjectOptions{Name: String("x")},
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			projectID:  1,
			opts:       &UpdateProjectOptions{Name: String("x")},
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			projectID:  1,
			opts:       &UpdateProjectOptions{Name: String("x")},
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/projects/%d", tt.projectID)
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
			project, _, err := client.Projects.UpdateProject(context.Background(), 100, tt.projectID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateProject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && project.ID != tt.wantID {
				t.Errorf("UpdateProject() ID = %d, want %d", project.ID, tt.wantID)
			}
		})
	}
}

func TestProjectsService_DeleteProject(t *testing.T) {
	tests := []struct {
		name           string
		projectID      int
		teDeletionMode *string
		statusCode     int
		response       string
		wantErr        bool
	}{
		{
			name:       "success",
			projectID:  1,
			statusCode: http.StatusOK,
			response:   "",
			wantErr:    false,
		},
		{
			name:           "success with delete mode",
			projectID:      1,
			teDeletionMode: String("delete"),
			statusCode:     http.StatusOK,
			response:       "",
			wantErr:        false,
		},
		{
			name:           "success with unassign mode",
			projectID:      1,
			teDeletionMode: String("unassign"),
			statusCode:     http.StatusOK,
			response:       "",
			wantErr:        false,
		},
		{
			name:       "not found",
			projectID:  999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			projectID:  1,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			projectID:  1,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/projects/%d", tt.projectID)
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
				if tt.teDeletionMode != nil {
					if got := r.URL.Query().Get("teDeletionMode"); got != *tt.teDeletionMode {
						t.Errorf("teDeletionMode = %q, want %q", got, *tt.teDeletionMode)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			})

			client := testClient(t, handler)
			_, err := client.Projects.DeleteProject(context.Background(), 100, tt.projectID, tt.teDeletionMode)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteProject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
