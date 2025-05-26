package main

import "github.com/cdvelop/tinystring"

func main() {
	// Same operations using TinyString (no standard library imports)
	text := "MÍ téxtO cön AcÉntos"

	// Text transformations using TinyString
	ts := tinystring.Convert(text)
	lower := ts.ToLower().String()
	upper := ts.ToUpper().String()

	// Number conversions
	intNum := 42
	floatNum := 3.14159
	intStr := tinystring.Convert(intNum).String()
	floatStr := tinystring.Convert(floatNum).RoundDecimals(2).String()

	// String formatting using TinyString
	formatted := tinystring.Format("Text: %s, Numbers: %s, %s", lower, intStr, floatStr)

	// String manipulation
	contains := tinystring.Contains(text, "tÉx")
	replaced := tinystring.Convert(text).Replace("ö", "o").String()

	// Simulate output without standard library
	result := tinystring.Format("%s | %s | %t | %s", formatted, upper, contains, replaced)
	_ = result
}
