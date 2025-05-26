#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BENCHMARK_DIR="$(dirname "$SCRIPT_DIR")"
EXAMPLES_DIR="$BENCHMARK_DIR/examples"

echo "🚀 Starting binary size benchmark..."

# Check if TinyGo is installed
if ! command -v tinygo &> /dev/null; then
    echo "❌ TinyGo is not installed. Building only standard Go binaries."
    echo "   Install TinyGo from: https://tinygo.org/getting-started/install/"
    TINYGO_AVAILABLE=false
else
    echo "✅ TinyGo found: $(tinygo version)"
    TINYGO_AVAILABLE=true
fi

# Clean previous files
echo "🧹 Cleaning previous files..."
find "$EXAMPLES_DIR" -name "*.exe" -delete 2>/dev/null || true
find "$EXAMPLES_DIR" -name "*.wasm" -delete 2>/dev/null || true
find "$EXAMPLES_DIR" -name "standard*" ! -name "*.go" ! -name "*.mod" -delete 2>/dev/null || true
find "$EXAMPLES_DIR" -name "tinystring*" ! -name "*.go" ! -name "*.mod" -delete 2>/dev/null || true

# Define optimization configurations
SUFFIXES=("" "-ultra" "-speed" "-debug")
declare -A OPT_FLAGS=(
    ["default"]=""
    ["ultra"]="-no-debug -panic=trap -scheduler=none -gc=leaking"
    ["speed"]="-opt=2"
    ["debug"]="-opt=0"
)
declare -A OPT_DESCRIPTIONS=(
    ["default"]="Default optimization (-opt=z)"
    ["ultra"]="Maximum size optimization"
    ["speed"]="Optimized for speed over size"
    ["debug"]="No optimization, best for debugging"
)

# Build standard library example
echo "📦 Building standard library example with multiple optimizations..."
cd "$EXAMPLES_DIR/standard-lib"

# Standard Go build (only default)
go build -ldflags="-s -w" -o standard main.go

# TinyGo builds with different optimizations
if [ "$TINYGO_AVAILABLE" = true ]; then
    for suffix in "" "-ultra" "-speed" "-debug"; do
        case $suffix in
            "")
                key="default"
                flags=""
                ;;
            "-ultra")
                key="ultra"
                flags="${OPT_FLAGS[ultra]}"
                ;;
            "-speed")
                key="speed"
                flags="${OPT_FLAGS[speed]}"
                ;;
            "-debug")
                key="debug"
                flags="${OPT_FLAGS[debug]}"
                ;;
        esac
        
        echo "  Building with optimization: $suffix (${OPT_DESCRIPTIONS[$key]})"
        
        if [ -z "$flags" ]; then
            tinygo build -o "standard${suffix}.wasm" -target wasm main.go
        else
            tinygo build $flags -o "standard${suffix}.wasm" -target wasm main.go
        fi
    done
    echo "✅ Standard: Go binary and WebAssembly variants created"
else
    echo "⚠️  Standard: only Go binary created"
fi

# Build TinyString example
echo "📦 Building TinyString example with multiple optimizations..."
cd "$EXAMPLES_DIR/tinystring-lib"
go mod tidy

# Standard Go build (only default)
go build -ldflags="-s -w" -o tinystring main.go

# TinyGo builds with different optimizations
if [ "$TINYGO_AVAILABLE" = true ]; then
    for suffix in "" "-ultra" "-speed" "-debug"; do
        case $suffix in
            "")
                key="default"
                flags=""
                ;;
            "-ultra")
                key="ultra"
                flags="${OPT_FLAGS[ultra]}"
                ;;
            "-speed")
                key="speed"
                flags="${OPT_FLAGS[speed]}"
                ;;
            "-debug")
                key="debug"
                flags="${OPT_FLAGS[debug]}"
                ;;
        esac
        
        echo "  Building with optimization: $suffix (${OPT_DESCRIPTIONS[$key]})"
        
        if [ -z "$flags" ]; then
            tinygo build -o "tinystring${suffix}.wasm" -target wasm main.go
        else
            tinygo build $flags -o "tinystring${suffix}.wasm" -target wasm main.go
        fi
    done
    echo "✅ TinyString: Go binary and WebAssembly variants created"
else
    echo "⚠️  TinyString: only Go binary created"
fi

# Run analysis and update
echo "📊 Analyzing sizes and updating README..."
cd "$BENCHMARK_DIR"
go run benchmark.go

echo ""
echo "🎉 Benchmark completed successfully!"
echo ""
echo "📁 Generated files:"
find "$EXAMPLES_DIR" -name "*.exe" -o -name "*.wasm" -o -name "standard" -o -name "tinystring" | while read file; do
    if [[ -f "$file" ]]; then
        if command -v numfmt &> /dev/null; then
            size=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null || echo "0")
            echo "  $(basename "$file"): $(numfmt --to=iec-i --suffix=B $size)"
        else
            size=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null || wc -c < "$file")
            echo "  $(basename "$file"): ${size} bytes"
        fi
    fi
done
