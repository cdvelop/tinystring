package tinystring

// Write appends any value to the buffer using unified type handling
// This is the core builder method that enables fluent chaining
//
// Usage:
//
//	c.Write("hello").Write(" ").Write("world")  // Strings
//	c.Write(42).Write(" items")                 // Numbers
//	c.Write('A').Write(" grade")                // Runes
func (c *conv) Write(v any) *conv {
	if c.hasContent(buffErr) { // Use buffer API
		return c // Error chain interruption
	}

	// BUILDER INTEGRATION: Only transfer initial value if buffer is empty
	// and we have a stored value that hasn't been converted yet
	if c.outLen == 0 && c.ptrValue != nil {
		// Convert current value to buffer using anyToBuff()
		c.anyToBuff(buffOut, c.ptrValue) // Use unified conversion
	}

	// Use unified anyToBuff() function to append new value
	c.anyToBuff(buffOut, v)
	return c
}

// Reset clears all conv fields and resets the buffer
// Useful for reusing the same conv object for multiple operations
func (c *conv) Reset() *conv {
	// Reset all conv fields to default state using buffer API
	c.resetAllBuffers()
	c.ptrValue = nil
	c.kind = KString
	return c
}

// END OF FILE - setVal() and val2Buf() eliminated per unified buffer architecture
