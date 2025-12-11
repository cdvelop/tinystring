package fmt

// Write appends any value to the buffer using unified type handling
// This is the core builder method that enables fluent chaining
//
// Usage:
//
//	c.Write("hello").Write(" ").Write("world")  // Strings
//	c.Write(42).Write(" items")                 // Numbers
//	c.Write('A').Write(" grade")                // Runes
func (c *Conv) Write(v any) *Conv {
	if c.hasContent(BuffErr) { // Use buffer API
		return c // Error chain interruption
	}

	// BUILDER INTEGRATION: Only transfer initial value if buffer is empty
	// and we have a stored value that hasn't been converted yet
	if c.outLen == 0 && c.dataPtr != nil {
		// Convert current value to buffer using AnyToBuff() - need to reconstruct any
		// For now, skip this optimization until we implement proper unsafe reconstruction
		// TODO: Implement unsafe.Pointer to any reconstruction
	}

	// Use unified AnyToBuff() function to append new value
	c.AnyToBuff(BuffOut, v)
	return c
}

// Reset clears all Conv fields and resets the buffer
// Useful for reusing the same Conv object for multiple operations
func (c *Conv) Reset() *Conv {
	// Reset all Conv fields to default state using buffer API
	c.resetAllBuffers()
	c.dataPtr = nil
	c.kind = K.String
	return c
}

// END OF FILE - setVal() and val2Buf() eliminated per unified buffer architecture
