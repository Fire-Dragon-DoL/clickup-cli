# ClickUp CLI

A command-line interface for ClickUp.

## Installation

```bash
go install github.com/Fire-Dragon-DoL/clickup-cli/cmd/clickup@latest
```

Or build from source:

```bash
git clone https://github.com/Fire-Dragon-DoL/clickup-cli.git
cd clickup-cli
go build -o clickup ./cmd/clickup
```

## Configuration

### Config File

Create a config file at `~/.config/clickup/config.json`:

```json
{
  "space_id": "your_space_id",
  "output_format": "text",
  "strict_resolve": false
}
```

### Environment Variables

```bash
export CLICKUP_SPACE_ID="your_space_id"
export CLICKUP_OUTPUT_FORMAT="json"
export CLICKUP_STRICT_RESOLVE="true"
```

### CLI Flags

```bash
clickup --space "space_id" --output json --strict
```

### Priority Order

Configuration is loaded in this order (later overrides earlier):

1. Config file
2. Environment variables
3. CLI flags

## API Key Setup

Store your ClickUp API key in the system keyring:

**Linux (using secret-tool):**

```bash
secret-tool store --label='ClickUp CLI API Key' service clickup-cli username api_key
```

**macOS (using security):**

```bash
security add-generic-password -s "clickup-cli" -a "api_key" -w "your_api_key"
```

**Windows (using cmdkey):**

```cmd
cmdkey /generic:clickup-cli /user:api_key /pass:your_api_key
```

Get your API key from: https://app.clickup.com/settings/apps

## Output Formats

- `text` (default): Human-readable output
- `json`: Machine-readable JSON output

## Commands

### Folders

#### List Folders

List all folders in the configured space.

```bash
clickup folders list
```

### Lists

#### List Lists

List all lists in a folder.

```bash
clickup lists list --folder <name|id|url>
clickup lists list -f "My Folder"
```

### Tasks

#### List Tasks

List tasks in a list.

```bash
clickup tasks list --list <name|id|url>
clickup tasks list -l "Backlog"
```

Include subtasks with `--recursive`:

```bash
clickup tasks list --list "Backlog" --recursive
clickup tasks list -l "Backlog" -r
```

Recursive output shows hierarchical indentation:

```
task1 | Parent Task | john | in progress | high
  task1.1 | Subtask 1 | jane | open | medium
    task1.1.1 | Sub-subtask | | completed |
```

#### Show Task

Display detailed information about a task.

```bash
clickup tasks show <task-id|name|url>
```

Shows: ID, Title, Description, Assignee, Status, Priority, Due Date, and Comments.

#### Create Task

Create a new task.

```bash
clickup tasks create --title <title> --list <name|id|url> [options]
```

**Required:**
- `--title, -t`: Task title
- `--list, -l`: List name, ID, or URL

**Optional:**
- `--description, -d`: Task description (markdown)
- `--priority, -p`: Priority (1-5, 0=none)
- `--status`: Task status
- `--due`: Due date
- `--assignee`: Assignee name, ID, or username
- `--parent`: Parent task ID or name

Example:

```bash
clickup tasks create \
  --title "Implement feature X" \
  --list "Backlog" \
  --description "This feature should do X, Y, and Z" \
  --priority 2 \
  --assignee "john"
```

#### Update Task

Update an existing task.

```bash
clickup tasks update <task-id|name|url> [options]
```

**Options:**
- `--title, -t`: Update title
- `--status, -s`: Update status
- `--priority, -p`: Update priority
- `--description, -d`: Update description
- `--due`: Update due date
- `--assignee, -a`: Update assignee
- `--parent`: Set parent task

Only specified fields are updated. Example:

```bash
clickup tasks update "Fix login bug" --status "done" --assignee "jane"
```

#### Delete Task

Delete a task permanently.

```bash
clickup tasks delete <task-id|name|url>
```

#### Archive Task

Archive a task.

```bash
clickup tasks archive <task-id|name|url>
```

## Resource Identifiers

Tasks, lists, folders, and users can be referenced by:
- **ID**: ClickUp internal ID (e.g., `abc123`)
- **Name**: Human-readable name (e.g., `"Fix login bug"`)
- **URL**: Browser URL (e.g., `https://app.clickup.com/t/abc123`)

### Ambiguous Name Resolution

When a name matches multiple resources:
- `strict_resolve: false` (default): Uses first match
- `strict_resolve: true`: Fails with error listing matches

## Integration Tests

Integration tests verify the CLI against a real ClickUp workspace.

### Prerequisites

1. Set up your API key in the system keyring (see API Key Setup above)
2. Set the `CLICKUP_SPACE_ID` environment variable:
   ```bash
   export CLICKUP_SPACE_ID="your_space_id"
   ```

### Running Tests

Run all integration tests:

```bash
./scripts/run-all.sh
```

Run individual tests:

```bash
./scripts/test-setup.sh      # Verify API key and configuration
./scripts/test-folders.sh    # Test folder operations
./scripts/test-lists.sh      # Test list operations
./scripts/test-tasks-list.sh # Test task listing
./scripts/test-task-crud.sh  # Test create/show/update/delete
```

**Note:** The CRUD test creates a temporary task and cleans it up after completion.

## License

MIT
