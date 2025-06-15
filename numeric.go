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

// Shared helper methods to reduce code duplication between numeric.go and format.go

// saveState saves the current string value and type for later restoration
func (t *conv) saveState() (string, vTpe) {
	return t.stringVal, t.vTpe
}

// restoreState restores previously saved string value and type
func (t *conv) restoreState(savedVal string, savedType vTpe) {
	t.stringVal = savedVal
	t.vTpe = savedType
	t.err = "" // Reset error when restoring state
}

// tryParseAs attempts to parse content as the specified numeric type with fallback to float
func (t *conv) tryParseAs(parseType vTpe, base int) bool {
	// Save original state
	oSV, oVT := t.saveState()
	// Try direct parsing based on type
	switch parseType {
	case typeInt:
		t.s2IntGeneric(base)
	case typeUint:
		t.s2Uint(oSV, base)
	}

	if t.err == "" {
		return true
	}

	// Check if the error is due to invalid base with negative numbers
	// If so, don't attempt float fallback as it would bypass base validation
	if base != 10 && len(oSV) > 0 && oSV[0] == '-' {
		return false
	}
	// If that fails, restore state and try to parse as float then convert
	t.restoreState(oSV, oVT)
	t.s2Float()
	if t.err == "" {
		switch parseType {
		case typeInt:
			t.intVal = int64(t.floatVal)
			t.vTpe = typeInt
		case typeUint:
			if t.floatVal < 0 {
				t.err = errNegativeUnsigned
				return false
			}
			t.uintVal = uint64(t.floatVal)
			t.vTpe = typeUint
		}
		return true
	}

	return false
}

