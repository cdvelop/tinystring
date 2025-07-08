package tinystring

// Custom error messages to avoid importing standard library packages like "errors" or "fmt"
// This keeps the binary size minimal for embedded systems and WebAssembly

// Err creates a new error message with support for multilingual translations
// Supports LocStr types for translations and lang types for language specification
// eg:
// tinystring.Err("invalid format") returns "invalid format"
// tinystring.Err(D.Format, D.Invalid) returns "invalid format"
// tinystring.Err(ES,D.Format, D.Invalid) returns "formato inválido"

func Err(msgs ...any) *conv {
	c := getConv() // Always obtain from pool
	c.Kind = KErr
	// UNIFIED PROCESSING: Use same intermediate function as T() but write to buffErr
	processTranslatedMessage(c, buffErr, msgs...)
	return c
}

// Errf creates a new conv instance with error formatting similar to fmt.Errf
// Example: tinystring.Errf("invalid value: %s", value).Error()
func Errf(format string, args ...any) *conv {
	c := getConv() // Always obtain from pool
	c.wrFormat(buffErr, format, args...)
	c.Kind = KErr
	return c
}

// StringErr returns the content of the conv along with any error and auto-releases to pool
func (c *conv) StringErr() (out string, err error) {
	// If there's an error, return empty string and the error object (do NOT release to pool)
	if c.hasContent(buffErr) {
		return "", c
	}

	// Otherwise return the string content and no error (safe to release to pool)
	out = c.getBuffString()
	c.putConv()
	return out, nil
}

// wrErr writes error messages with support for int, string and LocStr
// ENHANCED: Now supports int, string and LocStr parameters
// Used internally by anyToBuff for type error messages
func (c *conv) wrErr(msgs ...any) *conv {
	c.Kind = KErr // Set error Kind first

	// Write messages using default language (no detection needed)
	for i, msg := range msgs {
		if i > 0 {
			// Add space between words
			c.wrString(buffErr, " ")
		}
		// fmt.Printf("wrErr: Processing message part: %v\n", msg) // Depuración

		switch v := msg.(type) {
		case LocStr:
			// Translate LocStr using default language
			c.wrTranslation(v, defLang, buffErr)
		case string:
			// Direct string write
			c.wrString(buffErr, v)
		case int:
			// Convert int to string and write - simple conversion for errors
			if v == 0 {
				c.wrString(buffErr, "0")
			} else {
				// Simple int to string conversion for error messages
				var buf [20]byte // Enough for 64-bit int
				n := len(buf)
				negative := v < 0
				if negative {
					v = -v
				}
				for v > 0 {
					n--
					buf[n] = byte(v%10) + '0'
					v /= 10
				}
				if negative {
					n--
					buf[n] = '-'
				}
				c.wrString(buffErr, string(buf[n:]))
			}
		default:
			// For other types, convert to string representation
			c.wrString(buffErr, "<unsupported>")
		}
	}
	// fmt.Printf("wrErr: Final error buffer content: %q, errLen: %d\n", c.getString(buffErr), c.errLen) // Depuración
	return c
}

func (c *conv) getError() string {
	if !c.hasContent(buffErr) { // ✅ Use API method instead of len(c.err)
		return ""
	}
	return c.getString(buffErr) // ✅ Use API method instead of direct string(c.err)
}

func (c *conv) Error() string {
	return c.getError()
}
