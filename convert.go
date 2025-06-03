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

// convInit initializes a new conv struct with any type of value for string,bool and number manipulation.
// This is the centralized initialization function shared by Convert(), Format(), Sprintf(), etc.
// Uses optimized union-type storage to avoid unnecessary string conversions.
func convInit(v any) *conv {
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

// Convert initializes a new conv struct with any type of value for string,bool and number manipulation.
// Uses the centralized convInit function to avoid code duplication.
func Convert(v any) *conv {
	return convInit(v)
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

// convert to lower case eg: "HELLO WORLD" -> "hello world"
func (t *conv) ToLower() *conv {
	return t.transformWithMapping(lowerMappings)
}

// convert to upper case eg: "hello world" -> "HELLO WORLD"
func (t *conv) ToUpper() *conv {
	return t.transformWithMapping(upperMappings)
}

// converts conv to camelCase (first word lowercase) eg: "Hello world" -> "helloWorld"
func (t *conv) CamelCaseLower() *conv {
	return t.toCaseTransform(true, "")
}

// converts conv to PascalCase (all words capitalized) eg: "hello world" -> "HelloWorld"
func (t *conv) CamelCaseUpper() *conv {
	return t.toCaseTransform(false, "")
}

// snakeCase converts a string to snake_case format with optional separator.
// If no separator is provided, underscore "_" is used as default.
// Example:
//
//	Input: "camelCase" -> Output: "camel_case"
//	Input: "PascalCase", "-" -> Output: "pascal-case"
//	Input: "APIResponse" -> Output: "api_response"
//	Input: "user123Name", "." -> Output: "user123.name"
//
// ToSnakeCaseLower converts conv to snake_case format
func (t *conv) ToSnakeCaseLower(sep ...string) *conv {
	return t.toCaseTransform(true, t.separatorCase(sep...))
}

// ToSnakeCaseUpper converts conv to Snake_Case format
func (t *conv) ToSnakeCaseUpper(sep ...string) *conv {
	return t.toCaseTransform(false, t.separatorCase(sep...))
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

func (t *conv) toCaseTransform(firstWordLower bool, separator string) *conv {
	// Use local variable instead of struct field to avoid persistent allocation
	words := t.splitIntoWordsLocal()
	if len(words) == 0 {
		return t
	}

	// Use pooled string builder for efficient string construction
	builder := getBuilder()
	defer putBuilder(builder)

	str := t.getString()
	estimatedLen := len(str) + len(words)*len(separator)
	builder.grow(estimatedLen)
	var prevIsDigit bool
	var prevIsSeparator bool

	for i, word := range words {
		if len(word) == 0 {
			continue
		}

		// Add separator if needed
		if i > 0 && separator != "" {
			builder.writeByte(separator[0])
			prevIsSeparator = true
		}

		// Process each character in the word
		for j, r := range word {
			currentCaseTransform := toLower // Default to lower
			currIsDigit := isDigit(r)
			currIsLetter := isLetter(r)

			// Determine case transform
			if i == 0 && j == 0 { // First letter of first word
				if !firstWordLower {
					currentCaseTransform = toUpper
				}
			} else if i > 0 && j == 0 && separator == "" { // Start of new word in camelCase
				currentCaseTransform = toUpper
			} else if prevIsDigit && currIsLetter { // Letter after digit
				// For snake_case with separator, this is handled by adding separator later
				// For camelCase, new word starts, so apply upper if not firstWordLower
				if separator == "" && !firstWordLower {
					currentCaseTransform = toUpper
				} else if separator == "" && firstWordLower {
					// Maintain lower case if it's camelCaseLower and after a digit within the same "word" part
					currentCaseTransform = toLower
				} else if separator != "" && !firstWordLower { // Snake_Case_Upper
					currentCaseTransform = toUpper
				}

			} else if prevIsSeparator && currIsLetter { // Letter after separator (for snake_case)
				if separator != "" && !firstWordLower { // Snake_Case_Upper
					currentCaseTransform = toUpper
				}
				// Default toLower is fine for snake_case_lower
			} else if currIsLetter && j > 0 { // Subsequent letters in a word part
				// Maintain lower case unless it's an uppercase letter that should remain uppercase (e.g. in "APIResponse")
				// This part is tricky without knowing the original casing or specific rules for acronyms.
				// For now, default toLower is applied unless specific conditions for toUpper are met.
				// If the global transform is toUpper (e.g. ToSnakeCaseUpper), this will be handled.
				if separator != "" && !firstWordLower { // Snake_Case_Upper
					// This condition might be too broad.
					// We only want to uppercase the first letter after a separator.
					// Let's rely on the j==0 condition for this.
					// For subsequent letters in Snake_Case_Upper, they should be lower.
					// So, if j > 0, it should be toLower for Snake_Case_Upper.
					// This means the currentCaseTransform should be toLower here.
				}
			}

			// Add underscore for number to letter transition in snake_case
			if separator != "" && prevIsDigit && currIsLetter {
				builder.writeByte(separator[0])
			}

			if currIsLetter {
				var transformedRune rune
				var mapped bool
				if currentCaseTransform == toLower {
					transformedRune, mapped = transformSingleRune(r, lowerMappings)
				} else { // toUpper
					transformedRune, mapped = transformSingleRune(r, upperMappings)
				}
				if mapped {
					builder.writeRune(transformedRune)
				} else {
					builder.writeRune(r) // Write original if no mapping found (e.g. already correct case)
				}
			} else {
				builder.writeRune(r)
			}

			prevIsDigit = currIsDigit
			prevIsSeparator = false // Reset after processing the character
		}
	}

	t.setString(builder.string())
	// Clear the separator field after use to avoid memory overhead
	t.separator = ""
	return t
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

// formatIntInternal converts integer to string and stores in cachedString
func (t *conv) formatIntInternal(val int64, base int) {
	if val == 0 {
		t.cachedString = "0"
		t.stringVal = t.cachedString
		return
	}

	// Use pooled builder for conversion
	builder := getBuilder()
	defer putBuilder(builder)

	negative := val < 0
	if negative {
		val = -val
	}

	// Convert digits in reverse order
	for val > 0 {
		digit := val % int64(base)
		if digit < 10 {
			builder.writeByte(byte('0' + digit))
		} else {
			builder.writeByte(byte('a' + digit - 10))
		}
		val /= int64(base)
	}

	if negative {
		builder.writeByte('-')
	}

	// Reverse the string since we built it backwards
	buf := builder.buf
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	t.cachedString = builder.string()
	t.stringVal = t.cachedString
}

// intToStringOptimizedInternal converts int64 to string with minimal allocations and stores in cachedString
func (t *conv) intToStringOptimizedInternal(val int64) {
	// Handle common small numbers using lookup table
	if val >= 0 && val < int64(len(smallInts)) {
		t.cachedString = smallInts[val]
		t.stringVal = t.cachedString
		return
	}

	// Handle special cases
	if val == 0 {
		t.cachedString = zeroString
		t.stringVal = t.cachedString
		return
	}
	if val == 1 {
		t.cachedString = oneString
		t.stringVal = t.cachedString
		return
	}

	// Fall back to standard conversion for larger numbers
	t.formatIntInternal(val, 10)
}

// formatUintInternal converts unsigned integer to string and stores in cachedString
func (t *conv) formatUintInternal(val uint64, base int) {
	if val == 0 {
		t.cachedString = "0"
		t.stringVal = t.cachedString
		return
	}

	// Use pooled builder for conversion
	builder := getBuilder()
	defer putBuilder(builder)

	// Convert digits in reverse order
	for val > 0 {
		digit := val % uint64(base)
		if digit < 10 {
			builder.writeByte(byte('0' + digit))
		} else {
			builder.writeByte(byte('a' + digit - 10))
		}
		val /= uint64(base)
	}

	// Reverse the string since we built it backwards
	buf := builder.buf
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	t.cachedString = builder.string()
	t.stringVal = t.cachedString
}

// uintToStringOptimizedInternal converts uint64 to string with minimal allocations and stores in cachedString
func (t *conv) uintToStringOptimizedInternal(val uint64) {
	// Handle common small numbers using lookup table
	if val < uint64(len(smallInts)) {
		t.cachedString = smallInts[val]
		t.stringVal = t.cachedString
		return
	}

	// Handle special cases
	if val == 0 {
		t.cachedString = zeroString
		t.stringVal = t.cachedString
		return
	}
	if val == 1 {
		t.cachedString = oneString
		t.stringVal = t.cachedString
		return
	}

	// Fall back to standard conversion for larger numbers
	t.uintToStringWithBaseInternal(val, 10)
}

// uintToStringWithBaseInternal converts unsigned integer to string with specified base
// and stores the result in the conv struct fields
func (t *conv) uintToStringWithBaseInternal(number uint64, base int) {
	if number == 0 {
		t.cachedString = "0"
		t.stringVal = t.cachedString
		return
	}

	// Max uint64 is 18446744073709551615, which has 20 digits.
	// For base 2, uint64 needs up to 64 bits.
	var buf [64]byte // Max buffer size for uint64 in base 2
	i := len(buf)    // Start from the end of the buffer

	const digitChars = "0123456789abcdefghijklmnopqrstuvwxyz"

	for number > 0 {
		i--
		buf[i] = digitChars[number%uint64(base)]
		number /= uint64(base)
	}

	t.cachedString = unsafeString(buf[i:])
	t.stringVal = t.cachedString
}

// formatFloatInternal converts float to string and stores in cachedString
func (t *conv) formatFloatInternal(val float64) {
	// Use pooled builder for conversion
	builder := getBuilder()
	defer putBuilder(builder)

	// Handle special cases
	if val != val { // NaN
		t.cachedString = "NaN"
		t.stringVal = t.cachedString
		return
	}
	if val == 0 {
		t.cachedString = "0"
		t.stringVal = t.cachedString
		return
	}

	negative := val < 0
	if negative {
		val = -val
		builder.writeByte('-')
	}

	// Handle infinity
	if val > 1e308 {
		builder.writeString("Inf")
		t.cachedString = builder.string()
		t.stringVal = t.cachedString
		return
	}

	// Simple float to string conversion (basic implementation)
	intPart := int64(val)
	fracPart := val - float64(intPart)

	// Convert integer part
	if intPart == 0 {
		builder.writeByte('0')
	} else {
		// Reuse the integer conversion logic
		tempConv := &conv{}
		tempConv.formatIntInternal(intPart, 10)
		builder.writeString(tempConv.cachedString)
	}

	// Add decimal point and fractional part if needed
	if fracPart > 0 {
		builder.writeByte('.')
		// Simple fractional part handling (could be improved)
		for i := 0; i < 6 && fracPart > 0; i++ {
			fracPart *= 10
			digit := int(fracPart)
			builder.writeByte(byte('0' + digit))
			fracPart -= float64(digit)
		}
	}

	t.cachedString = builder.string()
	t.stringVal = t.cachedString
}

// parseIntInternal parses string to int with specified base
// and stores result in conv struct fields
func (t *conv) parseIntInternal(input string, base int) error {
	if input == "" {
		return newEmptyStringError()
	}

	isNegative := false
	if input[0] == '-' {
		if base != 10 {
			return newError(errInvalidBase, "negative numbers are not supported for non-decimal bases")
		}
		isNegative = true
		input = input[1:]
	}

	number, err := t.parseNumberHelperInternal(input, base)
	if err != nil {
		return err
	}

	if isNegative {
		t.intVal = -int64(number)
	} else {
		t.intVal = int64(number)
	}
	return nil
}

// parseFloatInternal parses string to float64
// and stores result in conv struct fields
func (t *conv) parseFloatInternal(input string) error {
	if input == "" {
		return newEmptyStringError()
	}

	// Handle special cases using math constants
	switch input {
	case "NaN", "nan":
		// Use a variable to create NaN since direct division causes compile error
		var zero float64 = 0.0
		t.floatVal = zero / zero // NaN
		return nil
	case "Inf", "+Inf", "inf", "+inf":
		// Use a variable to create +Inf
		var one float64 = 1.0
		var zero float64 = 0.0
		t.floatVal = one / zero // +Inf
		return nil
	case "-Inf", "-inf":
		// Use a variable to create -Inf
		var negone float64 = -1.0
		var zero float64 = 0.0
		t.floatVal = negone / zero // -Inf
		return nil
	}

	// Simple manual parsing for common cases
	negative := false
	if input[0] == '-' {
		negative = true
		input = input[1:]
	} else if input[0] == '+' {
		input = input[1:]
	}

	// Find decimal point
	dotIndex := -1
	for i, ch := range input {
		if ch == '.' {
			dotIndex = i
			break
		}
	}

	var integerPart, fractionalPart uint64
	var err error

	if dotIndex == -1 {
		// No decimal point, parse as integer
		integerPart, err = t.parseNumberHelperInternal(input, 10)
		if err != nil {
			return err
		}
	} else {
		// Parse integer part
		if dotIndex > 0 {
			integerPart, err = t.parseNumberHelperInternal(input[:dotIndex], 10)
			if err != nil {
				return err
			}
		}

		// Parse fractional part
		if dotIndex+1 < len(input) {
			fractionalPart, err = t.parseNumberHelperInternal(input[dotIndex+1:], 10)
			if err != nil {
				return err
			}
		}
	}

	// Convert to float64
	result := float64(integerPart)
	if dotIndex != -1 && fractionalPart > 0 {
		// Calculate fractional part
		fractionalDigits := len(input) - dotIndex - 1
		divisor := float64(1)
		for i := 0; i < fractionalDigits; i++ {
			divisor *= 10
		}
		result += float64(fractionalPart) / divisor
	}

	if negative {
		result = -result
	}

	t.floatVal = result
	return nil
}

// parseNumberHelperInternal is an internal helper for parsing digits
// and stores result in conv struct fields
func (t *conv) parseNumberHelperInternal(input string, base int) (uint64, error) {
	if input == "" {
		return 0, newEmptyStringError()
	}

	if base < 2 || base > 36 {
		return 0, newError(errInvalidBase, "base must be between 2 and 36")
	}

	var result uint64

	for _, ch := range input {
		var digit int

		switch {
		case '0' <= ch && ch <= '9':
			digit = int(ch - '0')
		case 'a' <= ch && ch <= 'z':
			digit = int(ch - 'a' + 10)
		case 'A' <= ch && ch <= 'Z':
			digit = int(ch - 'A' + 10)
		default:
			return 0, newError(errInvalidNumber, "invalid character in number")
		}

		if digit >= base {
			return 0, newError(errInvalidNumber, "digit out of range for base")
		}

		// Check for overflow
		if result > (^uint64(0)-uint64(digit))/uint64(base) {
			return 0, newError(errOverflow, "number too large")
		}

		result = result*uint64(base) + uint64(digit)
	}

	return result, nil
}
