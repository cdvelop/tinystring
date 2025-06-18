# TinyString Builder Integration - Technical Design Document

## Overview

This document outlines the integration of a high-performance string builder API into TinyString's existing `conv` structure, along with the creation of a new `T` (Translate) function to optimize memory usage in high-demand processes.

## Problem Statement

Current memory hotspots identified in TinyString:
- **Error message construction** (`Err` function) with dictionary translations
- **Format operations** (`Fmt` function) with multiple concatenations  
- **String transformation chains** (ToLower, CamelCase, etc.)
- **Repeated string allocations** in complex operations

## Solution Architecture

### 1. Builder API Integration

#### Design Decision: Extend `conv` with Builder Methods
- **Rationale**: Reuse existing buffer (`buf []byte`) and pooling infrastructure
- **Location**: All builder methods implemented in `memory.go`
- **API**: Short method names for efficiency

#### Builder API Methods
```go
// String builder methods - all return *conv for chaining
func (c *conv) WS(s string) *conv     // WriteString equivalent
func (c *conv) WB(b byte) *conv       // WriteByte equivalent  
func (c *conv) WR(r rune) *conv       // WriteRune equivalent
func (c *conv) RS() *conv             // Reset equivalent

// Build result - reuses existing String() method
func (c *conv) String() string        // Already exists - auto-releases to pool
```

### 2. Translation Function `T`

#### New Function Signature
```go
func T(values ...any) string
```

#### Usage Examples
```go
// Simple translation
msg := T(D.Invalid, D.Format)              // "invalid format"
msg := T(ES, D.Invalid, D.Format)          // "formato inválido" 

// Error function becomes wrapper
func Err(values ...any) *conv {
    return &conv{err: T(values...), vTpe: typeErr}
}
```

### 3. Memory Optimization Strategy

#### Responsibility Separation
- **`RS()` method**: Reset complete conv state (all fields + buffer)
- **`setString()` method**: Convert to string type + full conv cleanup
- **Builder methods**: Focus on efficient string construction
- **Pool integration**: Seamless reuse of existing `convPool`

#### Buffer Reuse Pattern
```go
// High-performance pattern for complex operations
func complexOperation() string {
    c := getConv()
    result := c.WS("prefix").WS(" ").WS("suffix").String() // Auto-releases
    return result
}
```

## Implementation Plan

### Phase 1: Builder Methods Implementation
**File**: `memory.go`

1. **Add Builder API methods**:
   - `WS(s string) *conv` - Append string to buffer
   - `WB(b byte) *conv` - Append byte to buffer  
   - `WR(r rune) *conv` - Append rune to buffer (UTF-8 encoded)
   - `RS() *conv` - Reset complete conv state

2. **Optimize buffer management**:
   - Reuse existing `cap()` (rename from ensureCapacity) and `resetBuffer()`
   - Enhance `setStringFromBuffer()` for builder workflow

### Phase 2: Translation Function `T`
**File**: `translation.go` (new file)

1. **Implement `T` function**:
   - Same signature as `Err` but returns `string`
   - Reuse translation logic from existing `Err`
   - Use builder API internally for efficiency

2. **Refactor `Err` function**:
   - Simplify to wrapper around `T`
   - Maintain backward compatibility

### Phase 3: Integration & Optimization
**Files**: Throughout library

1. **Create translation.go**: New file containing `T` function
2. **Identify high-demand processes** using memory profiling tools
3. **Replace string concatenation** with builder API
4. **Optimize transformation chains** to use single buffer

## Memory Analysis Requirements

### Benchmarking Strategy
Using tools from `ISSUE_MEMORY_TOOLS.md`:

1. **Baseline measurements**:
   ```bash
   go test -bench=BenchmarkStringOperations -benchmem -memprofile=before.prof
   ```

2. **Post-implementation comparison**:
   ```bash
   go test -bench=BenchmarkStringOperations -benchmem -memprofile=after.prof
   benchstat before.txt after.txt
   ```

3. **Heap escape analysis**:
   ```bash
   go build -gcflags="-m" ./...
   ```

### Target Metrics
- **Reduce allocations** in `Err` function by 50%+
- **Eliminate temporary strings** in transformation chains
- **Maintain or improve performance** vs `strings.Builder`
- **Zero additional binary size** impact

## Technical Specifications

### Builder Method Details

#### `WS(s string) *conv` - WriteString
```go
func (c *conv) WS(s string) *conv {
    c.cap(len(c.buf) + len(s))
    c.buf = append(c.buf, s...)
    return c
}
```

#### `WB(b byte) *conv` - WriteByte  
```go
func (c *conv) WB(b byte) *conv {
    c.cap(len(c.buf) + 1)
    c.buf = append(c.buf, b)
    return c
}
```

#### `WR(r rune) *conv` - WriteRune
```go
func (c *conv) WR(r rune) *conv {
    c.buf = addRne2Buf(c.buf, r) // Reuse existing UTF-8 encoder
    return c
}
```

