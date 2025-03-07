package tinystring

// Split divides a string by a separator and returns a slice of substrings
// Note: Strings shorter than 3 characters are returned as is
func Split(data, separator string) (result []string) {
	// Don't split short strings
	if len(data) < 3 {
		return []string{data}
	}

	// Handle empty separator
	if separator == "" {
		for _, ch := range data {
			result = append(result, string(ch))
		}
		return
	}

	start := 0
	sepLen := len(separator)

	for i := 0; i <= len(data)-sepLen; i++ {
		if data[i:i+sepLen] == separator {
			result = append(result, data[start:i])
			start = i + sepLen
			i += sepLen - 1 // Skip the characters we just checked
		}
	}

	// Add the remaining substring
	result = append(result, data[start:])
	return
}
