# TinyString - Benchmark Summary for LLM Context

## Binary Size Results - After Type Consolidation Refactoring

**WebAssembly Performance (Post-Refactoring):**
- Peak reduction: 52.8% (Ultra WASM build: 200.6 KB → 94.8 KB)
- Average WebAssembly reduction: 31.4%
- Default WASM: 23.6% reduction (879.1 KB → 671.9 KB)
- Total savings across all builds: 1.4 MB

**Native Performance (Post-Refactoring):**
- Native builds: 4.2% reduction (1.6 MB → 1.5 MB)

## Memory Usage Results - After Type Consolidation Refactoring

**Runtime Memory Impact:**
- Memory overhead: +103.3% average (significantly higher memory usage)
- Allocation efficiency: -7.4% average (slightly fewer allocations)
- Execution time: Generally 2-4x slower than standard library

**Key Findings:**
- String processing: +96.7% memory, -4.2% allocations
- Number processing: +110.7% memory, -9.1% allocations  
- Mixed operations: +119.8% memory, +4.5% allocations
- Pointer optimization: +86.0% memory, -20.8% allocations

## Refactoring Impact Analysis

**Code Quality Improvements:**
- Consolidated 42 repetitive type case statements into 12 helper functions
- Reduced code duplication by approximately 60-70% in core functions
- Improved maintainability through centralized type handling

**Performance Improvements from Refactoring:**
- Binary size: **IMPROVED** - Peak reduction increased from 50.7% to 52.8%
- Default WASM: **IMPROVED** - Reduction increased from 21.6% to 23.6%
- Native builds: **IMPROVED** - Reduction increased from 3.1% to 4.2%
- Memory patterns: **MAINTAINED** - Similar allocation patterns preserved

**Functions Successfully Refactored:**
- `convInit()`: 10 type cases → 3 helper calls
- `formatValue()`: 10 type cases → 3 helper calls  
- `any2s()`: 10 type cases → 3 helper calls
- `toInt()`: 12 type cases → 3 helper calls

## Key Insights

- Primary benefit: Dramatic WebAssembly binary size reduction (31.4% average post-refactoring)
- Primary cost: 2-4x higher runtime memory usage and slower execution
- Refactoring impact: **Positive** - Improved binary size without memory regression
- Code quality: **Significantly improved** - Major reduction in code duplication
- Optimization target: Distribution size over runtime efficiency
- Best fit: Size-constrained deployment environments
