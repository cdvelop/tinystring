package tinystring

// =============================================================================
// INTEGER OPERATIONS - All integer parsing, conversion and formatting
// =============================================================================

// ToInt converts the value to an integer with optional base specification.
// If no base is provided, base 10 is used. Supports bases 2-36.
// Returns the converted integer and any error that occurred during conversion.
func (c *conv) ToInt(base ...int) (int, error) {
	// Check for existing error
	if c.hasContent(buffErr) {
		return 0, c
	}

	// Validate base parameter
	baseVal := 10 // default base
	if len(base) > 0 {
		baseVal = base[0]
		if baseVal < 2 || baseVal > 36 {
			return 0, c.wrErr(D.Base, D.Invalid)
		}
	}

	// If base is not 10, we MUST parse the string representation
	// because the numeric value might have been parsed with base 10 initially
	if baseVal != 10 {
		// We need to parse the original string representation with the specified base
		inp := c.ensureStringInOut()
		c.rstBuffer(buffErr) // Clear any previous errors
		if stringToInt(c, inp, baseVal, buffWork) {
			if val, ok := c.anyValue.(int64); ok {
				// Check if int64 fits in int
				if val < -2147483648 || val > 2147483647 {
					return 0, c.wrErr(D.Number, D.Overflow)
				}
				return int(val), nil
			}
		}

		// If parsing failed, ensure we have an appropriate error
		if !c.hasContent(buffErr) {
			// Check for specific cases that should error
			if len(inp) > 0 && inp[0] == '-' {
				// Negative numbers are not valid in non-decimal bases
				c.wrErr(D.Base, D.Decimal, D.Invalid)
			} else {
				c.wrErr(D.Format, D.Invalid)
			}
		}
		return 0, c
	}

	// For base 10, direct conversion for numeric types already in memory is OK
	if c.anyValue != nil {
		switch v := c.anyValue.(type) {
		case int:
			return v, nil
		case int8:
			return int(v), nil
		case int16:
			return int(v), nil
		case int32:
			return int(v), nil
		case int64:
			// Check if int64 fits in int
			if v < -2147483648 || v > 2147483647 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return int(v), nil
		case uint:
			if v > 2147483647 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return int(v), nil
		case uint8:
			return int(v), nil
		case uint16:
			return int(v), nil
		case uint32:
			if v > 2147483647 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return int(v), nil
		case uint64:
			if v > 2147483647 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return int(v), nil
		case float32:
			return int(v), nil // Truncate decimal part
		case float64:
			return int(v), nil // Truncate decimal part
		}
	}

	// For string inputs, try parsing as float first if it contains decimal point
	if c.kind == KString {
		inp := c.ensureStringInOut()

		// If it contains a decimal point, try parsing as float first
		if Contains(inp, ".") {
			// Try to parse as float and truncate
			floatVal := c.parseFloat(inp)
			if !c.hasContent(buffErr) {
				return int(floatVal), nil
			}
			// Clear any parsing errors and continue with integer parsing
			c.rstBuffer(buffErr)
		}
	}

	// Try parsing current content as integer with base 10
	if c.tryParseAs(KInt, 10) {
		if val, ok := c.anyValue.(int64); ok {
			// Check if int64 fits in int
			if val < -2147483648 || val > 2147483647 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return int(val), nil
		}
	}

	// If parsing failed, return error
	if c.hasContent(buffErr) {
		return 0, c
	}

	return 0, c.wrErr(D.Format, D.Invalid)
}

// ToInt64 converts the value to a 64-bit integer with optional base specification.
// If no base is provided, base 10 is used. Supports bases 2-36.
// Returns the converted int64 and any error that occurred during conversion.
func (c *conv) ToInt64(base ...int) (int64, error) {
	// Validate base parameter
	baseVal := 10 // default base
	if len(base) > 0 {
		baseVal = base[0]
		if baseVal < 2 || baseVal > 36 {
			return 0, c.wrErr(D.Base, D.Invalid)
		}
	}

	// Check for existing error
	if c.hasContent(buffErr) {
		return 0, c.wrErr(c.getString(buffErr))
	}

	// Try parsing current content as integer
	if c.tryParseAs(KInt, baseVal) {
		if val, ok := c.anyValue.(int64); ok {
			return val, nil
		}
	}

	// If parsing failed, return error
	if c.hasContent(buffErr) {
		return 0, c.wrErr(c.getString(buffErr))
	}

	return 0, c.wrErr(D.Format, D.Invalid)
}

