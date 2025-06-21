package tinystring

// Quote wraps a string in double quotes and escapes any special characters
// Example: Quote("hello \"world\"") returns "\"hello \\\"world\\\"\""
func (t *conv) Quote() *conv {
	inp := t.ensureStringInOut()
	if len(inp) == 0 {
		t.setString(quoteStr)
		return t
	}

	// Pre-allocate with estimated size (input length + 20% buffer for escapes + 2 for quotes)
	eSz := len(inp) + (len(inp) / 5) + 2
	// Inline makeBuf logic
	estimatedSize := eSz
	if estimatedSize < defaultBufCap {
		estimatedSize = defaultBufCap
	}
	out := make([]byte, 0, estimatedSize)

	out = append(out, '"')
	for _, char := range inp {
		switch char {
		case '"':
			out = append(out, '\\', '"')
		case '\\':
			out = append(out, '\\', '\\')
		case '\n':
			out = append(out, '\\', 'n')
		case '\r':
			out = append(out, '\\', 'r')
		case '\t':
			out = append(out, '\\', 't')
		default:
			out = append(out, string(char)...)
		}
	}
	out = append(out, '"')
	t.setString(string(out))
	return t
}
