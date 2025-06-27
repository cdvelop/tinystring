package main

import "github.com/cdvelop/tinystring"

func main() {
	// EQUIVALENT FUNCTIONALITY TEST - Same operations, same complexity
	// Both implementations should do EXACTLY the same work

	// Test 1: Basic string operations
	text1 := "Hello World Example"
	result1 := tinystring.Convert(text1).ToLower().Replace(" ", "_").String()

	// Test 2: Number formatting
	num1 := 1234.567
	result2 := tinystring.Convert(num1).Round(2).String()

	// Test 3: Multiple string operations
	text2 := "Processing Multiple Strings"
	result3 := tinystring.Convert(text2).ToUpper().Replace(" ", "-").String()

	// Test 4: Join operations
	items := []string{"item1", "item2", "item3"}
	result4 := tinystring.Convert(items).Join(", ").String()

	// Test 5: Fmt operations
	result5 := tinystring.Fmt("Result: %s | Number: %s | Upper: %s | List: %s",
		result1, result2, result3, result4)

	// Use results to prevent optimization
	_ = result5
}
