package toggl

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

const tagJSON = `{
	"id": 1,
	"name": "backend",
	"workspace_id": 100,
	"creator_id": 42,
	"at": "2024-01-15T10:00:00Z",
	"deleted_at": null,
	"integration_ext_id": null,
	"integration_ext_type": null,
	"permissions": ["read", "write"]
}`

const tagJSONFull = `{
	"id": 2,
	"name": "frontend",
	"workspace_id": 100,
	"creator_id": 7,
	"at": "2024-06-01T12:30:00Z",
	"deleted_at": "2024-09-01T08:00:00Z",
	"integration_ext_id": "EXT-TAG-001",
	"integration_ext_type": "jira",
	"permissions": ["read"]
}`

func TestTagsService_ListTags(t *testing.T) {
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
			response:   "[" + tagJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "multiple tags",
			statusCode: http.StatusOK,
			response:   "[" + tagJSON + "," + tagJSON + "]",
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
			wantPath := "/api/v9/workspaces/100/tags"
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
			tags, _, err := client.Tags.ListTags(context.Background(), 100, nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ListTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(tags) != tt.wantCount {
				t.Errorf("ListTags() count = %d, want %d", len(tags), tt.wantCount)
			}
		})
	}
}

func TestTagsService_ListTags_QueryParams(t *testing.T) {
	tests := []struct {
		name        string
		opts        *ListTagsOptions
		wantParams  map[string]string
		wantNoParam []string
	}{
		{
			name:        "nil opts — no query params",
			opts:        nil,
			wantNoParam: []string{"search", "page", "per_page"},
		},
		{
			name:       "search",
			opts:       &ListTagsOptions{Search: String("back")},
			wantParams: map[string]string{"search": "back"},
		},
		{
			name:       "page",
			opts:       &ListTagsOptions{Page: Int(2)},
			wantParams: map[string]string{"page": "2"},
		},
		{
			name:       "per_page",
			opts:       &ListTagsOptions{PerPage: Int(50)},
			wantParams: map[string]string{"per_page": "50"},
		},
		{
			name:       "all params",
			opts:       &ListTagsOptions{Search: String("api"), Page: Int(3), PerPage: Int(25)},
			wantParams: map[string]string{"search": "api", "page": "3", "per_page": "25"},
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
			if _, _, err := client.Tags.ListTags(context.Background(), 100, tt.opts); err != nil {
				t.Fatalf("ListTags() unexpected error: %v", err)
			}
		})
	}
}

func TestTagsService_ListTags_Fields(t *testing.T) {
	t.Run("null optional fields", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[" + tagJSON + "]"))
		})

		client := testClient(t, handler)
		tags, _, err := client.Tags.ListTags(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("ListTags() error = %v", err)
		}
		if len(tags) != 1 {
			t.Fatalf("ListTags() count = %d, want 1", len(tags))
		}

		tag := tags[0]
		if tag.ID != 1 {
			t.Errorf("ID = %d, want 1", tag.ID)
		}
		if tag.Name != "backend" {
			t.Errorf("Name = %q, want %q", tag.Name, "backend")
		}
		if tag.WorkspaceID != 100 {
			t.Errorf("WorkspaceID = %d, want 100", tag.WorkspaceID)
		}
		if tag.Workspace != 100 {
			t.Errorf("Workspace = %d, want 100", tag.Workspace)
		}
		if tag.CreatorID != 42 {
			t.Errorf("CreatorID = %d, want 42", tag.CreatorID)
		}
		wantAt := "2024-01-15 10:00:00 +0000 UTC"
		if tag.At.String() != wantAt {
			t.Errorf("At = %q, want %q", tag.At.String(), wantAt)
		}
		if tag.DeletedAt != nil {
			t.Errorf("DeletedAt = %v, want nil", tag.DeletedAt)
		}
		if tag.IntegrationExtID != nil {
			t.Errorf("IntegrationExtID = %v, want nil", tag.IntegrationExtID)
		}
		if tag.IntegrationExtType != nil {
			t.Errorf("IntegrationExtType = %v, want nil", tag.IntegrationExtType)
		}
		if len(tag.Permissions) != 2 {
			t.Errorf("Permissions count = %d, want 2", len(tag.Permissions))
		}
	})

	t.Run("populated optional fields", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[" + tagJSONFull + "]"))
		})

		client := testClient(t, handler)
		tags, _, err := client.Tags.ListTags(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("ListTags() error = %v", err)
		}
		if len(tags) != 1 {
			t.Fatalf("ListTags() count = %d, want 1", len(tags))
		}

		tag := tags[0]
		if tag.DeletedAt == nil {
			t.Fatal("DeletedAt = nil, want non-nil")
		}
		wantDeletedAt := "2024-09-01 08:00:00 +0000 UTC"
		if tag.DeletedAt.String() != wantDeletedAt {
			t.Errorf("DeletedAt = %q, want %q", tag.DeletedAt.String(), wantDeletedAt)
		}
		if tag.IntegrationExtID == nil || *tag.IntegrationExtID != "EXT-TAG-001" {
			t.Errorf("IntegrationExtID = %v, want %q", tag.IntegrationExtID, "EXT-TAG-001")
		}
		if tag.IntegrationExtType == nil || *tag.IntegrationExtType != "jira" {
			t.Errorf("IntegrationExtType = %v, want %q", tag.IntegrationExtType, "jira")
		}
		wantAt := "2024-06-01 12:30:00 +0000 UTC"
		if tag.At.String() != wantAt {
			t.Errorf("At = %q, want %q", tag.At.String(), wantAt)
		}
	})
}

