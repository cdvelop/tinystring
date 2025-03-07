package tinystring_test

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestStringOperations(t *testing.T) {
	t.Run("Replace", func(t *testing.T) {
		tests := []struct {
			input, old, newStr, expected string
		}{
			{"Este es un ejemplo de texto de prueba.", "ejemplo", "cambio", "Este es un cambio de texto de prueba."},
			{"Hola mundo!", "mundo", "Gophers", "Hola Gophers!"},
			{"abc abc abc", "abc", "123", "123 123 123"},
			{"abc", "xyz", "123", "abc"},
			{"", "", "123", ""},
			{"abcdabcdabcd", "cd", "12", "ab12ab12ab12"},
			{"palabra, punto,", ",", ".", "palabra. punto."},
		}

		for _, test := range tests {
			result := Convert(test.input).Replace(test.old, test.newStr).String()
			if result != test.expected {
				t.Errorf("Para input '%s', old '%s', new '%s', esperado '%s', pero obtenido '%s'", test.input, test.old, test.newStr, test.expected, result)
			}
		}
	})

	t.Run("TrimSuffix", func(t *testing.T) {
		tests := []struct {
			input, suffix, expected string
		}{
			{"hello.txt", ".txt", "hello"},
			{"example", "123", "example"},
			{"file.txt.txt", ".txt", "file.txt"},
			{"", "", ""},
			{"abc", "xyz", "abc"},
		}

		for _, test := range tests {
			result := Convert(test.input).TrimSuffix(test.suffix).String()
			if result != test.expected {
				t.Errorf("Para input '%s', suffix '%s', esperado '%s', pero obtenido '%s'", test.input, test.suffix, test.expected, result)
			}
		}
	})

	t.Run("Trim", func(t *testing.T) {
		tests := []struct {
			input, expected string
		}{
			{"  hello world  ", "hello world"},
			{"abc123", "abc123"},
			{"  trim me  ", "trim me"},
			{"", ""},
			{"  ", ""},
			{"    mucho espacio\n\n\t\tcon salto\n\n\t\tde linea     \n\t\t\t\t\t\t\n\t\t\t\t", "mucho espacio\n\n\t\tcon salto\n\n\t\tde linea"},
			{`    mucho espacio
		
		con salto

		de linea     
		              
		
		`, `mucho espacio
		
		con salto

		de linea`},
		}

		for _, test := range tests {
			result := Convert(test.input).Trim().String()
			if result != test.expected {
				t.Errorf("Para input '%s', esperado '%s', pero obtenido '%s'", test.input, test.expected, result)
			}
		}
	})

	t.Run("ChainableMethods", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
			chain    func(input string) string
		}{
			{
				name:     "Replace and Trim",
				input:    "  hello world  ",
				expected: "hello universe",
				chain: func(input string) string {
					return Convert(input).Trim().Replace("world", "universe").String()
				},
			},
			{
				name:     "Replace, Trim, and TrimSuffix",
				input:    "  filename.txt  ",
				expected: "file",
				chain: func(input string) string {
					return Convert(input).Trim().Replace("name", "").TrimSuffix(".txt").String()
				},
			},
			{
				name:     "Multiple Replaces",
				input:    "replace multiple words in this text",
				expected: "change many terms in this content",
				chain: func(input string) string {
					return Convert(input).
						Replace("replace", "change").
						Replace("multiple", "many").
						Replace("words", "terms").
						Replace("text", "content").
						String()
				},
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := test.chain(test.input)
				if result != test.expected {
					t.Errorf("Chain test '%s': expected '%s', got '%s'", test.name, test.expected, result)
				}
			})
		}
	})
}
