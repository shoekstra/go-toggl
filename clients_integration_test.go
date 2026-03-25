//go:build integration

package toggl_test

import (
	"testing"

	toggl "github.com/shoekstra/go-toggl"
)

func TestIntegration_Clients_CRUD(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	// Create.
	name := uniqueName("client")
	created, _, err := client.Clients.CreateClient(ctx, wsID, &toggl.CreateClientOptions{
		Name:  name,
		Notes: toggl.String("integration test client"),
	})
	if err != nil {
		t.Fatalf("CreateClient: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("created client has ID=0")
	}
	t.Cleanup(func() {
		client.Clients.DeleteClient(ctx, wsID, created.ID) //nolint:errcheck
	})

	// Get.
	got, _, err := client.Clients.GetClient(ctx, wsID, created.ID)
	if err != nil {
		t.Fatalf("GetClient: %v", err)
	}
	if got.Name != name {
		t.Errorf("Name = %q, want %q", got.Name, name)
	}

	// Update.
	newName := uniqueName("client-renamed")
	updated, _, err := client.Clients.UpdateClient(ctx, wsID, created.ID, &toggl.UpdateClientOptions{
		Name: toggl.String(newName),
	})
	if err != nil {
		t.Fatalf("UpdateClient: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("updated Name = %q, want %q", updated.Name, newName)
	}

	// List — client should appear.
	clients, _, err := client.Clients.ListClients(ctx, wsID, nil)
	if err != nil {
		t.Fatalf("ListClients: %v", err)
	}
	found := false
	for _, c := range clients {
		if c.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created client %d not found in ListClients", created.ID)
	}

	// Delete.
	if _, err := client.Clients.DeleteClient(ctx, wsID, created.ID); err != nil {
		t.Fatalf("DeleteClient: %v", err)
	}
}
