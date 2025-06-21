package tinystring

// Custom error messages to avoid importing standard library packages like "errors" or "fmt"
// This keeps the binary size minimal for embedded systems and WebAssembly

// Err creates a new error message with support for multilingual translations
// Supports LocStr types for translations and lang types for language specification
// eg:
// tinystring.Err("invalid format") returns "invalid format"
// tinystring.Err(D.Format, D.Invalid) returns "invalid format"
// tinystring.Err(ES,D.Format, D.Invalid) returns "formato inválido"

func Err(values ...any) *conv {
	c := getConv() // Always obtain from pool
	return c.wrErr(values...)
}

// Errf creates a new conv instance with error formatting similar to fmt.Errf
// Example: tinystring.Errf("invalid value: %s", value).Error()
func Errf(format string, args ...any) *conv {
	c := getConv() // Always obtain from pool
	c.sprintf(format, args...)
	c.kind = KErr
	return c
}

// StringError returns the content of the conv along with any error and auto-releases to pool
func (t *conv) StringError() (string, error) {
	var out string
	var err error
	// BUILDER INTEGRATION: Check for error condition more comprehensively
	if t.hasError() { // ✅ Use new buffer state checking method
		// If there's an error, return empty string and the error
		out = ""
		err = &simpleError{message: string(t.err[:t.errLen])} // ✅ Use errLen for length control
	} else {
		out = t.ensureStringInOut() // ✅ Use new centralized method
		err = nil
	}

	// Auto-release back to pool for memory efficiency
	t.putConv()
	return out, err
}

// simpleError implements error interface without importing errors package
type simpleError struct {
	message string
}

func (e *simpleError) Error() string {
	return e.message
}

// =============================================================================
// LANGUAGE-AWARE ERROR SYSTEM - ARCHITECTURAL IMPLEMENTATION
// NO T() DEPENDENCY - DIRECT BUFFER WRITING - NO ERROR RETURN
// =============================================================================

// wrErr writes error messages using language detection and LocStr translation
// ARCHITECTURAL SPECIFICATION: NO error return, direct buffer writing
// REUSES: detectLanguage, getTranslation, writeStringToErr
func (c *conv) wrErr(msgs ...any) *conv {
	if len(msgs) == 0 {
		return c
	}

	c.kind = KErr // Set error kind first

	// Reset error buffer using API ONLY
	c.clearError() // ✅ Use API method instead of manual c.errLen = 0

	// STEP 1: Language detection (no c.language field dependency)
	currentLang := detectLanguage(c)

	// STEP 2: Process each message argument
	for i, msg := range msgs {
		if i > 0 {
			// Add space between words for readability
			c.writeStringToErr(" ") // ✅ Use API
		}

		switch v := msg.(type) {
		case LocStr:
			// STEP 3: Translate LocStr using detected language
			translation := getTranslation(v, currentLang)
			c.writeStringToErr(translation) // ✅ Use API

		case string:
			// Direct string - write as-is
			c.writeStringToErr(v) // ✅ Use API

		default:
			// Convert other types to string (int, float, etc.)
			// Use anyToBuff to convert to work buffer, then transfer result
			anyToBuff(c, buffWork, v) // Convert to work buffer of CURRENT conv

			if c.hasWorkContent() { // ✅ Use API method
				// Transfer work result to error buffer using API
				c.writeStringToErr(c.getWorkString()) // ✅ Use API
			}
		}
	}
	return c
}

func (c *conv) getError() string {
	if !c.hasError() { // ✅ Use API method instead of len(c.err)
		return ""
	}
	return c.getErrorString() // ✅ Use API method instead of direct string(c.err)
}

func (c *conv) Error() string {
	return c.getError()
}
