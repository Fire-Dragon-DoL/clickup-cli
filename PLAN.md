# ClickUp CLI Development Plan

## Requirements

Develop a CLI for ClickUp with support for both human-readable text output and JSON output.

### Features

**List Operations:**
- List tasks in a list, top level only → **List View**
- List tasks in a list, recursive with subtasks → **List View**
- List lists in a folder
- List folders in the current space

**Show Operations:**
- Show task by identifier (id/name/url) → **Details View**

**Create/Update/Delete Operations:**
- Create task → **Details View** (shows created task)
  - Required: title, list (by name or ID)
  - Optional: parent task, assignee, description (markdown), due date, status, priority
- Update task → **Details View** (shows updated task)
  - Same optional inputs as create
- Delete task → confirmation message only
- Archive task → confirmation message only

**Space Configuration:**
- Set via config file, environment variable, or CLI argument (global)

**Configuration Priority (lowest to highest):**
1. Config file
2. Environment variable
3. Command line argument (overwrites everything)

**Output Format:** Text (human-readable) or JSON (machine-readable)
- Configurable via config file, env var, or CLI arg (same priority chain)

**Resource Identifiers:**
Tasks, lists, folders, and users can be referenced by:
- **ID** - ClickUp internal ID (e.g., `abc123`)
- **Name** - Human-readable name (e.g., `"Fix login bug"`)
- **URL** - Browser URL pasted directly (e.g., `https://app.clickup.com/t/abc123`)

Resolver should be simple but structured to easily add support for new string formats.

**Ambiguous Name Resolution:**
- Config option: `strict_resolve` (default: `false`)
- When `false`: search returns multiple results → use first match
- When `true`: search returns multiple results → fail with error listing matches

> **Open Question (defer for now):**
> - How to differentiate IDs from names (prefix like `#abc123`? detect format?)

**API Key Retrieval:** System keyring (see Dependencies)

**Tech Stack:** Go

---

## Future: TUI (out of scope)

A terminal UI (`lazyclickup`) inspired by lazygit is planned for the future. Not part of this plan, but code should be written with reuse in mind:

- **CLI executable:** `clickup` (this plan)
- **TUI executable:** `lazyclickup` (future)
- Both in `cmd/` directory
- `internal/` packages should be UI-agnostic (no direct stdout writes from business logic)
- Minimize code in `main.go` (any "main" function) so that code is easier to test and reuse

---

## Dependencies

**Production:**
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management (file, env, flags with precedence)
- `github.com/zalando/go-keyring` - Cross-platform keyring access

**Development:**
- Go stdlib `testing` package only

**Config file format:** JSON (default for examples)

---

## Task Output Forms

### List View (Short Form)
Used when displaying multiple tasks. Fields:
- ID
- Title
- Assignee
- Status
- Priority

### Details View (Long Form)
Used when displaying a single task. Fields:
- All fields from List View
- Due Date
- Description (markdown)
- Comments
- Dependencies (blocked by / blocking)

---

## Workflow Per Phase

### Sequential Phases (0, 1, 2, 10, 11)
1. Write tests first
2. Run all tests, verify new tests fail
3. Implement code
4. Run tests, verify success
5. Update README.md with usage instructions (if applicable)
6. Mark phase complete in Progress section
7. Commit: `feat: phase N - <description>`
8. Push to remote

### Parallel Phases (3-9)
Phases 3-9 can be executed in parallel by subagents. Rules:
1. Write tests first
2. Run all tests, verify new tests fail
3. Implement code
4. Run tests, verify success
5. Write documentation to `docs/phase_N.md` (NOT README.md)
6. **Do NOT** mark phase complete in PLAN.md
7. **Do NOT** commit - leave changes uncommitted
8. `docs/phase_N.md` serves as proof of completion

**Important: User must confirm before starting each phase.**

---

## Progress

### Sequential
- [x] Phase 0: Project Foundation
- [ ] Phase 1: API Client
- [ ] Phase 2: Resolver

### Parallel (run after Phase 2 completes)
- [ ] Phase 3: List Folders & Lists
- [ ] Phase 4: List Tasks (Top-Level) - List View
- [ ] Phase 5: List Tasks (Recursive) - List View
- [ ] Phase 6: Task Details - Details View
- [ ] Phase 7: Create Task - Details View
- [ ] Phase 8: Update Task - Details View
- [ ] Phase 9: Delete & Archive Task

