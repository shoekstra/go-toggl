//go:build integration

package toggl_test

import (
	"testing"

	toggl "github.com/shoekstra/go-toggl"
)

func TestIntegration_Workspaces_ListWorkspaces(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationCtx(t)

	workspaces, _, err := client.Workspaces.ListWorkspaces(ctx)
	if err != nil {
		t.Fatalf("ListWorkspaces: %v", err)
	}
	if len(workspaces) == 0 {
		t.Fatal("expected at least one workspace")
	}
	for _, ws := range workspaces {
		if ws.ID == 0 {
			t.Errorf("workspace has ID=0")
		}
		if ws.Name == "" {
			t.Errorf("workspace %d has empty name", ws.ID)
		}
	}
}

func TestIntegration_Workspaces_GetWorkspace(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	ws, _, err := client.Workspaces.GetWorkspace(ctx, wsID)
	if err != nil {
		t.Fatalf("GetWorkspace: %v", err)
	}
	if ws.ID != wsID {
		t.Errorf("ID = %d, want %d", ws.ID, wsID)
	}
	if ws.Name == "" {
		t.Error("Name is empty")
	}
}

func TestIntegration_Workspaces_UpdateWorkspace(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	// Read the current name so we can restore it.
	ws, _, err := client.Workspaces.GetWorkspace(ctx, wsID)
	if err != nil {
		t.Fatalf("GetWorkspace: %v", err)
	}
	originalName := ws.Name

	t.Cleanup(func() {
		if _, _, err := client.Workspaces.UpdateWorkspace(ctx, wsID, &toggl.UpdateWorkspaceOptions{
			Name: toggl.String(originalName),
		}); err != nil {
			t.Errorf("cleanup: failed to restore workspace name: %v", err)
		}
	})

	expectedName := uniqueName("ws")
	updated, _, err := client.Workspaces.UpdateWorkspace(ctx, wsID, &toggl.UpdateWorkspaceOptions{
		Name: toggl.String(expectedName),
	})
	if err != nil {
		t.Fatalf("UpdateWorkspace: %v", err)
	}
	if updated.Name != expectedName {
		t.Errorf("Name = %q, want %q", updated.Name, expectedName)
	}
}
