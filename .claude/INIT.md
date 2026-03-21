# Claude Code Initialization for go-toggl

## Before Starting Work

When you run Claude Code for the first time, give this prompt:

```
"I'm working on go-toggl, a service-oriented Go client for Toggl Track API v9.

Before I give you work to do, please read and understand these context files:
- .claude/ARCHITECTURE.md (service design, what to build)
- .claude/CONVENTIONS.md (code style, commit messages, branch names)
- .claude/IMPLEMENTATION_PATTERNS.md (code templates to follow)
- .claude/V8_REFERENCE.md (patterns from the v8 SDK)
- .claude/swagger-reports.json (Reports OpenAPI spec)
- .claude/swagger-toggl-api.json (Toggl API OpenAPI spec)
- .claude/swagger-webhooks.json (Webhooks OpenAPI spec)

Then tell me:
1. What you understand about the project
2. The services to implement
3. The coding conventions you'll follow
4. Confirm you understand the template patterns
5. Ask any clarifying questions

Then I'll give you implementation tasks."
```

Claude will:

1. Read all context files
2. Confirm understanding
3. Be ready to implement services correctly

## Subsequent Sessions

For later sessions, shorter version:

```
"Refresh your understanding of go-toggl by re-reading:
- .claude/CONVENTIONS.md (for commit format)
- .claude/IMPLEMENTATION_PATTERNS.md (for code templates)

Then tell me you're ready and what service to implement."
```

Or if continuing work:

```
"Continue from where we left off. Review git log to see what's been implemented,
then I'll tell you what to work on next."
```

## What Claude Will Do

With proper context, Claude Code will:

- ✅ Write service methods following patterns
- ✅ Write comprehensive unit tests
- ✅ Use conventional commit messages
- ✅ Format code properly (task fmt)
- ✅ Check linters (task lint)
- ✅ Run tests and fix failures
- ✅ Commit to git automatically

## Typical Workflow

```bash
# 1. Start Claude Code
cd ~/projects/go-toggl
claude code

# 2. Give initialization prompt (copy from above)
# Claude reads context and confirms understanding

# 3. Give implementation task
# "Implement TimeEntriesService with full CRUD, comprehensive tests, and commit"

# 4. Claude:
#    - Writes time_entries.go with all methods
#    - Writes time_entries_test.go with tests
#    - Runs task fmt
#    - Runs task lint
#    - Runs task test
#    - Commits to git with conventional message

# 5. Review output and give next task
# "Implement ProjectsService with full CRUD..."

# 6. Repeat for each service
```

## Tips for Success

- **Be specific**: Don't say "implement services", say "implement TimeEntriesService with ListTimeEntries, GetTimeEntry, CreateTimeEntry, UpdateTimeEntry, DeleteTimeEntry"
- **Reference files**: "Following IMPLEMENTATION_PATTERNS.md, implement..."
- **Ask for validation**: "Run 'task all' and fix any failures"
- **Request commits**: "Commit with 'feat(ServiceName): add CRUD operations'"
