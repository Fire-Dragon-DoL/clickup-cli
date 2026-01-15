#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CLI="$PROJECT_DIR/clickup"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=== ClickUp CLI Integration Test: Lists ==="
echo ""

if [ ! -f "$CLI" ]; then
    echo "Building CLI..."
    cd "$PROJECT_DIR"
    go build -o clickup ./cmd/clickup
fi

if [ -f "$SCRIPT_DIR/.test-env" ]; then
    source "$SCRIPT_DIR/.test-env"
fi

if [ -z "$CLICKUP_TEST_FOLDER_ID" ]; then
    echo "Getting first folder..."
    FOLDER_OUTPUT=$($CLI folders list --output json 2>&1)
    CLICKUP_TEST_FOLDER_ID=$(echo "$FOLDER_OUTPUT" | grep -o '"ID":"[^"]*"' | head -1 | sed 's/"ID":"//;s/"//')
fi

if [ -z "$CLICKUP_TEST_FOLDER_ID" ]; then
    echo -e "${YELLOW}SKIP${NC}: No folders found, cannot test lists"
    exit 0
fi

echo "Using folder ID: $CLICKUP_TEST_FOLDER_ID"
echo ""

echo "Test: List lists in folder (text output)..."
OUTPUT=$($CLI lists list --folder "$CLICKUP_TEST_FOLDER_ID" 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: lists list (text)"
else
    echo -e "${RED}FAIL${NC}: lists list (text)"
    echo "$OUTPUT"
    exit 1
fi

echo "Test: List lists in folder (JSON output)..."
OUTPUT=$($CLI lists list --folder "$CLICKUP_TEST_FOLDER_ID" --output json 2>&1)
if [ $? -eq 0 ] && echo "$OUTPUT" | grep -q '\['; then
    echo -e "${GREEN}PASS${NC}: lists list (json)"
else
    echo -e "${RED}FAIL${NC}: lists list (json)"
    echo "$OUTPUT"
    exit 1
fi

FIRST_LIST=$(echo "$OUTPUT" | grep -o '"ID":"[^"]*"' | head -1 | sed 's/"ID":"//;s/"//')
if [ -n "$FIRST_LIST" ]; then
    echo ""
    echo "First list ID: $FIRST_LIST"
    echo "export CLICKUP_TEST_FOLDER_ID=\"$CLICKUP_TEST_FOLDER_ID\"" > "$SCRIPT_DIR/.test-env"
    echo "export CLICKUP_TEST_LIST_ID=\"$FIRST_LIST\"" >> "$SCRIPT_DIR/.test-env"
fi

echo ""
echo -e "${GREEN}=== Lists tests complete ===${NC}"
