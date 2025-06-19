package tinystring

// Quote wraps a string in double quotes and escapes any special characters
// Example: Quote("hello \"world\"") returns "\"hello \\\"world\\\"\""
func (t *conv) Quote() *conv {
	inp := t.getString()
	if isEmptySt(inp) {
		t.setString(quoteStr)
		return t
	}

	// Pre-allocate with estimated size (input length + 20% buffer for escapes + 2 for quotes)
	eSz := len(inp) + (len(inp) / 5) + 2
	result := makeBuf(eSz)

	result = append(result, '"')
	for _, char := range inp {
		switch char {
		case '"':
			result = append(result, '\\', '"')
		case '\\':
			result = append(result, '\\', '\\')
		case '\n':
			result = append(result, '\\', 'n')
		case '\r':
			result = append(result, '\\', 'r')
		case '\t':
			result = append(result, '\\', 't')
		default:
			result = append(result, string(char)...)
		}
	}
	result = append(result, '"')
	t.setString(string(result))
	return t
}
