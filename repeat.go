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

	// Use builder from pool to reduce allocations
	builder := getBuilder()
	defer putBuilder(builder)

	// Pre-grow to exact size needed
	builder.grow(len(str) * n)

	// Write string n times
	for i := 0; i < n; i++ {
		builder.writeString(str)
	}

	t.setString(string(builder.buf))
	return t
}
