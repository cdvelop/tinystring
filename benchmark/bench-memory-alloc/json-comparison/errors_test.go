package jsoncomparison

import (
	"encoding/json"
	"testing"

	"github.com/cdvelop/tinystring"
)

// Benchmarks para casos de error en Marshal

func BenchmarkJsonMarshalErrors_Standard(b *testing.B) {
	// Crear un tipo que cause error al marshalling
	ch := make(chan int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(ch)
	}
}

func BenchmarkJsonMarshalErrors_TinyString(b *testing.B) {
	// Crear un tipo que cause error al marshalling
	ch := make(chan int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tinystring.Convert(ch).JsonEncode()
	}
}

// Benchmarks para casos de error en Unmarshal

func BenchmarkJsonUnmarshalErrors_Standard(b *testing.B) {
	var result ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, invalidJSON := range invalidData {
			json.Unmarshal([]byte(invalidJSON), &result)
		}
	}
}

func BenchmarkJsonUnmarshalErrors_TinyString(b *testing.B) {
	var result ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, invalidJSON := range invalidData {
			tinystring.Convert(invalidJSON).JsonDecode(&result)
		}
	}
}
