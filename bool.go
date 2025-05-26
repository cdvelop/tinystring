package tinystring

import "errors"

// ToBool converts the text content to a boolean value
// Returns the boolean value and any error that occurred
func (t *Text) ToBool() (bool, error) {
	if t.err != nil {
		return false, t.err
	}
	
	// Since Convert() stores the original value as string in content,
	// we need to handle the original input type that was passed to Convert()
	// For this, we'll parse the content and try to infer the type
	return stringToBool(t.content)
}

// FromBool creates a new Text instance from a boolean value
func FromBool(value bool) *Text {
	return &Text{content: boolToString(value)}
}

// boolToString converts a boolean to its string representation (integrated from tinystrconv)
func boolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

// stringToBool converts a string to a boolean (integrated from tinystrconv)
func stringToBool(input string) (bool, error) {
	switch input {
	case "true", "True", "TRUE", "1", "t", "T":
		return true, nil
	case "false", "False", "FALSE", "0", "f", "F":
		return false, nil
	default:
		// Try to parse as numeric - non-zero numbers are true
		if val, err := stringToInt(input, 10); err == nil {
			return val != 0, nil
		}
		if val, err := stringToFloat(input); err == nil {
			return val != 0.0, nil
		}
		return false, errors.New("invalid boolean string: " + input)
	}
}
