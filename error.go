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
	if len(t.err) > 0 {
		// If there's an error, return empty string and the error
		out = ""
		err = &customError{message: string(t.err)}
	} else {
		out = t.getString()
		err = nil
	}

	// Auto-release back to pool for memory efficiency
	t.putConv()
	return out, err
}

// Phase 13.3: Helper methods for dynamic buffer management
func (c *conv) addToErrBuf(s string) {
	// Añadir al buffer dinámico de errores
	c.err = append(c.err, s...)
}

// wrErr - método privado para migración de asignaciones de error
// eg: c.wrErr(D.String, D.Empty) // Setear error de cadena vacía
func (c *conv) wrErr(values ...any) *conv {
	c.kind = KErr           // Setear ANTES de llamar T()
	T(append(values, c)...) // T() escribirá directamente al err
	return c
}

func (c *conv) getError() string {
	if len(c.err) == 0 {
		return ""
	}
	return string(c.err)
}

func (c *conv) Error() string {
	return c.getError()
}
