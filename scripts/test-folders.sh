#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CLI="$PROJECT_DIR/clickup"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "=== ClickUp CLI Integration Test: Folders ==="
echo ""

if [ ! -f "$CLI" ]; then
    echo "Building CLI..."
    cd "$PROJECT_DIR"
    go build -o clickup ./cmd/clickup
fi

echo "Test: List folders (text output)..."
OUTPUT=$($CLI folders list 2>&1)
if [ $? -eq 0 ] && [ -n "$OUTPUT" ]; then
    echo -e "${GREEN}PASS${NC}: folders list (text)"
else
    echo -e "${RED}FAIL${NC}: folders list (text)"
    echo "$OUTPUT"
    exit 1
fi

echo "Test: List folders (JSON output)..."
OUTPUT=$($CLI folders list --output json 2>&1)
if [ $? -eq 0 ] && echo "$OUTPUT" | grep -q '\['; then
    echo -e "${GREEN}PASS${NC}: folders list (json)"
else
    echo -e "${RED}FAIL${NC}: folders list (json)"
    echo "$OUTPUT"
    exit 1
fi

FIRST_FOLDER=$(echo "$OUTPUT" | grep -o '"ID":"[^"]*"' | head -1 | sed 's/"ID":"//;s/"//')
if [ -n "$FIRST_FOLDER" ]; then
    echo ""
    echo "First folder ID: $FIRST_FOLDER"
    echo "export CLICKUP_TEST_FOLDER_ID=\"$FIRST_FOLDER\"" > "$SCRIPT_DIR/.test-env"
fi

echo ""
echo -e "${GREEN}=== Folders tests complete ===${NC}"
