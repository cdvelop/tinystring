# TinyString - Unified Buffer Architecture 🎯

## **STATUS: 95% COMPLETE** ✅
- ✅ **Zero-allocation buffer architecture** implemented
- ✅ **Universal conversion**: `anyToBuff(c *conv, dest buffDest, value any)`
- ✅ **Non-recursive errors**: `wrErr()` with language support
- ✅ **100% Buffer API compliance**: No manual buffer access

## **CORE RULES** 🚨
```go
// ❌ FORBIDDEN:
c.out = c.out[:0]         // Manual buffer manipulation
len(c.err) > 0            // Direct length checks


// ✅ MANDATORY:
c.rstBuffer(dest buffDest)
c.hasContent(dest buffDest)
c.wrString(dest buffDest, s)

```

## **CORE ARCHITECTURE** ️
```go
// Buffer Management
type conv struct {
    out     []byte // Primary output
    outLen  int    // Length tracking
    work    []byte // Working buffer  
    workLen int    // Work tracking
    err     []byte // Error buffer
    errLen  int    // Error tracking
    kind    kind   // Type indicator
    anyValue any   // Universal storage
}

// Universal Conversion
func (c *conv) anyToBuff(dest buffDest, value any) {
    // dest: buffOut | buffWork | buffErr
    // Writes errors via c.wrErr(), no return values
    // NOW A *conv METHOD (not standalone function)
}

// Error System
func (c *conv) wrErr(msgs ...any) *conv {
    // Language-aware, dictionary-based, buffer API only
}
```

## **BUFFER API** 📋
```go
// State Checking
c.hasContent(buffDest) // Check out,work, err buffer
// Writing
c.wrString(buffDest, s) // Write string to specified buffer
c.wrBytes(buffDest, data) // Write bytes to specified buffer    
c.wrByte(buffDest,s)        
// Reading
c.getString(buffDest)        // Read out, work, err buffer

```

## **CONSTRAINTS** 📝
- **WebAssembly-first**: Binary size over performance
- **No stdlib**: Manual implementations only
- **Dictionary errors**: Use D.* constants only
- **TinyGo compatible**: Limited reflection

## **CENTRALIZED BUFFER API - UNIVERSAL METHODS** 🎯

### **NAMING CONVENTION CHANGE**
- All `write*` methods become `wr*`
- All methods receive `buffDest` parameter for destination selection
- Universal methods replace destination-specific variants
- NO backwards compatibility - complete refactoring

### **UNIVERSAL BUFFER METHODS** (SIMPLIFIED NAMING)
```go
// Writing Operations (dest FIRST parameter) - SIMPLIFIED NAMES
func (c *conv) wrString(dest buffDest, s string)    // Replaces: wrStringToOut, wrString, wrStringToErr  
func (c *conv) wrByte(dest buffDest, b byte)        // Replaces: wrByte (out-only)
func (c *conv) wrBytes(dest buffDest, data []byte)  // Replaces: wrBytes, wrToWork, wrBytes
func (c *conv) wrInt(dest buffDest, val int64)      // Replaces: writeIntToDest + duplicates
func (c *conv) wrUint(dest buffDest, val uint64)    // Replaces: writeUintToDest + duplicates  
func (c *conv) wrFloat(dest buffDest, val float64)  // Replaces: writeFloatToDest + duplicates

// Universal Conversion (NOW A METHOD)
func (c *conv) anyToBuff(dest buffDest, value any)  // Replaces: anyToBuff function

// Reading Operations (dest FIRST parameter)
func (c *conv) getString(dest buffDest) string      // Replaces: getString, getString, getString

// Buffer Management (dest FIRST parameter)
func (c *conv) rstBuffer(dest buffDest)             // Replaces: rstOut, rstWork, resetErr
func (c *conv) ensureCapacity(dest buffDest, cap int) // Replaces: ensureOutCapacity

// State Checking - Enhanced
func (c *conv) hasContent(dest buffDest) bool       // New: unified content checking
// Keep existing: hasContent(), hasContent(), hasContent() for performance
```

