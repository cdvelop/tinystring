package tinystring

import "testing"

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "hello world",
			expected: "Hello World",
		},
		{
			name:     "Already capitalized",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Mixed case",
			input:    "hELLo wORLd",
			expected: "Hello World",
		},
		{
			name:     "Extra spaces",
			input:    "  hello   world  ",
			expected: "Hello World",
		},
		{
			name:     "With numbers",
			input:    "hello 123 world",
			expected: "Hello 123 World",
		},
		{
			name:     "With special characters",
			input:    "héllö wörld",
			expected: "Héllö Wörld",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single word",
			input:    "hello",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).Capitalize().String()
			if result != tt.expected {
				t.Errorf("Capitalize() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCapitalizeChaining(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		chain    func(*Text) *Text
	}{
		{
			name:     "With RemoveTilde",
			input:    "hólá múndo",
			expected: "Hola Mundo",
			chain: func(text *Text) *Text {
				return text.RemoveTilde().Capitalize()
			},
		},
		{
			name:     "After ToLower",
			input:    "HELLO WORLD",
			expected: "Hello World",
			chain: func(text *Text) *Text {
				return text.ToLower().Capitalize()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.chain(Convert(tt.input)).String()
			if result != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}
