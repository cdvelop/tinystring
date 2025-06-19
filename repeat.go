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

	str, buf := t.newBuf(n)
	if len(str) == 0 {
		// Clear both buffer and stringVal for empty result
		t.buf = t.buf[:0]
		t.stringVal = ""
		return t
	}

	// Write string n times
	for range n {
		buf = append(buf, str...)
	}

	// Update buffer instead of using setString for buffer-first strategy
	t.buf = append(t.buf[:0], buf...)
	return t
}
