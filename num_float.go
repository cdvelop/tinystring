package tinystring

// =============================================================================
// FLOAT OPERATIONS - All float parsing, conversion and formatting
// =============================================================================

// ToFloat64 converts the value to a float64.
// Returns the converted float64 and any error that occurred during conversion.
func (c *conv) ToFloat64() (float64, error) {
	val := c.parseFloatBase()
	if c.hasContent(buffErr) {
		return 0, c
	}
	return val, nil
}

// ToFloat32 converts the value to a float32.
// Returns the converted float32 and any error that occurred during conversion.
func (c *conv) ToFloat32() (float32, error) {
	val := c.parseFloatBase()
	if c.hasContent(buffErr) {
		return 0, c
	}
	if val > 3.4028235e+38 {
		return 0, c.wrErr(D.Number, D.Overflow)
	}
	return float32(val), nil
}

// parseFloatBase parses the buffer as a float64, similar to parseIntBase for ints.
// It always uses the buffer output and handles errors internally.
func (c *conv) parseFloatBase() float64 {
	c.rstBuffer(buffErr)

	s := c.getString(buffOut)
	if len(s) == 0 {
		c.wrErr(D.String, D.Empty)
		return 0
	}

	var result float64
	var negative bool
	var hasDecimal bool
	var decimalPlaces int
	i := 0

	// Handle sign
	if s[0] == '-' {
		negative = true
		i = 1
		if len(s) == 1 {
			c.wrErr(D.Format, D.Invalid)
			return 0
		}
	} else if s[0] == '+' {
		i = 1
		if len(s) == 1 {
			c.wrErr(D.Format, D.Invalid)
			return 0
		}
	}

	// Parse integer part
	for ; i < len(s) && s[i] != '.'; i++ {
		if s[i] < '0' || s[i] > '9' {
			c.wrErr(D.Character, D.Invalid)
			return 0
		}
		result = result*10 + float64(s[i]-'0')
	}

	// Parse decimal part if present
	if i < len(s) && s[i] == '.' {
		hasDecimal = true
		i++ // Skip decimal point
		for ; i < len(s); i++ {
			if s[i] < '0' || s[i] > '9' {
				c.wrErr(D.Character, D.Invalid)
				return 0
			}
			decimalPlaces++
			result = result*10 + float64(s[i]-'0')
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

// DEPRECATED
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

// wrFloat32 writes a float32 to the buffer destination.
func (c *conv) wrFloat32(dest buffDest, val float32) {
	c.wrFloatBase(dest, float64(val), 3.4028235e+38)
}

// wrFloat64 writes a float64 to the buffer destination.
func (c *conv) wrFloat64(dest buffDest, val float64) {
	c.wrFloatBase(dest, float64(val), 1.7976931348623157e+308)
}

// wrFloatBase contains the shared logic for writing float values.
func (c *conv) wrFloatBase(dest buffDest, val float64, maxInf float64) {
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
	if val > maxInf {
		c.wrString(dest, "+Inf")
		return
	}
	if val < -maxInf {
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
		c.wrIntBase(dest, int64(val), 10, false)
		return
	}

	// For numbers with decimal places, use a precision-limited approach
	// Round to 6 decimal places to avoid precision issues
	scaled := val * 1000000
	rounded := int64(scaled + 0.5)

	intPart := rounded / 1000000
	fracPart := rounded % 1000000

	// Write integer part
	c.wrIntBase(dest, intPart, 10, false)

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
