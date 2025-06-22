package tinystring

// Phase 13: Static error messages to eliminate T() allocations
var (
	errBaseInvalid       = "base invalid"
	errNumberOverflow    = "number overflow"
	errInvalidNonNumeric = "invalid non-numeric character"
	errCharacterInvalid  = "character invalid"
	errDigitOutOfRange   = "digit out of range for base"
)

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
func parseSmallInt(s string) (int, error) {
	if len(s) == 0 {
		return 0, Err(D.String, D.Empty)
	}

	var out int
	var negative bool

	// Check for negative sign
	i := 0
	if s[0] == '-' {
		negative = true
		i = 1
		if len(s) == 1 {
			return 0, Err(D.Format, D.Invalid)
		}
	} else if s[0] == '+' {
		i = 1
		if len(s) == 1 {
			return 0, Err(D.Format, D.Invalid)
		}
	}

	// Parse digits
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, Err(D.Format, D.Invalid)
		}

		digit := int(s[i] - '0')

		// Phase 11: Extended overflow check for 99999 limit
		if out > (99999-digit)/10 {
			return 0, Err(D.Number, D.Overflow)
		}

		out = out*10 + digit
	}

	// Apply negative sign
	if negative {
		out = -out
	}

	return out, nil
}

// Shared helper methods to reduce code duplication between numeric.go and format.go

// tryParseAs attempts to parse content as the specified numeric type with fallback to float
func (c *conv) tryParseAs(parseType kind, base int) bool {
	// Save original state (inline saveState)
	oBuf := make([]byte, c.outLen)
	copy(oBuf, c.out[:c.outLen])
	oVT := c.kind
	
	// Try direct parsing based on type
	switch parseType {
	case KInt:
		inp := c.ensureStringInOut()
		c.clearError() // Clear error before parsing
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
		}
	case KUint:
		// For stringToUint, we need the string content
		c.stringToUint(string(c.out[:c.outLen]), base)
	}

	if !c.hasError() {
		return true
	}

	// Check if the error is due to invalid base with negative numbers
	// If so, don't attempt float fallback as it would bypass base validation
	if base != 10 && c.outLen > 0 && c.out[0] == '-' {
		return false
	}

	// If that fails, restore state and try to parse as float then convert
	// ✅ Inline restoreState logic using API
	c.rstOut()      // Clear buffer using API
	c.wrToOut(oBuf) // Write using API
	c.kind = oVT
	c.clearError() // Reset error when restoring state using API

	// Use new independent stringToFloat function
	inp := c.ensureStringInOut()
	if floatVal, ok := stringToFloat(c, inp, buffErr); ok {
		switch parseType {
		case KInt:
			// Store result using anyToBuff (no temp fields needed)
			c.kind = KInt
			anyToBuff(c, buffOut, int64(floatVal))
		case KUint:
			if floatVal < 0 {
				c.wrErr(D.Number, D.Negative, D.Not, D.Supported)
				return false
			}
			// Store result using anyToBuff (no temp fields needed)
			c.kind = KUint
			anyToBuff(c, buffOut, uint64(floatVal))
		}
		return true
	}

	return false
}

