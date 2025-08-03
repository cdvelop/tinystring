package tinystring

// Capitalize transforms the first letter of each word to uppercase and the rest to lowercase.
// Preserves all whitespace formatting (spaces, tabs, newlines) without normalization.
// OPTIMIZED: Uses work buffer efficiently to minimize allocations
// For example: "  hello   world  " -> "  Hello   World  "
func (t *conv) Capitalize() *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}

	if t.outLen == 0 {
		return t
	}

	// Fast path for ASCII-only content (common case)
	if t.isASCIIOnly() {
		t.capitalizeASCIIOptimized()
		return t
	}

	// Unicode fallback
	return t.capitalizeUnicode()
}

// capitalizeASCIIOptimized processes ASCII text preserving all formatting
func (t *conv) capitalizeASCIIOptimized() {
	// Use work buffer for processing
	t.rstBuffer(buffWork)

	inWord := false

	for i := 0; i < t.outLen; i++ {
		ch := t.out[i]

		// Use centralized word separator detection
		if isWordSeparator(ch) {
			// Preserve all separator characters as-is
			t.work = append(t.work, ch)
			t.workLen++
			inWord = false
		} else {
			if !inWord {
				// Start of new word - capitalize first letter
				if ch >= 'a' && ch <= 'z' {
					ch -= 32 // Convert to uppercase
				}
				inWord = true
			} else {
				// Rest of word - lowercase other letters
				if ch >= 'A' && ch <= 'Z' {
					ch += 32 // Convert to lowercase
				}
			}
			t.work = append(t.work, ch)
			t.workLen++
		}
	}

	// Swap processed content to output
	t.swapBuff(buffWork, buffOut)
}

// capitalizeUnicode handles full Unicode capitalization preserving formatting
func (t *conv) capitalizeUnicode() *conv {
	str := t.getString(buffOut)

	// Use internal work buffer for intermediate processing
	t.rstBuffer(buffWork)

	inWord := false

	for _, r := range str {
		// Use centralized word separator detection
		if isWordSeparator(r) {
			// Preserve all separator characters as-is
			t.wrString(buffWork, string(r))
			inWord = false
		} else {
			if !inWord {
				// Start of new word - capitalize first letter
				t.wrString(buffWork, string(toUpperRune(r)))
				inWord = true
			} else {
				// Rest of word - lowercase other letters
				t.wrString(buffWork, string(toLowerRune(r)))
			}
		}
	}

	// Copy result from work buffer to output buffer
	result := t.getString(buffWork)
	t.rstBuffer(buffOut)
	t.wrString(buffOut, result)
	return t
}

// convert to lower case eg: "HELLO WORLD" -> "hello world"
func (t *conv) Low() *conv {
	return t.changeCaseOptimized(true)
}

// convert to upper case eg: "hello world" -> "HELLO WORLD"
func (t *conv) Up() *conv {
	return t.changeCaseOptimized(false)
}

// changeCaseOptimized implements fast ASCII path with fallback to full Unicode
func (t *conv) changeCaseOptimized(toLower bool) *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}

	if t.outLen == 0 {
		return t
	}

	// Fast path: ASCII-only optimization (covers 85% of use cases)
	if t.isASCIIOnly() {
		t.changeCaseASCIIInPlace(toLower)
		return t
	}

	// Fallback: Full Unicode support for complex cases
	return t.changeCaseUnicode(toLower)
}

// changeCaseASCIIInPlace processes ASCII characters directly in buffer (zero allocations)
func (t *conv) changeCaseASCIIInPlace(toLower bool) {
	for i := 0; i < t.outLen; i++ {
		if toLower {
			// A-Z (65-90) → a-z (97-122): add 32
			if t.out[i] >= 'A' && t.out[i] <= 'Z' {
				t.out[i] += 32
			}
		} else {
			// a-z (97-122) → A-Z (65-90): subtract 32
			if t.out[i] >= 'a' && t.out[i] <= 'z' {
				t.out[i] -= 32
			}
		}
	}
}

// isASCIIOnly checks if buffer contains only ASCII characters (fast check)
func (t *conv) isASCIIOnly() bool {
	for i := 0; i < t.outLen; i++ {
		if t.out[i] > 127 {
			return false
		}
	}
	return true
}

// changeCaseUnicode handles full Unicode case conversion (legacy method)
func (t *conv) changeCaseUnicode(toLower bool) *conv {
	return t.changeCase(toLower, buffOut)
}

// changeCase consolidates Low and Up functionality - now accepts a destination buffer for internal reuse
func (t *conv) changeCase(toLower bool, dest buffDest) *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}

	str := t.getString(dest)
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
	t.rstBuffer(dest)     // Clear buffer using API
	t.wrString(dest, out) // Write using API

	return t
}

// converts conv to camelCase (first word lowercase) eg: "Hello world" -> "helloWorld"
func (t *conv) CamelLow() *conv {
	return t.toCaseTransformMinimal(true, "")
}

// converts conv to PascalCase (all words capitalized) eg: "hello world" -> "HelloWorld"
func (t *conv) CamelUp() *conv {
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
// SnakeLow converts conv to snake_case format
func (t *conv) SnakeLow(sep ...string) *conv {
	// Phase 4.3: Use local variable instead of struct field
	separator := "_" // underscore default
	if len(sep) > 0 {
		separator = sep[0]
	}
	return t.toCaseTransformMinimal(true, separator)
}

// SnakeUp converts conv to Snake_Case format
func (t *conv) SnakeUp(sep ...string) *conv {
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

	str := t.getString(buffOut)
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
			// - For PascalCase (CamelUp): start new word
			// - For camelCase (CamelLow): don't start new word
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
