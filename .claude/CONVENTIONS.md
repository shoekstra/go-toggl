# go-toggl Development Conventions

## Git Workflow

### Branch Naming

Use conventional branch names:

```
feat/feature-name        # New feature
fix/bug-description      # Bug fix
docs/what-changed        # Documentation
test/what-tested         # Test additions
refactor/what-changed    # Code refactoring
chore/task-description   # Build, deps, tooling
```

Examples:

```
feat/time-entries-service
fix/bearer-token-header
docs/readme-examples
test/time-entries-mocks
refactor/http-client
chore/update-dependencies
```

### Commit Messages

Use conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `test` - Tests
- `refactor` - Code refactoring
- `perf` - Performance improvement
- `chore` - Build, dependencies, tooling
- `ci` - CI/CD configuration

**Scope** (optional but recommended):

- Service name: `TimeEntriesService`, `ProjectsService`
- Component: `request`, `types`, `client`

**Subject**:

- Use imperative mood: "add", "fix", "update" (not "added", "fixed", "updated")
- Don't capitalize first letter
- No period at end
- Max 50 characters

**Examples**:

```
feat(TimeEntriesService): add list with filtering

Add ListTimeEntries method with optional date range, billable, and project filters.
Includes comprehensive table-driven unit tests with mocked HTTP responses.

Fixes #42
```

```
test(TimeEntriesService): add error case tests

Add tests for 401, 404, and 500 responses to ensure proper error handling.
```

```
fix(request): handle empty response body

Previously crashed on 204 No Content responses. Now safely handles empty bodies.
```

```
docs: add TimeEntriesService examples
```

## Code Style

### Go Idioms

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting (automatic via `task fmt`)
- Use `goimports` to organize imports

### Naming Conventions

**Types**:

- Use PascalCase: `TimeEntry`, `ProjectsService`, `ListTimeEntriesOptions`
- Service types end with `Service`: `TimeEntriesService`, `ProjectsService`
- Options types end with `Options`: `ListTimeEntriesOptions`, `CreateTimeEntryOptions`

**Functions/Methods**:

- Use PascalCase for exported: `GetTimeEntry`, `ListTimeEntries`
- Use camelCase for unexported: `newRequest`, `do`

**Variables**:

- Use camelCase: `timeEntry`, `projectID`, `workspace`
- Avoid single letters except loop counters and `ctx` for context

**Constants**:

- Use MixedCase: `defaultBaseURL`, `defaultTimeout`

### Godoc Comments

All exported items must have godoc comments:

```go
// TimeEntriesService handles operations related to time entries.
type TimeEntriesService struct {
    client *Client
}

// ListTimeEntries lists all time entries in a workspace with optional filters.
//
// The opts parameter is optional. If nil, all entries are returned with default
// pagination. Use opts to filter by date range, billable status, project, etc.
//
// API: GET /me/time_entries
//
// See: https://engineering.toggl.com/docs/track/api/time_entries#list-all-time-entries
func (s *TimeEntriesService) ListTimeEntries(ctx context.Context, opts *ListTimeEntriesOptions) ([]*TimeEntry, *Response, error) {
    // ...
}
```

**Godoc Style**:

- Start with the function name: "ListTimeEntries lists..."
- Add explanation of parameters if complex
- Include API endpoint: `// API: GET /path`
- Add link to API documentation if available
- Mention any important behavior (e.g., optional parameters, pagination)

### Error Handling

- Always wrap errors with context: `fmt.Errorf("failed to decode response: %w", err)`
- Return (result, response, error) from service methods
- Don't ignore error responses from HTTP calls
- Use custom ErrorResponse type for API errors

### Testing

**Table-Driven Tests**:

```go
func TestServiceService_Method(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        statusCode  int
        response    string
        want        string
        wantErr     bool
    }{
        {
            name:       "happy path description",
            input:      "test input",
            statusCode: 200,
            response:   `{"id": 1, "name": "test"}`,
            want:       "expected result",
            wantErr:    false,
        },
        {
            name:       "error case description",
            input:      "bad input",
            statusCode: 400,
            response:   `{"error": "invalid"}`,
            want:       "",
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Implementation
        })
    }
}
```

**Test Coverage**:

- Happy path (200, data returned)
- Client errors (400, 401, 404)
- Server errors (500, 502, 503)
- Edge cases (empty results, nil options)

**Naming**:

- Test files: `service_test.go` (lowercase, suffix `_test.go`)
- Test functions: `TestServiceService_Method` (follows `Test<Type>_<Method>` pattern)

## Code Review Checklist

Before committing, verify:

- [ ] `task fmt` - Code formatted
- [ ] `task lint` - No linter errors
- [ ] `task test` - All tests pass
- [ ] `task coverage` - At least 80% coverage
- [ ] Godoc comments on all exported items
- [ ] Error handling is explicit
- [ ] No TODOs left (except intentional placeholders)
- [ ] Conventional commit message used
- [ ] Branch name follows convention

## Pre-Commit Workflow

```bash
# Make your changes
# ...

# Format and lint
task fmt
task lint

# Run tests
task test

# Check coverage
task coverage

# Review changes
git diff

# Stage and commit
git add .
git commit -m "feat(ServiceName): add new capability"

# Push
git push origin feat/feature-name
```

## Documentation Requirements

Every service needs:

1. **Godoc comments** on all exported items
2. **README.md** - Installation, usage, services list
3. **CONTRIBUTING.md** - How to contribute
4. **examples/** - Working example for each service
5. **testdata/** - JSON fixtures for tests (if complex)
