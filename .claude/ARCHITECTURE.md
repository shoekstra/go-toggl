# go-toggl Architecture Guide

## Overview

go-toggl is a modern, service-oriented Go client for the Toggl Track API v9.
It follows the architectural patterns from go-gitlab.

## Design Philosophy

### Service-Based Architecture

Each API resource area (TimeEntries, Projects, Tags, etc.) is handled by its own service struct:

```go
type TimeEntriesService struct {
    client *Client
}

type ProjectsService struct {
    client *Client
}
```

All services are aggregated in a main Client struct:

```go
type Client struct {
    baseURL     string
    httpClient  *http.Client
    token       string

    TimeEntries *TimeEntriesService
    Projects    *ProjectsService
    Tags        *TagsService
    Clients     *ClientsService
    Workspaces  *WorkspacesService
    Reports     *ReportsService
}
```

### Key Design Decisions

1. **Context Required**: Every method accepts `context.Context` as first parameter

   - Enables timeouts, cancellation, tracing
   - Example: `GetTimeEntry(ctx context.Context, id int) (*TimeEntry, *Response, error)`

2. **Functional Options Pattern**: Client configuration via options

   ```go
   client, _ := toggl.NewClient("token",
       toggl.WithBaseURL("https://custom.url"),
       toggl.WithTimeout(60 * time.Second),
   )
   ```

3. **Centralized HTTP Layer**: Single place for all HTTP operations

   - `get(ctx, path, v)` - GET request
   - `post(ctx, path, body, v)` - POST request
   - `put(ctx, path, body, v)` - PUT request
   - `delete(ctx, path)` - DELETE request

4. **Pointer Helpers**: For optional parameters

   ```go
   toggl.String("value")
   toggl.Int(42)
   toggl.Bool(true)
   toggl.Time(time.Now())
   ```

5. **Error Handling**: Custom ErrorResponse type
   - StatusCode, Message, HTTP Response
   - All methods return (result, \*Response, error)

## Services to Implement

### 1. TimeEntriesService

