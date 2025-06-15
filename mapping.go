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
	// ASCII case conversion constant
	asciiCaseDiff = 32
	// Buffer capacity constants
	smallBufCap   = 4  // small arrays/words
	mediumBufCap  = 10 // medium text operations
	defaultBufCap = 16 // default buffer size
)

// Helper function to check if error is empty (reduces repeated == "" checks)
// Helper function to check if string is empty (reduces repeated len() calls)
func isEmpty(s string) bool {
	return len(s) == 0
}

// Helper function to check if slice/string has content (reduces repeated len() > 0 checks)
func hasLength(s any) bool {
	switch v := s.(type) {
	case string:
		return len(v) > 0
	case []string:
		return len(v) > 0
	case []int:
		return len(v) > 0
	case []any:
		return len(v) > 0
	default:
		return false
	}
}

// Helper function to create byte buffer with estimated capacity
func makeBuf(cap int) []byte {
	if cap < defaultBufCap {
		cap = defaultBufCap
	}
	return make([]byte, 0, cap)
}

// Helper function to create rune buffer with estimated capacity
func makeRuneBuf(cap int) []rune {
	if cap < defaultBufCap {
		cap = defaultBufCap
	}
	return make([]rune, 0, cap)
}

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
	str, buf := t.newBuf(2)
	hc := false
	for _, r := range str {
		// Find accent and replace with base character using index lookup
		found := false
		// Check lowercase accents
		for i, char := range aL {
			if r == char {
				buf = addRne2Buf(buf, bL[i])
				hc = true
				found = true
				break
			}
		}
		// Check uppercase accents if not found in lowercase
		if !found {
			for i, char := range aU {
				if r == char {
					buf = addRne2Buf(buf, bU[i])
					hc = true
					found = true
					break
				}
			}
		}
		if !found {
			buf = addRne2Buf(buf, r)
		}
	}
	if !hc {
		return t
	}
	t.setString(string(buf))
	return t
}
