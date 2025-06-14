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

func TestBooleanToIntegerConversion(t *testing.T) {
	// Test that boolean values cannot be converted to integers (should return error)
	t.Run("Boolean true to int should fail", func(t *testing.T) {
		result, err := Convert(true).ToInt()
		if err == nil {
			t.Error("ToInt() should fail for boolean input")
		}
		if result != 0 {
			t.Errorf("ToInt() should return 0 for failed conversion, got %v", result)
		}
	})

	t.Run("Boolean false to int should fail", func(t *testing.T) {
		result, err := Convert(false).ToInt()
		if err == nil {
			t.Error("ToInt() should fail for boolean input")
		}
		if result != 0 {
			t.Errorf("ToInt() should return 0 for failed conversion, got %v", result)
		}
	})

	t.Run("Boolean to uint should fail", func(t *testing.T) {
		result, err := Convert(true).ToUint()
		if err == nil {
			t.Error("ToUint() should fail for boolean input")
		}
		if result != 0 {
			t.Errorf("ToUint() should return 0 for failed conversion, got %v", result)
		}
	})

	// Test that boolean to string conversion works correctly
	t.Run("Boolean to string conversion", func(t *testing.T) {
		trueResult := Convert(true).String()
		if trueResult != "true" {
			t.Errorf("Convert(true).String() = %q, want \"true\"", trueResult)
		}

		falseResult := Convert(false).String()
		if falseResult != "false" {
			t.Errorf("Convert(false).String() = %q, want \"false\"", falseResult)
		}
	})
}

// Test error conditions for ToBool
func TestToBoolErrorHandling(t *testing.T) {
	tests := []struct {
		name         string
		input        any
		shouldError  bool
		expectedBool bool
	}{
		{"Float string with error", "invalid_float", true, false},
		{"Complex string", "complex_value", true, false},
		{"String with numeric content - positive", "123", false, true},
		{"String with numeric content - zero", "0", false, false},
		{"String with float content - positive", "3.14", false, true},
		{"String with float content - zero", "0.0", false, false},
		{"String with float content - negative", "-1.5", false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Convert(test.input).ToBool()

			if test.shouldError {
				if err == nil {
					t.Errorf("Expected error for input %v, but got none", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %v: %v", test.input, err)
				}
				if result != test.expectedBool {
					t.Errorf("Expected %v, got %v for input %v", test.expectedBool, result, test.input)
				}
			}
		})
	}
}

// Test ToBool with pre-existing error in conv
func TestToBoolWithPreExistingError(t *testing.T) {
	// Create a conv with an error
	c := &conv{err: "pre-existing error"}

	result, err := c.ToBool()

	if err == nil {
		t.Error("Expected error to be returned when conv already has error")
	}

	if result != false {
		t.Errorf("Expected false when error exists, got %v", result)
	}
}

// Test ToBool with numeric types and zero values
func TestToBoolNumericTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"Int zero", int(0), false},
		{"Int positive", int(42), true},
		{"Int negative", int(-1), true},
		{"Int8 zero", int8(0), false},
		{"Int8 positive", int8(1), true},
		{"Int16 zero", int16(0), false},
		{"Int16 positive", int16(100), true},
		{"Int32 zero", int32(0), false},
		{"Int32 positive", int32(1000), true},
		{"Int64 zero", int64(0), false},
		{"Int64 positive", int64(999999), true},
		{"Uint zero", uint(0), false},
		{"Uint positive", uint(42), true},
		{"Uint8 zero", uint8(0), false},
		{"Uint8 positive", uint8(255), true},
		{"Uint16 zero", uint16(0), false},
		{"Uint16 positive", uint16(65535), true},
		{"Uint32 zero", uint32(0), false},
		{"Uint32 positive", uint32(4294967295), true},
		{"Uint64 zero", uint64(0), false},
		{"Uint64 positive", uint64(18446744073709551615), true},
		{"Uintptr zero", uintptr(0), false},
		{"Uintptr positive", uintptr(0x1000), true},
		{"Float32 zero", float32(0.0), false},
		{"Float32 positive", float32(3.14), true},
		{"Float32 negative", float32(-2.71), true},
		{"Float64 zero", float64(0.0), false},
		{"Float64 positive", float64(3.14159), true},
		{"Float64 negative", float64(-1.414), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Convert(test.input).ToBool()

			if err != nil {
				t.Errorf("Unexpected error for %v: %v", test.input, err)
			}

			if result != test.expected {
				t.Errorf("Expected %v for input %v, got %v", test.expected, test.input, result)
			}
		})
	}
}

// Test ToBool with string variations
func TestToBoolStringVariations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Lowercase t", "t", true},
		{"Uppercase T", "T", true},
		{"Lowercase f", "f", false},
		{"Uppercase F", "F", false},
		{"Numeric string positive", "456", true},
		{"Numeric string negative", "-789", true},
		{"Float string positive", "2.718", true},
		{"Float string negative", "-3.14", true},
		{"Zero string", "0", false},
		{"Zero float string", "0.0", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Convert(test.input).ToBool()

			if err != nil {
				t.Errorf("Unexpected error for %q: %v", test.input, err)
			}

			if result != test.expected {
				t.Errorf("Expected %v for input %q, got %v", test.expected, test.input, result)
			}
		})
	}
}
