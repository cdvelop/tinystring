package tinystring

// vTpe represents the type of value stored in conv
type vTpe uint8

const (
	typeStr vTpe = iota
	typeInt
	typeUint
	typeFloat
	typeBool
	typeStrSlice
	typeStrPtr
	typeErr
)

// Generic type interfaces for consolidating repetitive type switches
type anyInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type anyUint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type anyFloat interface {
	~float32 | ~float64
}

type conv struct {
	stringVal      string
	intVal         int64
	uintVal        uint64
	floatVal       float64
	boolVal        bool
	stringSliceVal []string
	stringPtrVal   *string
	vTpe           vTpe
	roundDown      bool
	separator      string
	tmpStr         string // Cache for temp string conversion to avoid repeated work
	lastConvType   vTpe   // Track last converted type for cache validation
	err            string // Error message string (changed from errorType for dictionary integration)

	// Phase 6.2: Buffer reuse optimization
	buf []byte // Reusable buffer for string operations
}

// Convert initializes a new conv struct with optional value for string,bool and number manipulation.
// REFACTORED: Now accepts variadic parameters - Convert() or Convert(value)
// Phase 7: Uses object pool internally for memory optimization (transparent to user)
func Convert(v ...any) *conv {
	c := getConv()

	// Validation: Only accept 0 or 1 parameter
	if len(v) > 1 {
		c.err = T(D.Invalid, D.Number, D.Of, D.Argument) // Consistent error handling pattern
		return c
	}

	// Initialize with value if provided, empty otherwise
	if len(v) == 1 {
		// Inlined withValue logic for performance
		val := v[0]
		if val == nil {
			c.err = Err(D.String, D.Empty).err
			c.vTpe = typeErr
		} else {
			switch typedVal := val.(type) {
			case string:
				c.stringVal = typedVal
				c.vTpe = typeStr
			case []string:
				c.stringSliceVal = typedVal
				c.vTpe = typeStrSlice
			case *string:
				c.stringVal = *typedVal
				c.stringPtrVal = typedVal
				c.vTpe = typeStrPtr
			case bool:
				c.boolVal = typedVal
				c.vTpe = typeBool
			case error:
				c.err = typedVal.Error()
				c.vTpe = typeErr
			default:
				// Handle numeric types using generics
				c.handleAnyType(typedVal)
			}
		}
	}
	// If no value provided, conv is ready for builder pattern

	return c
}

// Consolidated generic functions - single function per type with operation parameter
func genInt[T anyInt](c *conv, v T, op int) {
	c.intVal = int64(v)
	switch op {
	case 0:
		c.vTpe = typeInt // setValue
	case 1:
		c.fmtInt(10) // any2s
	case 2:
		c.i2s() // format
	}
}

func genUint[T anyUint](c *conv, v T, op int) {
	c.uintVal = uint64(v)
	switch op {
	case 0:
		c.vTpe = typeUint // setValue
	case 1:
		c.fmtUint(10) // any2s
	case 2:
		c.u2s() // format
	}
}

func genFloat[T anyFloat](c *conv, v T, op int) {
	c.floatVal = float64(v)
	switch op {
	case 0:
		c.vTpe = typeFloat // setValue
	case 1:
		c.f2s() // any2s
	case 2:
		c.f2sMan(-1) // format
	}
}

// Apply updates the original string pointer with the current content and auto-releases to pool.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *conv) Apply() {
	if t.vTpe == typeStrPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = t.getString()
	}
	// Auto-release back to pool for memory efficiency
	t.putConv()
}

// String method to return the content of the conv and automatically returns object to pool
// Phase 7: Auto-release makes pool usage completely transparent to user
func (t *conv) String() string {
	result := t.getString()
	// Auto-release back to pool for memory efficiency
	t.putConv()
	return result
}

// StringError returns the content of the conv along with any error and auto-releases to pool
func (t *conv) StringError() (string, error) {
	var result string
	var err error

	// BUILDER INTEGRATION: Check for error condition more comprehensively
	if t.err != "" {
		// If there's an error, return empty string and the error
		result = ""
		err = &customError{message: t.err}
	} else {
		result = t.getString()
		err = nil
	}

	// Auto-release back to pool for memory efficiency
	t.putConv()
	return result, err
}

// customError implements error interface for StringError
type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}

