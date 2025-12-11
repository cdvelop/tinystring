package fmt

import "testing"

func TestToUint(t *testing.T) {
	tests := []struct {
		name       string
		input      any
		expected   uint64
		hasContent bool
	}{
		{
			name:       "String positive number",
			input:      "123",
			expected:   123,
			hasContent: false,
		},
		// Añadir un caso de depuración simple
		{
			name:       "Debug_Simple_String",
			input:      "42",
			expected:   42,
			hasContent: false,
		},
		{
			name:       "Integer positive",
			input:      456,
			expected:   456,
			hasContent: false,
		},
		{
			name:       "Uint value",
			input:      uint(789),
			expected:   789,
			hasContent: false,
		},
		{
			name:       "Float positive",
			input:      123.45,
			expected:   123,
			hasContent: false,
		},
		{
			name:       "String negative number",
			input:      "-123",
			expected:   0,
			hasContent: true,
		},
		{
			name:       "Integer negative",
			input:      -456,
			expected:   0,
			hasContent: true,
		},
		{
			name:       "Invalid string",
			input:      "invalid",
			expected:   0,
			hasContent: true,
		},
		{
			name:       "Nil input",
			input:      nil,
			expected:   0,
			hasContent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Convert(tt.input).Uint()

			if tt.hasContent {
				if err == nil {
					t.Errorf("Expected %v, got %v", tt.expected, out)
				}
			}
		})
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		name       string
		input      any
		expected   float64
		hasContent bool
	}{
		{
			name:       "String float",
			input:      "123.456",
			expected:   123.456,
			hasContent: false,
		},
		{
			name:       "Integer",
			input:      123,
			expected:   123.0,
			hasContent: false,
		},
		{
			name:       "Float value",
			input:      456.789,
			expected:   456.789,
			hasContent: false,
		},
		{
			name:       "String negative",
			input:      "-123.456",
			expected:   -123.456,
			hasContent: false,
		},
		{
			name:       "Invalid string",
			input:      "invalid",
			expected:   0,
			hasContent: true,
		},
		{
			name:       "Nil input",
			input:      nil,
			expected:   0,
			hasContent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Convert(tt.input).Float64()

			if tt.hasContent {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Use tolerance for floating-point comparison
				tolerance := 1e-5 // Increased tolerance for floating-point precision issues
				if out < tt.expected-tolerance || out > tt.expected+tolerance {
					t.Errorf("Expected %v, got %v", tt.expected, out)
				}
			}
		})
	}
}

func TestToIntConversion(t *testing.T) {
	tests := []struct {
		name       string
		input      any
		expected   int64
		hasContent bool
	}{
		{name: "int value", input: 42, expected: 42, hasContent: false},
		{name: "int8 value", input: int8(8), expected: 8, hasContent: false},
		{name: "int16 value", input: int16(16), expected: 16, hasContent: false},
		{name: "int32 value", input: int32(32), expected: 32, hasContent: false},
		{name: "int64 value", input: int64(64), expected: 64, hasContent: false},
		{name: "uint value", input: uint(42), expected: 42, hasContent: false},
		{name: "uint8 value", input: uint8(8), expected: 8, hasContent: false},
		{name: "uint16 value", input: uint16(16), expected: 16, hasContent: false},
		{name: "uint32 value", input: uint32(32), expected: 32, hasContent: false},
		{name: "uint64 value", input: uint64(64), expected: 64, hasContent: false},
		{name: "float32 value (truncation)", input: float32(3.14), expected: 3, hasContent: false},
		{name: "float64 value (truncation)", input: float64(6.28), expected: 6, hasContent: false},
		{name: "string numeric value", input: "12345", expected: 12345, hasContent: false},
		{name: "string negative numeric value", input: "-50", expected: -50, hasContent: false},
		{name: "string float numeric value (truncation)", input: "123.789", expected: 123, hasContent: false},
		{name: "string value (invalid)", input: "not a number", expected: 0, hasContent: true},
		{name: "boolean value (invalid)", input: true, expected: 0, hasContent: true},
		{name: "nil value (invalid)", input: nil, expected: 0, hasContent: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Convert(tt.input).Int()

			if tt.hasContent {
				if err == nil {
					t.Errorf("Convert(%v).Int() expected error, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Convert(%v).Int() unexpected error: %v", tt.input, err)
				}
				if int64(out) != tt.expected { // Cast out to int64 for comparison
					t.Errorf("Convert(%v).Int() = %v, want %v",
						tt.input, out, tt.expected)
				}
			}
		})
	}
}

func TestFromNumeric(t *testing.T) {
	t.Run("Convert from int", func(t *testing.T) {
		out := Convert(-123).String()
		expected := "-123"
		if out != expected {
			t.Errorf("Expected %q, got %q", expected, out)
		}
	})

	t.Run("Convert from uint", func(t *testing.T) {
		out := Convert(uint(456)).String()
		expected := "456"
		if out != expected {
			t.Errorf("Expected %q, got %q", expected, out)
		}
	})

	t.Run("Convert from float", func(t *testing.T) {
		out := Convert(123.5).String() // Use a value exactly representable in binary
		expected := "123.5"
		if out != expected {
			t.Errorf("Expected %q, got %q", expected, out)
		}
	})
}

func TestNumericChaining(t *testing.T) {
	original := 12345
	converted, err := Convert(original).Int()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if int64(converted) != int64(original) {
		t.Errorf("Expected %d, got %d", original, converted)
	}

	// Test with formatting
	c := Convert(123.456)
	c.Round(2)
	out := c.String()
	expected := "123.46"
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Test with formatting numbers (EU default)
	out = Convert(1234567).Thousands().String()
	expected = "1.234.567"
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}

	// Test with formatting numbers (Anglo)
	out = Convert(1234567).Thousands(true).String()
	expected = "1,234,567"
	if out != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func TestFixedNegativeNumbers(t *testing.T) {
	// Test negative numbers in s2Int
	out, err := Convert("-123").Int()
	if err != nil {
		t.Errorf("Int(-123) failed: %v", err)
	}
	if out != -123 {
		t.Errorf("Int(-123) = %d, want -123", out)
	}

	// Test negative numbers in s2Int64
	result64, err := Convert("-9223372036854775807").Int64()
	if err != nil {
		t.Errorf("Int64(-9223372036854775807) failed: %v", err)
	}
	if result64 != -9223372036854775807 {
		t.Errorf("Int64(-9223372036854775807) = %d, want -9223372036854775807", result64)
	}

	// Test negative numbers should fail for non-decimal bases
	_, err = Convert("-123").Int(16)
	if err == nil {
		t.Error("Int(-123, base 16) should have failed but didn't")
	}
}
