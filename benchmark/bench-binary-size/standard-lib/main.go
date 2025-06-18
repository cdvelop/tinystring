package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// EQUIVALENT FUNCTIONALITY TEST - Same operations, same complexity
	// Both implementations should do EXACTLY the same work

	// Test 1: Basic string operations
	text1 := "Hello World Example"
	result1 := strings.ToLower(text1)
	result1 = strings.ReplaceAll(result1, " ", "_")

	// Test 2: Number formatting
	num1 := 1234.567
	result2 := strconv.FormatFloat(num1, 'f', 2, 64)

	// Test 3: Multiple string operations
	text2 := "Processing Multiple Strings"
	result3 := strings.ToUpper(text2)
	result3 = strings.ReplaceAll(result3, " ", "-")

	// Test 4: Join operations
	items := []string{"item1", "item2", "item3"}
	result4 := strings.Join(items, ", ")

	// Test 5: Fmt operations
	result5 := fmt.Sprintf("Result: %s | Number: %s | Upper: %s | List: %s",
		result1, result2, result3, result4)

	// Use results to prevent optimization
	_ = result5
}
