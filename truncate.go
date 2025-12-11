package fmt

// Using shared constants from mapping.go for consistency

// validateIntParam validates and converts any numeric type to int
// Universal method that follows buffer API architecture - eliminates code duplication
func (c *Conv) validateIntParam(param any, allowZero bool) (int, bool) {
	var val int
	var ok bool
	switch v := param.(type) {
	case int, int8, int16, int32, int64:
		// Use type assertion to handle all integer types
		if i, isInt := param.(int); isInt {
			val, ok = i, true
		} else if i8, isInt8 := param.(int8); isInt8 {
			val, ok = int(i8), true
		} else if i16, isInt16 := param.(int16); isInt16 {
			val, ok = int(i16), true
		} else if i32, isInt32 := param.(int32); isInt32 {
			val, ok = int(i32), true
		} else if i64, isInt64 := param.(int64); isInt64 {
			val, ok = int(i64), true
		}
	case uint, uint8, uint16, uint32, uint64:
		if u, isUint := param.(uint); isUint {
			val, ok = int(u), true
		} else if u8, isUint8 := param.(uint8); isUint8 {
			val, ok = int(u8), true
		} else if u16, isUint16 := param.(uint16); isUint16 {
			val, ok = int(u16), true
		} else if u32, isUint32 := param.(uint32); isUint32 {
			val, ok = int(u32), true
		} else if u64, isUint64 := param.(uint64); isUint64 {
			val, ok = int(u64), true
		}
	case float32:
		val, ok = int(v), true
	case float64:
		val, ok = int(v), true
	default:
		val, ok = 0, false
	}

	if !ok {
		return 0, false
	}
	// Unified validation logic
	if allowZero {
		return val, val >= 0
	}
	return val, val > 0
}

// truncateWithEllipsis helper method to reduce code duplication
// Handles the common pattern of truncating content and adding ellipsis
func (c *Conv) truncateWithEllipsis(content string, maxWidth int) {
	ellipsisLen := len(ellipsisStr)
	if maxWidth >= ellipsisLen {
		contentToKeep := min(max(maxWidth-ellipsisLen, 0), len(content))
		c.ResetBuffer(BuffOut)                       // Clear buffer using API
		c.WrString(BuffOut, content[:contentToKeep]) // Write content using API
		c.WrString(BuffOut, ellipsisStr)             // Append ellipsis using API
	} else {
		// Ellipsis doesn't fit, just truncate
		contentToKeep := min(maxWidth, len(content))
		c.ResetBuffer(BuffOut)                       // Clear buffer using API
		c.WrString(BuffOut, content[:contentToKeep]) // Write using API
	}
}

// Truncate truncates a Conv so that it does not exceed the specified width.
// If the Conv is longer, it truncates it and adds "..." if there is space.
// If the Conv is shorter or equal to the width, it remains unchanged.
// The reservedChars parameter indicates how many characters should be reserved for suffixes.
// This parameter is optional - if not provided, no characters are reserved (equivalent to passing 0).
// eg: Convert("Hello, World!").Truncate(10) => "Hello, ..."
// eg: Convert("Hello, World!").Truncate(10, 3) => "Hell..."
// eg: Convert("Hello").Truncate(10) => "Hello"
func (t *Conv) Truncate(maxWidth any, reservedChars ...any) *Conv {
	if t.hasContent(BuffErr) {
		return t // Error chain interruption
	}

	// OPTIMIZED: Direct buffer processing
	if t.outLen == 0 {
		return t
	}

	// Validate maxWidth parameter
	mWI, ok := t.validateIntParam(maxWidth, false)
	if !ok {
		return t
	}

	// OPTIMIZED: Use direct buffer length check
	if t.outLen > mWI {
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
		} // Calculate the width available for the Conv itself, excluding reserved chars
		eW := max(mWI-rCI, 0)
		ellipsisLen := len(ellipsisStr)
		if rCI > 0 && mWI >= ellipsisLen && eW >= ellipsisLen {
			// Case 1: Reserved chars specified, and ellipsis fits within the effective width
			// Need string for ellipsis methods - fallback to GetString
			Conv := t.GetString(BuffOut)
			t.truncateWithEllipsis(Conv, eW)
		} else if rCI == 0 && mWI >= ellipsisLen {
			// Case 2: No reserved chars, ellipsis fits within maxWidth
			// Need string for ellipsis methods - fallback to GetString
			Conv := t.GetString(BuffOut)
			t.truncateWithEllipsis(Conv, mWI)
		} else {
			// Case 3: Ellipsis doesn't fit or reserved chars prevent it, just truncate
			// OPTIMIZED: Direct buffer truncation
			cTK := min(mWI, t.outLen)
			t.outLen = cTK
			t.out = t.out[:cTK]
		}
	}

	return t
}

// TruncateName truncates names and surnames in a user-friendly way for display in limited spaces
// like chart labels. It adds abbreviation dots where appropriate. This method processes the first
// word differently if there are more than 2 words in the Conv.
//
// Parameters:
//   - maxCharsPerWord: maximum number of characters to keep per word (any numeric type)
//   - maxWidth: maximum total length for the final string (any numeric type)
//
// Examples:
//   - Convert("Jeronimo Dominguez").TruncateName(3, 15) => "Jer. Dominguez"
//   - Convert("Ana Maria Rodriguez").TruncateName(2, 10) => "An. Mar..."
//   - Convert("Juan").TruncateName(3, 5) => "Juan"
func (t *Conv) TruncateName(maxCharsPerWord, maxWidth any) *Conv {
	if t.hasContent(BuffErr) {
		return t // Error chain interruption
	}

	// OPTIMIZED: Direct length check
	if t.outLen == 0 {
		return t
	} // Validate parameters
	mC, ok := t.validateIntParam(maxCharsPerWord, false)
	if !ok {
		return t
	}

	mT, ok := t.validateIntParam(maxWidth, false)
	if !ok {
		return t
	}

	words := t.Split()
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
	} // Step 2: Check if the processed out fits within maxWidth
	if len(res) <= mT {
		// ✅ Update buffer using API instead of direct manipulation
		t.ResetBuffer(BuffOut)   // Clear buffer using API
		t.WrString(BuffOut, res) // Write using API
		return t
	}

	// Step 3: Apply maxWidth constraint with ellipsis - inline applyMaxWidthConstraint logic
	// Check if we can fit at least two words with abbreviations
	if len(words) > 1 {
		// Calculate minimum space needed for normal abbreviation pattern
		minNeeded := mC + 1 + 1 + min(mC+1, len(words[1])) // "Abc. D..." pattern
		if len(words) > 2 {
			minNeeded = mC + 1 + 1 + mC + 1 // "Abc. D..." for 3+ words
		}
		// If we can't fit the normal pattern, use all space for first word
		if mT < minNeeded && mT >= 4 { // minimum "X..." is 4 chars
			if len(words[0]) > mT-len(ellipsisStr) {
				t.truncateWithEllipsis(words[0], mT)
				return t
			}
		}
	}
	// Build out with remaining space tracking
	var out string
	remaining := mT - len(ellipsisStr) // Reserve space for "..." suffix

	for i, word := range words { // Check if we need to add a space
		if i > 0 {
			if remaining > 0 {
				out += spaceStr
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
			out += prW
			remaining -= len(prW)
		} else {
			// We can only include part of the word
			out += prW[:remaining]
			remaining = 0
			break
		}
	} // Add the suffix
	out += ellipsisStr
	// ✅ Update buffer using API instead of direct manipulation
	t.ResetBuffer(BuffOut)   // Clear buffer using API
	t.WrString(BuffOut, out) // Write using API
	return t
}
