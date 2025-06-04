package tinystring

// Capitalize transforms the first letter of each word to uppercase and the rest to lowercase.
// Also normalizes whitespace (collapses multiple spaces into single space and trims).
// For example: "  hello   world  " -> "Hello World"
func (t *conv) Capitalize() *conv {
	str := t.getString()
	if len(str) == 0 {
		return t
	}

	// First pass: normalize whitespace and build word list
	runes := []rune(str)
	words := make([][]rune, 0, 4) // estimate 4 words
	currentWord := make([]rune, 0, 10) // estimate 10 chars per word
	
	for _, r := range runes {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if len(currentWord) > 0 {
				words = append(words, currentWord)
				currentWord = make([]rune, 0, 10)
			}
		} else {
			currentWord = append(currentWord, r)
		}
	}
	if len(currentWord) > 0 {
		words = append(words, currentWord)
	}
	
	if len(words) == 0 {
		t.setString("")
		return t
	}

	// Second pass: capitalize words and calculate total length
	totalLen := 0
	for i, word := range words {
		// Capitalize first letter, lowercase the rest
		if len(word) > 0 {
			transformedRune, mapped := transformSingleRune(word[0], upperMappings)
			if mapped {
				word[0] = transformedRune
			}
			for j := 1; j < len(word); j++ {
				transformedRune, mapped := transformSingleRune(word[j], lowerMappings)
				if mapped {
					word[j] = transformedRune
				}
			}
			totalLen += len(word)
			if i > 0 {
				totalLen++ // space between words
			}
		}
	}

	// Third pass: build final string
	buf := make([]rune, 0, totalLen)
	for i, word := range words {
		if i > 0 {
			buf = append(buf, ' ')
		}
		buf = append(buf, word...)
	}

	t.setString(string(buf))
	return t
}

// convert to lower case eg: "HELLO WORLD" -> "hello world"
func (t *conv) ToLower() *conv {
	str := t.getString()
	if len(str) == 0 {
		return t
	}

	// Optimized: use rune slice and single string conversion - Phase 3
	runes := []rune(str)
	for i, r := range runes {
		transformedRune, mapped := transformSingleRune(r, lowerMappings)
		if mapped {
			runes[i] = transformedRune
		}
	}

	t.setString(string(runes))
	return t
}

// convert to upper case eg: "hello world" -> "HELLO WORLD"
func (t *conv) ToUpper() *conv {
	str := t.getString()
	if len(str) == 0 {
		return t
	}

	// Optimized: use rune slice and single string conversion - Phase 3
	runes := []rune(str)
	for i, r := range runes {
		transformedRune, mapped := transformSingleRune(r, upperMappings)
		if mapped {
			runes[i] = transformedRune
		}
	}

	t.setString(string(runes))
	return t
}

// converts conv to camelCase (first word lowercase) eg: "Hello world" -> "helloWorld"
func (t *conv) CamelCaseLower() *conv {
	return t.toCaseTransformMinimal(true, "")
}

// converts conv to PascalCase (all words capitalized) eg: "hello world" -> "HelloWorld"
func (t *conv) CamelCaseUpper() *conv {
	return t.toCaseTransformMinimal(false, "")
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
	return t.toCaseTransformMinimal(true, t.separatorCase(sep...))
}

// ToSnakeCaseUpper converts conv to Snake_Case format
func (t *conv) ToSnakeCaseUpper(sep ...string) *conv {
	return t.toCaseTransformMinimal(false, t.separatorCase(sep...))
}

// Minimal implementation without pools or builders - optimized for minimal allocations
func (t *conv) toCaseTransformMinimal(firstWordLower bool, separator string) *conv {
	str := t.getString()
	if len(str) == 0 {
		return t
	}

	// Pre-allocate buffer with estimated size
	estimatedSize := len(str) + (len(separator) * 5) // Extra space for separators
	result := make([]byte, 0, estimatedSize)
	
	// Advanced word boundary detection for camelCase and snake_case
	wordIndex := 0
	var prevWasUpper, prevWasLower, prevWasDigit, prevWasSpace bool

	for i, r := range str {
		currIsUpper := isLetter(r) && isUpper(r)
		currIsLower := isLetter(r) && isLower(r)
		currIsDigit := isDigit(r)
		currIsSpace := r == ' ' || r == '\t' || r == '\n' || r == '\r'

		// Determine if we're starting a new word
		isWordStart := false
		if i == 0 {
			isWordStart = true
		} else if currIsSpace {
			// Skip spaces but mark that we had a space
			prevWasSpace = true
			continue
		} else if prevWasSpace {
			// After space - new word
			isWordStart = true
			prevWasSpace = false
		} else if prevWasLower && currIsUpper {
			// camelCase transition: "camelCase" -> "camel" + "Case"
			isWordStart = true
		} else if prevWasDigit && (currIsUpper || currIsLower) {
			// Digit to letter transition:
			// - For snake_case: always start new word
			// - For PascalCase (CamelCaseUpper): start new word
			// - For camelCase (CamelCaseLower): don't start new word
			if separator != "" || !firstWordLower {
				isWordStart = true
			}
		} else if (prevWasUpper || prevWasLower) && currIsDigit {
			// Letter to digit: no new word - numbers continue the word
		}

		// Add separator if starting new word (except first word)
		if isWordStart && wordIndex > 0 && separator != "" {
			result = append(result, separator...)
		}

		// Determine case transformation
		var transformedRune rune
		var mapped bool

		if isWordStart {
			// First letter of word
			if wordIndex == 0 && firstWordLower {
				// First word lowercase (camelCase)
				transformedRune, mapped = transformSingleRune(r, lowerMappings)
			} else if separator != "" && firstWordLower {
				// snake_case_lower - all words lowercase
				transformedRune, mapped = transformSingleRune(r, lowerMappings)
			} else {
				// PascalCase, camelCase subsequent words, or Snake_Case_Upper
				transformedRune, mapped = transformSingleRune(r, upperMappings)
			}
			wordIndex++
		} else {
			// Rest of letters in word - always lowercase
			transformedRune, mapped = transformSingleRune(r, lowerMappings)
		}

		// Add the character
		if isLetter(r) {
			if mapped {
				result = append(result, string(transformedRune)...)
			} else if isWordStart && !firstWordLower && separator == "" {
				// For camelCase/PascalCase, ensure first letter of non-first words is uppercase
				result = append(result, string(toUpperChar(r))...)
			} else if isWordStart && firstWordLower && separator == "" && wordIndex > 1 {
				// For camelCase, ensure first letter of subsequent words is uppercase
				result = append(result, string(toUpperChar(r))...)
			} else {
				result = append(result, string(toLowerChar(r))...)
			}
		} else {
			result = append(result, string(r)...)
		}

		prevWasUpper = currIsUpper
		prevWasLower = currIsLower
		prevWasDigit = currIsDigit
	}

	t.setString(string(result))
	return t
}

// Helper functions for simple case conversion
func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func toUpperChar(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}

func toLowerChar(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + 32
	}
	return r
}
