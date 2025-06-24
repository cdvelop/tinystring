package tinystring

// Repeat returns the string s repeated n times.
// If n is less than or equal to zero, or if s is empty, it returns an empty string.
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *conv) Repeat(n int) *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}
	if n <= 0 {
		// Clear buffer for empty out and clear ptrValue to prevent reconstruction
		t.rstBuffer(buffOut)
		t.ptrValue = "" // Set to empty string to prevent getBuffString from reconstructing
		return t
	}
	// Phase 4.2: Inline newBuf method to eliminate function call overhead
	str := t.getBuffString()
	if len(str) == 0 {
		// Clear buffer for empty out
		t.rstBuffer(buffOut)
		return t
	}
	// Use buffer API ONLY - no direct buffer manipulation
	t.rstBuffer(buffOut) // Clear output buffer using API

	// Write string n times using buffer API
	for range n {
		t.wrString(buffOut, str) // Use API to write to output buffer
	}

	return t
}
