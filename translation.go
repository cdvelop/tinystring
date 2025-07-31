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

	// UNIFIED PROCESSING: Use shared intermediate function
	processTranslatedMessage(c, buffOut, values...)
	// Return the constructed string
	return c.getString(buffOut)
}

// =============================================================================
// FUNCIÓN INTERMEDIA UNIFICADA - REUTILIZADA POR T() Y ERR()
// =============================================================================

// processTranslatedMessage procesa argumentos variádicos con traducción y escribe al buffer especificado
// FUNCIÓN UNIFICADA: Reduce duplicación de código entre T() y Err()
// Maneja detección de idioma, traducción de LocStr, y escritura al buffer destino
func processTranslatedMessage(c *conv, dest buffDest, values ...any) {
	if len(values) == 0 {
		return
	}

	// PASO 1: Detección unificada de idioma
	currentLang, startIdx := detectLanguage(c, values)

	// PASO 2: Procesamiento unificado de argumentos
	processTranslatedArgs(c, dest, values, currentLang, startIdx)
}

// =============================================================================
// SHARED LANGUAGE SYSTEM FUNCTIONS - REUSED BY ERROR.GO AND TRANSLATION.GO
// =============================================================================

// detectLanguage determines the current language and start index from variadic arguments
// UNIFIED FUNCTION: Handles language detection for both T() and wrErr()
// Returns: (language, startIndex) where startIndex skips the language argument if present
func detectLanguage(c *conv, args []any) (lang, int) {
	if len(args) == 0 {
		return defLang, 0
	}

	// Check if first argument is a language specifier
	if langVal, ok := args[0].(lang); ok {
		return langVal, 1 // Skip the language argument in processing
	}

	// If first argument is a string of length 2, treat as language code
	if strVal, ok := args[0].(string); ok && len(strVal) == 2 {

		return c.mapLangCode(strVal), 1 // Skip the language argument in processing
	}

	// No language specified, use default
	return defLang, 0
}

// processTranslatedArgs processes arguments with language-aware translation
// UNIFIED FUNCTION: Handles argument processing for both T() and wrErr()
// Eliminates code duplication between T() and wrErr()
// REFACTORED: Uses wrString instead of direct buffer access
func processTranslatedArgs(c *conv, dest buffDest, args []any, currentLang lang, startIndex int) {
	for i := startIndex; i < len(args); i++ {
		if i > startIndex {
			// Add space between words for readability using wrString
			c.wrString(dest, " ")
		}

		arg := args[i]
		switch v := arg.(type) {
		case LocStr:
			// Translate LocStr using specified/detected language
			c.wrTranslation(v, currentLang, dest)

		case string:
			// Direct string - write using wrString
			c.wrString(dest, v)

		default:
			// Convert other types to string using anyToBuff
			c.anyToBuff(buffWork, v) // Convert to work buffer first
			if c.hasContent(buffWork) {
				workResult := c.getString(buffWork)
				c.wrString(dest, workResult)
				c.rstBuffer(buffWork) // Clear work buffer for next use
			}
		}
	}
}

// wrTranslation extracts translation for specific language from LocStr and writes to destination buffer
// REUSES: existing LocStr array indexing logic
// METHOD: Now a conv method that writes directly to buffer without returning anything
func (c *conv) wrTranslation(locStr LocStr, currentLang lang, dest buffDest) {
	// Get translation for current language with fallback
	var translation string
	if int(currentLang) < len(locStr) && locStr[currentLang] != "" {
		translation = locStr[currentLang]
	} else {
		// Fallback to English if translation not available
		translation = locStr[EN]
	}

	// Write directly to destination buffer
	c.wrString(dest, translation)
}
