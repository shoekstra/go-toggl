//go:build integration

package toggl_test

import (
	"testing"
)

func TestIntegration_Me_GetMe(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationCtx(t)

	me, _, err := client.Me.GetMe(ctx)
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}
	if me == nil {
		t.Fatal("GetMe returned nil")
	}

	if me.ID == 0 {
		t.Error("ID = 0, want non-zero")
	}
	if me.Email == "" {
		t.Error("Email is empty")
	}
	if me.DefaultWorkspaceID == 0 {
		t.Error("DefaultWorkspaceID = 0, want non-zero")
	}
}
