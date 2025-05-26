package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Common string operations that TinyString replaces
	text := "MÍ téxtO cön AcÉntos"

	// Text transformations using standard library
	lower := strings.ToLower(text)
	upper := strings.ToUpper(text)

	// Number conversions
	intNum := 42
	floatNum := 3.14159
	intStr := strconv.Itoa(intNum)
	floatStr := strconv.FormatFloat(floatNum, 'f', 2, 64)

	// String formatting
	formatted := fmt.Sprintf("Text: %s, Numbers: %s, %s", lower, intStr, floatStr)

	// String manipulation
	contains := strings.Contains(text, "tÉx")
	replaced := strings.ReplaceAll(text, "ö", "o")

	// Simulate output without actually printing (for WASM compatibility)
	result := fmt.Sprintf("%s | %s | %t | %s", formatted, upper, contains, replaced)
	_ = result
}
