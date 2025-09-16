# Fprintf Implementation Plan for TinyString

## Overview
Implement `Fprintf` functionality to provide `fmt.Fprintf`-compatible formatting that writes directly to an `io.Writer` interface. This will extend the existing `Fmt` function to support output streaming instead of returning strings.

## Target Signature
```go
func Fprintf(w io.Writer, format string, args ...any) (n int, err error)
```

## Current Code Analysis

### Existing Foundation
1. **Fmt Function** (`fmt_template.go:9-14`): Already implements printf-style formatting using the internal `wrFormat` method
2. **Memory Management** (`memory.go`): Comprehensive buffer management system with:
   - Object pooling via `convPool` and `GetConv()/putConv()`
   - Universal buffer methods: `WrString()`, `wrBytes()`, `wrByte()`
   - Three buffer destinations: `BuffOut`, `BuffWork`, `BuffErr`
3. **Format Engine** (`wrFormat`): Complete printf implementation supporting all standard format specifiers
4. **Type Conversion**: Rich set of `wr*` methods for different data types:
   - `wrIntBase()` for integers with base conversion
   - `wrFloat32()`, `wrFloat64()` for floating point numbers
   - `wrBool()` for boolean values

### Key Existing Methods to Reuse
- `GetConv()`: Get pooled converter object
- `putConv()`: Return object to pool
- `wrFormat(dest, format, args...)`: Core formatting engine
- `GetString(dest)`: Extract string from buffer
- `hasContent(BuffErr)`: Error checking

## Implementation Strategy

### 1. Core Function Structure
```go
func Fprintf(w io.Writer, format string, args ...any) (n int, err error) {
    // Obtain converter from pool
    c := GetConv()
    defer c.putConv() // Ensure cleanup
    
    // Use existing wrFormat to populate buffer
    c.wrFormat(BuffOut, format, args...)
    
    // Check for formatting errors
    if c.hasContent(BuffErr) {
        return 0, c
    }
    
    // Write to io.Writer
    data := c.getBytes(BuffOut)
    return w.Write(data)
}
```

### 2. Import Requirements
- Add `io` import to access `io.Writer` interface
- Consider adding `errors` import or use tinystring's error system

### 3. Error Handling Strategy
Two approaches to consider:
- **Option A**: Use Go's standard `errors.New()` for io.Writer compatibility
- **Option B**: Convert tinystring errors to standard errors for consistency

### 4. Memory Optimization
- Leverage existing object pooling system
- Reuse all existing buffer management
- No additional allocations beyond what `Fmt` already uses

### 5. Performance Considerations
- Zero additional overhead compared to `Fmt` + `w.Write()`
- Same formatting performance as existing `wrFormat`
- Benefit from existing memory pool optimizations

## Implementation Plan

### Phase 1: Basic Implementation
1. Add necessary imports (`io` and error handling)
2. Implement basic `Fprintf` function using existing `wrFormat`
3. Add proper error handling for both formatting and writing errors
4. Ensure memory cleanup with `defer c.putConv()`

### Phase 2: Error Integration
1. Decide on error handling approach (standard vs tinystring errors)
2. Implement proper error conversion if using tinystring errors
3. Test error propagation scenarios

### Phase 3: Optimization
1. Consider direct buffer writing to avoid string conversion
2. Evaluate if streaming write is beneficial for large outputs
3. Benchmark against standard `fmt.Fprintf`

### Phase 4: Documentation & Testing
1. Add comprehensive tests covering all format specifiers
2. Test with various `io.Writer` implementations
3. Add examples to README.md
4. Performance benchmarks

## File Modifications Required

### New Function Location
Add `Fprintf` to `fmt_template.go` alongside the existing `Fmt` function for consistency.

### Import Changes
```go
import (
    "io"
    "unsafe"
)
```

## API Integration

### README.md Updates
Add to the "fmt Package" equivalents section:
```markdown
| Go Standard | TinyString Equivalent |
|-------------|----------------------|
| `fmt.Fprintf()` | `Fprintf(w, format, args...)` |
```

### Usage Examples
```go
// Write to file
file, _ := os.Create("output.txt")
Fprintf(file, "Hello %s, count: %d\n", "world", 42)

// Write to buffer
var buf bytes.Buffer
Fprintf(&buf, "Formatted: %v", data)

// Write to HTTP response
Fprintf(w, "JSON: %s", jsonData)
```

## Risk Assessment

### ToLower Risk
- Reusing existing, tested formatting engine
- Leveraging proven memory management
- Following established patterns in the codebase

### Medium Risk
- Error handling integration between tinystring and standard library
- Import dependencies (minimal impact given `unsafe` already imported)

### Mitigation
- Comprehensive testing with various io.Writer implementations
- Benchmark to ensure no performance regression
- Clear documentation of error handling behavior

## Success Criteria
1. **Functionality**: Full compatibility with `fmt.Fprintf` format specifiers
2. **Performance**: No more than 5% overhead compared to `fmt.Fprintf`
3. **Memory**: Reuse existing pool system, no additional permanent allocations
4. **Integration**: Seamless addition to existing API without breaking changes
5. **Documentation**: Clear examples and usage patterns

## Future Enhancements
- Consider `Fprintln` variant for consistency
- Evaluate streaming capabilities for very large outputs
- Potential integration with tinystring's multilingual error system
