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
	tmpStr         string    // Cache for temp string conversion to avoid repeated work
	lastConvType   vTpe      // Track last converted type for cache validation
	err            errorType // Error type from error.go
}

// Functional options pattern for conv construction
type convOpt func(*conv)

// withValue initializes conv with any value type
func withValue(v any) convOpt {
	return func(c *conv) {
		if v == nil {
			c.stringVal = ""
			c.vTpe = typeStr
			return
		}
		switch val := v.(type) {
		case string:
			c.stringVal = val
			c.vTpe = typeStr
		case []string:
			c.stringSliceVal = val
			c.vTpe = typeStrSlice
		case *string:
			c.stringVal = *val
			c.stringPtrVal = val
			c.vTpe = typeStrPtr
		case bool:
			c.setBoolVal(val)
		case errorType:
			c.setErrorVal(val)
		default:
			// Handle numeric types using generics
			c.handleAnyType(val)
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

// Convert initializes a new conv struct with any type of value for string,bool and number manipulation.
// Uses the functional options pattern internally.
func Convert(v any) *conv {
	return newConv(withValue(v))
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

// setBoolVal sets the bool value and updates the vTpe
func (c *conv) setBoolVal(val bool) {
	c.boolVal = val
	c.vTpe = typeBool
}

// setErrorVal sets the errorType value and updates the vTpe
func (c *conv) setErrorVal(val errorType) {
	c.err = val
	c.vTpe = typeErr
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
	if t.vTpe == typeStrPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = t.getString()
	}
}

// String method to return the content of the conv without modifying any original pointers
func (t *conv) String() string {
	return t.getString()
}

// StringError returns the content of the conv along with any error that occurred during processing
func (t *conv) StringError() (string, error) {
	if t.vTpe == typeErr {
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
// Optimized with string caching to avoid repeated conversions
func (t *conv) getString() string {
	if t.vTpe == typeErr {
		return ""
	}

	// If we already have a string value and haven't changed types, reuse it
	if t.vTpe == typeStr && t.stringVal != "" {
		return t.stringVal
	}

	// Use cached string if available and type hasn't changed
	if t.tmpStr != "" && t.lastConvType == t.vTpe {
		return t.tmpStr
	}

	// Convert to string using internal methods to avoid allocations
	switch t.vTpe {
	case typeStr:
		t.tmpStr = t.stringVal
	case typeStrPtr:
		t.tmpStr = t.stringVal // Already stored during creation
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

// newBuf creates an optimally-sized buffer for common string operations
func (t *conv) newBuf(sizeMultiplier int) (string, []byte) {
	str := t.getString()
	if isEmpty(str) {
		return str, nil
	}
	bufSize := len(str) * sizeMultiplier
	if bufSize < 16 {
		bufSize = 16 // Minimum useful buffer size
	}
	return str, make([]byte, 0, bufSize)
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
	result := makeBuf(tl) // result

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

// any2s converts any type to string and stores in tmpStr
// default set to "" if no conversion is possible
// supports int, uint, float, bool, string and error types
func (t *conv) any2s(v any) {
	switch val := v.(type) {
	case errorType:
		t.err = val
	case error:
		t.err = errorType(val.Error())
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

// Generic helper functions are all defined above
