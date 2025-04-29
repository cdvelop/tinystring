package tinystring

// Join concatenates the elements of a string slice to create a single string.
// If no separator is provided, it uses a space as default.
// Can be called with varargs to specify a custom separator.
// eg: Convert([]string{"Hello", "World"}).Join() => "Hello World"
// eg: Convert([]string{"Hello", "World"}).Join("-") => "Hello-World"
func (t *Text) Join(sep ...string) *Text {
	separator := " " // default separator is space
	if len(sep) > 0 {
		separator = sep[0]
	}

	// Handle case when we've received a string slice directly
	if t.contentSlice != nil {
		if len(t.contentSlice) == 0 {
			t.content = ""
		} else {
			// Manually join the elements with the separator
			var totalLen int
			for _, s := range t.contentSlice {
				totalLen += len(s)
			}
			// Add length for separators between elements
			if len(t.contentSlice) > 1 {
				totalLen += len(separator) * (len(t.contentSlice) - 1)
			}

			// Build the result
			var result string
			for i, s := range t.contentSlice {
				result += s
				if i < len(t.contentSlice)-1 {
					result += separator
				}
			}
			t.content = result
		}
		return t
	}

	// If content is already a string, we split it and join it again with the new separator
	if t.content != "" {
		// Split content by whitespace
		var parts []string
		var currentWord string
		for _, r := range t.content {
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				if currentWord != "" {
					parts = append(parts, currentWord)
					currentWord = ""
				}
			} else {
				currentWord += string(r)
			}
		}
		if currentWord != "" {
			parts = append(parts, currentWord)
		}

		// Join parts with the separator
		var result string
		for i, part := range parts {
			result += part
			if i < len(parts)-1 {
				result += separator
			}
		}
		t.content = result
	}

	return t
}
