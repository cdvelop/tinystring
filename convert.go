package tinystring

import "unsafe"

// Buffer destination selection for anyToBuff universal conversion function
type buffDest int

const (
	buffOut  buffDest = iota // Primary output buffer
	buffWork                 // Working/temporary buffer
	buffErr                  // Error message buffer
)

type conv struct {
	// Buffers with initial capacity 64, grow as needed (no truncation)
	out     []byte // Buffer principal - make([]byte, 0, 64)
	outLen  int    // Longitud actual en out
	work    []byte // Buffer temporal - make([]byte, 0, 64)
	workLen int    // Longitud actual en work
	err     []byte // Buffer de errores - make([]byte, 0, 64)
	errLen  int    // Longitud actual en err
	// Type indicator - most frequently accessed	// Type indicator - most frequently accessed
	Kind Kind // Hot path: type checking

	// ✅ OPTIMIZED MEMORY ARCHITECTURE - unsafe.Pointer for complex types
	dataPtr unsafe.Pointer // Direct unsafe pointer to data (replaces ptrValue)
}

// Convert initializes a new conv struct with optional value for string,bool and number manipulation.
// REFACTORED: Now accepts variadic parameters - Convert() or Convert(value)
// Phase 7: Uses object pool internally for memory optimization (transparent to user)
func Convert(v ...any) *conv {
	c := getConv()
	c.resetAllBuffers() // Asegurar que el objeto conv esté completamente limpio
	// Validation: Only accept 0 or 1 parameter
	if len(v) > 1 {
		return c.wrErr(D.Invalid, D.Number, D.Of, D.Argument)
	}
	// Initialize with value if provided, empty otherwise
	if len(v) == 1 {
		val := v[0]
		if val == nil {
			return c.wrErr(D.String, D.Empty)
		}

		// Special case: error type should return immediately with error state
		if _, isError := val.(error); isError {
			return c.wrErr(val.(error).Error())
		}

		// Use anyToBuff for ALL other conversions - eliminates all duplication
		c.rstBuffer(buffOut)
		c.anyToBuff(buffOut, val)

		// anyToBuff handles everything:
		// - Setting c.Kind and c.dataPtr for all types
		// - String pointer handling (*string)
		// - Complex types ([]string, map, etc.) with deferred conversion
		// - All numeric and boolean type conversions
		// - Error handling for unsupported types
	}
	// If no value provided, conv is ready for builder pattern
	return c
}

// =============================================================================
// UNIVERSAL CONVERSION FUNCTION - REUSES EXISTING IMPLEMENTATIONS
// =============================================================================