// validateBase validates that base is within acceptable range (2-36)
func (t *conv) validateBase(base int) bool {
	if base < 2 || base > 36 {
		t.NewErr(errInvalidBase)
		return false
	}
	return true
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
// negative signs will result in an error.
func (t *conv) ToInt(base ...int) (int, error) {
	if t.err != "" {
		return 0, t
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	switch t.vTpe {
	case typeInt:
		return int(t.intVal), nil // Direct return for integer values
	case typeUint:
		return int(t.uintVal), nil // Direct conversion from uint
	case typeFloat:
		return int(t.floatVal), nil // Direct truncation from float
	default:
		// For string types and other types, use shared helper method for parsing with fallback
		if t.tryParseAs(typeInt, b) {
			return int(t.intVal), nil
		}
		// Return error if parsing failed
		return 0, t
	}
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
func (t *conv) ToInt64(base ...int) (int64, error) {
	if t.err != "" {
		return 0, t
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	// Use shared helper method for parsing with fallback
	if t.tryParseAs(typeInt, b) {
		return t.intVal, nil
	}

	// Return error if parsing failed
	return 0, t
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
// always result in an error, regardless of the base.
func (t *conv) ToUint(base ...int) (uint, error) {
	if t.err != "" {
		return 0, t
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}

	switch t.vTpe {
	case typeUint:
		return uint(t.uintVal), nil // Direct return for uint values
	case typeInt:
		if t.intVal < 0 {
			t.NewErr(errNegativeUnsigned, t.intVal)
			return 0, t
		}
		return uint(t.intVal), nil // Direct conversion from int if positive
	case typeFloat:
		if t.floatVal < 0 {
			t.NewErr(errNegativeUnsigned, t.floatVal)
			return 0, t
		}
		return uint(t.floatVal), nil // Direct truncation from float if positive
	default:
		// For string types and other types, use shared helper method for parsing with fallback
		if t.tryParseAs(typeUint, b) {
			return uint(t.uintVal), nil
		}
		// Return error if parsing failed
		return 0, t
	}
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
func (t *conv) ToFloat() (float64, error) {
	if t.err != "" {
		return 0, t
	}

	switch t.vTpe {
	case typeFloat:
		return t.floatVal, nil // Direct return for float values	case typeInt:
	case typeUint:
		return float64(t.uintVal), nil // Direct conversion from uint
	default: // For string types and other types, parse as float
		t.s2Float()
		if t.err != "" {
			return 0, t
		}
		return t.floatVal, nil
	}
}

// s2IntGeneric converts string to signed integer with specified base and stores in conv struct.
// This unified method handles both int and int64 conversions.
func (t *conv) s2IntGeneric(base int) {
	inp := t.getString()
	if t.isEmptyString(inp) {
		return
	}

	isNeg := false
	if inp[0] == '-' {
		if base != 10 {
			t.NewErr(errInvalidBase, "negative numbers are not supported for non-decimal bases")
			return
		}
		isNeg = true // Update the conv struct with the string without the negative sign
		t.setString(inp[1:])
	}
	t.s2n(base)
	if t.err != "" {
		return
	}

	if isNeg {
		t.intVal = -int64(t.uintVal)
	} else {
		t.intVal = int64(t.uintVal)
	}
	t.vTpe = typeInt
}

// s2Uint converts string to uint with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) s2Uint(input string, base int) {
	if t.isEmptyString(input) {
		return
	}

	if input[0] == '-' {
		t.NewErr(errNegativeUnsigned)
		return
	}
	// Update the conv struct with the provided input string
	t.setString(input)
	t.s2n(base)
	// Result already stored in t.uintVal by s2n
}

// s2Float converts string to float64 and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) s2Float() {
	inp := t.getString()
	if t.isEmptyString(inp) {
		return
	}

	isNeg := false
	sIdx := 0
	if inp[0] == '-' {
		isNeg = true
		sIdx = 1
		if len(inp) == 1 { // Just a "-" sign
			t.err = errInvalidFloat
			return
		}
	} else if inp[0] == '+' {
		sIdx = 1
		if len(inp) == 1 { // Just a "+" sign
			t.err = errInvalidFloat
			return
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
				t.err = errInvalidFloat
				return
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
					t.err = errOverflow
					return
				}
				ip = ip*10 + dgt
			}
		} else {
			t.err = errInvalidFloat
			return
		}
	}

	if !hd {
		t.err = errInvalidFloat
		return
	}

	result := float64(ip)
	if fp > 0 {
		result += float64(fp) / fd
	}

	if isNeg {
		result = -result
	}
	t.floatVal = result
	t.vTpe = typeFloat
}

// s2n converts string to number with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) s2n(base int) {
	inp := t.getString()
	if !t.validateBase(base) {
		return
	}

	var res uint64 // result

	for _, ch := range inp { // char
		var d int // digit
		if ch >= '0' && ch <= '9' {
			d = int(ch - '0')
		} else if ch >= 'a' && ch <= 'z' {
			d = int(ch-'a') + 10
		} else if ch >= 'A' && ch <= 'Z' {
			d = int(ch-'A') + 10
		} else {
			t.NewErr(string(ch))
			return
		}

		if d >= base {
			t.NewErr(errInvalidBase, "digit out of range for base")
			return
		}

		// Check for overflow
		if res > (^uint64(0))/uint64(base) {
			t.NewErr(errOverflow)
			return
		}

		res = res*uint64(base) + uint64(d)
	}

	t.uintVal = res
	t.vTpe = typeUint
}

// isEmptyString checks if the input string is empty and sets an error if it is.
func (t *conv) isEmptyString(inp string) bool {
	if isEmpty(inp) {
		t.err = errEmptyString
		return true
	}
	return false
}

// fmtUint2Str converts an unsigned integer to a string in the given base,
// writing it into the provided buffer and returning the resulting string and its start index in the buffer.
func fmtUint2Str(val uint64, base int, buf []byte) (string, int) {
	idx := len(buf) // i

	for val > 0 {
		idx--
		buf[idx] = digs[val%uint64(base)]
		val /= uint64(base)
	}
	return string(buf[idx:]), idx
}

// fmtInt converts integer to string and stores in tmpStr
func (t *conv) fmtInt(base int) {
	t.fmtIntGeneric(t.intVal, base, true)
}

// fmtUint converts unsigned integer to string and stores in tmpStr
func (t *conv) fmtUint(base int) {
	t.fmtIntGeneric(int64(t.uintVal), base, false)
}

// fmtIntGeneric converts integer to string with unified logic
func (t *conv) fmtIntGeneric(val int64, base int, allowNegative bool) {
	if val == 0 {
		t.tmpStr = zeroStr
		t.stringVal = t.tmpStr
		return
	}

	var buf [64]byte // Max int64 needs 20 digits + sign, base 2 needs 64 digits
	negative := allowNegative && val < 0
	if negative {
		val = -val
	}

	_, idx := fmtUint2Str(uint64(val), base, buf[:])

	if negative {
		idx--
		buf[idx] = '-'
	}

	t.tmpStr = string(buf[idx:])
	t.stringVal = t.tmpStr
}

// i2s converts int64 to string with minimal allocations and stores in tmpStr
func (t *conv) i2s() {
	val := t.intVal
	// Handle common small numbers using lookup table
	if val >= 0 && val < int64(len(smallInts)) {
		t.tmpStr = smallInts[val]
		t.stringVal = t.tmpStr
		return
	}

	// Handle special cases
	if val == 0 {
		t.tmpStr = zeroStr
		t.stringVal = t.tmpStr
		return
	}
	if val == 1 {
		t.tmpStr = oneStr
		t.stringVal = t.tmpStr
		return
	}
	// Fall back to standard conversion for larger numbers
	t.fmtInt(10)
}

// u2s converts uint64 to string with minimal allocations and stores in tmpStr
func (t *conv) u2s() {
	val := t.uintVal
	// Handle common small numbers using lookup table
	if val < uint64(len(smallInts)) {
		t.tmpStr = smallInts[val]
		t.stringVal = t.tmpStr
		return
	}

	// Handle special cases
	if val == 0 {
		t.tmpStr = zeroStr
		t.stringVal = t.tmpStr
		return
	}
	if val == 1 {
		t.tmpStr = oneStr
		t.stringVal = t.tmpStr
		return
	}

	// Fall back to standard conversion for larger numbers
	t.fmtUint(10) // Use the unified fmtUint
}

// f2s converts float to string and stores in tmpStr
// Uses simplified float-to-string conversion for basic numeric.go operations
func (t *conv) f2s() {
	val := t.floatVal
	// Handle special cases
	if val != val { // NaN
		t.tmpStr = "NaN"
		t.stringVal = t.tmpStr
		return
	}

	// Handle infinity
	if val > 1e308 || val < -1e308 {
		if val < 0 {
			t.tmpStr = "-Inf"
		} else {
			t.tmpStr = "Inf"
		}
		t.stringVal = t.tmpStr
		return
	}
	// Handle zero
	if val == 0 {
		t.tmpStr = zeroStr
		t.stringVal = t.tmpStr
		return
	}

	// Simple float-to-string conversion for basic cases
	isNegative := val < 0
	if isNegative {
		val = -val
	}
	// Extract integer and fractional parts
	integerPart := int64(val)
	fractionalPart := val - float64(integerPart)

	// Convert integer part
	var result string
	if integerPart == 0 {
		result = zeroStr
	} else {
		// Convert integer to string
		temp := integerPart
		digits := ""
		for temp > 0 {
			digits = string(byte('0'+temp%10)) + digits
			temp /= 10
		}
		result = digits
	}

	// Add fractional part if significant
	if fractionalPart > 0 {
		result += "."
		// Add up to 6 significant fractional digits, removing trailing zeros
		multiplier := 1e6
		fracPart := int64(fractionalPart*multiplier + 0.5)

		fracStr := ""
		for i := 0; i < 6; i++ {
			fracStr = string(byte('0'+fracPart%10)) + fracStr
			fracPart /= 10
		}

		// Remove trailing zeros
		for len(fracStr) > 1 && fracStr[len(fracStr)-1] == '0' {
			fracStr = fracStr[:len(fracStr)-1]
		}

		result += fracStr
	}

	if isNegative {
		result = "-" + result
	}

	t.tmpStr = result
	t.stringVal = result
}

// validateIntParam validates and converts an any parameter to int
func (t *conv) validateIntParam(param any, allowZero bool) (int, bool) {
	val, ok := extractInt(param)
	if !ok {
		return 0, false
	}
	// Unified validation logic
	if allowZero {
		return val, val >= 0
	}
	return val, val > 0
}

// extractInt extracts integer value from any numeric type using generics
func extractInt(v any) (int, bool) {
	switch val := v.(type) {
	case int, int8, int16, int32, int64:
		// Use type assertion to handle all integer types
		if i, ok := v.(int); ok {
			return i, true
		}
		if i8, ok := v.(int8); ok {
			return int(i8), true
		}
		if i16, ok := v.(int16); ok {
			return int(i16), true
		}
		if i32, ok := v.(int32); ok {
			return int(i32), true
		}
		if i64, ok := v.(int64); ok {
			return int(i64), true
		}
	case uint, uint8, uint16, uint32, uint64:
		if u, ok := v.(uint); ok {
			return int(u), true
		}
		if u8, ok := v.(uint8); ok {
			return int(u8), true
		}
		if u16, ok := v.(uint16); ok {
			return int(u16), true
		}
		if u32, ok := v.(uint32); ok {
			return int(u32), true
		}
		if u64, ok := v.(uint64); ok {
			return int(u64), true
		}
	case float32:
		return int(val), true
	case float64:
		return int(val), true
	}
	return 0, false
}
