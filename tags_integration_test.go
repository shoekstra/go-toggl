//go:build integration

package toggl_test

import (
	"testing"

	toggl "github.com/shoekstra/go-toggl"
)

// TestIntegration_Tags_ListFindsCreatedTag verifies that a tag created via
// CreateTag is immediately visible in ListTags. Because the Toggl API does not
// support GET /workspaces/{id}/tags/{id} (returns HTTP 405), ListTags is the
// correct way to look up a tag by ID.
func TestIntegration_Tags_ListFindsCreatedTag(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	created, _, err := client.Tags.CreateTag(ctx, wsID, &toggl.CreateTagOptions{
		Name: uniqueName("tag-probe"),
	})
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	t.Cleanup(func() {
		client.Tags.DeleteTag(ctx, wsID, created.ID) //nolint:errcheck
	})

	tags, _, err := client.Tags.ListTags(ctx, wsID, nil)
	if err != nil {
		t.Fatalf("ListTags: %v", err)
	}
	found := false
	for _, tag := range tags {
		if tag.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created tag %d not found in ListTags", created.ID)
	}
}

func TestIntegration_Tags_CRUD(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	// Create.
	name := uniqueName("tag")
	created, _, err := client.Tags.CreateTag(ctx, wsID, &toggl.CreateTagOptions{Name: name})
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("created tag has ID=0")
	}
	t.Cleanup(func() {
		client.Tags.DeleteTag(ctx, wsID, created.ID) //nolint:errcheck
	})

	if created.Name != name {
		t.Errorf("Name = %q, want %q", created.Name, name)
	}

	// Update.
	newName := uniqueName("tag-renamed")
	updated, _, err := client.Tags.UpdateTag(ctx, wsID, created.ID, &toggl.UpdateTagOptions{Name: newName})
	if err != nil {
		t.Fatalf("UpdateTag: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("updated Name = %q, want %q", updated.Name, newName)
	}

	// Delete.
	if _, err := client.Tags.DeleteTag(ctx, wsID, created.ID); err != nil {
		t.Fatalf("DeleteTag: %v", err)
	}
}
