package tinystring

import "testing"

func TestToBool(t *testing.T) {
	tests := []struct {
		name       string
		input      any
		expected   bool
		hasContent bool
	}{
		{
			name:       "String true",
			input:      "true",
			expected:   true,
			hasContent: false,
		},
		{
			name:       "String false",
			input:      "false",
			expected:   false,
			hasContent: false,
		},
		{
			name:       "String 1",
			input:      "1",
			expected:   true,
			hasContent: false,
		},
		{
			name:       "String 0",
			input:      "0",
			expected:   false,
			hasContent: false,
		},
		{
			name:       "Boolean true",
			input:      true,
			expected:   true,
			hasContent: false,
		},
		{
			name:       "Boolean false",
			input:      false,
			expected:   false,
			hasContent: false,
		},
		{
			name:       "Integer 1",
			input:      1,
			expected:   true,
			hasContent: false,
		},
		{
			name:       "Integer 0",
			input:      0,
			expected:   false,
			hasContent: false,
		},
		{
			name:       "Integer non-zero",
			input:      42,
			expected:   true,
			hasContent: false,
		},
		{
			name:       "Invalid string",
			input:      "invalid",
			expected:   false,
			hasContent: true,
		},
		{
			name:       "Nil input",
			input:      nil,
			expected:   false,
			hasContent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Convert(tt.input).Bool()

			if tt.hasContent {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if out != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, out)
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
			out := Convert(tt.input).String()
			if out != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, out)
			}
		})
	}
}

func TestBoolChaining(t *testing.T) {
	// Test chaining with boolean operations
	out := Convert(true).Up().String()
	expected := "TRUE"
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Test converting back
	boolVal, err := Convert("TRUE").Low().Bool()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !boolVal {
		t.Errorf("Expected true, got false")
	}
}
