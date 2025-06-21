package tinystring

// Replace replaces up to n occurrences of old with new in the conv content
// If n < 0, there is no limit on the number of replacements
// eg: "hello world" with old "world" and new "universe" will return "hello universe"
// Old and new can be any type, they will be converted to string using Convert
func (t *conv) Replace(oldAny, newAny any, n ...int) *conv {
	if t.hasError() {
		return t // Error chain interruption
	}
	tmp := getConv()
	tmp.any2s(oldAny) // Convert oldAny to string
	old := string(tmp.work[:tmp.workLen])

	tmp.any2s(newAny)      // Convert newAny to string
	newStr := tmp.String() // return tmp to pool

	// Convert parameters to strings using the consistent Convert pattern
	str := t.ensureStringInOut()
	if len(old) == 0 || len(str) == 0 {
		return t
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
	t.rstOut()     // Clear buffer using API
	t.wrToOut(out) // Write using API
	return t
}

// TrimSuffix removes the specified suffix from the conv content if it exists
// eg: "hello.txt" with suffix ".txt" will return "hello"
func (t *conv) TrimSuffix(suffix string) *conv {
	if t.hasError() {
		return t // Error chain interruption
	}

	str := t.ensureStringInOut()
	if len(str) < len(suffix) || str[len(str)-len(suffix):] != suffix {
		return t
	} // ✅ Update buffer using API instead of direct manipulation
	out := str[:len(str)-len(suffix)]
	t.rstOut()           // Clear buffer using API
	t.wrStringToOut(out) // Write using API
	return t
}

// TrimPrefix removes the specified prefix from the conv content if it exists
// eg: "prefix-hello" with prefix "prefix-" will return "hello"
func (t *conv) TrimPrefix(prefix string) *conv {
	if t.hasError() {
		return t // Error chain interruption
	}

	str := t.ensureStringInOut()
	if len(str) < len(prefix) || str[:len(prefix)] != prefix {
		return t
	} // ✅ Update buffer using API instead of direct manipulation
	out := str[len(prefix):]
	t.rstOut()           // Clear buffer using API
	t.wrStringToOut(out) // Write using API
	return t
}

// Trim removes spaces at the beginning and end of the conv content
// eg: "  hello world  " will return "hello world"
func (t *conv) Trim() *conv {
	if t.hasError() {
		return t // Error chain interruption
	}

	str := t.ensureStringInOut()
	if len(str) == 0 {
		return t
	}

	// Remove spaces at the beginning
	start := 0
	for start < len(str) && str[start] == ' ' {
		start++
	}

	// Remove spaces at the end and at the end of each line
	end := len(str) - 1
	for end >= 0 && (str[end] == ' ' || str[end] == '\n' || str[end] == '\t') {
		end--
	} // Special case: empty string
	if start > end {
		// Clear buffer for empty out
		t.rstOut()
		return t
	}
	// ✅ Set the substring without spaces using API
	out := str[start : end+1]
	t.rstOut()           // Clear buffer using API
	t.wrStringToOut(out) // Write using API
	return t
}
