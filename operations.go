package fmt

// LastIndex returns the index of the last instance of substr in s, or -1 if substr is not present in s.
//
// Special cases:
//   - If substr is empty, LastIndex returns len(s).
//   - If substr is not found in s, LastIndex returns -1.
//   - If substr is longer than s, LastIndex returns -1.
//
// Examples:
//
//	LastIndex("hello world", "world")     // returns 6
//	LastIndex("hello world hello", "hello") // returns 12 (last occurrence)
//	LastIndex("image.backup.jpg", ".")    // returns 12 (useful for file extensions)
//	LastIndex("hello", "xyz")             // returns -1 (not found)
//	LastIndex("hello", "")                // returns 5 (len("hello"))
//	LastIndex("", "hello")                // returns -1 (not found in empty string)
//
// Common use case - extracting file extensions:
//
//	filename := "document.backup.pdf"
//	pos := LastIndex(filename, ".")
//	if pos >= 0 {
//	    extension := filename[pos+1:] // "pdf"
//	}
func LastIndex(s, substr string) int {
	n := len(substr)

	// Handle edge cases
	if n == 0 {
		return len(s)
	}
	if n > len(s) {
		return -1
	}

	// Simple reverse loop
	for i := len(s) - n; i >= 0; i-- {
		if s[i:i+n] == substr {
			return i
		}
	}

	return -1
}
