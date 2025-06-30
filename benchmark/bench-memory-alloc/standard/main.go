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
		lowered := strings.ToLower(text)
		replaced := strings.ReplaceAll(lowered, "á", "a")
		replaced = strings.ReplaceAll(replaced, "é", "e")
		replaced = strings.ReplaceAll(replaced, "í", "i")
		replaced = strings.ReplaceAll(replaced, "ó", "o")
		replaced = strings.ReplaceAll(replaced, "ú", "u")
		replaced = strings.ReplaceAll(replaced, "ñ", "n")

		// Capitalizar solo la primera letra del string completo
		out := replaced
		if len(out) > 0 {
			out = strings.ToUpper(out[:1]) + out[1:]
		}
		results[i] = out
	}
	return results
}

// processNumbersWithStandardLib simulates number processing
func processNumbersWithStandardLib(numbers []float64) []string {
	results := make([]string, len(numbers))
	for i, num := range numbers {
		// EQUIVALENT OPERATIONS: Fmt with 2 decimals + add thousand separators
		formatted := strconv.FormatFloat(num, 'f', 2, 64)

		// Add thousand separators (equivalent to Thousands)
		parts := strings.Split(formatted, ".")
		integer := parts[0]
		decimal := parts[1]

		// Simple thousand separator logic
		if len(integer) > 3 {
			var out strings.Builder
			for j, char := range integer {
				if j > 0 && (len(integer)-j)%3 == 0 {
					out.WriteString(".")
				}
				out.WriteRune(char)
			}
			formatted = out.String() + "," + decimal
		}
		results[i] = formatted
	}
	return results
}

func main() {
	fmt.Println("Standard library benchmark main function - used for testing only")
}
