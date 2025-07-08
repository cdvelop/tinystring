package tinystring

// Bool converts the conv content to a boolean value using internal implementations
// Returns the boolean value and any error that occurred
func (t *conv) Bool() (bool, error) {
	if t.hasContent(buffErr) {
		return false, t
	}

	// Get string representation using buffer API
	inp := t.getBuffString()

	// Direct boolean string matches
	switch inp {
	case "true", "True", "TRUE", "1", "t", "T":
		t.Kind = KBool
		return true, nil
	case "false", "False", "FALSE", "0", "f", "F":
		t.Kind = KBool
		return false, nil
	}
	// Try to parse as integer using parseIntString (base 10, signed)
	intVal := t.parseIntString(inp, 10, true)
	if !t.hasContent(buffErr) {
		t.Kind = KBool
		return intVal != 0, nil
	} else {
		// Limpia el error generado por el intento fallido usando la API
		t.rstBuffer(buffErr)
	}

	// Try basic float patterns (simple cases)
	if inp == "0.0" || inp == "0.00" || inp == "+0" || inp == "-0" {
		t.Kind = KBool
		return false, nil
	}
	if inp != "0" && (len(inp) > 0 && (inp[0] >= '1' && inp[0] <= '9')) {
		// Non-zero number starting with digit 1-9, likely true
		t.Kind = KBool
		return true, nil
	}

	t.wrErr(D.Boolean, D.Value, D.Invalid, inp)
	return false, t
}

// wrBool writes boolean value to specified buffer destination
func (c *conv) wrBool(dest buffDest, val bool) {
	if val {
		c.wrString(dest, "true")
	} else {
		c.wrString(dest, "false")
	}
}
