#!/bin/bash

set -e

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjBkMjkzYzdjLWU1OWMtNDBlMi1iMThlLWY5Y2MxMzVlM2VlNyIsImV4cGlyZXNBdCI6IjIwMjYtMDktMTlUMDM6NDM6MzEuMzA1WiIsInR5cGUiOiJzZXNzaW9uIiwiaWF0IjoxNzgxODQwNjExfQ.MS-LseF7tvNJTei-PfKdHomVKh6Xr7lYHj41pp8adbw"
BIN="./bin/outline"
TEST_DIR="/tmp/outline-test-$(date +%s)"

echo "=== Outline CLI - All Scenarios Test ==="
echo "Test directory: $TEST_DIR"
echo

# Test 1: Clone with document URL
echo "Test 1: Clone with document URL"
OUTLINE_TOKEN="$TOKEN" $BIN clone https://outline-rbi.jatismobile.com/doc/naufal-fsNXr2Zxj4 $TEST_DIR/test1 > /dev/null 2>&1
if [ -f "$TEST_DIR/test1/naufal.md" ]; then
    echo "✓ Clone with document URL: SUCCESS"
else
    echo "✗ Clone with document URL: FAILED"
    exit 1
fi

# Test 2: Clone with collection UUID
echo "Test 2: Clone with collection UUID"
OUTLINE_TOKEN="$TOKEN" $BIN clone 25299a17-a07d-48d2-b0df-4c5b7827a719 $TEST_DIR/test2 > /dev/null 2>&1
if [ -f "$TEST_DIR/test2/naufal.md" ]; then
    echo "✓ Clone with collection UUID: SUCCESS"
else
    echo "✗ Clone with collection UUID: FAILED"
    exit 1
fi

# Test 3: Clone with collection URL
echo "Test 3: Clone with collection URL"
OUTLINE_TOKEN="$TOKEN" $BIN clone https://outline-rbi.jatismobile.com/collection/catatan-dev-LtgQVZCHeI $TEST_DIR/test3 > /dev/null 2>&1
if [ -f "$TEST_DIR/test3/naufal.md" ]; then
    echo "✓ Clone with collection URL: SUCCESS"
else
    echo "✗ Clone with collection URL: FAILED"
    exit 1
fi

# Test 4: Verify frontmatter
echo "Test 4: Verify frontmatter"
if grep -q "outline_id:" "$TEST_DIR/test1/naufal.md" && grep -q "outline_collection:" "$TEST_DIR/test1/naufal.md"; then
    echo "✓ Frontmatter injection: SUCCESS"
else
    echo "✗ Frontmatter injection: FAILED"
    exit 1
fi

# Test 5: Verify manifest
echo "Test 5: Verify manifest"
if [ -f "$TEST_DIR/test1/.outline/manifest.json" ]; then
    COUNT=$(cat "$TEST_DIR/test1/.outline/manifest.json" | python3 -c "import sys,json; print(len(json.load(sys.stdin)))")
    if [ "$COUNT" -gt 0 ]; then
        echo "✓ Manifest generation: SUCCESS ($COUNT entries)"
    else
        echo "✗ Manifest generation: FAILED (empty)"
        exit 1
    fi
else
    echo "✗ Manifest generation: FAILED (not found)"
    exit 1
fi

# Cleanup
echo
echo "Cleaning up test directory..."
rm -rf "$TEST_DIR"

echo
echo "=== All tests PASSED ==="
