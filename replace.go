package tinystring

// Replace replaces all occurrences of old with new in the text content
// eg: "hello world" with old "world" and new "universe" will return "hello universe"
func (t *Text) Replace(old, newStr string) *Text {
	if len(old) == 0 || t.content == "" {
		return t
	}

	var result string
	for i := 0; i < len(t.content); i++ {
		// Check for occurrence of old in the text
		if i+len(old) <= len(t.content) && t.content[i:i+len(old)] == old {
			// Add the new word to the result
			result += newStr
			// Skip the length of the old word in the original text
			i += len(old) - 1
		} else {
			// Add the current character to the result
			result += string(t.content[i])
		}
	}

	t.content = result
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
