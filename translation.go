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
			// REUSE getTranslation() function
			translation := getTranslation(v, currentLang)
			c.out = append(c.out, translation...)
		case string:
			// Direct string
			c.out = append(c.out, v...)
		default:
			// Convert other types to string using anyToBuff()
			anyToBuff(c, buffOut, v) // Use unified conversion function
			str := c.ensureStringInOut()
			c.out = append(c.out, str...)
		}
	}
	// Return the constructed string
	return string(c.out)
}

// =============================================================================
// SHARED LANGUAGE SYSTEM FUNCTIONS - REUSED BY ERROR.GO AND TRANSLATION.GO
// =============================================================================

// detectLanguage determines the current language for translations and errors
// REUSES: existing defLang from language.go
func detectLanguage(c *conv) lang {
	// STEP 1: Use default language (can be extended later for auto-detection)
	// Note: c.language field doesn't exist in current struct, use global default
	return defLang // REUSE existing global language setting
}

// getTranslation extracts translation for specific language from LocStr
// REUSES: existing LocStr array indexing logic
func getTranslation(locStr LocStr, currentLang lang) string {
	// STEP 2: Get translation for current language with fallback
	// REUSE logic from T() function
	if int(currentLang) < len(locStr) && locStr[currentLang] != "" {
		return locStr[currentLang]
	}
	// Fallback to English if translation not available
	return locStr[EN]
}
