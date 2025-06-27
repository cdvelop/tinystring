package main

import (
	"github.com/cdvelop/tinystring"
)

// processTextWithTinyString simulates text processing using TinyString (equivalent to standard lib)
func processTextWithTinyString(texts []string) []string {
	results := make([]string, len(texts))
	for i, text := range texts {
		// EQUIVALENT OPERATIONS: Same logic as standard library but using TinyString
		processed := tinystring.Convert(text).
			Low().
			Tilde().
			String()

		// Split into words and capitalize first letter of first word, then join
		words := tinystring.Split(processed, " ")
		if len(words) > 0 && len(words[0]) > 0 {
			words[0] = tinystring.Convert(words[0]).Capitalize().String()
		}

		// Join words without spaces (equivalent to standard lib behavior)
		out := tinystring.Convert(words).Join().String()
		results[i] = out
	}
	return results
}

// processNumbersWithTinyString simulates number processing (equivalent to standard lib)
func processNumbersWithTinyString(numbers []float64) []string {
	results := make([]string, len(numbers))
	for i, num := range numbers {
		// EQUIVALENT OPERATIONS: Same formatting as standard library
		formatted := tinystring.Convert(num).
			Round(2).
			Thousands().
			String()
		results[i] = formatted
	}
	return results
}

func main() {
	println("TinyString benchmark main function - used for testing only")
}
