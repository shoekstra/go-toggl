package toggl

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

const clientJSON = `{
	"id": 1,
	"name": "Acme Corp",
	"wid": 100,
	"notes": null,
	"archived": false,
	"at": "2024-01-15T10:00:00Z",
	"creator_id": 42,
	"external_reference": null,
	"permissions": ["read", "write"]
}`

const clientJSONFull = `{
	"id": 2,
	"name": "Beta Ltd",
	"wid": 100,
	"notes": "Important client",
	"archived": true,
	"at": "2024-06-01T12:30:00Z",
	"creator_id": 7,
	"external_reference": "EXT-001",
	"permissions": ["read"]
}`

func TestClientsService_ListClients(t *testing.T) {
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
			response:   "[" + clientJSON + "]",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "multiple clients",
			statusCode: http.StatusOK,
			response:   "[" + clientJSON + "," + clientJSON + "]",
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
			wantPath := "/api/v9/workspaces/100/clients"
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
			clients, _, err := client.Clients.ListClients(context.Background(), 100, nil)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ListClients() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(clients) != tt.wantCount {
				t.Errorf("ListClients() count = %d, want %d", len(clients), tt.wantCount)
			}
		})
	}
}

func TestClientsService_ListClients_QueryParams(t *testing.T) {
	tests := []struct {
		name        string
		opts        *ListClientsOptions
		wantParams  map[string]string
		wantNoParam []string
	}{
		{
			name:        "nil opts — no query params",
			opts:        nil,
			wantNoParam: []string{"status", "name"},
		},
		{
			name:       "status active",
			opts:       &ListClientsOptions{Status: String("active")},
			wantParams: map[string]string{"status": "active"},
		},
		{
			name:       "status archived",
			opts:       &ListClientsOptions{Status: String("archived")},
			wantParams: map[string]string{"status": "archived"},
		},
		{
			name:       "status both",
			opts:       &ListClientsOptions{Status: String("both")},
			wantParams: map[string]string{"status": "both"},
		},
		{
			name:       "name filter",
			opts:       &ListClientsOptions{Name: String("Acme")},
			wantParams: map[string]string{"name": "Acme"},
		},
		{
			name:       "status and name",
			opts:       &ListClientsOptions{Status: String("both"), Name: String("Corp")},
			wantParams: map[string]string{"status": "both", "name": "Corp"},
		},
		{
			name:       "pagination",
			opts:       &ListClientsOptions{Page: Int(2), PerPage: Int(25)},
			wantParams: map[string]string{"page": "2", "per_page": "25"},
		},
		{
			name:        "no pagination params when not set",
			opts:        &ListClientsOptions{},
			wantNoParam: []string{"page", "per_page"},
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
			if _, _, err := client.Clients.ListClients(context.Background(), 100, tt.opts); err != nil {
				t.Fatalf("ListClients() unexpected error: %v", err)
			}
		})
	}
}

func TestClientsService_ListClients_InvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		opts *ListClientsOptions
	}{
		{
			name: "invalid status",
			opts: &ListClientsOptions{Status: String("invalid")},
		},
		{
			name: "page zero",
			opts: &ListClientsOptions{Page: Int(0)},
		},
		{
			name: "page negative",
			opts: &ListClientsOptions{Page: Int(-1)},
		},
		{
			name: "per_page zero",
			opts: &ListClientsOptions{PerPage: Int(0)},
		},
		{
			name: "per_page negative",
			opts: &ListClientsOptions{PerPage: Int(-5)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("handler should not be called for invalid options")
				w.WriteHeader(http.StatusBadRequest)
			})

			c := testClient(t, handler)
			_, _, err := c.Clients.ListClients(context.Background(), 100, tt.opts)
			if err == nil {
				t.Fatalf("ListClients() expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestClientsService_ListClients_Fields(t *testing.T) {
	t.Run("null optional fields", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[" + clientJSON + "]"))
		})

		c := testClient(t, handler)
		clients, _, err := c.Clients.ListClients(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("ListClients() error = %v", err)
		}
		if len(clients) != 1 {
			t.Fatalf("ListClients() count = %d, want 1", len(clients))
		}

		tc := clients[0]
		if tc.ID != 1 {
			t.Errorf("ID = %d, want 1", tc.ID)
		}
		if tc.Name != "Acme Corp" {
			t.Errorf("Name = %q, want %q", tc.Name, "Acme Corp")
		}
		if tc.WorkspaceID != 100 {
			t.Errorf("WorkspaceID = %d, want 100", tc.WorkspaceID)
		}
		if tc.Notes != nil {
			t.Errorf("Notes = %v, want nil", tc.Notes)
		}
		if tc.Archived {
			t.Errorf("Archived = true, want false")
		}
		if tc.CreatorID != 42 {
			t.Errorf("CreatorID = %d, want 42", tc.CreatorID)
		}
		if tc.ExternalReference != nil {
			t.Errorf("ExternalReference = %v, want nil", tc.ExternalReference)
		}
		if len(tc.Permissions) != 2 {
			t.Errorf("Permissions count = %d, want 2", len(tc.Permissions))
		}
		wantAt := "2024-01-15 10:00:00 +0000 UTC"
		if tc.At.String() != wantAt {
			t.Errorf("At = %q, want %q", tc.At.String(), wantAt)
		}
	})

	t.Run("populated optional fields", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[" + clientJSONFull + "]"))
		})

		c := testClient(t, handler)
		clients, _, err := c.Clients.ListClients(context.Background(), 100, nil)
		if err != nil {
			t.Fatalf("ListClients() error = %v", err)
		}
		if len(clients) != 1 {
			t.Fatalf("ListClients() count = %d, want 1", len(clients))
		}

		tc := clients[0]
		if tc.Notes == nil || *tc.Notes != "Important client" {
			t.Errorf("Notes = %v, want %q", tc.Notes, "Important client")
		}
		if tc.ExternalReference == nil || *tc.ExternalReference != "EXT-001" {
			t.Errorf("ExternalReference = %v, want %q", tc.ExternalReference, "EXT-001")
		}
		if !tc.Archived {
			t.Errorf("Archived = false, want true")
		}
		wantAt := "2024-06-01 12:30:00 +0000 UTC"
		if tc.At.String() != wantAt {
			t.Errorf("At = %q, want %q", tc.At.String(), wantAt)
		}
	})
}

