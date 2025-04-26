package tinystring

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		maxWidth      any
		reservedChars any
		useReserved   bool
		expected      string
	}{
		{
			name:        "Text shorter than max width",
			input:       "Hello",
			maxWidth:    10,
			useReserved: false,
			expected:    "Hello     ",
		},
		{
			name:        "Text longer than max width with ellipsis",
			input:       "Hello, World!",
			maxWidth:    10,
			useReserved: false,
			expected:    "Hello, ...",
		},
		{
			name:          "Text longer with reserved chars",
			input:         "Hello, World!",
			maxWidth:      10,
			reservedChars: 3,
			useReserved:   true,
			expected:      "Hell...",
		},
		{
			name:        "Very short max width",
			input:       "Hello",
			maxWidth:    2,
			useReserved: false,
			expected:    "He",
		},
		{
			name:        "maxWith is zero",
			input:       "Test",
			maxWidth:    0,
			useReserved: false,
			expected:    "Test",
		},
		{
			name:        "Empty string",
			input:       "",
			maxWidth:    5,
			useReserved: false,
			expected:    "     ",
		},
		{
			name:        "With uint8 maxWidth",
			input:       "Hello",
			maxWidth:    uint8(8),
			useReserved: false,
			expected:    "Hello   ",
		},
		{
			name:          "With uint16 maxWidth and uint8 reservedChars",
			input:         "Hello, World!",
			maxWidth:      uint16(10),
			reservedChars: uint8(3),
			useReserved:   true,
			expected:      "Hell...",
		},
		{
			name:          "With float32 maxWidth and int64 reservedChars",
			input:         "Hello, World!",
			maxWidth:      float32(10.5), // Should truncate to 10
			reservedChars: int64(2),
			useReserved:   true,
			expected:      "Hello...",
		},
		{
			name:        "With float64 maxWidth",
			input:       "Testing",
			maxWidth:    float64(9.9), // Should truncate to 9
			useReserved: false,
			expected:    "Testing  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.useReserved {
				result = Convert(tt.input).Truncate(tt.maxWidth, tt.reservedChars).String()
			} else {
				result = Convert(tt.input).Truncate(tt.maxWidth).String()
			}

			if result != tt.expected {
				t.Errorf("Convert(%q).Truncate() = %q, want %q",
					tt.input, result, tt.expected)
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
				return t.ToUpper().Truncate(10)
			},
		},
		{
			name:  "Lowercase and truncate",
			input: "HELLO WORLD",
			want:  "hello...",
			function: func(t *Text) *Text {
				return t.ToLower().Truncate(8)
			},
		},
		{
			name:  "Remove tilde and truncate",
			input: "Ñandú está corriendo",
			want:  "Nandu esta ...",
			function: func(t *Text) *Text {
				return t.RemoveTilde().Truncate(14)
			},
		},
		{
			name:  "CamelCase and truncate",
			input: "hello world example",
			want:  "helloWorld...",
			function: func(t *Text) *Text {
				return t.CamelCaseLower().Truncate(13)
			},
		},
		{
			name:  "Truncate and repeat",
			input: "hello",
			want:  "hello hello ",
			function: func(t *Text) *Text {
				return t.Truncate(6).Repeat(2)
			},
		},
		{
			name:  "SnakeCase and truncate",
			input: "Hello World Example",
			want:  "hello_world_...",
			function: func(t *Text) *Text {
				return t.ToSnakeCaseLower().Truncate(15)
			},
		},
		{
			name:  "Truncate with custom separator",
			input: "Hello World Example",
			want:  "hello-world-ex...",
			function: func(t *Text) *Text {
				return t.ToSnakeCaseLower("-").Truncate(17)
			},
		},
		{
			name:  "Complex chaining with truncate",
			input: "Él Múrcielago Rápido",
			want:  "ELMURC...",
			function: func(t *Text) *Text {
				return t.RemoveTilde().CamelCaseLower().ToUpper().Truncate(9)
			},
		},
		{
			name:  "Using explicit reserved chars",
			input: "Hello, World!",
			want:  "Hell...",
			function: func(t *Text) *Text {
				return t.Truncate(10, 3)
			},
		},
		{
			name:  "Using different numeric types",
			input: "Testing different types",
			want:  "Testing...",
			function: func(t *Text) *Text {
				return t.Truncate(uint8(12), float64(2))
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
