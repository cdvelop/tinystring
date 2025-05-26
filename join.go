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
			// Use string builder for efficient concatenation
			totalLen := len(t.contentSlice) * len(separator)
			for _, s := range t.contentSlice {
				totalLen += len(s)
			}
			builder := newTinyStringBuilder(totalLen)

			// Build the result efficiently
			for i, s := range t.contentSlice {
				builder.writeString(s)
				if i < len(t.contentSlice)-1 {
					builder.writeString(separator)
				}
			}
			t.content = builder.string()
		}
		return t
	}
	// If content is already a string, we split it and join it again with the new separator
	if t.content != "" {
		// Split content by whitespace using string builder for efficiency
		var parts []string
		builder := newTinyStringBuilder(64)

		for _, r := range t.content {
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				if builder.len() > 0 {
					parts = append(parts, builder.string())
					builder.reset()
				}
			} else {
				builder.writeRune(r)
			}
		}
		if builder.len() > 0 {
			parts = append(parts, builder.string())
		}

		// Join parts with the separator using builder
		builder.reset()
		for i, part := range parts {
			builder.writeString(part)
			if i < len(parts)-1 {
				builder.writeString(separator)
			}
		}
		t.content = builder.string()
	}

	return t
}
