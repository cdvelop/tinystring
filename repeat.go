package tinystring

// Repeat returns the string s repeated n times.
// If n is less than or equal to zero, or if s is empty, it returns an empty string.
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *conv) Repeat(n int) *conv {
	if len(t.err) > 0 {
		return t // Error chain interruption
	}
	if n <= 0 {
		// Clear buffer for empty out
		t.out = t.out[:0]
		t.outLen = 0
		return t
	}
	// Phase 4.2: Inline newBuf method to eliminate function call overhead
	str := t.ensureStringInOut()
	if len(str) == 0 {
		// Clear buffer for empty out
		t.out = t.out[:0]
		t.outLen = 0
		return t
	}

	bufSize := len(str) * n
	if bufSize < 16 {
		bufSize = 16 // Minimum useful buffer size
	}
	out := make([]byte, 0, bufSize)

	// Write string n times
	for range n {
		out = append(out, str...)
	}

	// Update buffer instead of using setString for buffer-first strategy
	t.out = append(t.out[:0], out...)
	return t
}
