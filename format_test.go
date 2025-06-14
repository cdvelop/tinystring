package tinystring

import "testing"

func TestRoundDecimals(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		decimals int
		want     string
	}{{
		name:     "Round to 2 decimals",
		input:    3.12221,
		decimals: 2,
		want:     "3.13", // Corrected: now uses up rounding by default
	},
		{
			name:     "Round to 3 decimals",
			input:    3.1415926,
			decimals: 3,
			want:     "3.142",
		},
		{
			name:     "Round to 0 decimals",
			input:    3.6,
			decimals: 0,
			want:     "4",
		}, {
			name:     "Round negative to 2 decimals",
			input:    -3.12221,
			decimals: 2,
			want:     "-3.13", // Corrected: now uses up rounding by default (away from zero)
		},
		{
			name:     "Round without decimal",
			input:    3,
			decimals: 2,
			want:     "3.00",
		}, {
			name:     "Round string input",
			input:    "3.12221",
			decimals: 2,
			want:     "3.13", // Corrected: now uses up rounding by default
		},
		{
			name:     "Non-numeric input",
			input:    "hello",
			decimals: 2,
			want:     "0.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).RoundDecimals(tt.decimals).String()
			if result != tt.want {
				t.Errorf("RoundDecimals() got = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "Format integer with thousand separators",
			input: 2189009,
			want:  "2.189.009",
		},
		{
			name:  "Format decimal number with trailing zeros",
			input: 2189009.00,
			want:  "2.189.009",
		},
		{
			name:  "Format decimal number",
			input: 2189009.123,
			want:  "2.189.009.123",
		},
		{
			name:  "Format string number",
			input: "2189009.00",
			want:  "2.189.009",
		},
		{
			name:  "Format negative number",
			input: -2189009,
			want:  "-2.189.009",
		},
		{
			name:  "Format small number",
			input: 123,
			want:  "123",
		},
		{
			name:  "Format zero",
			input: 0,
			want:  "0",
		},
		{
			name:  "Non-numeric input",
			input: "hello",
			want:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).FormatNumber().String()
			if result != tt.want {
				t.Errorf("FormatNumber() got = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
		hasError bool
	}{
		{
			name:     "String formatting",
			format:   "Hello %s!",
			args:     []any{"World"},
			expected: "Hello World!",
			hasError: false,
		},
		{
			name:     "Integer formatting",
			format:   "Value: %d",
			args:     []any{42},
			expected: "Value: 42",
			hasError: false,
		},
		{
			name:     "Float formatting",
			format:   "Pi: %.2f",
			args:     []any{3.14159},
			expected: "Pi: 3.14",
			hasError: false,
		},
		{
			name:     "Multiple arguments",
			format:   "Hello %s, you have %d messages",
			args:     []any{"Alice", 5},
			expected: "Hello Alice, you have 5 messages",
			hasError: false,
		},
		{
			name:     "Binary formatting",
			format:   "Binary: %b",
			args:     []any{7},
			expected: "Binary: 111",
			hasError: false,
		},
		{
			name:     "Hexadecimal formatting",
			format:   "Hex: %x",
			args:     []any{255},
			expected: "Hex: ff",
			hasError: false,
		},
		{
			name:     "Octal formatting",
			format:   "Octal: %o",
			args:     []any{64},
			expected: "Octal: 100",
			hasError: false,
		},
		{
			name:     "Value formatting",
			format:   "Bool: %v",
			args:     []any{true},
			expected: "Bool: true",
			hasError: false,
		},
		{
			name:     "Percent sign",
			format:   "100%% complete",
			args:     []any{},
			expected: "100% complete",
			hasError: false,
		},
		{
			name:     "Missing argument",
			format:   "Value: %d",
			args:     []any{},
			expected: "",
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Format(test.format, test.args...).String()
			resultWithError, err := Format(test.format, test.args...).StringError()

			if test.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != test.expected {
					t.Errorf("Expected %q, got %q", test.expected, result)
				}
				if resultWithError != test.expected {
					t.Errorf("StringError result: Expected %q, got %q", test.expected, resultWithError)
				}
			}
		})
	}
}

func TestRoundDecimalsEnhanced(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		decimals int
		down     bool
		expected string
	}{
		{
			name:     "Round up default",
			input:    "3.154",
			decimals: 2,
			down:     false,
			expected: "3.16",
		},
		{
			name:     "Round down explicit",
			input:    "3.154",
			decimals: 2,
			down:     true,
			expected: "3.15",
		},
		{
			name:     "Round up default zero decimals",
			input:    "3.7",
			decimals: 0,
			down:     false,
			expected: "4",
		},
		{
			name:     "Round down zero decimals",
			input:    "3.7",
			decimals: 0,
			down:     true,
			expected: "3",
		},
		{
			name:     "Negative number round up",
			input:    "-3.154",
			decimals: 2,
			down:     false,
			expected: "-3.16",
		},
		{
			name:     "Negative number round down",
			input:    "-3.154",
			decimals: 2,
			down:     true,
			expected: "-3.15",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var result string
			if test.down {
				result = Convert(test.input).RoundDecimals(test.decimals).Down().String()
			} else {
				result = Convert(test.input).RoundDecimals(test.decimals).String()
			}

			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestRoundDecimalsAPI(t *testing.T) {
	// Test the corrected API as specified
	t.Run("Default up rounding", func(t *testing.T) {
		result := Convert(3.154).RoundDecimals(2).String()
		expected := "3.16"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Explicit down rounding", func(t *testing.T) {
		result := Convert(3.154).RoundDecimals(2).Down().String()
		expected := "3.15"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}

// Test for internal formatting functions
func TestFormatValueInternal(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		// Test formatAny2Int
		{
			name:     "int type",
			input:    int(42),
			expected: "42",
		},
		{
			name:     "int8 type",
			input:    int8(42),
			expected: "42",
		},
		{
			name:     "int16 type",
			input:    int16(42),
			expected: "42",
		},
		{
			name:     "int32 type",
			input:    int32(42),
			expected: "42",
		},
		{
			name:     "int64 type",
			input:    int64(42),
			expected: "42",
		},
		// Test formatAny2Uint
		{
			name:     "uint type",
			input:    uint(42),
			expected: "42",
		},
		{
			name:     "uint8 type",
			input:    uint8(42),
			expected: "42",
		},
		{
			name:     "uint16 type",
			input:    uint16(42),
			expected: "42",
		},
		{
			name:     "uint32 type",
			input:    uint32(42),
			expected: "42",
		},
		{
			name:     "uint64 type",
			input:    uint64(42),
			expected: "42",
		},
		// Test formatAny2Float
		{
			name:     "float32 type",
			input:    float32(3.14),
			expected: "3.14",
		},
		{
			name:     "float64 type",
			input:    float64(3.14159),
			expected: "3.14159",
		},
		// Test formatUnsupported
		{
			name:     "complex type",
			input:    complex(1, 2),
			expected: "<unsupported>",
		},
		{
			name:     "struct type",
			input:    struct{ Name string }{Name: "test"},
			expected: "<unsupported>",
		},
		{
			name:     "slice type",
			input:    []int{1, 2, 3},
			expected: "<unsupported>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use Format with %v to trigger formatValue internally
			result := Format("%v", tt.input).String()
			if result != tt.expected {
				t.Errorf("Format(%%v, %T(%v)) = %q, want %q", tt.input, tt.input, result, tt.expected)
			}
		})
	}
}

// Test for Errorf function
func TestErrorfFormatting(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{
			name:     "Simple error message",
			format:   "error occurred",
			args:     []any{},
			expected: "error occurred",
		},
		{
			name:     "Error with string argument",
			format:   "invalid value: %s",
			args:     []any{"test"},
			expected: "invalid value: test",
		},
		{
			name:     "Error with integer argument",
			format:   "error code: %d",
			args:     []any{404},
			expected: "error code: 404",
		},
		{
			name:     "Error with multiple arguments",
			format:   "error at line %d, column %d: %s",
			args:     []any{10, 5, "syntax error"},
			expected: "error at line 10, column 5: syntax error",
		},
		{
			name:     "Error with float argument",
			format:   "temperature %.1f°C is too high",
			args:     []any{85.7},
			expected: "temperature 85.7°C is too high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Errorf(tt.format, tt.args...)

			// Test that it returns an error type
			if err.vTpe != tpErr {
				t.Errorf("Errorf should set vTpe to tpErr, got %v", err.vTpe)
			}

			// Test Error() method
			result := err.Error()
			if result != tt.expected {
				t.Errorf("Errorf(%q, %v).Error() = %q, want %q", tt.format, tt.args, result, tt.expected)
			}

			// Test String() method for compatibility
			resultStr := err.String()
			if resultStr != tt.expected {
				t.Errorf("Errorf(%q, %v).String() = %q, want %q", tt.format, tt.args, resultStr, tt.expected)
			}
		})
	}
}
