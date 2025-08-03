# TinyString Benchmark Suite

Automated benchmark tools to measure and compare performance between standard Go libraries and TinyString implementations.



## Binary Size Comparison

[Standard Library Example](bench-binary-size/standard-lib/main.go) | [TinyString Example](bench-binary-size/tinystring-lib/main.go)

<!-- This table is automatically generated from build-and-measure.sh -->
*Last updated: 2025-08-02 22:39:11*

| Build Type | Parameters | Standard Library<br/>`go build` | TinyString<br/>`tinygo build` | Size Reduction | Performance |
|------------|------------|------------------|------------|----------------|-------------|
| 🖥️ **Default Native** | `-ldflags="-s -w"` | 1.4 MB | 1.2 MB | **-152.0 KB** | ➖ **10.6%** |
| 🌐 **Default WASM** | `(default -opt=z)` | 590.2 KB | 257.0 KB | **-333.2 KB** | ✅ **56.5%** |
| 🌐 **Ultra WASM** | `-no-debug -panic=trap -scheduler=none -gc=leaking -target wasm` | 141.2 KB | 32.9 KB | **-108.2 KB** | 🏆 **76.7%** |
| 🌐 **Speed WASM** | `-opt=2 -target wasm` | 837.2 KB | 366.7 KB | **-470.5 KB** | ✅ **56.2%** |
| 🌐 **Debug WASM** | `-opt=0 -target wasm` | 1.8 MB | 868.0 KB | **-951.8 KB** | ✅ **52.3%** |

### 🎯 Performance Summary

- 🏆 **Peak Reduction: 76.7%** (Best optimization)
- ✅ **Average WebAssembly Reduction: 60.4%**
- ✅ **Average Native Reduction: 10.6%**
- 📦 **Total Size Savings: 2.0 MB across all builds**

#### Performance Legend
- ❌ Poor (<5% reduction)
- ➖ Fair (5-15% reduction)
- ✅ Good (15-70% reduction)
- 🏆 Outstanding (>70% reduction)


## Memory Usage Comparison

[Standard Library Example](bench-memory-alloc/standard) | [TinyString Example](bench-memory-alloc/tinystring)

<!-- This table is automatically generated from memory-benchmark.sh -->
*Last updated: 2025-08-02 22:39:27*

Performance benchmarks comparing memory allocation patterns between standard Go library and TinyString:

| 🧪 **Benchmark Category** | 📚 **Library** | 💾 **Memory/Op** | 🔢 **Allocs/Op** | ⏱️ **Time/Op** | 📈 **Memory Trend** | 🎯 **Alloc Trend** | 🏆 **Performance** |
|----------------------------|----------------|-------------------|-------------------|-----------------|---------------------|---------------------|--------------------|
| 📝 **String Processing** | 📊 Standard | `808 B / 535.737 OP` | `32` | `2.0μs` | - | - | - |
| | 🚀 TinyString | `464 B / 203.563 OP` | `17` | `5.9μs` | 🏆 **42.6% less** | 🏆 **46.9% less** | 🏆 **Excellent** |
| 🔢 **Number Processing** | 📊 Standard | `912 B / 511.083 OP` | `42` | `2.3μs` | - | - | - |
| | 🚀 TinyString | `320 B / 611.653 OP` | `17` | `1.9μs` | 🏆 **64.9% less** | 🏆 **59.5% less** | 🏆 **Excellent** |
| 🔄 **Mixed Operations** | 📊 Standard | `416 B / 808.946 OP` | `22` | `1.4μs` | - | - | - |
| | 🚀 TinyString | `192 B / 427.737 OP` | `12` | `2.7μs` | 🏆 **53.8% less** | 🏆 **45.5% less** | 🏆 **Excellent** |

### 🎯 Performance Summary

- 💾 **Memory Efficiency**: 🏆 **Excellent** (Lower memory usage) (-53.8% average change)
- 🔢 **Allocation Efficiency**: 🏆 **Excellent** (Fewer allocations) (-50.6% average change)
- 📊 **Benchmarks Analyzed**: 3 categories
- 🎯 **Optimization Focus**: Binary size reduction vs runtime efficiency

### ⚖️ Trade-offs Analysis

The benchmarks reveal important trade-offs between **binary size** and **runtime performance**:

#### 📦 **Binary Size Benefits** ✅
- 🏆 **16-84% smaller** compiled binaries
- 🌐 **Superior WebAssembly** compression ratios
- 🚀 **Faster deployment** and distribution
- 💾 **Lower storage** requirements

#### 🧠 **Runtime Memory Considerations** ⚠️
- 📈 **Higher allocation overhead** during execution
- 🗑️ **Increased GC pressure** due to allocation patterns
- ⚡ **Trade-off optimizes** for distribution size over runtime efficiency
- 🔄 **Different optimization strategy** than standard library

#### 🎯 **Optimization Recommendations**
| 🎯 **Use Case** | 💡 **Recommendation** | 🔧 **Best For** |
|-----------------|------------------------|------------------|
| 🌐 WebAssembly Apps | ✅ **TinyString** | Size-critical web deployment |
| 📱 Embedded Systems | ✅ **TinyString** | Resource-constrained devices |
| ☁️ Edge Computing | ✅ **TinyString** | Fast startup and deployment |
| 🏢 Memory-Intensive Server | ⚠️ **Standard Library** | High-throughput applications |
| 🔄 High-Frequency Processing | ⚠️ **Standard Library** | Performance-critical workloads |

#### 📊 **Performance Legend**
- 🏆 **Excellent** (Better performance)
- ✅ **Good** (Acceptable trade-off)
- ⚠️ **Caution** (Higher resource usage)
- ❌ **Poor** (Significant overhead)


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