// ToInt converts the conv content to an integer with optional base specification.
//
// Parameters:
//   - base (optional): The numeric base for conversion (2-36). Default is 10 (decimal).
//     Common bases: 2 (binary), 8 (octal), 10 (decimal), 16 (hexadecimal)
//
// Returns:
//   - int: The converted integer value
//   - error: Any error that occurred during conversion
//
// Conversion behavior:
//  1. First attempts direct integer parsing with the specified base
//  2. If that fails, tries to parse as float and truncates to integer
//  3. Returns error if both methods fail
//
// Supported input formats:
//   - Integer strings: "123", "-456"
//   - Float strings (truncated): "123.45" -> 123, "99.99" -> 99
//   - Different bases: "1010" (base 2) -> 10, "FF" (base 16) -> 255
//   - Negative numbers: Only supported for base 10
//
// Usage examples:
//
//	// Basic decimal conversion
//	val, err := Convert("123").ToInt()
//	// val: 123, err: nil
//
//	// Binary conversion
//	val, err := Convert("1010").ToInt(2)
//	// val: 10, err: nil
//
//	// Hexadecimal conversion
//	val, err := Convert("FF").ToInt(16)
//	// val: 255, err: nil
//
//	// Float truncation
//	val, err := Convert("123.99").ToInt()
//	// val: 123, err: nil
//
//	// Error handling
//	val, err := Convert("invalid").ToInt()
//	// val: 0, err: conversion error
//
// Note: Negative numbers are only supported for base 10. For other bases,
// negative signs will out in an error.
func (c *conv) ToInt(base ...int) (int, error) {
	if c.hasError() { // Use buffer API
		return 0, c
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}

	// For typed values, convert using anyToBuff to work buffer then parse back
	if c.kind != KString {
		// If we already have an integer type, use the stored value directly
		if c.kind == KInt {
			if intVal, ok := c.pointerVal.(int64); ok {
				return int(intVal), nil
			}
		} else if c.kind == KUint {
			if uintVal, ok := c.pointerVal.(uint64); ok {
				return int(uintVal), nil
			}
		} else if c.kind == KFloat64 {
			if floatVal, ok := c.pointerVal.(float64); ok {
				return int(floatVal), nil // Truncate float to int
			}
		}
		// Otherwise, convert current value to string using anyToBuff
		c.rstWork()
		anyToBuff(c, buffWork, c.pointerVal)
		str := c.getWorkString()

		// Parse the string representation
		if intVal, err := parseSmallInt(str); err == nil {
			return intVal, nil
		}
		// If parseSmallInt failed, continue to tryParseAs for fallback
	}

	// For string types and other types, use shared helper method for parsing with fallback
	if c.tryParseAs(KInt, b) {
		// tryParseAs succeeded and stringToInt stored result in pointerVal
		if intVal, ok := c.pointerVal.(int64); ok {
			return int(intVal), nil
		}
	}
	// Return error if parsing failed
	return 0, c
}

// ToInt64 converts the conv content to a 64-bit integer with optional base specification.
//
// Parameters:
//   - base (optional): The numeric base for conversion (2-36). Default is 10 (decimal).
//     Common bases: 2 (binary), 8 (octal), 10 (decimal), 16 (hexadecimal)
//
// Returns:
//   - int64: The converted 64-bit integer value
//   - error: Any error that occurred during conversion
//
// Conversion behavior:
//  1. First attempts direct 64-bit integer parsing with the specified base
//  2. If that fails, tries to parse as float and truncates to 64-bit integer
//  3. Returns error if both methods fail
//
// Supported input formats:
//   - Integer strings: "123", "-456", "9223372036854775807" (max int64)
//   - Float strings (truncated): "123.45" -> 123, "99.99" -> 99
//   - Different bases: "1010" (base 2) -> 10, "FF" (base 16) -> 255
//   - Negative numbers: Only supported for base 10
//   - Large numbers: Supports full int64 range (-9223372036854775808 to 9223372036854775807)
//
// Usage examples:
//
//	// Basic decimal conversion
//	val, err := Convert("123").ToInt64()
//	// val: 123, err: nil
//
//	// Large number conversion
//	val, err := Convert("9223372036854775807").ToInt64()
//	// val: 9223372036854775807 (max int64), err: nil
//
//	// Binary conversion
//	val, err := Convert("1010").ToInt64(2)
//	// val: 10, err: nil
//
//	// Hexadecimal conversion
//	val, err := Convert("7FFFFFFFFFFFFFFF").ToInt64(16)
//	// val: 9223372036854775807, err: nil
//
//	// Float truncation
//	val, err := Convert("123.99").ToInt64()
//	// val: 123, err: nil
//
//	// Error handling
//	val, err := Convert("invalid").ToInt64()
//	// val: 0, err: conversion error
//
// Note: This method provides the full range of 64-bit integers, which is useful
// for large numeric values that exceed the range of regular int type.
func (c *conv) ToInt64(base ...int) (int64, error) {
	if c.hasError() { // Use buffer API
		return 0, c
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}

	// Use shared helper method for parsing with fallback
	if c.tryParseAs(KInt, b) {
		// Parse the resulting buffer content as int64
		str := c.ensureStringInOut()
		if intVal, err := parseSmallInt(str); err == nil {
			return int64(intVal), nil
		}
	}

	// Return error if parsing failed
	return 0, c
}

