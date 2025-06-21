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
	if c.hasError() { // Use buffer API
		return c // Error chain interruption
	}
	// BUILDER INTEGRATION: If buffer is empty but we have initial value, transfer it first
	hasStringPointer := false
	if c.kind == KPointer && c.pointerVal != nil {
		if strPtr, ok := c.pointerVal.(*string); ok && *strPtr != "" {
			hasStringPointer = true
		}
	}

	if !c.hasOutContent() && ((c.kind == KString && c.outLen > 0) ||
		hasStringPointer ||
		(c.kind == KSliceStr && len(c.stringSliceVal) > 0) ||
		(c.kind == KInt || c.kind == KUint || c.kind == KFloat64 || c.kind == KBool)) {
		// Convert current value to buffer using anyToBuff()
		anyToBuff(c, buffOut, c.pointerVal) // Use unified conversion
	}

	// Use unified anyToBuff() function instead of complex setVal()
	anyToBuff(c, buffOut, v)
	return c
}

// Reset clears all conv fields and resets the buffer
// Useful for reusing the same conv object for multiple operations
func (c *conv) Reset() *conv {
	// Reset all conv fields to default state using buffer API
	c.rstOut()     // Clear main buffer using API
	c.rstWork()    // Clear temp buffer using API
	c.clearError() // Clear error buffer using API
	c.intVal = 0
	c.uintVal = 0
	c.floatVal = 0
	c.boolVal = false
	c.stringSliceVal = nil
	c.pointerVal = nil
	c.kind = KString
	return c
}

// END OF FILE - setVal() and val2Buf() eliminated per unified buffer architecture
