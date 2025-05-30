package tinystring

import "testing"

func TestToUint(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint64
		hasError bool
	}{
		{
			name:     "String positive number",
			input:    "123",
			expected: 123,
			hasError: false,
		},
		{
			name:     "Integer positive",
			input:    456,
			expected: 456,
			hasError: false,
		},
		{
			name:     "Uint value",
			input:    uint(789),
			expected: 789,
			hasError: false,
		},
		{
			name:     "Float positive",
			input:    123.45,
			expected: 123,
			hasError: false,
		},
		{
			name:     "String negative number",
			input:    "-123",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Integer negative",
			input:    -456,
			expected: 0,
			hasError: true,
		},
		{
			name:     "Invalid string",
			input:    "invalid",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {			result, err := Convert(tt.input).ToUint()
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if uint64(result) != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
		hasError bool
	}{
		{
			name:     "String float",
			input:    "123.456",
			expected: 123.456,
			hasError: false,
		},
		{
			name:     "Integer",
			input:    123,
			expected: 123.0,
			hasError: false,
		},
		{
			name:     "Float value",
			input:    456.789,
			expected: 456.789,
			hasError: false,
		},
		{
			name:     "String negative",
			input:    "-123.456",
			expected: -123.456,
			hasError: false,
		},
		{
			name:     "Invalid string",
			input:    "invalid",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input).ToFloat()
			
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

func TestFromNumeric(t *testing.T) {
	t.Run("FromInt", func(t *testing.T) {
		result := FromInt(-123).String()
		expected := "-123"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("FromUint", func(t *testing.T) {
		result := FromUint(456).String()
		expected := "456"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("FromFloat", func(t *testing.T) {
		result := FromFloat(123.456).String()
		expected := "123.456"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}

func TestNumericChaining(t *testing.T) {
	// Test converting number to string and back
	original := 12345
	converted, err := FromInt(original).ToInt()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if int64(converted) != int64(original) {
		t.Errorf("Expected %d, got %d", original, converted)
	}

	// Test with formatting
	result := FromFloat(123.456).RoundDecimals(2).String()
	expected := "123.46"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test with formatting numbers
	result = FromInt(1234567).FormatNumber().String()
	expected = "1.234.567"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
