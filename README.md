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

## License

MIT
