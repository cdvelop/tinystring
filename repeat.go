package tinystring

// Repeat repeats the conv content n times
// If n is 0 or negative, it clears the conv content
// eg: Convert("abc").Repeat(3) => "abcabcabc"
func (t *conv) Repeat(n int) *conv {
	if t.hasContent(buffErr) {
		return t // Error chain interruption
	}
	if n <= 0 {
		// Clear buffer for empty out and clear dataPtr to prevent reconstruction
		t.rstBuffer(buffOut)
		t.dataPtr = nil // Clear pointer to prevent getString from reconstructing
		return t
	}

	// OPTIMIZED: Direct length check
	if t.outLen == 0 {
		// Clear buffer for empty out
		t.rstBuffer(buffOut)
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
