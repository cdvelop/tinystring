package tinystring

// Custom error messages to avoid importing standard library packages like "errors" or "fmt"
// This keeps the binary size minimal for embedded systems and WebAssembly

// Err creates a new error message with support for multilingual translations
// REFACTORED: Uses T function and pool for optimal performance
// Supports OL types for translations and lang types for language specification
// Maintains backward compatibility with existing string-based errors
// eg:
// tinystring.Err("invalid format") returns "invalid format"
// tinystring.Err(D.Invalid, D.Fmt) returns "invalid format"
// tinystring.Err(ES,D.Fmt, D.Invalid) returns "formato inválido"
func Err(values ...any) *conv {
	c := getConv() // Always obtain from pool
	c.err = T(values...)
	c.vTpe = typeErr
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
