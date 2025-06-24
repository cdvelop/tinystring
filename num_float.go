package tinystring

// =============================================================================
// FLOAT OPERATIONS - All float parsing, conversion and formatting
// =============================================================================

// ToFloat converts the value to a float64.
// Returns the converted float64 and any error that occurred during conversion.
func (c *conv) ToFloat() (float64, error) {
	// Check for existing error
	if c.hasContent(buffErr) {
		return 0, c
	}

	// Try parsing current content as float
	if c.tryParseAs(KFloat64, 10) {
		if val, ok := c.anyValue.(float64); ok {
			return val, nil
		}
	}

	// If parsing failed, return error
	if c.hasContent(buffErr) {
		return 0, c
	}

	return 0, c.wrErr(D.Format, D.Invalid)
}

// parseFloat parses a string as a float64 and returns the result
// Universal method that follows buffer API architecture
func (c *conv) parseFloat(inp string) float64 {
	if len(inp) == 0 {
		c.wrErr(D.String, D.Empty)
		return 0
	}

	var result float64
	var negative bool
	var hasDecimal bool
	var decimalPlaces int
	i := 0

	// Handle sign
	if inp[0] == '-' {
		negative = true
		i = 1
	} else if inp[0] == '+' {
		i = 1
	}

	if i >= len(inp) {
		c.wrErr(D.Format, D.Invalid)
		return 0
	}

	// Parse integer part
	for ; i < len(inp) && inp[i] != '.'; i++ {
		if inp[i] < '0' || inp[i] > '9' {
			c.wrErr(D.Character, D.Invalid)
			return 0
		}
		result = result*10 + float64(inp[i]-'0')
	}

	// Parse decimal part if present
	if i < len(inp) && inp[i] == '.' {
		hasDecimal = true
		i++ // Skip decimal point

		for ; i < len(inp); i++ {
			if inp[i] < '0' || inp[i] > '9' {
				c.wrErr(D.Character, D.Invalid)
				return 0
			}
			decimalPlaces++
			result = result*10 + float64(inp[i]-'0')
		}
	}

	// Apply decimal places
	if hasDecimal {
		for j := 0; j < decimalPlaces; j++ {
			result /= 10
		}
	}

	if negative {
		result = -result
	}

	return result
}

// wrFloat converts float64 to string and writes to specified buffer destination
// Universal method with dest-first parameter order - follows buffer API architecture
func (c *conv) wrFloat(dest buffDest, val float64) {
	c.kind = KFloat64 // Set type
	c.anyValue = val  // Store original value

	// Handle special cases
	if val != val { // NaN
		c.wrString(dest, "NaN")
		return
	}
	if val == 0 {
		c.wrString(dest, "0")
		return
	}

	// Handle infinity
	if val > 1.7976931348623157e+308 {
		c.wrString(dest, "+Inf")
		return
	}
	if val < -1.7976931348623157e+308 {
		c.wrString(dest, "-Inf")
		return
	}

	// Handle negative numbers
	negative := val < 0
	if negative {
		c.wrString(dest, "-")
		val = -val
	}

	// Check if it's effectively an integer
	if val < 1e15 && val == float64(int64(val)) {
		c.fmtIntToDest(dest, int64(val), 10, false)
		return
	}

	// For numbers with decimal places, use a precision-limited approach
	// Round to 6 decimal places to avoid precision issues
	scaled := val * 1000000
	rounded := int64(scaled + 0.5)

	intPart := rounded / 1000000
	fracPart := rounded % 1000000

	// Write integer part
	c.fmtIntToDest(dest, intPart, 10, false)

	// Write fractional part if non-zero
	if fracPart > 0 {
		c.wrString(dest, ".")

		// Build fractional string using local array to avoid buffer conflicts
		var digits [6]byte
		temp := fracPart
		for i := 0; i < 6; i++ {
			digits[i] = byte(temp%10) + '0'
			temp /= 10
		}

		// Find the start position (skip leading zeros in the array)
		// Since digits[0] corresponds to the rightmost digit, we need to skip zeros from the left
		start := 0
		for start < 6 && digits[start] == '0' {
			start++
		}

		// Write digits in reverse order (correct order), skipping leading zeros
		if start < 6 {
			for i := 5; i >= start; i-- {
				c.wrByte(dest, digits[i])
			}
		}
	}
}
