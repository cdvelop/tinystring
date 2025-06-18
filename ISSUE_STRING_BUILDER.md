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
// REFACTORED: Convert function with variadic parameters (only accepts 0 or 1 value)
func Convert(v ...any) *conv          // Convert() or Convert(initialValue)

// Unified builder method - detects type automatically
func (c *conv) Write(v any) *conv     // Universal write method (string, byte, rune, numbers)
func (c *conv) Reset() *conv          // Reset complete conv state

// Build result - reuses existing String() method  
func (c *conv) String() string        // Already exists - auto-releases to pool

// CLEAN API: Convert() automatically uses builder internally
// Example: Convert("hello ").Write("tiny").Write(" string").ToUpper().String()
// Output: "HELLO TINY STRING"

// BUILDER PATTERN: Empty initialization for loops
// Example: 
// c := Convert()
// for _, item := range items {
//     c.Write(item).Write(" ")
// }
// result := c.String()
```

### 2. Translation Function `T`

#### New Function Signature
```go
func T(values ...any) string
```

#### Usage Examples
```go
// Convert automatically initializes buffer 
c := Convert("hello ")                    // buf = "hello "
c.Write("tiny").Write(" string")         // buf = "hello tiny string"  
c.ToUpper()                              // buf = "HELLO TINY STRING"
result := c.String()                     // "HELLO TINY STRING" + auto-release

// Chaining with mixed operations (should work seamlessly)
result := Convert("hello").ToUpper().Write(" WORLD").ToLower().String()
// Output: "hello world"

// Error chain interruption - operations after error are omitted
c := Convert("valid").Write("ok").Write(make(chan int)) // chan int = unsupported type
c.ToUpper().Write(" MORE")  // ToUpper and Write(" MORE") are OMITTED internally
result, err := c.StringError() // result = "validok", err = "unsupported type error"

// Clean error handling
if err != nil {
    log.Printf("Chain error: %v", err) // Developer responsibility
}
```

### 3. Memory Optimization Strategy

#### ARCHITECTURAL CHANGE: Single Source of Truth + Clean API
- **ELIMINATE `stringVal` AND `tmpStr`**: Use only `buf []byte` as string storage
- **Unified `Write()` method**: Automatically detects and handles string, byte, rune, numbers
- **`Reset()` method**: Reset complete conv state (all fields + buffer)
- **`Convert()` auto-builder**: Transparent internal builder usage for all operations
- **Memory trade-off**: Slightly more memory for small texts, zero allocations for complex operations
- **Perfect for TinyGo/WebAssembly**: Optimized for web applications with medium/large text processing

#### Buffer-Only Strategy
```go
// UNIFIED approach - everything uses buf as single source of truth
func complexOperation() string {
    c := Convert()
    // All operations append to buf directly
    c.Write("prefix").Write(" ").Write("suffix") // buf contains: "prefix suffix"
    result := c.String() // Convert buf to string and auto-release
    return result
}

