package tinystring

import "testing"

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
			expected:        "An. Ma. Rodriguez...",
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
			expected:        "Bob Alex.",
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
			result := Convert(tt.input).TruncateName(tt.maxCharsPerWord, tt.maxWidth).String()

			if result != tt.expected {
				t.Errorf("Convert(%q).TruncateName(%v, %v) = %q, want %q",
					tt.input, tt.maxCharsPerWord, tt.maxWidth, result, tt.expected)
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
		function func(*Text) *Text
	}{{
		name:  "Uppercase and truncate name",
		input: "carlos mendez",
		want:  "CARLOS MENDEZ", // No truncation needed within maxWidth 15
		function: func(t *Text) *Text {
			return t.ToUpper().TruncateName(3, 15)
		},
	}, {
		name:  "Remove tilde and truncate name",
		input: "José Martínez",
		want:  "Jose Martinez", // No truncation needed within maxWidth 15
		function: func(t *Text) *Text {
			return t.RemoveTilde().TruncateName(3, 15)
		},
	}, {
		name:  "Complex chaining",
		input: "MARÍA del carmen GARCÍA",
		want:  "maria del carmen garcia", // No truncation needed within maxWidth 25
		function: func(t *Text) *Text {
			return t.ToLower().TruncateName(3, 25)
		},
	}, {
		name:  "With total length limit",
		input: "Roberto Carlos Fernandez",
		want:  "Rob. Car...", // Truncation needed at maxWidth 11
		function: func(t *Text) *Text {
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
