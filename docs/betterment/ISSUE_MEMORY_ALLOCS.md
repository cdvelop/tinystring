# MEMORY ALLOCATION OPTIMIZATION ISSUE

## **CONSTRAINTS** 📝
- **WebAssembly-first**: Binary size over performance
- **No stdlib**: Manual implementations only
- **Dictionary errors**: Use D.* constants only
- **TinyGo compatible**: Limited reflection

## Current Problem
- **TinyString has 189 allocs/op vs Standard Library 48 allocs/op**
- **4050 B/op vs 1200 B/op** - 3.3x more memory usage
- **25945 ns/op vs 6597 ns/op** - 3.9x slower

## Root Cause Analysis
1. **Excessive ptrValue usage**: Storing simple types (string, int, bool) in `any` interface causes boxing allocations
2. **Double storage**: Data stored in both buffer AND ptrValue for simple types
3. **getBuffString() redundancy**: Re-converts from ptrValue even when buffer has data

## Optimization Plan

### Phase 1: Restrict ptrValue Usage
- **KEEP ptrValue for**: `[]string`, `map[string]string`, `map[string]any`, `*string` (for Apply())
- **ELIMINATE ptrValue for**: `string`, `int*`, `uint*`, `float*`, `bool`, `[]byte`

### Phase 2: Optimize anyToBuff()
```go
// BEFORE (causes allocation)
case string:
    c.kind = KString
    c.ptrValue = v      // ❌ Remove this boxing
    c.wrString(dest, v)

// AFTER (buffer-only)
case string:
    c.kind = KString
    c.wrString(dest, v) // ✅ Buffer only
    // No ptrValue assignment
```

### Phase 3: Optimize getBuffString()
```go
// BEFORE (double work)
if c.ptrValue != nil {
    // Re-converts even if buffer has data
}

// AFTER (buffer-first)
if c.outLen > 0 {
    return c.getString(buffOut) // ✅ Use existing buffer
}
// Only fallback for complex types
```

### Phase 4: Update struct definition
```go
type conv struct {
    // ...buffers...
    kind kind
    
    // RESTRICTED: Only for complex types and pointers
    ptrValue any // []string, map, *string ONLY
}
```

## Expected Results
- **Target**: Reduce from 189 to <50 allocs/op
- **Memory**: Reduce from 4050B to <1500B/op  
- **Speed**: Improve from 25945ns to <10000ns/op

## Decisions Made

### 1. String Pointer Handling (*string)
- **KEEP reference**: Maintain ptrValue for *string to support Apply()
- **IMMEDIATE conversion**: Dereference *string and store in buffer immediately
- **Apply() method**: Read from buffer out, write back to original pointer

### 2. Complex Types ([]string, map)
- **LAZY conversion**: Keep ptrValue, convert on-demand
- **Specific handlers**: Create dedicated methods for each complex type

### 3. Simple Types (string, int, bool, float)
- **ELIMINATE ptrValue**: No boxing for simple types
- **Buffer-only storage**: Direct write to buffers

### 4. Number Optimization Strategy
**DECISION**: Use Opción B - Immediate string conversion (APPROVED)
- **Rationale**: Treating all as string/bytes is simpler and more flexible
- **Benefits**: Easier type conversion, unified buffer architecture, less memory overhead
- **Implementation**: Convert numbers to string immediately in anyToBuff(), store only in buffers

### 5. API Refactoring Decisions
- **getBuffString()**: Must use wrString API internally
- **Method renaming**: `fmtIntToOut` → `wrInt` with buffDest parameter
- **Consistency**: All write methods follow pattern: `wr{Type}(dest buffDest, value)`
- **File organization**: Each file handles its responsibilities (bool.go, numeric, etc.)

### 6. Apply() Method Scope
- **Limitation**: Only works with same type matching (kind validation)
- ***string only**: Apply() limited to string pointers, not numeric pointers
- **Buffer-to-pointer**: Apply buffer out content to original pointer if types match
- **Error handling**: Silent no-op if kind != KPointer (dev can use StringError() if needed)

### 7. Memory Management
- **ptrValue cleanup**: Continue setting `ptrValue = nil` in putConv() (most efficient)
- **API Pattern**: All methods follow `wr{Type}(dest buffDest, value {type})` exactly

## Clean Code & Minimal Binary Policy
- **No dead code**: All unused, unreachable, or obsolete code must be removed immediately after refactoring.
- **No legacy branches**: After migration, fallback or legacy code paths must be deleted, not commented.
- **Minimal binary**: Every line must be justified for runtime or API; no helpers or debug code left behind.
- **Review**: Each refactor must include a pass to remove all code that is no longer needed.
- **Commit rule**: No commit is allowed if dead/unused code remains after a refactor.

## Implementation Strategy
1. **API Refactoring First**: 
   - Rename `fmtIntToOut` → `wrInt(dest buffDest, val int64)`
   - Create `wrBool(dest buffDest, val bool)` in bool.go
   - Create `wrUint(dest buffDest, val uint64)` for consistency
   - Ensure ALL wr* methods follow exact pattern: `wr{Type}(dest buffDest, value {type})`
2. **Modify anyToBuff()**: 
   - Eliminate ptrValue assignment for simple types (string, int, bool, float)
   - Keep ptrValue only for: *string, []string, map
   - Convert numbers to string immediately via wr* methods
3. **Replace getBuffString()**: 
   - Rename to `getBuffString(dest buffDest)`
   - Remove ptrValue logic for simple types
   - Use buffer-only approach
4. **Update putConv()**: Keep `ptrValue = nil` (most efficient cleanup)
5. **Test validation**: Run tests after each change
6. **Benchmark comparison**: Measure improvement
7. **Commit**: Only if tests pass and performance improves

## Final Architecture
```go
type conv struct {
    // Buffer storage (main approach)
    out, work, err []byte
    outLen, workLen, errLen int
    kind kind
    
    // ONLY for complex types and pointers needing Apply()
    ptrValue any // *string, []string, map ONLY
}
```

## 🛠️ **TOOLS & COMMANDS**
```bash
cd /c/Users/Cesar/Packages/Internal/tinystring

go build -gcflags="-m" ./... # Detect variables escaping to the heap
go test -bench=. -benchmem -memprofile=mem.prof # Profile memory from benchmarks
go tool pprof -text ./yourbinary.test mem.prof # Analyze memory profile
benchstat old.txt new.txt # Compare memory usage between versions
./memory-benchmark.sh  # Run memory benchmarks results in benchmark/benchmark_results.md
```

## Success Metrics
- ✅ Allocation count reduced by 70%+
- ✅ Memory usage reduced by 60%+
- ✅ All existing tests pass
- ✅ Apply() method still works for *string

## Continuous Development Notes
- Dead code in `fmt_number.go` has been removed and tests have been run and passed.
- Rounding logic in `fmt_precision.go` has been adjusted to align with test expectations.
- Los intentos de optimización de capacidad de buffer no fueron efectivos y se revirtieron.
- El enfoque actual se centra en la optimización de la manipulación de cadenas/bytes en métodos como `Capitalize`, `ToLower`, `ToUpper` para reducir las asignaciones.
