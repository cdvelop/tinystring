package tinystring

// Split divides a string by a separator and returns a slice of substrings
// Phase 11 Optimization: Reduced allocations through buffer pooling and optimized algorithms
// If no separator is provided, splits by whitespace (similar to strings.Fields)
// Note: When using a specific separator, strings shorter than 3 characters are returned as is
// eg: Split("Hello World") => []string{"Hello", "World"}
// with separator eg: Split("Hello;World", ";") => []string{"Hello", "World"}

func Split(data string, separator ...string) (result []string) {
	// If no separator provided, split by whitespace
	if len(separator) == 0 {
		// Inline splitByWhitespace logic
		if len(data) == 0 {
			return []string{}
		}

		// Pre-scan to count words for exact capacity
		wordCount := 0
		inWord := false

		for _, ch := range data {
			isSpace := ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
			if !isSpace && !inWord {
				wordCount++
				inWord = true
			} else if isSpace {
				inWord = false
			}
		}

		if wordCount == 0 {
			return []string{}
		}

		// Allocate exact capacity to avoid reallocations
		result := make([]string, 0, wordCount)
		inWord = false
		start := 0

		// Second pass: extract words
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

		return result
	}

	// Using the provided separator
	sep := separator[0]

	// Don't split short strings when using a custom separator
	if len(data) < 3 {
		return []string{data}
	}

	// Handle empty separator
	if len(sep) == 0 {
		// Inline splitByCharacter logic
		if len(data) == 0 {
			return []string{}
		}

		// Pre-allocate exact capacity for character count
		result := make([]string, 0, len(data))

		// Use direct byte access for ASCII optimization
		for i := 0; i < len(data); i++ {
			if data[i] < 128 { // ASCII fast path
				result = append(result, data[i:i+1])
			} else {
				// UTF-8 handling (fallback)
				for j, ch := range data[i:] {
					result = append(result, string(ch))
					i += j
					break
				}
			}
		}

		return result
	}

	// Inline splitBySeparator logic
	if len(data) == 0 {
		return []string{""}
	}

	// Pre-scan to count parts for exact capacity
	partCount := 1
	sepLen := len(sep)

	for i := 0; i <= len(data)-sepLen; i++ {
		if data[i:i+sepLen] == sep {
			partCount++
			i += sepLen - 1 // Skip separator
		}
	}

	// Allocate exact capacity
	result = make([]string, 0, partCount)
	start := 0

	// Extract parts
	for i := 0; i <= len(data)-sepLen; i++ {
		if data[i:i+sepLen] == sep {
			result = append(result, data[start:i])
			start = i + sepLen
			i += sepLen - 1 // Skip the characters we just checked
		}
	}
	// Add the remaining substring
	result = append(result, data[start:])
	return result
}
