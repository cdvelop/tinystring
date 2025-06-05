#!/bin/bash

# update-readme.sh - Update README.md with latest benchmark results
# This script runs benchmarks and updates documentation automatically

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Function to get the correct analyzer binary name based on OS
get_analyzer_name() {
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
        echo "analyzer.exe"
    else
        echo "analyzer"
    fi
}

ANALYZER_BINARY=$(get_analyzer_name)

echo "📚 Updating TinyString Benchmark Documentation"
echo "=============================================="

# Check if analyzer exists
if [[ ! -f "$ANALYZER_BINARY" ]]; then
    echo "📋 Building analyzer tool..."
    if ! go build -o "$ANALYZER_BINARY" *.go; then
        echo "❌ Failed to build analyzer tool"
        exit 1
    fi
    echo "✅ Analyzer built successfully"
fi

# Function to update with binary results
update_binary_results() {
    echo "📋 Updating README with binary size results..."
    
    if [[ ! -d "bench-binary-size" ]]; then
        echo "⚠️  Binary benchmark directory not found"
        echo "ℹ️  Run build-and-measure.sh first to generate binary samples"
        return 1
    fi
    
    # Check if binaries exist
    binary_count=$(find bench-binary-size -type f \( -name "standard*" -o -name "tinystring*" \) ! -name "*.go" ! -name "*.mod" | wc -l)
    if [[ $binary_count -eq 0 ]]; then
        echo "⚠️  No binary files found in bench-binary-size"
        echo "ℹ️  Run build-and-measure.sh first to generate binary samples"
        return 1
    fi
    
    if ./"$ANALYZER_BINARY" binary; then
        echo "✅ README updated with binary size results"
        return 0
    else
        echo "❌ Failed to update README with binary results"
        return 1
    fi
}

# Function to update with memory results
update_memory_results() {
    echo "📋 Updating README with memory allocation results..."
    
    if [[ ! -d "bench-memory-alloc" ]]; then
        echo "⚠️  Memory benchmark directory not found"
        echo "ℹ️  Memory benchmark structure may need setup"
        return 1
    fi
    
    if ./"$ANALYZER_BINARY" memory; then
        echo "✅ README updated with memory allocation results"
        return 0
    else
        echo "❌ Failed to update README with memory results"
        return 1
    fi
}

# Function to update with all results
update_all_results() {
    echo "� Updating README with complete benchmark results..."
    
    if ./"$ANALYZER_BINARY" all; then
        echo "✅ README updated with complete benchmark results"
        return 0
    else
        echo "❌ Failed to update README with complete results"
        return 1
    fi
}

# Function to validate README
validate_readme() {
    echo "📋 Validating README.md..."
    
    if [[ ! -f "README.md" ]]; then
        echo "❌ README.md not found"
        return 1
    fi
    
    # Check for required sections
    if ! grep -q "## Binary Size Comparison" README.md; then
        echo "⚠️  Binary Size Comparison section not found in README"
    fi
    
    if ! grep -q "## Memory Usage Comparison" README.md; then
        echo "⚠️  Memory Usage Comparison section not found in README"
    fi
    
    # Check for recent updates
    current_date=$(date +"%Y-%m-%d")
    if grep -q "Last updated: $current_date" README.md; then
        echo "✅ README contains today's results"
    else
        echo "ℹ️  README may not contain today's benchmark results"
    fi
    
    return 0
}

# Function to backup README
backup_readme() {
    if [[ -f "README.md" ]]; then
        backup_file="README.md.backup.$(date +%Y%m%d-%H%M%S)"
        cp README.md "$backup_file"
        echo "💾 README backed up to: $backup_file"
    fi
}

# Function to show status
show_status() {
    echo ""
    echo "📊 Benchmark Status Summary:"
    echo "============================"
    
    # Check binary benchmarks
    if [[ -d "bench-binary-size" ]]; then
        binary_count=$(find bench-binary-size -type f \( -name "standard*" -o -name "tinystring*" \) ! -name "*.go" ! -name "*.mod" | wc -l)
        echo "📦 Binary samples: $binary_count files found"
    else
        echo "📦 Binary samples: ❌ Directory not found"
    fi
    
    # Check memory benchmarks
    if [[ -d "bench-memory-alloc" ]]; then
        echo "🧠 Memory benchmarks: ✅ Directory exists"
    else
        echo "🧠 Memory benchmarks: ❌ Directory not found"
    fi
    
    # Check README
    if [[ -f "README.md" ]]; then
        echo "📚 README.md: ✅ File exists"
        if grep -q "Last updated:" README.md; then
            last_update=$(grep "Last updated:" README.md | head -1 | sed 's/.*Last updated: //' | sed 's/\*.*//')
            echo "📅 Last update: $last_update"
        fi
    else
        echo "📚 README.md: ❌ File not found"
    fi
}

# Parse command line arguments
case "${1:-all}" in
    "binary")
        backup_readme
        update_binary_results
        ;;
    "memory")
        backup_readme
        update_memory_results
        ;;
    "all")
        backup_readme
        update_all_results
        ;;
    "validate")
        validate_readme
        ;;
    "status")
        show_status
        exit 0
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [binary|memory|all|validate|status|help]"
        echo ""
        echo "Commands:"
        echo "  binary    - Update README with binary size results only"
        echo "  memory    - Update README with memory allocation results only"
        echo "  all       - Update README with complete benchmark results (default)"
        echo "  validate  - Validate README.md structure and content"
        echo "  status    - Show current benchmark status"
        echo "  help      - Show this help message"
        echo ""
        echo "Note: A backup of README.md will be created before updates"
        echo ""
        echo "Prerequisites:"
        echo "  - Run build-and-measure.sh first for binary benchmarks"
        echo "  - Memory benchmark directories should be properly set up"
        exit 0
        ;;
    *)
        echo "❌ Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac

RESULT=$?

echo ""
if [[ $RESULT -eq 0 ]]; then
    echo "✅ README update completed successfully!"
    validate_readme
else
    echo "❌ README update completed with errors"
    echo "ℹ️  Check the output above for details"
fi

echo "=============================================="
exit $RESULT
