package tinystring

// Repeat returns the string s repeated n times.
// If n is less than or equal to zero, or if s is empty, it returns an empty string.
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *conv) Repeat(n int) *conv {
	if t.hasError() {
		return t // Error chain interruption
	}
	if n <= 0 {
		// Clear buffer for empty out
		t.rstOut()
		return t
	}
	// Phase 4.2: Inline newBuf method to eliminate function call overhead
	str := t.ensureStringInOut()
	if len(str) == 0 {
		// Clear buffer for empty out
		t.rstOut()
		return t
	}
	// Use buffer API ONLY - no direct buffer manipulation
	t.rstOut() // Clear output buffer using API

	// Write string n times using buffer API
	for range n {
		t.wrStringToOut(str) // Use API to write to output buffer
	}

	return t
}