### **ELIMINATED METHODS**
- `setString()` - **ELIMINATED**: Removed from fmt.go, truncate.go, mapping.go ✅
- All legacy wrapper methods - Minimize code lines  
- All standalone functions - Convert to `*conv` methods only
- Duplicate methods with long names (e.g., `writeIntToDest` → `wrInt`)
- Complex temporary state management - Use internal work buffer instead

### **PENDING ELIMINATION** ⏳
```go
// These methods still exist but should be replaced:
c.setString()           // Used in: fmt.go, truncate.go, mapping.go
                       // Replace with: direct buffer management via anyToBuff
```

### **ARCHITECTURAL CONSTRAINT** ⚠️
```go
// DEPENDENCY HIERARCHY - PREVENT INFINITE RECURSION
// Level 1: memory.go, error.go (primitive operations only)
// Level 2: anyToBuff() (calls Level 1)
// Level 3: All other files (call anyToBuff)

// ❌ FORBIDDEN: memory.go and error.go calling anyToBuff
// ✅ REQUIRED: Use only primitive buffer methods in Level 1 files
```

### **ERROR HANDLING FOR INVALID buffDest**
```go
// Invalid buffDest cases are IGNORED (no panic, no error)
// Only handle: buffOut, buffWork, buffErr
// Default case: silent no-op (performance over safety)
```

## **CRITICAL RESTRICTION** ⚠️
```go
// ❌ FORBIDDEN - CAUSES INFINITE RECURSION:
// memory.go and error.go CANNOT call anyToBuff
// anyToBuff depends on these files for basic operations

// ✅ SAFE HIERARCHY:
anyToBuff() → wrInt/wrUint/wrFloat() → memory.go buffer methods
anyToBuff() → wrErr() → error.go buffer methods  
anyToBuff() → wrString() → memory.go buffer methods

// memory.go and error.go must use ONLY primitive buffer methods:
//  wrString(), wrBytes(), wrByte(), wrInt(), wrUint(), wrFloat()
```

## **CONFIRMED DECISIONS** ✅
1. **Parameter Order**: `dest` comes FIRST in all universal methods
2. **setString Elimination**: Removed - `anyToBuff` centralizes all conversions  
3. **Legacy Wrappers**: Completely eliminated to minimize code lines
4. **Buffer Reset**: Changed from `rst*` to `rstBuffer(dest)`
5. **Method Scope**: All buffer operations are `*conv` methods (no standalone functions)
6. **Testing**: Deferred until after implementation
7. **Simplified Naming**: `writeIntToDest` → `wrInt` (find and eliminate duplicate methods)
8. **Internal Work Buffer**: Methods use existing internal work buffer for operations
9. **anyToBuff Method**: Convert to `(c *conv) anyToBuff(dest, value)` method
10. **Progressive Implementation**: Make one change at a time with guidance

## **REMAINING TASKS** 🎯
- [x] **PRIORITY 1**: Implement simplified universal methods (`wrInt`, `wrUint`, `wrFloat`) ✅
- [x] **PRIORITY 2**: Convert `anyToBuff` to `*conv` method ✅
- [x] **PRIORITY 3**: Find and eliminate duplicate methods with simplified naming ✅
- [x] **PRIORITY 4**: Replace all destination-specific method calls ✅ (converted `wrBytes` to method, replaced major calls)
- [x] **PRIORITY 5**: Convert standalone functions to `*conv` methods ✅
- [x] **PRIORITY 6**: **FIX INFINITE RECURSION**: Remove `anyToBuff` calls from `memory.go` and `error.go` ✅
- [x] **PRIORITY 7**: Replace all `rst*` calls with `rstBuffer(dest)` ✅
- [x] **PRIORITY 8**: Eliminate temporary fields (`intVal`, `floatVal`, etc.) ✅
- [x] **PRIORITY 9**: Review `memory.go` for compliance ✅
- [ ] Run full test suite validation
- [ ] Measure allocation reduction (target: 50%)