### Sequential (run after Phases 3-9 complete)
- [ ] Phase 10: Documentation Consolidation
- [ ] Phase 11: Integration Test Scripts

---

## Phase 0: Project Foundation

### Tests
- Test config loading priority (file < env < CLI arg)
- Test API key retrieval from keyring (mock `go-keyring` interface)
- Test output formatter (text vs JSON)

### Implementation
1. Initialize Go module (`go mod init`)
2. Set up project structure:
   ```
   cmd/clickup/main.go
   internal/config/config.go
   internal/keyring/keyring.go
   internal/api/client.go
   internal/api/mock.go
   internal/output/formatter.go
   ```
3. Add CLI framework (cobra)
4. Implement config loading with priority chain
5. Implement keyring retrieval via `github.com/zalando/go-keyring`
6. Implement output formatter (text/JSON)

### README
- Installation instructions
- Config file location and format
- Environment variables
- API key setup with system keyring

---

## Phase 1: API Client

### Tests
- Test client initialization with API key and base URL
- Test HTTP request/response handling with generic types
- Test authentication header injection
- Test error response parsing

### Implementation
1. Define `Client` struct holding base URL and API key
2. Generic request method: `Do[Req, Res](method, path, body) (Res, error)`
3. Authentication: inject `Authorization` header on all requests
4. Error handling: parse ClickUp API error responses
5. Base URL configurable (default: `https://api.clickup.com/api/v2`)

### README
- N/A (internal infrastructure)

---

## Phase 2: Resolver

### Tests
- Test identifier type detection (ID, Name, URL)
- Test URL parsing for tasks, lists, folders, users
- Test `strict_resolve` behavior (first match vs error on ambiguous)
- Test name resolution with mock API responses

### Implementation
1. Implement `internal/resolver` package:
   - Detect identifier type from string format
   - Parse ClickUp URLs to extract resource IDs
   - Resolve names to IDs via API search
   - Respect `strict_resolve` config for ambiguous name handling
2. Resource types: Task, List, Folder, User

### README
- N/A (internal infrastructure)

---

## Phase 3: List Folders & Lists

### Tests
- Test `folders list` command parsing
- Test `lists list` command parsing
- Test API client `GetFolders(spaceID)` with mock
- Test API client `GetLists(folderID)` with mock
- Test text/JSON output formatting for folders and lists

### Implementation
1. ClickUp API client methods: `GetFolders`, `GetLists`
2. Commands: `clickup folders list`, `clickup lists list --folder <name|id|url>`
3. Resolve folder by name/id/url (requires lookup)

### docs/phase_3.md
- `folders list` usage
- `lists list` usage

---

## Phase 4: List Tasks (Top-Level)

### Tests
- Test `tasks list` command parsing
- Test API client `GetTasks(listID)` with mock
- Test List View output formatting

### Implementation
1. API client method: `GetTasks(listID, recursive=false)`
2. Command: `clickup tasks list --list <name|id|url>` → **List View**
3. Resolve list by name/id/url
4. List View formatter for tasks

### docs/phase_4.md
- `tasks list` usage
- Output columns explanation

---

## Phase 5: List Tasks (Recursive with Subtasks)

### Tests
- Test `tasks list --recursive` flag
- Test API client fetching subtasks
- Test hierarchical List View output

### Implementation
1. Extend `GetTasks` to fetch subtasks recursively
2. Add `--recursive` flag to `tasks list` → **List View** (with indentation)
3. Hierarchical display with indentation for subtasks

### docs/phase_5.md
- `--recursive` flag usage

---

## Phase 6: Task Details

### Tests
- Test `tasks show` command parsing
- Test API client `GetTask(taskID)` with mock
- Test Details View output formatting
- Test API client `GetTaskComments(taskID)` with mock

### Implementation
1. API client methods: `GetTask`, `GetTaskComments`
2. Command: `clickup tasks show <task-id|name|url>` → **Details View**
3. Details View formatter

### docs/phase_6.md
- `tasks show` usage

---

## Phase 7: Create Task

### Tests
- Test `tasks create` command parsing (required + optional args)
- Test API client `CreateTask` with mock
- Test status and priority name resolution

