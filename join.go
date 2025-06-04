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
		// Split content by whitespace using simple string operations
		var parts []string
		runes := []rune(str)
		start := 0

		for i, r := range runes {
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				if i > start {
					parts = append(parts, string(runes[start:i]))
				}
				start = i + 1
			}
		}
		if start < len(runes) {
			parts = append(parts, string(runes[start:]))
		}

		// Join parts with the separator using pre-allocated buffer
		if len(parts) > 0 {
			totalLen := 0
			for _, part := range parts {
				totalLen += len(part)
			}
			totalLen += (len(parts) - 1) * len(separator)

			buf := make([]byte, 0, totalLen)
			for i, part := range parts {
				buf = append(buf, part...)
				if i < len(parts)-1 {
					buf = append(buf, separator...)
				}
			}
			t.setString(string(buf))
		}
	}

	return t
}
