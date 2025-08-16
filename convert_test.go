package tinystring

import (
	"testing"
)

func TestConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		function func(*Conv) *Conv
	}{
		{
			name:     "Tilde does not remove Ñ/ñ",
			input:    "Ñandú ñandú",
			want:     "Ñandu ñandu",
			function: (*Conv).Tilde,
		},
		{
			name:     "Remove tildes",
			input:    "áéíóúÁÉÍÓÚ",
			want:     "aeiouAEIOU",
			function: (*Conv).Tilde,
		},
		{
			name:     "Remove tildes with mixed Conv",
			input:    "Hôlà Mündó",
			want:     "Hola Mundo",
			function: (*Conv).Tilde,
		},
		{
			name:  "CamelLow",
			input: "hello world example",
			want:  "helloWorldExample",
			function: func(t *Conv) *Conv {
				return t.CamelLow()
			},
		},
		{
			name:  "Convert to lower with tildes",
			input: "HÓLA MÚNDO",
			want:  "hola mundo",
			function: func(t *Conv) *Conv {
				return t.Tilde().ToLower()
			},
		},
		{
			name:  "Convert to upper with tildes",
			input: "hóla múndo",
			want:  "HOLA MUNDO",
			function: func(t *Conv) *Conv {
				return t.Tilde().ToUpper()
			},
		},
		{
			name:     "Special characters",
			input:    "ñÑàèìòùÀÈÌÒÙ",
			want:     "ñÑaeiouAEIOU",
			function: (*Conv).Tilde,
		},
		{
			name:     "Tilde does not remove Ñ/ñ",
			input:    "Ñandú ñandú",
			want:     "Ñandu ñandu",
			function: (*Conv).Tilde,
		},
		{
			name:  "Complete transformation",
			input: "Él Múrcielago Rápido",
			want:  "elMurcielagoRapido",
			function: func(t *Conv) *Conv {
				return t.Tilde().CamelLow()
			},
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
			function: func(t *Conv) *Conv {
				return t.Tilde().ToLower().ToUpper().CamelLow()
			},
		},
		{
			name:  "Single character",
			input: "A",
			want:  "a",
			function: func(t *Conv) *Conv {
				return t.ToLower()
			},
		},
		{
			name:  "Multiple spaces in camelCase",
			input: "hello    world    example",
			want:  "helloWorldExample",
			function: func(t *Conv) *Conv {
				return t.CamelLow()
			},
		},
		{
			name:  "Non-mappable characters",
			input: "Hello! @#$%^ World 123",
			want:  "hello!@#$%^World123",
			function: func(t *Conv) *Conv {
				return t.CamelLow()
			},
		},
		{
			name:  "Mixed transformations",
			input: "HÉLLÔ WórLD",
			want:  "HELLO WORLD",
			function: func(t *Conv) *Conv {
				return t.Tilde().ToUpper()
			},
		},
		{
			name:  "CamelCase with accents",
			input: "él múrcielago RÁPIDO vuela",
			want:  "elMurcielagoRapidoVuela",
			function: func(t *Conv) *Conv {
				return t.Tilde().CamelLow()
			},
		},
		{
			name:  "CamelLow",
			input: "hello world example",
			want:  "helloWorldExample",
			function: func(t *Conv) *Conv {
				return t.CamelLow()
			},
		},
		{
			name:  "CamelUp",
			input: "hello world example",
			want:  "HelloWorldExample",
			function: func(t *Conv) *Conv {
				return t.CamelUp()
			},
		},
		{
			name:  "SnakeLow",
			input: "hello world example",
			want:  "hello_world_example",
			function: func(t *Conv) *Conv {
				return t.SnakeLow()
			},
		},
		{
			name:  "Mixed case with numbers to CamelLow",
			input: "User123Name",
			want:  "user123name",
			function: func(t *Conv) *Conv {
				return t.CamelLow()
			},
		},
		{
			name:  "Mixed case with numbers to CamelUp",
			input: "User123Name",
			want:  "User123Name",
			function: func(t *Conv) *Conv {
				return t.CamelUp()
			},
		},
		{
			name:  "Mixed case with numbers to SnakeLow",
			input: "User123Name",
			want:  "user123_name",
			function: func(t *Conv) *Conv {
				return t.SnakeLow()
			},
		},
		{
			name:  "Accented Conv to camelCase",
			input: "Él Múrcielago Rápido",
			want:  "elMurcielagoRapido",
			function: func(t *Conv) *Conv {
				return t.Tilde().CamelLow()
			},
		},
		{
			name:  "Accented Conv to PascalCase",
			input: "Él Múrcielago Rápido",
			want:  "ElMurcielagoRapido",
			function: func(t *Conv) *Conv {
				return t.Tilde().CamelUp()
			},
		},
		{
			name:  "Accented Conv to snake_case",
			input: "Él Múrcielago Rápido",
			want:  "el_murcielago_rapido",
			function: func(t *Conv) *Conv {
				return t.Tilde().SnakeLow()
			},
		},
		{
			name:  "Accented Conv to snake-case",
			input: "Él Múrcielago Rápido",
			want:  "el-murcielago-rapido",
			function: func(t *Conv) *Conv {
				return t.Tilde().SnakeLow("-")
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
			out := Convert(tc.input).String()
			if out != tc.expected {
				t.Errorf("Convert(%v).String() = %q, want %q", tc.input, out, tc.expected)
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
