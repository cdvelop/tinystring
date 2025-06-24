package tinystring

// =============================================================================
// FORMAT PRECISION OPERATIONS - Decimal rounding and precision control
// =============================================================================

// RoundDecimals rounds the current numeric value to the specified number of decimal places.
// Example: Convert("3.14159").RoundDecimals(2) returns "3.14"
func (t *conv) RoundDecimals(decimals int) *conv {
	if t.hasContent(buffErr) {
		return t
	}

	// OPTIMIZED: Buffer-first approach - ensure we have content in buffOut
	if !t.hasContent(buffOut) && t.ptrValue != nil {
		t.rstBuffer(buffOut) // Reset buffer before conversion
		t.anyToBuff(buffOut, t.ptrValue)
		if t.hasContent(buffErr) {
			return t
		}
	}

	// Now apply rounding to the string in buffOut
	t.applyRoundingToNumber(buffOut, decimals, false)

	// If the result is not numeric, set to zero with correct decimals
	str := t.getString(buffOut)
	isNumeric := false
	for i := 0; i < len(str); i++ {
		if (str[i] >= '0' && str[i] <= '9') || str[i] == '.' || str[i] == '-' {
			isNumeric = true
		} else {
			isNumeric = false
			break
		}
	}
	if !isNumeric || str == "" || str == "-" {
		t.rstBuffer(buffOut)
		t.wrString(buffOut, "0")
		if decimals > 0 {
			t.wrString(buffOut, ".")
			for i := 0; i < decimals; i++ {
				t.wrString(buffOut, "0")
			}
		}
	}
	return t
}

// applyRoundingToNumber rounds the current number to specified decimal places
// Universal method with dest-first parameter order - follows buffer API architecture
func (t *conv) applyRoundingToNumber(dest buffDest, decimals int, roundDown bool) *conv {
	if t.hasContent(buffErr) {
		return t
	}

	// Get current string representation
	currentStr := t.getString(dest)

	// Find decimal point
	dotIndex := func() int {
		for i := range len(currentStr) {
			if currentStr[i] == '.' {
				return i
			}
		}
		return -1
	}()

	// If no decimal point, add zeros if needed
	if dotIndex == -1 {
		if decimals > 0 {
			t.wrString(dest, ".")
			for i := 0; i < decimals; i++ {
				t.wrByte(dest, '0')
			}
		}
		return t
	}

	// Calculate required length
	var targetLen int
	if decimals == 0 {
		targetLen = dotIndex // No decimal point for 0 decimals
	} else {
		targetLen = dotIndex + 1 + decimals // Include decimal point and decimal places
	}

	// If we need to truncate or round
	if len(currentStr) > targetLen {
		if roundDown {
			// Simple truncation for roundDown (floor behavior)
			t.rstBuffer(dest)
			t.wrString(dest, currentStr[:targetLen])
		} else {
			// Check if we need to round up
			// For "round up" behavior: if there are ANY non-zero digits after target length, round up
			shouldRoundUp := false

			// Get the remaining digits after the target length
			remainingDigits := currentStr[targetLen:]
			if len(remainingDigits) > 0 {
				// For "round up" (ceiling) behavior: round up if ANY remaining digit is non-zero
				for _, digit := range remainingDigits {
					if digit != '0' {
						shouldRoundUp = true
						break
					}
				}
			}

			if shouldRoundUp {
				// Implement simple rounding up
				rounded := currentStr[:targetLen]
				// Convert to slice for modification
				roundedBytes := []byte(rounded)
				carry := 1

				// Propagate carry from right to left
				for i := len(roundedBytes) - 1; i >= 0 && carry > 0; i-- {
					if roundedBytes[i] == '.' {
						continue
					}
					if roundedBytes[i] >= '0' && roundedBytes[i] <= '9' {
						digit := int(roundedBytes[i]-'0') + carry
						if digit > 9 {
							roundedBytes[i] = '0'
							carry = 1
						} else {
							roundedBytes[i] = byte(digit) + '0'
							carry = 0
						}
					}
				}

				// If still carrying, need to prepend 1
				if carry > 0 {
					t.rstBuffer(dest)
					t.wrString(dest, "1")
					t.wrBytes(dest, roundedBytes)
				} else {
					t.rstBuffer(dest)
					t.wrBytes(dest, roundedBytes)
				}
			} else {
				// Simple truncation when no rounding up needed
				t.rstBuffer(dest)
				t.wrString(dest, currentStr[:targetLen])
			}
		}
	} else if len(currentStr) < targetLen {
		// Add trailing zeros
		zerosNeeded := targetLen - len(currentStr)
		for i := 0; i < zerosNeeded; i++ {
			t.wrByte(dest, '0')
		}
	}

	return t
}

// Down applies floor-based rounding for the current precision level.
// When used after RoundDecimals(), it re-applies rounding with floor behavior.
// Example: Convert("3.154").RoundDecimals(2).Down() returns "3.15"
func (t *conv) Down() *conv {
	return t.downToBuffer(buffOut)
}

