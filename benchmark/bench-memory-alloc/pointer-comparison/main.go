package main

import (
	"fmt"

	"github.com/cdvelop/tinystring"
)

// processTextWithoutPointer simulates traditional processing without pointers
func processTextWithoutPointer(text string) string {
	return tinystring.Convert(text).RemoveTilde().CamelCaseLower().ToSnakeCaseLower().String()
}

// processTextWithPointer simulates processing using string pointers
func processTextWithPointer(text *string) {
	tinystring.Convert(text).RemoveTilde().CamelCaseLower().ToSnakeCaseLower().Apply()
}

// processMassiveTextsWithoutPointer simulates massive processing without pointers
func processMassiveTextsWithoutPointer(texts []string) []string {
	results := make([]string, len(texts))
	for j, text := range texts {
		results[j] = tinystring.Convert(text).RemoveTilde().CamelCaseLower().String()
	}
	return results
}

// processMassiveTextsWithPointer simulates massive processing using pointers
func processMassiveTextsWithPointer(texts []string) {
	for j := range texts {
		tinystring.Convert(&texts[j]).RemoveTilde().CamelCaseLower().Apply()
	}
}

func main() {
	fmt.Println("Pointer comparison benchmark")
}
