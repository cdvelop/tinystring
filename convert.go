package tinystring

import (
	"unsafe"
)

// Import unified types from abi.go - no more duplication
// kind is now defined in abi.go with tp prefix

type conv struct {
	// PRIMARY: Reflection fields integrated from refValue
	typ  *refType       // Reflection type information
	ptr  unsafe.Pointer // Pointer to the actual data
	flag refFlag        // Reflection flags for memory layout

	// ESSENTIAL: Core operation fields only
	vTpe         kind      // Type cache for performance (redundant with flag but kept for compatibility)
	roundDown    bool      // Operation flags
	separator    string    // String operations
	tmpStr       string    // String cache for performance
	lastConvType kind      // Cache validation
	err          errorType // Error handling

	// SPECIAL CASES: Complex types that need direct storage
	stringSliceVal []string // Slice operations
	stringPtrVal   *string  // Pointer operations
}

// Functional options pattern for conv construction
type convOpt func(*conv)

// withValue initializes conv with any value type using unified reflection approach
func withValue(v any) convOpt {
	return func(c *conv) {
		if v == nil {
			c.initFromValue(nil)
			c.vTpe = tpString
			return
		}

		c.initFromValue(v)
		if !c.refIsValid() {
			c.vTpe = tpString
			return
		}

		// For Convert() function, automatically dereference pointers for convenience
		// This allows Convert(ptr).JsonEncode() to work the same as Convert(value).JsonEncode()
		originalKind := c.refKind()
		if originalKind == tpPointer {
			// Dereference the pointer and use the underlying value
			elem := c.refElem()
			if elem.refIsValid() {
				// Copy the dereferenced value to our conv
				c.typ = elem.typ
				c.ptr = elem.ptr
				c.flag = elem.flag
				c.vTpe = elem.refKind()
			} else {
				c.vTpe = tpString // Handle nil pointer
			}
		} else {
			c.vTpe = originalKind
		}

		// Handle special cases that need direct storage
		switch val := v.(type) {
		case []string:
			c.stringSliceVal = val
			c.vTpe = tpStrSlice
		case *string:
			c.stringPtrVal = val
			c.vTpe = tpStrPtr
		default:
			// All other types handled via reflection - no need for type switches
			switch c.vTpe {
			case tpStruct, tpSlice, tpArray, tpPointer:
				// Complex types - value stored in integrated reflection
			}
		}
	}
}

// newConv creates a new conv with functional options
func newConv(opts ...convOpt) *conv {
	c := &conv{
		separator: "_", // default separator
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// length returns the length of the current string
func (c *conv) length() int {
	return len(c.getString())
}

// reset clears the conv builder
func (c *conv) reset() *conv {
	c.setString("")
	return c
}

// Builder creates a new string builder instance
func Builder() *conv {
	return &conv{
		vTpe:      tpString,
		separator: "_",
	}
}

// Convert initializes a new conv struct with any type of value for string,bool and number manipulation.
// Uses the functional options pattern internally.
func Convert(v any) *conv {
	return newConv(withValue(v))
}

// Unified reflection-based value access methods
func (c *conv) getInt64() int64 {
	if c.refIsValid() && (c.vTpe >= tpInt && c.vTpe <= tpInt64) {
		return c.refInt()
	}
	return 0
}

func (c *conv) getUint64() uint64 {
	if c.refIsValid() && (c.vTpe >= tpUint && c.vTpe <= tpUintptr) {
		return c.refUint()
	}
	return 0
}

func (c *conv) getFloat64() float64 {
	if !c.refIsValid() {
		return 0
	}

	switch c.vTpe {
	case tpFloat32, tpFloat64:
		return c.refFloat()
	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		return float64(c.refInt())
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr:
		return float64(c.refUint())
	default:
		return 0
	}
}

func (c *conv) getBool() bool {
	if c.refIsValid() && c.vTpe == tpBool {
		return c.refBool()
	}
	return false
}

func (c *conv) getStringDirect() string {
	if c.refIsValid() && c.refKind() == tpString {
		return c.refString()
	}
	return ""
}

func (t *conv) separatorCase(sep ...string) string {
	t.separator = "_" // underscore default
	if len(sep) > 0 {
		t.separator = sep[0]
	}
	return t.separator
}

// Apply updates the original string pointer with the current content.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *conv) Apply() {
	if t.vTpe == tpStrPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = t.getString()
	}
}

