#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=========================================="
echo "  ClickUp CLI Integration Test Suite"
echo "=========================================="
echo ""

cd "$PROJECT_DIR"

echo "Building CLI..."
go build -o clickup ./cmd/clickup
echo ""

PASSED=0
FAILED=0
SKIPPED=0

run_test() {
    local script=$1
    local name=$2

    echo "----------------------------------------"
    echo "Running: $name"
    echo "----------------------------------------"

    if bash "$script"; then
        ((PASSED++))
    else
        EXIT_CODE=$?
        if [ $EXIT_CODE -eq 0 ]; then
            ((SKIPPED++))
        else
            ((FAILED++))
        fi
    fi
    echo ""
}

run_test "$SCRIPT_DIR/test-setup.sh" "Setup Verification"
run_test "$SCRIPT_DIR/test-folders.sh" "Folders"
run_test "$SCRIPT_DIR/test-lists.sh" "Lists"
run_test "$SCRIPT_DIR/test-tasks-list.sh" "List Tasks"
run_test "$SCRIPT_DIR/test-task-crud.sh" "Task CRUD"

rm -f "$SCRIPT_DIR/.test-env"

echo "=========================================="
echo "  Test Results"
echo "=========================================="
echo -e "Passed:  ${GREEN}$PASSED${NC}"
echo -e "Failed:  ${RED}$FAILED${NC}"
echo -e "Skipped: ${YELLOW}$SKIPPED${NC}"
echo "=========================================="

if [ $FAILED -gt 0 ]; then
    exit 1
fi
