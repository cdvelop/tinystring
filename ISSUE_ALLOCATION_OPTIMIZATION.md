# tinystring Memory Allocation Optimization - UNIFIED BUFFER ARCHITECTURE

## 🎯 **CURRENT STATUS**

**Context:** WebAssembly-first library with manual implementations (no stdlib dependencies)

**Performance Targets:**
- String Processing: 2.8KB/op, 119 allocs/op → **Reduce 50%** 🚧 **IN PROGRESS**
- Mixed Operations: 1.7KB/op, 54 allocs/op → **Reduce 40%** 🚧 **IN PROGRESS**  
- Binary Size: 55.1% better than stdlib ✅ **MAINTAINED**

**Current Phase:** Implement unified buffer architecture with single conversion function

⚠️ **CRITICAL CONSTRAINTS:** TinyString operates at **133% higher memory usage** and **173% more allocations** than standard library as baseline. Optimization targets are relative to current TinyString performance, prioritizing binary size reduction for WebAssembly deployment over runtime efficiency.

## 🏗️ **UNIFIED BUFFER ARCHITECTURE**

### **✅ COMPLETED: Core Foundation**
- **Centralized Buffer Management:** All operations in `memory.go` ✅
- **Three-Buffer Strategy:** `out`, `work`, `err` with length control ✅  
- **getString() Elimination:** Replaced with `ensureStringInOut()` ✅
- **Pool Management:** Optimized with centralized reset ✅
- **Error System Refactor:** Non-recursive error handling ✅

### **🚧 IMPLEMENTING: Universal Conversion**
```go
// TARGET: Single Universal Conversion Function
func anyToBuff(c *conv, dest buffDest, value any) 

// Buffer Destination Selection
type buffDest int
const (
    buffOut buffDest = iota  // Primary output
    buffWork                 // Working/temporary  
    buffErr                  // Error messages
)

// Language-Aware Error Reporting  
func (c *conv)wrErr(values ...any) 
```

### **✅ CURRENT `conv` Structure (Migration Phase):**
```go
type conv struct {
    // CENTRALIZED BUFFERS ✅ IMPLEMENTED
    out     []byte // Primary output - main result
    outLen  int    // Output length control ✅ 
    work    []byte // Work/temporary operations  
    workLen int    // Work length control ✅
    err     []byte // Error messages
    errLen  int    // Error length control ✅ ADDED
    
    // TYPE INDICATOR ✅ KEEP - Hot path for type checking
    kind    kind   // Type differentiation for conversion logic ✅ REQUIRED
    pointerVal   *string  // ✅ Keep (Apply() support)
    
    // TYPE VALUES ⏳ TO ELIMINATE WITH anyToBuff()
    intVal         int64    // → direct parameter to anyToBuff()
    uintVal        uint64   // → direct parameter to anyToBuff()
    floatVal       float64  // → direct parameter to anyToBuff()
    boolVal        bool     // → direct parameter to anyToBuff()
}
```

### **✅ IMPLEMENTED: Buffer API Foundation**
```go
// MAIN BUFFER API ✅ WORKING
func (c *conv) wrStringToOut(s string)  // ✅ Primary output writing
func (c *conv) wrToOut(data []byte)     // ✅ Byte writing  
func (c *conv) rstOut()                  // ✅ Position reset
func (c *conv) getOutString() string       // ✅ String reading
func (c *conv) ensureStringInOut() string  // ✅ Conversion + reading

// UNIVERSAL CONVERSION ENTRY POINTS ✅ WORKING
func ensureStringInOut(c *conv) string             // ✅ Buffer-to-string reading
```

## ✅ **ARCHITECTURAL DECISIONS CORRECTED**

### **1. Universal Conversion Function - NO ERROR RETURN**
```go
// CORRECTED SIGNATURE - NO ERROR RETURN
func anyToBuff(c *conv, dest buffDest, value any)

// Supported Types: string, int, float, bool, []byte, LocStr
// Buffer Selection: buffOut, buffWork, buffErr
// Error Handling: Writes to c.err using c.wrErr(...), caller checks len(c.err) > 0
```

### **2. Language-Aware Error System - NO ERROR RETURN**  
```go
// CORRECTED SIGNATURE - NO ERROR RETURN
func wrErr(c *conv, dest buffDest, lang lang, msgs ...LocStr)

// Features:
// - Direct buffer writing (no T() dependency) 
// - Language detection integration
// - LocStr translation support
// - Writes error to c.err, no return value
```

