package tinystring

import (
	"strconv"
)

// RoundDecimals formats a float value by rounding it to the specified number of decimal places
// Example: Convert(3.12221).RoundDecimals(2).String() returns "3.12"
func (t *Text) RoundDecimals(decimals int) *Text {
	// Try to parse the string as float
	val, err := strconv.ParseFloat(t.content, 64)
	if err != nil {
		return t // Return unmodified if not a valid float
	}

	// Optimization for zero decimal case
	if decimals <= 0 {
		// Round to integer
		if val < 0 {
			val -= 0.5
		} else {
			val += 0.5
		}
		intVal := int64(val)
		t.content = intToStringOptimized(intVal)
		return t
	}

	// Manual rounding implementation without math package
	factor := float64(1)
	for i := 0; i < decimals; i++ {
		factor *= 10
	}

	// Scale, round, then scale back
	scaled := val * factor
	if scaled < 0 {
		scaled -= 0.5
	} else {
		scaled += 0.5
	}

	// Convert to integer and back to remove excess precision
	rounded := float64(int64(scaled)) / factor

	// Format with the specified number of decimal places
	t.content = formatFloatWithDecimals(rounded, decimals)
	return t
}

// FormatNumber formats a numeric value with thousand separators (dots) and
// removes trailing zeros after the decimal point
// Example: Convert(2189009.00).FormatNumber().String() returns "2.189.009"
func (t *Text) FormatNumber() *Text {
	// Try to parse the string to ensure it's a number
	val, err := strconv.ParseFloat(t.content, 64)
	if err != nil {
		return t // Return unmodified if not a valid float
	}

	// Handle zero case separately
	if val == 0 {
		t.content = "0"
		return t
	}

	// Handle negative numbers
	negative := val < 0
	if negative {
		val = -val
	}

	// Check if the number has decimal part
	hasDecimal := val != float64(int64(val))

	var result string
	if hasDecimal {
		// Format with decimal part, but remove trailing zeros
		strVal := strconv.FormatFloat(val, 'f', 6, 64)
		result = removeTrailingZeros(strVal)
	} else {
		// Format as integer
		intVal := int64(val)
		result = addThousandSeparators(intToStringOptimized(intVal))
	}

	// Add negative sign if needed
	if negative {
		result = "-" + result
	}

	t.content = result
	return t
}

// formatFloatWithDecimals formats a float with a specific number of decimal places
// without using the math library
func formatFloatWithDecimals(val float64, decimals int) string {
	if decimals <= 0 {
		return intToStringOptimized(int64(val))
	}

	// Use strconv with fixed precision
	return strconv.FormatFloat(val, 'f', decimals, 64)
}

// removeTrailingZeros removes unnecessary trailing zeros after the decimal point
// and the decimal point itself if all decimals are zero
func removeTrailingZeros(s string) string {
	// Find the decimal point position
	decimalPos := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			decimalPos = i
			break
		}
	}

	// If no decimal point, add separators to the integer part and return
	if decimalPos == -1 {
		return addThousandSeparators(s)
	}

	// Find the position of the last non-zero digit after decimal point
	lastNonZero := len(s) - 1
	for lastNonZero > decimalPos {
		if s[lastNonZero] != '0' {
			break
		}
		lastNonZero--
	}

	// If all digits after decimal point are zeros, remove the decimal point too
	if lastNonZero == decimalPos {
		return addThousandSeparators(s[:decimalPos])
	}

	// Add separators to the integer part, then append the decimal part
	intPart := s[:decimalPos]
	decPart := s[decimalPos : lastNonZero+1]
	return addThousandSeparators(intPart) + decPart
}

// addThousandSeparators adds dot separators to a numeric string
// Example: "2189009" -> "2.189.009"
func addThousandSeparators(s string) string {
	// Handle empty string
	if s == "" {
		return ""
	}

	// Get the length of the string
	length := len(s)

	// If length is 3 or less, no separators needed
	if length <= 3 {
		return s
	}

	// Calculate result size (original + separators)
	resultSize := length + ((length - 1) / 3)
	resultRunes := make([]rune, resultSize)

	// Fill result from right to left
	sourceIdx := length - 1
	resultIdx := resultSize - 1

	// Add separators every 3 digits
	count := 0
	for sourceIdx >= 0 {
		resultRunes[resultIdx] = rune(s[sourceIdx])
		resultIdx--
		sourceIdx--
		count++

		// Add separator if we've added 3 digits (except at the beginning)
		if count == 3 && sourceIdx >= 0 {
			resultRunes[resultIdx] = '.'
			resultIdx--
			count = 0
		}
	}

	return string(resultRunes)
}
