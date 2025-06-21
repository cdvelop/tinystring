package tinystring

// Shared constants for maximum code reuse and minimal binary size
const (
	// Digit characters for base conversion (supports bases 2-36)
	digs = "0123456789abcdefghijklmnopqrstuvwxyz"
	// Common string constants to avoid allocations for frequently used values
	emptyStr = ""
	trueStr  = "true"
	falseStr = "false"
	zeroStr  = "0"
	oneStr   = "1"
	// Common punctuation
	dotStr      = "."
	spaceStr    = " "
	ellipsisStr = "..."
	quoteStr    = "\"\""
	// Numeric special values - Phase 3G.7: String constant consolidation
	nanStr    = "NaN"
	infStr    = "Inf"
	negInfStr = "-Inf"
	// ASCII case conversion constant
	asciiCaseDiff = 32
	// Buffer capacity constants
	smallBufCap   = 4  // small arrays/words
	mediumBufCap  = 10 // medium text operations
	defaultBufCap = 16 // default buffer size
)

// Index-based character mapping for maximum efficiency
var (
	// Accented characters (lowercase)
	aL = []rune{'á', 'à', 'ã', 'â', 'ä', 'é', 'è', 'ê', 'ë', 'í', 'ì', 'î', 'ï', 'ó', 'ò', 'õ', 'ô', 'ö', 'ú', 'ù', 'û', 'ü', 'ý', 'ñ'}
	// Base characters (lowercase)
	bL = []rune{'a', 'a', 'a', 'a', 'a', 'e', 'e', 'e', 'e', 'i', 'i', 'i', 'i', 'o', 'o', 'o', 'o', 'o', 'u', 'u', 'u', 'u', 'y', 'n'}
	// Accented characters (uppercase)
	aU = []rune{'Á', 'À', 'Ã', 'Â', 'Ä', 'É', 'È', 'Ê', 'Ë', 'Í', 'Ì', 'Î', 'Ï', 'Ó', 'Ò', 'Õ', 'Ô', 'Ö', 'Ú', 'Ù', 'Û', 'Ü', 'Ý', 'Ñ'}
	// Base characters (uppercase)
	bU = []rune{'A', 'A', 'A', 'A', 'A', 'E', 'E', 'E', 'E', 'I', 'I', 'I', 'I', 'O', 'O', 'O', 'O', 'O', 'U', 'U', 'U', 'U', 'Y', 'N'}
)

// toUpperRune converts a single rune to uppercase using optimized lookup
func toUpperRune(r rune) rune {
	// ASCII fast path
	if r >= 'a' && r <= 'z' {
		return r - asciiCaseDiff
	}
	// Accent conversion using index lookup
	for i, char := range aL {
		if r == char {
			return aU[i]
		}
	}
	return r
}

// toLowerRune converts a single rune to lowercase using optimized lookup
func toLowerRune(r rune) rune {
	// ASCII fast path
	if r >= 'A' && r <= 'Z' {
		return r + asciiCaseDiff
	}
	// Accent conversion using index lookup
	for i, char := range aU {
		if r == char {
			return aL[i]
		}
	}
	return r
}

// RemoveTilde removes accents and diacritics using index-based lookup
func (t *conv) RemoveTilde() *conv {
	// Check for error chain interruption
	if len(t.err) > 0 {
		return t
	}

	str := t.ensureStringInOut()
	if len(str) == 0 {
		return t
	}

	// Use buffer-first strategy
	tempBuf := make([]byte, 0, len(str)*2)

	for _, r := range str {
		// Find accent and replace with base character using index lookup
		found := false
		// Check lowercase accents
		for i, char := range aL {
			if r == char {
				tempBuf = addRne2Buf(tempBuf, bL[i])
				found = true
				break
			}
		}
		// Check uppercase accents if not found in lowercase
		if !found {
			for i, char := range aU {
				if r == char {
					tempBuf = addRne2Buf(tempBuf, bU[i])
					found = true
					break
				}
			}
		}
		if !found {
			tempBuf = addRne2Buf(tempBuf, r)
		}
	}

	// Always update the buffer, even if no changes were made
	// This ensures consistency with buffer-first strategy
	t.out = t.out[:0]
	t.out = append(t.out, tempBuf...)
	t.setStringFromBuffer()

	return t
}
