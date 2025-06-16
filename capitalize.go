package tinystring

// Capitalize transforms the first letter of each word to uppercase and the rest to lowercase.
// Also normalizes whitespace (collapses multiple spaces into single space and trims).
// Phase 11 Optimization: Reduced buffer allocations through pooling
// For example: "  hello   world  " -> "Hello World"
func (t *conv) Capitalize() *conv {
	str := t.getString()
	if isEmpty(str) {
		return t
	}

	// Phase 11: Use single buffer approach to reduce allocations
	runes := []rune(str)
	if len(runes) == 0 {
		return t
	}
	// Get pooled buffer for result
	buf := getRuneBuffer(len(runes))
	defer putRuneBuffer(&buf)
	inWord := false
	addSpace := false // Flag to add space before next word

	for _, r := range runes {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if inWord {
				// End of word, mark that we need a space before next word
				addSpace = true
				inWord = false
			}
			// Skip multiple whitespaces
		} else {
			if !inWord {
				// Start of new word
				if addSpace && len(buf) > 0 {
					buf = append(buf, ' ')
				}
				buf = append(buf, toUpperRune(r))
				inWord = true
				addSpace = false
			} else {
				// Lowercase other letters in word
				buf = append(buf, toLowerRune(r))
			}
		}
	}

	t.setString(string(buf))
	return t
}

// convert to lower case eg: "HELLO WORLD" -> "hello world"
func (t *conv) ToLower() *conv {
	return t.changeCase(true)
}

// convert to upper case eg: "hello world" -> "HELLO WORLD"
func (t *conv) ToUpper() *conv {
	return t.changeCase(false)
}

// changeCase consolidates ToLower and ToUpper functionality
func (t *conv) changeCase(toLower bool) *conv {
	str := t.getString()
	if isEmpty(str) {
		return t
	}

	// Use rune buffer pool for better memory efficiency
	buf := getRuneBuffer(len(str))
	defer putRuneBuffer(&buf)

	// Convert to runes and process
	for _, r := range str {
		if toLower {
			buf = append(buf, toLowerRune(r))
		} else {
			buf = append(buf, toUpperRune(r))
		}
	}

	t.setString(string(buf))
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
	if isEmpty(str) {
		return t
	}
	// Pre-allocate buffer with estimated size
	eSz := len(str) + (len(separator) * 5) // Extra space for separators
	result := makeBuf(eSz)
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
		// Add the character - use bytes buffer for efficiency
		// Convert rune to bytes directly instead of string(rune)
		if transformedRune < 128 {
			// ASCII optimization
			result = append(result, byte(transformedRune))
		} else {
			// Use UTF-8 encoding for multi-byte runes
			result = append(result, string(transformedRune)...)
		}

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
