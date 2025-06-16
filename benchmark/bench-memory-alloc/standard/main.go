package main

import (
	"fmt"
	"strconv"
	"strings"
)

// processTextWithStandardLib simulates text processing using standard library
func processTextWithStandardLib(texts []string) []string {
	results := make([]string, len(texts))
	for i, text := range texts {
		// EQUIVALENT OPERATIONS: ToLower + simple character replacement + basic word processing
		lowered := strings.ToLower(text)

		// Simple character replacements (equivalent to RemoveTilde)
		replaced := strings.ReplaceAll(lowered, "á", "a")
		replaced = strings.ReplaceAll(replaced, "é", "e")
		replaced = strings.ReplaceAll(replaced, "í", "i")
		replaced = strings.ReplaceAll(replaced, "ó", "o")
		replaced = strings.ReplaceAll(replaced, "ú", "u")
		replaced = strings.ReplaceAll(replaced, "ñ", "n")

		// Simple word processing - split, capitalize first word, join
		words := strings.Fields(replaced)
		if len(words) > 0 && len(words[0]) > 0 {
			words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
		}
		result := strings.Join(words, "")
		results[i] = result
	}
	return results
}

// processNumbersWithStandardLib simulates number processing
func processNumbersWithStandardLib(numbers []float64) []string {
	results := make([]string, len(numbers))
	for i, num := range numbers {
		// EQUIVALENT OPERATIONS: Format with 2 decimals + add thousand separators
		formatted := strconv.FormatFloat(num, 'f', 2, 64)

		// Add thousand separators (equivalent to FormatNumber)
		parts := strings.Split(formatted, ".")
		integer := parts[0]
		decimal := parts[1]

		// Simple thousand separator logic
		if len(integer) > 3 {
			var result strings.Builder
			for j, char := range integer {
				if j > 0 && (len(integer)-j)%3 == 0 {
					result.WriteString(".")
				}
				result.WriteRune(char)
			}
			formatted = result.String() + "," + decimal
		}
		results[i] = formatted
	}
	return results
}

func main() {
	fmt.Println("Standard library benchmark main function - used for testing only")
}
