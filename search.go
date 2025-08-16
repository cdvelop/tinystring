package tinystring

// Index finds the first occurrence of substr in s, returns -1 if not found.
// This is the base primitive that other functions will reuse.
//
// Examples:
//
//	Index("hello world", "world")  // returns 6
//	Index("hello world", "lo")     // returns 3
//	Index("hello world", "xyz")    // returns -1 (not found)
//	Index("hello world", "")       // returns 0 (empty string)
//	Index("data\x00more", "\x00")  // returns 4 (null byte)
func Index(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0 // Comportamiento estándar: cadena vacía se encuentra en posición 0
	}
	if n == 1 {
		// Optimized single byte search
		for i := 0; i < len(s); i++ {
			if s[i] == substr[0] {
				return i
			}
		}
		return -1
	}

	// Brute force for longer strings
	for i := 0; i <= len(s)-n; i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}

// Count checks how many times the string 'search' is present in 'Conv'.
// Uses Index internally for consistency and maintainability.
//
// Examples:
//
//	Count("abracadabra", "abra")    // returns 2
//	Count("hello world", "l")       // returns 3
//	Count("golang", "go")           // returns 1
//	Count("test", "xyz")            // returns 0 (not found)
//	Count("anything", "")           // returns 0 (empty search)
//	Count("a\x00b\x00c", "\x00")    // returns 2 (null bytes)
func Count(Conv, search string) int {
	if len(search) == 0 {
		return 0
	}

	count := 0
	s := Conv
	for {
		i := Index(s, search)
		if i == -1 {
			break
		}
		count++
		s = s[i+len(search):] // Skip past this match
	}
	return count
}

// Contains checks if the string 'search' is present in 'Conv'.
// Uses Index internally for efficient single-pass detection.
//
// Examples:
//
//	Contains("hello world", "world")  // returns true
//	Contains("hello world", "xyz")    // returns false
//	Contains("", "test")              // returns false (empty string)
//	Contains("test", "")              // returns false (empty search)
//	Contains("data\x00more", "\x00")  // returns true (null byte)
//	Contains("Case", "case")          // returns false (case sensitive)
func Contains(Conv, search string) bool {
	if len(search) == 0 {
		return false // Cadena vacía no se considera contenida
	}
	return Index(Conv, search) != -1
}
