package tinystring

// Custom error messages to avoid importing standard library packages like "errors" or "fmt"
// This keeps the binary size minimal for embedded systems and WebAssembly

// errorType represents an error as a string constant
type errorType string

// Error message constants
const (
	errNone              errorType = ""
	errEmptyString       errorType = "empty string"
	errInvalidNumber     errorType = "invalid number"
	errNegativeUnsigned  errorType = "negative numbers are not supported for unsigned integers"
	errInvalidBase       errorType = "invalid base"
	errOverflow          errorType = "number overflow"
	errInvalidFormat     errorType = "invalid format"
	errFormatMissingArg  errorType = "missing argument"
	errFormatWrongType   errorType = "wrong argument type"
	errFormatUnsupported errorType = "unsupported format specifier"
	errIncompleteFormat  errorType = "incomplete format specifier at end of string"
	errCannotRound       errorType = "cannot round non-numeric value"
	errCannotFormat      errorType = "cannot format non-numeric value"
	errInvalidFloat      errorType = "invalid float string"
	errInvalidBool       errorType = "invalid boolean value"
)

// set new error message eg: tinystring.Err(errInvalidFormat, "custom message")
func (c *conv) NewErr(values ...any) *conv {
	var sep, out errorType
	c.cachedString = ""
	c.err = ""
	for _, v := range values {
		c.anyToStringInternal(v)
		if c.err != "" {
			out += sep + c.err
		}

		if c.cachedString != "" {
			out += sep + errorType(c.cachedString)
		}

		if c.err != "" || c.cachedString != "" {
			sep = " "
		}
	}
	c.err = out
	c.valType = valTypeErr
	return c
}

// Errorf creates a new conv instance with error formatting similar to fmt.Errorf
// Example: tinystring.Errorf("invalid value: %s", value).Error()
func Errorf(format string, args ...any) *conv {
	result := convInit(new(errorType))
	result.sprintf(format, args...)
	return result
}

// Err if NewErr replace lib errors.New
// supports multiple arguments
// eg: tinystring.Err(errInvalidFormat, "custom message")
// or tinystring.Err(errInvalidFormat, "custom message", "another message")
// or tinystring.Err("custom message", "another message")
func Err(args ...any) *conv {
	c := &conv{}
	return c.NewErr(args...)
}

// Error implements the error interface for conv
// Returns the error message stored in err
func (c *conv) Error() string {
	return string(c.err)
}