### **3. Error Checking Pattern - CORRECTED**
```go
// USAGE PATTERN FOR ALL OPERATIONS - Using Length Fields
c := getConv()
anyToBuff(c, buffOut, value)
if c.hasError() {  // ✅ Use errLen field via method
    // Handle error condition
    return c  // Return conv with error set
}
// Continue with normal operation
```

### **Buffer State Checking Methods**
```go
// ✅ REQUIRED: Buffer state checking methods using length fields
func (c *conv) hasError() bool      { return c.errLen > 0 }
func (c *conv) hasWorkContent() bool { return c.workLen > 0 }
func (c *conv) hasOutContent() bool  { return c.outLen > 0 }
func (c *conv) isEmpty() bool        { return c.outLen == 0 && c.workLen == 0 && c.errLen == 0 }
func (c *conv) clearError()          { c.errLen = 0 }

// ❌ INCORRECT - Never use direct len() checks:
if len(c.err) > 0 { }     // Wrong: doesn't respect errLen
if len(c.work) > 0 { }    // Wrong: doesn't respect workLen  
if len(c.out) > 0 { }     // Wrong: doesn't respect outLen

// ✅ CORRECT - Use state checking methods:
if c.hasError() { }       // Correct: uses errLen
if c.hasWorkContent() { } // Correct: uses workLen
if c.hasOutContent() { }  // Correct: uses outLen
```

### **4. Buffer Writing Logic for Convert() - CORRECTED**
```go
// CONVERT() BUFFER FLOW SPECIFICATION - USING DICTIONARY:
func Convert(v ...any) *conv {
    c := getConv()
    
    // STEP 1: Validation - ALWAYS USE DICTIONARY
    if len(v) > 1 {
        c.wrErr(D.Invalid, D.Number, D.Of, D.Argument)  // ✅ Dictionary usage
        return c
    }
    
    // STEP 2: Value conversion
    if len(v) == 1 {
        val := v[0]
        if val == nil {
            c.wrErr(D.String, D.Empty)  // ✅ Dictionary usage
            return c
        }
          // CONVERSION FLOW: value → work (first conversion)
        anyToBuff(c, buffWork, val)  // Convert to work buffer
        if c.hasError() {  // ✅ Use errLen field via method
            return c  // Return if conversion failed
        }
        
        // NO COPY TO OUT - First conversion stays in work
        // Second conversion (public API) will move work → out
    }
    
    return c
}
```

### **5. Buffer Destination Enum**
```go
// FINAL ENUM CONFIRMED
type buffDest int
const (
    buffOut buffDest = iota  // Primary output buffer
    buffWork                 // Working/temporary buffer  
    buffErr                  // Error message buffer
)
```

### **6. SMART TYPE HANDLING - OPTIMIZED APPROACH**
```go
// STRATEGY: Immediate conversion for simple types, pointer storage for complex types

// SIMPLE TYPES → Direct buffer conversion (anyToBuff)
// - string, int, float, bool, []byte → immediate conversion to buffer

// COMPLEX TYPES → Pointer storage + lazy conversion
// - []string, map[string]string, map[string]any → store pointer, convert on demand

// TARGET SIMPLIFIED STRUCT:
type conv struct {
    out     []byte   // Primary output buffer
    outLen  int      // Output length control
    work    []byte   // Working buffer  
    workLen int      // Working length control
    err     []byte   // Error buffer
    errLen  int      // Error length control
    
    // REQUIRED FIELDS
    kind        kind        // ✅ Type checking - hot path optimization
    pointerVal  any         // ✅ Universal pointer for complex types ([]string, map[string]any, etc.)
}

// USAGE PATTERNS:
// Convert(42)                    → anyToBuff(buffWork, "42")          // Direct conversion
// Convert([]string{"a", "b"})    → pointerVal = &slice, kind = KSliceStr // Pointer storage
// Convert(map[string]any{...})   → pointerVal = &map, kind = KMap     // Pointer storage
```

### **COMPLEX TYPE HANDLING STRATEGY**

### **Recommended Approach: Immediate vs Lazy Conversion**

