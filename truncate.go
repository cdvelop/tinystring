package tinystring

// Using shared constants from mapping.go for consistency

// Truncate truncates a conv so that it does not exceed the specified width.
// If the conv is longer, it truncates it and adds "..." if there is space.
// If the conv is shorter or equal to the width, it remains unchanged.
// The reservedChars parameter indicates how many characters should be reserved for suffixes.
// This parameter is optional - if not provided, no characters are reserved (equivalent to passing 0).
// eg: Convert("Hello, World!").Truncate(10) => "Hello, ..."
// eg: Convert("Hello, World!").Truncate(10, 3) => "Hell..."
// eg: Convert("Hello").Truncate(10) => "Hello"
func (t *conv) Truncate(maxWidth any, reservedChars ...any) *conv {
	if t.err != "" {
		return t // Error chain interruption
	}

	conv := t.getString()
	oL := len(conv)
	// Validate maxWidth parameter
	mWI, ok := t.validateIntParam(maxWidth, false)
	if !ok {
		return t
	}

	if oL > mWI {
		// Get reserved chars value
		rCI := 0
		if len(reservedChars) > 0 {
			if val, ok := t.validateIntParam(reservedChars[0], true); ok {
				rCI = val
			}
		}
		// Ensure rCI does not exceed mWI
		if rCI > mWI {
			rCI = mWI
		}
		// Calculate the width available for the conv itself, excluding reserved chars
		eW := max(mWI-rCI, 0)
		ellipsisLen := len(ellipsisStr)
		if rCI > 0 && mWI >= ellipsisLen && eW >= ellipsisLen {
			// Case 1: Reserved chars specified, and ellipsis fits within the effective width
			cTK := min(max(eW-ellipsisLen, 0), oL)
			// Phase 11: Use buffer instead of string concatenation to avoid allocation
			t.buf = t.getReusableBuffer(cTK + len(ellipsisStr))
			t.buf = append(t.buf, conv[:cTK]...)
			t.buf = append(t.buf, ellipsisStr...)
			t.setStringFromBuffer()
		} else if rCI == 0 && mWI >= ellipsisLen {
			// Case 2: No reserved chars, ellipsis fits within maxWidth
			cTK := min(max(mWI-ellipsisLen, 0), oL)
			// Phase 11: Use buffer instead of string concatenation to avoid allocation
			t.buf = t.getReusableBuffer(cTK + len(ellipsisStr))
			t.buf = append(t.buf, conv[:cTK]...)
			t.buf = append(t.buf, ellipsisStr...)
			t.setStringFromBuffer()
		} else {
			// Case 3: Ellipsis doesn't fit or reserved chars prevent it, just truncate
			cTK := min(mWI, oL)
			// Update buffer instead of using setString for buffer-first strategy
			t.buf = append(t.buf[:0], conv[:cTK]...)
		}
	}

	return t
}

// TruncateName truncates names and surnames in a user-friendly way for display in limited spaces
// like chart labels. It adds abbreviation dots where appropriate. This method processes the first
// word differently if there are more than 2 words in the conv.
//
// Parameters:
//   - maxCharsPerWord: maximum number of characters to keep per word (any numeric type)
//   - maxWidth: maximum total length for the final string (any numeric type)
//
// Examples:
//   - Convert("Jeronimo Dominguez").TruncateName(3, 15) => "Jer. Dominguez"
//   - Convert("Ana Maria Rodriguez").TruncateName(2, 10) => "An. Mar..."
//   - Convert("Juan").TruncateName(3, 5) => "Juan"
func (t *conv) TruncateName(maxCharsPerWord, maxWidth any) *conv {
	if t.err != "" {
		return t // Error chain interruption
	}

	if len(t.getString()) == 0 {
		return t
	}
	// Validate parameters
	mC, ok := t.validateIntParam(maxCharsPerWord, false)
	if !ok {
		return t
	}

	mT, ok := t.validateIntParam(maxWidth, false)
	if !ok {
		return t
	}

	words := Split(t.getString())
	if len(words) == 0 {
		return t
	} // Step 1: Apply maxCharsPerWord rule to each word
	var res string
	for i, word := range words {
		if i > 0 {
			res += spaceStr // Add space separator
		}
		// Inline processWordForName logic
		var processedWord string
		if i < len(words)-1 && len(word) > mC {
			processedWord = word[:mC] + dotStr
		} else if i == 0 && len(word) == 1 {
			// Special case: single letter first word gets a period
			processedWord = word + dotStr
		} else {
			processedWord = word
		}
		res += processedWord
	}
	// Step 2: Check if the processed result fits within maxWidth
	if len(res) <= mT {
		// Update buffer instead of using setString for buffer-first strategy
		t.buf = append(t.buf[:0], res...)
		return t
	}

	// Step 3: Apply maxWidth constraint with ellipsis - inline applyMaxWidthConstraint logic
	// Check if we can fit at least two words with abbreviations
	if len(words) > 1 {
		// Calculate minimum space needed for normal abbreviation pattern
		minNeeded := mC + 1 + 1 + min(mC+1, len(words[1])) // "Abc. D..." pattern
		if len(words) > 2 {
			minNeeded = mC + 1 + 1 + mC + 1 // "Abc. D..." for 3+ words
		} // If we can't fit the normal pattern, use all space for first word
		if mT < minNeeded && mT >= 4 { // minimum "X..." is 4 chars
			availableForFirstWord := mT - len(ellipsisStr)
			if len(words[0]) > availableForFirstWord {
				// Phase 11: Use buffer instead of string concatenation to avoid allocation
				t.buf = t.getReusableBuffer(availableForFirstWord + len(ellipsisStr))
				t.buf = append(t.buf, words[0][:availableForFirstWord]...)
				t.buf = append(t.buf, ellipsisStr...)
				t.setStringFromBuffer()
				return t
			}
		}
	}
	// Build result with remaining space tracking
	var result string
	remaining := mT - len(ellipsisStr) // Reserve space for "..." suffix

	for i, word := range words { // Check if we need to add a space
		if i > 0 {
			if remaining > 0 {
				result += spaceStr
				remaining--
			} else {
				break // No more space left
			}
		} // Inline processWordForName logic
		var prW string
		if i < len(words)-1 && len(word) > mC {
			prW = word[:mC] + dotStr
		} else if i == 0 && len(word) == 1 {
			// Special case: single letter first word gets a period
			prW = word + dotStr
		} else {
			prW = word
		}

		// Check how much of this word we can include
		if len(prW) <= remaining {
			// We can include the entire word
			result += prW
			remaining -= len(prW)
		} else {
			// We can only include part of the word
			result += prW[:remaining]
			remaining = 0
			break
		}
	}
	// Add the suffix
	result += ellipsisStr
	// Update buffer instead of using setString for buffer-first strategy
	t.buf = append(t.buf[:0], result...)
	return t
}