// anyToBuff converts any supported type to buffer using existing conversion logic
// REUSES: floatToOut, wrStringToOut, wrStringToErr
// Supports: string, int variants, uint variants, float variants, bool, []byte, LocStr
func (c *conv) anyToBuff(dest buffDest, value any) {
	// Limpiar buffer de error antes de cualquier conversión inmediata
	c.rstBuffer(buffErr)

	switch v := value.(type) {
	// IMMEDIATE CONVERSION - Simple Types (ordered as in kind.go)

	// KBool
	case bool:
		c.Kind = KBool
		c.wrBool(dest, v)

	// KFloat32
	case float32:
		c.Kind = KFloat32
		c.wrFloat32(dest, v)

	// KFloat64
	case float64:
		c.Kind = KFloat64
		c.wrFloat64(dest, v)

	// KInt
	case int:
		c.Kind = KInt
		c.wrIntBase(dest, int64(v), 10, true)

	// KInt16
	case int16:
		c.Kind = KInt16
		c.wrIntBase(dest, int64(v), 10, true)

	// KInt32
	case int32:
		c.Kind = KInt32
		c.wrIntBase(dest, int64(v), 10, true)

	// KInt64
	case int64:
		c.Kind = KInt64
		c.wrIntBase(dest, v, 10, true)

	// KInt8
	case int8:
		c.Kind = KInt8
		c.wrIntBase(dest, int64(v), 10, true)

	// KMap - Maps with supported types
	case map[string]string:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KMap
	case map[string]any:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KMap
	case map[string]int:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KMap
	case map[int]string:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KMap

	// KPointer - Pointers to supported types
	case *string:
		// String pointer - verify not nil before dereferencing
		if v == nil {
			c.wrErr(D.String, D.Empty)
			return
		}
		// Store content relationship
		c.Kind = KPointer             // Correctly set Kind to KPointer for *string
		c.dataPtr = unsafe.Pointer(v) // Store the pointer itself for Apply()
		c.wrString(dest, *v)
	case *int:
		c.Kind = KPointer
		c.dataPtr = unsafe.Pointer(v)
		if v != nil {
			c.wrIntBase(dest, int64(*v), 10, true)
		}
	case *bool:
		c.Kind = KPointer
		c.dataPtr = unsafe.Pointer(v)
		if v != nil {
			c.wrBool(dest, *v)
		}
	case *float64:
		c.Kind = KPointer
		c.dataPtr = unsafe.Pointer(v)
		if v != nil {
			c.wrFloat64(dest, *v)
		}
	case **string:
		// Double pointer to string - pointer to pointer
		c.Kind = KPointer
		c.dataPtr = unsafe.Pointer(v)
		if v != nil && *v != nil {
			c.wrString(dest, **v)
		}

	// KSlice - All basic slices as specified in README.md
	case []bool:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []int:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []int8:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []int16:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []int32:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []int64:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []uint:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []byte: // []byte is same as []uint8 but commonly used
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KByte
	case []uint16:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []uint32:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []uint64:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []float32:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice
	case []float64:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSlice

	// KString
	case string:
		c.Kind = KString
		c.wrString(dest, v)

	// KSliceStr - Special case for []string
	case []string:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = KSliceStr

	// KUint
	case uint:
		c.Kind = KUint
		c.wrIntBase(dest, int64(v), 10, false)

	// KUint16
	case uint16:
		c.Kind = KUint16
		c.wrIntBase(dest, int64(v), 10, false)

	// KUint32
	case uint32:
		c.Kind = KUint32
		c.wrIntBase(dest, int64(v), 10, false)

	// KUint64
	case uint64:
		c.Kind = KUint64
		c.wrIntBase(dest, int64(v), 10, false)

	// KUint8
	case uint8:
		c.Kind = KUint8
		c.wrIntBase(dest, int64(v), 10, false)

	// Special cases
	case error:
		c.wrErr(v.Error())
		return // Early return since this sets error state

	default:
		// For complex types, use a more sophisticated check
		// Check if it's a struct by examining its characteristics
		if c.looksLikeStruct(value) {
			c.Kind = KStruct
			c.dataPtr = unsafe.Pointer(&value)
			// For structs, we don't convert to string immediately
			// The actual field access is handled by reflection methods
			return
		}

		// Unknown/unsupported type - write error using DICTIONARY (REUSE existing wrErr)
		c.wrErr(D.Type, D.Not, D.Supported)
	}
}

// GetKind returns the Kind of the value stored in the conv
// This allows external packages to reuse tinystring's type detection logic
func (c *conv) GetKind() Kind {
	return c.Kind
}

// Apply updates the original string pointer with the current content and auto-releases to pool.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *conv) Apply() {
	if t.Kind == KPointer && t.dataPtr != nil {
		// Type assert to *string for Apply() functionality using unsafe pointer
		if strPtr := (*string)(t.dataPtr); strPtr != nil {
			*strPtr = t.getBuffString()
		}
	}
	// Auto-release back to pool for memory efficiency
	t.putConv()
}

// String method to return the content of the conv and automatically returns object to pool
// Phase 7: Auto-release makes pool usage completely transparent to user
func (c *conv) String() string {
	// If there's an error, return empty string (error available via StringErr())
	if c.hasContent(buffErr) {
		c.putConv() // Auto-release back to pool for memory efficiency
		return ""
	}

	out := c.getBuffString()
	// Auto-release back to pool for memory efficiency
	c.putConv()
	return out
}
