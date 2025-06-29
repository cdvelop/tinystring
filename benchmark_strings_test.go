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

	b.Run("Low", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).Low().String()
				_ = out
			}
		}
	})

	b.Run("Up", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).Up().String()
				_ = out
			}
		}
	})

	b.Run("Capitalize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).Capitalize().String()
				_ = out
			}
		}
	})

	b.Run("Tilde", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).Tilde().String()
				_ = out
			}
		}
	})

	b.Run("CamelLow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).CamelLow().String()
				_ = out
			}
		}
	})

	b.Run("Replace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).Replace(" ", "_").String()
				_ = out
			}
		}
	})

	b.Run("Trim", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).Trim().String()
				_ = out
			}
		}
	})

	b.Run("Split", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testData {
				out := Convert(str).Split()
				_ = out
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
				out := Convert(str).Tilde().Low().Capitalize().String()
				_ = out
			}
		}
	})

	b.Run("ComplexChain2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				out := Convert(str).Trim().Replace(" ", "_").Low().String()
				_ = out
			}
		}
	})

	b.Run("CamelCaseChain", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, str := range testStrings {
				out := Convert(str).Tilde().CamelLow().String()
				_ = out
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
				out := ""
				for _, word := range words {
					out += Convert(word).String() + " "
				}
				_ = out
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
			out := Convert("hello").Up().Write(" WORLD").Low().String()
			_ = out
		}
	})

	b.Run("JoinSliceOptimization", func(b *testing.B) {
		items := []string{"item1", "item2", "item3", "item4", "item5"}

		for i := 0; i < b.N; i++ {
			out := Convert(items).String()
			_ = out
		}
	})
}

// BenchmarkHighDemandProcesses tests the critical optimization targets
func BenchmarkHighDemandProcesses(b *testing.B) {
	b.Run("TransformationChains", func(b *testing.B) {
		testStr := "Hello World Testing String"

		for i := 0; i < b.N; i++ {
			out := Convert(testStr).Low().Capitalize().Up().String()
			_ = out
		}
	})

	b.Run("FormatOperations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out := Fmt("User %s has %d messages with %.2f%% completion", "Alice", 42, 85.5)
			_ = out
		}
	})

	b.Run("ErrorMessageConstruction", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out := Err("Operation failed for user", "Alice", "with error code", 500).String()
			_ = out
		}
	})
}
