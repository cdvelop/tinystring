#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BENCHMARK_DIR="$(dirname "$SCRIPT_DIR")"

echo "📊 Updating README with latest benchmark data..."

cd "$BENCHMARK_DIR"
go run benchmark.go

echo "✅ README updated successfully!"
