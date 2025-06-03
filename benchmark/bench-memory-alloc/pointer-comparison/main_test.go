package main

import (
	"testing"
	"strings"
	"strconv"
	
	"github.com/cdvelop/tinystring"
)

// Shared test data - centralized for consistency
var testTexts = []string{
	"Él Múrcielago Rápido",
	"PROCESANDO textos LARGOS", 
	"Optimización de MEMORIA",
	"Rendimiento en APLICACIONES",
	"Reducción de ASIGNACIONES",
}

// TinyString benchmarks - traditional approach (creates new strings)
func BenchmarkTinyStringWithoutPointer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make([]string, len(testTexts))
		for j, text := range testTexts {
			results[j] = tinystring.Convert(text).RemoveTilde().CamelCaseLower().String()
		}
		_ = results
	}
}

// TinyString benchmarks - pointer optimization (modifies in-place)
func BenchmarkTinyStringWithPointer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create fresh copy for each iteration since pointers modify in-place
		textsCopy := make([]string, len(testTexts))
		copy(textsCopy, testTexts)
		for j := range textsCopy {
			tinystring.Convert(&textsCopy[j]).RemoveTilde().CamelCaseLower().Apply()
		}
		_ = textsCopy
	}
}

// Standard library benchmarks - traditional approach
func BenchmarkStandardLibWithoutPointer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make([]string, len(testTexts))
		for j, text := range testTexts {
			// Simulate equivalent operations using standard library
			processed := strings.ToLower(text)
			processed = strings.ReplaceAll(processed, "á", "a")
			processed = strings.ReplaceAll(processed, "é", "e")
			processed = strings.ReplaceAll(processed, "í", "i")
			processed = strings.ReplaceAll(processed, "ó", "o")
			processed = strings.ReplaceAll(processed, "ú", "u")
			// Convert to camelCase manually
			words := strings.Fields(processed)
			if len(words) > 0 {
				result := words[0]
				for k := 1; k < len(words); k++ {
					if len(words[k]) > 0 {
						result += strings.ToUpper(words[k][:1]) + words[k][1:]
					}
				}
				results[j] = result
			}
		}
		_ = results
	}
}

// Standard library benchmarks - pointer simulation (still creates new strings)
func BenchmarkStandardLibWithPointer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		textsCopy := make([]string, len(testTexts))
		copy(textsCopy, testTexts)
		for j := range textsCopy {
			// Standard library doesn't have true in-place operations
			processed := strings.ToLower(textsCopy[j])
			processed = strings.ReplaceAll(processed, "á", "a")
			processed = strings.ReplaceAll(processed, "é", "e")
			processed = strings.ReplaceAll(processed, "í", "i")
			processed = strings.ReplaceAll(processed, "ó", "o")
			processed = strings.ReplaceAll(processed, "ú", "u")
			// Convert to camelCase manually
			words := strings.Fields(processed)
			if len(words) > 0 {
				result := words[0]
				for k := 1; k < len(words); k++ {
					if len(words[k]) > 0 {
						result += strings.ToUpper(words[k][:1]) + words[k][1:]
					}
				}
				textsCopy[j] = result
			}
		}
		_ = textsCopy
	}
}

// Numeric processing comparison
func BenchmarkTinyStringNumeric(b *testing.B) {
	numbers := []float64{123.456, 789.012, 345.678}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make([]string, len(numbers))
		for j, num := range numbers {
			results[j] = tinystring.Convert(num).String()
		}
		_ = results
	}
}

func BenchmarkStandardLibNumeric(b *testing.B) {
	numbers := []float64{123.456, 789.012, 345.678}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make([]string, len(numbers))
		for j, num := range numbers {
			results[j] = strconv.FormatFloat(num, 'f', -1, 64)
		}
		_ = results
	}
}