```go
// IMMEDIATE CONVERSION (Simple Types)
// These types convert directly to buffer at Convert() time
Convert("hello")        → anyToBuff(buffWork, "hello")     // Direct to work
Convert(42)            → anyToBuff(buffWork, "42")        // Convert & write  
Convert(true)          → anyToBuff(buffWork, "true")      // Convert & write
Convert([]byte{...})   → anyToBuff(buffWork, data)       // Direct copy

// LAZY CONVERSION (Complex Types) 
// These types store pointer, convert on first operation
Convert([]string{"a", "b"})     → pointerVal = slice, kind = KSliceStr
Convert(map[string]string{...}) → pointerVal = map, kind = KMap
Convert(map[string]any{...})    → pointerVal = map, kind = KMap
```

### **Why This Hybrid Approach is Best:**

1. **Performance**: Simple types convert once, complex types convert on-demand
2. **Memory**: No unnecessary string allocations for unused complex conversions  
3. **Flexibility**: Complex types can be converted differently based on operation
4. **Code Simplicity**: Less conditional logic, cleaner implementation

### **Complex Type Conversion Examples:**

```go
// []string handling
slice := []string{"apple", "banana", "cherry"}
c := Convert(slice)  // pointerVal = slice, kind = KSliceStr

c.Join()          → "apple banana cherry"     // Default space join
c.Join(",")       → "apple,banana,cherry"     // Custom separator  
c.Count()         → "3"                       // Count elements
c.First()         → "apple"                   // First element
```

### **Map Handling Options:**

```go
// map[string]string
data := map[string]string{"name": "John", "age": "30"}
c := Convert(data)  // pointerVal = data, kind = KMap

// Different conversion strategies based on operation:
c.String()        → "name=John age=30"        // Key=value pairs
c.ToJSON()        → `{"name":"John","age":"30"}` // JSON format
c.Keys()          → "name age"                // Keys only
c.Values()        → "John 30"                 // Values only
```

### **Implementation Benefits:**

1. **No Kind Proliferation**: Use existing `KSliceStr`, `KMap` instead of creating `KSliceStrPtr`
2. **Universal `pointerVal`**: Single field handles all complex types via `any`
3. **Lazy Conversion**: Only convert when needed, based on actual operation
4. **Operation-Specific**: Same data can be converted differently per operation
```go
// These fields will be removed after anyToBuff() implementation:
// - intVal, uintVal, floatVal, boolVal (replaced by direct parameters)
// - kind (replaced by buffDest + type reflection)
// - All buffer-specific methods (fmtIntToWork, etc.)
```

## 🚧 **IMPLEMENTATION ROADMAP - CORRECTED**

### **Phase 1: Universal Conversion Implementation - UPDATED**
```go
// STEP 1: Implement anyToBuff() with Hybrid Conversion Strategy
func anyToBuff(c *conv, dest buffDest, value any) {
    switch v := value.(type) {
    // IMMEDIATE CONVERSION - Simple Types
    case string:
        writeStringToDest(c, dest, v)
    case int, int8, int16, int32, int64:
        str := intToString(v)
        writeStringToDest(c, dest, str)
    case uint, uint8, uint16, uint32, uint64:
        str := uintToString(v)
        writeStringToDest(c, dest, str)
    case float32, float64:
        str := floatToString(v)
        writeStringToDest(c, dest, str)
    case bool:
        str := boolToString(v)
        writeStringToDest(c, dest, str)
    case []byte:
        writeBytesToDest(c, dest, v)
    case LocStr:
        str := translateLocStr(v)
        writeStringToDest(c, dest, str)
        
    // LAZY CONVERSION - Complex Types
    case []string:
        c.pointerVal = v
        c.kind = KSliceStr
        // No immediate conversion - wait for operation
        
    case map[string]string, map[string]any:
        c.pointerVal = v
        c.kind = KMap
        // No immediate conversion - wait for operation
        
    default:
        // Unknown type - write error using DICTIONARY
        c.wrErr(D.Type, D.Unsupported)
    }
}

// Helper function to write to correct destination
func writeStringToDest(c *conv, dest buffDest, s string) {
    switch dest {
    case buffOut:
        c.wrStringToOut(s)
    case buffWork:
        c.wrStringToWork(s)
    case buffErr:
        c.writeStringToErr(s)
    }
}
```

