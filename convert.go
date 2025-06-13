package tinystring

// Import unified types from abi.go - no more duplication
// kind is now defined in abi.go with tp prefix

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
	vTpe           kind // Updated to use unified kind type from abi.go
	roundDown      bool
	separator      string
	tmpStr         string    // Cache for temp string conversion to avoid repeated work
	lastConvType   kind      // Track last converted type for cache validation
	err            errorType // Error type from error.go
	anyVal         any       // Generic field to store any value type
}

// struct to store mappings to remove accents and diacritics
type charMapping struct {
	from rune
	to   rune
}

type wordTransform int

const (
	toLower wordTransform = iota
	toUpper
)

// Functional options pattern for conv construction
type convOpt func(*conv)

// withValue initializes conv with any value type
func withValue(v any) convOpt {
	return func(c *conv) {
		if v == nil {
			c.stringVal = ""
			c.vTpe = tpString
			return
		}
		switch val := v.(type) {
		case string:
			c.stringVal = val
			c.vTpe = tpString
		case []string:
			c.stringSliceVal = val
			c.vTpe = tpStrSlice
		case *string:
			c.stringVal = *val
			c.stringPtrVal = val
			c.vTpe = tpStrPtr
		case bool:
			c.setBoolVal(val)
		case errorType:
			c.setErrorVal(val)
		default:
			// Handle numeric types using generics
			switch val := val.(type) {
			case int, int8, int16, int32, int64:
				c.handleIntType(val)
			case uint, uint8, uint16, uint32, uint64:
				c.handleUintType(val)
			case float32, float64:
				c.handleFloatType(val)
			default: // Check for struct or slice types using reflection
				rv := refValueOf(v)
				switch rv.Kind() {
				case tpStruct:
					// Store the value for struct encoding
					c.anyVal = v
					c.vTpe = tpStruct
				case tpSlice, tpArray:
					// Store the value for slice encoding
					c.anyVal = v
					c.vTpe = tpSlice
				case tpPointer:
					// Handle pointer types - dereference and store the pointed value
					elem := rv.Elem()
					if !elem.IsValid() {
						// Nil pointer
						c.stringVal = ""
						c.vTpe = tpString
					} else {
						// Dereference and check the pointed-to type
						switch elem.Kind() {
						case tpStruct:
							// Store the pointer for struct encoding (encoder will handle dereferencing)
							c.anyVal = v
							c.vTpe = tpStruct
						case tpSlice, tpArray:
							// Store the pointer for slice encoding
							c.anyVal = v
							c.vTpe = tpSlice
						default:
							// For other pointer types, convert to string
							c.vTpe = tpString
							c.any2s(val)
						}
					}
				default:
					// Fallback to string conversion for unknown types
					c.vTpe = tpString
					c.any2s(val)
				}
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
		stringVal: "",
		vTpe:      tpString,
		separator: "_",
	}
}

// Convert initializes a new conv struct with any type of value for string,bool and number manipulation.
// Uses the functional options pattern internally.
func Convert(v any) *conv {
	return newConv(withValue(v))
}

// Generic functions to handle numeric types without repetitive switches
func genInt[T anyInt](c *conv, v T) {
	c.intVal = int64(v)
	c.vTpe = tpInt
}

func genUint[T anyUint](c *conv, v T) {
	c.uintVal = uint64(v)
	c.vTpe = tpUint
}

func genFloat[T anyFloat](c *conv, v T) {
	c.floatVal = float64(v)
	c.vTpe = tpFloat64
}

// Generic functions for any2s operations
func genAny2sInt[T anyInt](c *conv, v T) {
	c.intVal = int64(v)
	c.fmtInt(10)
}

func genAny2sUint[T anyUint](c *conv, v T) {
	c.uintVal = uint64(v)
	c.fmtUint(10)
}

func genAny2sFloat[T anyFloat](c *conv, v T) {
	c.floatVal = float64(v)
	c.f2s()
}

// Generic functions for format operations
func genFormatInt[T anyInt](c *conv, v T) {
	c.intVal = int64(v)
	c.i2s()
}

func genFormatUint[T anyUint](c *conv, v T) {
	c.uintVal = uint64(v)
	c.u2s()
}

func genFormatFloat[T anyFloat](c *conv, v T) {
	c.floatVal = float64(v)
	c.f2sMan(-1)
}

// setBoolVal sets the bool value and updates the vTpe
func (c *conv) setBoolVal(val bool) {
	c.boolVal = val
	c.vTpe = tpBool
}

// setErrorVal sets the errorType value and updates the vTpe
func (c *conv) setErrorVal(val errorType) {
	c.err = val
	c.vTpe = tpErr
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
// Optimized with string caching to avoid repeated conversions
func (t *conv) getString() string {
	if t.vTpe == tpErr {
		return ""
	}

	// If we already have a string value and haven't changed types, reuse it
	if t.vTpe == tpString && t.stringVal != "" {
		return t.stringVal
	}

	// Use cached string if available and type hasn't changed
	if t.tmpStr != "" && t.lastConvType == t.vTpe {
		return t.tmpStr
	}

	// Convert to string using internal methods to avoid allocations
	switch t.vTpe {
	case tpString:
		t.tmpStr = t.stringVal
	case tpStrPtr:
		t.tmpStr = t.stringVal // Already stored during creation
	case tpStrSlice:
		if len(t.stringSliceVal) == 0 {
			t.tmpStr = ""
		} else {
			// Join with space as default - use internal method
			t.tmpStr = t.joinSlice(" ")
		}
	case tpInt:
		// Use internal method instead of external function
		t.fmtInt(10)
	case tpUint:
		// Use internal method instead of external function
		t.fmtUint(10)
	case tpFloat64:
		// Use internal method instead of external function
		t.f2s()
	case tpBool:
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
	if t.vTpe == tpStrPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = s
		// Keep the vTpe as stringPtr to maintain the pointer relationship
	} else {
		t.vTpe = tpString
	}

	// Clear other values to save memory
	t.intVal = 0
	t.uintVal = 0
	t.floatVal = 0
	t.boolVal = false
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
		switch v := v.(type) {
		case int, int8, int16, int32, int64:
			t.handleIntTypeForAny2s(v)
		case uint, uint8, uint16, uint32, uint64:
			t.handleUintTypeForAny2s(v)
		case float32, float64:
			t.handleFloatTypeForAny2s(v)
		default:
			// Fallback for unknown types
			t.tmpStr = ""
			t.stringVal = t.tmpStr
		}
	}
}

// handleIntType processes integer types using generics
func (c *conv) handleIntType(val any) {
	switch v := val.(type) {
	case int:
		genInt(c, v)
	case int8:
		genInt(c, v)
	case int16:
		genInt(c, v)
	case int32:
		genInt(c, v)
	case int64:
		genInt(c, v)
	}
}

// handleUintType processes unsigned integer types using generics
func (c *conv) handleUintType(val any) {
	switch v := val.(type) {
	case uint:
		genUint(c, v)
	case uint8:
		genUint(c, v)
	case uint16:
		genUint(c, v)
	case uint32:
		genUint(c, v)
	case uint64:
		genUint(c, v)
	}
}

// handleFloatType processes float types using generics
func (c *conv) handleFloatType(val any) {
	switch v := val.(type) {
	case float32:
		genFloat(c, v)
	case float64:
		genFloat(c, v)
	}
}

// handleIntTypeForAny2s processes integer types for any2s using generics
func (c *conv) handleIntTypeForAny2s(val any) {
	switch v := val.(type) {
	case int:
		genAny2sInt(c, v)
	case int8:
		genAny2sInt(c, v)
	case int16:
		genAny2sInt(c, v)
	case int32:
		genAny2sInt(c, v)
	case int64:
		genAny2sInt(c, v)
	}
}

// handleUintTypeForAny2s processes unsigned integer types for any2s using generics
func (c *conv) handleUintTypeForAny2s(val any) {
	switch v := val.(type) {
	case uint:
		genAny2sUint(c, v)
	case uint8:
		genAny2sUint(c, v)
	case uint16:
		genAny2sUint(c, v)
	case uint32:
		genAny2sUint(c, v)
	case uint64:
		genAny2sUint(c, v)
	}
}

// handleFloatTypeForAny2s processes float types for any2s using generics
func (c *conv) handleFloatTypeForAny2s(val any) {
	switch v := val.(type) {
	case float32:
		genAny2sFloat(c, v)
	case float64:
		genAny2sFloat(c, v)
	}
}

// Generic helper functions are all defined above
