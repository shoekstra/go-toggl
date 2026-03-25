//go:build integration

package toggl_test

import (
	"testing"

	toggl "github.com/shoekstra/go-toggl"
)

func TestIntegration_Projects_CRUD(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	// Create.
	name := uniqueName("project")
	created, _, err := client.Projects.CreateProject(ctx, wsID, &toggl.CreateProjectOptions{
		Name:   name,
		Color:  toggl.String("#06aaf5"),
		Active: toggl.Bool(true),
	})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("created project has ID=0")
	}
	t.Cleanup(func() {
		client.Projects.DeleteProject(ctx, wsID, created.ID, nil) //nolint:errcheck
	})

	// Get.
	got, _, err := client.Projects.GetProject(ctx, wsID, created.ID)
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if got.Name != name {
		t.Errorf("Name = %q, want %q", got.Name, name)
	}
	if got.Color != "#06aaf5" {
		t.Errorf("Color = %q, want %q", got.Color, "#06aaf5")
	}

	// Update.
	newName := uniqueName("project-renamed")
	updated, _, err := client.Projects.UpdateProject(ctx, wsID, created.ID, &toggl.UpdateProjectOptions{
		Name: toggl.String(newName),
	})
	if err != nil {
		t.Fatalf("UpdateProject: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("updated Name = %q, want %q", updated.Name, newName)
	}

	// List — project should appear.
	projects, _, err := client.Projects.ListProjects(ctx, wsID, nil)
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	found := false
	for _, p := range projects {
		if p.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created project %d not found in ListProjects", created.ID)
	}

	// Delete.
	if _, err := client.Projects.DeleteProject(ctx, wsID, created.ID, nil); err != nil {
		t.Fatalf("DeleteProject: %v", err)
	}
}
