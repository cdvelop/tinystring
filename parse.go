package tinystring

import (
	"errors"
)

// ParseKeyValue extracts the value part from a "key:value" formatted string.
// By default, it uses ":" as the delimiter but accepts an optional custom delimiter.
// The function returns the value part and an error (nil if successful).
//
// Examples:
//
//	value, err := ParseKeyValue("name:John")
//	// value = "John", err = nil
//
//	value, err := ParseKeyValue("data=123", "=")
//	// value = "123", err = nil
//
//	value, err := ParseKeyValue("invalid-string")
//	// value = "", err = error containing "delimiter ':' not found in string invalid-string"
func ParseKeyValue(input string, delimiters ...string) (value string, err error) {
	// Default delimiter is ":"
	delimiter := ":"

	// Check for a custom delimiter
	if len(delimiters) > 0 {
		if len(delimiters) > 1 {
			return "", errors.New("only one delimiter is allowed")
		}
		if delimiters[0] != "" {
			delimiter = delimiters[0]
		}
	}

	// Special case: if the input is exactly the delimiter, return empty value without error
	if input == delimiter {
		return "", nil
	}

	// Check if delimiter exists in the input
	if !Contains(input, delimiter) {
		errorMsg := "delimiter '" + delimiter + "' not found in string " + input
		return "", errors.New(errorMsg)
	}
	// Extract value part (everything after the first occurrence of the delimiter)
	// Find the position of the first delimiter
	delimIndex := -1
	for i := 0; i <= len(input)-len(delimiter); i++ {
		if input[i:i+len(delimiter)] == delimiter {
			delimIndex = i
			break
		}
	}

	// Return everything after the first delimiter
	if delimIndex >= 0 {
		return input[delimIndex+len(delimiter):], nil
	}

	// This should never happen if Contains returned true
	return "", nil
}
