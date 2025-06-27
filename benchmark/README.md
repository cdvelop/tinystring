# TinyString Benchmark Suite

Automated benchmark tools to measure and compare performance between standard Go libraries and TinyString implementations.



## Binary Size Comparison

[Standard Library Example](benchmark/bench-binary-size/standard-lib/main.go) | [TinyString Example](benchmark/bench-binary-size/tinystring-lib/main.go)

<!-- This table is automatically generated from build-and-measure.sh -->
*Last updated: 2025-06-26 22:25:49*

| Build Type | Parameters | Standard Library<br/>`go build` | TinyString<br/>`tinygo build` | Size Reduction | Performance |
|------------|------------|------------------|------------|----------------|-------------|
| 🖥️ **Default Native** | `-ldflags="-s -w"` |  |  | **-** | ➖ **14.0%** |
| 🌐 **Default WASM** | `(default -opt=z)` |  |  | **-** | ✅ **61.0%** |
| 🌐 **Ultra WASM** | `-no-debug -panic=trap -scheduler=none -gc=leaking -target wasm` |  |  | **-** | 🏆 **81.8%** |
| 🌐 **Speed WASM** | `-opt=2 -target wasm` |  |  | **-** | ✅ **61.8%** |
| 🌐 **Debug WASM** | `-opt=0 -target wasm` |  |  | **-** | ✅ **61.4%** |

### 🎯 Performance Summary

- 🏆 **Peak Reduction: 81.8%** (Best optimization)
- ✅ **Average WebAssembly Reduction: 66.5%**
- ✅ **Average Native Reduction: 14.0%**
- 📦 **Total Size Savings:  across all builds**

#### Performance Legend
- ❌ Poor (<5% reduction)
- ➖ Fair (5-15% reduction)
- ✅ Good (15-70% reduction)
- 🏆 Outstanding (>70% reduction)


## Memory Usage Comparison

[Standard Library Example](benchmark/bench-memory-alloc/standard) | [TinyString Example](benchmark/bench-memory-alloc/tinystring)

<!-- This table is automatically generated from memory-benchmark.sh -->
*Last updated: 2025-06-26 22:26:07*

Performance benchmarks comparing memory allocation patterns between standard Go library and TinyString:

| 🧪 **Benchmark Category** | 📚 **Library** | 💾 **Memory/Op** | 🔢 **Allocs/Op** | ⏱️ **Time/Op** | 📈 **Memory Trend** | 🎯 **Alloc Trend** | 🏆 **Performance** |
|----------------------------|----------------|-------------------|-------------------|-----------------|---------------------|---------------------|--------------------|
| 📝 **String Processing** | 📊 Standard | ` / 375,354 OP` | `48` | `3.0μs` | - | - | - |
| | 🚀 TinyString | ` / 89,876 OP` | `135` | `13.8μs` | ❌ **165.5% more** | ❌ **181.3% more** | ❌ **Poor** |
| 🔢 **Number Processing** | 📊 Standard | `912 B / 498,667 OP` | `42` | `2.5μs` | - | - | - |
| | 🚀 TinyString | `320 B / 1,000,000 OP` | `17` | `2.0μs` | 🏆 **64.9% less** | 🏆 **59.5% less** | 🏆 **Excellent** |
| 🔄 **Mixed Operations** | 📊 Standard | `512 B / 677,230 OP` | `26` | `1.7μs` | - | - | - |
| | 🚀 TinyString | `816 B / 292,696 OP` | `40` | `4.1μs` | ❌ **59.4% more** | ❌ **53.8% more** | ❌ **Poor** |

### 🎯 Performance Summary

- 💾 **Memory Efficiency**: ❌ **Poor** (Significant overhead) (53.3% average change)
- 🔢 **Allocation Efficiency**: ❌ **Poor** (Excessive allocations) (58.5% average change)
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