// String method to return the content of the conv without modifying any original pointers
func (t *conv) String() string {
	return t.getString()
}

// StringError returns the content of the conv along with any error that occurred during processing
func (t *conv) StringError() (string, error) {
	if t.vTpe == tpErr {
		return t.getString(), t
	}
	return t.getString(), nil
}

// Helper function to check if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= 'À' && r <= 'ÿ' && r != 'x' && r != '÷')
}

// getString converts the current value to string only when needed
// Optimized with string caching to avoid repeated conversions using reflection-only approach
func (t *conv) getString() string {
	if t.vTpe == tpErr {
		return string(t.err)
	} // Use cached string if available and type hasn't changed
	// For TinyString special types (tpStrPtr, tpStrSlice, etc.), always use vTpe
	// For struct fields and reflection values, use refKind()
	var currentKind kind
	if t.vTpe == tpStrPtr || t.vTpe == tpStrSlice || t.vTpe == tpErr {
		// TinyString special types - use vTpe
		currentKind = t.vTpe
	} else if t.refIsValid() && t.typ != nil {
		// Reflection-based types - use refKind()
		currentKind = t.refKind()
	} else {
		// Fallback to vTpe
		currentKind = t.vTpe
	}

	if t.tmpStr != "" && t.lastConvType == currentKind {
		return t.tmpStr
	}
	// Convert to string using reflection-based methods
	switch currentKind {
	case tpString:
		t.tmpStr = t.getStringDirect()
	case tpStrPtr:
		if t.stringPtrVal != nil {
			t.tmpStr = *t.stringPtrVal
		} else {
			t.tmpStr = ""
		}
	case tpStrSlice:
		if len(t.stringSliceVal) == 0 {
			t.tmpStr = ""
		} else {
			// Join with space as default - use internal method
			t.tmpStr = t.joinSlice(" ")
		}
	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		// Use manual implementation instead of strconv
		intVal := t.getInt64()
		if intVal == 0 {
			t.tmpStr = "0"
		} else {
			t.fmtIntGeneric(intVal, 10, true)
		}
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr:
		// Use manual implementation instead of strconv
		uintVal := t.getUint64()
		if uintVal == 0 {
			t.tmpStr = "0"
		} else {
			t.fmtIntGeneric(int64(uintVal), 10, false)
		}
	case tpFloat32:
		// Use manual implementation with float32 precision
		t.f2s()
	case tpFloat64:
		// Use manual implementation for float64
		t.f2s()
	case tpBool:
		if t.getBool() {
			t.tmpStr = trueStr
		} else {
			t.tmpStr = falseStr
		}
	case tpErr:
		// For error types, return the error message
		t.tmpStr = string(t.err)
	default:
		t.tmpStr = ""
	}
	// Update cache state
	t.lastConvType = currentKind
	return t.tmpStr
}

// addRne2Buf manually encodes a rune to UTF-8 and appends it to the byte slice.
// This avoids importing the unicode/utf8 package for size optimization.
func addRne2Buf(buf []byte, r rune) []byte {
	if r < 0x80 {
		return append(buf, byte(r))
	} else if r < 0x800 {
		return append(buf, byte(0xC0|(r>>6)), byte(0x80|(r&0x3F)))
	} else if r < 0x10000 {
		return append(buf, byte(0xE0|(r>>12)), byte(0x80|((r>>6)&0x3F)), byte(0x80|(r&0x3F)))
	} else {
		return append(buf, byte(0xF0|(r>>18)), byte(0x80|((r>>12)&0x3F)), byte(0x80|((r>>6)&0x3F)), byte(0x80|(r&0x3F)))
	}
}

