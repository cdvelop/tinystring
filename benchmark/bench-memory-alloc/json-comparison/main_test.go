package jsoncomparison

import (
	"encoding/json"
	"testing"

	"github.com/cdvelop/tinystring"
)

var (
	singleUser  = generateTestData(1)[0]
	batch100    = generateTestData(100)
	batch1000   = generateTestData(1000)
	batch10000  = generateTestData(10000)
	invalidData = generateInvalidData()
)

// Benchmarks para Marshal (encoding)

func BenchmarkJsonMarshalSingle_Standard(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&singleUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshalSingle_TinyString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tinystring.Convert(&singleUser).JsonEncode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshalBatch100_Standard(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&batch100)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshalBatch100_TinyString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tinystring.Convert(&batch100).JsonEncode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshalBatch1000_Standard(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&batch1000)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshalBatch1000_TinyString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tinystring.Convert(&batch1000).JsonEncode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshalBatch10000_Standard(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&batch10000)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonMarshalBatch10000_TinyString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tinystring.Convert(&batch10000).JsonEncode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmarks para Unmarshal (decoding)

func BenchmarkJsonUnmarshalSingle_Standard(b *testing.B) {
	data, _ := json.Marshal(&singleUser)
	var result ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalSingle_TinyString(b *testing.B) {
	data, err := json.Marshal(&singleUser) // Usar json.Marshal para generar el JSON
	if err != nil {
		b.Fatal(err)
	}
	var result ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := tinystring.Convert(string(data)).JsonDecode(&result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalBatch100_Standard(b *testing.B) {
	data, _ := json.Marshal(&batch100)
	var result []ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalBatch100_TinyString(b *testing.B) {
	data, err := json.Marshal(&batch100)
	if err != nil {
		b.Fatal(err)
	}
	var result []ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := tinystring.Convert(string(data)).JsonDecode(&result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalBatch1000_Standard(b *testing.B) {
	data, _ := json.Marshal(&batch1000)
	var result []ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalBatch1000_TinyString(b *testing.B) {
	data, err := json.Marshal(&batch1000)
	if err != nil {
		b.Fatal(err)
	}
	var result []ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := tinystring.Convert(string(data)).JsonDecode(&result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalBatch10000_Standard(b *testing.B) {
	data, _ := json.Marshal(&batch10000)
	var result []ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsonUnmarshalBatch10000_TinyString(b *testing.B) {
	data, err := json.Marshal(&batch10000)
	if err != nil {
		b.Fatal(err)
	}
	var result []ComplexUser
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := tinystring.Convert(string(data)).JsonDecode(&result)
		if err != nil {
			b.Fatal(err)
		}
	}
}
