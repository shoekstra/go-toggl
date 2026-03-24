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

	// List active projects.
	projects, _, err := client.Projects.ListProjects(ctx, wsID, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Active projects: %d\n", len(projects))
	for _, p := range projects {
		fmt.Printf("  %d: %s\n", p.ID, p.Name)
	}

	// Create a project.
	created, _, err := client.Projects.CreateProject(ctx, wsID, &toggl.CreateProjectOptions{
		Name:  "Example Project",
		Color: toggl.String("#06aaf5"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCreated: %s (ID: %d)\n", created.Name, created.ID)

	// Get the project by ID.
	project, _, err := client.Projects.GetProject(ctx, wsID, created.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched: %s (color: %s)\n", project.Name, project.Color)

	// Update the project.
	updated, _, err := client.Projects.UpdateProject(ctx, wsID, created.ID, &toggl.UpdateProjectOptions{
		Name: toggl.String("Example Project (renamed)"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated: %s\n", updated.Name)

	// Delete the project.
	if _, err := client.Projects.DeleteProject(ctx, wsID, created.ID, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted.")
}
