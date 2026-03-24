package toggl

import (
	"context"
	"net/http"
	"testing"
)

const meJSON = `{
	"id": 1234567,
	"email": "user@example.com",
	"fullname": "Ada Lovelace",
	"timezone": "Europe/Amsterdam",
	"default_workspace_id": 100,
	"beginning_of_week": 1,
	"image_url": "https://assets.toggl.com/images/profile.png",
	"country_id": null,
	"has_password": true,
	"openid_enabled": false,
	"api_token": "abc123token",
	"at": "2024-03-01T09:00:00Z",
	"created_at": "2020-06-15T12:00:00Z",
	"updated_at": "2024-03-01T09:00:00Z"
}`

func TestMeService_GetMe(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response:   meJSON,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
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
				if r.URL.Path != "/api/v9/me" {
					t.Errorf("expected path /api/v9/me, got %s", r.URL.Path)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})

			client := testClient(t, handler)
			me, _, err := client.Me.GetMe(context.Background())

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetMe() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && me.ID != 1234567 {
				t.Errorf("GetMe() ID = %d, want 1234567", me.ID)
			}
		})
	}
}

func TestMeService_GetMe_Fields(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(meJSON))
	})

	client := testClient(t, handler)
	me, _, err := client.Me.GetMe(context.Background())
	if err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}

	if me.ID != 1234567 {
		t.Errorf("ID = %d, want 1234567", me.ID)
	}
	if me.Email != "user@example.com" {
		t.Errorf("Email = %q, want %q", me.Email, "user@example.com")
	}
	if me.Fullname != "Ada Lovelace" {
		t.Errorf("Fullname = %q, want %q", me.Fullname, "Ada Lovelace")
	}
	if me.Timezone != "Europe/Amsterdam" {
		t.Errorf("Timezone = %q, want %q", me.Timezone, "Europe/Amsterdam")
	}
	if me.DefaultWorkspaceID != 100 {
		t.Errorf("DefaultWorkspaceID = %d, want 100", me.DefaultWorkspaceID)
	}
	if me.BeginningOfWeek != 1 {
		t.Errorf("BeginningOfWeek = %d, want 1", me.BeginningOfWeek)
	}
	if me.ImageURL != "https://assets.toggl.com/images/profile.png" {
		t.Errorf("ImageURL = %q, want %q", me.ImageURL, "https://assets.toggl.com/images/profile.png")
	}
	if me.CountryID != nil {
		t.Errorf("CountryID = %v, want nil", me.CountryID)
	}
	if !me.HasPassword {
		t.Error("HasPassword = false, want true")
	}
	if me.OpenIDEnabled {
		t.Error("OpenIDEnabled = true, want false")
	}
	if me.APIToken != "abc123token" {
		t.Errorf("APIToken = %q, want %q", me.APIToken, "abc123token")
	}
	wantAt := "2024-03-01 09:00:00 +0000 UTC"
	if me.At.String() != wantAt {
		t.Errorf("At = %q, want %q", me.At.String(), wantAt)
	}
	wantCreatedAt := "2020-06-15 12:00:00 +0000 UTC"
	if me.CreatedAt.String() != wantCreatedAt {
		t.Errorf("CreatedAt = %q, want %q", me.CreatedAt.String(), wantCreatedAt)
	}
}
