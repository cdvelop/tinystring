package tinystring

// CountOccurrences checks how many times the string 'search' is present in 'text'
// eg: "hello world" with search "world" will return 1
func CountOccurrences(text, search string) int {
	// If the search string is empty, there can be no matches
	if search == "" {
		return 0
	}

	// Get the length of the search string
	searchLen := len(search)

	// Initialize the match counter
	count := 0

	// Traverse the text and count the number of matches
	for i := 0; i <= len(text)-searchLen; i++ {
		if text[i:i+searchLen] == search {
			count++
		}
	}

	// Return the number of matches found
	return count
}

// Contains checks if the string 'search' is present in 'text'
// Returns true if found, false otherwise
// This matches the behavior of the standard library strings.Contains
func Contains(text, search string) bool {
	return CountOccurrences(text, search) > 0
}
