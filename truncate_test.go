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
			name:        "conv shorter than max width",
			input:       "Hello",
			maxWidth:    10,
			useReserved: false,
			expected:    "Hello", // No padding expected
		},
		{
			name:        "conv longer than max width with ellipsis",
			input:       "Hello, World!",
			maxWidth:    10,
			useReserved: false,
			expected:    "Hello, ...",
		},
		{
			name:          "conv longer with reserved chars",
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
			expected:    "", // No padding expected
		},
		{
			name:        "With uint8 maxWidth",
			input:       "Hello",
			maxWidth:    uint8(8),
			useReserved: false,
			expected:    "Hello", // No padding expected
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
			expected:    "Testing", // No padding expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out string
			if tt.useReserved {
				out = Convert(tt.input).Truncate(tt.maxWidth, tt.reservedChars).String()
			} else {
				out = Convert(tt.input).Truncate(tt.maxWidth).String()
			}

			if out != tt.expected {
				t.Errorf("Convert(%q).Truncate() = %q, want %q",
					tt.input, out, tt.expected)
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
		function func(*conv) *conv
	}{
		{
			name:  "Uppercase and truncate",
			input: "hello world",
			want:  "HELLO W...",
			function: func(t *conv) *conv {
				return t.ToUpper().Truncate(10)
			},
		},
		{
			name:  "Lowercase and truncate",
			input: "HELLO WORLD",
			want:  "hello...",
			function: func(t *conv) *conv {
				return t.ToLower().Truncate(8)
			},
		},
		{
			name:  "Remove tilde and truncate",
			input: "Ñandú está corriendo",
			want:  "Nandu esta ...",
			function: func(t *conv) *conv {
				return t.RemoveTilde().Truncate(14)
			},
		},
		{
			name:  "CamelCase and truncate",
			input: "hello world example",
			want:  "helloWorld...",
			function: func(t *conv) *conv {
				return t.CamelCaseLower().Truncate(13)
			},
		},
		{
			name:  "Truncate and repeat",
			input: "hello",
			want:  "hellohello", // No padding expected before repeat
			function: func(t *conv) *conv {
				return t.Truncate(6).Repeat(2)
			},
		},
		{
			name:  "SnakeCase and truncate",
			input: "Hello World Example",
			want:  "hello_world_...",
			function: func(t *conv) *conv {
				return t.ToSnakeCaseLower().Truncate(15)
			},
		},
		{
			name:  "Truncate with custom separator",
			input: "Hello World Example",
			want:  "hello-world-ex...",
			function: func(t *conv) *conv {
				return t.ToSnakeCaseLower("-").Truncate(17)
			},
		},
		{
			name:  "Complex chaining with truncate",
			input: "Él Múrcielago Rápido",
			want:  "ELMURC...",
			function: func(t *conv) *conv {
				return t.RemoveTilde().CamelCaseLower().ToUpper().Truncate(9)
			},
		},
		{
			name:  "Using explicit reserved chars",
			input: "Hello, World!",
			want:  "Hell...",
			function: func(t *conv) *conv {
				return t.Truncate(10, 3)
			},
		},
		{
			name:  "Using different numeric types",
			input: "Testing different types",
			want:  "Testing...",
			function: func(t *conv) *conv {
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

func TestTruncateName(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		maxCharsPerWord any
		maxWidth        any
		expected        string
	}{
		{
			name:            "Basic example with one name one surname",
			input:           "Jeronimo Dominguez",
			maxCharsPerWord: 3,
			maxWidth:        15,
			expected:        "Jer. Dominguez",
		},
		{
			name:            "Basic example with name exceeding maxWidth",
			input:           "Jeronimo Dominguez",
			maxCharsPerWord: 3,
			maxWidth:        10,
			expected:        "Jer. Do...",
		},
		{
			name:            "Multiple names and surnames",
			input:           "Ana Maria Rodriguez Sanchez",
			maxCharsPerWord: 2,
			maxWidth:        20,
			expected:        "An. Ma. Ro. Sanchez",
		},
		{
			name:            "With total length constraint",
			input:           "Ana Maria Rodriguez Sanchez",
			maxCharsPerWord: 2,
			maxWidth:        12,
			expected:        "An. Ma. R...",
		},
		{
			name:            "Short names no truncation needed",
			input:           "Ana Gil",
			maxCharsPerWord: 3,
			maxWidth:        10,
			expected:        "Ana Gil",
		},
		{
			name:            "Mixed length names",
			input:           "Bob Alexandrovich",
			maxCharsPerWord: 4,
			maxWidth:        15,
			expected:        "Bob Alexandr...",
		},
		{
			name:            "Very short max chars",
			input:           "John Smith",
			maxCharsPerWord: 1,
			maxWidth:        10,
			expected:        "J. Smith",
		},
		{
			name:            "Empty string",
			input:           "",
			maxCharsPerWord: 3,
			maxWidth:        10,
			expected:        "",
		},
		{
			name:            "Single name only",
			input:           "Alexander",
			maxCharsPerWord: 4,
			maxWidth:        9,
			expected:        "Alexander",
		},
		{
			name:            "Very restrictive total length",
			input:           "Alexander Graham Bell",
			maxCharsPerWord: 3,
			maxWidth:        7,
			expected:        "Alex...",
		},
		{
			name:            "With uint8 max chars",
			input:           "Manuel Rodriguez",
			maxCharsPerWord: uint8(3),
			maxWidth:        15,
			expected:        "Man. Rodriguez",
		},
		{
			name:            "With float64 max length",
			input:           "Pedro Gonzalez",
			maxCharsPerWord: 3,
			maxWidth:        float64(10.5), // Should truncate to 10
			expected:        "Ped. Go...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := Convert(tt.input).TruncateName(tt.maxCharsPerWord, tt.maxWidth).String()

			if out != tt.expected {
				t.Errorf("Convert(%q).TruncateName(%v, %v) = %q, want %q",
					tt.input, tt.maxCharsPerWord, tt.maxWidth, out, tt.expected)
			}
		})
	}
}

// TestTruncateNameChain tests the chaining capabilities of the TruncateName method with other methods
func TestTruncateNameChain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		function func(*conv) *conv
	}{{
		name:  "Uppercase and truncate name",
		input: "carlos mendez",
		want:  "CAR. MENDEZ", // only truncation first word
		function: func(t *conv) *conv {
			return t.ToUpper().TruncateName(3, 15)
		},
	}, {
		name:  "Remove tilde and truncate name",
		input: "José Martínez",
		want:  "Jose Martinez", // No truncation (4) needed within maxWidth 15
		function: func(t *conv) *conv {
			return t.RemoveTilde().TruncateName(4, 15)
		},
	}, {
		name:  "Complex chaining",
		input: "MARÍA del carmen GARCÍA",
		want:  "mar. del car. garcía", // truncation per word needed within maxWidth 25
		function: func(t *conv) *conv {
			return t.ToLower().TruncateName(3, 25)
		},
	}, {
		name:  "With total length limit",
		input: "Roberto Carlos Fernandez",
		want:  "Rob. Car...", // Truncation needed at maxWidth 11
		function: func(t *conv) *conv {
			return t.TruncateName(3, 11)
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
