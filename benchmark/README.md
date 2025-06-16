# TinyString Benchmark Suite

Automated benchmark tools to measure and compare performance between standard Go libraries and TinyString implementations.

## Quick Usage 🚀

```bash
# Run complete benchmark (recommended)
./build-and-measure.sh

# Clean generated files
./clean-all.sh

# Update README with existing data only (does not re-run benchmarks)
./update-readme.sh

# Run all memory and binary size benchmarks (without updating README)
./run-all-benchmarks.sh

# Run only memory benchmarks
./memory-benchmark.sh
```

## What Gets Measured 📊

1.  **Binary Size Comparison**: Native + WebAssembly builds with multiple optimization levels. This compares the compiled output size of projects using the standard Go library versus TinyString.
2.  **Memory Allocation**: Measures Bytes/op, Allocations/op, and execution time (ns/op) for benchmark categories. This helps in understanding the memory efficiency of TinyString compared to standard library operations.
    *   **String Processing**: Benchmarks operations like case conversion, text manipulation, etc.
    *   **Number Processing**: Benchmarks numeric formatting, conversion operations, etc.
    *   **Mixed Operations**: Benchmarks scenarios involving a combination of string and numeric operations.

## Current Performance Status

**Target**: Achieve memory usage close to standard library while maintaining binary size benefits.

**Latest Results** (Run `./build-and-measure.sh` to update):
- ✅ **Binary Size**: TinyString is 20-50% smaller than stdlib for WebAssembly.
- ⚠️ **Memory Usage**: Number Processing uses 1000% more memory (needs optimization).

📋 **Memory Optimization Guide**: See [`MEMORY_REDUCTION.md`](./MEMORY_REDUCTION.md) for comprehensive techniques and best practices to replace Go standard libraries with TinyString's optimized implementations. Essential reading for efficient string and numeric processing in TinyGo WebAssembly applications.

## Requirements

- **Go 1.21+**
- **TinyGo** (optional, but recommended for full WebAssembly testing and to achieve smallest binary sizes).

## Directory Structure

```
benchmark/
├── analyzer.go               # Main analysis program that processes benchmark results and generates reports.
├── common.go                 # Shared utilities used by benchmark scripts and analysis tools.
├── reporter.go               # Logic for formatting and updating the README.md with benchmark results.
├── MEMORY_REDUCTION.md       # Detailed guide for memory optimization techniques in TinyGo applications.
├── build-and-measure.sh      # 🎯 MAIN SCRIPT: Comprehensive benchmark that builds binaries, measures sizes, 
│                             #    runs memory tests, and updates README.md with latest results.
├── memory-benchmark.sh       # Executes only memory allocation benchmarks without building binaries or 
│                             #    updating documentation. Useful for focused memory optimization work.
├── clean-all.sh              # Cleanup script that removes all generated binaries (.exe, .wasm) and 
│                             #    temporary analysis files to free disk space.
├── update-readme.sh          # Updates benchmark sections in README.md using existing data files without 
│                             #    re-running benchmarks. Only reformats previously generated results.
├── run-all-benchmarks.sh     # Executes all benchmark tests (binary size + memory allocation) and generates 
│                             #    raw data files but does NOT update the README.md automatically.
├── validate-shared-data.sh   # Validation script that ensures test data consistency across all benchmark suites.
├── shared/                   # 🔄 SHARED TEST DATA: Centralized test data for consistent benchmarking.
│   ├── go.mod               #    Module definition for shared data package used by all benchmarks.
│   └── testdata.go          #    Common test data (TestTexts, TestNumbers, TestMixedData) ensuring 
│                             #    identical inputs for fair TinyString vs standard library comparisons.
├── bench-binary-size/        # Binary size comparison projects for measuring compiled output sizes.
│   ├── standard-lib/         #    Example project using only standard Go library functions.
│   │   ├── go.mod           #    Module with standard library dependencies.
│   │   └── main.go          #    Implementation using fmt, strconv, strings packages.
│   └── tinystring-lib/       #    Equivalent project using TinyString library instead.
│       ├── go.mod           #    Module with TinyString dependency and local replace directive.
│       └── main.go          #    Implementation using TinyString functions (same logic as standard-lib).
└── bench-memory-alloc/       # Memory allocation benchmark suites for runtime performance comparison.
    ├── standard/             #    Memory benchmarks using standard Go library (fmt, strconv, strings).
    │   ├── go.mod           #    Module with shared data dependency and standard library imports.
    │   ├── main.go          #    Processing functions using standard library implementations.
    │   └── main_test.go     #    Benchmark tests measuring Bytes/op, Allocs/op, ns/op for standard lib.
    └── tinystring/           #    Equivalent memory benchmarks using TinyString library functions.
        ├── go.mod           #    Module with TinyString and shared data dependencies.
        ├── main.go          #    Processing functions using TinyString implementations (same logic).        └── main_test.go     #    Benchmark tests measuring memory metrics for TinyString (identical to standard).
```

## Example Output

```
🚀 Starting binary size benchmark...
✅ TinyGo found: tinygo version 0.37.0
🧹 Cleaning previous files...
📦 Building standard library example with multiple optimizations...
📦 Building TinyString example with multiple optimizations...
📊 Analyzing sizes and updating README...
🧠 Running memory allocation benchmarks...
✅ Binary size analysis completed and README updated
✅ Memory benchmarks completed and README updated

🎉 Benchmark completed successfully!

📁 Generated files:
  standard: 1.3MiB
  tinystring: 1.1MiB  
  standard.wasm: 581KiB
  tinystring.wasm: 230KiB
  standard-ultra.wasm: 142KiB
  tinystring-ultra.wasm: 23KiB
```

## Troubleshooting

**TinyGo Not Found:**
```
❌ TinyGo is not installed. Building only standard Go binaries.
```
Install TinyGo from: https://tinygo.org/getting-started/install/

**Permission Issues (Linux/macOS/WSL):**
If you encounter permission errors when trying to run the shell scripts, make them executable:
```bash
chmod +x *.sh
```

**Build Failures:**
- Ensure you're in the `benchmark/` directory
- Verify TinyString library is available in the parent directory

