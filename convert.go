package tinytext

type wordTransform int

const (
	toLower wordTransform = iota
	toUpper
	keepCase
)

// initialize the text struct
func Convert(s string) *Text {
	return &Text{content: s}
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
	return t.toCaseTransform(true)
}

// converts text to PascalCase (all words capitalized) eg: "hello world" -> "HelloWorld"
func (t *Text) CamelCaseUpper() *Text {
	return t.toCaseTransform(false)
}

// snakeCase converts a string to snake_case format with optional separator.
// If no separator is provided, underscore "_" is used as default.
// Example:
//
//	Input: "camelCase" -> Output: "camel_case"
//	Input: "PascalCase", "-" -> Output: "pascal-case"
//	Input: "APIResponse" -> Output: "api_response"
//	Input: "user123Name", "." -> Output: "user123.name"
func (t *Text) ToSnakeCase(sep ...string) *Text {
	separator := "_"
	if len(sep) > 0 {
		separator = sep[0]
	}

	runes := []rune(t.content)
	var result []rune

	for i, r := range runes {
		if r == ' ' {
			// Replace space with separator
			result = append(result, []rune(separator)...)
			continue
		}

		// Check if we need to add separator for uppercase letters
		if i > 0 && t.isUpperOrDigit(r) && t.isLowerOrDigit(runes[i-1]) {
			result = append(result, []rune(separator)...)
		}

		// Convert to lowercase
		result = append(result, t.transformWord([]rune{r}, toLower)...)
	}

	t.content = string(result)
	return t
}

// String method to return the content of the text
func (t *Text) String() string {
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

func (t *Text) toCaseTransform(firstWordLower bool) *Text {
	t.splitIntoWords()
	if len(t.words) == 0 {
		return t
	}

	var result []rune

	// Process first word
	firstWord := t.words[0]
	if len(firstWord) > 0 {
		transform := toUpper
		if firstWordLower {
			transform = toLower
		}
		// Procesar carácter por carácter para manejar números
		for i, r := range firstWord {
			if i == 0 {
				result = append(result, t.transformWord([]rune{r}, transform)...)
			} else if i > 0 && isDigit(firstWord[i-1]) {
				// Si viene después de un número, mantener mayúscula
				result = append(result, t.transformWord([]rune{r}, toUpper)...)
			} else {
				result = append(result, t.transformWord([]rune{r}, toLower)...)
			}
		}
	}

	// Process remaining words (always capitalize first letter)
	for _, word := range t.words[1:] {
		if len(word) > 0 {
			// Procesar carácter por carácter
			for i, r := range word {
				if i == 0 {
					// Primera letra siempre mayúscula
					result = append(result, t.transformWord([]rune{r}, toUpper)...)
				} else if i > 0 && isDigit(word[i-1]) {
					// Si viene después de un número, mantener mayúscula
					result = append(result, t.transformWord([]rune{r}, toUpper)...)
				} else {
					// En otro caso, minúscula
					result = append(result, t.transformWord([]rune{r}, toLower)...)
				}
			}
		}
	}

	t.content = string(result)
	return t
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

func (t *Text) isUpperOrDigit(r rune) bool {
	for _, mapping := range upperMappings {
		if r == mapping.from {
			return true
		}
	}
	return isDigit(r)
}

func (t *Text) isLowerOrDigit(r rune) bool {
	for _, mapping := range lowerMappings {
		if r == mapping.from {
			return true
		}
	}
	return isDigit(r)
}

// Helper function to check if a rune is a digit
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
