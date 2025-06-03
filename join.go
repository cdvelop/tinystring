package tinystring

// Join concatenates the elements of a string slice to create a single string.
// If no separator is provided, it uses a space as default.
// Can be called with varargs to specify a custom separator.
// eg: Convert([]string{"Hello", "World"}).Join() => "Hello World"
// eg: Convert([]string{"Hello", "World"}).Join("-") => "Hello-World"
func (t *conv) Join(sep ...string) *conv {
	separator := " " // default separator is space
	if len(sep) > 0 {
		separator = sep[0]
	}

	// Handle case when we have a string slice stored
	if t.valType == valTypeStringSlice {
		if len(t.stringSliceVal) == 0 {
			t.setString("")
		} else {
			result := t.joinSlice(separator)
			t.setString(result)
		}
		return t
	}

	// If content is already a string, we split it and join it again with the new separator
	str := t.getString()
	if str != "" {
		// Split content by whitespace using string builder for efficiency
		var parts []string
		builder := newTinyStringBuilder(64)

		for _, r := range str {
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
		t.setString(builder.string())
	}

	return t
}
