package tinystring

import (
	"errors"
)

// ToInt converts the text content to an integer with optional base specification.
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
func (t *Text) ToInt(base ...int) (int, error) {
	if t.err != nil {
		return 0, t.err
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}

	// First try to parse as int directly
	val, err := stringToInt(t.content, b)
	if err == nil {
		return val, nil
	}

	// If that fails, try to parse as float and then convert to int (for truncation)
	floatVal, floatErr := stringToFloat(t.content)
	if floatErr == nil {
		return int(floatVal), nil
	}

	// Return the original int parsing error
	return 0, err
}

// ToUint converts the text content to an unsigned integer with optional base specification.
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
func (t *Text) ToUint(base ...int) (uint, error) {
	if t.err != nil {
		return 0, t.err
	}

	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}

	// First try to parse as uint directly
	val, err := stringToUint(t.content, b)
	if err == nil {
		return uint(val), nil
	}

	// If that fails, try to parse as float and then convert to uint (for truncation)
	floatVal, floatErr := stringToFloat(t.content)
	if floatErr == nil {
		if floatVal < 0 {
			return 0, errors.New("negative numbers are not supported for unsigned integers")
		}
		return uint(floatVal), nil
	}

	// Return the original uint parsing error
	return 0, err
}

// ToFloat converts the text content to a float64 (double precision floating point).
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
func (t *Text) ToFloat() (float64, error) {
	return stringToFloat(t.content)
}

// FromInt creates a new Text instance from an integer with optional base specification.
//
// Parameters:
//   - value: The integer value to convert to Text
//   - base (optional): The numeric base for string representation (2-36). Default is 10 (decimal).
//     Common bases: 2 (binary), 8 (octal), 10 (decimal), 16 (hexadecimal)
//
// Returns:
//   - *Text: A new Text instance containing the string representation of the integer
//
// Usage examples:
//
//	// Basic decimal conversion
//	text := FromInt(123)
//	result := text.String() // "123"
//
//	// Binary representation
//	text := FromInt(10, 2)
//	result := text.String() // "1010"
//
//	// Hexadecimal representation
//	text := FromInt(255, 16)
//	result := text.String() // "ff"
//
//	// Octal representation
//	text := FromInt(64, 8)
//	result := text.String() // "100"
//
//	// Negative numbers (only base 10)
//	text := FromInt(-42)
//	result := text.String() // "-42"
//
// Note: The resulting Text instance can be used for further string manipulations
// like case conversion, joining, etc.
func FromInt(value int, base ...int) *Text {
	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	str := intToStringWithBase(int64(value), b)
	return &Text{content: str}
}

// FromUint creates a new Text instance from an unsigned integer with optional base specification.
//
// Parameters:
//   - value: The unsigned integer value to convert to Text
//   - base (optional): The numeric base for string representation (2-36). Default is 10 (decimal).
//     Common bases: 2 (binary), 8 (octal), 10 (decimal), 16 (hexadecimal)
//
// Returns:
//   - *Text: A new Text instance containing the string representation of the unsigned integer
//
// Usage examples:
//
//	// Basic decimal conversion
//	text := FromUint(123)
//	result := text.String() // "123"
//
//	// Binary representation
//	text := FromUint(10, 2)
//	result := text.String() // "1010"
//
//	// Hexadecimal representation
//	text := FromUint(255, 16)
//	result := text.String() // "ff"
//
//	// Large unsigned values
//	text := FromUint(18446744073709551615) // max uint64
//	result := text.String() // "18446744073709551615"
//
//	// Chain with other operations
//	text := FromUint(255, 16).ToUpper()
//	result := text.String() // "FF"
//
// Note: Unlike FromInt, this function only works with non-negative values.
// The resulting Text instance can be used for further string manipulations.
func FromUint(value uint, base ...int) *Text {
	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	str := uintToStringWithBase(uint64(value), b)
	return &Text{content: str}
}

// FromFloat creates a new Text instance from a float64 with optional precision specification.
//
// Parameters:
//   - value: The float64 value to convert to Text
//   - precision (optional): Number of decimal places to include. Default is -1 (full precision).
//     Use 0 for no decimal places, positive values for fixed decimal places.
//
// Returns:
//   - *Text: A new Text instance containing the string representation of the float
//
// Usage examples:
//
//	// Full precision (default)
//	text := FromFloat(123.456789)
//	result := text.String() // "123.456789" (or full precision representation)
//
//	// No decimal places
//	text := FromFloat(123.456, 0)
//	result := text.String() // "123"
//
//	// Fixed decimal places
//	text := FromFloat(123.456, 2)
//	result := text.String() // "123.46" (rounded)
//
//	// Scientific notation values
//	text := FromFloat(1.23e-4, 6)
//	result := text.String() // "0.000123"
//
//	// Negative numbers
//	text := FromFloat(-99.99, 1)
//	result := text.String() // "-100.0"
//
//	// Chain with other operations
//	text := FromFloat(123.456, 2).ToUpper()
//	result := text.String() // "123.46" (case conversion doesn't affect numbers)
//
// Note: The precision parameter controls the number of digits after the decimal point.
// The resulting Text instance can be used for further string manipulations.
func FromFloat(value float64, precision ...int) *Text {
	p := -1 // default precision (full precision)
	if len(precision) > 0 {
		p = precision[0]
	}
	str := formatFloatToString(value, p)
	return &Text{content: str}
}

