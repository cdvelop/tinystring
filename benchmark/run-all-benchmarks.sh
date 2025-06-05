#!/bin/bash

# run-all-benchmarks.sh - Run all TinyString benchmarks and generate reports
# This script runs both binary size analysis and memory allocation benchmarks

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

echo "🚀 Starting TinyString Benchmark Suite"
echo "======================================="

# Build the analyzer tool
echo "📋 Building analyzer tool..."
if ! go build -o "$ANALYZER_BINARY" *.go; then
    echo "❌ Failed to build analyzer tool"
    exit 1
fi
echo "✅ Analyzer built successfully"

# Function to run binary size benchmarks
run_binary_benchmarks() {
    echo ""
    echo "📋 Running binary size benchmarks..."
    
    # Check if binary directories exist
    if [[ ! -d "bench-binary-size" ]]; then
        echo "❌ Binary benchmark directory not found: bench-binary-size"
        echo "ℹ️  Run build-and-measure.sh first to create binary samples"
        return 1
    fi
    
    # Run binary analysis
    if ./"$ANALYZER_BINARY" binary; then
        echo "✅ Binary size analysis completed"
        return 0
    else
        echo "❌ Binary size analysis failed"
        return 1
    fi
}

# Function to run memory allocation benchmarks
run_memory_benchmarks() {
    echo ""
    echo "📋 Running memory allocation benchmarks..."
    
    # Check if memory benchmark directories exist
    if [[ ! -d "bench-memory-alloc" ]]; then
        echo "❌ Memory benchmark directory not found: bench-memory-alloc"
        echo "ℹ️  Memory benchmarks directory structure may need setup"
        return 1
    fi
    
    # Run memory analysis
    if ./"$ANALYZER_BINARY" memory; then
        echo "✅ Memory allocation analysis completed"
        return 0
    else
        echo "❌ Memory allocation analysis failed"
        return 1
    fi
}

# Function to run all benchmarks
run_all_benchmarks() {
    echo ""
    echo "📋 Running complete benchmark suite..."
    
    if ./"$ANALYZER_BINARY" all; then
        echo "✅ Complete benchmark suite completed"
        return 0
    else
        echo "❌ Complete benchmark suite failed"
        return 1
    fi
}

# Parse command line arguments
case "${1:-all}" in
    "binary")
        run_binary_benchmarks
        ;;
    "memory")
        run_memory_benchmarks
        ;;
    "all")
        run_all_benchmarks
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [binary|memory|all|help]"
        echo ""
        echo "Commands:"
        echo "  binary  - Run only binary size benchmarks"
        echo "  memory  - Run only memory allocation benchmarks" 
        echo "  all     - Run all benchmarks (default)"
        echo "  help    - Show this help message"
        echo ""
        echo "Prerequisites:"
        echo "  - Go 1.19+ installed"
        echo "  - TinyGo installed (for binary size benchmarks)"
        echo "  - Binary samples built (run build-and-measure.sh first)"
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
    echo "🎉 Benchmark suite completed successfully!"
    echo "📋 Results have been updated in README.md"
else
    echo "❌ Benchmark suite completed with errors"
    echo "ℹ️  Check the output above for details"
fi

echo "======================================="
exit $RESULT
