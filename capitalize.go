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
	words := make([][]rune, 0, 4)      // estimate 4 words
	currentWord := make([]rune, 0, 10) // estimate 10 chars per word

	for _, r := range str {
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
			word[0] = toUpperRune(word[0])
			for j := 1; j < len(word); j++ {
				word[j] = toLowerRune(word[j])
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
		runes[i] = toLowerRune(r)
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
		runes[i] = toUpperRune(r)
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
	var pWU, pWL, pWD, pWS bool
	for i, r := range str {
		cIU := isLetter(r) && isUpper(r)
		cIL := isLetter(r) && isLower(r)
		cID := isDigit(r)
		cIS := r == ' ' || r == '\t' || r == '\n' || r == '\r'

		// Determine if we're starting a new word
		iWS := false
		if i == 0 {
			iWS = true
		} else if cIS {
			// Skip spaces but mark that we had a space
			pWS = true
			continue
		} else if pWS {
			// After space - new word
			iWS = true
			pWS = false
		} else if pWL && cIU {
			// camelCase transition: "camelCase" -> "camel" + "Case"
			iWS = true
		} else if pWD && (cIU || cIL) {
			// Digit to letter transition:
			// - For snake_case: always start new word
			// - For PascalCase (CamelCaseUpper): start new word
			// - For camelCase (CamelCaseLower): don't start new word
			if separator != "" || !firstWordLower {
				iWS = true
			}
		} else if (pWU || pWL) && cID {
			// Letter to digit: no new word - numbers continue the word
		}

		// Add separator if starting new word (except first word)
		if iWS && wordIndex > 0 && separator != "" {
			result = append(result, separator...)
		}

		// Determine case transformation
		var transformedRune rune
		if iWS {
			// First letter of word
			if wordIndex == 0 && firstWordLower {
				// First word lowercase (camelCase)
				transformedRune = toLowerRune(r)
			} else if separator != "" && firstWordLower {
				// snake_case_lower - all words lowercase
				transformedRune = toLowerRune(r)
			} else {
				// PascalCase, camelCase subsequent words, or Snake_Case_Upper
				transformedRune = toUpperRune(r)
			}
			wordIndex++
		} else {
			// Rest of letters in word - always lowercase
			transformedRune = toLowerRune(r)
		}

		// Add the character
		result = append(result, string(transformedRune)...)

		// Update state for next iteration
		pWU = cIU
		pWL = cIL
		pWD = cID
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
