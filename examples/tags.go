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

	// List tags, optionally filtered by name.
	tags, _, err := client.Tags.ListTags(ctx, wsID, &toggl.ListTagsOptions{
		Search: toggl.String("backend"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tags matching 'backend': %d\n", len(tags))
	for _, t := range tags {
		fmt.Printf("  %d: %s\n", t.ID, t.Name)
	}

	// Create a tag.
	created, _, err := client.Tags.CreateTag(ctx, wsID, &toggl.CreateTagOptions{
		Name: "example-tag",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nCreated: %s (ID: %d)\n", created.Name, created.ID)

	// Rename the tag.
	updated, _, err := client.Tags.UpdateTag(ctx, wsID, created.ID, &toggl.UpdateTagOptions{
		Name: "example-tag-renamed",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated: %s\n", updated.Name)

	// Delete the tag.
	if _, err := client.Tags.DeleteTag(ctx, wsID, created.ID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted.")
}
