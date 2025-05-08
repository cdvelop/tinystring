package tinystring

// Text struct to store the content of the text
type Text struct {
	content      string
	contentSlice []string // slice of strings for the Join method
	words        [][]rune // words split into runes eg: "hello world" -> [][]rune{{'h','e','l','l','o'}, {'w','o','r','l','d'}}
	separator    string   // eg "_" "-"
	stringPtr    *string  // pointer to original string (if one was provided)
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

// initialize the text struct with any type of value
// supports string, *string, int, float, bool, []string and their variants
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
	runes := []rune(t.content)
	for i, r := range runes {
		for _, mapping := range mappings {
			if r == mapping.from {
				runes[i] = mapping.to
				break
			}
		}
	}
	t.content = string(runes)
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

// String method to return the content of the text
func (t *Text) String() string {
	// If contentSlice exists but not yet joined, join it with a space
	if t.contentSlice != nil && t.content == "" {
		t.content = t.Join().content
	}

	// If we have a string pointer, update it with the final content
	if t.stringPtr != nil {
		*t.stringPtr = t.content
	}

	return t.content
}

// splitIntoWords splits the text content into a slice of words, where each word is represented
// as a slice of runes. Words are separated by spaces, and empty spaces are ignored.
//
// Example:
//
//	text := NewText("hello world")
//	words := text.splitIntoWords()
//	// Returns: [][]rune{{'h','e','l','l','o'}, {'w','o','r','l','d'}}
//
//	text := NewText("  multiple   spaces  ")
//	words := text.splitIntoWords()
//	// Returns: [][]rune{{'m','u','l','t','i','p','l','e'}, {'s','p','a','c','e','s'}}
func (t *Text) splitIntoWords() {
	t.words = make([][]rune, 0)
	var currentWord []rune

	// Iterate directly over the string instead of converting to []rune
	for _, r := range t.content {
		if r == ' ' {
			if len(currentWord) > 0 {
				t.words = append(t.words, currentWord)
				currentWord = nil
			}
			continue
		}
		currentWord = append(currentWord, r)
	}

	if len(currentWord) > 0 {
		t.words = append(t.words, currentWord)
	}
}

func (t *Text) transformWord(word []rune, transform wordTransform) []rune {
	if len(word) == 0 {
		return word
	}

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

	return result
}

// Helper function to check if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= 'À' && r <= 'ÿ' && r != '×' && r != '÷')
}

func (t *Text) toCaseTransform(firstWordLower bool, separator string) *Text {

	t.splitIntoWords()
	if len(t.words) == 0 {
		return t
	}

	var result []rune
	var prevIsDigit bool
	var prevIsSeparator bool

	for i, word := range t.words {
		if len(word) == 0 {
			continue
		}

		// Add separator if needed
		if i > 0 && separator != "" {
			result = append(result, rune(separator[0]))
			prevIsSeparator = true
		}

		// Process each character in the word
		for j, r := range word {
			transform := toLower
			currIsDigit := isDigit(r)
			currIsLetter := isLetter(r)

			// Determine case transform
			if i == 0 && j == 0 {
				// First letter of first word
				if !firstWordLower {
					transform = toUpper
				}
			} else if i > 0 && j == 0 && separator == "" { // Start of new word in camelCase
				transform = toUpper
			} else if prevIsDigit && currIsLetter { // Letter after digit
				if firstWordLower {
					transform = toLower
				} else {
					transform = toUpper
				}
			} else if prevIsSeparator && currIsLetter { // Letter after separator
				if separator != "" && !firstWordLower {
					transform = toUpper
				}
			}

			// Add underscore for number to letter transition in snake_case
			if separator != "" && prevIsDigit && currIsLetter {
				result = append(result, rune(separator[0]))
			}

			// Only transform letters, leave other characters as-is
			if currIsLetter {
				result = append(result, t.transformWord([]rune{r}, transform)...)
			} else {
				result = append(result, r)
			}

			prevIsDigit = currIsDigit
			prevIsSeparator = false
		}
	}

	t.content = string(result)
	return t
}
