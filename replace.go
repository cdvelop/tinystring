package tinystring

// Replace replaces up to n occurrences of old with new in the conv content
// If n < 0, there is no limit on the number of replacements
// eg: "hello world" with old "world" and new "universe" will return "hello universe"
// Old and new can be any type, they will be converted to string using Convert
func (t *conv) Replace(oldAny, newAny any, n ...int) *conv {
	// Convert parameters to strings using the consistent Convert pattern
	old := Convert(oldAny).String()
	newStr := Convert(newAny).String()
	str := t.getString()
	if isEmpty(old) || isEmpty(str) {
		return t
	}
	// Estimate buffer capacity based on replacement patterns
	estimatedCap := len(str)
	if len(newStr) > len(old) {
		// If new string is longer, estimate extra space needed
		estimatedCap += (len(newStr) - len(old)) * 5 // Assume up to 5 replacements
	}
	buf := makeBuf(estimatedCap)

	// Default behavior: replace all occurrences (n = -1)
	maxReps := -1
	if hasLength(n) {
		maxReps = n[0]
	}

	rep := 0
	for i := 0; i < len(str); i++ {
		// Check for occurrence of old in the string and if we haven't reached the maximum rep
		if i+len(old) <= len(str) && str[i:i+len(old)] == old && (maxReps < 0 || rep < maxReps) {
			// Add the new word to the result
			buf = append(buf, newStr...)
			// Skip the length of the old word in the original string
			i += len(old) - 1
			// Increment replacement counter
			rep++
		} else {
			// Add the current character to the result
			buf = append(buf, str[i])
		}
	}

	t.setString(string(buf))
	return t
}

// TrimSuffix removes the specified suffix from the conv content if it exists
// eg: "hello.txt" with suffix ".txt" will return "hello"
func (t *conv) TrimSuffix(suffix string) *conv {
	str := t.getString()
	if len(str) < len(suffix) || str[len(str)-len(suffix):] != suffix {
		return t
	}
	t.setString(str[:len(str)-len(suffix)])
	return t
}

// TrimPrefix removes the specified prefix from the conv content if it exists
// eg: "prefix-hello" with prefix "prefix-" will return "hello"
func (t *conv) TrimPrefix(prefix string) *conv {
	str := t.getString()
	if len(str) < len(prefix) || str[:len(prefix)] != prefix {
		return t
	}
	t.setString(str[len(prefix):])
	return t
}

// Trim removes spaces at the beginning and end of the conv content
// eg: "  hello world  " will return "hello world"
func (t *conv) Trim() *conv {
	str := t.getString()
	if isEmpty(str) {
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
		t.setString("")
		return t
	}

	// Set the substring without spaces
	t.setString(str[start : end+1])
	return t
}
