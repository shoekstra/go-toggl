# go-toggl

A modern Go client library for the Toggl Track API v9.

[![Go Reference](https://pkg.go.dev/badge/github.com/shoekstra/go-toggl.svg)](https://pkg.go.dev/github.com/shoekstra/go-toggl)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoekstra/go-toggl)](https://goreportcard.com/report/github.com/shoekstra/go-toggl)
[![codecov](https://codecov.io/github/shoekstra/go-toggl/graph/badge.svg?token=K6DEQB183Y)](https://codecov.io/github/shoekstra/go-toggl)

## Installation

```bash
go get github.com/shoekstra/go-toggl
```

## Usage

```go
package main

import (
  "context"
  "log"

  "github.com/shoekstra/go-toggl"
)

func main() {
  client, err := toggl.NewClient("your-api-token")
  if err != nil {
    log.Fatal(err)
  }

  te, resp, err := client.TimeEntries.GetTimeEntry(context.Background(), 12345)
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("time entry %d (status %d)", te.ID, resp.StatusCode)
}
```

## Pagination

List endpoints that support pagination accept `Page` and `PerPage` options.
The `Response` returned by every method includes a `Pagination` field with
metadata parsed from response headers.

### Page-based pagination

```go
perPage := 50
page := 1
for {
    projects, resp, err := client.Projects.ListProjects(ctx, workspaceID, &toggl.ListProjectsOptions{
        Page:    toggl.Int(page),
        PerPage: toggl.Int(perPage),
    })
    if err != nil {
        log.Fatal(err)
    }
    // process projects...

    // TotalPages is 0 when the endpoint does not return X-Pages header;
    // in that case stop when the page returns fewer items than requested.
    if resp.Pagination.TotalPages > 0 && resp.Pagination.CurrentPage >= resp.Pagination.TotalPages {
        break
    }
    if resp.Pagination.TotalPages == 0 && len(projects) < perPage {
        break
    }
    page++
}
```

### Cursor-based pagination (detailed reports)

The detailed reports endpoint uses cursor-based pagination via the
`X-Next-ID` and `X-Next-Row-Number` response headers.

```go
var firstID, firstRowNumber *int
for {
    entries, resp, err := client.Reports.DetailedReport(ctx, workspaceID, &toggl.DetailedReportOptions{
        FirstID:        firstID,
        FirstRowNumber: firstRowNumber,
    })
    if err != nil {
        log.Fatal(err)
    }
    // process entries...

    if resp.Pagination.NextID == 0 {
        break
    }
    firstID = toggl.Int(resp.Pagination.NextID)
    firstRowNumber = toggl.Int(resp.Pagination.NextRowNumber)
}
```

## Services

- **Me** - Authenticated user profile
- **TimeEntries** - Time entry management
- **Projects** - Project management
- **Tags** - Tag management
- **Clients** - Client management
- **Workspaces** - Workspace information
- **Reports** - Report generation (Reports API)

## Development

### Prerequisites

- Go 1.22+
- Task (<https://taskfile.dev>)
- golangci-lint

### Common Tasks

```bash
# Run all tests
task test

# Run integration tests against the real Toggl API
TOGGL_API_TOKEN=<token> TOGGL_WORKSPACE_ID=<id> task test:integration

# Format code
task fmt

# Run linter
task lint

# Run everything: fmt, vet, lint, test
task all

# View coverage
task test:coverage

# Clean build artifacts
task clean
```

Integration tests are gated with `//go:build integration` and run sequentially
(`-p 1`) to stay within Toggl's API rate limits. They require a valid API token
and workspace ID and will create and delete real data.

## License

Apache License 2.0