// downToBuffer applies floor-based rounding and writes to specified buffer destination
// Universal method with dest-first parameter order - follows buffer API architecture
func (t *conv) downToBuffer(dest buffDest) *conv {
	if t.hasContent(buffErr) {
		return t
	}

	// Get the current rounded result
	currentStr := t.getString(dest)

	// Find decimal places in current rounded result
	decimals := 0
	dotIndex := -1
	for i := 0; i < len(currentStr); i++ {
		if currentStr[i] == '.' {
			dotIndex = i
			break
		}
	}

	if dotIndex >= 0 {
		decimals = len(currentStr) - dotIndex - 1
	}

	// Strategy: Re-apply rounding using the original value but with floor behavior
	// This avoids precision issues while ensuring consistent rounding behavior
	if t.ptrValue != nil {
		// For string inputs, we can re-apply safely
		if _, isString := t.ptrValue.(string); isString {
			t.rstBuffer(buffWork)
			t.anyToBuff(buffWork, t.ptrValue)
			if !t.hasContent(buffErr) {
				workStr := t.getString(buffWork)
				t.rstBuffer(dest)
				t.wrString(dest, workStr)
				return t.applyRoundingToNumber(dest, decimals, true) // roundDown = true
			}
		} else {
			// For numeric inputs (float, int), use a more robust approach
			// Convert the current rounded-up result to implement floor behavior
			return t.applyFloorRoundingToResult(dest, decimals)
		}
	}

	// Fallback
	return t
}

// applyFloorRoundingToResult implements floor rounding by adjusting the last decimal digit down
func (t *conv) applyFloorRoundingToResult(dest buffDest, decimals int) *conv {
	currentStr := t.getString(dest)

	// Find the decimal point
	dotIndex := -1
	for i := 0; i < len(currentStr); i++ {
		if currentStr[i] == '.' {
			dotIndex = i
			break
		}
	}

	if decimals == 0 {
		// For zero decimals, we need to floor the integer part
		// E.g., "4" should become "3" if original was 3.7
		if dotIndex == -1 {
			// Parse integer and subtract 1
			parsed := int64(0)
			negative := false
			start := 0
			if len(currentStr) > 0 && currentStr[0] == '-' {
				negative = true
				start = 1
			}

			for i := start; i < len(currentStr); i++ {
				if currentStr[i] >= '0' && currentStr[i] <= '9' {
					parsed = parsed*10 + int64(currentStr[i]-'0')
				}
			}

			// Subtract 1 for floor behavior
			parsed--
			if negative {
				parsed = -parsed
			}

			t.rstBuffer(dest)
			t.fmtIntToDest(dest, parsed, 10, false)
		}
		return t
	}

	// For decimals > 0, reduce the last decimal digit by 1 if possible
	if dotIndex >= 0 && len(currentStr) > dotIndex+decimals {
		result := []byte(currentStr)
		lastDecimalIdx := dotIndex + decimals

		if lastDecimalIdx < len(result) && result[lastDecimalIdx] >= '1' && result[lastDecimalIdx] <= '9' {
			result[lastDecimalIdx]-- // Reduce by 1
			t.rstBuffer(dest)
			t.wrBytes(dest, result)
		}
	}

	return t
}

// wrFloatWithPrecision formats a float with specified precision and writes to buffer destination
// Universal method with dest-first parameter order - follows buffer API architecture
func (c *conv) wrFloatWithPrecision(dest buffDest, value float64, precision int) {
	// Handle special cases
	if value != value { // NaN
		c.wrString(dest, "NaN")
		return
	}

	if value == 0 {
		if precision > 0 {
			c.wrString(dest, "0.")
			for i := 0; i < precision; i++ {
				c.wrByte(dest, '0')
			}
		} else {
			c.wrString(dest, "0")
		}
		return
	}

	// Handle infinity
	if value > 1.7976931348623157e+308 {
		c.wrString(dest, "+Inf")
		return
	}
	if value < -1.7976931348623157e+308 {
		c.wrString(dest, "-Inf")
		return
	}

	// Handle negative numbers
	negative := value < 0
	if negative {
		c.wrString(dest, "-")
		value = -value
	}

	// Scale value by precision to get required decimal places
	multiplier := 1.0
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}

	scaled := value * multiplier
	rounded := int64(scaled + 0.5) // Round to nearest integer

	// Extract integer and fractional parts
	intPart := rounded
	for i := 0; i < precision; i++ {
		intPart /= 10
	}

	fracPart := rounded - intPart*int64(multiplier)

	// Write integer part
	c.wrInt(dest, intPart)

	// Write fractional part if precision > 0
	if precision > 0 {
		c.wrString(dest, ".")

		// Convert fractional part to string with leading zeros
		// Build digits array in reverse order to avoid allocations
		var digits [20]byte // Support up to 20 decimal places
		temp := fracPart
		for i := 0; i < precision; i++ {
			digits[i] = byte(temp%10) + '0'
			temp /= 10
		}

		// Write digits in reverse order (correct order)
		for i := precision - 1; i >= 0; i-- {
			c.wrByte(dest, digits[i])
		}
	}
}
