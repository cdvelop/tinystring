package fmt

// Quote wraps a string in double quotes and escapes any special characters
// Example: Quote("hello \"world\"") returns "\"hello \\\"world\\\"\""
func (c *Conv) Quote() *Conv {
	if c.hasContent(BuffErr) {
		return c // Error chain interruption
	}
	if c.outLen == 0 {
		c.ResetBuffer(BuffOut)
		c.WrString(BuffOut, quoteStr)
		return c
	}

	// Use work buffer to build quoted string, then swap to output
	c.ResetBuffer(BuffWork)
	c.wrByte(BuffWork, '"')

	// Process buffer directly without string allocation (like capitalizeASCIIOptimized)
	for i := 0; i < c.outLen; i++ {
		char := c.out[i]
		switch char {
		case '"':
			c.wrByte(BuffWork, '\\')
			c.wrByte(BuffWork, '"')
		case '\\':
			c.wrByte(BuffWork, '\\')
			c.wrByte(BuffWork, '\\')
		case '\n':
			c.wrByte(BuffWork, '\\')
			c.wrByte(BuffWork, 'n')
		case '\r':
			c.wrByte(BuffWork, '\\')
			c.wrByte(BuffWork, 'r')
		case '\t':
			c.wrByte(BuffWork, '\\')
			c.wrByte(BuffWork, 't')
		default:
			c.wrByte(BuffWork, char)
		}
	}

	c.wrByte(BuffWork, '"')
	c.swapBuff(BuffWork, BuffOut)
	return c
}
