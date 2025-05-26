#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BENCHMARK_DIR="$(dirname "$SCRIPT_DIR")"
EXAMPLES_DIR="$BENCHMARK_DIR/examples"

echo "🧹 Cleaning benchmark files..."

# Remove all generated binaries
find "$EXAMPLES_DIR" -name "*.exe" -delete 2>/dev/null || true
find "$EXAMPLES_DIR" -name "*.wasm" -delete 2>/dev/null || true
find "$EXAMPLES_DIR" -name "standard*" ! -name "*.go" ! -name "*.mod" -delete 2>/dev/null || true
find "$EXAMPLES_DIR" -name "tinystring*" ! -name "*.go" ! -name "*.mod" -delete 2>/dev/null || true

# Clean go modules cache in examples
cd "$EXAMPLES_DIR/standard-lib" && go clean 2>/dev/null || true
cd "$EXAMPLES_DIR/tinystring-lib" && go clean 2>/dev/null || true

echo "✅ Cleanup completed!"
