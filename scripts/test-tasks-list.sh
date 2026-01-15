#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CLI="$PROJECT_DIR/clickup"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=== ClickUp CLI Integration Test: List Tasks ==="
echo ""

if [ ! -f "$CLI" ]; then
    echo "Building CLI..."
    cd "$PROJECT_DIR"
    go build -o clickup ./cmd/clickup
fi

if [ -f "$SCRIPT_DIR/.test-env" ]; then
    source "$SCRIPT_DIR/.test-env"
fi

if [ -z "$CLICKUP_TEST_LIST_ID" ]; then
    echo "Getting first folder and list..."
    FOLDER_OUTPUT=$($CLI folders list --output json 2>&1)
    FOLDER_ID=$(echo "$FOLDER_OUTPUT" | grep -o '"ID":"[^"]*"' | head -1 | sed 's/"ID":"//;s/"//')

    if [ -n "$FOLDER_ID" ]; then
        LIST_OUTPUT=$($CLI lists list --folder "$FOLDER_ID" --output json 2>&1)
        CLICKUP_TEST_LIST_ID=$(echo "$LIST_OUTPUT" | grep -o '"ID":"[^"]*"' | head -1 | sed 's/"ID":"//;s/"//')
    fi
fi

if [ -z "$CLICKUP_TEST_LIST_ID" ]; then
    echo -e "${YELLOW}SKIP${NC}: No lists found, cannot test tasks"
    exit 0
fi

echo "Using list ID: $CLICKUP_TEST_LIST_ID"
echo ""

echo "Test: List tasks (text output)..."
OUTPUT=$($CLI tasks list --list "$CLICKUP_TEST_LIST_ID" 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: tasks list (text)"
else
    echo -e "${RED}FAIL${NC}: tasks list (text)"
    echo "$OUTPUT"
    exit 1
fi

echo "Test: List tasks (JSON output)..."
OUTPUT=$($CLI tasks list --list "$CLICKUP_TEST_LIST_ID" --output json 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: tasks list (json)"
else
    echo -e "${RED}FAIL${NC}: tasks list (json)"
    echo "$OUTPUT"
    exit 1
fi

echo "Test: List tasks recursive (text output)..."
OUTPUT=$($CLI tasks list --list "$CLICKUP_TEST_LIST_ID" --recursive 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: tasks list --recursive (text)"
else
    echo -e "${RED}FAIL${NC}: tasks list --recursive (text)"
    echo "$OUTPUT"
    exit 1
fi

echo "Test: List tasks recursive (JSON output)..."
OUTPUT=$($CLI tasks list --list "$CLICKUP_TEST_LIST_ID" --recursive --output json 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: tasks list --recursive (json)"
else
    echo -e "${RED}FAIL${NC}: tasks list --recursive (json)"
    echo "$OUTPUT"
    exit 1
fi

echo ""
echo -e "${GREEN}=== List Tasks tests complete ===${NC}"
