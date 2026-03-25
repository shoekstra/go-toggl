//go:build integration

package toggl_test

import (
	"testing"
	"time"

	toggl "github.com/shoekstra/go-toggl"
)

func TestIntegration_TimeEntries_GetRunningTimeEntry_NilWhenNoneRunning(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationCtx(t)

	// The API returns null/200 when nothing is running; the client must
	// return (nil, resp, nil) rather than a zero-value TimeEntry.
	entry, _, err := client.TimeEntries.GetRunningTimeEntry(ctx)
	if err != nil {
		t.Fatalf("GetRunningTimeEntry: %v", err)
	}
	// If something happens to be running, skip rather than fail — this test
	// is specifically about the no-entry-running case.
	if entry != nil && entry.ID != 0 {
		t.Skip("a time entry is currently running; cannot test the nil-response path")
	}
	if entry != nil {
		t.Errorf("GetRunningTimeEntry() = %+v, want nil when nothing is running", entry)
	}
}

func TestIntegration_TimeEntries_CRUD(t *testing.T) {
	client := integrationClient(t)
	wsID := integrationWorkspaceID(t)
	ctx := integrationCtx(t)

	// Start a running entry.
	started, _, err := client.TimeEntries.StartTimeEntry(ctx, wsID, &toggl.CreateTimeEntryOptions{
		Start:       time.Now().UTC(),
		Description: toggl.String(uniqueName("entry")),
	})
	if err != nil {
		t.Fatalf("StartTimeEntry: %v", err)
	}
	if started.ID == 0 {
		t.Fatal("started entry has ID=0")
	}
	t.Cleanup(func() {
		client.TimeEntries.DeleteTimeEntry(ctx, wsID, started.ID) //nolint:errcheck
	})
	if started.Duration != -1 {
		t.Errorf("Duration = %d, want -1 for running entry", started.Duration)
	}

	// GetRunningTimeEntry should now return the entry we just started.
	running, _, err := client.TimeEntries.GetRunningTimeEntry(ctx)
	if err != nil {
		t.Fatalf("GetRunningTimeEntry: %v", err)
	}
	if running == nil {
		t.Fatal("GetRunningTimeEntry() = nil, want running entry")
	}
	if running.ID != started.ID {
		t.Errorf("running ID = %d, want %d", running.ID, started.ID)
	}

	// Get by ID.
	got, _, err := client.TimeEntries.GetTimeEntry(ctx, started.ID)
	if err != nil {
		t.Fatalf("GetTimeEntry: %v", err)
	}
	if got.ID != started.ID {
		t.Errorf("GetTimeEntry ID = %d, want %d", got.ID, started.ID)
	}

	// Stop.
	stopped, _, err := client.TimeEntries.StopTimeEntry(ctx, wsID, started.ID)
	if err != nil {
		t.Fatalf("StopTimeEntry: %v", err)
	}
	if stopped.Duration < 0 {
		t.Errorf("stopped Duration = %d, want >= 0", stopped.Duration)
	}

	// Update.
	newDesc := uniqueName("entry-updated")
	updated, _, err := client.TimeEntries.UpdateTimeEntry(ctx, wsID, stopped.ID, &toggl.UpdateTimeEntryOptions{
		Description: toggl.String(newDesc),
	})
	if err != nil {
		t.Fatalf("UpdateTimeEntry: %v", err)
	}
	if updated.Description == nil || *updated.Description != newDesc {
		t.Errorf("updated Description = %v, want %q", updated.Description, newDesc)
	}

	// Delete.
	if _, err := client.TimeEntries.DeleteTimeEntry(ctx, wsID, stopped.ID); err != nil {
		t.Fatalf("DeleteTimeEntry: %v", err)
	}
}

func TestIntegration_TimeEntries_ListTimeEntries(t *testing.T) {
	client := integrationClient(t)
	ctx := integrationCtx(t)

	// Use a recent window — Toggl rejects queries older than ~3 months.
	now := time.Now().UTC()
	entries, _, err := client.TimeEntries.ListTimeEntries(ctx, &toggl.ListTimeEntriesOptions{
		StartDate: toggl.String(now.AddDate(0, -1, 0).Format("2006-01-02")),
		EndDate:   toggl.String(now.Format("2006-01-02")),
	})
	if err != nil {
		t.Fatalf("ListTimeEntries: %v", err)
	}
	// An empty list is valid; just verify no error and correct types.
	for _, e := range entries {
		if e.ID == 0 {
			t.Errorf("entry has ID=0")
		}
	}
	t.Logf("found %d time entries in the last 30 days", len(entries))
}
