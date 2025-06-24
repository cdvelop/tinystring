package tinystring

// Buffer destination selection for anyToBuff universal conversion function
type buffDest int

const (
	buffOut  buffDest = iota // Primary output buffer
	buffWork                 // Working/temporary buffer
	buffErr                  // Error message buffer
)

type conv struct {
	// Phase 13.3: Optimized memory allocation - all buffers are dynamic []byte
	// Buffers with initial capacity 64, grow as needed (no truncation)
	out     []byte // Buffer principal - make([]byte, 0, 64)
	outLen  int    // Longitud actual en out
	work    []byte // Buffer temporal - make([]byte, 0, 64)
	workLen int    // Longitud actual en work
	err     []byte // Buffer de errores - make([]byte, 0, 64)
	errLen  int    // Longitud actual en err
	// Type indicator - most frequently accessed	// Type indicator - most frequently accessed
	kind kind // Hot path: type checking

	// ✅ UNIFIED BUFFER ARCHITECTURE - Only essential fields remain
	ptrValue any // ✅ Universal pointer for complex types ([]string, map[string]any, etc.)
}

// Convert initializes a new conv struct with optional value for string,bool and number manipulation.
// REFACTORED: Now accepts variadic parameters - Convert() or Convert(value)
// Phase 7: Uses object pool internally for memory optimization (transparent to user)
func Convert(v ...any) *conv {
	c := getConv()
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
		// - Setting c.kind and c.ptrValue for all types
		// - String pointer handling (*string)
		// - Complex types ([]string, map, etc.) with lazy conversion
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
// REUSES: fmtIntToOut, floatToOut, wrStringToOut, wrStringToErr
// Supports: string, int variants, uint variants, float variants, bool, []byte, LocStr
func (c *conv) anyToBuff(dest buffDest, value any) {
	switch v := value.(type) {
	// IMMEDIATE CONVERSION - Simple Types (REUSE existing implementations)
	case string:
		c.kind = KString
		c.wrString(dest, v)
	case *string:
		// String pointer - verify not nil before dereferencing
		if v == nil {
			c.wrErr(D.String, D.Empty)
			return
		}
		// Store content and maintain pointer relationship
		c.kind = KPointer
		c.wrString(dest, *v)
		c.ptrValue = v
	case error:
		// Error type - write error message
		c.wrErr(v.Error())
		return // Early return since this sets error state
	case int:
		c.kind = KInt
		c.wrInt(dest, int64(v))
	case int8:
		c.kind = KInt
		c.wrInt(dest, int64(v))
	case int16:
		c.kind = KInt
		c.wrInt(dest, int64(v))
	case int32:
		c.kind = KInt
		c.wrInt(dest, int64(v))
	case int64:
		c.kind = KInt
		c.wrInt(dest, v)
	case uint:
		c.kind = KUint
		c.wrUint(dest, uint64(v))
	case uint8:
		c.kind = KUint
		c.wrUint(dest, uint64(v))
	case uint16:
		c.kind = KUint
		c.wrUint(dest, uint64(v))
	case uint32:
		c.kind = KUint
		c.wrUint(dest, uint64(v))
	case uint64:
		c.kind = KUint
		c.wrUint(dest, v)
	case float32:
		c.kind = KFloat32
		c.wrFloat(dest, float64(v))
	case float64:
		c.kind = KFloat64
		c.wrFloat(dest, v)

	case bool:
		c.kind = KBool
		c.wrBool(dest, v)

	case []byte:
		c.wrBytes(dest, v)
		c.kind = KString // Treat []byte as string

	// LAZY CONVERSION - Complex Types (store pointer, convert on demand)
	case []string:
		c.ptrValue = v
		c.kind = KSliceStr
		// No immediate conversion - wait for operation
	case map[string]string:
		c.ptrValue = v
		c.kind = KMap
		// No immediate conversion - wait for operation
	case map[string]any:
		c.ptrValue = v
		c.kind = KMap
		// No immediate conversion - wait for operation
	default:
		// Unknown type - write error using DICTIONARY (REUSE existing wrErr)
		c.wrErr(D.Type, D.Not, D.Supported)
	}
}

// Apply updates the original string pointer with the current content and auto-releases to pool.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *conv) Apply() {
	if t.kind == KPointer && t.ptrValue != nil {
		// Type assert to *string for Apply() functionality
		if strPtr, ok := t.ptrValue.(*string); ok && strPtr != nil {
			*strPtr = t.getBuffString()
		}
	}
	// Auto-release back to pool for memory efficiency
	t.putConv()
}

// String method to return the content of the conv and automatically returns object to pool
// Phase 7: Auto-release makes pool usage completely transparent to user
func (c *conv) String() string {
	// If there's an error, return empty string (error available via StringError())
	if c.hasContent(buffErr) {
		c.putConv() // Auto-release back to pool for memory efficiency
		return ""
	}

	out := c.getBuffString()
	// Auto-release back to pool for memory efficiency
	c.putConv()
	return out
}
