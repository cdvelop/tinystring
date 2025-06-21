package tinystring

// ToBool converts the conv content to a boolean value
// Returns the boolean value and any error that occurred
func (t *conv) ToBool() (bool, error) {
	if len(t.err) > 0 {
		return false, t
	}

	switch t.kind {
	case KBool:
		return t.boolVal, nil // Direct return for boolean values
	case KInt:
		return t.intVal != 0, nil // Non-zero integers are true
	case KUint:
		return t.uintVal != 0, nil // Non-zero unsigned integers are true
	case KFloat64:
		return t.floatVal != 0.0, nil // Non-zero floats are true
	default:
		// For string types, parse the string content
		inp := t.ensureStringInOut()
		switch inp {
		case "true", "True", "TRUE", "1", "t", "T":
			t.boolVal = true
			t.kind = KBool
			return true, nil
		case "false", "False", "FALSE", "0", "f", "F":
			t.boolVal = false
			t.kind = KBool
			return false, nil
		default:
			// Try to parse as numeric - non-zero numbers are true
			t.stringToInt(10)
			if len(t.err) == 0 {
				t.boolVal = t.intVal != 0
				t.kind = KBool
				t.err = t.err[:0] // Clear any errors since we successfully converted
				return t.boolVal, nil
			}

			// Reset error and try float
			t.err = t.err[:0]
			t.stringToFloat()
			if len(t.err) == 0 {
				t.boolVal = t.floatVal != 0.0
				t.kind = KBool
				t.err = t.err[:0] // Clear any errors since we successfully converted
				return t.boolVal, nil
			}

			return false, Err(D.Boolean, D.Value, D.Invalid, inp)
		}
	}
}
