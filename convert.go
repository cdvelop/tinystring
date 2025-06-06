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

// Convert initializes a new conv struct with any type of value for string,bool and number manipulation.
// Uses the centralized convInit function to avoid code duplication.
func Convert(v any) *conv {
	return convInit(v)
}

// convInit initializes a new conv struct with any type of value for string,bool and number manipulation.
// This is the centralized initialization function shared by Convert(), Format(), Sprintf(), etc.
// Uses optimized union-type storage to avoid unnecessary string conversions.
func convInit(value any) *conv {
	// Initialize conv struct
	c := &conv{}

	// Handle empty case or nil first value
	if value == nil {
		c.stringVal = ""
		c.vTpe = typeStr
		return c
	}
	switch val := value.(type) {
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
		// Try consolidated type handlers
		if c.handleIntTypes(value) {
			return c
		}
		if c.handleUintTypes(value) {
			return c
		}
		if c.handleFloatTypes(value) {
			return c
		}

		// Fallback to string conversion for unknown types - use internal method to avoid allocation
		c.vTpe = typeStr
		c.any2s(value)
	}

	return c
}

// setIntVal sets the int64 value and updates the vTpe
func (c *conv) setIntVal(val int64) {
	c.intVal = val
	c.vTpe = typeInt
}

// setUintVal sets the uint64 value and updates the vTpe
func (c *conv) setUintVal(val uint64) {
	c.uintVal = val
	c.vTpe = typeUint
}

// setFloatVal sets the float64 value and updates the vTpe
func (c *conv) setFloatVal(val float64) {
	c.floatVal = val
	c.vTpe = typeFloat
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

func (t *conv) tmap(mappings []charMapping) *conv {
	str := t.getString()

	// Use pre-allocated buffer for efficient string construction
	buf := make([]byte, 0, len(str)*2) // Allocate extra space for UTF-8 encoding

	hc := false // hasChanges
	for _, r := range str {
		mapped := false
		for _, m := range mappings { // mapping
			if r == m.from {
				buf = addRne2Buf(buf, m.to)
				hc = true
				mapped = true
				break
			}
		}
		if !mapped {
			buf = addRne2Buf(buf, r)
		}
	}

	// If no changes were made, return self to avoid allocation
	if !hc {
		return t
	}

	ns := string(buf) // newStr

	// Always modify in place to avoid creating new instances
	t.setString(ns)
	return t
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

// splitIntoWordsLocal returns words as local variable without storing in struct field
// Optimized for minimal allocations - Phase 3
func (t *conv) split() [][]rune {
	str := t.getString()
	if len(str) == 0 {
		return nil
	}

	// Pre-allocate based on estimated word count (rough heuristic: len/5)
	ew := (len(str) / 5) + 1 // estimatedWords
	if ew > 16 {
		ew = 16 // Cap reasonable maximum
	}

	words := make([][]rune, 0, ew)

	// Convert entire string to runes once
	runes := []rune(str)

	var start int
	inWord := false

	for i, r := range runes {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if inWord {
				// Extract word slice directly from runes - no copying
				if i > start {
					word := runes[start:i:i] // Limit capacity to avoid sharing
					words = append(words, word)
				}
				inWord = false
			}
		} else {
			if !inWord {
				start = i
				inWord = true
			}
		}
	}

	// Handle last word if string doesn't end with whitespace
	if inWord && len(runes) > start {
		word := runes[start:len(runes):len(runes)]
		words = append(words, word)
	}

	return words
}

func (t *conv) transformWord(word []rune, transform wordTransform) []rune {
	if len(word) == 0 {
		return word
	}

	// Transform in-place to avoid allocation, then copy once
	switch transform {
	case toLower:
		for i, r := range word {
			word[i] = toLowerRune(r)
		}
	case toUpper:
		for i, r := range word {
			word[i] = toUpperRune(r)
		}
	}

	// Return the transformed word (caller will handle copying if needed)
	return word
}

// Helper function to check if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= 'À' && r <= 'ÿ' && r != 'x' && r != '÷')
}

// trRune applies a character mapping to a single rune.
// It returns the transformed rune and true if a mapping was applied, otherwise the original rune and false.
func trRune(r rune, mappings []charMapping) (rune, bool) {
	for _, mapping := range mappings {
		if r == mapping.from {
			return mapping.to, true
		}
	}
	return r, false
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
		// Use consolidated type handlers
		if t.handleIntTypesForAny2s(v) {
			return
		}
		if t.handleUintTypesForAny2s(v) {
			return
		}
		if t.handleFloatTypesForAny2s(v) {
			return
		}

		// Fallback for unknown types
		t.tmpStr = ""
		t.stringVal = t.tmpStr
	}
}

// Helper functions to reduce code duplication in type handling

// handleIntTypes consolidates repetitive int type handling
func (c *conv) handleIntTypes(val any) bool {
	switch v := val.(type) {
	case int:
		c.setIntVal(int64(v))
		return true
	case int8:
		c.setIntVal(int64(v))
		return true
	case int16:
		c.setIntVal(int64(v))
		return true
	case int32:
		c.setIntVal(int64(v))
		return true
	case int64:
		c.setIntVal(v)
		return true
	default:
		return false
	}
}

// handleUintTypes consolidates repetitive uint type handling
func (c *conv) handleUintTypes(val any) bool {
	switch v := val.(type) {
	case uint:
		c.setUintVal(uint64(v))
		return true
	case uint8:
		c.setUintVal(uint64(v))
		return true
	case uint16:
		c.setUintVal(uint64(v))
		return true
	case uint32:
		c.setUintVal(uint64(v))
		return true
	case uint64:
		c.setUintVal(v)
		return true
	default:
		return false
	}
}

// handleFloatTypes consolidates repetitive float type handling
func (c *conv) handleFloatTypes(val any) bool {
	switch v := val.(type) {
	case float32:
		c.setFloatVal(float64(v))
		return true
	case float64:
		c.setFloatVal(v)
		return true
	default:
		return false
	}
}

// handleIntTypesForAny2s consolidates repetitive int type handling for any2s
func (c *conv) handleIntTypesForAny2s(val any) bool {
	switch v := val.(type) {
	case int:
		c.intVal = int64(v)
		c.fmtInt(10)
		return true
	case int8:
		c.intVal = int64(v)
		c.fmtInt(10)
		return true
	case int16:
		c.intVal = int64(v)
		c.fmtInt(10)
		return true
	case int32:
		c.intVal = int64(v)
		c.fmtInt(10)
		return true
	case int64:
		c.intVal = v
		c.fmtInt(10)
		return true
	default:
		return false
	}
}

// handleUintTypesForAny2s consolidates repetitive uint type handling for any2s
func (c *conv) handleUintTypesForAny2s(val any) bool {
	switch v := val.(type) {
	case uint:
		c.uintVal = uint64(v)
		c.fmtUint(10)
		return true
	case uint8:
		c.uintVal = uint64(v)
		c.fmtUint(10)
		return true
	case uint16:
		c.uintVal = uint64(v)
		c.fmtUint(10)
		return true
	case uint32:
		c.uintVal = uint64(v)
		c.fmtUint(10)
		return true
	case uint64:
		c.uintVal = v
		c.fmtUint(10)
		return true
	default:
		return false
	}
}

// handleFloatTypesForAny2s consolidates repetitive float type handling for any2s
func (c *conv) handleFloatTypesForAny2s(val any) bool {
	switch v := val.(type) {
	case float32:
		c.floatVal = float64(v)
		c.f2s()
		return true
	case float64:
		c.floatVal = v
		c.f2s()
		return true
	default:
		return false
	}
}
