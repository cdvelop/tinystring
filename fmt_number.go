package tinystring

// =============================================================================
// FORMAT NUMBER OPERATIONS - Number formatting with separators and display
// =============================================================================

// Thousands formats the number with thousand separators.
// By default (no param), uses EU style: 1.234.567,89
// If anglo is true, uses Anglo style: 1,234,567.89
func (t *conv) Thousands(anglo ...bool) *conv {
	if t.hasContent(buffErr) {
		return t
	}

	useAnglo := false
	if len(anglo) > 0 && anglo[0] {
		useAnglo = true
	}

	if t.hasContent(buffOut) {
		str := t.getString(buffOut)
		if t.isNumericString(str) {
			floatVal := t.parseFloat(str)
			if !t.hasContent(buffErr) {
				t.rstBuffer(buffOut)
				if floatVal == float64(int64(floatVal)) {
					t.wrIntBase(buffOut, int64(floatVal), 10, true)
				} else {
					t.wrFloat64(buffOut, floatVal)
					t.removeTrailingZeros(buffOut)
				}
			}
			t.addThousandSeparatorsCustom(buffOut, useAnglo)
		}
		return t
	}
	return t
}

// addThousandSeparatorsCustom adds thousand separators to the numeric string in buffer.
// If anglo is true: 1,234,567.89; if false: 1.234.567,89
func (c *conv) addThousandSeparatorsCustom(dest buffDest, anglo bool) {
	str := c.getString(dest)
	if len(str) <= 3 {
		return
	}

	// Find decimal point if it exists
	dotIndex := -1
	for i, char := range str {
		if char == '.' {
			dotIndex = i
			break
		}
	}

	intPart := str
	decPart := ""
	if dotIndex != -1 {
		intPart = str[:dotIndex]
		decPart = str[dotIndex+1:]
	}

	intLen := len(intPart)
	if intPart[0] == '-' {
		if intLen <= 4 {
			return
		}
	} else {
		if intLen <= 3 {
			return
		}
	}

	c.rstBuffer(dest)
	start := 0
	if intPart[0] == '-' {
		c.wrByte(dest, '-')
		start = 1
	}

	remainingDigits := intLen - start
	firstGroupSize := remainingDigits % 3
	if firstGroupSize == 0 {
		firstGroupSize = 3
	}

	for i := start; i < start+firstGroupSize; i++ {
		c.wrByte(dest, intPart[i])
	}

	sep := byte('.')
	if anglo {
		sep = ','
	}

	pos := start + firstGroupSize
	for pos < intLen {
		c.wrByte(dest, sep)
		for i := 0; i < 3 && pos < intLen; i++ {
			c.wrByte(dest, intPart[pos])
			pos++
		}
	}

	// Add decimal part if it exists
	if decPart != "" {
		if anglo {
			c.wrByte(dest, '.')
			c.wrString(dest, decPart)
		} else {
			c.wrByte(dest, ',')
			c.wrString(dest, decPart)
		}
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