// ToUint converts the conv content to an unsigned integer with optional base specification.
//
// Parameters:
//   - base (optional): The numeric base for conversion (2-36). Default is 10 (decimal).
//     Common bases: 2 (binary), 8 (octal), 10 (decimal), 16 (hexadecimal)
//
// Returns:
//   - uint: The converted unsigned integer value
//   - error: Any error that occurred during conversion
//
// Conversion behavior:
//  1. First attempts direct unsigned integer parsing with the specified base
//  2. If that fails, tries to parse as float and truncates to unsigned integer
//  3. Returns error if both methods fail or if the value is negative
//
// Supported input formats:
//   - Positive integer strings: "123", "456"
//   - Float strings (truncated): "123.45" -> 123, "99.99" -> 99
//   - Different bases: "1010" (base 2) -> 10, "FF" (base 16) -> 255
//   - Negative numbers: NOT supported, will return error
//
// Usage examples:
//
//	// Basic decimal conversion
//	val, err := Convert("123").ToUint()
//	// val: 123, err: nil
//
//	// Binary conversion
//	val, err := Convert("1010").ToUint(2)
//	// val: 10, err: nil
//
//	// Hexadecimal conversion
//	val, err := Convert("FF").ToUint(16)
//	// val: 255, err: nil
//
//	// Float truncation
//	val, err := Convert("123.99").ToUint()
//	// val: 123, err: nil
//
//	// Error with negative number
//	val, err := Convert("-123").ToUint()
//	// val: 0, err: "negative numbers are not supported for unsigned integers"
//
// Note: Negative numbers are never supported for unsigned integers and will
// always out in an error, regardless of the base.
func (c *conv) ToUint(base ...int) (uint, error) {
	if c.hasError() { // Use buffer API
		return 0, c
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}

	// For typed values, convert using anyToBuff and check for negatives
	if c.kind != KString {
		// Convert current value to string first
		c.rstWork()
		anyToBuff(c, buffWork, c.pointerVal)
		str := c.getWorkString()

		// Check for negative values
		if len(str) > 0 && str[0] == '-' {
			c.wrErr(D.Number, D.Negative, D.Not, D.Supported)
			return 0, c
		}

		// Parse as unsigned
		if intVal, err := parseSmallInt(str); err == nil && intVal >= 0 {
			return uint(intVal), nil
		}
	}

	// For string types and other types, use shared helper method for parsing with fallback
	if c.tryParseAs(KUint, b) {
		// Parse the resulting buffer content
		str := c.ensureStringInOut()
		if intVal, err := parseSmallInt(str); err == nil && intVal >= 0 {
			return uint(intVal), nil
		}
	}
	// Return error if parsing failed
	return 0, c
}

// ToFloat converts the conv content to a float64 (double precision floating point).
//
// Returns:
//   - float64: The converted floating point value
//   - error: Any error that occurred during conversion
//
// Conversion behavior:
//   - Parses the string content as a floating point number
//   - Supports both positive and negative numbers
//   - Handles decimal points and scientific notation (if implemented)
//   - Returns error for invalid number formats
//
// Supported input formats:
//   - Integer strings: "123" -> 123.0, "-456" -> -456.0
//   - Decimal numbers: "123.45", "-99.99", "0.001"
//   - Numbers with leading signs: "+123.45", "-0.99"
//   - Zero values: "0", "0.0", "0.000"
//
// Usage examples:
//
//	// Basic decimal conversion
//	val, err := Convert("123.45").ToFloat()
//	// val: 123.45, err: nil
//
//	// Integer to float
//	val, err := Convert("42").ToFloat()
//	// val: 42.0, err: nil
//
//	// Negative numbers
//	val, err := Convert("-99.99").ToFloat()
//	// val: -99.99, err: nil
//
//	// Error handling
//	val, err := Convert("invalid").ToFloat()
//	// val: 0.0, err: conversion error
//
// Note: This method uses a custom float parsing implementation that may have
// different precision characteristics compared to the standard library.
func (c *conv) ToFloat() (float64, error) {
	if c.hasError() { // Use buffer API
		return 0, c
	}

	// For typed values, convert using anyToBuff to work buffer then parse
	if c.kind != KString {
		c.rstWork()
		anyToBuff(c, buffWork, c.pointerVal)
		str := c.getWorkString()

		// Use stringToFloat to parse
		if floatVal, ok := stringToFloat(c, str, buffErr); ok {
			return floatVal, nil
		}
		return 0, c
	}

	// For string types, parse as float directly
	str := c.ensureStringInOut()
	if floatVal, ok := stringToFloat(c, str, buffErr); ok {
		return floatVal, nil
	}

	return 0, c
}

