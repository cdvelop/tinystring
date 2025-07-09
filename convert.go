package tinystring

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

	// ✅ UNIFIED BUFFER ARCHITECTURE - Only essential fields remain
	ptrValue any // ✅ Universal pointer for complex types ([]string, map[string]any, etc.)
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
		// - Setting c.Kind and c.ptrValue for all types
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
		c.ptrValue = v
		c.Kind = KMap
	case map[string]any:
		c.ptrValue = v
		c.Kind = KMap
	case map[string]int:
		c.ptrValue = v
		c.Kind = KMap
	case map[int]string:
		c.ptrValue = v
		c.Kind = KMap

	// KPointer - Pointers to supported types
	case *string:
		// String pointer - verify not nil before dereferencing
		if v == nil {
			c.wrErr(D.String, D.Empty)
			return
		}
		// Store content relationship
		c.Kind = KPointer // Correctly set Kind to KPointer for *string
		c.ptrValue = v    // Store the pointer itself for Apply()
		c.wrString(dest, *v)
	case *int:
		c.Kind = KPointer
		c.ptrValue = v
		if v != nil {
			c.wrIntBase(dest, int64(*v), 10, true)
		}
	case *bool:
		c.Kind = KPointer
		c.ptrValue = v
		if v != nil {
			c.wrBool(dest, *v)
		}
	case *float64:
		c.Kind = KPointer
		c.ptrValue = v
		if v != nil {
			c.wrFloat64(dest, *v)
		}
	case **string:
		// Double pointer to string - pointer to pointer
		c.Kind = KPointer
		c.ptrValue = v
		if v != nil && *v != nil {
			c.wrString(dest, **v)
		}

	// KSlice - All basic slices as specified in README.md
	case []bool:
		c.ptrValue = v
		c.Kind = KSlice
	case []int:
		c.ptrValue = v
		c.Kind = KSlice
	case []int8:
		c.ptrValue = v
		c.Kind = KSlice
	case []int16:
		c.ptrValue = v
		c.Kind = KSlice
	case []int32:
		c.ptrValue = v
		c.Kind = KSlice
	case []int64:
		c.ptrValue = v
		c.Kind = KSlice
	case []uint:
		c.ptrValue = v
		c.Kind = KSlice
	case []byte: // []byte is same as []uint8 but commonly used
		c.ptrValue = v
		c.Kind = KByte
	case []uint16:
		c.ptrValue = v
		c.Kind = KSlice
	case []uint32:
		c.ptrValue = v
		c.Kind = KSlice
	case []uint64:
		c.ptrValue = v
		c.Kind = KSlice
	case []float32:
		c.ptrValue = v
		c.Kind = KSlice
	case []float64:
		c.ptrValue = v
		c.Kind = KSlice

	// KString
	case string:
		c.Kind = KString
		c.wrString(dest, v)

	// KSliceStr - Special case for []string
	case []string:
		c.ptrValue = v
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
			c.ptrValue = value
			// For structs, we don't convert to string immediately
			// The actual field access is handled by tinyreflect
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
	if t.Kind == KPointer && t.ptrValue != nil {
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

// =============================================================================
// REFLECTION FUNCTIONALITY - CORE LOGIC MOVED FROM TINYREFLECT
// =============================================================================

// GetValue returns the interface{} value stored in ptrValue
// This allows external packages to access complex types managed by tinystring
func (c *conv) GetValue() any {
	return c.ptrValue
}

// GetFieldCount returns the number of fields in a struct
// Only works with KStruct types
func (c *conv) GetFieldCount() int {
	if c.Kind != KStruct || c.ptrValue == nil {
		return 0
	}

	// Use reflection-like interface detection for struct field counting
	// This is a minimal implementation - structs are complex and need careful handling

	// For now, return 0 as we need to implement proper struct introspection
	// TODO: Implement struct field enumeration
	return 0
}

// GetField returns the value of a struct field by index
// Only works with KStruct types
func (c *conv) GetField(index int) any {
	if c.Kind != KStruct || c.ptrValue == nil {
		return nil
	}

	// TODO: Implement struct field access by index
	return nil
}

// GetTag returns the tag of a struct field by index
// Only works with KStruct types
func (c *conv) GetTag(index int, tagName string) string {
	if c.Kind != KStruct || c.ptrValue == nil {
		return ""
	}

	// TODO: Implement struct tag reading
	return ""
}

// GetMapKeys returns all keys from a map
// Only works with KMap types
func (c *conv) GetMapKeys() []any {
	if c.Kind != KMap || c.ptrValue == nil {
		return nil
	}

	var keys []any

	switch m := c.ptrValue.(type) {
	case map[string]string:
		for k := range m {
			keys = append(keys, k)
		}
	case map[string]any:
		for k := range m {
			keys = append(keys, k)
		}
	case map[string]int:
		for k := range m {
			keys = append(keys, k)
		}
	case map[int]string:
		for k := range m {
			keys = append(keys, k)
		}
	}

	return keys
}

// GetMapValue returns a value from a map by key
// Only works with KMap types
func (c *conv) GetMapValue(key any) any {
	if c.Kind != KMap || c.ptrValue == nil {
		return nil
	}

	switch m := c.ptrValue.(type) {
	case map[string]string:
		if k, ok := key.(string); ok {
			return m[k]
		}
	case map[string]any:
		if k, ok := key.(string); ok {
			return m[k]
		}
	case map[string]int:
		if k, ok := key.(string); ok {
			return m[k]
		}
	case map[int]string:
		if k, ok := key.(int); ok {
			return m[k]
		}
	}

	return nil
}

// GetSliceLen returns the length of a slice
// Works with KSlice, KSliceStr, and KByte types
func (c *conv) GetSliceLen() int {
	if c.ptrValue == nil {
		return 0
	}

	switch c.Kind {
	case KSliceStr:
		if s, ok := c.ptrValue.([]string); ok {
			return len(s)
		}
	case KByte:
		if s, ok := c.ptrValue.([]byte); ok {
			return len(s)
		}
	case KSlice:
		switch s := c.ptrValue.(type) {
		case []bool:
			return len(s)
		case []int:
			return len(s)
		case []int8:
			return len(s)
		case []int16:
			return len(s)
		case []int32:
			return len(s)
		case []int64:
			return len(s)
		case []uint:
			return len(s)
		case []uint16:
			return len(s)
		case []uint32:
			return len(s)
		case []uint64:
			return len(s)
		case []float32:
			return len(s)
		case []float64:
			return len(s)
		}
	}

	return 0
}

// GetSliceElement returns an element from a slice by index
// Works with KSlice, KSliceStr, and KByte types
func (c *conv) GetSliceElement(index int) any {
	if c.ptrValue == nil {
		return nil
	}

	switch c.Kind {
	case KSliceStr:
		if s, ok := c.ptrValue.([]string); ok && index >= 0 && index < len(s) {
			return s[index]
		}
	case KByte:
		if s, ok := c.ptrValue.([]byte); ok && index >= 0 && index < len(s) {
			return s[index]
		}
	case KSlice:
		switch s := c.ptrValue.(type) {
		case []bool:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []int:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []int8:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []int16:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []int32:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []int64:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []uint:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []uint16:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []uint32:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []uint64:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []float32:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		case []float64:
			if index >= 0 && index < len(s) {
				return s[index]
			}
		}
	}

	return nil
}

// GetSliceElementKind returns the Kind of elements in a slice
// Uses fallback logic for TinyGo/WebAssembly compatibility
func (c *conv) GetSliceElementKind() Kind {
	if c.ptrValue == nil {
		return KInvalid
	}

	switch c.Kind {
	case KSliceStr:
		return KString
	case KByte:
		return KUint8
	case KSlice:
		switch c.ptrValue.(type) {
		case []bool:
			return KBool
		case []int:
			return KInt
		case []int8:
			return KInt8
		case []int16:
			return KInt16
		case []int32:
			return KInt32
		case []int64:
			return KInt64
		case []uint:
			return KUint
		case []uint16:
			return KUint16
		case []uint32:
			return KUint32
		case []uint64:
			return KUint64
		case []float32:
			return KFloat32
		case []float64:
			return KFloat64
		}
	}

	return KInvalid
}

// GetPointerElement returns the element pointed to by a pointer
// Only works with KPointer types
func (c *conv) GetPointerElement() any {
	if c.Kind != KPointer || c.ptrValue == nil {
		return nil
	}

	switch p := c.ptrValue.(type) {
	case *string:
		if p != nil {
			return *p
		}
	case *int:
		if p != nil {
			return *p
		}
	case *bool:
		if p != nil {
			return *p
		}
	case *float64:
		if p != nil {
			return *p
		}
	case **string:
		if p != nil && *p != nil {
			return **p
		}
	}

	return nil
}

// GetPointerElementKind returns the Kind of the element pointed to by a pointer
func (c *conv) GetPointerElementKind() Kind {
	if c.Kind != KPointer || c.ptrValue == nil {
		return KInvalid
	}

	switch c.ptrValue.(type) {
	case *string:
		return KString
	case *int:
		return KInt
	case *bool:
		return KBool
	case *float64:
		return KFloat64
	case **string:
		return KPointer // Pointer to pointer
	}

	return KInvalid
}

// looksLikeStruct is a helper function to identify struct types
// This replaces the previous isStruct method with better logic
func (c *conv) looksLikeStruct(value any) bool {
	if value == nil {
		return false
	}

	// Explicitly reject known non-struct types
	switch value.(type) {
	// Basic types - explicitly not structs
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return false
	// Slice types - explicitly not structs
	case []string, []bool, []int, []int8, []int16, []int32, []int64, []uint, []uint16, []uint32, []uint64, []float32, []float64, []byte:
		return false
	// Map types - explicitly not structs
	case map[string]string, map[string]any, map[string]int, map[int]string:
		return false
	// Pointer types - explicitly not structs
	case *string, *int, *bool, *float64, **string:
		return false
	// Error type - explicitly not struct
	case error:
		return false
	default:
		// For anything else that's not explicitly known, we conservatively assume it could be a struct
		// This allows structs to be processed, while unknown types will still get an error during reflection
		return true
	}
}
