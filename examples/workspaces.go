//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	toggl "github.com/shoekstra/go-toggl"
)

func main() {
	client, err := toggl.NewClient(os.Getenv("TOGGL_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List all workspaces for the authenticated user.
	workspaces, _, err := client.Workspaces.ListWorkspaces(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, ws := range workspaces {
		fmt.Printf("Workspace %d: %s\n", ws.ID, ws.Name)
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
		return
	}

	wsID := workspaces[0].ID

	// Get a single workspace.
	ws, _, err := client.Workspaces.GetWorkspace(ctx, wsID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nWorkspace:        %s\n", ws.Name)
	fmt.Printf("Currency:         %s\n", ws.DefaultCurrency)
	fmt.Printf("Active projects:  %d\n", ws.ActiveProjectCount)

	// Update a workspace — only when TOGGL_WORKSPACE_ID is explicitly set to
	// avoid accidentally renaming the wrong workspace. Restores the original
	// name afterward.
	if wsIDEnv := os.Getenv("TOGGL_WORKSPACE_ID"); wsIDEnv != "" {
		originalName := ws.Name
		updated, _, err := client.Workspaces.UpdateWorkspace(ctx, wsID, &toggl.UpdateWorkspaceOptions{
			Name: toggl.String("My Workspace"),
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nUpdated name: %s\n", updated.Name)

		if _, _, err := client.Workspaces.UpdateWorkspace(ctx, wsID, &toggl.UpdateWorkspaceOptions{
			Name: toggl.String(originalName),
		}); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Restored name: %s\n", originalName)
	} else {
		fmt.Println("\nSkipping UpdateWorkspace: set TOGGL_WORKSPACE_ID to enable.")
	}
}
