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
			name:  "Fmt integer with thousand separators",
			input: 2189009,
			want:  "2.189.009",
		},
		{
			name:  "Fmt decimal number with trailing zeros",
			input: 2189009.00,
			want:  "2.189.009",
		},
		{
			name:  "Fmt decimal number",
			input: 2189009.123,
			want:  "2.189.009.123",
		},
		{
			name:  "Fmt string number",
			input: "2189009.00",
			want:  "2.189.009",
		},
		{
			name:  "Fmt negative number",
			input: -2189009,
			want:  "-2.189.009",
		},
		{
			name:  "Fmt small number",
			input: 123,
			want:  "123",
		},
		{
			name:  "Fmt zero",
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
			result := Fmt(test.format, test.args...).String()
			resultWithError, err := Fmt(test.format, test.args...).StringError()

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

func TestReporterFormatting(t *testing.T) {
	// Test cases based on actual usage in reporter.go
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
		hasError bool
	}{
		{
			name:     "Peak reduction percentage",
			format:   "- 🏆 **Peak Reduction: %.1f%%** (Best optimization)\n",
			args:     []any{71.5},
			expected: "- 🏆 **Peak Reduction: 71.5%** (Best optimization)\n",
			hasError: false,
		},
		{
			name:     "Average WebAssembly reduction",
			format:   "- ✅ **Average WebAssembly Reduction: %.1f%%**\n",
			args:     []any{53.2},
			expected: "- ✅ **Average WebAssembly Reduction: 53.2%**\n",
			hasError: false,
		},
		{
			name:     "Size savings with string",
			format:   "- 📦 **Total Size Savings: %s across all builds**\n\n",
			args:     []any{"1.7 MB"},
			expected: "- 📦 **Total Size Savings: 1.7 MB across all builds**\n\n",
			hasError: false,
		},
		{
			name:     "Memory efficiency class",
			format:   "- 💾 **Memory Efficiency**: %s (%.1f%% average change)\n",
			args:     []any{"❌ **Poor** (Significant overhead)", 154.2},
			expected: "- 💾 **Memory Efficiency**: ❌ **Poor** (Significant overhead) (154.2% average change)\n",
			hasError: false,
		},
		{
			name:     "Allocation efficiency class",
			format:   "- 🔢 **Allocation Efficiency**: %s (%.1f%% average change)\n",
			args:     []any{"❌ **Poor** (Excessive allocations)", 118.4},
			expected: "- 🔢 **Allocation Efficiency**: ❌ **Poor** (Excessive allocations) (118.4% average change)\n",
			hasError: false,
		},
		{
			name:     "Benchmarks analyzed count",
			format:   "- 📊 **Benchmarks Analyzed**: %d categories\n",
			args:     []any{3},
			expected: "- 📊 **Benchmarks Analyzed**: 3 categories\n",
			hasError: false,
		},
		{
			name:     "Complex table row with multiple formats",
			format:   "| %s **%s** | 📊 Standard | `%s` | `%d` | `%s` | - | - | - |\n",
			args:     []any{"📝", "String Processing", "1.2 KB", 48, "3.4μs"},
			expected: "| 📝 **String Processing** | 📊 Standard | `1.2 KB` | `48` | `3.4μs` | - | - | - |\n",
			hasError: false,
		},
		{
			name:     "TinyString performance row",
			format:   "| | 🚀 TinyString | `%s` | `%d` | `%s` | %s **%s** | %s **%s** | %s |\n",
			args:     []any{"2.8 KB", 119, "13.7μs", "❌", "140.3% more", "❌", "147.9% more", "❌ **Poor**"},
			expected: "| | 🚀 TinyString | `2.8 KB` | `119` | `13.7μs` | ❌ **140.3% more** | ❌ **147.9% more** | ❌ **Poor** |\n",
			hasError: false,
		},
		{
			name:     "Binary size table row",
			format:   "| %s **%s Native** | `%s` | %s | %s | **-%s** | %s **%.1f%%** |\n",
			args:     []any{"🖥️", "Default", "-ldflags=\"-s -w\"", "1.3 MB", "1.1 MB", "176.0 KB", "➖", 13.4},
			expected: "| 🖥️ **Default Native** | `-ldflags=\"-s -w\"` | 1.3 MB | 1.1 MB | **-176.0 KB** | ➖ **13.4%** |\n",
			hasError: false,
		},
		{
			name:     "Error message formatting",
			format:   "Failed to read README: %v",
			args:     []any{Err("file not found")},
			expected: "Failed to read README: file not found",
			hasError: false,
		},
		{
			name:     "Memory improvement percentage",
			format:   "%.1f%% less",
			args:     []any{44.2},
			expected: "44.2% less",
			hasError: false,
		},
		{
			name:     "Memory improvement percentage more",
			format:   "%.1f%% more",
			args:     []any{140.3},
			expected: "140.3% more",
			hasError: false,
		},
		{
			name:     "Nanosecond formatting",
			format:   "%dns",
			args:     []any{int64(500)},
			expected: "500ns",
			hasError: false,
		},
		{
			name:     "Microsecond formatting",
			format:   "%.1fμs",
			args:     []any{3.4},
			expected: "3.4μs",
			hasError: false,
		},
		{
			name:     "Millisecond formatting",
			format:   "%.1fms",
			args:     []any{1.5},
			expected: "1.5ms",
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Fmt(test.format, test.args...).String()
			resultWithError, err := Fmt(test.format, test.args...).StringError()

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