// stringToInt converts string to signed integer with specified base and stores result in pointerVal
// Independent function that receives parameters instead of using temp fields
// This unified function handles both int and int64 conversions using buffer API only
// Returns success status - does not set errors (caller handles error reporting)
func stringToInt(c *conv, inp string, base int, dest buffDest) bool {
	if len(inp) == 0 {
		return false
	}

	isNeg := false
	if inp[0] == '-' {
		if base != 10 {
			return false
		}
		isNeg = true
		inp = inp[1:] // Process without negative sign
	}

	// Try parsing with parseSmallInt first
	if intVal, err := parseSmallInt(inp); err == nil {
		if isNeg {
			intVal = -intVal
		}
		// Store the result in pointerVal and update buffer state based on destination
		c.pointerVal = int64(intVal)
		switch dest {
		case buffOut:
			// Write the integer value as string to output buffer for string representation
			anyToBuff(c, buffOut, int64(intVal))
		case buffWork:
			// Write the integer value as string to work buffer
			anyToBuff(c, buffWork, int64(intVal))
		case buffErr:
			// This case shouldn't happen for successful parsing
			anyToBuff(c, buffErr, int64(intVal))
		}
		return true
	}

	// If parseSmallInt fails, return false (no error set)
	return false
}

// stringToUint converts string to uint with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) stringToUint(input string, base int) {
	if len(input) == 0 {
		return
	}
	if input[0] == '-' {
		t.wrErr(D.Number, D.Negative, D.Not, D.Supported)
		return
	}
	// Update the conv struct with the provided input string
	t.setString(input)
	t.stringToIn(base)
	// Result already stored in t.uintVal by stringToIn
}