### Implementation
1. API client method: `CreateTask`
2. Command: `clickup tasks create --title <title> --list <name|id|url>` → **Details View** (shows created task)
3. Optional flags: `--parent`, `--assignee`, `--description`, `--due`, `--status`, `--priority`

### docs/phase_7.md
- `tasks create` usage with all flags

---

## Phase 8: Update Task

### Tests
- Test `tasks update` command parsing
- Test API client `UpdateTask` with mock
- Test partial updates (only changed fields)

### Implementation
1. API client method: `UpdateTask`
2. Command: `clickup tasks update <task-id|name|url>` → **Details View** (shows updated task)
3. Same optional flags as create

### docs/phase_8.md
- `tasks update` usage

---

## Phase 9: Delete & Archive Task

### Tests
- Test `tasks delete` command parsing
- Test `tasks archive` command parsing
- Test API client `DeleteTask`, `ArchiveTask` with mock

### Implementation
1. API client methods: `DeleteTask`, `ArchiveTask`
2. Commands: `clickup tasks delete <task-id|name|url>`, `clickup tasks archive <task-id|name|url>`

### docs/phase_9.md
- `tasks delete` and `tasks archive` usage

---

## Phase 10: Documentation Consolidation

### Implementation
1. Read all `docs/phase_N.md` files (N = 3-9)
2. Consolidate into README.md under appropriate sections
3. Delete `docs/phase_N.md` files after consolidation
4. Review all parallel phase code for consistency
5. Run full test suite: `go test ./...`
6. Commit all parallel phase work: `feat: phases 3-9 - CLI commands`

### README
- Complete usage documentation for all commands
- Examples for common workflows

---

## Phase 11: Integration Test Scripts

### Implementation
Create `scripts/` directory with executable test scripts:

1. `scripts/test-setup.sh` - Verify API key in keyring, print space info
2. `scripts/test-folders.sh` - List folders, verify output
3. `scripts/test-lists.sh` - List lists in first folder
4. `scripts/test-tasks-list.sh` - List tasks (top-level and recursive)
5. `scripts/test-task-crud.sh` - Create, show, update, archive, delete a task
6. `scripts/run-all.sh` - Execute all tests in sequence

Each script:
- Uses real API
- Prints clear pass/fail status
- Cleans up created test data

### README
- How to run integration tests
- Required environment setup

---

## Project Structure (Final)

```
clickup-cli/
├── cmd/
│   ├── clickup/main.go        # CLI (this plan)
│   └── lazyclickup/main.go    # TUI (future)
├── internal/
│   ├── api/
│   │   ├── client.go
│   │   ├── client_test.go
│   │   ├── mock.go
│   │   ├── folders.go
│   │   ├── lists.go
│   │   ├── tasks.go
│   │   └── users.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── resolver/
│   │   ├── resolver.go
│   │   └── resolver_test.go
│   ├── keyring/
│   │   ├── keyring.go
│   │   └── keyring_test.go
│   ├── output/
│   │   ├── formatter.go
│   │   └── formatter_test.go
│   └── cmd/
│       ├── root.go
│       ├── folders.go
│       ├── lists.go
│       └── tasks.go
├── docs/                          # Temporary, deleted after Phase 10
│   ├── phase_3.md
│   ├── phase_4.md
│   ├── phase_5.md
│   ├── phase_6.md
│   ├── phase_7.md
│   ├── phase_8.md
│   └── phase_9.md
├── scripts/
│   ├── test-setup.sh
│   ├── test-folders.sh
│   ├── test-lists.sh
│   ├── test-tasks-list.sh
│   ├── test-task-crud.sh
│   └── run-all.sh
├── README.md
├── LICENSE
├── go.mod
└── go.sum
```

---

## Testing Strategy

### Mock Design
- Single `MockClient` implementing `APIClient` interface
- Configurable responses per method
- Tracks call history for assertions

### Test Execution
```bash
go test ./... -v
```

---

## After Each Phase

1. `go test ./...` - All tests pass
2. `go build ./cmd/clickup` - Binary builds
3. `./clickup --help` - CLI responds
4. Update `PLAN.md` - Mark phase as complete with checkbox
5. Commit with message: `feat: phase N - <description>`
6. Push to remote

---

## Final Verification

Run `scripts/run-all.sh` against real ClickUp workspace
