package tinystring

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
	if t.err != nil {
		return 0, t.err
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	// First try to parse as int directly
	t.stringToInt(b)
	if t.err == nil {
		return int(t.intVal), nil
	}
	intErr := t.err

	// If that fails, try to parse as float and then convert to int (for truncation)
	t.err = nil // Reset error for float parsing
	t.stringToFloat()
	if t.err == nil {
		return int(t.floatVal), nil
	}

	// Return the original int parsing error
	return 0, intErr
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
	if t.err != nil {
		return 0, t.err
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}

	// First try to parse as int64 directly
	t.stringToInt64(b)
	if t.err == nil {
		return t.intVal, nil
	}
	int64Err := t.err

	// If that fails, try to parse as float and then convert to int64 (for truncation)
	t.err = nil // Reset error for float parsing
	t.stringToFloat()
	if t.err == nil {
		return int64(t.floatVal), nil
	}

	// Return the original int64 parsing error
	return 0, int64Err
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
	if t.err != nil {
		return 0, t.err
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	// First try to parse as uint directly
	t.stringToUint(t.getString(), b)
	if t.err == nil {
		return uint(t.uintVal), nil
	}
	uintErr := t.err

	// If that fails, try to parse as float and then convert to uint (for truncation)
	t.err = nil // Reset error for float parsing
	t.stringToFloat()
	if t.err == nil {
		if t.floatVal < 0 {
			return 0, newError(errNegativeUnsigned)
		}
		return uint(t.floatVal), nil
	}

	// Return the original uint parsing error
	return 0, uintErr
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
	t.stringToFloat()
	if t.err != nil {
		return 0, t.err
	}
	return t.floatVal, nil
}

// stringToInt converts string to int with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) stringToInt(base int) {
	input := t.getString()
	if input == "" {
		t.err = newEmptyStringError()
		return
	}

	isNegative := false
	if input[0] == '-' {
		if base != 10 {
			t.err = newError(errInvalidBase, "negative numbers are not supported for non-decimal bases")
			return
		}
		isNegative = true
		// Update the conv struct with the string without the negative sign
		t.setString(input[1:])
	}

	t.stringToNumberHelper(base)
	if t.err != nil {
		return
	}

	if isNegative {
		t.intVal = -int64(t.uintVal)
	} else {
		t.intVal = int64(t.uintVal)
	}
	t.valType = valTypeInt
}

// stringToInt64 converts string to int64 with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) stringToInt64(base int) {
	input := t.getString()
	if input == "" {
		t.err = newEmptyStringError()
		return
	}

	isNegative := false
	if input[0] == '-' {
		if base != 10 {
			t.err = newInvalidNumberError("negative numbers are not supported for non-decimal bases")
			return
		}
		isNegative = true
		// Update the conv struct with the string without the negative sign
		t.setString(input[1:])
	}

	t.stringToNumberHelper(base)
	if t.err != nil {
		return
	}

	if isNegative {
		t.intVal = -int64(t.uintVal)
	} else {
		t.intVal = int64(t.uintVal)
	}
	t.valType = valTypeInt
}

// stringToUint converts string to uint with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) stringToUint(input string, base int) {
	if input == "" {
		t.err = newEmptyStringError()
		return
	}

	if input[0] == '-' {
		t.err = newError(errNegativeUnsigned)
		return
	}

	// Update the conv struct with the provided input string
	t.setString(input)
	t.stringToNumberHelper(base)
	// Result already stored in t.uintVal by stringToNumberHelper
}

// stringToFloat converts string to float64 and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) stringToFloat() {
	input := t.getString()
	if input == "" {
		t.err = newEmptyStringError()
		return
	}

	isNegative := false
	startIndex := 0
	if input[0] == '-' {
		isNegative = true
		startIndex = 1
		if len(input) == 1 { // Just a "-" sign
			t.err = newInvalidFloatError()
			return
		}
	} else if input[0] == '+' {
		startIndex = 1
		if len(input) == 1 { // Just a "+" sign
			t.err = newInvalidFloatError()
			return
		}
	}

	var integerPart uint64
	var fractionPart uint64
	var fractionDivisor float64 = 1.0
	decimalPointSeen := false
	parsingFraction := false
	hasDigits := false // To check if there's at least one digit

	for i := startIndex; i < len(input); i++ {
		char := input[i]
		if char == '.' {
			if decimalPointSeen {
				t.err = newInvalidFloatError()
				return
			}
			decimalPointSeen = true
			parsingFraction = true
		} else if char >= '0' && char <= '9' {
			hasDigits = true
			digit := uint64(char - '0')
			if parsingFraction {
				// Check for overflow before multiplying fractionPart
				if fractionPart > (^uint64(0))/10 {
					// This condition might be too strict for typical float precision needs,
					// but good for catching extreme cases.
					// Standard library might handle this by losing precision.
					// For now, let's consider it an error or rely on float64 limits.
				}
				fractionPart = fractionPart*10 + digit
				fractionDivisor *= 10
			} else {
				// Check for overflow before multiplying integerPart
				if integerPart > (^uint64(0)-digit)/10 {
					t.err = newError(errOverflow)
					return
				}
				integerPart = integerPart*10 + digit
			}
		} else {
			// Invalid character
			t.err = newInvalidNumberError(string(char))
			return
		}
	}

	if !hasDigits {
		t.err = newInvalidFloatError()
		return
	}

	result := float64(integerPart)
	if fractionPart > 0 {
		result += float64(fractionPart) / fractionDivisor
	}

	if isNegative {
		result = -result
	}

	t.floatVal = result
	t.valType = valTypeFloat
}

// stringToNumberHelper converts string to number with specified base and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) stringToNumberHelper(base int) {
	input := t.getString()
	if base < 2 || base > 36 {
		t.err = newError(errInvalidBase)
		return
	}

	var result uint64

	for _, char := range input {
		var digit int
		if char >= '0' && char <= '9' {
			digit = int(char - '0')
		} else if char >= 'a' && char <= 'z' {
			digit = int(char-'a') + 10
		} else if char >= 'A' && char <= 'Z' {
			digit = int(char-'A') + 10
		} else {
			t.err = newInvalidNumberError(string(char))
			return
		}

		if digit >= base {
			t.err = newError(errInvalidBase, "digit out of range for base")
			return
		}

		// Check for overflow
		if result > (^uint64(0))/uint64(base) {
			t.err = newError(errOverflow)
			return
		}

		result = result*uint64(base) + uint64(digit)
	}

	t.uintVal = result
	t.valType = valTypeUint
}
