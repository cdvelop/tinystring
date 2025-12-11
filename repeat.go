package fmt

// Repeat repeats the Conv content n times
// If n is 0 or negative, it clears the Conv content
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *Conv) Repeat(n int) *Conv {
	if t.hasContent(BuffErr) {
		return t // Error chain interruption
	}
	if n <= 0 {
		// Clear buffer for empty out and clear dataPtr to prevent reconstruction
		t.ResetBuffer(BuffOut)
		t.dataPtr = nil // Clear pointer to prevent GetString from reconstructing
		return t
	}

	// OPTIMIZED: Direct length check
	if t.outLen == 0 {
		// Clear buffer for empty out
		t.ResetBuffer(BuffOut)
		return t
	}

	// OPTIMIZED: Use buffer copy for efficiency
	originalLen := t.outLen
	originalData := make([]byte, originalLen)
	copy(originalData, t.out[:originalLen])

	// Calculate total size needed
	totalSize := originalLen * n
	if cap(t.out) < totalSize {
		// Expand buffer if needed
		newBuf := make([]byte, 0, totalSize)
		t.out = newBuf
	}

	// Reset and fill buffer efficiently
	t.outLen = 0
	t.out = t.out[:0]

	// Write original data n times
	for range n {
		t.out = append(t.out, originalData...)
	}
	t.outLen = len(t.out)

	return t
}
