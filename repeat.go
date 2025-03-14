package tinystring

// Repeat returns the string s repeated n times.
// If n is less than or equal to zero, or if s is empty, it returns an empty string.
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *Text) Repeat(n int) *Text {
	if n <= 0 || len(t.content) == 0 {
		t.content = ""
		return t
	}

	// Preallocate the necessary memory
	result := make([]byte, len(t.content)*n)

	// Copy the first segment
	copy(result, t.content)

	// Duplicate the segment in the rest of the slice
	for i := len(t.content); i < len(result); i *= 2 {
		copy(result[i:], result[:i])
	}

	t.content = string(result)
	return t
}
