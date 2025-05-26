package tinystring

// Replace replaces up to n occurrences of old with new in the text content
// If n < 0, there is no limit on the number of replacements
// eg: "hello world" with old "world" and new "universe" will return "hello universe"
// Old and new can be any type, they will be converted to string using anyToString
func (t *Text) Replace(oldAny, newAny any, n ...int) *Text {
	old := anyToString(oldAny)
	newStr := anyToString(newAny)

	if len(old) == 0 || t.content == "" {
		return t
	}
	// Default behavior: replace all occurrences (n = -1)
	maxReplacements := -1
	if len(n) > 0 {
		maxReplacements = n[0]
	}

	// Use string builder for efficient string construction
	estimatedLen := len(t.content) + len(newStr)*10 // Conservative estimate
	builder := newTinyStringBuilder(estimatedLen)

	replacements := 0
	for i := 0; i < len(t.content); i++ {
		// Check for occurrence of old in the text and if we haven't reached the maximum replacements
		if i+len(old) <= len(t.content) && t.content[i:i+len(old)] == old && (maxReplacements < 0 || replacements < maxReplacements) {
			// Add the new word to the result
			builder.writeString(newStr)
			// Skip the length of the old word in the original text
			i += len(old) - 1
			// Increment replacement counter
			replacements++
		} else {
			// Add the current character to the result
			builder.writeByte(t.content[i])
		}
	}

	t.content = builder.string()
	return t
}

// TrimSuffix removes the specified suffix from the text content if it exists
// eg: "hello.txt" with suffix ".txt" will return "hello"
func (t *Text) TrimSuffix(suffix string) *Text {
	if len(t.content) < len(suffix) || t.content[len(t.content)-len(suffix):] != suffix {
		return t
	}
	t.content = t.content[:len(t.content)-len(suffix)]
	return t
}

// TrimPrefix removes the specified prefix from the text content if it exists
// eg: "prefix-hello" with prefix "prefix-" will return "hello"
func (t *Text) TrimPrefix(prefix string) *Text {
	if len(t.content) < len(prefix) || t.content[:len(prefix)] != prefix {
		return t
	}
	t.content = t.content[len(prefix):]
	return t
}

// Trim removes spaces at the beginning and end of the text content
// eg: "  hello world  " will return "hello world"
func (t *Text) Trim() *Text {
	if t.content == "" {
		return t
	}

	// Remove spaces at the beginning
	start := 0
	for start < len(t.content) && t.content[start] == ' ' {
		start++
	}

	// Remove spaces at the end and at the end of each line
	end := len(t.content) - 1
	for end >= 0 && (t.content[end] == ' ' || t.content[end] == '\n' || t.content[end] == '\t') {
		end--

	}

	// Special case: empty string
	if start > end {
		t.content = ""
		return t
	}

	// Set the substring without spaces
	t.content = t.content[start : end+1]
	return t
}
