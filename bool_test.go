package tinystring

import "testing"

func TestToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
		hasError bool
	}{
		{
			name:     "String true",
			input:    "true",
			expected: true,
			hasError: false,
		},
		{
			name:     "String false",
			input:    "false",
			expected: false,
			hasError: false,
		},
		{
			name:     "String 1",
			input:    "1",
			expected: true,
			hasError: false,
		},
		{
			name:     "String 0",
			input:    "0",
			expected: false,
			hasError: false,
		},
		{
			name:     "Boolean true",
			input:    true,
			expected: true,
			hasError: false,
		},
		{
			name:     "Boolean false",
			input:    false,
			expected: false,
			hasError: false,
		},
		{
			name:     "Integer 1",
			input:    1,
			expected: true,
			hasError: false,
		},
		{
			name:     "Integer 0",
			input:    0,
			expected: false,
			hasError: false,
		},
		{
			name:     "Integer non-zero",
			input:    42,
			expected: true,
			hasError: false,
		},
		{
			name:     "Invalid string",
			input:    "invalid",
			expected: false,
			hasError: true,
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input).ToBool()

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestFromBool(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{
			name:     "True to string",
			input:    true,
			expected: "true",
		},
		{
			name:     "False to string",
			input:    false,
			expected: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).String()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBoolChaining(t *testing.T) {
	// Test chaining with boolean operations
	result := Convert(true).ToUpper().String()
	expected := "TRUE"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test converting back
	boolVal, err := Convert("TRUE").ToLower().ToBool()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !boolVal {
		t.Errorf("Expected true, got false")
	}
}
