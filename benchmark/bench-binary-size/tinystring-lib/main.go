package main

import "github.com/cdvelop/tinystring"

func main() {
	// Realistic complex operations using TinyString's chaining capabilities
	conv := "MÍ téxtO cön AcÉntos Y MÁS TEXTO"
	// Complex chained transformations (TinyString's strength)
	processed := tinystring.Convert(conv).
		ToLower().
		RemoveTilde().
		Replace(" ", "_").
		CamelCaseLower().
		Capitalize().
		String()

	// Number processing with chaining
	prices := []any{1234.567, 9876.54, 42.0}
	formattedPrices := make([]string, len(prices))
	for i, price := range prices {
		formattedPrices[i] = tinystring.Convert(price).
			RoundDecimals(2).
			FormatNumber().
			String()
	}

	// Complex string operations with chaining
	userInput := "  Hello@World#2024!  "
	cleaned := tinystring.Convert(userInput).
		Trim().
		Replace("@", "_at_").
		Replace("#", "_hash_").
		Replace("!", "").
		ToLower().
		String()

	// Advanced formatting and joining
	priceList := tinystring.Convert(formattedPrices).Join(", ")
	finalResult := tinystring.Format(
		"Processed: %s | Cleaned: %s | Prices: %s",
		processed, cleaned, priceList,
	)
	// Additional complex transformations
	mixedText := "José María-González_2024"
	normalized := tinystring.Convert(mixedText).
		RemoveTilde().
		Replace("-", "_").
		ToSnakeCaseLower().
		String()

	// Final comprehensive result
	result := tinystring.Format("%s | Normalized: %s", finalResult, normalized)
	_ = result
}
