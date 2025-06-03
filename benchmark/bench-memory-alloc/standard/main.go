package main

import (
	"fmt"
	"strings"
)

// processTextWithStandardLib simulates conv processing using standard library
func processTextWithStandardLib(texts []string) []string {
	results := make([]string, len(texts))
	for i, conv := range texts {
		// Simulate common string operations using standard library
		lowered := strings.ToLower(conv)
		replaced := strings.ReplaceAll(lowered, "á", "a")
		replaced = strings.ReplaceAll(replaced, "é", "e")
		replaced = strings.ReplaceAll(replaced, "í", "i")
		replaced = strings.ReplaceAll(replaced, "ó", "o")
		replaced = strings.ReplaceAll(replaced, "ú", "u")

		// Additional processing
		words := strings.Fields(replaced)
		if len(words) > 0 {
			// Capitalize first letter of first word
			if len(words[0]) > 0 {
				words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
			}
			replaced = strings.Join(words, "")
		}
		results[i] = replaced
	}
	return results
}

// processNumbersWithStandardLib simulates number processing
func processNumbersWithStandardLib(numbers []float64) []string {
	results := make([]string, len(numbers))
	for i, num := range numbers {
		// Convert to string and format
		formatted := fmt.Sprintf("%.2f", num)
		// Add thousand separators using standard library
		parts := strings.Split(formatted, ".")
		integer := parts[0]
		decimal := parts[1]

		// Add thousand separators
		if len(integer) > 3 {
			result := ""
			for j, char := range integer {
				if j > 0 && (len(integer)-j)%3 == 0 {
					result += ","
				}
				result += string(char)
			}
			formatted = result + "." + decimal
		}
		results[i] = formatted
	}
	return results
}

func main() {
	fmt.Println("Standard library benchmark main function - used for testing only")
}