func TestClientsService_GetClient(t *testing.T) {
	tests := []struct {
		name       string
		clientID   int
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			clientID:   1,
			statusCode: http.StatusOK,
			response:   clientJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:       "not found",
			clientID:   999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantID:     0,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			clientID:   1,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantID:     0,
			wantErr:    true,
		},
		{
			name:       "server error",
			clientID:   1,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantID:     0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/clients/%d", tt.clientID)
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

			c := testClient(t, handler)
			tc, _, err := c.Clients.GetClient(context.Background(), 100, tt.clientID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tc.ID != tt.wantID {
				t.Errorf("GetClient() ID = %d, want %d", tc.ID, tt.wantID)
			}
		})
	}
}

func TestClientsService_CreateClient(t *testing.T) {
	tests := []struct {
		name       string
		opts       *CreateClientOptions
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			opts:       &CreateClientOptions{Name: "Acme Corp"},
			statusCode: http.StatusOK,
			response:   clientJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name: "success with all options",
			opts: &CreateClientOptions{
				Name:              "Acme Corp",
				Notes:             String("Important client"),
				ExternalReference: String("EXT-001"),
			},
			statusCode: http.StatusOK,
			response:   clientJSON,
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
			opts:    &CreateClientOptions{Name: ""},
			wantErr: true,
		},
		{
			name:       "unauthorized",
			opts:       &CreateClientOptions{Name: "x"},
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			opts:       &CreateClientOptions{Name: "x"},
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := "/api/v9/workspaces/100/clients"
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
				if tt.opts.Notes != nil {
					if got, ok := body["notes"].(string); !ok || got != *tt.opts.Notes {
						t.Errorf("notes = %v, want %q", body["notes"], *tt.opts.Notes)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				if tt.opts.ExternalReference != nil {
					if got, ok := body["external_reference"].(string); !ok || got != *tt.opts.ExternalReference {
						t.Errorf("external_reference = %v, want %q", body["external_reference"], *tt.opts.ExternalReference)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			c := testClient(t, handler)
			tc, _, err := c.Clients.CreateClient(context.Background(), 100, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tc.ID != tt.wantID {
				t.Errorf("CreateClient() ID = %d, want %d", tc.ID, tt.wantID)
			}
		})
	}
}

func TestClientsService_UpdateClient(t *testing.T) {
	tests := []struct {
		name       string
		clientID   int
		opts       *UpdateClientOptions
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success",
			clientID:   1,
			opts:       &UpdateClientOptions{Name: String("Renamed Corp")},
			statusCode: http.StatusOK,
			response:   clientJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:     "success with all options",
			clientID: 1,
			opts: &UpdateClientOptions{
				Name:              String("Updated Corp"),
				Notes:             String("Updated notes"),
				ExternalReference: String("EXT-002"),
			},
			statusCode: http.StatusOK,
			response:   clientJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:     "nil options",
			clientID: 1,
			opts:     nil,
			wantErr:  true,
		},
		{
			name:     "empty options",
			clientID: 1,
			opts:     &UpdateClientOptions{},
			wantErr:  true,
		},
		{
			name:       "not found",
			clientID:   999,
			opts:       &UpdateClientOptions{Name: String("x")},
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			clientID:   1,
			opts:       &UpdateClientOptions{Name: String("x")},
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			clientID:   1,
			opts:       &UpdateClientOptions{Name: String("x")},
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/clients/%d", tt.clientID)
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
				if tt.opts.Notes != nil {
					if got, ok := body["notes"].(string); !ok || got != *tt.opts.Notes {
						t.Errorf("notes = %v, want %q", body["notes"], *tt.opts.Notes)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				if tt.opts.ExternalReference != nil {
					if got, ok := body["external_reference"].(string); !ok || got != *tt.opts.ExternalReference {
						t.Errorf("external_reference = %v, want %q", body["external_reference"], *tt.opts.ExternalReference)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			c := testClient(t, handler)
			tc, _, err := c.Clients.UpdateClient(context.Background(), 100, tt.clientID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tc.ID != tt.wantID {
				t.Errorf("UpdateClient() ID = %d, want %d", tc.ID, tt.wantID)
			}
		})
	}
}

func TestClientsService_DeleteClient(t *testing.T) {
	tests := []struct {
		name       string
		clientID   int
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			clientID:   1,
			statusCode: http.StatusOK,
			response:   "",
			wantErr:    false,
		},
		{
			name:       "not found",
			clientID:   999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			clientID:   1,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			clientID:   1,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/clients/%d", tt.clientID)
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

			c := testClient(t, handler)
			_, err := c.Clients.DeleteClient(context.Background(), 100, tt.clientID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientsService_ArchiveClient(t *testing.T) {
	tests := []struct {
		name        string
		clientID    int
		statusCode  int
		response    string
		wantItems   []int
		wantErr     bool
	}{
		{
			name:       "success — no archived projects",
			clientID:   1,
			statusCode: http.StatusOK,
			response:   `{"items": []}`,
			wantItems:  []int{},
			wantErr:    false,
		},
		{
			name:       "success — with archived projects",
			clientID:   1,
			statusCode: http.StatusOK,
			response:   `{"items": [10, 20, 30]}`,
			wantItems:  []int{10, 20, 30},
			wantErr:    false,
		},
		{
			name:       "not found",
			clientID:   999,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			clientID:   1,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			clientID:   1,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/clients/%d/archive", tt.clientID)
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
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			c := testClient(t, handler)
			ar, _, err := c.Clients.ArchiveClient(context.Background(), 100, tt.clientID)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ArchiveClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(ar.Items) != len(tt.wantItems) {
					t.Fatalf("ArchiveClient() items count = %d, want %d", len(ar.Items), len(tt.wantItems))
				}
				for i, id := range ar.Items {
					if id != tt.wantItems[i] {
						t.Errorf("ArchiveClient() items[%d] = %d, want %d", i, id, tt.wantItems[i])
					}
				}
			}
		})
	}
}

func TestClientsService_RestoreClient(t *testing.T) {
	tests := []struct {
		name       string
		clientID   int
		opts       *RestoreClientOptions
		statusCode int
		response   string
		wantID     int
		wantErr    bool
	}{
		{
			name:       "success — nil opts",
			clientID:   1,
			opts:       nil,
			statusCode: http.StatusOK,
			response:   clientJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:       "success — restore all projects",
			clientID:   1,
			opts:       &RestoreClientOptions{RestoreAllProjects: Bool(true)},
			statusCode: http.StatusOK,
			response:   clientJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:       "success — specific projects",
			clientID:   1,
			opts:       &RestoreClientOptions{Projects: []int{10, 20}},
			statusCode: http.StatusOK,
			response:   clientJSON,
			wantID:     1,
			wantErr:    false,
		},
		{
			name:       "not found",
			clientID:   999,
			opts:       nil,
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			clientID:   1,
			opts:       nil,
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			clientID:   1,
			opts:       nil,
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantPath := fmt.Sprintf("/api/v9/workspaces/100/clients/%d/restore", tt.clientID)
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
				if tt.opts != nil {
					body := assertBody(t, r)
					if tt.opts.RestoreAllProjects != nil {
						if got, ok := body["restore_all_projects"].(bool); !ok || got != *tt.opts.RestoreAllProjects {
							t.Errorf("restore_all_projects = %v, want %v", body["restore_all_projects"], *tt.opts.RestoreAllProjects)
							w.WriteHeader(http.StatusBadRequest)
							return
						}
					}
					if len(tt.opts.Projects) > 0 {
						raw, ok := body["projects"].([]interface{})
						if !ok || len(raw) != len(tt.opts.Projects) {
							t.Errorf("projects = %v, want %v", body["projects"], tt.opts.Projects)
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						for i, v := range raw {
							if int(v.(float64)) != tt.opts.Projects[i] {
								t.Errorf("projects[%d] = %v, want %d", i, v, tt.opts.Projects[i])
								w.WriteHeader(http.StatusBadRequest)
								return
							}
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			c := testClient(t, handler)
			tc, _, err := c.Clients.RestoreClient(context.Background(), 100, tt.clientID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("RestoreClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tc.ID != tt.wantID {
				t.Errorf("RestoreClient() ID = %d, want %d", tc.ID, tt.wantID)
			}
		})
	}
}
