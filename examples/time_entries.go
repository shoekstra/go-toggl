//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

	// List recent time entries (adjust dates as needed; Toggl limits queries to ~3 months back).
	entries, _, err := client.TimeEntries.ListTimeEntries(ctx, &toggl.ListTimeEntriesOptions{
		StartDate: toggl.String("2026-01-01"),
		EndDate:   toggl.String("2026-03-01"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Time entries Jan–Mar 2026: %d\n", len(entries))
	for _, e := range entries {
		desc := ""
		if e.Description != nil {
			desc = *e.Description
		}
		fmt.Printf("  %d: %s (%ds)\n", e.ID, desc, e.Duration)
	}

	// Get the currently running time entry. The API returns null (not an error)
	// when no entry is running, so check the ID to distinguish the two cases.
	running, _, err := client.TimeEntries.GetRunningTimeEntry(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if running == nil || running.ID == 0 {
		fmt.Println("\nNo entry currently running.")
	} else {
		fmt.Printf("\nRunning entry: %d (started %s)\n", running.ID, running.Start.Format(time.RFC3339))
	}

	// Start a new time entry.
	started, _, err := client.TimeEntries.StartTimeEntry(ctx, wsID, &toggl.CreateTimeEntryOptions{
		Start:       time.Now(),
		Description: toggl.String("Working on feature X"),
		Tags:        []string{"backend"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nStarted entry: %d\n", started.ID)

	// Stop the running entry.
	stopped, _, err := client.TimeEntries.StopTimeEntry(ctx, wsID, started.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Stopped. Duration: %ds\n", stopped.Duration)

	// Update the time entry.
	updated, _, err := client.TimeEntries.UpdateTimeEntry(ctx, wsID, stopped.ID, &toggl.UpdateTimeEntryOptions{
		Description: toggl.String("Working on feature X (updated)"),
		Billable:    toggl.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated: %s\n", *updated.Description)

	// Delete the time entry.
	if _, err := client.TimeEntries.DeleteTimeEntry(ctx, wsID, stopped.ID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted.")
}
