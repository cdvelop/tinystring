package tinystring

// Replace replaces up to n occurrences of old with new in the conv content
// If n < 0, there is no limit on the number of replacements
// eg: "hello world" with old "world" and new "universe" will return "hello universe"
// Old and new can be any type, they will be converted to string using Convert
func (c *conv) Replace(oldAny, newAny any, n ...int) *conv {
	if c.hasContent(buffErr) {
		return c // Error chain interruption
	}

	// Get the original string before any conversions
	str := c.getBuffString()

	// Preserve original state before temporary conversions
	originalDataPtr := c.dataPtr
	originalKind := c.Kind

	// Use internal work buffer instead of getConv() for zero-allocation
	c.rstBuffer(buffWork)         // Clear work buffer
	c.anyToBuff(buffWork, oldAny) // Convert oldAny to work buffer
	old := c.getString(buffWork)  // Get old string from work buffer

	c.rstBuffer(buffWork)           // Clear work buffer for next conversion
	c.anyToBuff(buffWork, newAny)   // Convert newAny to work buffer
	newStr := c.getString(buffWork) // Get new string from work buffer

	// Restore original state after temporary conversions
	c.dataPtr = originalDataPtr
	c.Kind = originalKind

	// Check early return condition
	if len(old) == 0 || len(str) == 0 {
		return c
	}

	// Estimate buffer capacity based on replacement patterns
	estimatedCap := len(str)
	if len(newStr) > len(old) {
		// If new string is longer, estimate extra space needed
		estimatedCap += (len(newStr) - len(old)) * 5 // Assume up to 5 replacements
	}
	// Inline makeBuf logic
	bufCap := estimatedCap
	if bufCap < defaultBufCap {
		bufCap = defaultBufCap
	}
	out := make([]byte, 0, bufCap)
	// Default behavior: replace all occurrences (n = -1)
	maxReps := -1
	if len(n) > 0 {
		maxReps = n[0]
	}

	rep := 0
	for i := 0; i < len(str); i++ {
		// Check for occurrence of old in the string and if we haven't reached the maximum rep
		if i+len(old) <= len(str) && str[i:i+len(old)] == old && (maxReps < 0 || rep < maxReps) {
			// Add the new word to the out
			out = append(out, newStr...)
			// Skip the length of the old word in the original string
			i += len(old) - 1
			// Increment replacement counter
			rep++
		} else {
			// Add the current character to the out
			out = append(out, str[i])
		}
	}
	// ✅ Update buffer using API instead of direct manipulation
	c.rstBuffer(buffOut)    // Clear buffer using API
	c.wrBytes(buffOut, out) // Write using API
	return c
}

// TrimSuffix removes the specified suffix from the conv content if it exists
// eg: "hello.txt" with suffix ".txt" will return "hello"
func (c *conv) TrimSuffix(suffix string) *conv {
	if c.hasContent(buffErr) {
		return c // Error chain interruption
	}

	str := c.getBuffString()
	if len(str) < len(suffix) || str[len(str)-len(suffix):] != suffix {
		return c
	} // ✅ Update buffer using API instead of direct manipulation
	out := str[:len(str)-len(suffix)]
	c.rstBuffer(buffOut)     // Clear buffer using API
	c.wrString(buffOut, out) // Write using API
	return c
}

// TrimPrefix removes the specified prefix from the conv content if it exists
// eg: "prefix-hello" with prefix "prefix-" will return "hello"
func (c *conv) TrimPrefix(prefix string) *conv {
	if c.hasContent(buffErr) {
		return c // Error chain interruption
	}

	str := c.getBuffString()
	if len(str) < len(prefix) || str[:len(prefix)] != prefix {
		return c
	} // ✅ Update buffer using API instead of direct manipulation
	out := str[len(prefix):]
	c.rstBuffer(buffOut)     // Clear buffer using API
	c.wrString(buffOut, out) // Write using API
	return c
}

// Trim removes spaces at the beginning and end of the conv content
// eg: "  hello world  " will return "hello world"
func (c *conv) Trim() *conv {
	if c.hasContent(buffErr) {
		return c // Error chain interruption
	}

	str := c.getBuffString()
	if len(str) == 0 {
		return c
	}

	// Remove spaces at the beginning
	start := 0
	for start < len(str) && (str[start] == ' ' || str[start] == '\t' || str[start] == '\n' || str[start] == '\r') {
		start++
	}

	// Remove spaces at the end
	end := len(str) - 1
	for end >= 0 && (str[end] == ' ' || str[end] == '\t' || str[end] == '\n' || str[end] == '\r') {
		end--
	}

	// Debug: For "  " input, start should be 2, end should be -1, so start > end
	// Special case: empty string (all whitespace)
	if start > end {
		// Clear buffer and write empty string
		c.rstBuffer(buffOut)
		c.wrString(buffOut, "")
		// Also clear dataPtr to prevent fallback
		c.dataPtr = nil
		c.Kind = K.String
		return c
	}

	// Set the substring without spaces using API
	out := str[start : end+1]
	c.rstBuffer(buffOut)     // Clear buffer using API
	c.wrString(buffOut, out) // Write using API
	return c
}
