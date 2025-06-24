package tinystring

// ToBool converts the conv content to a boolean value using internal implementations
// Returns the boolean value and any error that occurred
func (t *conv) ToBool() (bool, error) {
	if t.hasContent(buffErr) {
		return false, t
	}

	// Get string representation using buffer API
	inp := t.ensureStringInOut()

	// Direct boolean string matches
	switch inp {
	case "true", "True", "TRUE", "1", "t", "T":
		t.kind = KBool
		return true, nil
	case "false", "False", "FALSE", "0", "f", "F":
		t.kind = KBool
		return false, nil
	} // Try to parse as integer using internal parseSmallInt
	if intVal, err := t.parseSmallInt(inp); err == nil {
		t.kind = KBool
		return intVal != 0, nil
	}

	// Try basic float patterns (simple cases)
	if inp == "0.0" || inp == "0.00" || inp == "+0" || inp == "-0" {
		t.kind = KBool
		return false, nil
	}
	if inp != "0" && (len(inp) > 0 && (inp[0] >= '1' && inp[0] <= '9')) {
		// Non-zero number starting with digit 1-9, likely true
		t.kind = KBool
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
