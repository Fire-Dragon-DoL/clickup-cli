#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CLI="$PROJECT_DIR/clickup"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=== ClickUp CLI Integration Test: Setup Verification ==="
echo ""

if [ ! -f "$CLI" ]; then
    echo -e "${YELLOW}Building CLI...${NC}"
    cd "$PROJECT_DIR"
    go build -o clickup ./cmd/clickup
fi

echo "Checking API key in keyring..."
if ! $CLI folders list >/dev/null 2>&1; then
    echo -e "${RED}FAIL${NC}: API key not found or invalid"
    echo ""
    echo "Please set up your API key using one of these methods:"
    echo ""
    echo "Linux (secret-tool):"
    echo "  secret-tool store --label='ClickUp CLI API Key' service clickup-cli username api_key"
    echo ""
    echo "macOS (security):"
    echo "  security add-generic-password -s \"clickup-cli\" -a \"api_key\" -w \"your_api_key\""
    echo ""
    exit 1
fi
echo -e "${GREEN}PASS${NC}: API key found and valid"

echo ""
echo "Checking CLICKUP_SPACE_ID environment variable..."
if [ -z "$CLICKUP_SPACE_ID" ]; then
    echo -e "${YELLOW}WARN${NC}: CLICKUP_SPACE_ID not set"
    echo "Some tests may fail without a configured space"
    echo "Set it with: export CLICKUP_SPACE_ID=\"your_space_id\""
else
    echo -e "${GREEN}PASS${NC}: CLICKUP_SPACE_ID is set to: $CLICKUP_SPACE_ID"
fi

echo ""
echo "Verifying CLI can list folders..."
OUTPUT=$($CLI folders list 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: CLI successfully connected to ClickUp API"
    echo ""
    echo "Folders in space:"
    echo "$OUTPUT"
else
    echo -e "${RED}FAIL${NC}: Could not list folders"
    echo "$OUTPUT"
    exit 1
fi

echo ""
echo -e "${GREEN}=== Setup verification complete ===${NC}"
