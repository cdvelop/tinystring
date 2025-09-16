# TinyString - Unified Buffer Architecture 🎯

## **STATUS: 95% COMPLETE** ✅
- ✅ **Zero-allocation buffer architecture** implemented
- ✅ **Universal conversion**: `anyToBuff(c *Conv, dest BuffDest, value any)`
- ✅ **Non-recursive errors**: `wrErr()` with language support
- ✅ **100% Buffer API compliance**: No manual buffer access

## **CORE RULES** 🚨
```go
// ❌ FORBIDDEN:
c.out = c.out[:0]         // Manual buffer manipulation
len(c.err) > 0            // Direct length checks


// ✅ MANDATORY:
c.ResetBuffer(dest BuffDest)
c.hasContent(dest BuffDest)
c.WrString(dest BuffDest, s)

```

## **CORE ARCHITECTURE** ️
```go
// Buffer Management
type Conv struct {
    out     []byte // Primary output
    outLen  int    // Length tracking
    work    []byte // Working buffer  
    workLen int    // Work tracking
    err     []byte // Error buffer
    errLen  int    // Error tracking
    Kind    Kind   // Type indicator
    ptrValue any   // Universal storage
}

// Universal Conversion
func (c *Conv) anyToBuff(dest BuffDest, value any) {
    // dest: BuffOut | BuffWork | BuffErr
    // Writes errors via c.wrErr(), no return values
    // NOW A *Conv METHOD (not standalone function)
}

// Error System
func (c *Conv) wrErr(msgs ...any) *Conv {
    // Language-aware, dictionary-based, buffer API only
}
```

## **BUFFER API** 📋
```go
// State Checking
c.hasContent(BuffDest) // Check out,work, err buffer
// Writing
c.WrString(BuffDest, s) // Write string to specified buffer
c.wrBytes(BuffDest, data) // Write bytes to specified buffer    
c.wrByte(BuffDest,s)        
// Reading
c.GetString(BuffDest)        // Read out, work, err buffer

```

## **CONSTRAINTS** 📝
- **WebAssembly-first**: Binary size over performance
- **No stdlib**: Manual implementations only
- **Dictionary errors**: Use D.* constants only
- **TinyGo compatible**: Limited reflection

## **CENTRALIZED BUFFER API - UNIVERSAL METHODS** 🎯

### **NAMING CONVENTION CHANGE**
- All `write*` methods become `wr*`
- All methods receive `BuffDest` parameter for destination selection
- Universal methods replace destination-specific variants
- NO backwards compatibility - complete refactoring

### **UNIVERSAL BUFFER METHODS** (SIMPLIFIED NAMING)
```go
// Writing Operations (dest FIRST parameter) - SIMPLIFIED NAMES
func (c *Conv) WrString(dest BuffDest, s string)    // Replaces: wrStringToOut, WrString, wrStringToErr  
func (c *Conv) wrByte(dest BuffDest, b byte)        // Replaces: wrByte (out-only)
func (c *Conv) wrBytes(dest BuffDest, data []byte)  // Replaces: wrBytes, wrToWork, wrBytes
func (c *Conv) wrInt(dest BuffDest, val int64)      // Replaces: writeIntToDest + duplicates
func (c *Conv) wrUint(dest BuffDest, val uint64)    // Replaces: writeUintToDest + duplicates  
func (c *Conv) wrFloat(dest BuffDest, val float64)  // Replaces: writeFloatToDest + duplicates

// Universal Conversion (NOW A METHOD)
func (c *Conv) anyToBuff(dest BuffDest, value any)  // Replaces: anyToBuff function

// Reading Operations (dest FIRST parameter)
func (c *Conv) GetString(dest BuffDest) string      // Replaces: GetString, GetString, GetString

// Buffer Management (dest FIRST parameter)
func (c *Conv) ResetBuffer(dest BuffDest)             // Replaces: rstOut, rstWork, resetErr
func (c *Conv) ensureCapacity(dest BuffDest, cap int) // Replaces: ensureOutCapacity

// State Checking - Enhanced
func (c *Conv) hasContent(dest BuffDest) bool       // New: unified content checking
// Keep existing: hasContent(), hasContent(), hasContent() for performance
```

### **ELIMINATED METHODS**
- `setString()` - **ELIMINATED**: Removed from fmt.go, truncate.go, mapping.go ✅
- All legacy wrapper methods - Minimize code lines  
- All standalone functions - Convert to `*Conv` methods only
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

### **ERROR HANDLING FOR INVALID BuffDest**
```go
// Invalid BuffDest cases are IGNORED (no panic, no error)
// Only handle: BuffOut, BuffWork, BuffErr
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
anyToBuff() → WrString() → memory.go buffer methods

// memory.go and error.go must use ONLY primitive buffer methods:
//  WrString(), wrBytes(), wrByte(), wrInt(), wrUint(), wrFloat()
```

## **CONFIRMED DECISIONS** ✅
1. **Parameter Order**: `dest` comes FIRST in all universal methods
2. **setString Elimination**: Removed - `anyToBuff` centralizes all conversions  
3. **Legacy Wrappers**: Completely eliminated to minimize code lines
4. **Buffer Reset**: Changed from `rst*` to `ResetBuffer(dest)`
5. **Method Scope**: All buffer operations are `*Conv` methods (no standalone functions)
6. **Testing**: Deferred until after implementation
7. **Simplified Naming**: `writeIntToDest` → `wrInt` (find and eliminate duplicate methods)
8. **Internal Work Buffer**: Methods use existing internal work buffer for operations
9. **anyToBuff Method**: Convert to `(c *Conv) anyToBuff(dest, value)` method
10. **Progressive Implementation**: Make one change at a time with guidance

## **REMAINING TASKS** 🎯
- [x] **PRIORITY 1**: Implement simplified universal methods (`wrInt`, `wrUint`, `wrFloat`) ✅
- [x] **PRIORITY 2**: Convert `anyToBuff` to `*Conv` method ✅
- [x] **PRIORITY 3**: Find and eliminate duplicate methods with simplified naming ✅
- [x] **PRIORITY 4**: Replace all destination-specific method calls ✅ (converted `wrBytes` to method, replaced major calls)
- [x] **PRIORITY 5**: Convert standalone functions to `*Conv` methods ✅
- [x] **PRIORITY 6**: **FIX INFINITE RECURSION**: Remove `anyToBuff` calls from `memory.go` and `error.go` ✅
- [x] **PRIORITY 7**: Replace all `rst*` calls with `ResetBuffer(dest)` ✅
- [x] **PRIORITY 8**: Eliminate temporary fields (`intVal`, `floatVal`, etc.) ✅
- [x] **PRIORITY 9**: Review `memory.go` for compliance ✅
- [ ] Run full test suite validation
- [ ] Measure allocation reduction (target: 50%)

