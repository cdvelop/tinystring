package tinystring

// Capitalize transforms the first letter of each word to uppercase and the rest to lowercase.
// Also normalizes whitespace (collapses multiple spaces into single space and trims).
// Optimized to use internal work buffer instead of separate pool allocations
// For example: "  hello   world  " -> "Hello World"
func (t *conv) Capitalize() *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}

	str := t.getBuffString()
	if len(str) == 0 {
		return t
	}

	// Use internal work buffer for intermediate processing (follows Unified Buffer Architecture)
	t.rstBuffer(buffWork)

	inWord := false
	addSpace := false // Flag to add space before next word

	for _, r := range str {
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
				if addSpace && t.hasContent(buffWork) {
					t.wrString(buffWork, " ")
				}
				t.wrString(buffWork, string(toUpperRune(r)))
				inWord = true
				addSpace = false
			} else {
				// Lowercase other letters in word
				t.wrString(buffWork, string(toLowerRune(r)))
			}
		}
	}

	// Copy result from work buffer to output buffer using API
	result := t.getString(buffWork)
	t.rstBuffer(buffOut)        // Clear output buffer using API
	t.wrString(buffOut, result) // Write using API
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

// changeCase consolidates ToLower and ToUpper functionality - optimized with buffer-first strategy
func (t *conv) changeCase(toLower bool) *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}

	str := t.getBuffString()
	if len(str) == 0 {
		return t
	}

	// Convert to runes for proper Unicode handling
	runes := []rune(str)

	// Process runes for case conversion
	for i, r := range runes {
		if toLower {
			runes[i] = toLowerRune(r)
		} else {
			runes[i] = toUpperRune(r)
		}
	}
	// Convert back to string and store in buffer using API
	out := string(runes)
	t.rstBuffer(buffOut)     // Clear buffer using API
	t.wrString(buffOut, out) // Write using API

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
	// Phase 4.3: Use local variable instead of struct field
	separator := "_" // underscore default
	if len(sep) > 0 {
		separator = sep[0]
	}
	return t.toCaseTransformMinimal(true, separator)
}

// ToSnakeCaseUpper converts conv to Snake_Case format
func (t *conv) ToSnakeCaseUpper(sep ...string) *conv {
	// Phase 4.3: Use local variable instead of struct field
	separator := "_" // underscore default
	if len(sep) > 0 {
		separator = sep[0]
	}
	return t.toCaseTransformMinimal(false, separator)
}

// Minimal implementation without pools or builders - optimized for minimal allocations
func (t *conv) toCaseTransformMinimal(firstWordLower bool, separator string) *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}

	str := t.getBuffString()
	if len(str) == 0 {
		return t
	} // Pre-allocate buffer with estimated size
	eSz := len(str) + (len(separator) * 5) // Extra space for separators
	// Inline makeBuf logic
	resultCap := eSz
	if resultCap < defaultBufCap {
		resultCap = defaultBufCap
	}
	out := make([]byte, 0, resultCap)
	// Advanced word boundary detection for camelCase and snake_case
	wordIndex := 0
	var pWU, pWL, pWD, pWS bool
	for i, r := range str {
		// Inline isLetter logic
		isLetterR := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= 'À' && r <= 'ÿ' && r != 'x' && r != '÷')
		cIU := isLetterR && isUpper(r)
		cIL := isLetterR && isLower(r)
		// Inline isDigit logic
		cID := r >= '0' && r <= '9'
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
			out = append(out, separator...)
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
			out = append(out, byte(transformedRune))
		} else {
			// Use UTF-8 encoding for multi-byte runes
			out = append(out, string(transformedRune)...)
		}

		// Update state for next iteration
		pWU = cIU
		pWL = cIL
		pWD = cID
	}

	// ✅ Update buffer using API instead of direct manipulation
	t.rstBuffer(buffOut)    // Clear buffer using API
	t.wrBytes(buffOut, out) // Write bytes using API
	return t
}

// Helper functions for simple case conversion
func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}
