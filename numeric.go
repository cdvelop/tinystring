package tinystring

import (
	"errors"
)

// ToInt converts the text content to an integer with specified base
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

// ToUint converts the text content to an unsigned integer with specified base
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

// ToFloat converts the text content to a float64
func (t *Text) ToFloat() (float64, error) {
	return stringToFloat(t.content)
}

// FromInt creates a new Text instance from an integer with specified base
func FromInt(value int, base ...int) *Text {
	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	str := intToStringWithBase(int64(value), b)
	return &Text{content: str}
}

// FromUint creates a new Text instance from an unsigned integer with specified base
func FromUint(value uint, base ...int) *Text {
	b := 10 // default base
	if len(base) > 0 {
		b = base[0]
	}
	str := uintToStringWithBase(uint64(value), b)
	return &Text{content: str}
}

// FromFloat creates a new Text instance from a float with specified precision
func FromFloat(value float64, precision ...int) *Text {
	p := -1 // default precision (full precision)
	if len(precision) > 0 {
		p = precision[0]
	}
	str := formatFloatToString(value, p)
	return &Text{content: str}
}

// stringToInt converts a string to an integer with specified base (integrated from tinystrconv)
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

// stringToUint converts a string to an unsigned integer with specified base (integrated from tinystrconv)
func stringToUint(input string, base int) (uint64, error) {
	if input == "" {
		return 0, errors.New("empty string")
	}

	if input[0] == '-' {
		return 0, errors.New("negative numbers are not supported for unsigned integers")
	}

	return stringToNumberHelper(input, base)
}

// stringToFloat converts a string to a float64 (manual implementation, integrated from tinystrconv)
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

// stringToNumberHelper converts a string to a number with specified base (integrated from tinystrconv)
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

// uintToStringWithBase converts an unsigned integer to string with specified base
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
