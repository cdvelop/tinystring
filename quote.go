package tinystring

// Quote wraps a string in double quotes and escapes any special characters
// Example: Quote("hello \"world\"") returns "\"hello \\\"world\\\"\""
func (c *conv) Quote() *conv {
	if c.hasContent(buffErr) {
		return c // Error chain interruption
	}
	if c.outLen == 0 {
		c.rstBuffer(buffOut)
		c.wrString(buffOut, quoteStr)
		return c
	}

	// Use work buffer to build quoted string, then swap to output
	c.rstBuffer(buffWork)
	c.wrByte(buffWork, '"')

	// Process buffer directly without string allocation (like capitalizeASCIIOptimized)
	for i := 0; i < c.outLen; i++ {
		char := c.out[i]
		switch char {
		case '"':
			c.wrByte(buffWork, '\\')
			c.wrByte(buffWork, '"')
		case '\\':
			c.wrByte(buffWork, '\\')
			c.wrByte(buffWork, '\\')
		case '\n':
			c.wrByte(buffWork, '\\')
			c.wrByte(buffWork, 'n')
		case '\r':
			c.wrByte(buffWork, '\\')
			c.wrByte(buffWork, 'r')
		case '\t':
			c.wrByte(buffWork, '\\')
			c.wrByte(buffWork, 't')
		default:
			c.wrByte(buffWork, char)
		}
	}

	c.wrByte(buffWork, '"')
	c.swapBuff(buffWork, buffOut)
	return c
}
