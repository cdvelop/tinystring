package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	// Realistic complex operations using standard library (multiple separate calls)
	conv := "MÍ téxtO cön AcÉntos Y MÁS TEXTO"

	// Complex transformations using multiple standard library calls
	step1 := strings.ToLower(conv)
	step2 := removeTildes(step1)
	step3 := strings.ReplaceAll(step2, " ", "_")
	step4 := toCamelCase(step3)
	processed := capitalize(step4)

	// Number processing with multiple operations
	prices := []float64{1234.567, 9876.54, 42.0}
	formattedPrices := make([]string, len(prices))
	for i, price := range prices {
		rounded := roundFloat(price, 2)
		formattedPrices[i] = formatNumber(rounded)
	}

	// Complex string operations with multiple calls
	userInput := "  Hello@World#2024!  "
	trimmed := strings.TrimSpace(userInput)
	replaced1 := strings.ReplaceAll(trimmed, "@", "_at_")
	replaced2 := strings.ReplaceAll(replaced1, "#", "_hash_")
	replaced3 := strings.ReplaceAll(replaced2, "!", "")
	cleaned := strings.ToLower(replaced3)

	// Manual joining and formatting
	priceList := strings.Join(formattedPrices, ", ")
	finalResult := fmt.Sprintf(
		"Processed: %s | Cleaned: %s | Prices: %s",
		processed, cleaned, priceList,
	)

	// Additional complex transformations
	mixedText := "José María-González_2024"
	normalized1 := removeTildes(mixedText)
	normalized2 := strings.ReplaceAll(normalized1, "-", "_")
	normalized := toSnakeCase(normalized2)

	// Final comprehensive result
	result := fmt.Sprintf("%s | Normalized: %s", finalResult, normalized)
	_ = result
}

// Helper functions to simulate equivalent operations
func removeTildes(s string) string {
	replacements := map[rune]rune{
		'á': 'a', 'é': 'e', 'í': 'i', 'ó': 'o', 'ú': 'u',
		'Á': 'A', 'É': 'E', 'Í': 'I', 'Ó': 'O', 'Ú': 'U',
		'ñ': 'n', 'Ñ': 'N', 'ü': 'u', 'Ü': 'U',
		'ç': 'c', 'Ç': 'C',
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, ok := replacements[r]; ok {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		return s
	}

	var result strings.Builder
	result.WriteString(words[0])
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result.WriteString(strings.ToUpper(words[i][:1]))
			if len(words[i]) > 1 {
				result.WriteString(words[i][1:])
			}
		}
	}
	return result.String()
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func roundFloat(val float64, precision int) float64 {
	ratio := 1.0
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}

func formatNumber(val float64) string {
	str := strconv.FormatFloat(val, 'f', 2, 64)
	parts := strings.Split(str, ".")

	// Add thousand separators
	intPart := parts[0]
	var result strings.Builder
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}

	if len(parts) > 1 {
		result.WriteString(".")
		result.WriteString(parts[1])
	}

	return result.String()
}

func toSnakeCase(s string) string {
	// Convert to snake_case
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)
	return strings.ToLower(snake)
}
