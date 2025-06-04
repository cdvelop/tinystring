package tinystring

// valType represents the type of value stored in conv
type valType uint8

const (
	valTypeString valType = iota
	valTypeInt
	valTypeUint
	valTypeFloat
	valTypeBool
	valTypeStringSlice
	valTypeStringPtr
)

type conv struct {
	stringVal         string
	intVal            int64
	uintVal           uint64
	floatVal          float64
	boolVal           bool
	stringSliceVal    []string
	stringPtrVal      *string
	valType           valType
	err               error
	roundDown         bool
	separator         string
	cachedString      string  // Cache for string conversion to avoid repeated work
	lastConvertedType valType // Track last converted type for cache validation
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
func convInit(v any) *conv {
	// Handle nil case explicitly
	if v == nil {
		return &conv{stringVal: "", valType: valTypeString}
	}

	switch val := v.(type) {
	case []string:
		return &conv{stringSliceVal: val, valType: valTypeStringSlice}
	case *string:
		return &conv{stringVal: *val, stringPtrVal: val, valType: valTypeStringPtr}
	case string:
		return &conv{stringVal: val, valType: valTypeString}
	case int:
		return &conv{intVal: int64(val), valType: valTypeInt}
	case int8:
		return &conv{intVal: int64(val), valType: valTypeInt}
	case int16:
		return &conv{intVal: int64(val), valType: valTypeInt}
	case int32:
		return &conv{intVal: int64(val), valType: valTypeInt}
	case int64:
		return &conv{intVal: val, valType: valTypeInt}
	case uint:
		return &conv{uintVal: uint64(val), valType: valTypeUint}
	case uint8:
		return &conv{uintVal: uint64(val), valType: valTypeUint}
	case uint16:
		return &conv{uintVal: uint64(val), valType: valTypeUint}
	case uint32:
		return &conv{uintVal: uint64(val), valType: valTypeUint}
	case uint64:
		return &conv{uintVal: val, valType: valTypeUint}
	case float32:
		return &conv{floatVal: float64(val), valType: valTypeFloat}
	case float64:
		return &conv{floatVal: val, valType: valTypeFloat}
	case bool:
		return &conv{boolVal: val, valType: valTypeBool}
	default:
		// Fallback to string conversion for unknown types - use internal method to avoid allocation
		c := &conv{valType: valTypeString}
		c.anyToStringInternal(v)
		return c
	}
}

func (t *conv) transformWithMapping(mappings []charMapping) *conv {
	str := t.getString()

	// Use pooled builder for efficient string construction
	builder := getBuilder()
	defer putBuilder(builder)

	// Pre-allocate builder with exact string length
	builder.grow(len(str))

	hasChanges := false
	for _, r := range str {
		mapped := false
		for _, mapping := range mappings {
			if r == mapping.from {
				builder.writeRune(mapping.to)
				mapped = true
				hasChanges = true
				break
			}
		}
		if !mapped {
			builder.writeRune(r)
		}
	}

	// If no changes were made, return self to avoid allocation
	if !hasChanges {
		return t
	}

	newStr := builder.string()

	// Always modify in place to avoid creating new instances
	t.setString(newStr)
	return t
}

// Remueve tildes y diacríticos
func (t *conv) RemoveTilde() *conv {
	return t.transformWithMapping(accentMappings)
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
	if t.valType == valTypeStringPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = t.getString()
	}
}

// String method to return the content of the conv without modifying any original pointers
func (t *conv) String() string {
	return t.getString()
}

// StringError returns the content of the conv along with any error that occurred during processing
func (t *conv) StringError() (string, error) {
	return t.getString(), t.err
}

