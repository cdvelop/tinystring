package tinystring

// TruncateName truncates names and surnames in a user-friendly way for display in limited spaces
// like chart labels. It adds abbreviation dots where appropriate. This method processes the first
// word differently if there are more than 2 words in the text.
//
// Parameters:
//   - maxCharsPerWord: maximum number of characters to keep per word (any numeric type)
//   - maxWidth: maximum total length for the final string (any numeric type)
//   - text: the text to truncate (string)
//
// Examples:
//   - Convert("").TruncateName(3, 15, "Jeronimo Dominguez") => "Jer. Dominguez"
//   - Convert("").TruncateName(2, 10, "Ana Maria Rodriguez") => "An. Mar..."
//   - Convert("").TruncateName(3, 5, "Juan") => "Juan"
func (t *Text) TruncateName(maxCharsPerWord, maxWidth any, inputText string) *Text {
	if inputText == "" {
		return t
	}

	// Convert maxCharsPerWord to integer using the toInt utility function
	maxChars, ok := toInt(maxCharsPerWord)
	if !ok || maxChars <= 0 {
		return t // Return original if maxCharsPerWord is invalid
	}

	// Convert maxWidth to integer
	maxTotal, ok := toInt(maxWidth)
	if !ok || maxTotal <= 0 {
		return t // Return original if maxWidth is invalid
	}

	// Default suffix
	suffixStr := "..."
	// Split by spaces to get individual words
	words := SplitBySpace(inputText)
	result := make([]string, 0, len(words))
	// Special processing - first word gets truncated if there are more than 2 words
	// Calculate total length before applying any truncation
	totalLength := 0
	// Count spaces
	if len(words) > 0 {
		totalLength = len(words) - 1 // Add spaces between words
	}
	// Add length of each word
	for _, word := range words {
		totalLength += len(word)
	}

	// Only truncate words if the total length would exceed maxTotal
	if totalLength > maxTotal {
		// Now apply truncation rules on each word
		for i, word := range words {
			// First word gets special treatment if there are more than 2 words
			if i == 0 && len(words) > 2 {
				if len(word) > maxChars {
					result = append(result, word[:min(maxChars, len(word))]+".")
				} else {
					result = append(result, word)
				}
			} else if len(word) > maxChars {
				// Truncate each word to maxChars and add a dot
				result = append(result, word[:min(maxChars, len(word))]+".")
			} else {
				// Keep short words as is
				result = append(result, word)
			}
		}
	} else {
		// No need to truncate words if the total length fits
		result = words
	}

	// Join the words back together with spaces
	truncated := JoinWithSpace(result)

	// Check if the truncated text exceeds maxTotal
	if len(truncated) > maxTotal {
		// Create a temporary Text to use the Truncate method
		temp := &Text{content: truncated}

		// Use the Truncate method, ensuring we leave space for the suffix
		temp.Truncate(maxTotal - len(suffixStr))

		// Set the final result with suffix
		t.content = temp.content + suffixStr
	} else {
		// No need for suffix if it fits
		t.content = truncated
	}

	return t
}
