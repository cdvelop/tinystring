package tinystring

import "unsafe"

// Buffer destination selection for anyToBuff universal conversion function
type buffDest int

const (
	buffOut  buffDest = iota // Primary output buffer
	buffWork                 // Working/temporary buffer
	buffErr                  // Error message buffer
)

type Conv struct {
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

// Convert initializes a new Conv struct with optional value for string,bool and number manipulation.
// REFACTORED: Now accepts variadic parameters - Convert() or Convert(value)
// Phase 7: Uses object pool internally for memory optimization (transparent to user)
func Convert(v ...any) *Conv {
	c := getConv()
	c.resetAllBuffers() // Asegurar que el objeto Conv esté completamente limpio
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
	// If no value provided, Conv is ready for builder pattern
	return c
}

// =============================================================================
// UNIVERSAL CONVERSION FUNCTION - REUSES EXISTING IMPLEMENTATIONS

// =============================================================================

// anyToBuff converts any supported type to buffer using existing conversion logic
// REUSES: floatToOut, wrStringToOut, wrStringToErr
// Supports: string, int variants, uint variants, float variants, bool, []byte, LocStr
func (c *Conv) anyToBuff(dest buffDest, value any) {
	// Limpiar buffer de error antes de cualquier conversión inmediata
	c.rstBuffer(buffErr)

	switch v := value.(type) {
	// IMMEDIATE CONVERSION - Simple Types (ordered as in Kind.go)

	// K.Bool
	case bool:
		c.Kind = K.Bool
		c.wrBool(dest, v)

	// K.Float32
	case float32:
		c.Kind = K.Float32
		c.wrFloat32(dest, v)

	// K.Float64
	case float64:
		c.Kind = K.Float64
		c.wrFloat64(dest, v)

	// K.Int
	case int:
		c.Kind = K.Int
		c.wrIntBase(dest, int64(v), 10, true)

	// K.Int8
	case int8:
		c.Kind = K.Int8
		c.wrIntBase(dest, int64(v), 10, true)

	// K.Int16
	case int16:
		c.Kind = K.Int16
		c.wrIntBase(dest, int64(v), 10, true)

	// K.Int32
	case int32:
		c.Kind = K.Int32
		c.wrIntBase(dest, int64(v), 10, true)

	// K.Int64
	case int64:
		c.Kind = K.Int64
		c.wrIntBase(dest, v, 10, true)

	// K.Pointer - Only *string pointer supported
	case *string:
		// String pointer - verify not nil before dereferencing
		if v == nil {
			c.wrErr(D.String, D.Empty)
			return
		}
		// Store content relationship
		c.Kind = K.Pointer            // Correctly set Kind to K.Pointer for *string
		c.dataPtr = unsafe.Pointer(v) // Store the pointer itself for Apply()
		c.wrString(dest, *v)

	// K.String
	case string:
		c.Kind = K.String
		c.wrString(dest, v)

	// K.SliceStr - Special case for []string
	case []string:
		c.dataPtr = unsafe.Pointer(&v)
		c.Kind = K.Slice

	// K.Uint
	case uint:
		c.Kind = K.Uint
		c.wrIntBase(dest, int64(v), 10, false)

	// K.Uint8
	case uint8:
		c.Kind = K.Uint8
		c.wrIntBase(dest, int64(v), 10, false)

	// K.Uint16
	case uint16:
		c.Kind = K.Uint16
		c.wrIntBase(dest, int64(v), 10, false)

	// K.Uint32
	case uint32:
		c.Kind = K.Uint32
		c.wrIntBase(dest, int64(v), 10, false)

	// K.Uint64
	case uint64:
		c.Kind = K.Uint64
		c.wrIntBase(dest, int64(v), 10, false)

	// Special cases
	case error:
		c.wrErr(v.Error())

	default:
		// Unknown/unsupported type - write error using DICTIONARY (REUSE existing wrErr)
		c.wrErr(D.Type, D.Not, D.Supported)
	}
}

// GetKind returns the Kind of the value stored in the Conv
// This allows external packages to reuse tinystring's type detection logic
func (c *Conv) GetKind() Kind {
	return c.Kind
}

// Apply updates the original string pointer with the current content and auto-releases to pool.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *Conv) Apply() {
	if t.Kind == K.Pointer && t.dataPtr != nil {
		// Type assert to *string for Apply() functionality using unsafe pointer
		if strPtr := (*string)(t.dataPtr); strPtr != nil {
			*strPtr = t.getString(buffOut)
		}
	}
	// Auto-release back to pool for memory efficiency
	t.putConv()
}

// String method to return the content of the Conv and automatically returns object to pool
// Phase 7: Auto-release makes pool usage completely transparent to user
func (c *Conv) String() string {
	// If there's an error, return empty string (error available via StringErr())
	if c.hasContent(buffErr) {
		c.putConv() // Auto-release back to pool for memory efficiency
		return ""
	}

	out := c.getString(buffOut)
	// Auto-release back to pool for memory efficiency
	c.putConv()
	return out
}

// Bytes returns the content of the Conv as a byte slice
func (c *Conv) Bytes() []byte {
	return c.getBytes(buffOut)
}
