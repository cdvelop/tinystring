package main

import (
	"testing"

	"github.com/cdvelop/tinystring"
)

// Benchmark para comparar la eficiencia de uso de punteros vs. método tradicional
func BenchmarkStringWithoutPointer(b *testing.B) {
	texts := []string{
		"Él Múrcielago Rápido",
		"PROCESANDO textos LARGOS",
		"Optimización de MEMORIA",
		"Rendimiento en APLICACIONES",
		"Reducción de ASIGNACIONES",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make([]string, len(texts))
		for j, text := range texts {
			results[j] = tinystring.Convert(text).RemoveTilde().CamelCaseLower().String()
		}
		_ = results
	}
}

func BenchmarkStringWithPointer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		texts := []string{
			"Él Múrcielago Rápido",
			"PROCESANDO textos LARGOS",
			"Optimización de MEMORIA",
			"Rendimiento en APLICACIONES",
			"Reducción de ASIGNACIONES",
		}
		for j := range texts {
			tinystring.Convert(&texts[j]).RemoveTilde().CamelCaseLower().Apply()
		}
		_ = texts
	}
}

// Benchmark para escenarios de procesamiento masivo (simulado)
func BenchmarkMassiveProcessingWithoutPointer(b *testing.B) {
	// Simular procesamiento de una gran cantidad de textos
	texts := []string{
		"Él Múrcielago Rápido",
		"PROCESANDO textos LARGOS",
		"Optimización de MEMORIA",
		"Rendimiento en APLICACIONES",
		"Reducción de ASIGNACIONES",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := make([]string, len(texts))
		for j, text := range texts {
			results[j] = tinystring.Convert(text).RemoveTilde().CamelCaseLower().String()
		}
		_ = results
	}
}

func BenchmarkMassiveProcessingWithPointer(b *testing.B) {
	// Simular procesamiento de una gran cantidad de textos
	// pero modificando directamente los strings originales
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		texts := []string{
			"Él Múrcielago Rápido",
			"PROCESANDO textos LARGOS",
			"Optimización de MEMORIA",
			"Rendimiento en APLICACIONES",
			"Reducción de ASIGNACIONES",
		}
		for j := range texts {
			tinystring.Convert(&texts[j]).RemoveTilde().CamelCaseLower().Apply()
		}
		_ = texts
	}
}
