package tinystring

// Split divides a string by a separator and returns a slice of substrings
// If no separator is provided, splits by whitespace (similar to strings.Fields)
// Note: When using a specific separator, strings shorter than 3 characters are returned as is
// eg: Split("Hello World") => []string{"Hello", "World"}
// with separator eg: Split("Hello;World", ";") => []string{"Hello", "World"}

func Split(data string, separator ...string) (result []string) {
	// If no separator provided, split by whitespace
	if len(separator) == 0 {
		inWord := false
		start := 0

		// Iterate through the string character by character
		for i, ch := range data {
			isSpace := ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'

			if !isSpace && !inWord {
				// Start of a new word
				inWord = true
				start = i
			} else if isSpace && inWord {
				// End of a word
				inWord = false
				result = append(result, data[start:i])
			}
		}

		// Handle the last word if the string doesn't end with whitespace
		if inWord {
			result = append(result, data[start:])
		}

		return
	}

	// Using the provided separator
	sep := separator[0]

	// Don't split short strings when using a custom separator
	if len(data) < 3 {
		return []string{data}
	}

	// Handle empty separator
	if sep == "" {
		for _, ch := range data {
			result = append(result, string(ch))
		}
		return
	}

	start := 0
	sepLen := len(sep)

	for i := 0; i <= len(data)-sepLen; i++ {
		if data[i:i+sepLen] == sep {
			result = append(result, data[start:i])
			start = i + sepLen
			i += sepLen - 1 // Skip the characters we just checked
		}
	}

	// Add the remaining substring
	result = append(result, data[start:])
	return
}
