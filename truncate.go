package tinystring

// Truncate truncates a text so that it does not exceed the specified width.
// If the text is longer, it truncates it and adds "..." if there is space.
// If the text is shorter or equal to the width, it remains unchanged.
// The reservedChars parameter indicates how many characters should be reserved for suffixes.
// This parameter is optional - if not provided, no characters are reserved (equivalent to passing 0).
// eg: Convert("Hello, World!").Truncate(10) => "Hello, ..."
// eg: Convert("Hello, World!").Truncate(10, 3) => "Hell..."
// eg: Convert("Hello").Truncate(10) => "Hello"
func (t *Text) Truncate(maxWidth any, reservedChars ...any) *Text {
	text := t.content
	originalLength := len(text)

	// Convert maxWidth to integer using the toInt utility function
	maxWidthInt, ok := toInt(maxWidth)
	if !ok || maxWidthInt <= 0 {
		// If maxWidth is zero or invalid, return the original text without modification.
		return t
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
		// Ensure reservedCharsInt does not exceed maxWidthInt
		if reservedCharsInt > maxWidthInt {
			reservedCharsInt = maxWidthInt
		}

		// Calculate the width available for the text itself, excluding reserved chars
		effectiveWidth := max(maxWidthInt-reservedCharsInt, 0)

		canFitEllipsisInEffectiveWidth := effectiveWidth >= ellipsisLen

		if reservedCharsInt > 0 && canFitEllipsisInMaxWidth && canFitEllipsisInEffectiveWidth {
			// Case 1: Reserved chars specified, and ellipsis fits within the effective width
			charsToKeep := min(max(effectiveWidth-ellipsisLen, 0), originalLength)
			t.content = text[:charsToKeep] + ellipsis

		} else if reservedCharsInt == 0 && canFitEllipsisInMaxWidth {
			// Case 2: No reserved chars, ellipsis fits within maxWidth
			charsToKeep := min(max(maxWidthInt-ellipsisLen, 0), originalLength)
			t.content = text[:charsToKeep] + ellipsis
		} else {
			// Case 3: Ellipsis doesn't fit or reserved chars prevent it, just truncate
			charsToKeep := min(maxWidthInt, originalLength)
			t.content = text[:charsToKeep]
		}
		// Truncation happened, no padding needed.

	} // Remove the entire else block that handled padding
	// If originalLength <= maxWidthInt, the text remains unchanged.

	return t
}
