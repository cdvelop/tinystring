# TinyString - Benchmark Summary for LLM Context

## Binary Size Results - Current Performance Data

**WebAssembly Performance:**
- Peak reduction: 87.6% (Ultra WASM build: 200.6 KB → 24.9 KB)
- Average WebAssembly reduction: 80.3%
- Default WASM: 76.5% reduction (879.1 KB → 206.2 KB)
- Speed WASM: 79.0% reduction (1.3 MB → 271.7 KB)
- Debug WASM: 78.3% reduction (3.0 MB → 666.1 KB)
- Total savings across all builds: 4.8 MB

**Native Performance:**
- Native builds: 38.1% reduction (1.6 MB → 983.5 KB)

## Memory Usage Results - Current Performance Data

**Runtime Memory Impact:**
- Memory overhead: +103.3% average (significantly higher memory usage)
- Allocation efficiency: -7.4% average (slightly fewer allocations)
- Execution time: Generally 2-4x slower than standard library

**Detailed Memory Benchmarks:**
- String processing: +96.7% memory, -4.2% allocations (2.3 KB vs 1.2 KB, 46 vs 48 allocs)
- Number processing: +110.7% memory, -9.1% allocations (2.5 KB vs 1.2 KB, 120 vs 132 allocs)
- Mixed operations: +119.8% memory, +4.5% allocations (1.2 KB vs 546 B, 46 vs 44 allocs)
- Pointer optimization: +86.0% memory, -20.8% allocations (2.2 KB vs 1.2 KB, 38 vs 48 allocs)

**Performance Categories:**
- Memory Efficiency: Poor (Significant overhead - 103.3% average increase)
- Allocation Efficiency: Good (Allocation efficient - 7.4% average decrease)
- Benchmarks Analyzed: 4 categories
- Optimization Focus: Binary size reduction vs runtime efficiency

## Performance Analysis Summary

**Trade-offs Analysis:**
The benchmarks reveal important trade-offs between binary size and runtime performance:

**Binary Size Benefits:**
- 38.1% smaller native compiled binaries (1.6 MB → 983.5 KB)
- Superior WebAssembly compression ratios (76.5-87.6% reduction)
- Faster deployment and distribution
- Lower storage requirements

**Runtime Memory Considerations:**
- Higher allocation overhead during execution (+103.3% average)
- Increased GC pressure due to allocation patterns
- Trade-off optimizes for distribution size over runtime efficiency
- Different optimization strategy than standard library

**Current Performance Achievements:**
- Peak WebAssembly reduction: 87.6% (Ultra WASM: 24.9 KB final size)
- Average WebAssembly reduction: 80.3%
- Total size savings across all builds: 4.8 MB
- Outstanding performance in size-constrained environments

## Key Insights

**Primary Benefits:**
- Dramatic WebAssembly binary size reduction (80.3% average)
- Outstanding peak performance: 87.6% size reduction (Ultra WASM)
- Excellent size optimization for web deployment and embedded systems
- Total size savings: 4.8 MB across all build configurations

**Primary Costs:**
- 2-4x higher runtime memory usage (+103.3% average)
- Slower execution performance compared to standard library
- Higher GC pressure due to allocation patterns

**Optimization Recommendations:**
- **Use TinyString for:** WebAssembly apps, embedded systems, edge computing, size-critical deployments
- **Use Standard Library for:** Memory-intensive servers, high-frequency processing, performance-critical workloads

**Performance Categories:**
- Binary Size: Outstanding (>70% reduction)
- Memory Efficiency: Poor (Significant overhead)
- Allocation Efficiency: Good (Fewer allocations)
- Target Use Case: Distribution size over runtime efficiency