// setString converts to string type and stores the value using reflection
func (t *conv) setString(s string) {
	// Preserve original pointer information before reinitializing
	originalVTpe := t.vTpe
	originalStringPtrVal := t.stringPtrVal

	// Update conv to hold the new string value directly
	t.initFromValue(s)

	// If working with string pointer, restore pointer info and update the original string
	if originalVTpe == tpStrPtr && originalStringPtrVal != nil {
		t.vTpe = tpStrPtr
		t.stringPtrVal = originalStringPtrVal
		*originalStringPtrVal = s
		// Keep the vTpe as tpStrPtr to maintain the pointer relationship
	} else {
		t.vTpe = tpString
	}

	// Clear slice values to save memory - other values handled by reflection
	t.stringSliceVal = nil

	// Invalidate cache since we changed the string
	t.tmpStr = ""
	t.lastConvType = kind(0)
}

// joinSlice joins string slice with separator - optimized for minimal allocations
func (t *conv) joinSlice(separator string) string {
	if len(t.stringSliceVal) == 0 {
		return ""
	}
	if len(t.stringSliceVal) == 1 {
		return t.stringSliceVal[0]
	}

	// Calculate total length to minimize allocations
	tl := 0 // totalLen
	for _, s := range t.stringSliceVal {
		tl += len(s)
	}
	tl += len(separator) * (len(t.stringSliceVal) - 1)

	// Build result string efficiently using slice of bytes
	result := make([]byte, 0, tl) // result

	for i, s := range t.stringSliceVal {
		if i > 0 {
			result = append(result, separator...)
		}
		result = append(result, s...)
	}

	return string(result)
}

// refEface is the header for an interface{} value
type refEface struct {
	typ  *refType
	data unsafe.Pointer
}

// initFromValue initializes conv fields from any value (replaces refValueOf)
func (c *conv) initFromValue(v any) {
	if v == nil {
		c.typ = nil
		c.ptr = nil
		c.flag = 0
		c.vTpe = tpString
		return
	}

	e := (*refEface)(unsafe.Pointer(&v))
	c.typ = e.typ
	c.ptr = e.data
	c.flag = refFlag(c.typ.Kind())

	// Determine flagIndir according to type
	switch c.typ.Kind() {
	case tpBool, tpInt, tpInt8, tpInt16, tpInt32, tpInt64,
		tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr,
		tpFloat32, tpFloat64, tpPointer, tpUnsafePointer:
		// These basic types are stored directly in interface - no flagIndir
	case tpString:
		// Strings are stored directly in interface on most architectures
	default:
		// For other types (struct, slice, array, etc.), check if stored indirectly
		if ifaceIndir(c.typ) {
			c.flag |= flagIndir
		}
	}
	// Cache vTpe for compatibility
	c.vTpe = c.typ.Kind()
}

// Internal conversion methods - centralized in conv to minimize allocations
// These methods modify the conv struct directly instead of returning values

// any2s converts any type to string using reflection-only approach
// default set to "" if no conversion is possible
// supports int, uint, float, bool, string and error types
func (t *conv) any2s(v any) {
	switch val := v.(type) {
	case errorType:
		t.err = val
	case error:
		t.err = errorType(val.Error())
	case string:
		t.tmpStr = val
	case bool:
		if val {
			t.tmpStr = trueStr
		} else {
			t.tmpStr = falseStr
		}
	default:
		// Handle all other types using reflection
		t.initFromValue(v)
		if !t.refIsValid() {
			t.tmpStr = ""
			return
		}
		switch t.refKind() {
		case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
			intVal := t.refInt()
			t.fmtIntGeneric(intVal, 10, true)
		case tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr:
			uintVal := t.refUint()
			t.fmtIntGeneric(int64(uintVal), 10, false)
		case tpFloat32, tpFloat64:
			// Float value already stored in t by initFromValue
			t.f2s()
		default:
			// Fallback for unknown types
			t.tmpStr = ""
		}
	}
}

// Generic helper functions are all defined above
