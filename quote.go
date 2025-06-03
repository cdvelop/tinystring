package tinystring

// Quote wraps a string in double quotes and escapes any special characters
// Example: Quote("hello \"world\"") returns "\"hello \\\"world\\\"\""
func (t *conv) Quote() *conv {
	t.quoteString()
	return t
}

// quoteString quotes a string by wrapping it in double quotes and escaping special characters
// Integrated from tinystrconv QuoteString function
func (c *conv) quoteString() {
	input := c.getString()
	if input == "" {
		c.setString("\"\"")
		return
	}

	result := "\""
	for _, char := range input {
		switch char {
		case '"':
			result += "\\\""
		case '\\':
			result += "\\\\"
		case '\n':
			result += "\\n"
		case '\r':
			result += "\\r"
		case '\t':
			result += "\\t"
		default:
			result += string(char)
		}
	}
	result += "\""
	c.setString(result)
}
