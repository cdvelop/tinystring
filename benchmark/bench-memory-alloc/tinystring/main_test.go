package main

import (
	"testing"
)

func BenchmarkStringProcessing(b *testing.B) {
	testTexts := []string{
		"Él Múrcielago Rápido",
		"PROCESANDO textos LARGOS",
		"Optimización de MEMORIA",
		"Rendimiento en APLICACIONES",
		"Reducción de ASIGNACIONES",
		"Análisis de RENDIMIENTO",
		"Gestión de RECURSOS",
		"Eficiencia OPERACIONAL",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processTextWithTinyString(testTexts)
	}
}

func BenchmarkStringProcessingWithPointers(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testTexts := []string{
			"Él Múrcielago Rápido",
			"PROCESANDO textos LARGOS",
			"Optimización de MEMORIA",
			"Rendimiento en APLICACIONES",
			"Reducción de ASIGNACIONES",
			"Análisis de RENDIMIENTO",
			"Gestión de RECURSOS",
			"Eficiencia OPERACIONAL",
		}
		processTextWithTinyStringPointers(testTexts)
		_ = testTexts
	}
}

func BenchmarkNumberProcessing(b *testing.B) {
	testNumbers := []float64{
		123456.789,
		987654.321,
		555555.555,
		111111.111,
		999999.999,
		777777.777,
		333333.333,
		888888.888,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processNumbersWithTinyString(testNumbers)
	}
}

func BenchmarkMixedOperations(b *testing.B) {
	testData := map[string]interface{}{
		"Número": 12345.67,
		"Texto":  "Información IMPORTANTE",
		"Valor":  98765.43,
		"Título": "Análisis de RENDIMIENTO",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make(map[string]string)
		for key, value := range testData {
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
