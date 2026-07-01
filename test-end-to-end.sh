#!/bin/bash
set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║          Outline CLI - End-to-End Test Suite              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo

# Ensure token is set
if [ -z "$OUTLINE_TOKEN" ]; then
    echo "❌ Error: OUTLINE_TOKEN not set"
    exit 1
fi

echo "✓ Environment: OUTLINE_TOKEN is set"
echo

# Test 1: Binary exists and runs
echo "Test 1: Binary validation"
if [ -x "./bin/outline-linux-amd64" ]; then
    echo "  ✓ Binary exists and is executable"
else
    echo "  ❌ Binary not found or not executable"
    exit 1
fi
echo

# Test 2: Clone with document URL
echo "Test 2: Clone collection"
TEST_DIR="/tmp/outline-test-$(date +%s)"
./bin/outline-linux-amd64 clone 25299a17-a07d-48d2-b0df-4c5b7827a719 "$TEST_DIR" > /dev/null 2>&1
if [ -d "$TEST_DIR/.outline" ]; then
    echo "  ✓ Collection cloned successfully"
    DOC_COUNT=$(find "$TEST_DIR" -name "*.md" | wc -l)
    echo "  ✓ Found $DOC_COUNT documents"
else
    echo "  ❌ Clone failed"
    exit 1
fi
echo

# Test 3: Status check
echo "Test 3: Status command"
cd "$TEST_DIR"
STATUS_OUTPUT=$(../zero/research/outline-cli/bin/outline-linux-amd64 status 2>&1)
if echo "$STATUS_OUTPUT" | grep -q "working tree clean"; then
    echo "  ✓ Status shows clean working tree"
else
    echo "  ❌ Status check failed"
    exit 1
fi
echo

# Test 4: Add new file
echo "Test 4: Add new document"
echo "# End-to-End Test

This document was created during automated testing.

Timestamp: $(date -u +%Y-%m-%dT%H:%M:%S.000Z)
" > e2e-test.md

~/zero/research/outline-cli/bin/outline-linux-amd64 add e2e-test.md > /dev/null 2>&1
if ~/zero/research/outline-cli/bin/outline-linux-amd64 status | grep -q "e2e-test.md"; then
    echo "  ✓ File added to tracking"
else
    echo "  ❌ Add command failed"
    exit 1
fi
echo

# Test 5: Push (create)
echo "Test 5: Push new document"
PUSH_OUTPUT=$(~/zero/research/outline-cli/bin/outline-linux-amd64 push 2>&1)
if echo "$PUSH_OUTPUT" | grep -q "created"; then
    echo "  ✓ Document created successfully"
    if grep -q "outline_id:" e2e-test.md; then
        echo "  ✓ Frontmatter added automatically"
    else
        echo "  ❌ Frontmatter not added"
        exit 1
    fi
else
    echo "  ❌ Push failed"
    echo "$PUSH_OUTPUT"
    exit 1
fi
echo

# Test 6: Modify and push (update)
echo "Test 6: Update existing document"
echo "

## Updated Section

This section was added during update test.
" >> e2e-test.md

PUSH_UPDATE=$(~/zero/research/outline-cli/bin/outline-linux-amd64 push 2>&1)
if echo "$PUSH_UPDATE" | grep -q "updated"; then
    echo "  ✓ Document updated successfully"
else
    echo "  ❌ Update failed"
    exit 1
fi
echo

# Cleanup
cd ~
rm -rf "$TEST_DIR"

echo "╔════════════════════════════════════════════════════════════╗"
echo "║              ✅ ALL TESTS PASSED! ✅                        ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo
echo "Summary:"
echo "  ✓ Binary validation"
echo "  ✓ Clone collection"
echo "  ✓ Status command"
echo "  ✓ Add new document"
echo "  ✓ Push (create with parentDocumentId fallback)"
echo "  ✓ Push (update existing document)"
echo
echo "All features working correctly with session token!"
