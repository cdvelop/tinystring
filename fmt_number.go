package tinystring

// =============================================================================
// FORMAT NUMBER OPERATIONS - Number formatting with separators and display
// =============================================================================

// FormatNumber intelligently formats the current value as a number with thousand separators.
// Handles both integers and floats, preserving decimal places when appropriate.
// Example: Convert("1234567.89").FormatNumber() returns "1.234.567,89"
func (t *conv) FormatNumber() *conv {
	if t.hasContent(buffErr) {
		return t
	}

	// OPTIMIZED: Buffer-first approach - check if we already have content
	if t.hasContent(buffOut) {
		str := t.getString(buffOut)
		// Only add thousand separators if it's actually a number
		if t.isNumericString(str) {
			// For numeric strings, try to parse as float first to handle .00 cases
			floatVal := t.parseFloat(str)
			if !t.hasContent(buffErr) {
				t.rstBuffer(buffOut)
				// Check if it's effectively an integer (.00 case)
				if floatVal == float64(int64(floatVal)) {
					t.wrInt(buffOut, int64(floatVal))
				} else {
					t.wrFloat(buffOut, floatVal)
					t.removeTrailingZeros(buffOut)
				}
			}
			t.addThousandSeparators(buffOut)
		}
		return t
	}

	// Fallback: Use ptrValue for complex types that need lazy conversion
	if t.ptrValue != nil {
		// Clear any existing errors before processing
		t.rstBuffer(buffErr)
		t.rstBuffer(buffOut)
		t.anyToBuff(buffOut, t.ptrValue)

		// Check if conversion was successful
		if !t.hasContent(buffErr) {
			str := t.getString(buffOut)
			// Only add thousand separators if it's actually a number
			if t.isNumericString(str) {
				// For numeric strings, try to parse as float first to handle .00 cases
				floatVal := t.parseFloat(str)
				if !t.hasContent(buffErr) {
					t.rstBuffer(buffOut)
					// Check if it's effectively an integer (.00 case)
					if floatVal == float64(int64(floatVal)) {
						t.wrInt(buffOut, int64(floatVal))
					} else {
						t.wrFloat(buffOut, floatVal)
						t.removeTrailingZeros(buffOut)
					}
				}
				t.addThousandSeparators(buffOut)
			}
		}
		return t
	}

	// Get current string and attempt number formatting
	inp := t.getString(buffOut)

	// Save the current error state
	errState := t.hasContent(buffErr)

	// Try integer parsing first
	if intVal, err := t.parseSmallInt(inp); err == nil {
		t.rstBuffer(buffOut)
		t.wrInt(buffOut, int64(intVal))
		t.addThousandSeparators(buffOut)
		return t
	}

	// Reset error state before trying float parsing
	if !errState {
		t.rstBuffer(buffErr)
	}

	// Try float parsing - unified approach
	floatVal := t.parseFloat(inp)
	if !t.hasContent(buffErr) {
		t.rstBuffer(buffOut)
		// Check if it's effectively an integer (.00 case)
		if floatVal == float64(int64(floatVal)) {
			t.wrInt(buffOut, int64(floatVal))
		} else {
			t.wrFloat(buffOut, floatVal)
			t.removeTrailingZeros(buffOut)
		}
		t.addThousandSeparators(buffOut)
		return t
	}

	// Fallback: restore original string if not a number
	t.rstBuffer(buffOut)
	t.wrString(buffOut, inp)
	return t
}

// addThousandSeparators adds thousand separators to the numeric string in buffer.
// Universal method with dest-first parameter order - follows buffer API architecture
func (c *conv) addThousandSeparators(dest buffDest) {
	str := c.getString(dest)
	if len(str) <= 3 {
		return // No separators needed for numbers with 3 or fewer digits
	}

	// Find decimal point if it exists
	dotIndex := -1
	for i, char := range str {
		if char == '.' {
			dotIndex = i
			break
		}
	}

	// Determine the integer part
	intPart := str
	decPart := ""
	if dotIndex != -1 {
		intPart = str[:dotIndex]
		decPart = str[dotIndex:] // Include the decimal point
	}

	// Skip if integer part is too short or starts with negative sign and is too short
	intLen := len(intPart)
	if intPart[0] == '-' {
		if intLen <= 4 { // -123 or shorter
			return
		}
	} else {
		if intLen <= 3 { // 123 or shorter
			return
		}
	}

	// Build result with separators
	c.rstBuffer(dest)

	// Handle negative sign
	start := 0
	if intPart[0] == '-' {
		c.wrByte(dest, '-')
		start = 1
	}

	// Calculate number of full groups of 3
	remainingDigits := intLen - start
	firstGroupSize := remainingDigits % 3
	if firstGroupSize == 0 {
		firstGroupSize = 3
	}

	// Write first group
	for i := start; i < start+firstGroupSize; i++ {
		c.wrByte(dest, intPart[i])
	}

	// Write remaining groups with separators
	pos := start + firstGroupSize
	for pos < intLen {
		c.wrByte(dest, ',')
		for i := 0; i < 3 && pos < intLen; i++ {
			c.wrByte(dest, intPart[pos])
			pos++
		}
	}

	// Add decimal part if it exists
	if decPart != "" {
		c.wrString(dest, decPart)
	}
}

// removeTrailingZeros removes trailing zeros from decimal numbers in buffer
// Universal method with dest-first parameter order - follows buffer API architecture
func (c *conv) removeTrailingZeros(dest buffDest) {
	str := c.getString(dest)
	if len(str) == 0 {
		return
	}

	// Find decimal point
	dotIndex := -1
	for i := 0; i < len(str); i++ {
		if str[i] == '.' {
			dotIndex = i
			break
		}
	}

	if dotIndex == -1 {
		return // No decimal point
	}

	// Find last non-zero digit
	lastNonZero := len(str) - 1
	for i := len(str) - 1; i > dotIndex; i-- {
		if str[i] != '0' {
			lastNonZero = i
			break
		}
	}

	// Remove trailing zeros (and decimal point if all zeros)
	var result string
	if lastNonZero == dotIndex {
		result = str[:dotIndex] // Remove decimal point too
	} else {
		result = str[:lastNonZero+1]
	}

	c.rstBuffer(dest)
	c.wrString(dest, result)
}

// isNumericString checks if a string represents a valid number
// Universal helper method - follows buffer API architecture
func (c *conv) isNumericString(str string) bool {
	if len(str) == 0 {
		return false
	}

	i := 0
	// Handle sign
	if str[0] == '-' || str[0] == '+' {
		i = 1
		if i >= len(str) {
			return false // Just a sign is not a number
		}
	}

	hasDigit := false
	hasDecimal := false

	for ; i < len(str); i++ {
		if str[i] >= '0' && str[i] <= '9' {
			hasDigit = true
		} else if str[i] == '.' && !hasDecimal {
			hasDecimal = true
		} else {
			return false // Invalid character
		}
	}

	return hasDigit // Must have at least one digit
}
