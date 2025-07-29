# TinyString Benchmark Suite

Automated benchmark tools to measure and compare performance between standard Go libraries and TinyString implementations.



## Binary Size Comparison

[Standard Library Example](bench-binary-size/standard-lib/main.go) | [TinyString Example](bench-binary-size/tinystring-lib/main.go)

<!-- This table is automatically generated from build-and-measure.sh -->
*Last updated: 2025-07-29 15:11:57*

| Build Type | Parameters | Standard Library<br/>`go build` | TinyString<br/>`tinygo build` | Size Reduction | Performance |
|------------|------------|------------------|------------|----------------|-------------|
| 🖥️ **Default Native** | `-ldflags="-s -w"` | 1.4 MB | 1.2 MB | **-160.0 KB** | ➖ **11.2%** |
| 🌐 **Default WASM** | `(default -opt=z)` | 590.2 KB | 244.6 KB | **-345.5 KB** | ✅ **58.5%** |
| 🌐 **Ultra WASM** | `-no-debug -panic=trap -scheduler=none -gc=leaking -target wasm` | 141.2 KB | 31.8 KB | **-109.3 KB** | 🏆 **77.5%** |
| 🌐 **Speed WASM** | `-opt=2 -target wasm` | 837.2 KB | 350.1 KB | **-487.1 KB** | ✅ **58.2%** |
| 🌐 **Debug WASM** | `-opt=0 -target wasm` | 1.8 MB | 825.4 KB | **-994.4 KB** | ✅ **54.6%** |

### 🎯 Performance Summary

- 🏆 **Peak Reduction: 77.5%** (Best optimization)
- ✅ **Average WebAssembly Reduction: 62.2%**
- ✅ **Average Native Reduction: 11.2%**
- 📦 **Total Size Savings: 2.0 MB across all builds**

#### Performance Legend
- ❌ Poor (<5% reduction)
- ➖ Fair (5-15% reduction)
- ✅ Good (15-70% reduction)
- 🏆 Outstanding (>70% reduction)


## Memory Usage Comparison

[Standard Library Example](bench-memory-alloc/standard) | [TinyString Example](bench-memory-alloc/tinystring)

<!-- This table is automatically generated from memory-benchmark.sh -->
*Last updated: 2025-07-29 15:12:12*

Performance benchmarks comparing memory allocation patterns between standard Go library and TinyString:

| 🧪 **Benchmark Category** | 📚 **Library** | 💾 **Memory/Op** | 🔢 **Allocs/Op** | ⏱️ **Time/Op** | 📈 **Memory Trend** | 🎯 **Alloc Trend** | 🏆 **Performance** |
|----------------------------|----------------|-------------------|-------------------|-----------------|---------------------|---------------------|--------------------|
| 📝 **String Processing** | 📊 Standard | `808 B / 609.450 OP` | `32` | `2.3μs` | - | - | - |
| | 🚀 TinyString | `464 B / 225.711 OP` | `17` | `5.2μs` | 🏆 **42.6% less** | 🏆 **46.9% less** | 🏆 **Excellent** |
| 🔢 **Number Processing** | 📊 Standard | `912 B / 535.048 OP` | `42` | `2.4μs` | - | - | - |
| | 🚀 TinyString | `320 B / 526.324 OP` | `17` | `1.9μs` | 🏆 **64.9% less** | 🏆 **59.5% less** | 🏆 **Excellent** |
| 🔄 **Mixed Operations** | 📊 Standard | `416 B / 722.950 OP` | `22` | `1.4μs` | - | - | - |
| | 🚀 TinyString | `192 B / 468.020 OP` | `12` | `2.5μs` | 🏆 **53.8% less** | 🏆 **45.5% less** | 🏆 **Excellent** |

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




