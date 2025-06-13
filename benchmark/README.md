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

2.  **Memory Allocation**: Measures Bytes/op, Allocations/op, and execution time (ns/op) for benchmark categories:
    *   **String Processing**: Basic string operations (case conversion, manipulation)
    *   **Number Processing**: Numeric formatting and conversion operations
    *   **Mixed Operations**: Combined string and numeric operations
    *   **JSON Operations**: JSON marshaling and unmarshaling operations
        - Complex nested structures
        - Batch processing with different sizes
        - Error handling scenarios

## JSON Benchmarks Overview

The JSON benchmarking system is located in `bench-memory-alloc/json-comparison/` and consists of:

- `data.go`: Contains test data structures and generation logic
- `main_test.go`: Core benchmarking functions for marshaling/unmarshaling
- `errors_test.go`: Error case handling tests
- `README.md`: Detailed documentation for JSON benchmark suite

For specific implementation details and examples, refer to the documentation in the JSON comparison directory.

## Current Performance Status

**Target**: Achieve memory usage close to standard library while maintaining binary size benefits.

**Latest Results** (Run `./build-and-measure.sh` to update):
- ✅ **Binary Size**: TinyString is 20-50% smaller than stdlib for WebAssembly.
- ⚠️ **Memory Usage**: Number Processing uses 1000% more memory (needs optimization).
- 📊 **JSON Performance**:
  - ✅ Marshal: 25-30% better performance and memory usage
  - ⚠️ Unmarshal: 50% higher memory usage and slower processing
  - 🎯 Error Handling: Mixed results (better for Marshal, worse for Unmarshal)

📋 **Memory Optimization Guide**: See [`MEMORY_REDUCTION.md`](./MEMORY_REDUCTION.md) for comprehensive techniques and best practices to replace Go standard libraries with TinyString's optimized implementations. Essential reading for efficient string and numeric processing in TinyGo WebAssembly applications.

## Requirements

- **Go 1.21+**
- **TinyGo** (optional, but recommended for full WebAssembly testing and to achieve smallest binary sizes).

## Directory Structure

```
benchmark/
├── analyzer.go               # Main analysis program for benchmark results.
├── common.go                # Shared utilities used by benchmark scripts and tools.
├── reporter.go              # Logic for updating the README.md with benchmark results.
├── MEMORY_REDUCTION.md      # Detailed guide for memory optimization techniques in TinyGo.
├── build-and-measure.sh     # Main comprehensive script: compiles apps with TinyGo optimizations,
│                            # measures binary sizes (native + WebAssembly), runs memory benchmarks,
│                            # and updates README.md with latest results. Use for full performance overview.
├── memory-benchmark.sh      # Executes only memory allocation benchmarks without building binaries.
│                            # Runs 'go test -bench=. -benchmem' in standard/tinystring directories.
│                            # Useful for focused memory optimization efforts.
├── clean-all.sh            # Removes compiled binaries (.exe, .wasm) and temporary analysis files.
│                            # Run before fresh benchmark runs or to free disk space.
├── update-readme.sh        # Updates README.md benchmark sections using existing data only.
│                            # Does NOT re-run benchmarks or recompile code. Only reformats/inserts
│                            # previously generated data into documentation.
├── run-all-benchmarks.sh   # Executes all benchmark tests (binary size + memory allocation) but
│                            # does NOT update README.md. Generates raw data for manual analysis.
├── bench-binary-size/      # Contains Go programs for binary size testing.
│   ├── standard-lib/       # Example project using standard Go library.
│   └── tinystring-lib/     # Example project using TinyString library.
└── bench-memory-alloc/     # Contains Go programs for memory allocation benchmarks.
    ├── standard/           # Memory benchmark tests for standard Go library.
    ├── tinystring/        # Memory benchmark tests for TinyString library.
    ├── pointer-comparison/ # Specific tests for pointer optimization in TinyString.
    └── json-comparison/    # JSON functionality benchmarks.
        ├── data.go        # Test data structures and JSON generators.
        ├── main_test.go   # Core marshal/unmarshal benchmarks.
        ├── errors_test.go # Error handling benchmarks.
        └── README.md      # JSON benchmark documentation.
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

📊 Latest Results: See generated benchmark reports in respective test directories
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