#### `RS() *conv` - Reset Builder
```go
func (c *conv) RS() *conv {
    // Reset all conv fields to default state
    c.stringVal = ""
    c.intVal = 0
    c.uintVal = 0
    c.floatVal = 0
    c.boolVal = false
    c.stringSliceVal = nil
    c.stringPtrVal = nil
    c.vTpe = typeStr
    c.roundDown = false
    c.separator = "_"
    c.tmpStr = ""
    c.lastConvType = typeStr
    c.err = ""
    c.buf = c.buf[:0] // Reset buffer length, keep capacity
    return c
}
```

### Method Naming Conventions

#### Short API Names
- **`WS`**: WriteString - Optimized for frequent use
- **`WB`**: WriteByte - Single byte append
- **`WR`**: WriteRune - UTF-8 rune encoding  
- **`RS`**: Reset - Complete conv reset (all fields + buffer)
- **`cap`**: Capacity management (renamed from `ensureCapacity` for brevity)

#### Implementation Notes
- **`RS()` behavior**: Resets all conv fields to default state, similar to `putConv()` but without returning to pool
- **Complete reset**: Clears all values, buffer, and state for fresh builder usage
- **Capacity optimization**: `cap()` method replaces `ensureCapacity()` with shorter name
- **Memory efficiency**: All methods reuse existing buffer infrastructure

### Integration Points

#### High-Demand Process Optimization
1. **`Err` function**: Use `T` with builder internally
2. **`Fmt` function**: Replace string concatenation with builder
3. **Transformation chains**: Single buffer for multiple operations
4. **Dictionary translations**: Efficient multi-part message construction

#### Pool Integration
- Builder methods work seamlessly with existing `convPool`
- `String()` method handles auto-release (existing behavior)
- No additional pool required

### File Structure Updates

#### New File: `translation.go`
**Purpose**: Dedicated translation functionality
```go
package tinystring

// T creates a translated string with support for multilingual translations
// Same functionality as Err but returns string directly instead of *conv
// eg:
// T(D.Invalid, D.Format) returns "invalid format"
// T(ES, D.Invalid, D.Format) returns "formato inválido"
func T(values ...any) string {
    // Implementation using builder API for efficiency
    // Reuses logic from existing Err function
}
```

#### Updated `error.go`
**Purpose**: Simplified error handling using T function
```go
// Err becomes a simple wrapper around T function
func Err(values ...any) *conv {
    return &conv{err: T(values...), vTpe: typeErr}
}
```

#### Builder Integration
- **T function** uses builder API internally for zero-allocation string construction
- **Separation of concerns**: Translation logic isolated from error handling
- **Reusability**: T function can be used independently for any translation needs

## Testing Strategy

### New Benchmarks
Add to `benchmark_strings_test.go`:

```go
func BenchmarkBuilderOperations(b *testing.B) {
    b.Run("BuilderVsConcat", func(b *testing.B) {
        // Compare builder vs string concatenation
    })
    
    b.Run("BuilderVsStringsBuilder", func(b *testing.B) {
        // Compare against stdlib strings.Builder
    })
    
    b.Run("TranslationFunction", func(b *testing.B) {
        // Benchmark T() vs current Err() approach
    })
}
```

### Memory Validation
- Profile before/after with `pprof`
- Verify zero escape to heap with `-gcflags="-m"`
- Confirm pool efficiency with allocation counters

## Success Criteria

### Performance Targets
- **Allocations**: 50%+ reduction in high-demand processes
- **Speed**: Match or exceed `strings.Builder` performance  
- **Memory**: Zero additional persistent memory usage
- **Binary size**: No increase in compiled size

### API Requirements
- **Backward compatibility**: All existing APIs unchanged
- **Fluent interface**: Chainable builder methods
- **Pool transparency**: Automatic memory management
- **TinyGo compatibility**: Full support for embedded targets

## Risk Mitigation

### Potential Issues
1. **State corruption**: Builder and conv state conflicts
   - **Mitigation**: Clear separation of responsibilities
   
2. **Pool efficiency**: Builder usage affecting pool performance
   - **Mitigation**: Existing pool patterns proven effective

3. **API confusion**: Too many ways to build strings
   - **Mitigation**: Clear documentation and focused use cases

### Validation Steps
1. **Unit tests**: All builder methods with edge cases
2. **Integration tests**: Builder + existing functionality  
3. **Translation tests**: T function with all supported languages
4. **Benchmark validation**: Performance meets targets
5. **Memory analysis**: No regressions in allocations

## Conclusion

This design leverages TinyString's existing memory optimization infrastructure while adding a high-performance string builder API. The `T` function provides efficient translation functionality, and the builder methods enable zero-allocation string construction for high-demand processes.

The implementation maintains full backward compatibility while providing significant performance improvements for memory-intensive operations.
