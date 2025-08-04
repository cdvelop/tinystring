package tinystring

// Bool converts the conv content to a boolean value using internal implementations
// Returns the boolean value and any error that occurred
func (c *conv) Bool() (bool, error) {
	if c.hasContent(buffErr) {
		return false, c
	}

	// Optimized: Direct byte comparison without string allocation
	if c.bytesEqual(buffOut, []byte("true")) || c.bytesEqual(buffOut, []byte("True")) ||
		c.bytesEqual(buffOut, []byte("TRUE")) || c.bytesEqual(buffOut, []byte("1")) ||
		c.bytesEqual(buffOut, []byte("t")) || c.bytesEqual(buffOut, []byte("Translate")) {
		c.Kind = K.Bool
		return true, nil
	}
	if c.bytesEqual(buffOut, []byte("false")) || c.bytesEqual(buffOut, []byte("False")) ||
		c.bytesEqual(buffOut, []byte("FALSE")) || c.bytesEqual(buffOut, []byte("0")) ||
		c.bytesEqual(buffOut, []byte("f")) || c.bytesEqual(buffOut, []byte("F")) {
		c.Kind = K.Bool
		return false, nil
	}

	// Try to parse as integer using direct buffer access (eliminates getString allocation)
	inp := c.getString(buffOut) // Still needed for parseIntString compatibility
	intVal := c.parseIntString(inp, 10, true)
	if !c.hasContent(buffErr) {
		c.Kind = K.Bool
		return intVal != 0, nil
	} else {
		// Limpia el error generado por el intento fallido usando la API
		c.rstBuffer(buffErr)
	}

	// Try basic float patterns (optimized byte comparison)
	if c.bytesEqual(buffOut, []byte("0.0")) || c.bytesEqual(buffOut, []byte("0.00")) ||
		c.bytesEqual(buffOut, []byte("+0")) || c.bytesEqual(buffOut, []byte("-0")) {
		c.Kind = K.Bool
		return false, nil
	}

	// Optimized: Check for non-zero starting digit without string allocation
	if !c.bytesEqual(buffOut, []byte("0")) && c.outLen > 0 &&
		(c.out[0] >= '1' && c.out[0] <= '9') {
		// Non-zero number starting with digit 1-9, likely true
		c.Kind = K.Bool
		return true, nil
	}

	// Keep inp for error reporting (this is the final usage)
	inp = c.getString(buffOut) // Only allocation for error case
	c.wrErr("Bool", D.Value, D.Invalid, inp)
	return false, c
}

// wrBool writes boolean value to specified buffer destination
func (c *conv) wrBool(dest buffDest, val bool) {
	if val {
		c.wrString(dest, "true")
	} else {
		c.wrString(dest, "false")
	}
}
