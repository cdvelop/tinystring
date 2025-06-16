package tinystring

import (
	"testing"
)

// Phase 11: String Operations Benchmarks
// Focus on string manipulation functions that could benefit from optimization

func BenchmarkStringOperations(b *testing.B) {
	testData := []string{
		"hello world test string",
		"CONVERT TO LOWERCASE",
		"convert_to_camelCase",
		"remove-special-chars",
		"trim   spaces   ",
		"José María González", // With accents
		"test@email.com",
		"long string with multiple words for processing",
	}

	b.Run("ToLower", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Convert(str).ToLower().String()
				_ = result
			}
		}
	})

	b.Run("ToUpper", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Convert(str).ToUpper().String()
				_ = result
			}
		}
	})

	b.Run("Capitalize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Convert(str).Capitalize().String()
				_ = result
			}
		}
	})

	b.Run("RemoveTilde", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Convert(str).RemoveTilde().String()
				_ = result
			}
		}
	})

	b.Run("CamelCaseLower", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Convert(str).CamelCaseLower().String()
				_ = result
			}
		}
	})

	b.Run("Replace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Convert(str).Replace(" ", "_").String()
				_ = result
			}
		}
	})

	b.Run("Trim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Convert(str).Trim().String()
				_ = result
			}
		}
	})

	b.Run("Split", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				result := Split(str, " ")
				_ = result
			}
		}
	})
}

// Benchmark specific string processing chains
func BenchmarkStringChains(b *testing.B) {
	testStrings := []string{
		"José María González Pérez",
		"CONVERT_TO_PROPER_CASE",
		"email@domain.com test string",
		"  trim and process  ",
	}

	b.Run("ComplexChain1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				result := Convert(str).RemoveTilde().ToLower().Capitalize().String()
				_ = result
			}
		}
	})

	b.Run("ComplexChain2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				result := Convert(str).Trim().Replace(" ", "_").ToLower().String()
				_ = result
			}
		}
	})

	b.Run("CamelCaseChain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				result := Convert(str).RemoveTilde().CamelCaseLower().String()
				_ = result
			}
		}
	})
}
