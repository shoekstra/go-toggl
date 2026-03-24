# Claude Code Initialization for go-toggl

## Project Status

go-toggl is a complete, released Go client library for the Toggl Track API v9.
**Current version: v0.1.0** (tagged and released via release-please).

All six services are implemented and merged to main:

- TimeEntriesService — full CRUD + Start/Stop
- ProjectsService — full CRUD
- TagsService — List, Create, Update, Delete
- ClientsService — full CRUD + Archive/Restore
- WorkspacesService — List, Get, Update
- ReportsService — Summary, Detailed, Weekly + PDF/CSV exports

## Starting a New Session

Tell Claude:

```
"I'm continuing work on go-toggl, a released Go client library for the Toggl
Track API v9. Check git log and your memory for recent context, then tell me
you're ready."
```

Claude should then check:

- `git log --oneline -20` to see recent activity
- Memory files for decisions and feedback
- Any open PRs or branches

## Conventions (summarised)

- Branch: `feat/X`, `fix/X`, `chore/X`, `ci/X`, `refactor/X`
- Always branch from `main`
- Conventional commits with `-s` (DCO sign-off)
- One logical change per PR
- After each commit, user pushes and opens PR; Claude waits for review feedback
- CodeRabbit reviews PRs — process findings with "Verify each finding against
  the current code and only fix it if needed"
- Amend commits rather than adding fix commits (user force-pushes)
- Use `git rebase -i` to squash, never `git reset --soft`

## CI / Tooling

- Tests: `go test ./...` — must pass before committing
- Lint: golangci-lint via `task lint` (improvements pending — see next tasks)
- Coverage: Codecov — patch coverage is checked on PRs
- Release: release-please (GitHub App token via `shoekstra-ci` app)
- Renovate: Mend GitHub App — manages dependency and Actions updates,
  config at `.github/renovate.json`

## Architecture

- Service-oriented (go-gitlab pattern): `client.TimeEntries.X`, `client.Projects.X`, etc.
- Base URL: `https://api.track.toggl.com`; main API: `/api/v9/`; Reports API: `/reports/api/v3/`
- All Reports endpoints use POST; export methods return `[]byte`
- Auth: Basic auth with token + "api_token"
- Pagination: `Pagination` struct on every `Response`; headers parsed automatically
- Options structs use pointer fields + `json:",omitempty"` for body serialisation

## Pending Work

1. **golangci-lint improvements** — branch `chore/golangci-lint`

   - Add `errcheck`, `staticcheck`, `gosimple` linters at minimum
   - Add `godot` for godoc comment consistency
   - Set explicit `line-length` under `linters-settings.lll`
   - Fix `develop` branch trigger leftover in `lint.yaml`
   - Fix any new findings surfaced by the additional linters

2. **GetTag** — not yet implemented (documented in ARCHITECTURE.md but missing
   from tags.go). Should be done via `feat/GetTag` branch.

3. **Examples** — `examples/` directory is empty; deferred until after v0.1.0.
   One file per service (e.g. `examples/time_entries.go`), not one directory per service.

## Key Files

- `toggl.go` — Client struct, NewClient, functional options pattern
- `client_options.go` — WithBaseURL, WithHTTPClient, WithTimeout
- `request.go` — HTTP layer, Response struct, Pagination parsing
- `types.go` — all shared types including Tag.UnmarshalJSON for backwards compat
- `errors.go` — ErrorResponse
- `.github/renovate.json` — Renovate config (pinDigests, grouping rules)
- `release-please-config.json` — release-please config (go release type)
