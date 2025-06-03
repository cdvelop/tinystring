package tinystring

// Repeat returns the string s repeated n times.
// If n is less than or equal to zero, or if s is empty, it returns an empty string.
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *conv) Repeat(n int) *conv {
	str := t.getString()
	if n <= 0 || len(str) == 0 {
		t.setString("")
		return t
	}

	// Preallocate the necessary memory
	result := make([]byte, len(str)*n)

	// Copy the first segment
	copy(result, str)

	// Duplicate the segment in the rest of the slice
	for i := len(str); i < len(result); i *= 2 {
		copy(result[i:], result[:i])
	}

	t.setString(string(result))
	return t
}
