//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	toggl "github.com/shoekstra/go-toggl"
)

func main() {
	client, err := toggl.NewClient(os.Getenv("TOGGL_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	wsID, err := strconv.Atoi(os.Getenv("TOGGL_WORKSPACE_ID"))
	if err != nil {
		log.Fatal("TOGGL_WORKSPACE_ID must be set to a valid workspace ID (run examples/workspaces.go to find yours)")
	}

	ctx := context.Background()

	// List active clients.
	clients, _, err := client.Clients.ListClients(ctx, wsID, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Active clients: %d\n", len(clients))
	for _, c := range clients {
		fmt.Printf("  %d: %s\n", c.ID, c.Name)
	}

	// Create a client.
	created, _, err := client.Clients.CreateClient(ctx, wsID, &toggl.CreateClientOptions{
		Name:  "Example Corp",
		Notes: toggl.String("Created by go-toggl example"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCreated: %s (ID: %d)\n", created.Name, created.ID)

	// Get the client by ID.
	wc, _, err := client.Clients.GetClient(ctx, wsID, created.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched: %s\n", wc.Name)

	// Update the client.
	updated, _, err := client.Clients.UpdateClient(ctx, wsID, created.ID, &toggl.UpdateClientOptions{
		Name: toggl.String("Example Corp (renamed)"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated: %s\n", updated.Name)

	// Archive the client (premium workspaces only).
	ar, _, err := client.Clients.ArchiveClient(ctx, wsID, created.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Archived. Related project IDs: %v\n", ar.Items)

	// Restore the archived client.
	restored, _, err := client.Clients.RestoreClient(ctx, wsID, created.ID, &toggl.RestoreClientOptions{
		RestoreAllProjects: toggl.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Restored: %s\n", restored.Name)

	// Delete the client.
	if _, err := client.Clients.DeleteClient(ctx, wsID, created.ID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted.")
}
