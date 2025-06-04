# TinyString Benchmark Suite

Automated benchmark tools to measure and compare performance between standard Go libraries and TinyString implementations.

## Quick Usage 🚀

```bash
# Run complete benchmark (recommended)
./build-and-measure.sh

# Clean generated files
./clean-all.sh

# Update README with existing data only
./update-readme.sh
```

## What Gets Measured 📊

1. **Binary Size Comparison**: Native + WebAssembly builds with multiple optimization levels
2. **Memory Allocation**: Bytes/op, allocations/op, execution time for 3 benchmark categories:
   - **String Processing**: Case conversion, text manipulation
   - **Number Processing**: Numeric formatting, conversion operations  
   - **Mixed Operations**: Combined string + numeric operations

## Current Performance Status

**Target**: Achieve memory usage close to standard library while maintaining binary size benefits

**Latest Results** (Run `./build-and-measure.sh` to update):
- ✅ **Binary Size**: TinyString is 20-50% smaller than stdlib for WebAssembly
- ⚠️ **Memory Usage**: Number Processing uses 1000% more memory (needs optimization)

## Requirements

- **Go 1.21+**
- **TinyGo** (optional, but recommended for full WebAssembly testing)

## Directory Structure

```
benchmark/
├── analyzer.go               # Main analysis program
├── common.go                # Shared utilities
├── reporter.go              # README update logic
├── build-and-measure.sh     # Main benchmark script
├── memory-benchmark.sh      # Memory benchmarks
├── clean-all.sh            # Cleanup script  
├── update-readme.sh        # README updater
├── run-all-benchmarks.sh   # Run all benchmarks
├── bench-binary-size/      # Binary size test programs
│   ├── standard-lib/       # Standard library example
│   └── tinystring-lib/     # TinyString example
└── bench-memory-alloc/     # Memory allocation benchmarks
    ├── standard/           # Standard library memory tests
    ├── tinystring/         # TinyString memory tests
    └── pointer-comparison/ # Pointer optimization tests
```

## What the Scripts Do

### `build-and-measure.sh`
- Compiles examples with multiple TinyGo optimization levels (-ultra, -speed, -debug)
- Measures binary sizes for both native and WebAssembly targets
- Runs memory allocation benchmarks
- Updates the main project README.md with current results

### `memory-benchmark.sh`
- Executes Go memory benchmarks for standard library vs TinyString
- Measures bytes/operation, allocations/operation, and execution time
- Generates memory comparison data for README

### `clean-all.sh`
- Removes all generated binaries (.exe, .wasm files)
- Cleans up temporary analysis files

### `update-readme.sh`
- Updates only the benchmark sections in the main README
- Uses existing binary files (doesn't rebuild)

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
```bash
chmod +x *.sh
```

**Build Failures:**
- Ensure you're in the `benchmark/` directory
- Verify TinyString library is available in the parent directory

