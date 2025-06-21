package tinystring

// T creates a translated string with support for multilingual translations
// Same functionality as Err but returns string directly instead of *conv
// This function is used internally by the builder API for efficient string construction
//
// Usage examples:
// T(D.Format, D.Invalid) returns "invalid format"
// T(ES, D.Format, D.Invalid) returns "formato inválido"
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
			c.out = append(c.out, ' ') // Add space between words
		}

		switch v := values[i].(type) {
		case LocStr:
			// Dictionary term - get translation for current language
			// Inline get logic
			translation := func() string {
				if int(currentLang) < len(v) && v[currentLang] != "" {
					return v[currentLang]
				}
				return v[EN] // Fallback to English
			}()
			c.out = append(c.out, translation...)
		case string:
			// Direct string
			c.out = append(c.out, v...)
		default:
			// Convert other types to string and append
			c.setVal(v, 0) // Use mode 0 for initial conversion
			str := c.ensureStringInOut()
			c.out = append(c.out, str...)
		}
	}

	// Return the constructed string
	return string(c.out)
}
