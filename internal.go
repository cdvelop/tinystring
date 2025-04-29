package tinystring

// SplitBySpace divides a string by whitespace and returns a slice of substrings
// This is an internal function that provides similar functionality to strings.Fields
// eg: SplitBySpace("Hello World") => []string{"Hello", "World"}
func SplitBySpace(s string) []string {
	var result []string
	inWord := false
	start := 0

	// Iterate through the string character by character
	for i, ch := range s {
		isSpace := ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'

		if !isSpace && !inWord {
			// Start of a new word
			inWord = true
			start = i
		} else if isSpace && inWord {
			// End of a word
			inWord = false
			result = append(result, s[start:i])
		}
	}

	// Handle the last word if the string doesn't end with whitespace
	if inWord {
		result = append(result, s[start:])
	}

	return result
}

// JoinWithSpace concatenates the elements of a string slice to create a single string with the elements
// separated by spaces
// eg: JoinWithSpace([]string{"Hello", "World"}) => "Hello World"
func JoinWithSpace(elements []string) string {
	if len(elements) == 0 {
		return ""
	}

	// Calculate the total length of the result string
	totalLen := len(elements) - 1 // For spaces between elements
	for _, e := range elements {
		totalLen += len(e)
	}

	// Build the result
	var result string
	for i, e := range elements {
		result += e
		if i < len(elements)-1 {
			result += " "
		}
	}

	return result
}
