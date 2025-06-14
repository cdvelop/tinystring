package tinystring

// ToBool converts the conv content to a boolean value using reflection-only approach
// Returns the boolean value and any error that occurred
func (t *conv) ToBool() (bool, error) {
	if t.err != "" {
		return false, t
	}

	switch t.vTpe {
	case tpBool:
		return t.getBool(), nil // Use reflection to get boolean value
	case tpInt, tpInt8, tpInt16, tpInt32, tpInt64:
		return t.getInt64() != 0, nil // Non-zero integers are true
	case tpUint, tpUint8, tpUint16, tpUint32, tpUint64, tpUintptr:
		return t.getUint64() != 0, nil // Non-zero unsigned integers are true
	case tpFloat32, tpFloat64:
		return t.getFloat64() != 0.0, nil // Non-zero floats are true
	default:
		// For string types, parse the string content
		inp := t.getString()
		switch inp {
		case "true", "True", "TRUE", "1", "t", "T":
			// Set boolean value via reflection
			t.refVal = refValueOf(true)
			t.vTpe = tpBool
			return true, nil
		case "false", "False", "FALSE", "0", "f", "F":
			// Set boolean value via reflection
			t.refVal = refValueOf(false)
			t.vTpe = tpBool
			return false, nil
		default:
			// Try to parse as numeric - non-zero numbers are true
			t.s2Int(10)
			if t.err == "" {
				intVal := t.getInt64()
				boolResult := intVal != 0
				t.refVal = refValueOf(boolResult)
				t.vTpe = tpBool
				t.err = "" // Clear any errors since we successfully converted
				return boolResult, nil
			}

			// Reset error and try float
			t.err = ""
			t.s2Float()
			if t.err == "" {
				floatVal := t.getFloat64()
				boolResult := floatVal != 0.0
				t.refVal = refValueOf(boolResult)
				t.vTpe = tpBool
				t.err = "" // Clear any errors since we successfully converted
				return boolResult, nil
			}

			t.NewErr(errInvalidBool, inp)
			return false, t
		}
	}
}
