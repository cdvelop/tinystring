package tinystring

// toInt converts any numeric interface value to integer
// Returns the integer value and true if conversion was successful,
// or 0 and false if the input is not a valid numeric type
func toInt(value any) (int, bool) {
	// Handle int types
	if result, ok := convertIntTypes(value); ok {
		return result, true
	}

	// Handle uint types
	if result, ok := convertUintTypes(value); ok {
		return result, true
	}

	// Handle float types
	if result, ok := convertFloatTypes(value); ok {
		return result, true
	}

	return 0, false
}

// convertIntTypes handles all signed integer type conversions
func convertIntTypes(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

// convertUintTypes handles all unsigned integer type conversions
func convertUintTypes(value any) (int, bool) {
	switch v := value.(type) {
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true
	default:
		return 0, false
	}
}

// convertFloatTypes handles all floating point type conversions
func convertFloatTypes(value any) (int, bool) {
	switch v := value.(type) {
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}
