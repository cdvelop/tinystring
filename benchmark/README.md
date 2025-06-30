# TinyString Benchmark Suite

Automated benchmark tools to measure and compare performance between standard Go libraries and TinyString implementations.



## Binary Size Comparison

[Standard Library Example](bench-binary-size/standard-lib/main.go) | [TinyString Example](bench-binary-size/tinystring-lib/main.go)

<!-- This table is automatically generated from build-and-measure.sh -->
*Last updated: 2025-06-29 20:31:51*

| Build Type | Parameters | Standard Library<br/>`go build` | TinyString<br/>`tinygo build` | Size Reduction | Performance |
|------------|------------|------------------|------------|----------------|-------------|
| 🖥️ **Default Native** | `-ldflags="-s -w"` | 1.3 MB | 1.1 MB | **-164.0 KB** | ➖ **12.5%** |
| 🌐 **Default WASM** | `(default -opt=z)` | 580.8 KB | 252.4 KB | **-328.4 KB** | ✅ **56.5%** |
| 🌐 **Ultra WASM** | `-no-debug -panic=trap -scheduler=none -gc=leaking -target wasm` | 141.3 KB | 34.3 KB | **-107.0 KB** | 🏆 **75.8%** |
| 🌐 **Speed WASM** | `-opt=2 -target wasm` | 827.0 KB | 362.5 KB | **-464.5 KB** | ✅ **56.2%** |
| 🌐 **Debug WASM** | `-opt=0 -target wasm` | 1.8 MB | 800.8 KB | **-1.0 MB** | ✅ **56.3%** |

### 🎯 Performance Summary

- 🏆 **Peak Reduction: 75.8%** (Best optimization)
- ✅ **Average WebAssembly Reduction: 61.2%**
- ✅ **Average Native Reduction: 12.5%**
- 📦 **Total Size Savings: 2.0 MB across all builds**

#### Performance Legend
- ❌ Poor (<5% reduction)
- ➖ Fair (5-15% reduction)
- ✅ Good (15-70% reduction)
- 🏆 Outstanding (>70% reduction)


## Memory Usage Comparison

[Standard Library Example](bench-memory-alloc/standard) | [TinyString Example](bench-memory-alloc/tinystring)

<!-- This table is automatically generated from memory-benchmark.sh -->
*Last updated: 2025-06-29 20:32:08*

Performance benchmarks comparing memory allocation patterns between standard Go library and TinyString:

| 🧪 **Benchmark Category** | 📚 **Library** | 💾 **Memory/Op** | 🔢 **Allocs/Op** | ⏱️ **Time/Op** | 📈 **Memory Trend** | 🎯 **Alloc Trend** | 🏆 **Performance** |
|----------------------------|----------------|-------------------|-------------------|-----------------|---------------------|---------------------|--------------------|
| 📝 **String Processing** | 📊 Standard | `752 B / 530.773 OP` | `32` | `2.1μs` | - | - | - |
| | 🚀 TinyString | `1.5 KB / 110.544 OP` | `49` | `10.6μs` | ❌ **100.0% more** | ❌ **53.1% more** | ❌ **Poor** |
| 🔢 **Number Processing** | 📊 Standard | `912 B / 481.078 OP` | `42` | `2.4μs` | - | - | - |
| | 🚀 TinyString | `320 B / 596.006 OP` | `17` | `2.0μs` | 🏆 **64.9% less** | 🏆 **59.5% less** | 🏆 **Excellent** |
| 🔄 **Mixed Operations** | 📊 Standard | `400 B / 758.782 OP` | `22` | `1.5μs` | - | - | - |
| | 🚀 TinyString | `432 B / 341.202 OP` | `20` | `3.4μs` | ⚠️ **8.0% more** | ✅ **9.1% less** | ➖ **Fair** |

### 🎯 Performance Summary

- 💾 **Memory Efficiency**: ➖ **Fair** (Acceptable overhead) (14.4% average change)
- 🔢 **Allocation Efficiency**: ✅ **Good** (Allocation efficient) (-5.2% average change)
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




