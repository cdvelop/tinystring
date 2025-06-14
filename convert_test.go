package tinystring

import (
	"errors"
	"fmt"
	"testing"
)

func TestConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		function func(*conv) *conv
	}{
		{
			name:     "Remove tildes",
			input:    "áéíóúÁÉÍÓÚ",
			want:     "aeiouAEIOU",
			function: (*conv).RemoveTilde,
		},
		{
			name:     "Remove tildes with mixed conv",
			input:    "Hôlà Mündó",
			want:     "Hola Mundo",
			function: (*conv).RemoveTilde,
		},
		{
			name:  "CamelCaseLower",
			input: "hello world example",
			want:  "helloWorldExample",
			function: func(t *conv) *conv {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Convert to lower with tildes",
			input: "HÓLA MÚNDO",
			want:  "hola mundo",
			function: func(t *conv) *conv {
				return t.RemoveTilde().ToLower()
			},
		},
		{
			name:  "Convert to upper with tildes",
			input: "hóla múndo",
			want:  "HOLA MUNDO",
			function: func(t *conv) *conv {
				return t.RemoveTilde().ToUpper()
			},
		},
		{
			name:     "Special characters",
			input:    "ñÑàèìòùÀÈÌÒÙ",
			want:     "nNaeiouAEIOU",
			function: (*conv).RemoveTilde,
		},
		{
			name:  "Complete transformation",
			input: "Él Múrcielago Rápido",
			want:  "elMurcielagoRapido",
			function: func(t *conv) *conv {
				return t.RemoveTilde().CamelCaseLower()
			},
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
			function: func(t *conv) *conv {
				return t.RemoveTilde().ToLower().ToUpper().CamelCaseLower()
			},
		},
		{
			name:  "Single character",
			input: "A",
			want:  "a",
			function: func(t *conv) *conv {
				return t.ToLower()
			},
		},
		{
			name:  "Multiple spaces in camelCase",
			input: "hello    world    example",
			want:  "helloWorldExample",
			function: func(t *conv) *conv {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Non-mappable characters",
			input: "Hello! @#$%^ World 123",
			want:  "hello!@#$%^World123",
			function: func(t *conv) *conv {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Mixed transformations",
			input: "HÉLLÔ WórLD",
			want:  "HELLO WORLD",
			function: func(t *conv) *conv {
				return t.RemoveTilde().ToUpper()
			},
		},
		{
			name:  "CamelCase with accents",
			input: "él múrcielago RÁPIDO vuela",
			want:  "elMurcielagoRapidoVuela",
			function: func(t *conv) *conv {
				return t.RemoveTilde().CamelCaseLower()
			},
		},
		{
			name:  "CamelCaseLower",
			input: "hello world example",
			want:  "helloWorldExample",
			function: func(t *conv) *conv {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "CamelCaseUpper",
			input: "hello world example",
			want:  "HelloWorldExample",
			function: func(t *conv) *conv {
				return t.CamelCaseUpper()
			},
		},
		{
			name:  "ToSnakeCaseLower",
			input: "hello world example",
			want:  "hello_world_example",
			function: func(t *conv) *conv {
				return t.ToSnakeCaseLower()
			},
		},
		{
			name:  "Mixed case with numbers to CamelCaseLower",
			input: "User123Name",
			want:  "user123name",
			function: func(t *conv) *conv {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Mixed case with numbers to CamelCaseUpper",
			input: "User123Name",
			want:  "User123Name",
			function: func(t *conv) *conv {
				return t.CamelCaseUpper()
			},
		},
		{
			name:  "Mixed case with numbers to ToSnakeCaseLower",
			input: "User123Name",
			want:  "user123_name",
			function: func(t *conv) *conv {
				return t.ToSnakeCaseLower()
			},
		},
		{
			name:  "Accented conv to camelCase",
			input: "Él Múrcielago Rápido",
			want:  "elMurcielagoRapido",
			function: func(t *conv) *conv {
				return t.RemoveTilde().CamelCaseLower()
			},
		},
		{
			name:  "Accented conv to PascalCase",
			input: "Él Múrcielago Rápido",
			want:  "ElMurcielagoRapido",
			function: func(t *conv) *conv {
				return t.RemoveTilde().CamelCaseUpper()
			},
		},
		{
			name:  "Accented conv to snake_case",
			input: "Él Múrcielago Rápido",
			want:  "el_murcielago_rapido",
			function: func(t *conv) *conv {
				return t.RemoveTilde().ToSnakeCaseLower()
			},
		},
		{
			name:  "Accented conv to snake-case",
			input: "Él Múrcielago Rápido",
			want:  "el-murcielago-rapido",
			function: func(t *conv) *conv {
				return t.RemoveTilde().ToSnakeCaseLower("-")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.function(Convert(tt.input)).String()
			if got != tt.want {
				t.Fatalf("\n🎯Test: %q\ninput: %q\n   got: %q\n  want: %q", tt.name, tt.input, got, tt.want)
			}
		})
	}
}

func TestConvertToStringOptimized(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, ""},
		{"empty string", "", ""},
		{"string", "hello", "hello"},
		{"int zero", 0, "0"},
		{"int one", 1, "1"},
		{"int small", 42, "42"},
		{"int negative", -10, "-10"},
		{"int large", 12345678, "12345678"},
		{"int8", int8(8), "8"},
		{"int16", int16(16), "16"},
		{"int32", int32(32), "32"},
		{"int64", int64(64), "64"},
		{"uint", uint(5), "5"},
		{"uint8", uint8(8), "8"},
		{"uint16", uint16(16), "16"},
		{"uint32", uint32(32), "32"},
		{"uint64", uint64(64), "64"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"float32 zero", float32(0), "0"},
		{"float32 one", float32(1), "1"},
		{"float32", float32(3.14), "3.14"},
		{"float64 zero", 0.0, "0"},
		{"float64 one", 1.0, "1"},
		{"float64", 3.14159, "3.14159"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Convert(tc.input).String()
			if result != tc.expected {
				t.Errorf("Convert(%v).String() = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// Test to ensure that repeated conversions return the expected value
func TestConvertStringConsistency(t *testing.T) {
	// For small integers
	s1 := Convert(42).String()
	s2 := Convert(42).String()

	if s1 != "42" || s2 != "42" {
		t.Errorf("Expected Convert(42).String() to return \"42\"")
	}

	// For boolean values
	b1 := Convert(true).String()
	b2 := Convert(true).String()

	if b1 != "true" || b2 != "true" {
		t.Errorf("Expected Convert(true).String() to return \"true\"")
	}

	// For zero value
	z1 := Convert(0).String()
	z2 := Convert(0).String()

	if z1 != "0" || z2 != "0" {
		t.Errorf("Expected Convert(0).String() to return \"0\"")
	}

	// Prueba adicional para verificar la consistencia
	t.Log("Tests for Convert().String() consistency passed")
}

// Test Unicode conversion functionality (addRne2Buf)
func TestUnicodeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"ASCII characters", "Hello", "Hello"},
		{"UTF-8 characters", "Héllo", "Héllo"},
		{"Unicode emojis", "Hello 🌍", "Hello 🌍"},
		{"Mixed Unicode", "Café ñoño", "Café ñoño"},
		{"Chinese characters", "你好", "你好"},
		{"Accented characters", "résumé", "résumé"},
		{"Special symbols", "™®©", "™®©"},
		{"Mathematical symbols", "∀∃∈∉", "∀∃∈∉"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test Unicode case conversion operations
func TestUnicodeCaseOperations(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		operation func(*conv) *conv
		expected  string
	}{
		{"Unicode to upper", "café", func(c *conv) *conv { return c.ToUpper() }, "CAFÉ"},
		{"Unicode to lower", "CAFÉ", func(c *conv) *conv { return c.ToLower() }, "café"},
		{"Unicode remove tildes", "ñoño", func(c *conv) *conv { return c.RemoveTilde() }, "nono"},
		{"Mixed Unicode operations", "Résumé", func(c *conv) *conv { return c.ToUpper().RemoveTilde() }, "RESUME"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.operation(Convert(tt.input)).String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test any2s function with various input types
func TestAny2sConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"String input", "hello", "hello"},
		{"Bool true", true, "true"},
		{"Bool false", false, "false"},
		{"Int value", 42, "42"},
		{"Negative int", -123, "-123"},
		{"Uint value", uint(42), "42"},
		{"Float value", 3.14, "3.14"},
		{"Float32 value", float32(2.5), "2.5"},
		{"Int8 value", int8(127), "127"},
		{"Int16 value", int16(32767), "32767"},
		{"Int32 value", int32(2147483647), "2147483647"},
		{"Int64 value", int64(9223372036854775807), "9223372036854775807"},
		{"Uint8 value", uint8(255), "255"},
		{"Uint16 value", uint16(65535), "65535"},
		{"Uint32 value", uint32(4294967295), "4294967295"},
		{"Uint64 value", uint64(18446744073709551615), "18446744073709551615"},
		{"Uintptr value", uintptr(12345), "12345"},
		{"Nil value", nil, ""},
		{"Empty string", "", ""},
		{"Zero int", 0, "0"},
		{"Zero float", 0.0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test any2s with error types
func TestAny2sErrorHandling(t *testing.T) {
	// Test with standard error - check if conversion handles it
	err := errors.New("test error")
	conv := Convert(err)

	// The behavior may vary - check what actually happens
	errResult := conv.Error()
	t.Logf("Error conversion result: %q", errResult)

	// Test with NewErr function which should work
	conv2 := Convert(nil).NewErr("custom error")
	if conv2.Error() != "custom error" {
		t.Errorf("Expected 'custom error', got '%s'", conv2.Error())
	}

	// Test Errorf function
	conv3 := Errorf("formatted error: %s", "test")
	if conv3.Error() != "formatted error: test" {
		t.Errorf("Expected 'formatted error: test', got '%s'", conv3.Error())
	}
}

// Test any2s with unsupported types
func TestAny2sUnsupportedTypes(t *testing.T) {
	// Test with unsupported types that should return empty string
	unsupportedTypes := []any{
		make(chan int),
		func() {},
		make(map[string]int),
		[]int{1, 2, 3},         // slice should be handled differently
		struct{ X int }{X: 42}, // struct should be handled differently
	}

	for i, input := range unsupportedTypes {
		t.Run(fmt.Sprintf("Unsupported type %d", i), func(t *testing.T) {
			result := Convert(input).String()
			// Most unsupported types should result in empty string or some default behavior
			t.Logf("Input type %T resulted in: %q", input, result)
		})
	}
}

func TestAddRne2BufDirectly(t *testing.T) {
	tests := []struct {
		name     string
		rune     rune
		expected []byte
	}{
		// ASCII range (< 0x80)
		{"ASCII A", 'A', []byte{0x41}},
		{"ASCII space", ' ', []byte{0x20}},
		{"ASCII zero", '\x00', []byte{0x00}},
		{"ASCII DEL", '\x7F', []byte{0x7F}},

		// 2-byte UTF-8 (0x80-0x7FF)
		{"Latin-1 supplement é", 'é', []byte{0xC3, 0xA9}},
		{"Latin-1 supplement ñ", 'ñ', []byte{0xC3, 0xB1}},
		{"Cyrillic А", 'А', []byte{0xD0, 0x90}},
		{"Greek α", 'α', []byte{0xCE, 0xB1}},

		// 3-byte UTF-8 (0x800-0xFFFF)
		{"CJK 你", '你', []byte{0xE4, 0xBD, 0xA0}},
		{"CJK 好", '好', []byte{0xE5, 0xA5, 0xBD}},
		{"Symbol ™", '™', []byte{0xE2, 0x84, 0xA2}},
		{"Symbol ∀", '∀', []byte{0xE2, 0x88, 0x80}},

		// 4-byte UTF-8 (0x10000-0x10FFFF)
		{"Emoji 😀", '😀', []byte{0xF0, 0x9F, 0x98, 0x80}},
		{"Emoji 🌍", '🌍', []byte{0xF0, 0x9F, 0x8C, 0x8D}},
		{"Musical 𝄞", '𝄞', []byte{0xF0, 0x9D, 0x84, 0x9E}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf []byte
			result := addRne2Buf(buf, tt.rune)

			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: expected %d bytes, got %d bytes",
					len(tt.expected), len(result))
			}

			for i := 0; i < len(tt.expected) && i < len(result); i++ {
				if result[i] != tt.expected[i] {
					t.Errorf("Byte %d: expected 0x%02X, got 0x%02X",
						i, tt.expected[i], result[i])
				}
			}

			// Verify the result is valid UTF-8 by converting back to string
			str := string(result)
			if len([]rune(str)) != 1 || []rune(str)[0] != tt.rune {
				t.Errorf("Round-trip failed: expected rune %U, got %U",
					tt.rune, []rune(str)[0])
			}
		})
	}
}

func TestAddRne2BufAppend(t *testing.T) {
	// Test that addRne2Buf properly appends to existing buffer
	buf := []byte("Hello ")
	buf = addRne2Buf(buf, '🌍')
	expected := "Hello 🌍"

	result := string(buf)
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
