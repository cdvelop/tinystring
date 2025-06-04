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

// convert to lower case eg: "HELLO WORLD" -> "hello world"
func (t *conv) ToLower() *conv {
	return t.transformWithMapping(lowerMappings)
}

// convert to upper case eg: "hello world" -> "HELLO WORLD"
func (t *conv) ToUpper() *conv {
	return t.transformWithMapping(upperMappings)
}

// converts conv to camelCase (first word lowercase) eg: "Hello world" -> "helloWorld"
func (t *conv) CamelCaseLower() *conv {
	return t.toCaseTransform(true, "")
}

// converts conv to PascalCase (all words capitalized) eg: "hello world" -> "HelloWorld"
func (t *conv) CamelCaseUpper() *conv {
	return t.toCaseTransform(false, "")
}

// snakeCase converts a string to snake_case format with optional separator.
// If no separator is provided, underscore "_" is used as default.
// Example:
//
//	Input: "camelCase" -> Output: "camel_case"
//	Input: "PascalCase", "-" -> Output: "pascal-case"
//	Input: "APIResponse" -> Output: "api_response"
//	Input: "user123Name", "." -> Output: "user123.name"
//
// ToSnakeCaseLower converts conv to snake_case format
func (t *conv) ToSnakeCaseLower(sep ...string) *conv {
	return t.toCaseTransform(true, t.separatorCase(sep...))
}

// ToSnakeCaseUpper converts conv to Snake_Case format
func (t *conv) ToSnakeCaseUpper(sep ...string) *conv {
	return t.toCaseTransform(false, t.separatorCase(sep...))
}

func (t *conv) toCaseTransform(firstWordLower bool, separator string) *conv {
	// Use local variable instead of struct field to avoid persistent allocation
	words := t.splitIntoWordsLocal()
	if len(words) == 0 {
		return t
	}

	// Use pooled string builder for efficient string construction
	builder := getBuilder()
	defer putBuilder(builder)

	str := t.getString()
	estimatedLen := len(str) + len(words)*len(separator)
	builder.grow(estimatedLen)
	var prevIsDigit bool
	var prevIsSeparator bool

	for i, word := range words {
		if len(word) == 0 {
			continue
		}

		// Add separator if needed
		if i > 0 && separator != "" {
			builder.writeByte(separator[0])
			prevIsSeparator = true
		}

		// Process each character in the word
		for j, r := range word {
			currentCaseTransform := toLower // Default to lower
			currIsDigit := isDigit(r)
			currIsLetter := isLetter(r)

			// Determine case transform
			if i == 0 && j == 0 { // First letter of first word
				if !firstWordLower {
					currentCaseTransform = toUpper
				}
			} else if i > 0 && j == 0 && separator == "" { // Start of new word in camelCase
				currentCaseTransform = toUpper
			} else if prevIsDigit && currIsLetter { // Letter after digit
				// For snake_case with separator, this is handled by adding separator later
				// For camelCase, new word starts, so apply upper if not firstWordLower
				if separator == "" && !firstWordLower {
					currentCaseTransform = toUpper
				} else if separator == "" && firstWordLower {
					// Maintain lower case if it's camelCaseLower and after a digit within the same "word" part
					currentCaseTransform = toLower
				} else if separator != "" && !firstWordLower { // Snake_Case_Upper
					currentCaseTransform = toUpper
				}

			} else if prevIsSeparator && currIsLetter { // Letter after separator (for snake_case)
				if separator != "" && !firstWordLower { // Snake_Case_Upper
					currentCaseTransform = toUpper
				}
				// Default toLower is fine for snake_case_lower
			} else if currIsLetter && j > 0 { // Subsequent letters in a word part
				// Maintain lower case unless it's an uppercase letter that should remain uppercase (e.g. in "APIResponse")
				// This part is tricky without knowing the original casing or specific rules for acronyms.
				// For now, default toLower is applied unless specific conditions for toUpper are met.
				// If the global transform is toUpper (e.g. ToSnakeCaseUpper), this will be handled.
				if separator != "" && !firstWordLower { // Snake_Case_Upper
					// This condition might be too broad.
					// We only want to uppercase the first letter after a separator.
					// Let's rely on the j==0 condition for this.
					// For subsequent letters in Snake_Case_Upper, they should be lower.
					// So, if j > 0, it should be toLower for Snake_Case_Upper.
					// This means the currentCaseTransform should be toLower here.
				}
			}

			// Add underscore for number to letter transition in snake_case
			if separator != "" && prevIsDigit && currIsLetter {
				builder.writeByte(separator[0])
			}

			if currIsLetter {
				var transformedRune rune
				var mapped bool
				if currentCaseTransform == toLower {
					transformedRune, mapped = transformSingleRune(r, lowerMappings)
				} else { // toUpper
					transformedRune, mapped = transformSingleRune(r, upperMappings)
				}
				if mapped {
					builder.writeRune(transformedRune)
				} else {
					builder.writeRune(r) // Write original if no mapping found (e.g. already correct case)
				}
			} else {
				builder.writeRune(r)
			}

			prevIsDigit = currIsDigit
			prevIsSeparator = false // Reset after processing the character
		}
	}

	t.setString(builder.string())
	// Clear the separator field after use to avoid memory overhead
	t.separator = ""
	return t
}