func TestTagsService_GetTag(t *testing.T) {
	tests := []struct {
		name       string
		tagID      int
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			tagID:      1,
			statusCode: http.StatusOK,
			response:   tagJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:       "not found",
			tagID:      999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			tagID:      1,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			tagID:      1,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/tags/%d", tt.tagID)
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
			tag, _, err := client.Tags.GetTag(context.Background(), 100, tt.tagID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetTag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tag.ID != tt.wantID {
				t.Errorf("GetTag() ID = %d, want %d", tag.ID, tt.wantID)
			}
		})
	}
}

func TestTagsService_GetTag_Fields(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tagJSONFull))
	})

	client := testClient(t, handler)
	tag, _, err := client.Tags.GetTag(context.Background(), 100, 2)
	if err != nil {
		t.Fatalf("GetTag() error = %v", err)
	}

	if tag.ID != 2 {
		t.Errorf("ID = %d, want 2", tag.ID)
	}
	if tag.Name != "frontend" {
		t.Errorf("Name = %q, want %q", tag.Name, "frontend")
	}
	if tag.WorkspaceID != 100 {
		t.Errorf("WorkspaceID = %d, want 100", tag.WorkspaceID)
	}
	if tag.DeletedAt == nil {
		t.Fatal("DeletedAt = nil, want non-nil")
	}
	if tag.IntegrationExtID == nil || *tag.IntegrationExtID != "EXT-TAG-001" {
		t.Errorf("IntegrationExtID = %v, want %q", tag.IntegrationExtID, "EXT-TAG-001")
	}
	if tag.IntegrationExtType == nil || *tag.IntegrationExtType != "jira" {
		t.Errorf("IntegrationExtType = %v, want %q", tag.IntegrationExtType, "jira")
	}
}

func TestTagsService_CreateTag(t *testing.T) {
	tests := []struct {
		name       string
		opts       *CreateTagOptions
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			opts:       &CreateTagOptions{Name: "backend"},
			statusCode: http.StatusOK,
			response:   tagJSON,
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
			opts:    &CreateTagOptions{Name: ""},
			wantErr: true,
		},
		{
			name:       "unauthorized",
			opts:       &CreateTagOptions{Name: "x"},
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       &CreateTagOptions{Name: "x"},
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := "/api/v9/workspaces/100/tags"
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
			tag, _, err := client.Tags.CreateTag(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tag.ID != tt.wantID {
				t.Errorf("CreateTag() ID = %d, want %d", tag.ID, tt.wantID)
			}
		})
	}
}

func TestTagsService_UpdateTag(t *testing.T) {
	tests := []struct {
		name       string
		tagID      int
		opts       *UpdateTagOptions
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			tagID:      1,
			opts:       &UpdateTagOptions{Name: "renamed"},
			statusCode: http.StatusOK,
			response:   tagJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:    "nil options",
			tagID:   1,
			opts:    nil,
			wantErr: true,
		},
		{
			name:    "empty name",
			tagID:   1,
			opts:    &UpdateTagOptions{Name: ""},
			wantErr: true,
		},
		{
			name:       "not found",
			tagID:      999,
			opts:       &UpdateTagOptions{Name: "x"},
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			tagID:      1,
			opts:       &UpdateTagOptions{Name: "x"},
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			tagID:      1,
			opts:       &UpdateTagOptions{Name: "x"},
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/tags/%d", tt.tagID)
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
			tag, _, err := client.Tags.UpdateTag(context.Background(), 100, tt.tagID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tag.ID != tt.wantID {
				t.Errorf("UpdateTag() ID = %d, want %d", tag.ID, tt.wantID)
			}
		})
	}
}

func TestTagsService_DeleteTag(t *testing.T) {
	tests := []struct {
		name       string
		tagID      int
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			tagID:      1,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "not found",
			tagID:      999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			tagID:      1,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			tagID:      1,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/tags/%d", tt.tagID)
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
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					w.Write([]byte(tt.response))
				}
			})

			client := testClient(t, handler)
			_, err := client.Tags.DeleteTag(context.Background(), 100, tt.tagID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
