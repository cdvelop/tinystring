package tinystring

// Repeat returns the string s repeated n times.
// If n is less than or equal to zero, or if s is empty, it returns an empty string.
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *conv) Repeat(n int) *conv {
	str, buf := t.newBuf(n)
	if n <= 0 || isEmpty(str) {
		t.setString("")
		return t
	}

	// Write string n times
	for range n {
		buf = append(buf, str...)
	}

	t.setString(string(buf))
	return t
}