// splitIntoWordsLocal returns words as local variable without storing in struct field
// This avoids persistent memory allocation in the conv struct
func (t *conv) splitIntoWordsLocal() [][]rune {
	str := t.getString()
	if len(str) == 0 {
		return nil
	}

	// Pre-allocate slices with estimated capacity to reduce allocations
	words := make([][]rune, 0, 8) // Estimate 8 words max

	// Use a more efficient approach: build words directly without intermediate copies
	var start int
	inWord := false

	for i, r := range str {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if inWord {
				// Extract word directly from string range
				word := make([]rune, 0, i-start)
				for _, char := range str[start:i] {
					word = append(word, char)
				}
				words = append(words, word)
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
	if inWord {
		word := make([]rune, 0, len(str)-start)
		for _, char := range str[start:] {
			word = append(word, char)
		}
		words = append(words, word)
	}

	return words
}

func (t *conv) transformWord(word []rune, transform wordTransform) []rune {
	if len(word) == 0 {
		return word
	}

	// Create a copy to avoid modifying the original
	result := make([]rune, len(word))
	copy(result, word)

	switch transform {
	case toLower:
		for i, r := range result {
			for _, mapping := range lowerMappings {
				if r == mapping.from {
					result[i] = mapping.to
					break
				}
			}
		}
	case toUpper:
		for i, r := range result {
			for _, mapping := range upperMappings {
				if r == mapping.from {
					result[i] = mapping.to
					break
				}
			}
		}
	}

	// Create a copy to return
	resultCopy := make([]rune, len(result))
	copy(resultCopy, result)
	return resultCopy
}

// Helper function to check if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= 'À' && r <= 'ÿ' && r != '×' && r != '÷')
}

// transformSingleRune applies a character mapping to a single rune.
// It returns the transformed rune and true if a mapping was applied, otherwise the original rune and false.
func transformSingleRune(r rune, mappings []charMapping) (rune, bool) {
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
	if t.err != nil {
		return ""
	}

	// If we already have a string value and haven't changed types, reuse it
	if t.valType == valTypeString && t.stringVal != "" {
		return t.stringVal
	}

	// Use cached string if available and type hasn't changed
	if t.cachedString != "" && t.lastConvertedType == t.valType {
		return t.cachedString
	}

	// Convert to string using internal methods to avoid allocations
	switch t.valType {
	case valTypeString:
		t.cachedString = t.stringVal
	case valTypeStringPtr:
		t.cachedString = t.stringVal // Already stored during creation
	case valTypeStringSlice:
		if len(t.stringSliceVal) == 0 {
			t.cachedString = ""
		} else {
			// Join with space as default - use internal method
			t.cachedString = t.joinSlice(" ")
		}
	case valTypeInt:
		// Use internal method instead of external function
		t.formatIntInternal(t.intVal, 10)
	case valTypeUint:
		// Use internal method instead of external function
		t.formatUintInternal(t.uintVal, 10)
	case valTypeFloat:
		// Use internal method instead of external function
		t.formatFloatInternal(t.floatVal)
	case valTypeBool:
		if t.boolVal {
			t.cachedString = "true"
		} else {
			t.cachedString = "false"
		}
	default:
		t.cachedString = ""
	}

	// Update cache state
	t.lastConvertedType = t.valType
	return t.cachedString
}

// setString converts to string type and stores the value
func (t *conv) setString(s string) {
	t.stringVal = s

	// If working with string pointer, update the original string
	if t.valType == valTypeStringPtr && t.stringPtrVal != nil {
		*t.stringPtrVal = s
		// Keep the valType as stringPtr to maintain the pointer relationship
	} else {
		t.valType = valTypeString
	}

	// Clear other values to save memory
	t.intVal = 0
	t.uintVal = 0
	t.floatVal = 0
	t.boolVal = false
	t.stringSliceVal = nil

	// Invalidate cache since we changed the string
	t.cachedString = ""
	t.lastConvertedType = valType(0)
}

// joinSlice joins string slice with separator
func (t *conv) joinSlice(separator string) string {
	if len(t.stringSliceVal) == 0 {
		return ""
	}
	if len(t.stringSliceVal) == 1 {
		return t.stringSliceVal[0]
	}

	// Calculate total length to minimize allocations
	totalLen := 0
	for _, s := range t.stringSliceVal {
		totalLen += len(s)
	}
	totalLen += len(separator) * (len(t.stringSliceVal) - 1)

	// Build result string efficiently with pooled builder
	builder := getBuilder()
	defer putBuilder(builder)
	builder.grow(totalLen)

	for i, s := range t.stringSliceVal {
		if i > 0 {
			builder.writeString(separator)
		}
		builder.writeString(s)
	}

	return builder.string()
}

// Internal conversion methods - centralized in conv to minimize allocations
// These methods modify the conv struct directly instead of returning values

// anyToStringInternal converts any type to string and stores in cachedString
func (t *conv) anyToStringInternal(v any) {
	switch val := v.(type) {
	case string:
		t.stringVal = val
		t.cachedString = val
	case int:
		t.formatIntInternal(int64(val), 10)
	case int8:
		t.formatIntInternal(int64(val), 10)
	case int16:
		t.formatIntInternal(int64(val), 10)
	case int32:
		t.formatIntInternal(int64(val), 10)
	case int64:
		t.formatIntInternal(val, 10)
	case uint:
		t.formatUintInternal(uint64(val), 10)
	case uint8:
		t.formatUintInternal(uint64(val), 10)
	case uint16:
		t.formatUintInternal(uint64(val), 10)
	case uint32:
		t.formatUintInternal(uint64(val), 10)
	case uint64:
		t.formatUintInternal(val, 10)
	case float32:
		t.formatFloatInternal(float64(val))
	case float64:
		t.formatFloatInternal(val)
	case bool:
		if val {
			t.cachedString = "true"
		} else {
			t.cachedString = "false"
		}
		t.stringVal = t.cachedString
	default:
		t.cachedString = "unknown"
		t.stringVal = t.cachedString
	}
}
