# TinyString Benchmark System Usage

This enhanced benchmark system automates the measurement and reporting of binary sizes when comparing standard Go libraries vs TinyString library usage across **multiple optimization levels**.

## Quick Start

### Run Complete Benchmark
```bash
./benchmark/scripts/build-and-measure.sh
```
This will:
1. Build both standard library and TinyString examples
2. Compile with **4 different optimization levels** for WebAssembly
3. Measure actual binary sizes for all variants
4. Generate a comprehensive optimization comparison table
5. Update the main README.md with real data

### Update README Only
```bash
./benchmark/scripts/update-readme.sh
```
Use this when binaries already exist and you only want to update the README.

### Clean Generated Files
```bash
./benchmark/scripts/clean.sh
```
Removes all generated binaries and WebAssembly files while preserving source code.

## Optimization Levels

Based on the [TinyGo optimization guide](https://tinygo.org/docs/guides/optimizing-binaries/), the benchmark tests 4 different optimization configurations:

### 1. Default (`-opt=z`)
- **Purpose**: Default TinyGo optimization for size
- **Usage**: `tinygo build -target wasm main.go`
- **Results**: Standard: 580.6KB → TinyString: 229.5KB (**60.5% reduction**)

### 2. Ultra Size (`-no-debug -panic=trap -scheduler=none -gc=leaking`)
- **Purpose**: Maximum size optimization
- **Usage**: `tinygo build -no-debug -panic=trap -scheduler=none -gc=leaking -target wasm main.go`
- **Results**: Standard: 141.1KB → TinyString: 22.2KB (**84.3% reduction**)
- **Best for**: Production deployments where size is critical

### 3. Speed (`-opt=2`)
- **Purpose**: Optimized for speed over size
- **Usage**: `tinygo build -opt=2 -target wasm main.go`
- **Results**: Standard: 815.9KB → TinyString: 321.3KB (**60.6% reduction**)
- **Best for**: Performance-critical applications

### 4. Debug (`-opt=0`)
- **Purpose**: No optimization, best for debugging
- **Usage**: `tinygo build -opt=0 -target wasm main.go`
- **Results**: Standard: 1.8MB → TinyString: 664.3KB (**63.8% reduction**)
- **Best for**: Development and debugging

## Generated Binary Files

After running the benchmark, multiple variants are created:

### Standard Library Example
- `standard` - Native binary (~1.3MB)
- `standard.wasm` - Default WebAssembly (~580KB)
- `standard-ultra.wasm` - Ultra-optimized WebAssembly (~141KB)
- `standard-speed.wasm` - Speed-optimized WebAssembly (~816KB)
- `standard-debug.wasm` - Debug WebAssembly (~1.8MB)

### TinyString Example  
- `tinystring` - Native binary (~1.1MB)
- `tinystring.wasm` - Default WebAssembly (~230KB)
- `tinystring-ultra.wasm` - Ultra-optimized WebAssembly (~22KB)
- `tinystring-speed.wasm` - Speed-optimized WebAssembly (~321KB)
- `tinystring-debug.wasm` - Debug WebAssembly (~664KB)

## Key Results Summary

| Optimization | Standard Library | TinyString | Reduction |
|-------------|------------------|------------|-----------|
| **Default** | 580.6KB | 229.5KB | **60.5%** |
| **Ultra Size** | 141.1KB | 22.2KB | **84.3%** |
| **Speed** | 815.9KB | 321.3KB | **60.6%** |
| **Debug** | 1.8MB | 664.3KB | **63.8%** |

## README Integration

The system automatically generates a comprehensive table in the README:

**Before (Manual Estimates):**
```bash
tinygo build -o app-standard.wasm -target wasm main.go  # ~500KB+ WebAssembly
tinygo build -o app-tiny.wasm -target wasm main.go      # ~180KB WebAssembly
```

**After (Real Data with Optimization Table):**
- Basic size comparison with actual measurements
- **Complete optimization comparison table**
- Detailed reduction percentages for each optimization level
- Professional presentation of benchmarking data

## Example Programs

### Standard Library Example (`examples/standard-lib/main.go`)
Demonstrates common string operations using:
- `fmt` package for formatting
- `strings` package for manipulation
- `strconv` package for conversions

### TinyString Example (`examples/tinystring-lib/main.go`)
Equivalent operations using only TinyString:
- Manual implementations without standard library imports
- Same functionality with reduced binary size
- TinyGo compatible without compilation issues

## CI Integration

The benchmark system can be integrated into CI/CD pipelines:

1. **Manual Trigger**: Run `./benchmark/scripts/build-and-measure.sh` when needed
2. **Pre-commit Hook**: Automatically update README before commits
3. **Release Process**: Generate fresh data before releases

## Maintenance

- **Update Examples**: Modify `examples/*/main.go` to test new functionality
- **Add Platforms**: Extend `benchmark.go` to support additional build targets
- **Customize Output**: Modify the README template in `generateBinarySizeSection()`

## Requirements

- Go compiler (for native builds)
- TinyGo compiler (for WebAssembly builds)
- bash shell (for running scripts)
- Local TinyString module (automatic via go.mod replace directive)
