# Memory Reduction Techniques for TinyGo WebAssembly

## TinyString as Standard Library Replacement

**TinyString is designed to completely replace Go's standard libraries** (`fmt`, `strings`, `strconv`) for string and numeric manipulation in TinyGo WebAssembly applications. Since **string and number processing are the most common programming tasks**, this document serves as a **comprehensive guide for best practices** when transitioning from standard libraries to TinyString's optimized implementations.

### Why Replace Standard Libraries?

String and numeric operations are fundamental to most applications, but Go's standard libraries (`fmt.Sprintf`, `strings.Builder`, `strconv.Itoa`) introduce significant overhead in TinyGo:

- **Binary bloat**: Standard libraries add unnecessary code
- **Hidden allocations**: Unpredictable memory usage  
- **Performance penalties**: Reflection and parsing overhead
- **WebAssembly inefficiency**: Poor optimization for WASM targets

**TinyString eliminates these issues** by providing manual implementations that are specifically optimized for TinyGo and WebAssembly environments.

## TinyGo Memory Model

TinyGo configures WebAssembly **linear memory** with minimal initial size (2 pages = 128 KiB) and grows dynamically via `memory.grow` calls. Each Go allocation can expand WASM memory until resources are exhausted.

**GC Modes:**
- `-gc=conservative` (default): Mark-sweep GC with unpredictable performance
- `-gc=leaking`: Only allocates, never frees - ultra-fast but memory grows indefinitely
- Use `-print-allocs=.` to identify heap allocations at compile time

## Operations That Cause Allocations

**Common Standard Library Problems vs TinyString Solutions:**

| **Standard Library** | **Issue** | **TinyString Solution** |
|---------------------|-----------|-------------------------|
| `fmt.Sprintf("Value: %d", num)` | Parsing + reflection overhead | `tinystring.Format("Value: %d", num)` - Direct implementation |
| `strings.Builder` concatenation | Still uses standard library internally | `tinystring.Convert().Join()` - Manual implementation |
| `strconv.Itoa(num)` | Standard library dependency | `tinystring.Convert(num).String()` - Zero dependencies |
| `string(bytes)` / `[]byte(string)` | Duplicates data | Use `unsafe.String()` / `unsafe.SliceData()` |
| String concatenation (`s1 + s2`) | Creates new string | Use TinyString chaining methods |

## TinyString Best Practices for Standard Library Replacement

### 1. Replace `fmt.Sprintf` with TinyString.Format
```go
// ❌ Standard Library (bloated, slow)
result := fmt.Sprintf("User: %s, Age: %d", name, age)

// ✅ TinyString (optimized, zero dependencies)
result := tinystring.Format("User: %s, Age: %d", name, age)
```

### 2. Replace `strconv` with TinyString Conversion
```go
// ❌ Standard Library
numStr := strconv.Itoa(42)
floatStr := strconv.FormatFloat(3.14, 'f', 2, 64)

// ✅ TinyString  
numStr := tinystring.Convert(42).String()
floatStr := tinystring.Convert(3.14).RoundDecimals(2).String()
```

### 3. Replace `strings` Operations with TinyString Chaining
```go
// ❌ Standard Library (multiple allocations)
result := strings.ToUpper(strings.TrimSpace(input))
parts := strings.Split(input, ",")
joined := strings.Join(parts, "|")

// ✅ TinyString (single chain, fewer allocations)
result := tinystring.Convert(input).Trim().ToUpper().String()
joined := tinystring.Convert(input).Split(",").Join("|")
```

### 4. In-Place String Modification (Unique to TinyString)
```go
// ✅ TinyString exclusive feature - modify original string
text := "hello world"
tinystring.Convert(&text).ToUpper().RemoveTilde().Apply()
// text is now modified directly: "HELLO WORLD"
```

## TinyGo `unsafe.Pointer` for WebAssembly

For WebAssembly (WASM) targets, especially when aiming for minimal binary size and high performance, TinyString leverages TinyGo's `unsafe.Pointer` capabilities.

**`unsafe.Pointer` Support**: TinyGo 0.37.0 (compatible with Go 1.24) fully supports `unsafe.Pointer`, `unsafe.String()`, and `unsafe.SliceData()`. These are standard Go features (added in Go 1.20) and are recommended for conversions between pointers and slices/strings without data copying, replacing older methods like direct manipulation of `reflect.StringHeader` or `reflect.SliceHeader`.

**Memory Implications**: Using these `unsafe` functions allows for more efficient memory operations, which is crucial for performance-sensitive code and reducing the overhead of string manipulations in WASM environments. TinyGo compiles WASM code using these functions without issues when targeting `syscall/js` for JavaScript interoperability.

**Best Practices**:
- Prefer using `unsafe.String()` and `unsafe.SliceData()` for zero-copy conversions
- Avoid direct manipulation of `reflect.Header` types (TinyGo has internal differences)
- Use `unsafe` conversions to eliminate data copying in string operations

```go
// ✅ Zero-copy string from byte slice
data := []byte("hello world")
str := unsafe.String(&data[0], len(data))

// ✅ Zero-copy byte slice from string  
s := "hello world"
bytes := unsafe.Slice(unsafe.StringData(s), len(s))

// ✅ Buffer to string conversion (used internally by strings.Builder)
buf := []byte("processed data")
result := *(*string)(unsafe.Pointer(&buf))
```

## TinyGo Optimization Flags

| **Optimization Goal** | **Flags** | **Description** | **Notes** |
|----------------------|-----------|-----------------|-----------|
| **Minimal Size** | `-no-debug -panic=trap -scheduler=none -gc=leaking` | Remove debug info, disable goroutines, disable GC | Best for short-lived programs |
| **Code Size** | `-opt=z` (default), `-opt=s` | Optimize for size | Default TinyGo behavior |
| **Speed** | `-opt=2 -gc=leaking -scheduler=none -panic=trap` | Optimize for performance | Memory never freed with `-gc=leaking` |
| **Debugging** | `-opt=1` or `-opt=0` | Reduce optimizations for easier debugging | `-opt=0` may break some programs |
| **WebAssembly** | `-no-debug` | Remove debug symbols (⅔ size reduction) | Use `strip` utility on desktop systems |

**Tips:**
- Avoid large packages like `fmt` - use `println()` instead
- Example extreme optimization: 93K → 1.6K with all size flags
- `-gc=leaking` only for very short-lived programs

## Analysis Tools

```bash
# Identify allocations at compile time
tinygo test -print-allocs=.

# Memory profiling
go test -benchmem -bench=.
```

## Key Recommendations for TinyString Migration

1. **Replace `fmt.Sprintf`** with `tinystring.Format()` for all string formatting
2. **Replace `strconv` conversions** with `tinystring.Convert()` methods
3. **Replace `strings` operations** with TinyString's chainable methods
4. **Use pointer modification** with `Apply()` to avoid allocations entirely
5. **Leverage TinyString's numeric handling** for formatting and rounding
6. **Use `-gc=leaking`** for short-lived WebAssembly modules
7. **Always use `-no-debug`** in production to reduce binary size
8. **Combine operations in single chains** to minimize intermediate allocations

**The Result**: Complete elimination of standard library dependencies while maintaining familiar Go syntax and achieving optimal memory usage for WebAssembly string and numeric manipulation - the most common programming tasks.
