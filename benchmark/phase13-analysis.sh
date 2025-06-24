#!/bin/bash

# TinyString Phase 13 Memory Analysis Script
# Automated analysis for allocation optimization progress

set -e

echo "🎯 TinyString Phase 13 - Performance Recovery Analysis"
echo "=================================================="

WORKSPACE_DIR="/c/Users/Cesar/Packages/Internal/tinystring"
BENCHMARK_DIR="$WORKSPACE_DIR/benchmark"
RESULTS_DIR="$BENCHMARK_DIR/phase13-analysis"

# Create results directory
mkdir -p "$RESULTS_DIR"

cd "$WORKSPACE_DIR"

echo "📊 Step 1: Baseline Performance Measurement"
echo "-------------------------------------------"

# Run benchmarks with memory profiling
echo "Running memory allocation benchmarks..."
go test -bench=BenchmarkStringOperations -benchmem -memprofile="$RESULTS_DIR/memory_baseline.prof" > "$RESULTS_DIR/benchmark_baseline.txt" 2>&1

echo "Running number processing benchmarks..."  
go test -bench=BenchmarkNumberProcessing -benchmem -memprofile="$RESULTS_DIR/memory_numbers.prof" >> "$RESULTS_DIR/benchmark_baseline.txt" 2>&1

echo "Running builder operations benchmarks..."
go test -bench=BenchmarkBuilderOperations -benchmem >> "$RESULTS_DIR/benchmark_baseline.txt" 2>&1

echo "🔍 Step 2: Memory Allocation Analysis"
echo "-------------------------------------"

# Analyze memory allocations with pprof
echo "Analyzing memory hotspots..."
go tool pprof -text "$RESULTS_DIR/memory_baseline.prof" | head -30 > "$RESULTS_DIR/memory_hotspots.txt"

echo "Top allocation sources:" 
cat "$RESULTS_DIR/memory_hotspots.txt" | head -15

echo "🔬 Step 3: Escape Analysis"
echo "-------------------------"

# Analyze what escapes to heap
echo "Analyzing heap escapes..."
go build -gcflags="-m=2" . 2>&1 | grep -E "(escapes|moved to heap|too large)" > "$RESULTS_DIR/escape_analysis.txt" || true

echo "Top escape analysis results:"
head -20 "$RESULTS_DIR/escape_analysis.txt"

echo "🧪 Step 4: String Allocation Detection"  
echo "--------------------------------------"

# Focus on string allocations specifically
echo "Analyzing string allocations..."
go build -gcflags="-m=3" . 2>&1 | grep -E "string.*allocation|getString.*alloc|wrString.*alloc" > "$RESULTS_DIR/string_allocations.txt" || true

if [ -s "$RESULTS_DIR/string_allocations.txt" ]; then
    echo "String allocation hotspots found:"
    cat "$RESULTS_DIR/string_allocations.txt"
else
    echo "No explicit string allocation warnings found in escape analysis."
fi

echo "🎯 Step 5: Performance Baseline Summary"
echo "---------------------------------------"

# Extract key metrics from benchmark results
echo "Extracting baseline metrics..."

# Parse benchmark results for key metrics
grep -E "BenchmarkStringOperations.*allocs/op|BenchmarkBuilderOperations.*allocs/op" "$RESULTS_DIR/benchmark_baseline.txt" | head -10 > "$RESULTS_DIR/baseline_metrics.txt"

echo "Current Phase 12 Performance Baseline:"
echo "======================================"
cat "$RESULTS_DIR/baseline_metrics.txt"

# Calculate averages (simplified)
TOTAL_MEMORY=$(grep -oE "[0-9]+ B/op" "$RESULTS_DIR/baseline_metrics.txt" | grep -oE "[0-9]+" | awk '{sum+=$1; count++} END {if(count>0) printf "%.0f", sum/count}')
TOTAL_ALLOCS=$(grep -oE "[0-9]+ allocs/op" "$RESULTS_DIR/baseline_metrics.txt" | grep -oE "[0-9]+" | awk '{sum+=$1; count++} END {if(count>0) printf "%.0f", sum/count}')

echo ""
echo "📈 Baseline Averages:"
echo "Memory/op: ${TOTAL_MEMORY:-"N/A"} B"  
echo "Allocs/op: ${TOTAL_ALLOCS:-"N/A"}"

echo ""
echo "🎯 Phase 13 Targets:"
echo "Memory/op: 600 B (Target: 20% improvement)"
echo "Allocs/op: 38 (Target: 32% improvement)"
echo "Speed: 2900 ns/op (Target: 11% improvement)"

echo ""
echo "📁 Analysis files generated in: $RESULTS_DIR"
echo "   - benchmark_baseline.txt    (Full benchmark results)"
echo "   - memory_hotspots.txt       (Memory allocation analysis)" 
echo "   - escape_analysis.txt       (Heap escape analysis)"
echo "   - string_allocations.txt    (String-specific allocations)"
echo "   - baseline_metrics.txt      (Key performance metrics)"

echo ""
echo "🔧 Next Steps:"
echo "1. Review memory_hotspots.txt for optimization targets"
echo "2. Implement Stage 1: String Caching optimization"
echo "3. Re-run this script after each optimization stage"
echo "4. Compare results with: benchstat baseline.txt optimized.txt"

echo ""
echo "✅ Phase 13 baseline analysis complete!"
