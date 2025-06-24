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

	// Handle case when we have a string slice stored (LAZY CONVERSION)
	if t.kind == KSliceStr && t.anyValue != nil {
		if slice, ok := t.anyValue.([]string); ok {
			// Direct join using anyToBuff to output buffer
			t.rstBuffer(buffOut)
			for i, s := range slice {
				if i > 0 {
					t.anyToBuff( buffOut, separator)
				}
				t.anyToBuff( buffOut, s)
			}
		}
		return t
	}

	// For other types, convert to string first using anyToBuff through ensureStringInOut
	str := t.ensureStringInOut()
	if str != "" {
		// Split content by whitespace and rejoin with new separator
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

		// Join parts with the separator using anyToBuff only
		if len(parts) > 0 {
			t.rstBuffer(buffOut) // Reset output buffer
			for i, part := range parts {
				if i > 0 {
					t.anyToBuff( buffOut, separator)
				}
				t.anyToBuff( buffOut, part)
			}
		}
	}

	return t
}
