# go-toggl

A modern Go client library for the Toggl Track API v9.

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

## Services

- **TimeEntries** - Time entry management
- **Projects** - Project management
- **Tags** - Tag management
- **Clients** - Client management
- **Workspaces** - Workspace information
- **Reports** - Report generation (Reports API)

## Development

### Prerequisites

- Go 1.21+
- Task (<https://taskfile.dev>)
- golangci-lint

### Common Tasks

```bash
# Run all tests
task test

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

## License

Apache License 2.0
