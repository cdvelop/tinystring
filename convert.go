package tinystring

// Text struct to store the content of the text - optimized for memory efficiency
type Text struct {
	content   string  // main content string
	stringPtr *string // pointer to original string (if one was provided)
	err       error   // internal error for error handling
	roundDown bool    // flag for down rounding in RoundDecimals

	// Temporary fields - these should be cleared after use to avoid memory overhead
	contentSlice []string // slice of strings for the Join method
	separator    string   // eg "_" "-" (temporary for case transforms)
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

// Convert initializes a new Text struct with any type of value for string manipulation.
//
// Supported input types:
//   - string: Direct string value
//   - *string: Pointer to string (allows in-place modification with Apply())
//   - []string: Slice of strings (use Join() to combine)
//   - int, int8, int16, int32, int64: Integer types
//   - uint, uint8, uint16, uint32, uint64: Unsigned integer types
//   - float32, float64: Floating point types
//   - bool: Boolean values (true/false)
//   - any other type: Converted to string representation
//
// Usage patterns:
//
// 1. Basic string manipulation:
//     result := Convert("hello world").ToUpper().String()
//     // result: "HELLO WORLD"
//
// 2. In-place modification of string pointer:
//     original := "hello world"
//     Convert(&original).ToUpper().Apply()
//     // original is now: "HELLO WORLD"
//
// 3. Working with slices:
//     words := []string{"hello", "world"}
//     result := Convert(words).Join(" ").ToUpper().String()
//     // result: "HELLO WORLD"
//
// 4. Numeric conversions:
//     result := Convert(42).String()
//     // result: "42"
//
// The Convert function returns a *Text instance that supports method chaining
// for various string transformations like case conversion, joining, parsing, etc.
// Use String() to get the final result or Apply() to modify the original pointer.
func Convert(v any) *Text {
	switch val := v.(type) {
	case []string:
		return &Text{contentSlice: val}
	case *string:
		// Handle string pointer directly without creating a new allocation
		// Store both the content and a reference to the original pointer
		return &Text{content: *val, stringPtr: val}
	default:
		return &Text{content: anyToString(v)}
	}
}

func (t *Text) transformWithMapping(mappings []charMapping) *Text {
	// Use a strings.Builder for efficient string construction
	var builder tinyStringBuilder
	// Pre-allocate builder with a reasonable estimate of the final string length
	builder.grow(len(t.content))

	for _, r := range t.content {
		mapped := false
		for _, mapping := range mappings {
			if r == mapping.from {
				builder.writeRune(mapping.to)
				mapped = true
				break
			}
		}
		if !mapped {
			builder.writeRune(r)
		}
	}
	t.content = builder.string()
	return t
}

// Remueve tildes y diacríticos
func (t *Text) RemoveTilde() *Text {
	return t.transformWithMapping(accentMappings)
}

// convert to lower case eg: "HELLO WORLD" -> "hello world"
func (t *Text) ToLower() *Text {
	return t.transformWithMapping(lowerMappings)
}

// convert to upper case eg: "hello world" -> "HELLO WORLD"
func (t *Text) ToUpper() *Text {
	return t.transformWithMapping(upperMappings)
}

// converts text to camelCase (first word lowercase) eg: "Hello world" -> "helloWorld"
func (t *Text) CamelCaseLower() *Text {
	return t.toCaseTransform(true, "")
}

// converts text to PascalCase (all words capitalized) eg: "hello world" -> "HelloWorld"
func (t *Text) CamelCaseUpper() *Text {
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
// ToSnakeCaseLower converts text to snake_case format
func (t *Text) ToSnakeCaseLower(sep ...string) *Text {
	return t.toCaseTransform(true, t.separatorCase(sep...))
}

// ToSnakeCaseUpper converts text to Snake_Case format
func (t *Text) ToSnakeCaseUpper(sep ...string) *Text {
	return t.toCaseTransform(false, t.separatorCase(sep...))
}

func (t *Text) separatorCase(sep ...string) string {
	t.separator = "_" // underscore default
	if len(sep) > 0 {
		t.separator = sep[0]
	}
	return t.separator
}

// Apply updates the original string pointer with the current content.
// This method should be used when you want to modify the original string directly
// without additional allocations.
func (t *Text) Apply() {
	// If contentSlice exists but not yet joined, join it with a space
	if t.contentSlice != nil && t.content == "" {
		t.content = t.Join().content
	}

	// If we have a string pointer, update it with the final content
	if t.stringPtr != nil {
		*t.stringPtr = t.content
	}
}

// String method to return the content of the text without modifying any original pointers
func (t *Text) String() string {
	// If contentSlice exists but not yet joined, join it with a space
	if t.contentSlice != nil && t.content == "" {
		t.content = t.Join().content
	}

	return t.content
}

// StringError returns the content of the text along with any error that occurred during processing
func (t *Text) StringError() (string, error) {
	// If contentSlice exists but not yet joined, join it with a space
	if t.contentSlice != nil && t.content == "" {
		t.content = t.Join().content
	}

	return t.content, t.err
}

// splitIntoWordsLocal returns words as local variable without storing in struct field
// This avoids persistent memory allocation in the Text struct
func (t *Text) splitIntoWordsLocal() [][]rune {
	words := make([][]rune, 0)
	currentWord := make([]rune, 0, 64) // Pre-allocate with reasonable capacity

	// Iterate directly over the string instead of converting to []rune
	for _, r := range t.content {
		if r == ' ' {
			if len(currentWord) > 0 {
				// Create a copy of currentWord for words
				wordCopy := make([]rune, len(currentWord))
				copy(wordCopy, currentWord)
				words = append(words, wordCopy)
				currentWord = currentWord[:0] // Reset for reuse
			}
			continue
		}
		currentWord = append(currentWord, r)
	}

	if len(currentWord) > 0 {
		// Create a copy of currentWord for words
		wordCopy := make([]rune, len(currentWord))
		copy(wordCopy, currentWord)
		words = append(words, wordCopy)
	}

	return words
}

func (t *Text) transformWord(word []rune, transform wordTransform) []rune {
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

func (t *Text) toCaseTransform(firstWordLower bool, separator string) *Text {
	// Use local variable instead of struct field to avoid persistent allocation
	words := t.splitIntoWordsLocal()
	if len(words) == 0 {
		return t
	}

	// Use string builder for efficient string construction
	estimatedLen := len(t.content) + len(words)*len(separator)
	builder := newTinyStringBuilder(estimatedLen)
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

	t.content = builder.string()
	// Clear the separator field after use to avoid memory overhead
	t.separator = ""
	return t
}
