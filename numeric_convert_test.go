package tinystring

import "testing"

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int
		valid    bool
	}{
		{
			name:     "int value",
			input:    42,
			expected: 42,
			valid:    true,
		},
		{
			name:     "int8 value",
			input:    int8(8),
			expected: 8,
			valid:    true,
		},
		{
			name:     "int16 value",
			input:    int16(16),
			expected: 16,
			valid:    true,
		},
		{
			name:     "int32 value",
			input:    int32(32),
			expected: 32,
			valid:    true,
		},
		{
			name:     "int64 value",
			input:    int64(64),
			expected: 64,
			valid:    true,
		},
		{
			name:     "uint value",
			input:    uint(42),
			expected: 42,
			valid:    true,
		},
		{
			name:     "uint8 value",
			input:    uint8(8),
			expected: 8,
			valid:    true,
		},
		{
			name:     "uint16 value",
			input:    uint16(16),
			expected: 16,
			valid:    true,
		},
		{
			name:     "uint32 value",
			input:    uint32(32),
			expected: 32,
			valid:    true,
		},
		{
			name:     "uint64 value",
			input:    uint64(64),
			expected: 64,
			valid:    true,
		},
		{
			name:     "float32 value",
			input:    float32(3.14),
			expected: 3,
			valid:    true,
		},
		{
			name:     "float64 value",
			input:    float64(6.28),
			expected: 6,
			valid:    true,
		},
		{
			name:     "string value (invalid)",
			input:    "not a number",
			expected: 0,
			valid:    false,
		},
		{
			name:     "boolean value (invalid)",
			input:    true,
			expected: 0,
			valid:    false,
		},
		{
			name:     "nil value (invalid)",
			input:    nil,
			expected: 0,
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := toInt(tt.input)

			if ok != tt.valid {
				t.Errorf("toInt(%v) validity = %v, want %v",
					tt.input, ok, tt.valid)
			}

			if result != tt.expected {
				t.Errorf("toInt(%v) = %v, want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}
