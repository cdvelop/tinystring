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
├── analyzer.go               # Main analysis program for benchmark results.
├── common.go                 # Shared utilities used by benchmark scripts and tools.
├── reporter.go               # Logic for updating the README.md with benchmark results.
├── MEMORY_REDUCTION.md       # Detailed guide for memory optimization techniques in TinyGo.
├── build-and-measure.sh      # Main script: builds, measures binary sizes, runs memory benchmarks, and updates README.md.
├── memory-benchmark.sh       # Script to run only memory allocation benchmarks.
├── clean-all.sh              # Script to clean all generated binaries and temporary files.
├── update-readme.sh          # Script to update benchmark sections in README.md using existing data (does not run new benchmarks).
├── run-all-benchmarks.sh     # Script to execute all benchmark tests (binary size and memory) without updating the README.md.
├── bench-binary-size/        # Contains Go programs for binary size testing.
│   ├── standard-lib/         # Example project using standard Go library.
│   └── tinystring-lib/       # Example project using TinyString library.
└── bench-memory-alloc/       # Contains Go programs for memory allocation benchmarks.
    ├── standard/             # Memory benchmark tests for standard Go library.
    ├── tinystring/           # Memory benchmark tests for TinyString library.
    └── pointer-comparison/   # Specific tests for pointer optimization in TinyString.
```

## What the Scripts Do

This section provides a clear explanation of each script's function, expected behavior, and typical use case.

### `build-and-measure.sh`
*   **Purpose**: This is the main, comprehensive benchmark script. It orchestrates the entire benchmarking process.
*   **Actions**:
    1.  Compiles example applications (both standard library and TinyString versions) using various TinyGo optimization levels (e.g., -ultra, -speed, -debug) if TinyGo is available.
    2.  Measures the resulting binary sizes for both native and WebAssembly targets.
    3.  Executes memory allocation benchmarks (delegating to `go test` with `-benchmem`).
    4.  Calls `reporter.go` to update the main project's `README.md` (this file) with the latest binary size and memory allocation results.
*   **Output**: Updates `README.md` with new data tables and summaries. Generates compiled binaries.
*   **Use Case**: Run this script to get a full performance overview and update the documentation with the latest figures.

### `memory-benchmark.sh`
*   **Purpose**: Executes only the memory allocation benchmarks.
*   **Actions**:
    1.  Navigates to the relevant benchmark directories (`bench-memory-alloc/standard` and `bench-memory-alloc/tinystring`).
    2.  Runs `go test -bench=. -benchmem` to perform memory benchmarks.
    3.  The script itself does not directly update the README; results are typically fed into `reporter.go` by `build-and-measure.sh` or can be analyzed manually.
*   **Output**: Prints benchmark results (Bytes/op, Allocs/op, ns/op) to standard output.
*   **Use Case**: Use this script when you specifically want to test memory performance without rebuilding binaries or updating the README. Useful for focused optimization efforts.

### `clean-all.sh`
*   **Purpose**: Cleans up files generated by the benchmark and build processes.
*   **Actions**:
    1.  Removes all compiled binaries (e.g., `.exe`, `.wasm` files) from benchmark directories.
    2.  Deletes temporary analysis files or other artifacts created during benchmarking.
*   **Output**: A cleaner workspace, free of generated files.
*   **Use Case**: Run this before a fresh benchmark run or to free up disk space.

### `update-readme.sh`
*   **Purpose**: Updates the benchmark sections in this `README.md` file using previously generated benchmark data.
*   **Actions**:
    1.  Reads existing data files (if any, typically produced by `analyzer.go` after benchmarks have been run).
    2.  Calls `reporter.go` to re-format and insert this data into the `README.md`.
*   **Important**: This script **does not** re-run any benchmarks or re-compile any code. It only updates the documentation based on the last available data.
*   **Output**: Modifies `README.md` if existing data is found.
*   **Use Case**: Use this if you have new benchmark data processed by `analyzer.go` and only want to update the README without running the full `build-and-measure.sh` script again.

### `run-all-benchmarks.sh`
*   **Purpose**: Executes all available benchmark tests (both binary size and memory allocation) but does not automatically update the `README.md`.
*   **Actions**:
    1.  Performs binary builds and size measurements (similar to `build-and-measure.sh` but without the README update step).
    2.  Runs all memory allocation benchmarks (similar to `memory-benchmark.sh`).
*   **Output**: Prints benchmark results to standard output and generates compiled binaries and potentially raw data files.
*   **Use Case**: Useful for gathering all raw benchmark data for analysis without immediately changing the `README.md`. The results can then be used by `update-readme.sh` or analyzed separately.

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

