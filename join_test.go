package tinystring_test

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestJoinMethod(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		sep      string
		expected string
	}{
		{
			name:     "Join with default space separator",
			input:    []string{"Hello", "World"},
			sep:      "",
			expected: "Hello World",
		},
		{
			name:     "Join with custom separator",
			input:    []string{"hello", "world", "example"},
			sep:      "-",
			expected: "hello-world-example",
		},
		{
			name:     "Join with semicolon separator",
			input:    []string{"apple", "orange", "banana"},
			sep:      ";",
			expected: "apple;orange;banana",
		},
		{
			name:     "Empty slice",
			input:    []string{},
			sep:      ",",
			expected: "",
		},
		{
			name:     "Single element",
			input:    []string{"test"},
			sep:      ":",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if tt.sep == "" {
				result = Convert(tt.input).Join().String()
			} else {
				result = Convert(tt.input).Join(tt.sep).String()
			}

			if result != tt.expected {
				t.Errorf("Join test %q: expected %q, got %q",
					tt.name, tt.expected, result)
			}
		})
	}
}

func TestJoinChainMethods(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		function func([]string) string
		expected string
	}{
		{
			name:     "Join with ToUpper",
			input:    []string{"hello", "world"},
			expected: "HELLO WORLD",
			function: func(input []string) string {
				return Convert(input).Join().ToUpper().String()
			},
		},
		{
			name:     "Join with custom separator and ToLower",
			input:    []string{"HELLO", "WORLD"},
			expected: "hello-world",
			function: func(input []string) string {
				return Convert(input).Join("-").ToLower().String()
			},
		},
		{
			name:     "Join with Capitalize",
			input:    []string{"hello", "world"},
			expected: "Hello World",
			function: func(input []string) string {
				return Convert(input).Join().Capitalize().String()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.input)
			if result != tt.expected {
				t.Errorf("Chain test %q: expected %q, got %q",
					tt.name, tt.expected, result)
			}
		})
	}
}

// Test Join function with various scenarios to improve coverage
func TestJoinEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		sep      []string
		expected string
	}{
		{"Empty slice with default separator", []string{}, []string{}, ""},
		{"Empty slice with custom separator", []string{}, []string{","}, ""},
		{"Single element with default separator", []string{"Hello"}, []string{}, "Hello"},
		{"Single element with custom separator", []string{"Hello"}, []string{","}, "Hello"},
		{"Two elements with default separator", []string{"Hello", "World"}, []string{}, "Hello World"},
		{"Multiple elements with comma", []string{"A", "B", "C"}, []string{","}, "A,B,C"},
		{"Multiple elements with pipe", []string{"X", "Y", "Z"}, []string{"|"}, "X|Y|Z"},
		{"Elements with spaces", []string{"Hello World", "Test String"}, []string{" - "}, "Hello World - Test String"},
		{"Empty string elements", []string{"", "Middle", ""}, []string{","}, ",Middle,"},
		{"All empty strings", []string{"", "", ""}, []string{","}, ",,"},
		{"Multi-character separator", []string{"One", "Two", "Three"}, []string{" -> "}, "One -> Two -> Three"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := Convert(tt.input)
			var result string
			if len(tt.sep) == 0 {
				result = conv.Join().String()
			} else {
				result = conv.Join(tt.sep[0]).String()
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test Join with string splitting behavior
func TestJoinStringSplitting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		sep      string
		expected string
	}{
		{"String with spaces", "Hello World Test", ",", "Hello,World,Test"},
		{"String with tabs", "A\tB\tC", "|", "A|B|C"},
		{"String with newlines", "Line1\nLine2\nLine3", " - ", "Line1 - Line2 - Line3"},
		{"String with mixed whitespace", "A \tB\nC\rD", "_", "A_B_C_D"},
		{"String with multiple spaces", "A  B   C", ",", "A,B,C"},
		{"String with leading/trailing spaces", " Start End ", "-", "Start-End"},
		{"String with only whitespace", "   \t\n\r   ", ",", ""},
		{"Single word", "Word", ",", "Word"},
		{"Empty string", "", ",", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).Join(tt.sep).String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test Join chaining operations
func TestJoinChaining(t *testing.T) {
	// Test converting string slice to joined string and then transforming
	slice := []string{"hello", "world", "test"}

	result1 := Convert(slice).Join("_").ToUpper().String()
	expected1 := "HELLO_WORLD_TEST"
	if result1 != expected1 {
		t.Errorf("Expected %q, got %q", expected1, result1)
	}

	result2 := Convert(slice).Join("-").Capitalize().String()
	expected2 := "Hello-world-test"
	if result2 != expected2 {
		t.Errorf("Expected %q, got %q", expected2, result2)
	}

	// Test with empty slice
	emptySlice := []string{}
	result3 := Convert(emptySlice).Join(",").ToUpper().String()
	expected3 := ""
	if result3 != expected3 {
		t.Errorf("Expected %q, got %q", expected3, result3)
	}
}
