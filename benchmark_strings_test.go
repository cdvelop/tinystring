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

// BenchmarkBuilderOperations tests the new builder API performance
func BenchmarkBuilderOperations(b *testing.B) {
	b.Run("BuilderVsConcat", func(b *testing.B) {
		words := []string{"hello", "tiny", "string", "builder", "performance"}

		b.Run("BuilderPattern", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := Convert() // Empty initialization
				for _, word := range words {
					c.Write(word).Write(" ")
				}
				_ = c.String()
			}
		})

		b.Run("MultipleAllocations", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result := ""
				for _, word := range words {
					result += Convert(word).String() + " "
				}
				_ = result
			}
		})
	})

	b.Run("UnifiedWrite", func(b *testing.B) {
		testValues := []any{"string", 42, 3.14, true, 'x'}

		for i := 0; i < b.N; i++ {
			c := Convert()
			for _, val := range testValues {
				c.Write(val).Write(" ")
			}
			_ = c.String()
		}
	})

	b.Run("ChainedOperations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := Convert("hello").ToUpper().Write(" WORLD").ToLower().String()
			_ = result
		}
	})

	b.Run("JoinSliceOptimization", func(b *testing.B) {
		items := []string{"item1", "item2", "item3", "item4", "item5"}

		for i := 0; i < b.N; i++ {
			result := Convert(items).String()
			_ = result
		}
	})
}

// BenchmarkHighDemandProcesses tests the critical optimization targets
func BenchmarkHighDemandProcesses(b *testing.B) {
	b.Run("TransformationChains", func(b *testing.B) {
		testStr := "Hello World Testing String"

		for i := 0; i < b.N; i++ {
			result := Convert(testStr).ToLower().Capitalize().ToUpper().String()
			_ = result
		}
	})

	b.Run("FormatOperations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := Fmt("User %s has %d messages with %.2f%% completion", "Alice", 42, 85.5).String()
			_ = result
		}
	})

	b.Run("ErrorMessageConstruction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := Err("Operation failed for user", "Alice", "with error code", 500).String()
			_ = result
		}
	})
}