### **Phase 2: Error System Implementation**
```go
// STEP 2: Implement wrErr() - NO ERROR RETURN
func (c *conv) wrErr(msgs ...any) {
    // 1. Use detectLanguage() for language selection
    // 2. Translate each LocStr using getTranslation()
    // 3. Write directly to dest buffer (usually buffErr)
    // 4. No return value, no T() dependency
}

// STEP 3: Implement detectLanguage() helper
func detectLanguage(c *conv) lang {
    // 1. Check c.language if set
    // 2. Check environment variables  
    // 3. Return default fallback
}
```

### **Phase 3: Convert() Flow Specification - CORRECTED**
```go
// STEP 4: Convert() Buffer Flow Logic - USING DICTIONARY
func Convert(v ...any) *conv {
    c := getConv()
    
    // VALIDATION: errors written to c.err using DICTIONARY
    if len(v) > 1 {
        c.wrErr(D.Invalid, D.Number, D.Of, D.Argument)  // ✅ Dictionary
        return c
    }
    
    if len(v) == 1 {
        val := v[0]
        if val == nil {
            c.wrErr(D.String, D.Empty)  // ✅ Dictionary
            return c
        }
          // CONVERSION FLOW: value → work buffer ONLY
        // NO automatic copy to out - public API will handle that
        anyToBuff(c, buffWork, val)
        if c.hasError() {  // ✅ Use errLen field via method
            return c  // Return if conversion failed
        }
        
        // Set kind for type tracking
        c.kind = determineKind(val)
    }
    
    return c  // c.work contains converted value, c.out empty
}

// First public API call will transfer work → out
func (c *conv) AnyPublicMethod() *conv {
    // UPDATED: Consistent OUT-WORK-OUT pattern
    if c.hasOutContent() {
        // Standard case: out → work
        currentValue := c.getOutString()
        c.rstWork()
        c.wrStringToWork(currentValue)
    } else if c.hasWorkContent() {
        // First API after Convert(): work has initial value
        // Process directly in work, then transfer to out
    }
    
    // Perform operation in work buffer
    // ...operation logic...
    
    // Always end with work → out
    c.rstOut()
    c.wrStringToOut(c.getWorkString())
    c.rstWork()
    
    return c
}
```

### **Phase 4: Public API Migration**
```go
// STEP 4: Update all public methods to use anyToBuff()
func (t *conv) Convert(value any) *conv {
    c := getConv()
    anyToBuff(c, buffOut, value)
    t.c = c
    return t
}

func (t *conv) Fmt(format string, args ...any) *conv {
    c := getConv()
    // Use anyToBuff() for format processing
    t.c = c
    return t
}
```

### **Phase 4: Legacy Cleanup**
```go
// STEP 5: Remove legacy fields and methods
// - Remove: intVal, uintVal, floatVal, boolVal from conv struct
// - Remove: fmtIntToWork, floatToWork, etc. buffer-specific methods
// - Remove: kind-based logic, replace with direct value handling
// - Update: All remaining usages to use anyToBuff()
```

## ⚠️ **TINYSTRING LIBRARY LIMITATIONS & CONSTRAINTS**

### **📋 Architecture Design Limitations**

The TinyString library is specifically designed for **WebAssembly deployment** and **binary size optimization**, which creates inherent limitations that must be considered during the unified buffer architecture implementation:

#### **🎯 Performance Trade-offs - CRITICAL**
```go
// DOCUMENTED PERFORMANCE IMPACT - From benchmark results:
// Memory Usage: 133.3% more memory than standard library
// Allocations: 172.8% more allocations than standard library  
// Execution Time: 2-4x slower than standard library operations

// IMPACT ON OPTIMIZATION TARGETS:
// Current: 2.8KB/op, 119 allocs/op → Target: 1.4KB/op, 60 allocs/op
// Already operating at higher baseline than stdlib
```

#### **🔧 Manual Implementation Constraints**
- **No Standard Library**: Cannot use `fmt`, `strings`, `strconv`, `errors` packages
- **Custom Conversions**: All numeric/string conversions must be manually implemented
- **Limited Built-ins**: Restricted to basic Go built-in functions only
- **TinyGo Compatibility**: Must work within TinyGo's WebAssembly limitations

#### **💾 Memory Management Limitations**
```go
// BUFFER SIZE CONSTRAINTS
type conv struct {
    out  []byte  // Limited by available memory on target device
    work []byte  // Cannot use unlimited buffer growth
    err  []byte  // Must be conservative with error message length
}

// ALLOCATION PATTERNS
// ❌ Cannot rely on efficient GC patterns (embedded/WASM targets)
// ❌ Cannot use standard library's optimized buffer management
// ✅ Must implement custom pooling and reuse strategies
```

