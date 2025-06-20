package tinystring

// Repeat returns the string s repeated n times.
// If n is less than or equal to zero, or if s is empty, it returns an empty string.
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *conv) Repeat(n int) *conv {
	if t.err != "" {
		return t // Error chain interruption
	}

	if n <= 0 {
		// Clear both buffer and stringVal for empty result
		t.buf = t.buf[:0]
		t.stringVal = ""
		return t
	}
	// Phase 4.2: Inline newBuf method to eliminate function call overhead
	str := t.getString()
	if len(str) == 0 {
		// Clear both buffer and stringVal for empty result
		t.buf = t.buf[:0]
		t.stringVal = ""
		return t
	}

	bufSize := len(str) * n
	if bufSize < 16 {
		bufSize = 16 // Minimum useful buffer size
	}
	buf := make([]byte, 0, bufSize)

	// Write string n times
	for range n {
		buf = append(buf, str...)
	}

	// Update buffer instead of using setString for buffer-first strategy
	t.buf = append(t.buf[:0], buf...)
	return t
}
