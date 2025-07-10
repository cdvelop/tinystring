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
	kind kind // Hot path: type checking

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
		// - Setting c.kind and c.dataPtr for all types
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

	// Kind.Bool
	case bool:
		c.kind = Kind.Bool
		c.wrBool(dest, v)

	// Kind.Float32
	case float32:
		c.kind = Kind.Float32
		c.wrFloat32(dest, v)

	// Kind.Float64
	case float64:
		c.kind = Kind.Float64
		c.wrFloat64(dest, v)

	// Kind.Int
	case int:
		c.kind = Kind.Int
		c.wrIntBase(dest, int64(v), 10, true)

	// Kind.Int8
	case int8:
		c.kind = Kind.Int8
		c.wrIntBase(dest, int64(v), 10, true)

	// Kind.Int16
	case int16:
		c.kind = Kind.Int16
		c.wrIntBase(dest, int64(v), 10, true)

	// Kind.Int32
	case int32:
		c.kind = Kind.Int32
		c.wrIntBase(dest, int64(v), 10, true)

	// Kind.Int64
	case int64:
		c.kind = Kind.Int64
		c.wrIntBase(dest, v, 10, true)

	// Kind.Pointer - Only *string pointer supported
	case *string:
		// String pointer - verify not nil before dereferencing
		if v == nil {
			c.wrErr(D.String, D.Empty)
			return
		}
		// Store content relationship
		c.kind = Kind.Pointer         // Correctly set Kind to Kind.Pointer for *string
		c.dataPtr = unsafe.Pointer(v) // Store the pointer itself for Apply()
		c.wrString(dest, *v)

	// Kind.String
	case string:
		c.kind = Kind.String
		c.wrString(dest, v)

	// Kind.SliceStr - Special case for []string
	case []string:
		c.dataPtr = unsafe.Pointer(&v)
		c.kind = Kind.SliceStr

	// Kind.Uint
	case uint:
		c.kind = Kind.Uint
		c.wrIntBase(dest, int64(v), 10, false)

	// Kind.Uint16
	case uint16:
		c.kind = Kind.Uint16
		c.wrIntBase(dest, int64(v), 10, false)

	// Kind.Uint32
	case uint32:
		c.kind = Kind.Uint32
		c.wrIntBase(dest, int64(v), 10, false)

	// Kind.Uint64
	case uint64:
		c.kind = Kind.Uint64
		c.wrIntBase(dest, int64(v), 10, false)

	// Kind.Uint8
	case uint8:
		c.kind = Kind.Uint8
		c.wrIntBase(dest, int64(v), 10, false)

	// Special cases
	case error:
		c.wrErr(v.Error())

	default:
		// Unknown/unsupported type - write error using DICTIONARY (REUSE existing wrErr)
		c.wrErr(D.Type, D.Not, D.Supported)
	}
}

// GetKind returns the Kind of the value stored in the conv
// This allows external packages to reuse tinystring's type detection logic
func (c *conv) GetKind() kind {
	return c.kind
}

// Apply updates the original string pointer with the current content and auto-releases to pool.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *conv) Apply() {
	if t.kind == Kind.Pointer && t.dataPtr != nil {
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
