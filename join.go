package fmt

// Join concatenates the elements of a string slice to create a single string.
// If no separator is provided, it uses a space as default.
// Can be called with varargs to specify a custom separator.
// eg: Convert([]string{"Hello", "World"}).Join() => "Hello World"
// eg: Convert([]string{"Hello", "World"}).Join("-") => "Hello-World"
func (c *Conv) Join(sep ...string) *Conv {
	separator := " " // default separator is space
	if len(sep) > 0 {
		separator = sep[0]
	}

	// Handle case when we have a string slice stored (DEFERRED CONVERSION)
	if c.kind == K.Slice && c.dataPtr != nil {
		// Use proper unsafe.Pointer to []string reconstruction
		slice := *(*[]string)(c.dataPtr)
		if len(slice) > 0 {
			c.ResetBuffer(BuffOut)
			for i, s := range slice {
				if i > 0 {
					c.AnyToBuff(BuffOut, separator)
				}
				c.AnyToBuff(BuffOut, s)
			}
		}
		return c
	}

	// For other types, convert to string first using AnyToBuff through GetString
	// OPTIMIZED: Check if content is ASCII for fast path
	if c.outLen == 0 {
		return c
	}

	// Check if buffer contains only ASCII for fast processing
	isASCII := true
	for i := 0; i < c.outLen; i++ {
		if c.out[i] > 127 {
			isASCII = false
			break
		}
	}

	if isASCII {
		// Fast path: ASCII-only content - process buffer directly
		var parts [][]byte
		start := 0

		for i := 0; i < c.outLen; i++ {
			ch := c.out[i]
			if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
				if i > start {
					part := make([]byte, i-start)
					copy(part, c.out[start:i])
					parts = append(parts, part)
				}
				start = i + 1
			}
		}
		if start < c.outLen {
			part := make([]byte, c.outLen-start)
			copy(part, c.out[start:])
			parts = append(parts, part)
		}

		// Join parts with the separator using buffer operations
		if len(parts) > 0 {
			c.ResetBuffer(BuffOut) // Reset output buffer
			for i, part := range parts {
				if i > 0 {
					c.AnyToBuff(BuffOut, separator)
				}
				c.wrBytes(BuffOut, part)
			}
		}
	} else {
		// Unicode fallback: use string processing
		str := c.GetString(BuffOut)
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

			// Join parts with the separator using AnyToBuff only
			if len(parts) > 0 {
				c.ResetBuffer(BuffOut) // Reset output buffer
				for i, part := range parts {
					if i > 0 {
						c.AnyToBuff(BuffOut, separator)
					}
					c.AnyToBuff(BuffOut, part)
				}
			}
		}
	}

	return c
}
