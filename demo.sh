#!/bin/bash
#
# demo.sh - Demo script for Outline CLI
#
set -e

echo "=== Outline CLI Demo ==="
echo ""

# Check if binary exists
if [ ! -f "./bin/outline" ]; then
    echo "Building binary..."
    make build
fi

echo "Binary size:"
ls -lh bin/outline
echo ""

echo "Version:"
./bin/outline --version
echo ""

echo "Available commands:"
./bin/outline --help
echo ""

# Test init
echo "=== Test 1: Initialize repository ==="
TEST_DIR="/tmp/outline-demo-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

echo "Running: outline init"
~/zero/research/outline-cli/bin/outline init
echo ""

echo "Repository structure:"
find .outline -type f
echo ""

# Test status
echo "=== Test 2: Status (stub) ==="
echo "Running: outline status"
~/zero/research/outline-cli/bin/outline status
echo ""

# Cleanup
echo "=== Cleanup ==="
cd -
rm -rf "$TEST_DIR"
echo "Test directory removed: $TEST_DIR"
echo ""

echo "=== Demo Complete ==="
echo ""
echo "Next steps:"
echo "1. Set OUTLINE_API_KEY environment variable"
echo "2. Run: outline clone <collection-id> <directory>"
echo "3. Implement Phase 3: Status command"