// ToUint converts the value to an unsigned integer with optional base specification.
// If no base is provided, base 10 is used. Supports bases 2-36.
// Returns the converted uint and any error that occurred during conversion.
func (c *conv) ToUint(base ...int) (uint, error) {
	// Check for existing error
	if c.hasContent(buffErr) {
		return 0, c
	}

	// Direct conversion for numeric types already in memory
	if c.anyValue != nil {
		switch v := c.anyValue.(type) {
		case int:
			if v < 0 {
				return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			}
			return uint(v), nil
		case int8:
			if v < 0 {
				return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			}
			return uint(v), nil
		case int16:
			if v < 0 {
				return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			}
			return uint(v), nil
		case int32:
			if v < 0 {
				return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			}
			return uint(v), nil
		case int64:
			if v < 0 {
				return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			}
			if v > 4294967295 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return uint(v), nil
		case uint:
			return v, nil
		case uint8:
			return uint(v), nil
		case uint16:
			return uint(v), nil
		case uint32:
			return uint(v), nil
		case uint64:
			if v > 4294967295 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return uint(v), nil
		case float32:
			if v < 0 {
				return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			}
			return uint(v), nil // Truncate decimal part
		case float64:
			if v < 0 {
				return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
			}
			return uint(v), nil // Truncate decimal part
		}
	}

	// Validate base parameter for string parsing
	baseVal := 10 // default base
	if len(base) > 0 {
		baseVal = base[0]
		if baseVal < 2 || baseVal > 36 {
			return 0, c.wrErr(D.Base, D.Invalid)
		}
	}

	// For string inputs, try parsing as float first if it contains decimal point
	if c.kind == KString {
		inp := c.ensureStringInOut()

		// If it contains a decimal point, try parsing as float first
		if Contains(inp, ".") {
			// Try to parse as float and truncate
			floatVal := c.parseFloat(inp)
			if !c.hasContent(buffErr) {
				if floatVal < 0 {
					return 0, c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
				}
				return uint(floatVal), nil
			}
			// Clear any parsing errors and continue with integer parsing
			c.rstBuffer(buffErr)
		}
	}

	// Try parsing current content as unsigned integer
	if c.tryParseAs(KUint, baseVal) {
		if val, ok := c.anyValue.(uint64); ok {
			// Check if uint64 fits in uint
			if val > 4294967295 {
				return 0, c.wrErr(D.Number, D.Overflow)
			}
			return uint(val), nil
		}
	}

	// If parsing failed, return error
	if c.hasContent(buffErr) {
		return 0, c
	}

	return 0, c.wrErr(D.Format, D.Invalid)
}

// wrInt converts int64 to string and writes to specified buffer destination
// Universal method with dest-first parameter order - eliminates duplicate code
func (c *conv) wrInt(dest buffDest, v int64) {
	c.kind = KInt  // Set type
	c.anyValue = v // Store original value

	// Use existing fmtIntToDest for conversion
	c.fmtIntToDest(dest, v, 10, true)
}

// wrUint converts uint64 to string and writes to specified buffer destination
// Universal method with dest-first parameter order - eliminates duplicate code
func (c *conv) wrUint(dest buffDest, v uint64) {
	c.kind = KUint // Set type
	c.anyValue = v // Store original value

	// Use existing fmtIntToDest for conversion
	c.fmtIntToDest(dest, int64(v), 10, false)
}

// wrInt64Base converts an int64 to a string with specified base and writes to destination buffer
// Universal method that receives parameters instead of using temp fields
func (c *conv) wrInt64Base(dest buffDest, number int64, base int) {
	if number == 0 {
		c.wrString(dest, "0")
		return
	}

	// Use optimized wrInt() for decimal base
	if base == 10 {
		c.wrInt(dest, number)
		return
	}

	isNegative := number < 0
	if isNegative {
		number = -number
	}

	// Inline validateBase logic
	if base < 2 || base > 36 {
		c.wrString(buffErr, T(D.Base, " ", D.Invalid))
		return
	}

	// Convert to string with base
	digits := "0123456789abcdef"
	var out [64]byte // Maximum digits for int64 in base 2
	idx := len(out)

	// Build string backwards
	for number > 0 {
		idx--
		out[idx] = digits[number%int64(base)]
		number /= int64(base)
	}

	if isNegative {
		idx--
		out[idx] = '-'
	}

	c.wrBytes(dest, out[idx:])
}

// fmtIntToOut converts integer to string and writes to out buffer
// Replaces fmtIntGeneric() with centralized buffer management
func (c *conv) fmtIntToOut(val int64, base int, signed bool) {
	if val == 0 {
		c.wrString(buffOut, "0")
		return
	}
	// Handle negative numbers for signed integers
	if signed && val < 0 {
		val = -val
		c.wrString(buffOut, "-")
	}

	// Convert using existing manual implementation logic
	// Use work buffer for intermediate operations
	c.rstBuffer(buffWork)

	// Build digits in reverse order in work buffer
	for val > 0 {
		digit := byte(val%int64(base)) + '0'
		if digit > '9' {
			digit += 'a' - '9' - 1
		}
		c.work = append(c.work, digit)
		c.workLen++
		val /= int64(base)
	}

	// Reverse and write to out buffer
	for i := c.workLen - 1; i >= 0; i-- {
		c.wrByte(buffOut, c.work[i])
	}
}

// fmtIntToDest converts integer to string and writes to specified buffer destination
// Universal method with dest-first parameter order - eliminates duplicate code
func (c *conv) fmtIntToDest(dest buffDest, val int64, base int, signed bool) {
	if val == 0 {
		c.wrString(dest, "0")
		return
	}
	// Handle negative numbers for signed integers
	if signed && val < 0 {
		val = -val
		c.wrString(dest, "-")
	}

	// Convert using existing manual implementation logic
	// Use a different approach to avoid conflicting with dest buffer

	// Build digits directly without intermediate buffer when dest is buffWork
	var digits []byte
	tempVal := val
	for tempVal > 0 {
		digit := byte(tempVal%int64(base)) + '0'
		if digit > '9' {
			digit += 'a' - '9' - 1
		}
		digits = append(digits, digit)
		tempVal /= int64(base)
	}

	// Write digits in reverse order directly to destination
	for i := len(digits) - 1; i >= 0; i-- {
		c.wrByte(dest, digits[i])
	}
}
