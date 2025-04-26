package tinystring

// Truncate truncates a text so that it has exactly the specified width.
// If the text is longer, it truncates it and adds "..." if there is space.
// If the text is shorter, it pads it with spaces until the width is reached.
// The reservedChars parameter indicates how many characters should be reserved for suffixes.
// This parameter is optional - if not provided, no characters are reserved (equivalent to passing 0).
// eg: Convert("Hello, World!").Truncate(10) => "Hello, ..."
// eg: Convert("Hello, World!").Truncate(10, 3) => "Hel..."
func (t *Text) Truncate(maxWidth any, reservedChars ...any) *Text {
	text := t.content
	originalLength := len(text)

	// Convert maxWidth to integer using the toInt utility function
	maxWidthInt, ok := toInt(maxWidth)
	if !ok || maxWidthInt <= 0 {
		return t // Return original if maxWidth is invalid
	}

	if originalLength > maxWidthInt {
		// --- Truncation Logic ---
		ellipsis := "..."
		ellipsisLen := len(ellipsis)
		canFitEllipsisInMaxWidth := maxWidthInt >= ellipsisLen

		reservedCharsInt := 0
		if len(reservedChars) > 0 {
			if val, ok := toInt(reservedChars[0]); ok && val >= 0 {
				reservedCharsInt = val
			}
		}
		if reservedCharsInt > maxWidthInt {
			reservedCharsInt = maxWidthInt
		}

		effectiveWidth := max(maxWidthInt-reservedCharsInt, 0)

		canFitEllipsisInEffectiveWidth := effectiveWidth >= ellipsisLen

		if canFitEllipsisInMaxWidth && canFitEllipsisInEffectiveWidth {
			// Use ellipsis
			charsToKeep := max(effectiveWidth-ellipsisLen, 0)
			// Ensure slicing doesn't exceed original text length
			if charsToKeep > originalLength { // Added safety check
				charsToKeep = originalLength
			}
			t.content = text[:charsToKeep] + ellipsis
		} else {
			// No ellipsis, just truncate to maxWidth
			// Ensure slicing doesn't exceed original text length
			charsToKeep := maxWidthInt
			if charsToKeep > originalLength { // Added safety check
				charsToKeep = originalLength
			}
			t.content = text[:charsToKeep]
		}
		// Truncation happened, no padding needed.

	} else {
		// --- Padding Logic (only if originalLength <= maxWidthInt) ---
		currentLength := len(t.content) // Should be originalLength here
		if currentLength < maxWidthInt {
			paddingCount := maxWidthInt - currentLength
			padding := ""
			for range paddingCount {
				padding += " "
			}
			t.content = t.content + padding
		}
	}

	return t
}
