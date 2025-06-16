package main

import (
	"benchmark/shared"
	"testing"
)

func BenchmarkStringProcessing(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processTextWithTinyString(shared.TestTexts)
	}
}

func BenchmarkNumberProcessing(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processNumbersWithTinyString(shared.TestNumbers)
	}
}

func BenchmarkMixedOperations(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make(map[string]string)
		for key, value := range shared.TestMixedData {
			switch v := value.(type) {
			case string:
				processed := processTextWithTinyString([]string{v})[0]
				results[key] = processed
			case float64:
				processed := processNumbersWithTinyString([]float64{v})[0]
				results[key] = processed
			}
		}
		_ = results
	}
}
