package tinystring

// Truncate truncates a text so that it has exactly the specified width.
// If the text is longer, it truncates it and adds "..." if there is space.
// If the text is shorter, it pads it with spaces until the width is reached.
// The reservedChars parameter indicates how many characters should be reserved for suffixes.
// eg: Convert("Hello, World!").Truncate(10, 3) => "Hel..."
func (t *Text) Truncate(maxWidth int, reservedChars int) *Text {
	text := t.content

	if maxWidth == 0 {
		// Do not truncate
		return t
	}

	actualWidth := maxWidth - reservedChars

	if len(text) > actualWidth {
		// Truncate the text
		if actualWidth > 3 {
			// There is enough space for the truncated text and the ellipsis
			t.content = text[:actualWidth-3] + "..."
		} else {
			// Not enough space for an ellipsis
			t.content = text[:max(1, actualWidth)]
		}
	} else {
		// Pad with spaces
		padding := ""
		for range maxWidth - len(text) {
			padding += " "
		}
		t.content = text + padding
	}

	return t
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
