package tinystring

// Small number lookup table to avoid allocations for small integers
var smallInts = [...]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
}

// Phase 11: Extended fast parsing for common integers (0-99999)
// parseSmallInt optimizes parsing of small integers using direct byte access
// Returns the parsed integer and nil if successful, otherwise returns 0 and non-nil error
// Expanded from 999 to 99999 to handle more common integer patterns
func (c *conv) parseSmallInt(s string) (int, error) {
	if len(s) == 0 {
		return 0, c.wrErr(D.String, D.Empty)
	}

	var out int
	var negative bool

	// Check for negative sign
	i := 0
	if s[0] == '-' {
		negative = true
		i = 1
		if len(s) == 1 {
			return 0, c.wrErr(D.Format, D.Invalid)
		}
	} else if s[0] == '+' {
		i = 1
		if len(s) == 1 {
			return 0, c.wrErr(D.Format, D.Invalid)
		}
	}

	// Manual parsing with overflow check for up to 99999
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, c.wrErr(D.Format, D.Invalid)
		}
		digit := int(s[i] - '0')
		// Check for overflow before multiplication
		if out > (99999-digit)/10 {
			return 0, c.wrErr(D.Number, D.Overflow)
		}
		out = out*10 + digit
	}

	if negative {
		out = -out
	}

	return out, nil
}

// tryParseAs attempts to parse the current content as the specified type
// Universal validation method that follows buffer API architecture
func (c *conv) tryParseAs(parseType kind, base int) bool {
	// Save original state (inline saveState)
	oBuf := make([]byte, c.outLen)
	copy(oBuf, c.out[:c.outLen])
	oVT := c.kind

	// Try direct parsing based on type
	switch parseType {
	case KInt:
		inp := c.getString(buffOut)
		c.rstBuffer(buffErr)                    // Clear error before parsing
		if stringToInt(c, inp, base, buffOut) { // Write result directly to output buffer
			c.kind = KInt
		} else {
			// Set appropriate error for failed integer parsing
			if len(inp) == 0 {
				c.wrErr(D.String, D.Empty)
			} else if base != 10 && inp[0] == '-' {
				c.wrErr(D.Base, D.Decimal, D.Invalid)
			} else {
				c.wrErr(D.Format, D.Invalid)
			}
			// Restore original state if parsing failed
			c.rstBuffer(buffOut)     // Clear buffer using API
			c.wrBytes(buffOut, oBuf) // Write using API
			c.kind = oVT
			return false
		}
	case KUint:
		inp := c.getString(buffOut)
		c.rstBuffer(buffErr) // Clear error before parsing
		c.stringToUint(inp, base)
		if c.hasContent(buffErr) {
			// Restore original state if parsing failed
			c.rstBuffer(buffOut)     // Clear buffer using API
			c.wrBytes(buffOut, oBuf) // Write using API
			c.kind = oVT
			return false
		}
	case KFloat64:
		inp := c.getString(buffOut)
		c.rstBuffer(buffErr) // Clear error before parsing
		if _, success := stringToFloat(c, inp, buffOut); !success {
			// Restore original state if parsing failed
			c.rstBuffer(buffOut)     // Clear buffer using API
			c.wrBytes(buffOut, oBuf) // Write using API
			c.kind = oVT
			return false
		}
	default:
		// Unsupported parse type
		c.wrErr(D.Type, D.Not, D.Supported)
		return false
	}

	// Success - parsing completed successfully
	return true
}

