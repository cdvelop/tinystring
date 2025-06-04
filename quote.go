package tinystring

// Quote wraps a string in double quotes and escapes any special characters
// Example: Quote("hello \"world\"") returns "\"hello \\\"world\\\"\""
func (t *conv) Quote() *conv {
	t.quoteString()
	return t
}

// quoteString quotes a string by wrapping it in double quotes and escaping special characters
// Integrated from tinystrconv QuoteString function - optimized for minimal allocations
func (c *conv) quoteString() {
	input := c.getString()
	if input == "" {
		c.setString("\"\"")
		return
	}

	// Pre-allocate with estimated size (input length + 20% buffer for escapes + 2 for quotes)
	estimatedSize := len(input) + (len(input) / 5) + 2
	result := make([]byte, 0, estimatedSize)

	result = append(result, '"')
	for _, char := range input {
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
	c.setString(string(result))
}