// getString converts the current value to string only when needed
// BUILDER INTEGRATION: Prioritizes buffer content when available
// Optimized with string caching to avoid repeated conversions
func (t *conv) getString() string {
	if t.vTpe == typeErr {
		return ""
	}

	// BUILDER PRIORITY: If buffer has content, use it as source of truth
	if len(t.buf) > 0 {
		return string(t.buf)
	}

	// If we already have a string value and haven't changed types, reuse it
	if t.vTpe == typeStr && t.stringVal != "" {
		return t.stringVal
	}

	// For string pointers, always return the current value (don't use cache)
	if t.vTpe == typeStrPtr && t.stringPtrVal != nil {
		return *t.stringPtrVal
	}

	// Use cached string if available and type hasn't changed (but not for string pointers)
	if t.tmpStr != "" && t.lastConvType == t.vTpe && t.vTpe != typeStrPtr {
		return t.tmpStr
	}
	// Convert to string using internal methods to avoid allocations
	switch t.vTpe {
	case typeStr:
		t.tmpStr = t.stringVal
	case typeStrPtr:
		// For string pointers, always get the current value from the pointer
		if t.stringPtrVal != nil {
			t.tmpStr = *t.stringPtrVal
		} else {
			t.tmpStr = t.stringVal // Fallback to stored value
		}
	case typeStrSlice:
		if len(t.stringSliceVal) == 0 {
			t.tmpStr = ""
		} else {
			// Join with space as default - use internal method
			t.tmpStr = t.joinSlice(" ")
		}
	case typeInt:
		// Use internal method instead of external function
		t.fmtInt(10)
	case typeUint:
		// Use internal method instead of external function
		t.fmtUint(10)
	case typeFloat:
		// Use internal method instead of external function
		t.f2s()
	case typeBool:
		if t.boolVal {
			t.tmpStr = trueStr
		} else {
			t.tmpStr = falseStr
		}
	default:
		t.tmpStr = ""
	}

	// Update cache state
	t.lastConvType = t.vTpe
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

// setString converts to string type and stores the value
func (t *conv) setString(s string) {
	t.stringVal = s

	// If working with string pointer, update the original string
	if t.vTpe == typeStrPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = s
		// Keep the vTpe as stringPtr to maintain the pointer relationship
	} else {
		t.vTpe = typeStr
	}

	// Clear other values to save memory
	t.intVal = 0
	t.uintVal = 0
	t.floatVal = 0
	t.boolVal = false
	t.stringSliceVal = nil

	// Invalidate cache since we changed the string
	t.tmpStr = ""
	t.lastConvType = vTpe(0)
}

// joinSlice joins string slice with separator - optimized using builder API
func (t *conv) joinSlice(separator string) string {
	if len(t.stringSliceVal) == 0 {
		return ""
	}
	if len(t.stringSliceVal) == 1 {
		return t.stringSliceVal[0]
	}

	// Use builder API for zero-allocation string construction
	c := Convert() // Empty initialization for builder pattern

	for i, s := range t.stringSliceVal {
		if i > 0 {
			c.Write(separator)
		}
		c.Write(s)
	}

	return c.String() // Auto-release to pool
}

// Internal conversion methods - centralized in conv to minimize allocations
// These methods modify the conv struct directly instead of returning values

// any2s converts any type to string and stores in tmpStr
// default set to "" if no conversion is possible
// supports int, uint, float, bool, string and error types
func (t *conv) any2s(v any) {
	switch val := v.(type) {
	case error:
		t.err = val.Error()
	case string:
		t.stringVal = val
		t.tmpStr = val
	case bool:
		if val {
			t.tmpStr = trueStr
		} else {
			t.tmpStr = falseStr
		}
		t.stringVal = t.tmpStr
	default:
		// Handle numeric types using generics for any2s
		t.handleAnyTypeForAny2s(v)
	}
}

// handleIntType processes integer types using generics
// handleAnyType processes any numeric type using generics
func (c *conv) handleAnyType(val any) {
	switch v := val.(type) {
	case int:
		genInt(c, v, 0)
	case int8:
		genInt(c, v, 0)
	case int16:
		genInt(c, v, 0)
	case int32:
		genInt(c, v, 0)
	case int64:
		genInt(c, v, 0)
	case uint:
		genUint(c, v, 0)
	case uint8:
		genUint(c, v, 0)
	case uint16:
		genUint(c, v, 0)
	case uint32:
		genUint(c, v, 0)
	case uint64:
		genUint(c, v, 0)
	case float32:
		genFloat(c, v, 0)
	case float64:
		genFloat(c, v, 0)
	}
}

// handleAnyTypeForAny2s processes any numeric type for any2s using generics
func (c *conv) handleAnyTypeForAny2s(val any) {
	switch v := val.(type) {
	case int:
		genInt(c, v, 1)
	case int8:
		genInt(c, v, 1)
	case int16:
		genInt(c, v, 1)
	case int32:
		genInt(c, v, 1)
	case int64:
		genInt(c, v, 1)
	case uint:
		genUint(c, v, 1)
	case uint8:
		genUint(c, v, 1)
	case uint16:
		genUint(c, v, 1)
	case uint32:
		genUint(c, v, 1)
	case uint64:
		genUint(c, v, 1)
	case float32:
		genFloat(c, v, 1)
	case float64:
		genFloat(c, v, 1)
	}
}
