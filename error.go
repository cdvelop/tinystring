package tinystring

// Custom error messages to avoid importing standard library packages like "errors" or "fmt"
// This keeps the binary size minimal for embedded systems and WebAssembly

// Err creates a new error message with support for multilingual translations
// Supports OL types for translations and lang types for language specification
// Maintains backward compatibility with existing string-based errors
// eg:
// tinystring.Err("invalid format") returns "invalid format"
// tinystring.Err(D.Invalid, D.Fmt) returns "invalid format"
// tinystring.Err(ES,D.Fmt, D.Invalid) returns "formato inválido"
func Err(values ...any) *conv {
	c := &conv{
		vTpe: typeErr,
		err:  "",
	}

	var sep, out string
	c.tmpStr = ""

	// Determine target language
	targetLang := defLang

	for _, v := range values {
		switch val := v.(type) {
		case lang:
			// Language specified inline
			targetLang = val
			continue
		case OL:
			// Translation type detected - use translation
			out += sep + val.get(targetLang)
		default:
			// Handle other types normally (unchanged logic)
			c.any2s(v)
			if c.err != "" {
				out += sep + c.err
			}
			if c.tmpStr != "" {
				out += sep + c.tmpStr
			}
		}

		if c.err != "" || c.tmpStr != "" || (v != nil) {
			sep = " "
		}
	}

	c.err = out
	return c
}

// Errf creates a new conv instance with error formatting similar to fmt.Errf
// Example: tinystring.Errf("invalid value: %s", value).Error()
func Errf(format string, args ...any) *conv {
	result := unifiedFormat(format, args...)
	result.vTpe = typeErr
	return result
}

// Error implements the error interface for conv
// Returns the error message stored in err
func (c *conv) Error() string {
	return c.err
}
