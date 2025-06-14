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
			result := tt.chain(Convert(tt.input)).String()
			if result != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}

// Test for ToSnakeCaseUpper function
func TestToSnakeCaseUpper(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		sep      []string
	}{
		{
			name:     "Simple words",
			input:    "hello world",
			expected: "Hello_World",
		},
		{
			name:     "CamelCase input",
			input:    "helloWorld",
			expected: "Hello_World",
		},
		{
			name:     "PascalCase input",
			input:    "HelloWorld",
			expected: "Hello_World",
		},
		{
			name:     "Multiple words",
			input:    "hello world example",
			expected: "Hello_World_Example",
		},
		{
			name:     "With numbers",
			input:    "hello123World",
			expected: "Hello123_World",
		},
		{
			name:     "Custom separator",
			input:    "hello world",
			expected: "Hello-World",
			sep:      []string{"-"},
		},
		{
			name:     "Custom separator with dot",
			input:    "helloWorld",
			expected: "Hello.World",
			sep:      []string{"."},
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
			result := Convert(tt.input).ToSnakeCaseUpper(tt.sep...).String()
			if result != tt.expected {
				t.Errorf("ToSnakeCaseUpper() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test for case conversion edge cases
func TestCaseConversionEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		function func(*conv) *conv
		expected string
	}{
		{
			name:     "ToSnakeCaseUpper with already snake_case",
			input:    "hello_world_example",
			function: func(c *conv) *conv { return c.ToSnakeCaseUpper() },
			expected: "Hello_world_example", // underscores are not treated as word boundaries
		},
		{
			name:     "ToSnakeCaseUpper with mixed separators",
			input:    "hello-world_example",
			function: func(c *conv) *conv { return c.ToSnakeCaseUpper() },
			expected: "Hello-world_example", // only spaces are word boundaries by default
		},
		{
			name:     "ToSnakeCaseUpper with spaces and underscores",
			input:    "hello world_example",
			function: func(c *conv) *conv { return c.ToSnakeCaseUpper() },
			expected: "Hello_World_example", // space creates word boundary, underscore doesn't
		},
		{
			name:     "ToSnakeCaseUpper with special characters",
			input:    "hello@world#example",
			function: func(c *conv) *conv { return c.ToSnakeCaseUpper() },
			expected: "Hello@world#example", // special chars are preserved, not treated as boundaries
		},
		{
			name:     "ToSnakeCaseUpper with camelCase boundaries",
			input:    "helloWorldExample",
			function: func(c *conv) *conv { return c.ToSnakeCaseUpper() },
			expected: "Hello_World_Example", // camelCase boundaries work properly
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(Convert(tt.input)).String()
			if result != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, result, tt.expected)
			}
		})
	}
}
