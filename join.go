package tinystring

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