// stringToInt converts a string to an integer with specified base.
// This is an internal helper function integrated from tinystrconv.
//
// Parameters:
//   - input: The string to convert
//   - base: The numeric base (2-36)
//
// Returns:
//   - int: The converted integer value
//   - error: Any error that occurred during conversion
//
// Supports negative numbers only for base 10.
func stringToInt(input string, base int) (int, error) {
	if input == "" {
		return 0, errors.New("empty string")
	}

	isNegative := false
	if input[0] == '-' {
		if base != 10 {
			return 0, errors.New("negative numbers are not supported for non-decimal bases")
		}
		isNegative = true
		input = input[1:]
	}

	number, err := stringToNumberHelper(input, base)
	if err != nil {
		return 0, err
	}

	if isNegative {
		return -int(number), nil
	}
	return int(number), nil
}

// stringToUint converts a string to an unsigned integer with specified base.
// This is an internal helper function integrated from tinystrconv.
//
// Parameters:
//   - input: The string to convert
//   - base: The numeric base (2-36)
//
// Returns:
//   - uint64: The converted unsigned integer value
//   - error: Any error that occurred during conversion
//
// Does not support negative numbers.
func stringToUint(input string, base int) (uint64, error) {
	if input == "" {
		return 0, errors.New("empty string")
	}

	if input[0] == '-' {
		return 0, errors.New("negative numbers are not supported for unsigned integers")
	}

	return stringToNumberHelper(input, base)
}

// stringToFloat converts a string to a float64 using manual implementation.
// This is an internal helper function integrated from tinystrconv.
//
// Parameters:
//   - input: The string to convert to float64
//
// Returns:
//   - float64: The converted floating point value
//   - error: Any error that occurred during conversion
//
// Supports positive and negative numbers with decimal points.
// Handles basic floating point formats but may have different precision
// characteristics compared to the standard library.
func stringToFloat(input string) (float64, error) {
	if input == "" {
		return 0, errors.New("empty string")
	}

	isNegative := false
	startIndex := 0
	if input[0] == '-' {
		isNegative = true
		startIndex = 1
		if len(input) == 1 { // Just a "-" sign
			return 0, errors.New("invalid float string: only sign")
		}
	} else if input[0] == '+' {
		startIndex = 1
		if len(input) == 1 { // Just a "+" sign
			return 0, errors.New("invalid float string: only sign")
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
				return 0, errors.New("invalid float string: multiple decimal points")
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
					return 0, errors.New("integer part overflow")
				}
				integerPart = integerPart*10 + digit
			}
		} else {
			// Invalid character
			return 0, errors.New("invalid character in float string: " + string(char))
		}
	}

	if !hasDigits {
		return 0, errors.New("invalid float string: no digits found")
	}

	result := float64(integerPart)
	if fractionPart > 0 {
		result += float64(fractionPart) / fractionDivisor
	}

	if isNegative {
		result = -result
	}

	return result, nil
}

// stringToNumberHelper converts a string to a number with specified base.
// This is an internal helper function integrated from tinystrconv.
//
// Parameters:
//   - input: The string to convert (must contain only valid digits for the base)
//   - base: The numeric base (must be between 2 and 36)
//
// Returns:
//   - uint64: The converted number value
//   - error: Any error that occurred during conversion
//
// Used internally by stringToInt and stringToUint functions.
// Validates base range and character validity for the specified base.
func stringToNumberHelper(input string, base int) (uint64, error) {
	if base < 2 || base > 36 {
		return 0, errors.New("base must be between 2 and 36")
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
			return 0, errors.New("invalid character in number: " + string(char))
		}

		if digit >= base {
			return 0, errors.New("digit out of range for base")
		}

		// Check for overflow
		if result > (^uint64(0))/uint64(base) {
			return 0, errors.New("number too large")
		}

		result = result*uint64(base) + uint64(digit)
	}

	return result, nil
}

// uintToStringWithBase converts an unsigned integer to string with specified base.
// This is an internal helper function used by FromInt and FromUint.
//
// Parameters:
//   - number: The unsigned integer to convert
//   - base: The numeric base for the string representation (2-36)
//
// Returns:
//   - string: The string representation of the number in the specified base
//
// Uses lowercase letters (a-z) for digits above 9 in bases greater than 10.
// For example, base 16 uses digits 0-9 and letters a-f.
func uintToStringWithBase(number uint64, base int) string {
	if number == 0 {
		return "0" // Directly return "0" for zero value
	}

	// Max uint64 is 18446744073709551615, which has 20 digits.
	// For base 2, uint64 needs up to 64 bits.
	var buf [64]byte // Max buffer size for uint64 in base 2
	i := len(buf)    // Start from the end of the buffer

	const digitChars = "0123456789abcdefghijklmnopqrstuvwxyz"

	for number > 0 {
		i--
		buf[i] = digitChars[number%uint64(base)]
		number /= uint64(base)
	}

	return string(buf[i:]) // Convert the relevant part of the buffer to a string
}