// stringToFloat converts string to float64 and writes result to specified buffer destination
// Independent function that uses buffer API and destination selection
func stringToFloat(c *conv, inp string, dest buffDest) (float64, bool) {
	if len(inp) == 0 {
		return 0, false
	}

	isNeg := false
	sIdx := 0
	if inp[0] == '-' {
		isNeg = true
		sIdx = 1
		if len(inp) == 1 { // Just a "-" sign
			anyToBuff(c, dest, D.Float)
			anyToBuff(c, dest, " ")
			anyToBuff(c, dest, D.String)
			anyToBuff(c, dest, " ")
			anyToBuff(c, dest, D.Invalid)
			return 0, false
		}
	} else if inp[0] == '+' {
		sIdx = 1
		if len(inp) == 1 { // Just a "+" sign
			anyToBuff(c, dest, D.Float)
			anyToBuff(c, dest, " ")
			anyToBuff(c, dest, D.String)
			anyToBuff(c, dest, " ")
			anyToBuff(c, dest, D.Invalid)
			return 0, false
		}
	}

	var ip uint64        // integerPart
	var fp uint64        // fractionPart
	var fd float64 = 1.0 // fractionDivisor
	dps := false         // decimalPointSeen
	pf := false          // parsingFraction
	hd := false          // hasDigits

	for i := sIdx; i < len(inp); i++ {
		ch := inp[i] // char
		if ch == '.' {
			if dps {
				anyToBuff(c, dest, D.Float)
				anyToBuff(c, dest, " ")
				anyToBuff(c, dest, D.String)
				anyToBuff(c, dest, " ")
				anyToBuff(c, dest, D.Invalid)
				return 0, false
			}
			dps = true
			pf = true
		} else if ch >= '0' && ch <= '9' {
			hd = true
			dgt := uint64(ch - '0') // digit
			if pf {
				fp = fp*10 + dgt
				fd *= 10.0
			} else { // Check for overflow in integer part
				if ip > ^uint64(0)/10 || (ip == ^uint64(0)/10 && dgt > ^uint64(0)%10) {
					anyToBuff(c, dest, D.Number)
					anyToBuff(c, dest, " ")
					anyToBuff(c, dest, D.Overflow)
					return 0, false
				}
				ip = ip*10 + dgt
			}
		} else {
			anyToBuff(c, dest, D.Float)
			anyToBuff(c, dest, " ")
			anyToBuff(c, dest, D.String)
			anyToBuff(c, dest, " ")
			anyToBuff(c, dest, D.Invalid)
			return 0, false
		}
	}
	if !hd {
		anyToBuff(c, dest, D.Float)
		anyToBuff(c, dest, " ")
		anyToBuff(c, dest, D.String)
		anyToBuff(c, dest, " ")
		anyToBuff(c, dest, D.Invalid)
		return 0, false
	}

	result := float64(ip)
	if fp > 0 {
		result += float64(fp) / fd
	}

	if isNeg {
		result = -result
	}

	return result, true
}

// stringToIn converts string to number with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
// Phase 13: Optimized to eliminate 82.35% of allocations by removing T() calls
func (t *conv) stringToIn(base int) {
	inp := t.ensureStringInOut()
	// Inline validateBase logic - use static error string
	if base < 2 || base > 36 {
		t.wrErr(errBaseInvalid)
		return
	}

	// Phase 11: Extended fast path for common numbers (0-99999) in base 10
	if base == 10 && len(inp) <= 5 && len(inp) > 0 {
		if num, err := parseSmallInt(inp); err == nil {
			// Store result using anyToBuff to output buffer
			t.rstOut()
			anyToBuff(t, buffOut, uint64(num))
			t.kind = KUint
			return
		}
	}

	var res uint64 // out

	// Phase 8.4: Optimized path for base 10 (most common case)
	if base == 10 {
		// Use direct byte access instead of range to avoid UTF-8 overhead
		for i := 0; i < len(inp); i++ {
			ch := inp[i]
			if ch < '0' || ch > '9' {
				// Phase 13: Use static error string instead of T() + string(ch)
				t.wrErr(errInvalidNonNumeric)
				return
			}

			d := uint64(ch - '0')
			// Check for overflow - optimized for base 10
			if res > 1844674407370955161 || (res == 1844674407370955161 && d > 5) { // (2^64-1)/10
				t.wrErr(errNumberOverflow)
				return
			}

			res = res*10 + d
		}
	} else {
		// Original implementation for other bases - also optimized
		for _, ch := range inp { // char
			var d int // digit
			if ch >= '0' && ch <= '9' {
				d = int(ch - '0')
			} else if ch >= 'a' && ch <= 'z' {
				d = int(ch-'a') + 10
			} else if ch >= 'A' && ch <= 'Z' {
				d = int(ch-'A') + 10
			} else {
				// Phase 13: Use static error string
				t.wrErr(errCharacterInvalid)
				return
			}

			if d >= base {
				t.wrErr(errDigitOutOfRange)
				return
			}
			// Check for overflow
			if res > (^uint64(0))/uint64(base) {
				t.wrErr(errNumberOverflow)
				return
			}

			res = res*uint64(base) + uint64(d)
		}
	}

	// Store result using anyToBuff to output buffer
	t.rstOut()
	anyToBuff(t, buffOut, res)
	t.kind = KUint
}

