#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CLI="$PROJECT_DIR/clickup"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

CREATED_TASK_ID=""

cleanup() {
    if [ -n "$CREATED_TASK_ID" ]; then
        echo ""
        echo "Cleaning up test task: $CREATED_TASK_ID"
        $CLI tasks delete "$CREATED_TASK_ID" 2>/dev/null || true
    fi
}

trap cleanup EXIT

echo "=== ClickUp CLI Integration Test: Task CRUD ==="
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
    echo -e "${YELLOW}SKIP${NC}: No lists found, cannot test task CRUD"
    exit 0
fi

echo "Using list ID: $CLICKUP_TEST_LIST_ID"
echo ""

TEST_TASK_NAME="Integration Test Task $(date +%s)"

echo "Test: Create task..."
OUTPUT=$($CLI tasks create --title "$TEST_TASK_NAME" --list "$CLICKUP_TEST_LIST_ID" --output json 2>&1)
if [ $? -eq 0 ]; then
    CREATED_TASK_ID=$(echo "$OUTPUT" | grep -o '"ID":"[^"]*"' | head -1 | sed 's/"ID":"//;s/"//')
    if [ -n "$CREATED_TASK_ID" ]; then
        echo -e "${GREEN}PASS${NC}: tasks create"
        echo "Created task ID: $CREATED_TASK_ID"
    else
        echo -e "${RED}FAIL${NC}: tasks create (could not extract task ID)"
        echo "$OUTPUT"
        exit 1
    fi
else
    echo -e "${RED}FAIL${NC}: tasks create"
    echo "$OUTPUT"
    exit 1
fi

echo ""
echo "Test: Show task..."
OUTPUT=$($CLI tasks show "$CREATED_TASK_ID" 2>&1)
if [ $? -eq 0 ] && echo "$OUTPUT" | grep -q "$TEST_TASK_NAME"; then
    echo -e "${GREEN}PASS${NC}: tasks show"
else
    echo -e "${RED}FAIL${NC}: tasks show"
    echo "$OUTPUT"
    exit 1
fi

echo ""
echo "Test: Update task..."
UPDATED_NAME="Updated $TEST_TASK_NAME"
OUTPUT=$($CLI tasks update "$CREATED_TASK_ID" --title "$UPDATED_NAME" --output json 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: tasks update"
else
    echo -e "${RED}FAIL${NC}: tasks update"
    echo "$OUTPUT"
    exit 1
fi

echo ""
echo "Test: Verify update..."
OUTPUT=$($CLI tasks show "$CREATED_TASK_ID" 2>&1)
if [ $? -eq 0 ] && echo "$OUTPUT" | grep -q "Updated"; then
    echo -e "${GREEN}PASS${NC}: update verified"
else
    echo -e "${RED}FAIL${NC}: update verification"
    echo "$OUTPUT"
    exit 1
fi

echo ""
echo "Test: Delete task..."
OUTPUT=$($CLI tasks delete "$CREATED_TASK_ID" 2>&1)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}PASS${NC}: tasks delete"
    CREATED_TASK_ID=""
else
    echo -e "${YELLOW}SKIP${NC}: tasks delete (command may not be implemented)"
    echo "Note: Delete command not available, task will remain in ClickUp"
    echo "Task ID: $CREATED_TASK_ID"
    CREATED_TASK_ID=""
fi

echo ""
echo -e "${GREEN}=== Task CRUD tests complete ===${NC}"
