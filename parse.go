package tinystring

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
func ParseKeyValue(in string, delimiters ...string) (value string, err error) {
	// Default delimiter is ":"
	d := ":"

	// Check for a custom delimiter
	if len(delimiters) > 0 {
		if len(delimiters) > 1 {
			return "", Err(errInvalidFormat, "only one delimiter is allowed")
		}
		if delimiters[0] != "" {
			d = delimiters[0]
		}
	}

	// Special case: if the in is exactly the delimiter, return empty value without error
	if in == d {
		return "", nil
	}

	// Check if delimiter exists in the in
	if !Contains(in, d) {
		em := "delimiter '" + d + "' not found in string " + in
		return "", Err(errInvalidFormat, em)
	}
	// Extract value part (everything after the first occurrence of the delimiter)
	// Find the position of the first delimiter
	di := -1
	for i := 0; i <= len(in)-len(d); i++ {
		if in[i:i+len(d)] == d {
			di = i
			break
		}
	}

	// Return everything after the first delimiter
	if di >= 0 {
		return in[di+len(d):], nil
	}

	// This should never happen if Contains returned true
	return "", nil
}