### **🌍 Localization & Language Limitations**

#### **Dictionary Constraints**
```go
// SUPPORTED LANGUAGES - FIXED SET
const supportedLanguages = 9  // EN, ES, ZH, HI, AR, PT, FR, DE, RU

// DICTIONARY SIZE LIMITATIONS
// - Only 35+ essential words available
// - Cannot add unlimited vocabulary  
// - Must compose complex messages from limited word set
// - No dynamic translation capabilities

// ERROR MESSAGE CONSTRAINTS
wrErr(D.Invalid, D.Format)  // ✅ Available
wrErr("Complex custom message with details")  // ❌ Increases binary size
```

#### **Unicode Handling Limitations**
```go
// ACCENT/DIACRITIC SUPPORT - LIMITED
RemoveTilde()  // ✅ Handles common European accents
// ❌ Limited support for complex Unicode normalization
// ❌ No support for right-to-left languages (Arabic script layout)
// ❌ No support for complex script rendering (Devanagari, Thai)
```

### **🚫 Functional Limitations**

#### **Type Support Constraints**
```go
// SUPPORTED TYPES IN anyToBuff()
string, int, int8, int16, int32, int64           // ✅ Supported
uint, uint8, uint16, uint32, uint64              // ✅ Supported  
float32, float64, bool, []byte                   // ✅ Supported
[]string, map[string]string, map[string]any      // ✅ Supported

// UNSUPPORTED TYPES
complex64, complex128                            // ❌ Not supported
interface{} (general)                            // ❌ Limited support
channels, functions, struct types               // ❌ Not supported
time.Time, custom types                         // ❌ Not supported
```

#### **Numeric Precision Limitations**
```go
// FLOATING POINT CONSTRAINTS
// Manual implementation may have different precision than standard library
ToFloat()         // Limited to manual parsing precision
RoundDecimals()   // Custom rounding, may differ from math.Round()
FormatNumber()    // Basic thousand separators only

// INTEGER LIMITATIONS  
ToInt(base)       // Supports base 2-36, but manual validation
ToUint(base)      // No negative number detection for uint conversion
```

#### **String Processing Limitations**
```go
// REGEX SUPPORT
// ❌ No regex support (regexp package would increase binary size)
// ✅ Basic string matching only (Contains, IndexByte)

// FORMATTING LIMITATIONS
Fmt(format, args...)  // ✅ Basic sprintf-style, limited verb support
// ❌ No complex formatting verbs (%+v, %#v, %T, etc.)
// ❌ No width/precision modifiers for all types

// UNICODE NORMALIZATION
// ❌ No full Unicode normalization (NFC, NFD, NFKC, NFKD)
// ✅ Basic accent removal only
```

### **⚡ Concurrency & Thread Safety Limitations**

#### **Pool Management Constraints**
```go
// OBJECT POOLING LIMITATIONS
var pool sync.Pool  // ✅ Thread-safe pool available

// CONSTRAINTS:
// - Limited to simple reset/reuse patterns
// - Cannot use complex pooling strategies due to memory constraints
// - Must be conservative with pool size on embedded targets

// GOROUTINE LIMITATIONS
// ✅ Thread-safe operations supported
// ❌ No advanced concurrency patterns (worker pools, pipelines)
// ❌ Limited by TinyGo's goroutine implementation constraints
```

### **🌐 WebAssembly Specific Limitations**

#### **Binary Size vs Feature Trade-offs**
```go
// SIZE OPTIMIZATION TARGETS CONFLICT WITH FEATURES
// Every feature addition impacts binary size targets:

// CURRENT BENCHMARKS:
// TinyString WASM: 156.1 KB (Ultra optimization)  
// Standard Lib WASM: 141.3 KB
// SIZE PENALTY: +14.8 KB for TinyString features

// FEATURE ADDITION IMPACT:
// +1KB = Significant impact on size targets
// +New dependencies = Risk of size regression
// +Complex algorithms = Memory/speed penalties
```

#### **TinyGo Compiler Constraints**
```go
// COMPILATION LIMITATIONS
// ❌ Some Go features not supported in TinyGo
// ❌ Limited reflection capabilities
// ❌ Restricted standard library subset
// ❌ Memory management differences from standard Go

// PLATFORM CONSTRAINTS  
// ✅ WebAssembly (main target)
// ⚠️ Limited testing on all embedded platforms
// ⚠️ Performance characteristics vary by target
```

