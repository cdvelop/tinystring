package tinystring

import "testing"

func TestQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "String with spaces",
			input:    "hello world",
			expected: `"hello world"`,
		},
		{
			name:     "String with quotes",
			input:    `say "hello"`,
			expected: `"say \"hello\""`,
		},
		{
			name:     "String with backslash",
			input:    `path\to\file`,
			expected: `"path\\to\\file"`,
		},
		{
			name:     "String with newline",
			input:    "line1\nline2",
			expected: `"line1\nline2"`,
		},
		{
			name:     "String with tab",
			input:    "before\tafter",
			expected: `"before\tafter"`,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: `""`,
		},
		{
			name:     "String with carriage return",
			input:    "before\rafter",
			expected: `"before\rafter"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).Quote().String()
			if result != tt.expected {
				t.Errorf("Quote() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestQuoteWithError(t *testing.T) {
	// Test quote functionality with error handling
	result, err := Convert("test").Quote().StringError()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := `"test"`
	if result != expected {
		t.Errorf("Quote() = %q, want %q", result, expected)
	}
}

func TestQuoteChaining(t *testing.T) {
	// Test chaining quote with other operations
	result := Convert("hello").Quote().String()
	expected := `"hello"`
	if result != expected {
		t.Errorf("Quote chaining = %q, want %q", result, expected)
	}

	// Test quote after conversion
	result2 := FromInt(123).Quote().String()
	expected2 := `"123"`
	if result2 != expected2 {
		t.Errorf("Quote after conversion = %q, want %q", result2, expected2)
	}
}
