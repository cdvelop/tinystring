package tinystring

// Capitalize transforms the first letter of each word to uppercase and the rest to lowercase.
// For example: "hello world" -> "Hello World"
func (t *conv) Capitalize() *conv {
	// Use local variable instead of struct field to avoid persistent allocation
	words := t.splitIntoWordsLocal()
	if len(words) == 0 {
		return t
	}

	var result []rune
	for i, word := range words {
		if len(word) == 0 {
			continue
		}

		// Add space between words (not before first word)
		if i > 0 {
			result = append(result, ' ')
		}

		// Process each character in the word
		for j, r := range word {
			if j == 0 {
				// First letter of the word - convert to uppercase
				result = append(result, t.transformWord([]rune{r}, toUpper)...)
			} else {
				// Rest of the word - convert to lowercase
				result = append(result, t.transformWord([]rune{r}, toLower)...)
			}
		}
	}

	t.setString(string(result))
	return t
}