### **🔧 Implementation Impact on Buffer Architecture**

#### **Buffer Size Constraints**
```go
// MUST CONSIDER IN anyToBuff() IMPLEMENTATION
func anyToBuff(c *conv, dest buffDest, value any) {
    // ⚠️ CONSTRAINT: Cannot allocate unlimited buffer sizes
    // ⚠️ CONSTRAINT: Must handle buffer overflow gracefully  
    // ⚠️ CONSTRAINT: Error messages must be concise (dictionary words only)
    // ⚠️ CONSTRAINT: Cannot use stdlib for type conversion
}
```

#### **Error Handling Constraints**
```go
// wrErr() IMPLEMENTATION MUST CONSIDER:
func (c *conv) wrErr(msgs ...any) {
    // ✅ Must use dictionary words (D.Invalid, D.Format, etc.)
    // ❌ Cannot use detailed error descriptions (binary size)
    // ❌ Cannot use fmt.Sprintf for error formatting
    // ⚠️ Limited to 9 supported languages
    // ⚠️ Error message length impacts buffer size
}
```

#### **Type Conversion Constraints**
```go
// MANUAL IMPLEMENTATIONS REQUIRED:
// ❌ Cannot use strconv.ParseInt() → Manual integer parsing
// ❌ Cannot use strconv.FormatFloat() → Manual float formatting  
// ❌ Cannot use fmt.Sprintf() → Manual format implementation
// ❌ Cannot use strings.Builder → Manual buffer management

// PRECISION/COMPATIBILITY IMPACT:
// ⚠️ Results may differ slightly from standard library
// ⚠️ Edge cases may not be handled identically  
// ⚠️ Performance characteristics are different
```

### **📊 Optimization Target Reality Check**

#### **Baseline Performance Awareness**
```go
// CURRENT PERFORMANCE CONTEXT:
// TinyString is ALREADY 133% higher memory usage than stdlib
// TinyString is ALREADY 173% more allocations than stdlib

// OPTIMIZATION TARGET FEASIBILITY:
// From: 2.8KB/op, 119 allocs/op 
// To:   1.4KB/op, 60 allocs/op (50% reduction)

// REALITY CHECK:
// - Starting from higher baseline than stdlib
// - Manual implementations limit optimization potential
// - Binary size constraints limit algorithmic complexity
// - Must balance size vs performance trade-offs
```

#### **Success Metrics Adjustment**
```go
// REALISTIC OPTIMIZATION EXPECTATIONS:
// 🎯 PRIMARY: Binary size maintenance (WebAssembly deployment)
// 🎯 SECONDARY: Memory allocation reduction within constraints
// 🎯 TERTIARY: Performance improvement where possible

// ACCEPTABLE TRADE-OFFS:
// ✅ Slower execution vs smaller binary size
// ✅ Higher memory usage vs zero stdlib dependencies  
// ✅ Limited features vs TinyGo compatibility
// ✅ Manual implementations vs automatic optimizations
```

## 🚨 **CRITICAL CONSTRAINTS FOR IMPLEMENTATION**

### **⚠️ Must Remember During Development:**

1. **No Standard Library**: All conversions must be manual implementations
2. **Binary Size Priority**: Every byte counts for WebAssembly deployment
3. **Memory Constraints**: Target devices may have limited RAM
4. **TinyGo Compatibility**: Features must work in TinyGo compilation
5. **Dictionary Only**: Error messages must use existing dictionary words
6. **Type Limitations**: Only supported types can be handled in anyToBuff()
7. **Performance Baseline**: Already operating at higher resource usage than stdlib
8. **Unicode Limitations**: Basic accent support only, no complex Unicode

### **✅ Implementation Validation Checklist:**

- [ ] **Binary Size**: New features don't increase WASM size significantly
- [ ] **TinyGo Compatibility**: Code compiles and runs in TinyGo
- [ ] **Memory Constraints**: Allocations are bounded and predictable  
- [ ] **Error Dictionary**: All error messages use D.* constants
- [ ] **Type Support**: Only supported types are handled in conversions
- [ ] **Manual Implementation**: No standard library dependencies introduced
- [ ] **WebAssembly Testing**: Features work correctly in WASM environment
- [ ] **Performance Baseline**: Improvements are measured against current TinyString baseline, not stdlib
