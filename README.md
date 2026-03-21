# go-toggl

A modern Go client library for the Toggl Track API v9.

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

  // Use the client
  // te, resp, err := client.GetTimeEntry(context.Background(), 12345)
}
```

## Services

- **TimeEntries** - Time entry management

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