Manage time entries (Toggl's main resource)

- ListTimeEntries(ctx, opts) - List with filters (date range, billable, etc.)
- GetTimeEntry(ctx, id) - Get single entry
- CreateTimeEntry(ctx, opts) - Create new entry
- UpdateTimeEntry(ctx, id, opts) - Update entry
- DeleteTimeEntry(ctx, id) - Delete entry
- StartTimeEntry(ctx, opts) - Start tracking
- StopTimeEntry(ctx, id) - Stop tracking

### 2. ProjectsService

Manage projects

- ListProjects(ctx, opts) - List active/archived
- GetProject(ctx, id) - Get single project
- CreateProject(ctx, opts) - Create project
- UpdateProject(ctx, id, opts) - Update project
- DeleteProject(ctx, id) - Delete project

### 3. TagsService

Manage tags

- ListTags(ctx, workspace) - List workspace tags
- GetTag(ctx, workspace, id) - Get single tag
- CreateTag(ctx, workspace, opts) - Create tag
- UpdateTag(ctx, workspace, id, opts) - Update tag
- DeleteTag(ctx, workspace, id) - Delete tag

### 4. ClientsService

Manage clients (organizations)

- ListClients(ctx, workspace, opts) - List clients
- GetClient(ctx, workspace, id) - Get single client
- CreateClient(ctx, workspace, opts) - Create client
- UpdateClient(ctx, workspace, id, opts) - Update client
- DeleteClient(ctx, workspace, id) - Delete client

### 5. WorkspacesService

Workspace information (read-only)

- ListWorkspaces(ctx) - List user's workspaces
- GetWorkspace(ctx, id) - Get workspace details

### 6. ReportsService

Generate reports (Reports API v2)

- GenerateReport(ctx, opts) - Generate report
- GetReportPDF(ctx, reportID) - Get PDF report
- GetReportCSV(ctx, reportID) - Get CSV report

## File Structure

```
go-toggl/
├── .claude/                          # Context for Claude Code
│   ├── ARCHITECTURE.md              # This file
│   ├── CONVENTIONS.md               # Code style and commit rules
│   ├── IMPLEMENTATION_PATTERNS.md   # Code templates
│   ├── V8_REFERENCE.md              # v8 SDK patterns
│   ├── api-spec.yaml                # Toggl API v9 spec
│   └── .clauderc                    # Auto-loaded config (optional)
│
├── .github/workflows/               # CI/CD pipelines
│   ├── test.yml
│   └── lint.yml
│
├── examples/                        # Usage examples
│   ├── time_entries_example.go
│   ├── projects_example.go
│   └── README.md
│
├── testdata/                        # Test fixtures
│   ├── time_entry.json
│   ├── project.json
│   └── README.md
│
├── time_entries.go                  # Service implementation
├── time_entries_test.go
├── projects.go
├── projects_test.go
├── tags.go
├── tags_test.go
├── clients.go
├── clients_test.go
├── workspaces.go
├── workspaces_test.go
├── reports.go
├── reports_test.go
│
├── toggl.go                         # Main client
├── client_options.go                # Configuration options
├── request.go                       # HTTP layer
├── errors.go                        # Error types
├── types.go                         # Shared types
│
├── go.mod
├── go.sum
├── Taskfile.yaml                    # Build tasks
├── .golangci.yml                    # Linter config
├── README.md
├── LICENSE
└── CONTRIBUTING.md
```

## Testing Strategy

### Unit Tests - Table-Driven with Mocks

Every service method has table-driven unit tests using mocked HTTP:

```go
func TestTimeEntriesService_ListTimeEntries(t *testing.T) {
    tests := []struct {
        name       string
        opts       *ListTimeEntriesOptions
        statusCode int
        response   string
        wantCount  int
        wantErr    bool
    }{
        {
            name:       "success with filter",
            opts:       &ListTimeEntriesOptions{Limit: Int(50)},
            statusCode: 200,
            response:   `[{"id":1,"name":"test"}]`,
            wantCount:  1,
            wantErr:    false,
        },
        {
            name:       "unauthorized",
            opts:       nil,
            statusCode: 401,
            response:   `{"error": "unauthorized"}`,
            wantCount:  0,
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation with httptest mock
        })
    }
}
```

### Mocking HTTP Responses

Use `net/http/httptest` to mock API responses:

```go
func testClient(handler http.Handler) *Client {
    server := httptest.NewServer(handler)
    client, _ := NewClient("test-token", WithBaseURL(server.URL))
    return client
}

handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"id":1}`))
})

client := testClient(handler)
entry, _, err := client.TimeEntries.GetTimeEntry(ctx, 1)
```

### Test Coverage Goals

- Minimum 80% code coverage
- All service methods tested
- Happy path + error cases (400, 401, 404, 500)
- Options parsing and parameter passing

## API Endpoint Pattern

All Toggl API v9 endpoints follow this pattern:

**Base URL**: `https://api.track.toggl.com/api/v9`

**Authentication**: Bearer token in Authorization header

**Time Entries Endpoints**:

- `GET /me/time_entries` - List entries for current user
- `GET /me/time_entries/{id}` - Get single entry
- `POST /workspaces/{workspace_id}/time_entries` - Create entry
- `PUT /workspaces/{workspace_id}/time_entries/{id}` - Update entry
- `DELETE /workspaces/{workspace_id}/time_entries/{id}` - Delete entry

**Similar patterns for Projects, Tags, Clients, Workspaces, Reports**

## Implementation Workflow

1. **Define Types**: Add to types.go

   - Main resource type (e.g., `TimeEntry`)
   - Options types (e.g., `ListTimeEntriesOptions`, `CreateTimeEntryOptions`)

2. **Implement Service**: Create service_name.go

   - Service struct with client reference
   - All CRUD methods with full godoc comments
   - Follow the implementation patterns from IMPLEMENTATION_PATTERNS.md

3. **Write Tests**: Create service_name_test.go

   - Table-driven tests for each method
   - Mocked HTTP responses
   - Error cases

4. **Create Examples**: examples/service_name_example.go

   - Show how to use the service
   - Make runnable with comments

5. **Update Documentation**: README.md

   - Add service to list
   - Update example if needed

6. **Commit**: Use conventional commits
   - `feat(TimeEntriesService): add full CRUD operations`
   - `test(TimeEntriesService): add comprehensive unit tests`