// stringToUint parses a string as an unsigned integer and stores in output buffer
func (c *conv) stringToUint(input string, base int) {
	// Fast path for small integers in base 10
	if base == 10 {
		if val, err := c.parseSmallInt(input); err == nil && val >= 0 {
			uint64To(c, uint64(val), buffOut)
			c.anyValue = uint64(val)
			c.kind = KUint
			return
		}
	}

	// General parsing for larger numbers or other bases
	if len(input) == 0 {
		c.wrErr(D.String, D.Empty)
		return
	}

	var result uint64
	i := 0

	// Skip optional '+' sign (no negative for uint)
	if input[0] == '+' {
		i = 1
	} else if input[0] == '-' {
		c.wrErr(D.Number, D.Negative, D.Not, D.Allowed)
		return
	}

	if i >= len(input) {
		c.wrErr(D.Format, D.Invalid)
		return
	}

	// Parse digits
	for ; i < len(input); i++ {
		var digit uint64
		char := input[i]

		if char >= '0' && char <= '9' {
			digit = uint64(char - '0')
		} else if char >= 'a' && char <= 'z' {
			digit = uint64(char - 'a' + 10)
		} else if char >= 'A' && char <= 'Z' {
			digit = uint64(char - 'A' + 10)
		} else {
			c.wrErr(D.Character, D.Invalid)
			return
		}

		if digit >= uint64(base) {
			c.wrErr(D.Digit, D.Out, D.Of, D.Range, D.For, D.Base)
			return
		}

		// Check for overflow
		if result > (18446744073709551615-digit)/uint64(base) {
			c.wrErr(D.Number, D.Overflow)
			return
		}

		result = result*uint64(base) + digit
	}

	uint64To(c, result, buffOut)
	c.anyValue = result
	c.kind = KUint
}

// stringToInt parses a string as an integer and writes result to specified buffer destination
// Universal method with dest-first parameter order - follows buffer API architecture
func stringToInt(c *conv, inp string, base int, dest buffDest) bool {
	// Fast path for small integers in base 10
	if base == 10 {
		if val, err := c.parseSmallInt(inp); err == nil {
			c.wrInt(dest, int64(val))
			c.anyValue = int64(val)
			c.kind = KInt
			return true
		}
	}

	// General parsing for larger numbers or other bases
	var result int64
	var negative bool
	i := 0

	if len(inp) == 0 {
		c.wrErr(D.String, D.Empty)
		return false
	}

	// Handle sign
	if inp[0] == '-' {
		// Negative numbers are only valid in base 10
		if base != 10 {
			c.wrErr(D.Number, D.Negative, D.Not, D.Allowed, D.For, D.Base)
			return false
		}
		negative = true
		i = 1
	} else if inp[0] == '+' {
		i = 1
	}

	if i >= len(inp) {
		c.wrErr(D.Format, D.Invalid)
		return false
	}

	// Parse digits
	for ; i < len(inp); i++ {
		var digit int64
		char := inp[i]

		if char >= '0' && char <= '9' {
			digit = int64(char - '0')
		} else if char >= 'a' && char <= 'z' {
			digit = int64(char - 'a' + 10)
		} else if char >= 'A' && char <= 'Z' {
			digit = int64(char - 'A' + 10)
		} else {
			c.wrErr(D.Character, D.Invalid)
			return false
		}

		if digit >= int64(base) {
			c.wrErr(D.Digit, D.Out, D.Of, D.Range, D.For, D.Base)
			return false
		}

		// Check for overflow
		if result > (9223372036854775807-digit)/int64(base) {
			c.wrErr(D.Number, D.Overflow)
			return false
		}

		result = result*int64(base) + digit
	}

	if negative {
		result = -result
	}

	c.wrInt(dest, result)
	c.anyValue = result
	c.kind = KInt
	return true
}

// stringToFloat parses a string as a float and writes result to specified buffer destination
// Universal method with dest-first parameter order - follows buffer API architecture
func stringToFloat(c *conv, inp string, dest buffDest) (float64, bool) {
	val := c.parseFloat(inp)
	if c.hasContent(buffErr) {
		return 0, false
	}

	c.wrFloat(dest, val)
	c.anyValue = val
	c.kind = KFloat64
	return val, true
}

// uint64To converts uint64 to string and writes to specified buffer destination
// Universal method with dest-first parameter order - follows buffer API architecture
func uint64To(c *conv, val uint64, dest buffDest) {
	if val == 0 {
		c.wrString(dest, "0")
		return
	}

	// Use lookup table for small values
	if val < uint64(len(smallInts)) {
		c.wrString(dest, smallInts[val])
		return
	}

	// General case: convert using division
	var result [64]byte
	i := len(result)

	for val > 0 {
		i--
		result[i] = byte(val%10) + '0'
		val /= 10
	}

	c.wrBytes(dest, result[i:])
}