// Transparent builder usage in existing functions
func (t *conv) ToUpper() *conv {
    // Always work with buf, eliminate stringVal conflicts
    if len(t.buf) == 0 {
        // Initialize buf from current value
        t.initBufferFromValue()
    }
    t.processBufferToUpper() // Single processing path
    return t
}
```

## Implementation Plan

### Phase 1: Builder Methods Implementation + Refactoring
**File**: `builder.go`

1. **Convert function refactoring**:
   - **CRITICAL**: Change `Convert(v any)` to `Convert(v ...any)` with validation
   - **Rule**: Only accepts 0 or 1 parameters: `Convert()` or `Convert(value)`
   - **Purpose**: Enable empty initialization for builder pattern in loops
   - **Error handling**: Use `c.err = T(D.Only, D.One, D.Value, D.Supported)` for multiple parameters (consistent with library pattern)

2. **Unified type handling method**:
   - `type cm uint8` - Convert mode constants (mi=inicial, mb=buffer, ma=any)
   - `setVal(v any, mode cm)` - Consolidate ALL type switches (withValue, Write, any2s)
   - `Write(v any) *conv` - Universal write method using setVal(v, mb)
   - `Reset() *conv` - Reset complete conv state

2. **Core refactoring**:
   - Eliminate `withValue()` - replace with `setVal(v, mi)` 
   - Rename `ensureCapacity()` to `grow()` 
   - Replace `getString()` with `getBuf()` (get buffer value - heavily used)   - Create `val2Buf()` method for direct buffer conversion
   - Unify `i2sBuf()`, `f2sBuf()` with `appendValueToBuf(v any, typ vTpe)` method
   - **ELIMINATE `stringVal` AND `tmpStr` fields**: Use only `buf` as single source of truth

3. **setString() transition strategy**:
   - Mark `setString()` as **@deprecated** during transition
   - Gradual replacement throughout library with `setBuf()`

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

#### `Convert(v ...any) *conv` - Refactored Constructor
```go
func Convert(v ...any) *conv {
    c := getConv()
    
    // Validation: Only accept 0 or 1 parameter
    if len(v) > 1 {
        c.err = T(D.Only, D.One, D.Value, D.Supported) // Consistent error handling pattern
        return c
    }
    
    // Initialize with value if provided, empty otherwise
    if len(v) == 1 {
        c.setVal(v[0], mi) // Initial mode
    }
    // If no value provided, conv is ready for builder pattern
    
    return c
}
```

**Usage Examples**:
```go
// Traditional usage (backward compatible)
result := Convert("hello").ToUpper().String() // "HELLO"

// Builder pattern for loops (NEW - eliminates multiple allocations)
c := Convert() // Empty initialization
for _, word := range []string{"hello", "tiny", "string"} {
    c.Write(word).Write(" ")
}
result := c.String() // "hello tiny string "

// Mixed usage
c := Convert("prefix: ") // Initialize with value
for i := 0; i < 5; i++ {
    c.Write(i).Write(" ")
}
result := c.String() // "prefix: 0 1 2 3 4 "
```

#### `Write(v any) *conv` - Unified Write Operation
```go
func (c *conv) Write(v any) *conv {
    if c.err != nil {
        return c  // Error chain interruption
    }
    
    // Use unified type handler with write mode
    c.setVal(v, cmWrite)
    return c
}
```

The Write method provides a unified interface for appending any value to the buffer. It delegates to the unified `setVal()` handler with write mode, which ensures consistent type conversion and error handling across all operations.

#### `Reset() *conv` - Reset Builder
```go
func (c *conv) Reset() *conv {
    // Reset all conv fields to default state
    // ELIMINATED: c.stringVal = ""
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
    c.buf = c.buf[:0] // Reset buffer length, keep capacity - SINGLE SOURCE OF TRUTH
    return c
}
```

### Method Naming Conventions

#### Short API Names & Implementation Notes
- **`Write()`**: Universal write method - detects type automatically using `setVal(v, mb)`
- **`Reset()`**: Complete conv reset (all fields + buffer)
- **`grow()`**: Capacity management (renamed from `ensureCapacity` to avoid Go built-in conflicts)
- **`getBuf()`**: Get buffer value - heavily used internal method (documented with //)
- **`setVal()`**: **UNIFIED** type handling - consolidates ALL type switches into single method
- **`val2Buf()`**: Direct conversion to buffer using unified `appendValueToBuf(v any, typ vTpe)` method
- **Error checking**: Each operation internally verifies `c.err` before processing
- **`cm` type**: Convert mode constants (mi, mb, ma) for type-safe mode selection

#### Critical Implementation Details
- **Buffer initialization strategy**: `c.buf = append(c.buf[:0], val...)` (reset + set)
- **Convert mode constants**: `mi`=inicial, `mb`=buffer, `ma`=any (type cm uint8)
- **Error chain behavior**: **ALL operations internally check `c.err` before processing** (error contaminates chain)
- **Memory layout optimization**: Hot fields (`buf`, `vTpe`, `err`) placed first in struct for better cache locality
- **`val2Buf()` strategy**: Direct conversion to buffer using unified `appendValueToBuf(v any, typ vTpe)` method
- **Performance**: Each operation has minimal error check overhead: `if c.err != "" { return c }`

### Integration Points

#### High-Demand Process Optimization (PRIORITY ORDER)
**CRITICAL** (Immediate Implementation):
1. **`joinSlice()` function**: Replace makeBuf + append with builder (convert.go:299)
2. **`Err()` function**: Use `T` with builder internally for translation concatenation
3. **`Fmt()` function**: Replace string concatenation with builder throughout format.go
4. **Transformation chains**: Single buffer for ToUpper, ToLower, CamelCase operations

**MODERATE** (Secondary Implementation):
5. **`any2s()` conversions**: Use buffer for all type-to-string operations
6. **Dictionary translations**: Efficient multi-part message construction
7. **Error message construction**: Throughout library error handling

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
    c := getConv() // Always obtain from pool
    c.err = T(values...)
    c.vTpe = typeErr
    return c
}
```

