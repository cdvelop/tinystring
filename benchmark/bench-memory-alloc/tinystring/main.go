package main

import (
	"github.com/cdvelop/tinystring"
)

// processTextWithTinyString simulates text processing using TinyString
func processTextWithTinyString(texts []string) []string {
	results := make([]string, len(texts))
	for i, text := range texts {
		// Process using TinyString with chained operations
		processed := tinystring.Convert(text).
			RemoveTilde().
			CamelCaseLower().
			String()
		results[i] = processed
	}
	return results
}

// processTextWithTinyStringPointers uses pointer approach for efficiency
func processTextWithTinyStringPointers(texts []string) {
	for i := range texts {
		// Modify in place using pointer approach
		tinystring.Convert(&texts[i]).
			RemoveTilde().
			CamelCaseLower().
			Apply()
	}
}

// processNumbersWithTinyString simulates number processing
func processNumbersWithTinyString(numbers []float64) []string {
	results := make([]string, len(numbers))
	for i, num := range numbers {
		// Convert and format using TinyString
		formatted := tinystring.Convert(num).
			RoundDecimals(2).
			FormatNumber().
			String()
		results[i] = formatted
	}
	return results
}

func main() {
	println("TinyString benchmark main function - used for testing only")
}
