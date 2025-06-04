package tinystring

import (
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
