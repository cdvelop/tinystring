# TinyString Benchmark Scripts

Automated benchmark tools to measure and compare binary sizes and memory allocations between standard Go libraries and TinyString implementations.

## Quick Start

### Run Complete Benchmark

```bash
# Run full benchmark (binary size + memory allocation)
./build-and-measure.sh
```

This script will:
1. Build binaries with multiple TinyGo optimization levels
2. Measure all binary sizes 
3. Run memory allocation benchmarks
4. Update the main README.md with results

### Individual Scripts

```bash
# Clean all generated files
./clean-all.sh

# Update only README with existing data
./update-readme.sh

# Run specific benchmarks
./run-all-benchmarks.sh
```

## Requirements

- **Go 1.21+**: For building native binaries
- **TinyGo** (optional): For WebAssembly compilation
  - Install from: https://tinygo.org/getting-started/install/
  - If not installed, only native binaries will be measured

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

