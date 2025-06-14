package tinystring_test

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestStringOperations(t *testing.T) {
	t.Run("Replace", func(t *testing.T) {
		tests := []struct {
			input    string
			old      any
			newStr   any
			n        int
			expected string
		}{
			{"Este es un ejemplo de texto de prueba.", "ejemplo", "cambio", -1, "Este es un cambio de texto de prueba."},
			{"Hola mundo!", "mundo", "Gophers", -1, "Hola Gophers!"},
			{"abc abc abc", "abc", "123", -1, "123 123 123"},
			{"abc", "xyz", "123", -1, "abc"},
			{"", "", "123", -1, ""},
			{"abcdabcdabcd", "cd", "12", -1, "ab12ab12ab12"},
			{"palabra, punto,", ",", ".", -1, "palabra. punto."},
			// Pruebas con tipos diferentes de any
			{"Test 123 value", 123, 456, -1, "Test 456 value"},
			{"Boolean true in conv", true, false, -1, "Boolean false in conv"},
			{"Pi is 3.14159", 3.14159, 3.142, -1, "Pi is 3.142"},
			// Pruebas con límite de reemplazos
			{"abc abc abc", "abc", "123", 1, "123 abc abc"},
			{"abc abc abc", "abc", "123", 2, "123 123 abc"},
			{"abc abc abc", "abc", "123", 0, "abc abc abc"},
		}

		for _, test := range tests {
			var result string
			if test.n >= 0 {
				result = Convert(test.input).Replace(test.old, test.newStr, test.n).String()
			} else {
				result = Convert(test.input).Replace(test.old, test.newStr).String()
			}

			if result != test.expected {
				t.Errorf("Para input '%s', old '%v', new '%v', n '%d', esperado '%s', pero obtenido '%s'",
					test.input, test.old, test.newStr, test.n, test.expected, result)
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
			{"mi_directorio\\cmd", "\\cmd", "mi_directorio"},
		}

		for _, test := range tests {
			result := Convert(test.input).TrimSuffix(test.suffix).String()
			if result != test.expected {
				t.Errorf("Para input '%s', suffix '%s', esperado '%s', pero obtenido '%s'", test.input, test.suffix, test.expected, result)
			}
		}
	})

	t.Run("TrimPrefix", func(t *testing.T) {
		tests := []struct {
			input, prefix, expected string
		}{
			{"prefix-hello", "prefix-", "hello"},
			{"example", "123", "example"},
			{"txt.file", "txt.", "file"},
			{"", "", ""},
			{"abc", "xyz", "abc"},
		}

		for _, test := range tests {
			result := Convert(test.input).TrimPrefix(test.prefix).String()
			if result != test.expected {
				t.Errorf("Para input '%s', prefix '%s', esperado '%s', pero obtenido '%s'", test.input, test.prefix, test.expected, result)
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
			// Test case for JSON key trimming issue
			{"\n\t\t\"ID\"", "\"ID\""},
			{"\n\t\t\"Username\"", "\"Username\""},
			{"\t\t\"Email\"\n\t\t", "\"Email\""},
			{"    \"json_key\"    ", "\"json_key\""},
			{"\r\n\t \"nested_field\" \t\r\n", "\"nested_field\""},
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
				input:    "replace multiple words in this conv",
				expected: "change many terms in this content",
				chain: func(input string) string {
					return Convert(input).
						Replace("replace", "change").
						Replace("multiple", "many").
						Replace("words", "terms").
						Replace("conv", "content").
						String()
				},
			},
			{
				name:     "TrimPrefix and TrimSuffix",
				input:    "prefix-content.suffix",
				expected: "content",
				chain: func(input string) string {
					return Convert(input).TrimPrefix("prefix-").TrimSuffix(".suffix").String()
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

	t.Run("Split", func(t *testing.T) {
		tests := []struct {
			name      string
			input     string
			separator string
			expected  []string
		}{
			// Split by whitespace (no separator)
			{"Whitespace split", "hello world test", "", []string{"hello", "world", "test"}},
			{"Multiple spaces", "  hello   world  ", "", []string{"hello", "world"}},
			{"Tabs and newlines", "hello\tworld\ntest", "", []string{"hello", "world", "test"}},
			{"Empty string", "", "", []string{}},
			{"Single word", "hello", "", []string{"hello"}},
			
			// Split by custom separator
			{"Comma separator", "apple,banana,cherry", ",", []string{"apple", "banana", "cherry"}},
			{"Semicolon separator", "one;two;three", ";", []string{"one", "two", "three"}},
			{"Multi-char separator", "one::two::three", "::", []string{"one", "two", "three"}},
			{"No separator found", "hello world", ",", []string{"hello world"}},
			{"Empty separator characters", "hello", "", []string{"h", "e", "l", "l", "o"}},
			{"Short string", "hi", ",", []string{"hi"}}, // Less than 3 chars
			
			// Edge cases
			{"Separator at start", ",apple,banana", ",", []string{"", "apple", "banana"}},
			{"Separator at end", "apple,banana,", ",", []string{"apple", "banana", ""}},
			{"Consecutive separators", "apple,,banana", ",", []string{"apple", "", "banana"}},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var result []string
				if test.name == "Whitespace split" || test.name == "Multiple spaces" || 
				   test.name == "Tabs and newlines" || test.name == "Empty string" || 
				   test.name == "Single word" {
					// These tests should use Split without separator (whitespace split)
					result = Split(test.input)
				} else {
					// These tests should use Split with separator
					result = Split(test.input, test.separator)
				}
				
				if len(result) != len(test.expected) {
					t.Errorf("Length mismatch: expected %d, got %d", len(test.expected), len(result))
					return
				}
				
				for i, expected := range test.expected {
					if result[i] != expected {
						t.Errorf("Element %d: expected %q, got %q", i, expected, result[i])
					}
				}
			})
		}
	})

	t.Run("Join", func(t *testing.T) {
		tests := []struct {
			name      string
			input     []string
			separator string
			expected  string
		}{
			{"Default separator", []string{"hello", "world"}, "", "hello world"},
			{"Custom separator", []string{"apple", "banana", "cherry"}, ",", "apple,banana,cherry"},
			{"Dash separator", []string{"one", "two", "three"}, "-", "one-two-three"},
			{"Empty slice", []string{}, ",", ""},
			{"Single element", []string{"hello"}, ",", "hello"},
			{"Empty strings", []string{"", "", ""}, ",", ",,"},
			{"Multi-char separator", []string{"a", "b", "c"}, "::", "a::b::c"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var result string
				if test.separator == "" {
					result = Convert(test.input).Join().String()
				} else {
					result = Convert(test.input).Join(test.separator).String()
				}
				
				if result != test.expected {
					t.Errorf("Expected %q, got %q", test.expected, result)
				}
			})
		}
	})

	t.Run("Repeat", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			count    int
			expected string
		}{
			{"Basic repeat", "abc", 3, "abcabcabc"},
			{"Single char", "x", 5, "xxxxx"},
			{"Zero count", "hello", 0, ""},
			{"Negative count", "test", -1, ""},
			{"Empty string", "", 5, ""},
			{"Count 1", "hello", 1, "hello"},
			{"Large count", "a", 10, "aaaaaaaaaa"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Convert(test.input).Repeat(test.count).String()
				if result != test.expected {
					t.Errorf("Expected %q, got %q", test.expected, result)
				}
			})
		}
	})

	t.Run("TrimEdgeCases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			// Test different whitespace combinations
			{"Only spaces", "     ", ""},
			{"Only tabs", "\t\t\t", ""},
			{"Only newlines", "\n\n\n", ""},
			{"Mixed whitespace", " \t\n\r ", ""},
			{"Carriage return", "\r\nhello\r\n", "hello"},
			{"Unicode whitespace mixed", "  \t hello world \n\r  ", "hello world"},
			{"Single space", " ", ""},
			{"Content with internal whitespace", "  hello   world  ", "hello   world"},
			
			// JSON-specific trim cases
			{"JSON field name", "    \"fieldName\"    ", "\"fieldName\""},
			{"JSON value", "\t\t\"some value\"\n\n", "\"some value\""},
			{"Complex JSON key", "\r\n    \"user_profile\"  \t\r", "\"user_profile\""},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := Convert(test.input).Trim().String()
				if result != test.expected {
					t.Errorf("Input: %q, Expected: %q, Got: %q", test.input, test.expected, result)
				}
			})
		}
	})

	t.Run("ReplaceEdgeCases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			old      any
			newStr   any
			n        int
			expected string
		}{
			// Edge cases for replace
			{"Empty old string", "hello world", "", "X", -1, "hello world"},
			{"Empty input", "", "old", "new", -1, ""},
			{"Replace with empty", "hello", "hello", "", -1, ""},
			{"No match", "hello", "xyz", "abc", -1, "hello"},
			{"Replace at start", "hello world", "hello", "hi", -1, "hi world"},
			{"Replace at end", "hello world", "world", "earth", -1, "hello earth"},
			{"Overlapping matches", "aaa", "aa", "b", -1, "ba"},
			
			// Type conversion tests
			{"Number to string", "value: 42", 42, "forty-two", -1, "value: forty-two"},
			{"Bool to string", "enabled: true", true, "yes", -1, "enabled: yes"},
			{"Float to string", "pi: 3.14", 3.14, "π", -1, "pi: π"},
			
			// Limit tests
			{"Zero limit", "abc abc abc", "abc", "123", 0, "abc abc abc"},
			{"Partial replace", "the the the", "the", "a", 2, "a a the"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var result string
				if test.n >= 0 {
					result = Convert(test.input).Replace(test.old, test.newStr, test.n).String()
				} else {
					result = Convert(test.input).Replace(test.old, test.newStr).String()
				}

				if result != test.expected {
					t.Errorf("Input: %q, Expected: %q, Got: %q", test.input, test.expected, result)
				}
			})
		}
	})

	t.Run("TrimPrefixSuffixEdgeCases", func(t *testing.T) {
		tests := []struct {
			name     string
			function string // "prefix" or "suffix"
			input    string
			remove   string
			expected string
		}{
			// TrimPrefix edge cases
			{"Prefix - empty input", "prefix", "", "pre", ""},
			{"Prefix - empty remove", "prefix", "hello", "", "hello"},
			{"Prefix - no match", "prefix", "hello", "world", "hello"},
			{"Prefix - partial match", "prefix", "hello", "hel", "lo"},
			{"Prefix - exact match", "prefix", "hello", "hello", ""},
			{"Prefix - longer remove", "prefix", "hi", "hello", "hi"},
			{"Prefix - case sensitive", "prefix", "Hello", "hello", "Hello"},
			
			// TrimSuffix edge cases
			{"Suffix - empty input", "suffix", "", "ing", ""},
			{"Suffix - empty remove", "suffix", "hello", "", "hello"},
			{"Suffix - no match", "suffix", "hello", "world", "hello"},
			{"Suffix - partial match", "suffix", "hello", "llo", "he"},
			{"Suffix - exact match", "suffix", "hello", "hello", ""},
			{"Suffix - longer remove", "suffix", "hi", "hello", "hi"},
			{"Suffix - case sensitive", "suffix", "Hello", "hello", "Hello"},
			{"Suffix - multiple occurrences", "suffix", "test.txt.txt", ".txt", "test.txt"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var result string
				if test.function == "prefix" {
					result = Convert(test.input).TrimPrefix(test.remove).String()
				} else {
					result = Convert(test.input).TrimSuffix(test.remove).String()
				}

				if result != test.expected {
					t.Errorf("Input: %q, Remove: %q, Expected: %q, Got: %q", 
						test.input, test.remove, test.expected, result)
				}
			})
		}
	})

	t.Run("MethodChaining", func(t *testing.T) {
		// Test complex method chaining scenarios
		tests := []struct {
			name     string
			input    string
			expected string
			chain    func(input string) string
		}{
			{
				name:     "Complex cleaning",
				input:    "  /path/to/file.txt.backup  ",
				expected: "file",
				chain: func(input string) string {
					return Convert(input).
						Trim().
						TrimPrefix("/path/to/").
						TrimSuffix(".backup").
						TrimSuffix(".txt").
						String()
				},
			},
			{
				name:     "Text processing",
				input:    "Hello_World_Test",
				expected: "Hello World (processed)",
				chain: func(input string) string {
					return Convert(input).
						Replace("_", " ").
						Replace("Test", "(processed)").
						String()
				},
			},
			{
				name:     "Multiple operations",
				input:    "   prefix-data-suffix.old   ",
				expected: "DATA",
				chain: func(input string) string {
					return Convert(input).
						Trim().
						TrimPrefix("prefix-").
						TrimSuffix(".old").
						Replace("data", "DATA").
						TrimSuffix("-suffix").
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

	// Existing ChainableMethods test...
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
				input:    "replace multiple words in this conv",
				expected: "change many terms in this content",
				chain: func(input string) string {
					return Convert(input).
						Replace("replace", "change").
						Replace("multiple", "many").
						Replace("words", "terms").
						Replace("conv", "content").
						String()
				},
			},
			{
				name:     "TrimPrefix and TrimSuffix",
				input:    "prefix-content.suffix",
				expected: "content",
				chain: func(input string) string {
					return Convert(input).TrimPrefix("prefix-").TrimSuffix(".suffix").String()
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
