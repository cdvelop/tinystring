package tinystring

import "testing"

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple conv",
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
			out := Convert(tt.input).Capitalize().String()
			if out != tt.expected {
				t.Errorf("Capitalize() = %q, want %q", out, tt.expected)
			}
		})
	}
}

func TestCapitalizeChaining(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		chain    func(*conv) *conv
	}{
		{
			name:     "With RemoveTilde",
			input:    "hólá múndo",
			expected: "Hola Mundo",
			chain: func(conv *conv) *conv {
				return conv.RemoveTilde().Capitalize()
			},
		},
		{
			name:     "After ToLower",
			input:    "HELLO WORLD",
			expected: "Hello World",
			chain: func(conv *conv) *conv {
				return conv.ToLower().Capitalize()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.chain(Convert(tt.input)).String()
			if out != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, out, tt.expected)
			}
		})
	}
}
