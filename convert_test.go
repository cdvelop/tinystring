package tinystring_test

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestConversions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		function func(*Text) *Text
	}{
		{
			name:     "Remove tildes",
			input:    "áéíóúÁÉÍÓÚ",
			want:     "aeiouAEIOU",
			function: (*Text).RemoveTilde,
		},
		{
			name:     "Remove tildes with mixed text",
			input:    "Hôlà Mündó",
			want:     "Hola Mundo",
			function: (*Text).RemoveTilde,
		},
		{
			name:  "CamelCaseLower",
			input: "hello world example",
			want:  "helloWorldExample",
			function: func(t *Text) *Text {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Convert to lower with tildes",
			input: "HÓLA MÚNDO",
			want:  "hola mundo",
			function: func(t *Text) *Text {
				return t.RemoveTilde().ToLower()
			},
		},
		{
			name:  "Convert to upper with tildes",
			input: "hóla múndo",
			want:  "HOLA MUNDO",
			function: func(t *Text) *Text {
				return t.RemoveTilde().ToUpper()
			},
		},
		{
			name:     "Special characters",
			input:    "ñÑàèìòùÀÈÌÒÙ",
			want:     "nNaeiouAEIOU",
			function: (*Text).RemoveTilde,
		},
		{
			name:  "Complete transformation",
			input: "Él Múrcielago Rápido",
			want:  "elMurcielagoRapido",
			function: func(t *Text) *Text {
				return t.RemoveTilde().CamelCaseLower()
			},
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
			function: func(t *Text) *Text {
				return t.RemoveTilde().ToLower().ToUpper().CamelCaseLower()
			},
		},
		{
			name:  "Single character",
			input: "A",
			want:  "a",
			function: func(t *Text) *Text {
				return t.ToLower()
			},
		},
		{
			name:  "Multiple spaces in camelCase",
			input: "hello    world    example",
			want:  "helloWorldExample",
			function: func(t *Text) *Text {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Non-mappable characters",
			input: "Hello! @#$%^ World 123",
			want:  "hello!@#$%^World123",
			function: func(t *Text) *Text {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Mixed transformations",
			input: "HÉLLÔ WórLD",
			want:  "HELLO WORLD",
			function: func(t *Text) *Text {
				return t.RemoveTilde().ToUpper()
			},
		},
		{
			name:  "CamelCase with accents",
			input: "él múrcielago RÁPIDO vuela",
			want:  "elMurcielagoRapidoVuela",
			function: func(t *Text) *Text {
				return t.RemoveTilde().CamelCaseLower()
			},
		},
		{
			name:  "CamelCaseLower",
			input: "hello world example",
			want:  "helloWorldExample",
			function: func(t *Text) *Text {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "CamelCaseUpper",
			input: "hello world example",
			want:  "HelloWorldExample",
			function: func(t *Text) *Text {
				return t.CamelCaseUpper()
			},
		},
		{
			name:  "ToSnakeCaseLower",
			input: "hello world example",
			want:  "hello_world_example",
			function: func(t *Text) *Text {
				return t.ToSnakeCaseLower()
			},
		},
		{
			name:  "Mixed case with numbers to CamelCaseLower",
			input: "User123Name",
			want:  "user123name",
			function: func(t *Text) *Text {
				return t.CamelCaseLower()
			},
		},
		{
			name:  "Mixed case with numbers to CamelCaseUpper",
			input: "User123Name",
			want:  "User123Name",
			function: func(t *Text) *Text {
				return t.CamelCaseUpper()
			},
		},
		{
			name:  "Mixed case with numbers to ToSnakeCaseLower",
			input: "User123Name",
			want:  "user123_name",
			function: func(t *Text) *Text {
				return t.ToSnakeCaseLower()
			},
		},
		{
			name:  "Accented text to camelCase",
			input: "Él Múrcielago Rápido",
			want:  "elMurcielagoRapido",
			function: func(t *Text) *Text {
				return t.RemoveTilde().CamelCaseLower()
			},
		},
		{
			name:  "Accented text to PascalCase",
			input: "Él Múrcielago Rápido",
			want:  "ElMurcielagoRapido",
			function: func(t *Text) *Text {
				return t.RemoveTilde().CamelCaseUpper()
			},
		},
		{
			name:  "Accented text to snake_case",
			input: "Él Múrcielago Rápido",
			want:  "el_murcielago_rapido",
			function: func(t *Text) *Text {
				return t.RemoveTilde().ToSnakeCaseLower()
			},
		},
		{
			name:  "Accented text to snake-case",
			input: "Él Múrcielago Rápido",
			want:  "el-murcielago-rapido",
			function: func(t *Text) *Text {
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
