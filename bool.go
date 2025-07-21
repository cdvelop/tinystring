package tinystring

// Bool converts the conv content to a boolean value using internal implementations
// Returns the boolean value and any error that occurred
func (c *conv) Bool() (bool, error) {
	if c.hasContent(buffErr) {
		return false, c
	}

	// Get string representation using buffer API
	inp := c.getBuffString()

	// Direct boolean string matches
	switch inp {
	case "true", "True", "TRUE", "1", "t", "T":
		c.Kind = K.Bool
		return true, nil
	case "false", "False", "FALSE", "0", "f", "F":
		c.Kind = K.Bool
		return false, nil
	}
	// Try to parse as integer using parseIntString (base 10, signed)
	intVal := c.parseIntString(inp, 10, true)
	if !c.hasContent(buffErr) {
		c.Kind = K.Bool
		return intVal != 0, nil
	} else {
		// Limpia el error generado por el intento fallido usando la API
		c.rstBuffer(buffErr)
	}

	// Try basic float patterns (simple cases)
	if inp == "0.0" || inp == "0.00" || inp == "+0" || inp == "-0" {
		c.Kind = K.Bool
		return false, nil
	}
	if inp != "0" && (len(inp) > 0 && (inp[0] >= '1' && inp[0] <= '9')) {
		// Non-zero number starting with digit 1-9, likely true
		c.Kind = K.Bool
		return true, nil
	}

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
