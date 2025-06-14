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
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input).ToUint()

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
				// Use tolerance for floating-point comparison
				tolerance := 1e-5 // Increased tolerance for floating-point precision issues
				if result < tt.expected-tolerance || result > tt.expected+tolerance {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestToIntConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int64
		hasError bool
	}{
		{name: "int value", input: 42, expected: 42, hasError: false},
		{name: "int8 value", input: int8(8), expected: 8, hasError: false},
		{name: "int16 value", input: int16(16), expected: 16, hasError: false},
		{name: "int32 value", input: int32(32), expected: 32, hasError: false},
		{name: "int64 value", input: int64(64), expected: 64, hasError: false},
		{name: "uint value", input: uint(42), expected: 42, hasError: false},
		{name: "uint8 value", input: uint8(8), expected: 8, hasError: false},
		{name: "uint16 value", input: uint16(16), expected: 16, hasError: false},
		{name: "uint32 value", input: uint32(32), expected: 32, hasError: false},
		{name: "uint64 value", input: uint64(64), expected: 64, hasError: false},
		{name: "float32 value (truncation)", input: float32(3.14), expected: 3, hasError: false},
		{name: "float64 value (truncation)", input: float64(6.28), expected: 6, hasError: false},
		{name: "string numeric value", input: "12345", expected: 12345, hasError: false},
		{name: "string negative numeric value", input: "-50", expected: -50, hasError: false},
		{name: "string float numeric value (truncation)", input: "123.789", expected: 123, hasError: false},
		{name: "string value (invalid)", input: "not a number", expected: 0, hasError: true},
		{name: "boolean value (invalid)", input: true, expected: 0, hasError: true},
		{name: "nil value (invalid)", input: nil, expected: 0, hasError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input).ToInt()

			if tt.hasError {
				if err == nil {
					t.Errorf("Convert(%v).ToInt() expected error, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Convert(%v).ToInt() unexpected error: %v", tt.input, err)
				}
				if int64(result) != tt.expected { // Cast result to int64 for comparison
					t.Errorf("Convert(%v).ToInt() = %v, want %v",
						tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestFromNumeric(t *testing.T) {
	t.Run("Convert from int", func(t *testing.T) {
		result := Convert(-123).String()
		expected := "-123"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Convert from uint", func(t *testing.T) {
		result := Convert(uint(456)).String()
		expected := "456"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Convert from float", func(t *testing.T) {
		result := Convert(123.5).String() // Use a value exactly representable in binary
		expected := "123.5"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}

func TestNumericChaining(t *testing.T) {
	// Test converting number to string and back
	original := 12345
	converted, err := Convert(original).ToInt()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if int64(converted) != int64(original) {
		t.Errorf("Expected %d, got %d", original, converted)
	}

	// Test with formatting
	result := Convert(123.456).RoundDecimals(2).String()
	expected := "123.46"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test with formatting numbers
	result = Convert(1234567).FormatNumber().String()
	expected = "1.234.567"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFixedNegativeNumbers(t *testing.T) {
	// Test negative numbers in s2Int
	result, err := Convert("-123").ToInt()
	if err != nil {
		t.Errorf("ToInt(-123) failed: %v", err)
	}
	if result != -123 {
		t.Errorf("ToInt(-123) = %d, want -123", result)
	}

	// Test negative numbers in s2Int64
	result64, err := Convert("-9223372036854775807").ToInt64()
	if err != nil {
		t.Errorf("ToInt64(-9223372036854775807) failed: %v", err)
	}
	if result64 != -9223372036854775807 {
		t.Errorf("ToInt64(-9223372036854775807) = %d, want -9223372036854775807", result64)
	}

	// Test negative numbers should fail for non-decimal bases
	_, err = Convert("-123").ToInt(16)
	if err == nil {
		t.Error("ToInt(-123, base 16) should have failed but didn't")
	}
}

func TestFloatStringConversion(t *testing.T) {
	// Test string float conversion - validates the internal s2Float conversion path
	t.Run("String to float conversion", func(t *testing.T) {
		result, err := Convert("123.789").ToFloat()
		if err != nil {
			t.Errorf("ToFloat() failed: %v", err)
		}
		expected := 123.789
		tolerance := 1e-6
		if result < expected-tolerance || result > expected+tolerance {
			t.Errorf("ToFloat() = %v, want %v", result, expected)
		}
	})

	// Test string to int conversion with float input (truncation)
	t.Run("String float to int truncation", func(t *testing.T) {
		result, err := Convert("123.789").ToInt()
		if err != nil {
			t.Errorf("ToInt() failed: %v", err)
		}
		expected := 123
		if result != expected {
			t.Errorf("ToInt() = %v, want %v", result, expected)
		}
	})

	// Test edge case: empty string
	t.Run("Empty string conversion", func(t *testing.T) {
		_, err := Convert("").ToFloat()
		if err == nil {
			t.Error("ToFloat() should fail for empty string")
		}
	})

	// Test edge case: invalid float string
	t.Run("Invalid float string", func(t *testing.T) {
		_, err := Convert("not-a-number").ToFloat()
		if err == nil {
			t.Error("ToFloat() should fail for invalid string")
		}
	})
}

func TestRoundDecimalsNonNumericInput(t *testing.T) {
	// Test RoundDecimals behavior with non-numeric input
	testCases := []struct {
		name     string
		input    any
		decimals int
		expected string
	}{
		{
			name:     "Non-numeric string should become 0.00",
			input:    "hello",
			decimals: 2,
			expected: "0.00",
		},
		{
			name:     "Invalid string with 3 decimals",
			input:    "not-a-number",
			decimals: 3,
			expected: "0.000",
		},
		{
			name:     "Empty string with 1 decimal",
			input:    "",
			decimals: 1,
			expected: "0.0",
		},
		{
			name:     "Special characters with 0 decimals",
			input:    "!@#$%",
			decimals: 0,
			expected: "0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Convert(tc.input).RoundDecimals(tc.decimals).String()
			if result != tc.expected {
				t.Errorf("Convert(%v).RoundDecimals(%d).String() = %q, want %q",
					tc.input, tc.decimals, result, tc.expected)
			}
		})
	}
}

func TestRoundDecimalsEdgeCases(t *testing.T) {
	// Test additional edge cases for RoundDecimals
	t.Run("Mixed valid and invalid conversions", func(t *testing.T) {
		// Valid conversion first
		result1 := Convert("123.456").RoundDecimals(2).String()
		expected1 := "123.46"
		if result1 != expected1 {
			t.Errorf("Valid conversion failed: got %q, want %q", result1, expected1)
		}

		// Invalid conversion after valid one
		result2 := Convert("invalid").RoundDecimals(2).String()
		expected2 := "0.00"
		if result2 != expected2 {
			t.Errorf("Invalid conversion failed: got %q, want %q", result2, expected2)
		}
	})

	t.Run("Zero decimals with non-numeric", func(t *testing.T) {
		result := Convert("text").RoundDecimals(0).String()
		expected := "0"
		if result != expected {
			t.Errorf("Zero decimals failed: got %q, want %q", result, expected)
		}
	})
}

func TestToIntInvalidBase(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		base        int
		expectError bool
	}{
		// These tests use strings that can't fall back to float parsing
		// and must go through base validation
		{"base 1 invalid with ABC", "ABC", 1, true},
		{"base 0 invalid with ABC", "ABC", 0, true},
		{"base -1 invalid with ABC", "ABC", -1, true},
		{"base 37 invalid with ABC", "ABC", 37, true},
		{"base 100 invalid with ABC", "ABC", 100, true},
		{"base 2 valid with 101", "101", 2, false},
		{"base 16 valid with ABC", "ABC", 16, false},
		{"base 36 valid with ZZ", "ZZ", 36, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Convert(tt.input).ToInt(tt.base)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for base %d with input %q, but got none", tt.base, tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for base %d with input %q: %v", tt.base, tt.input, err)
				}
			}
		})
	}
}

func TestToUintInvalidBase(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		base        int
		expectError bool
	}{
		// These tests use strings that can't fall back to float parsing
		// and must go through base validation
		{"base 1 invalid with ABC", "ABC", 1, true},
		{"base 0 invalid with ABC", "ABC", 0, true},
		{"base -1 invalid with ABC", "ABC", -1, true},
		{"base 37 invalid with ABC", "ABC", 37, true},
		{"base 100 invalid with ABC", "ABC", 100, true},
		{"base 2 valid with 101", "101", 2, false},
		{"base 16 valid with ABC", "ABC", 16, false},
		{"base 36 valid with ZZ", "ZZ", 36, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Convert(tt.input).ToUint(tt.base)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for base %d with input %q, but got none", tt.base, tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for base %d with input %q: %v", tt.base, tt.input, err)
				}
			}
		})
	}
}
