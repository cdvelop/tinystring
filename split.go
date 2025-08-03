package tinystring

// Split divides a string by a separator and returns a slice of substrings.
// Usage: Convert("Hello World").Split() => []string{"Hello", "World"}
// Usage with separator: Convert("Hello;World").Split(";") => []string{"Hello", "World"}
// If no separator is provided, splits by whitespace (similar to strings.Fields).
// Uses the conv work buffer for memory efficiency. The global Split function is deprecated; always use Convert(...).Split(...).

func (c *conv) Split(separator ...string) []string {
	src := c.getString(buffOut)
	return c.splitStr(src, separator...)
}

// splitStr is a reusable internal method for splitting a string by a separator (empty = by character, default whitespace).
func (c *conv) splitStr(src string, separator ...string) []string {
	var sep string
	if len(separator) == 0 {
		// Whitespace split: mimic strings.Fields
		out := make([]string, 0, len(src)/2+1)
		fieldStart := -1
		for i, r := range src {
			isSpace := r == ' ' || r == '\t' || r == '\n' || r == '\r'
			if !isSpace {
				if fieldStart == -1 {
					fieldStart = i
				}
			} else {
				if fieldStart != -1 {
					out = append(out, src[fieldStart:i])
					fieldStart = -1
				}
			}
		}
		if fieldStart != -1 {
			out = append(out, src[fieldStart:])
		}
		return out
	} else {
		sep = separator[0]
	}
	// Special case: split by character (empty separator)
	if len(sep) == 0 {
		if len(src) == 0 {
			return []string{}
		}
		out := make([]string, 0, len(src))
		for _, ch := range src {
			// OPTIMIZED: Direct string conversion without buffer operations
			out = append(out, string(ch))
		}
		return out
	}
	// Handle string shorter than 3 chars (legacy behavior)
	if len(src) < 3 {
		return []string{src}
	}
	// If src is empty, return [""] (legacy behavior)
	if len(src) == 0 {
		return []string{""}
	}
	// Use splitByDelimiterWithBuffer for all splits
	var out []string
	first := true
	orig := src
	for {
		before, after, found := c.splitByDelimiterWithBuffer(src, sep)
		out = append(out, before)
		if !found {
			// Legacy: if separator not found at all, return original string as single element
			if first && len(out) == 1 && out[0] == orig {
				return []string{orig}
			}
			break
		}
		src = after
		first = false
	}
	return out
}

// splitByDelimiterWithBuffer splits a string by the first occurrence of the delimiter.
// Returns the parts (before and after the delimiter). If not found, found=false.
// OPTIMIZED: Direct substring operations without buffer usage
func (c *conv) splitByDelimiterWithBuffer(s, delim string) (before, after string, found bool) {
	di := -1
	for i := 0; i <= len(s)-len(delim); i++ {
		if s[i:i+len(delim)] == delim {
			di = i
			break
		}
	}
	if di < 0 {
		return s, "", false
	}
	// OPTIMIZED: Direct substring without buffer operations
	before = s[:di]
	after = s[di+len(delim):]
	return before, after, true
}
