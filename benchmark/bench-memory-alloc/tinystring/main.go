package main

import (
	. "github.com/cdvelop/tinystring"
)

// processTextWithTinyString simulates text processing using TinyString (equivalent to standard lib)
func processTextWithTinyString(texts []string) []string {
	results := make([]string, len(texts))
	for i, text := range texts {
		out := Convert(text).
			Low().
			Tilde().
			Capitalize().
			String()
		results[i] = out
	}
	return results
}

// processNumbersWithTinyString simulates number processing (equivalent to standard lib)
func processNumbersWithTinyString(numbers []float64) []string {
	results := make([]string, len(numbers))
	for i, num := range numbers {
		// EQUIVALENT OPERATIONS: Same formatting as standard library
		formatted := Convert(num).
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