#### Builder Integration
- **T function** uses builder API internally for zero-allocation string construction
- **Separation of concerns**: Translation logic isolated from error handling
- **Reusability**: T function can be used independently for any translation needs

## Testing Strategy

### New Benchmarks + Critical Tests
Add to `benchmark_strings_test.go`:

```go
func BenchmarkBuilderOperations(b *testing.B) {
    b.Run("BuilderVsConcat", func(b *testing.B) {
        // Compare builder vs string concatenation
    })
    
    b.Run("UnifiedWrite", func(b *testing.B) {
        // Benchmark Write() with different types
    })
    
    b.Run("ChainedOperations", func(b *testing.B) {
        // Test: Convert().Write().ToUpper().Write().String()
    })
    
    b.Run("EmptyConvertLoop", func(b *testing.B) {
        // CRITICAL: Test Convert() + Write() in loops vs multiple Convert(value)
        // This validates the main optimization goal
        items := []string{"a", "b", "c", "d", "e"}
        
        b.Run("BuilderPattern", func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                c := Convert() // Empty initialization
                for _, item := range items {
                    c.Write(item).Write(" ")
                }
                _ = c.String()
            }
        })
        
        b.Run("MultipleAllocations", func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                result := ""
                for _, item := range items {
                    result += Convert(item).String() + " "
                }
            }
        })
    })
}

func TestConvertVariadicValidation(t *testing.T) {
    // CRITICAL: Test Convert() parameter validation
    
    // Valid usage
    c1 := Convert()          // Empty - should work
    c2 := Convert("hello")   // Single value - should work
    
    if c1.err != "" {
        t.Errorf("Convert() should not have error, got: %s", c1.err)
    }
    if c2.err != "" {
        t.Errorf("Convert(value) should not have error, got: %s", c2.err)
    }
    
    // Invalid usage - should set error and continue chain
    c3 := Convert("hello", "world") // Multiple values - should set error
    if c3.err == "" {
        t.Error("Convert with multiple parameters should set error")
    }
    
    // Chain should continue but operations should be omitted due to error
    result := c3.Write(" more").ToUpper().String()
    if result != "" {
        t.Errorf("Operations after error should be omitted, got: %s", result)
    }
}

func TestChainedOperationsOrder(t *testing.T) {
    // CRITICAL: Test mixed chaining behavior
    result := Convert("hello").ToUpper().Write(" WORLD").ToLower().String()
    expected := "hello world"
    if result != expected {
        t.Errorf("Chained operations failed: got %q, want %q", result, expected)
    }
}

func TestErrorChainInterruption(t *testing.T) {
    // CRITICAL: Test error chain interruption behavior
    c := Convert("valid").Write("ok").Write(make(chan int)) // Unsupported type
    c.ToUpper().Write(" MORE") // These should be omitted internally
    
    result, err := c.StringError()
    expectedResult := "validok" // Only operations before error
    if result != expectedResult {
        t.Errorf("Error chain result: got %q, want %q", result, expectedResult)
    }
    if err == nil {
        t.Error("Expected error for unsupported type, got nil") 
    }
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
