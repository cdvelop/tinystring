package tinystring

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		maxWidth      int
		reservedChars int
		expected      string
	}{
		{
			name:          "Text shorter than max width",
			input:         "Hello",
			maxWidth:      10,
			reservedChars: 0,
			expected:      "Hello     ",
		},
		{
			name:          "Text longer than max width with ellipsis",
			input:         "Hello, World!",
			maxWidth:      10,
			reservedChars: 0,
			expected:      "Hello, ...",
		},
		{
			name:          "Text longer with reserved chars",
			input:         "Hello, World!",
			maxWidth:      10,
			reservedChars: 3,
			expected:      "Hell...",
		},
		{
			name:          "Very short max width",
			input:         "Hello",
			maxWidth:      2,
			reservedChars: 0,
			expected:      "He",
		},
		{
			name:          "maxWith is zero",
			input:         "Test",
			maxWidth:      0,
			reservedChars: 0,
			expected:      "Test",
		},
		{
			name:          "Empty string",
			input:         "",
			maxWidth:      5,
			reservedChars: 0,
			expected:      "     ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).Truncate(tt.maxWidth, tt.reservedChars).String()
			if result != tt.expected {
				t.Errorf("Convert(%q).Truncate(%d, %d) = %q, want %q",
					tt.input, tt.maxWidth, tt.reservedChars, result, tt.expected)
			}
		})
	}
}

// TestTruncateChain tests the chaining capabilities of the Truncate method with other methods
func TestTruncateChain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		function func(*Text) *Text
	}{
		{
			name:  "Uppercase and truncate",
			input: "hello world",
			want:  "HELLO W...",
			function: func(t *Text) *Text {
				return t.ToUpper().Truncate(10, 0)
			},
		},
		{
			name:  "Lowercase and truncate",
			input: "HELLO WORLD",
			want:  "hello...",
			function: func(t *Text) *Text {
				return t.ToLower().Truncate(8, 0)
			},
		},
		{
			name:  "Remove tilde and truncate",
			input: "Ñandú está corriendo",
			want:  "Nandu esta ...",
			function: func(t *Text) *Text {
				return t.RemoveTilde().Truncate(14, 0)
			},
		},
		{
			name:  "CamelCase and truncate",
			input: "hello world example",
			want:  "helloWorld...",
			function: func(t *Text) *Text {
				return t.CamelCaseLower().Truncate(13, 0)
			},
		},
		{
			name:  "Truncate and repeat",
			input: "hello",
			want:  "hello hello ",
			function: func(t *Text) *Text {
				return t.Truncate(6, 0).Repeat(2)
			},
		},
		{
			name:  "SnakeCase and truncate",
			input: "Hello World Example",
			want:  "hello_world_...",
			function: func(t *Text) *Text {
				return t.ToSnakeCaseLower().Truncate(15, 0)
			},
		},
		{
			name:  "Truncate with custom separator",
			input: "Hello World Example",
			want:  "hello-world-ex...",
			function: func(t *Text) *Text {
				return t.ToSnakeCaseLower("-").Truncate(17, 0)
			},
		},
		{
			name:  "Complex chaining with truncate",
			input: "Él Múrcielago Rápido",
			want:  "ELMURC...",
			function: func(t *Text) *Text {
				return t.RemoveTilde().CamelCaseLower().ToUpper().Truncate(9, 0)
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
