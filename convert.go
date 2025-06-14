package tinystring

import "strconv"

// Import unified types from abi.go - no more duplication
// kind is now defined in abi.go with tp prefix

type conv struct {
	// PRIMARY: Reflection-based value storage
	refVal refValue // All values accessed via reflection

	// ESSENTIAL: Core operation fields only
	vTpe         kind      // Type cache for performance
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

// withValue initializes conv with any value type using reflection-only approach
func withValue(v any) convOpt {
	return func(c *conv) {
		if v == nil {
			c.refVal = refValue{} // Invalid refValue for nil
			c.vTpe = tpString
			return
		}

		c.refVal = refValueOf(v)
		if !c.refVal.IsValid() {
			c.vTpe = tpString
			return
		}

		c.vTpe = c.refVal.Kind()

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
				// Complex types - value stored in refVal for JSON operations
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
	if c.refVal.IsValid() && (c.vTpe >= tpInt && c.vTpe <= tpInt64) {
		return c.refVal.Int()
	}
	return 0
}

func (c *conv) getUint64() uint64 {
	if c.refVal.IsValid() && (c.vTpe >= tpUint && c.vTpe <= tpUintptr) {
		return c.refVal.Uint()
	}
	return 0
}

func (c *conv) getFloat64() float64 {
	if !c.refVal.IsValid() {
		return 0
	}

	switch c.vTpe {
	case tpFloat32, tpFloat64:
		return c.refVal.Float()
	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		return float64(c.refVal.Int())
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr:
		return float64(c.refVal.Uint())
	default:
		return 0
	}
}

func (c *conv) getBool() bool {
	if c.refVal.IsValid() && c.vTpe == tpBool {
		return c.refVal.Bool()
	}
	return false
}

func (c *conv) getStringDirect() string {
	if c.refVal.IsValid() && c.vTpe == tpString {
		return c.refVal.String()
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
		return ""
	}

	// Use cached string if available and type hasn't changed
	if t.tmpStr != "" && t.lastConvType == t.vTpe {
		return t.tmpStr
	}

	// Convert to string using reflection-based methods
	switch t.vTpe {
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
		// Use reflection to get int value and direct string formatting
		intVal := t.getInt64()
		if intVal == 0 {
			t.tmpStr = "0"
		} else {
			t.tmpStr = strconv.FormatInt(intVal, 10)
		}
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr:
		// Use reflection to get uint value and direct string formatting
		uintVal := t.getUint64()
		if uintVal == 0 {
			t.tmpStr = "0"
		} else {
			t.tmpStr = strconv.FormatUint(uintVal, 10)
		}
	case tpFloat32:
		// Handle float32 with appropriate precision to match original behavior
		floatVal := t.getFloat64()
		t.tmpStr = strconv.FormatFloat(floatVal, 'g', 7, 32) // Use 32-bit precision
	case tpFloat64:
		// Use reflection to get float value and direct string formatting
		floatVal := t.getFloat64()
		t.tmpStr = strconv.FormatFloat(floatVal, 'g', -1, 64)
	case tpBool:
		if t.getBool() {
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

// setString converts to string type and stores the value using reflection
func (t *conv) setString(s string) {
	// Update refVal to hold the new string value
	t.refVal = refValueOf(s)

	// If working with string pointer, update the original string
	if t.vTpe == tpStrPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = s
		// Keep the vTpe as stringPtr to maintain the pointer relationship
	} else {
		t.vTpe = tpString
	}

	// Clear slice values to save memory - other values handled by refVal
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
		rv := refValueOf(v)
		if !rv.IsValid() {
			t.tmpStr = ""
			return
		}

		switch rv.Kind() {
		case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
			intVal := rv.Int()
			t.fmtIntGeneric(intVal, 10, true)
		case tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr:
			uintVal := rv.Uint()
			t.fmtIntGeneric(int64(uintVal), 10, false)
		case tpFloat32, tpFloat64:
			// Store value in refVal for f2s to work
			t.refVal = rv
			t.f2s()
		default:
			// Fallback for unknown types
			t.tmpStr = ""
		}
	}
}

// Generic helper functions are all defined above
