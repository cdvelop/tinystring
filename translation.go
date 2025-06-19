package tinystring

// T creates a translated string with support for multilingual translations
// Same functionality as Err but returns string directly instead of *conv
// This function is used internally by the builder API for efficient string construction
// 
// Usage examples:
// T(D.Invalid, D.Format) returns "invalid format"
// T(ES, D.Invalid, D.Format) returns "formato inválido"
func T(values ...any) string {
	if len(values) == 0 {
		return ""
	}

	// Use builder pattern for efficient string construction
	c := getConv()
	defer c.putConv()

	// Check if first argument is a language selector
	startIdx := 0
	currentLang := defLang
	
	// Language detection - check if first value is a language
	if len(values) > 0 {
		if l, ok := values[0].(lang); ok {
			currentLang = l
			startIdx = 1
		}
	}

	// Process remaining values and build the translated string
	for i := startIdx; i < len(values); i++ {
		if i > startIdx {
			c.buf = append(c.buf, ' ') // Add space between words
		}
		
		switch v := values[i].(type) {
		case OL:
			// Dictionary term - get translation for current language
			translation := v.get(currentLang)
			c.buf = append(c.buf, translation...)
		case string:
			// Direct string
			c.buf = append(c.buf, v...)
		default:
			// Convert other types to string and append
			c.setVal(v, 0) // Use mode 0 for initial conversion
			str := c.getString()
			c.buf = append(c.buf, str...)
		}
	}

	// Return the constructed string
	return string(c.buf)
}
