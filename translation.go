package tinystring

// Translate creates a translated string with support for multilingual translations
// Same functionality as Err but returns string directly instead of *Conv
// This function is used internally by the builder API for efficient string construction
//
// Usage examples:
// Translate(D.Format, D.Invalid) returns "invalid format"
// Translate(ES, D.Format, D.Invalid) returns "formato inválido"
// Translate creates a translated string with support for multilingual translations
// Same functionality as Err but returns *Conv for further formatting
// This function is used internally by the builder API for efficient string construction
//
// Usage examples:
// Translate(D.Format, D.Invalid) returns *Conv with "invalid format"
// Translate(ES, D.Format, D.Invalid) returns *Conv with "formato inválido"
//
// MEMORY MANAGEMENT:
// The returned *Conv object is pooled.
// - Calling .String() or .Apply() automatically returns it to the pool.
// - If you use .Bytes() or other methods, you MUST call .PutConv() manually to avoid memory leaks.
//
// PERFORMANCE OPTIMIZATION:
// For zero-allocation performance, pass pointers to LocStr instead of values:
//
//	Translate(D.Format)     // 1 alloc/op  (LocStr boxing)
//	Translate(&D.Format)    // 0 allocs/op (pointer fits in interface)
//
// This is useful in hot paths where allocation-free operation is critical.
func Translate(values ...any) *Conv {
	c := GetConv()
	// UNIFIED PROCESSING: Use shared intermediate function
	processTranslatedMessage(c, BuffOut, values...)
	return c
}

// =============================================================================
// FUNCIÓN INTERMEDIA UNIFICADA - REUTILIZADA POR Translate() Y ERR()
// =============================================================================

// processTranslatedMessage procesa argumentos variádicos con traducción y escribe al buffer especificado
// FUNCIÓN UNIFICADA: Reduce duplicación de código entre Translate() y Err()
// Maneja detección de idioma, traducción de LocStr, y escritura al buffer destino
func processTranslatedMessage(c *Conv, dest BuffDest, values ...any) {
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
// UNIFIED FUNCTION: Handles language detection for both Translate() and wrErr()
// Returns: (language, startIndex) where startIndex skips the language argument if present
func detectLanguage(c *Conv, args []any) (lang, int) {
	if len(args) == 0 {
		return getCurrentLang(), 0
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
	return getCurrentLang(), 0
}

// processTranslatedArgs processes arguments with language-aware translation
// UNIFIED FUNCTION: Handles argument processing for both Translate() and wrErr()
// Eliminates code duplication between Translate() and wrErr()
// REFACTORED: Uses WrString instead of direct buffer access
func processTranslatedArgs(c *Conv, dest BuffDest, args []any, currentLang lang, startIndex int) {
	for i := startIndex; i < len(args); i++ {
		arg := args[i]
		switch v := arg.(type) {
		case LocStr:
			c.wrTranslation(v, currentLang, dest)
		case string:
			c.WrString(dest, v)
		default:
			c.AnyToBuff(BuffWork, v)
			if c.hasContent(BuffWork) {
				workResult := c.GetString(BuffWork)
				c.WrString(dest, workResult)
				c.ResetBuffer(BuffWork)
			}
		}

		// Agregar espacio después, excepto si es el último o el siguiente es separador
		if shouldAddSpace(args, i) {
			c.WrString(dest, " ")
		}
	}
}

// shouldAddSpace determina si se debe agregar espacio después del argumento actual
func shouldAddSpace(args []any, currentIndex int) bool {
	// No agregar espacio si es el último argumento
	if currentIndex >= len(args)-1 {
		return false
	}

	// Si el argumento actual termina en newline, espacio, o ciertos separadores específicos, no agregar espacio
	if currentStr, ok := args[currentIndex].(string); ok {
		if len(currentStr) > 0 {
			lastChar := currentStr[len(currentStr)-1]
			// Solo ciertos separadores no necesitan espacio después (como '/')
			if lastChar == '\n' || lastChar == ' ' || lastChar == '/' {
				return false
			}
		}
	}

	// Si el siguiente argumento es un string separador, no agregar espacio
	if nextStr, ok := args[currentIndex+1].(string); ok {
		return !isWordSeparator(nextStr)
	}

	// Para otros tipos (LocStr, etc.) sí agregar espacio
	return true
}

// wrTranslation extracts translation for specific language from LocStr and writes to destination buffer
// REUSES: existing LocStr array indexing logic
// METHOD: Now a Conv method that writes directly to buffer without returning anything
func (c *Conv) wrTranslation(locStr LocStr, currentLang lang, dest BuffDest) {
	// Get translation for current language with fallback
	var translation string
	if int(currentLang) < len(locStr) && locStr[currentLang] != "" {
		translation = locStr[currentLang]
	} else {
		// Fallback to English if translation not available
		translation = locStr[EN]
	}

	// Write directly to destination buffer
	c.WrString(dest, translation)
}
