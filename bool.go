package tinystring

// ToBool converts the conv content to a boolean value
// Returns the boolean value and any error that occurred
func (t *conv) ToBool() (bool, error) {
	if t.err != "" {
		return false, t
	}

	switch t.vTpe {
	case typeBool:
		return t.boolVal, nil // Direct return for boolean values
	case typeInt:
		return t.intVal != 0, nil // Non-zero integers are true
	case typeUint:
		return t.uintVal != 0, nil // Non-zero unsigned integers are true
	case typeFloat:
		return t.floatVal != 0.0, nil // Non-zero floats are true
	default:
		// For string types, parse the string content
		t.stringToBool()
		if t.err != "" {
			return false, t
		}
		return t.boolVal, nil
	}
}

// stringToBool converts string to bool and stores in conv struct.
// This is an internal conv method that modifies the struct instead of returning values.
func (t *conv) stringToBool() {
	inp := t.getString()
	switch inp {
	case "true", "True", "TRUE", "1", "t", "T":
		t.boolVal = true
		t.vTpe = typeBool
		return
	case "false", "False", "FALSE", "0", "f", "F":
		t.boolVal = false
		t.vTpe = typeBool
		return
	default:
		// Try to parse as numeric - non-zero numbers are true
		t.s2Int(10)
		if t.err == "" {
			t.boolVal = t.intVal != 0
			t.vTpe = typeBool
			t.err = "" // Clear any errors since we successfully converted
			return
		}

		// Reset error and try float
		t.err = ""
		t.s2Float()
		if t.err == "" {
			t.boolVal = t.floatVal != 0.0
			t.vTpe = typeBool
			t.err = "" // Clear any errors since we successfully converted
			return
		}

		t.NewErr(errInvalidBool, inp)
	}
}