// fmtIntGeneric converts integer to string and writes to specified buffer destination
// Independent function that receives parameters instead of using temp fields
func fmtIntGeneric(c *conv, val int64, base int, allowNegative bool, dest buffDest) {
	if val == 0 {
		anyToBuff(c, dest, "0")
		return
	}

	var out [64]byte // Max int64 needs 20 digits + sign, base 2 needs 64 digits
	negative := allowNegative && val < 0
	if negative {
		val = -val
	}

	// Inlined fmtUint2Str logic
	uval := uint64(val)
	idx := len(out)

	for uval > 0 {
		idx--
		out[idx] = digs[uval%uint64(base)]
		uval /= uint64(base)
	}

	if negative {
		idx--
		out[idx] = '-'
	}

	writeBytesToDest(c, dest, out[idx:])
}

// intTo converts int64 to string and writes to specified buffer destination
// Independent function that receives parameters instead of using temp fields
func intTo(c *conv, val int64, dest buffDest) {
	// Handle common small numbers using lookup table
	if val >= 0 && val < int64(len(smallInts)) {
		anyToBuff(c, dest, smallInts[val])
		return
	}

	// Handle special cases
	if val == 0 {
		anyToBuff(c, dest, "0")
		return
	}
	if val == 1 {
		anyToBuff(c, dest, "1")
		return
	}
	// Fall back to standard conversion for larger numbers
	fmtIntGeneric(c, val, 10, true, dest)
}

// uint64To converts uint64 to string and writes to specified buffer destination
// Independent function that receives parameters instead of using temp fields
func uint64To(c *conv, val uint64, dest buffDest) {
	// Handle common small numbers using lookup table
	if val < uint64(len(smallInts)) {
		anyToBuff(c, dest, smallInts[val])
		return
	}

	// Handle special cases
	if val == 0 {
		anyToBuff(c, dest, "0")
		return
	}
	if val == 1 {
		anyToBuff(c, dest, "1")
		return
	}

	// Fall back to standard conversion for larger numbers
	fmtIntGeneric(c, int64(val), 10, false, dest)
}

// floatTo converts float to string and writes to specified buffer destination
// Independent function that receives parameters instead of using temp fields
func floatTo(c *conv, val float64, dest buffDest) {
	// Handle special cases using centralized buffer methods
	if val != val { // NaN
		anyToBuff(c, dest, "NaN")
		return
	}

	// Handle infinity
	if val > 1e308 || val < -1e308 {
		if val < 0 {
			anyToBuff(c, dest, "-Inf")
		} else {
			anyToBuff(c, dest, "+Inf")
		}
		return
	}

	// Handle zero (reuse existing constants)
	if val == 0 {
		anyToBuff(c, dest, "0")
		return
	}

	// Handle negative sign
	isNegative := val < 0
	if isNegative {
		val = -val
		anyToBuff(c, dest, "-")
	}

	// Extract integer and fractional parts
	integerPart := int64(val)
	fractionalPart := val - float64(integerPart)

	// Convert integer part using existing optimized methods
	if integerPart >= 0 && integerPart < int64(len(smallInts)) {
		// Use lookup table for small integers
		anyToBuff(c, dest, smallInts[integerPart])
	} else {
		// Direct integer conversion without temp fields
		writeIntToDest(c, dest, integerPart)
	}

	// Add fractional part if significant
	if fractionalPart > 1e-6 { // Avoid tiny fractions
		anyToBuff(c, dest, ".")

		// Convert fractional part directly to buffer (no intermediate strings)
		multiplier := 1e6
		fracPart := int64(fractionalPart*multiplier + 0.5)

		// Write fractional digits directly to buffer
		var fracBuf [6]byte // Fixed size for 6 decimal places
		digits := 0
		tempFrac := fracPart

		// Build digits in reverse order
		for i := 5; i >= 0 && tempFrac > 0; i-- {
			fracBuf[i] = byte('0' + tempFrac%10)
			tempFrac /= 10
			digits++
		}

		// Remove trailing zeros
		for digits > 1 && fracBuf[6-digits] == '0' {
			digits--
		}

		// Write only significant fractional digits
		if digits > 0 {
			writeBytesToDest(c, dest, fracBuf[6-digits:6])
		}
	}
}
