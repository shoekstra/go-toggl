# Examples

This directory contains example usage of the go-toggl client. Each file
demonstrates the key operations for a single service.

The examples are excluded from the normal build (`//go:build ignore`) but can
be run directly:

```bash
export TOGGL_API_TOKEN=your_api_token
export TOGGL_WORKSPACE_ID=your_workspace_id  # not needed for me.go or workspaces.go

go run examples/projects.go
```

Run `examples/workspaces.go` first if you need to find your workspace ID:

```bash
TOGGL_API_TOKEN=your_api_token go run examples/workspaces.go
```

## Files

| File | Service | Key operations shown |
|---|---|---|
| `me.go` | `Me` | `GetMe` |
| `workspaces.go` | `Workspaces` | `ListWorkspaces`, `GetWorkspace`, `UpdateWorkspace` |
| `projects.go` | `Projects` | `ListProjects`, `GetProject`, `CreateProject`, `UpdateProject`, `DeleteProject` |
| `tags.go` | `Tags` | `ListTags`, `CreateTag`, `UpdateTag`, `DeleteTag` |
| `clients.go` | `Clients` | `ListClients`, `GetClient`, `CreateClient`, `UpdateClient`, `ArchiveClient`, `RestoreClient`, `DeleteClient` |
| `time_entries.go` | `TimeEntries` | `ListTimeEntries`, `GetRunningTimeEntry`, `StartTimeEntry`, `StopTimeEntry`, `UpdateTimeEntry`, `DeleteTimeEntry` |
| `reports.go` | `Reports` | `SummaryReport`, `DetailedReport`, `DetailedReportTotals`, `WeeklyReport`, `ExportSummaryPDF`, `ExportDetailedCSV` |
