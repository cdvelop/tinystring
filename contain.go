package tinystring

// Count checks how many times the string 'search' is present in 'conv'
// eg: "hello world" with search "world" will return 1
func Count(conv, search string) int {
	// If the search string is empty, there can be no matches
	if len(search) == 0 {
		return 0
	}

	// Get the length of the search string
	searchLen := len(search)

	// Initialize the match counter
	count := 0

	// Traverse the conv and count the number of matches
	for i := 0; i <= len(conv)-searchLen; i++ {
		if conv[i:i+searchLen] == search {
			count++
		}
	}

	// Return the number of matches found
	return count
}

// Contains checks if the string 'search' is present in 'conv'
// Returns true if found, false otherwise
// This matches the behavior of the standard library strings.Contains
func Contains(conv, search string) bool {
	return Count(conv, search) != 0
}
