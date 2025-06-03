package tinystring

// Custom error types to avoid importing standard library packages like "errors" or "fmt"
// This keeps the binary size minimal for embedded systems and WebAssembly

// errorType represents different types of errors that can occur
type errorType uint8

const (
	errNone errorType = iota
	errEmptyString
	errInvalidNumber
	errNegativeUnsigned
	errInvalidBase
	errOverflow
	errInvalidFormat
	errFormatMissingArg
	errFormatWrongType
	errFormatUnsupported
	errIncompleteFormat
	errCannotRound
	errCannotFormat
	errInvalidFloat
)

// tinyError represents a lightweight error without external dependencies
type tinyError struct {
	errType errorType
	context string // minimal context without fmt.Sprintf
}

// Error implements the error interface
func (e *tinyError) Error() string {
	switch e.errType {
	case errEmptyString:
		return "empty string"
	case errInvalidNumber:
		if e.context != "" {
			return "invalid number: " + e.context
		}
		return "invalid number"
	case errNegativeUnsigned:
		return "negative numbers are not supported for unsigned integers"
	case errInvalidBase:
		return "invalid base"
	case errOverflow:
		return "number overflow"
	case errInvalidFormat:
		if e.context != "" {
			return e.context
		}
		return "invalid format"
	case errFormatMissingArg:
		if e.context != "" {
			return "missing argument for " + e.context
		}
		return "missing argument"
	case errFormatWrongType:
		if e.context != "" {
			return "wrong argument type for " + e.context
		}
		return "wrong argument type"
	case errFormatUnsupported:
		if e.context != "" {
			return "unsupported format specifier: " + e.context
		}
		return "unsupported format specifier"
	case errIncompleteFormat:
		return "incomplete format specifier at end of string"
	case errCannotRound:
		return "cannot round non-numeric value"
	case errCannotFormat:
		return "cannot format non-numeric value"
	case errInvalidFloat:
		return "invalid float string"
	default:
		return "unknown error"
	}
}

// newError creates a new tinyError without external dependencies
func newError(errType errorType, context ...string) error {
	err := &tinyError{errType: errType}
	if len(context) > 0 {
		err.context = context[0]
	}
	return err
}

// Error constructors for common cases
func newEmptyStringError() error {
	return newError(errEmptyString)
}

func newInvalidNumberError(input string) error {
	return newError(errInvalidNumber, input)
}

func newFormatMissingArgError(specifier string) error {
	return newError(errFormatMissingArg, specifier)
}

func newFormatWrongTypeError(specifier string) error {
	return newError(errFormatWrongType, specifier)
}

func newFormatUnsupportedError(specifier string) error {
	return newError(errFormatUnsupported, specifier)
}

func newInvalidFloatError() error {
	return newError(errInvalidFloat)
}
