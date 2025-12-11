package fmt

import (
	"testing"
)

func TestSscanf(t *testing.T) {
	tests := []struct {
		name        string
		src         string
		format      string
		args        []any
		expected    int
		shouldError bool
		validate    func(t *testing.T, args []any)
	}{
		// Decimal integer parsing
		{
			name:     "Parse single decimal integer",
			src:      "42",
			format:   "%d",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 42 {
					t.Errorf("Expected 42, got %d", val)
				}
			},
		},
		{
			name:     "Parse negative decimal integer",
			src:      "-123",
			format:   "%d",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != -123 {
					t.Errorf("Expected -123, got %d", val)
				}
			},
		},
		{
			name:     "Parse int64",
			src:      "9223372036854775807",
			format:   "%d",
			args:     []any{new(int64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int64); val != 9223372036854775807 {
					t.Errorf("Expected 9223372036854775807, got %d", val)
				}
			},
		},
		{
			name:     "Parse int32",
			src:      "2147483647",
			format:   "%d",
			args:     []any{new(int32)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int32); val != 2147483647 {
					t.Errorf("Expected 2147483647, got %d", val)
				}
			},
		},

		// Hexadecimal parsing
		{
			name:     "Parse hexadecimal lowercase",
			src:      "ff",
			format:   "%x",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 255 {
					t.Errorf("Expected 255, got %d", val)
				}
			},
		},
		{
			name:     "Parse hexadecimal uppercase",
			src:      "FF",
			format:   "%X",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 255 {
					t.Errorf("Expected 255, got %d", val)
				}
			},
		},
		{
			name:     "Parse hexadecimal mixed case",
			src:      "AbC",
			format:   "%x",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 2748 {
					t.Errorf("Expected 2748, got %d", val)
				}
			},
		},
		{
			name:     "Parse hex to uint64",
			src:      "1a2b3c",
			format:   "%x",
			args:     []any{new(uint64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*uint64); val != 1715004 {
					t.Errorf("Expected 1715004, got %d", val)
				}
			},
		},

		// Float parsing
		{
			name:     "Parse float64",
			src:      "3.14159",
			format:   "%f",
			args:     []any{new(float64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				val := *args[0].(*float64)
				expected := 3.14159
				if val < expected-0.0001 || val > expected+0.0001 {
					t.Errorf("Expected %f, got %f", expected, val)
				}
			},
		},
		{
			name:     "Parse negative float",
			src:      "-2.718",
			format:   "%f",
			args:     []any{new(float64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				val := *args[0].(*float64)
				expected := -2.718
				if val < expected-0.0001 || val > expected+0.0001 {
					t.Errorf("Expected %f, got %f", expected, val)
				}
			},
		},
		{
			name:     "Parse float32",
			src:      "1.5",
			format:   "%f",
			args:     []any{new(float32)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*float32); val != 1.5 {
					t.Errorf("Expected 1.5, got %f", val)
				}
			},
		},
		{
			name:     "Parse float with %g format",
			src:      "123.456",
			format:   "%g",
			args:     []any{new(float64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				val := *args[0].(*float64)
				expected := 123.456
				if val < expected-0.0001 || val > expected+0.0001 {
					t.Errorf("Expected %f, got %f", expected, val)
				}
			},
		},
		{
			name:     "Parse float with %e format",
			src:      "2.5",
			format:   "%e",
			args:     []any{new(float64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				val := *args[0].(*float64)
				expected := 2.5
				if val < expected-0.0001 || val > expected+0.0001 {
					t.Errorf("Expected %f, got %f", expected, val)
				}
			},
		},

		// String parsing
		{
			name:     "Parse string",
			src:      "hello",
			format:   "%s",
			args:     []any{new(string)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*string); val != "hello" {
					t.Errorf("Expected 'hello', got '%s'", val)
				}
			},
		},
		{
			name:     "Parse string with spaces (stops at first space)",
			src:      "hello world",
			format:   "%s",
			args:     []any{new(string)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*string); val != "hello" {
					t.Errorf("Expected 'hello', got '%s'", val)
				}
			},
		},

		// Character parsing
		{
			name:     "Parse character to rune",
			src:      "A",
			format:   "%c",
			args:     []any{new(rune)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*rune); val != 'A' {
					t.Errorf("Expected 'A', got '%c'", val)
				}
			},
		},
		{
			name:     "Parse character to byte",
			src:      "Z",
			format:   "%c",
			args:     []any{new(byte)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*byte); val != 'Z' {
					t.Errorf("Expected 'Z', got '%c'", val)
				}
			},
		},

		// Complex pattern from example in docstring
		{
			name:     "Complex Unicode pattern",
			src:      "!3F U+003F question",
			format:   "!%x U+%x %s",
			args:     []any{new(int), new(int), new(string)},
			expected: 3,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 0x3F {
					t.Errorf("Expected 63 (0x3F), got %d", val)
				}
				if val := *args[1].(*int); val != 0x003F {
					t.Errorf("Expected 63 (0x003F), got %d", val)
				}
				if val := *args[2].(*string); val != "question" {
					t.Errorf("Expected 'question', got '%s'", val)
				}
			},
		},

		// Multiple values
		{
			name:     "Parse multiple integers",
			src:      "123 456 789",
			format:   "%d %d %d",
			args:     []any{new(int), new(int), new(int)},
			expected: 3,
			validate: func(t *testing.T, args []any) {
				expected := []int{123, 456, 789}
				for i, exp := range expected {
					if val := *args[i].(*int); val != exp {
						t.Errorf("Expected %d at position %d, got %d", exp, i, val)
					}
				}
			},
		},
		{
			name:     "Parse mixed types",
			src:      "Name: John Age: 30 Height: 5.9",
			format:   "Name: %s Age: %d Height: %f",
			args:     []any{new(string), new(int), new(float64)},
			expected: 3,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*string); val != "John" {
					t.Errorf("Expected 'John', got '%s'", val)
				}
				if val := *args[1].(*int); val != 30 {
					t.Errorf("Expected 30, got %d", val)
				}
				if val := *args[2].(*float64); val < 5.9-0.0001 || val > 5.9+0.0001 {
					t.Errorf("Expected 5.9, got %f", val)
				}
			},
		},

		// Percent literal - now fully supported
		{
			name:     "Parse with percent literal (full support)",
			src:      "100% done",
			format:   "%d%% %s",
			args:     []any{new(int), new(string)},
			expected: 2, // Now parses both parts correctly
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 100 {
					t.Errorf("Expected 100, got %d", val)
				}
				if val := *args[1].(*string); val != "done" {
					t.Errorf("Expected 'done', got '%s'", val)
				}
			},
		},

		// Simple percent literal test
		{
			name:     "Parse simple percent literal",
			src:      "100%",
			format:   "%d%%",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 100 {
					t.Errorf("Expected 100, got %d", val)
				}
			},
		},

		// Edge cases and error conditions
		{
			name:        "Not enough arguments for format specifiers",
			src:         "42 hello 99",
			format:      "%d %s %d %d", // 4 format specifiers, only 3 args
			args:        []any{new(int), new(string), new(int)},
			expected:    3,    // Should parse all available arguments
			shouldError: true, // Should error when trying to parse 4th specifier
		},
		{
			name:        "Invalid format character",
			src:         "42",
			format:      "%z",
			args:        []any{new(int)},
			expected:    0,
			shouldError: true,
		},
		{
			name:        "Literal mismatch",
			src:         "hello world",
			format:      "hello universe",
			args:        []any{},
			expected:    0,
			shouldError: true,
		},
		{
			name:     "Partial parse - stops at first failure",
			src:      "123 abc 456",
			format:   "%d %d %d",
			args:     []any{new(int), new(int), new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 123 {
					t.Errorf("Expected 123, got %d", val)
				}
			},
		},
		{
			name:     "Empty string",
			src:      "",
			format:   "",
			args:     []any{},
			expected: 0,
		},
		{
			name:     "Only literals",
			src:      "hello world",
			format:   "hello world",
			args:     []any{},
			expected: 0,
		},

		// Special numeric cases
		{
			name:     "Parse zero",
			src:      "0",
			format:   "%d",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 0 {
					t.Errorf("Expected 0, got %d", val)
				}
			},
		},
		{
			name:     "Parse zero hex",
			src:      "0",
			format:   "%x",
			args:     []any{new(int)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*int); val != 0 {
					t.Errorf("Expected 0, got %d", val)
				}
			},
		},
		{
			name:     "Parse zero float",
			src:      "0.0",
			format:   "%f",
			args:     []any{new(float64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				val := *args[0].(*float64)
				if val < -0.0001 || val > 0.0001 {
					t.Errorf("Expected 0.0, got %f", val)
				}
			},
		},

		// Large numbers
		{
			name:     "Parse large hex number",
			src:      "deadbeef",
			format:   "%x",
			args:     []any{new(uint64)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*uint64); val != 0xdeadbeef {
					t.Errorf("Expected 3735928559 (0xdeadbeef), got %d", val)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n, err := Sscanf(test.src, test.format, test.args...)

			// Check parsed count
			if n != test.expected {
				t.Errorf("Expected %d items parsed, got %d", test.expected, n)
			}

			// Check error expectation
			if test.shouldError && err == nil {
				t.Error("Expected an error but got none")
			} else if !test.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Run custom validation if provided
			if test.validate != nil && !test.shouldError {
				test.validate(t, test.args)
			}
		})
	}
}

func TestSscanfTypeValidation(t *testing.T) {
	// Test type validation for different format specifiers
	tests := []struct {
		name        string
		src         string
		format      string
		arg         any
		shouldError bool
	}{
		{
			name:        "Integer to wrong type pointer",
			src:         "42",
			format:      "%d",
			arg:         new(string), // Wrong type
			shouldError: true,
		},
		{
			name:        "Float to wrong type pointer",
			src:         "3.14",
			format:      "%f",
			arg:         new(int), // Wrong type
			shouldError: true,
		},
		{
			name:        "String to wrong type pointer",
			src:         "hello",
			format:      "%s",
			arg:         new(int), // Wrong type
			shouldError: true,
		},
		{
			name:        "Character to wrong type pointer",
			src:         "A",
			format:      "%c",
			arg:         new(int), // Wrong type (not rune or byte)
			shouldError: true,
		},
		{
			name:        "Hex to wrong type pointer",
			src:         "ff",
			format:      "%x",
			arg:         new(string), // Wrong type
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n, err := Sscanf(test.src, test.format, test.arg)

			if test.shouldError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				if n != 0 {
					t.Errorf("Expected 0 items parsed on error, got %d", n)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if n != 1 {
					t.Errorf("Expected 1 item parsed, got %d", n)
				}
			}
		})
	}
}

func TestSscanfWhitespaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		format   string
		args     []any
		expected int
		validate func(t *testing.T, args []any)
	}{
		{
			name:     "String parsing stops at whitespace",
			src:      "first second third",
			format:   "%s",
			args:     []any{new(string)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*string); val != "first" {
					t.Errorf("Expected 'first', got '%s'", val)
				}
			},
		},
		{
			name:     "Multiple strings separated by spaces",
			src:      "one two three",
			format:   "%s %s %s",
			args:     []any{new(string), new(string), new(string)},
			expected: 3,
			validate: func(t *testing.T, args []any) {
				expected := []string{"one", "two", "three"}
				for i, exp := range expected {
					if val := *args[i].(*string); val != exp {
						t.Errorf("Expected '%s' at position %d, got '%s'", exp, i, val)
					}
				}
			},
		},
		{
			name:     "String with tabs and newlines",
			src:      "word1\tword2\nword3",
			format:   "%s",
			args:     []any{new(string)},
			expected: 1,
			validate: func(t *testing.T, args []any) {
				if val := *args[0].(*string); val != "word1" {
					t.Errorf("Expected 'word1', got '%s'", val)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n, err := Sscanf(test.src, test.format, test.args...)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if n != test.expected {
				t.Errorf("Expected %d items parsed, got %d", test.expected, n)
			}

			if test.validate != nil {
				test.validate(t, test.args)
			}
		})
	}
}

// Test specific to PDF library font encoding usage
func TestSscanfPDFFontEncoding(t *testing.T) {
	// This test simulates the exact usage pattern in the PDF library's loadMap function
	tests := []struct {
		name     string
		line     string
		expected struct {
			pos  int
			uv   int
			name string
		}
		shouldError bool
	}{
		{
			name: "PDF font encoding line - question mark",
			line: "!3F U+003F question",
			expected: struct {
				pos  int
				uv   int
				name string
			}{pos: 0x3F, uv: 0x003F, name: "question"},
		},
		{
			name: "PDF font encoding line - exclamation",
			line: "!21 U+0021 exclam",
			expected: struct {
				pos  int
				uv   int
				name string
			}{pos: 0x21, uv: 0x0021, name: "exclam"},
		},
		{
			name: "PDF font encoding line - space",
			line: "!20 U+0020 space",
			expected: struct {
				pos  int
				uv   int
				name string
			}{pos: 0x20, uv: 0x0020, name: "space"},
		},
		{
			name: "PDF font encoding line - A uppercase",
			line: "!41 U+0041 A",
			expected: struct {
				pos  int
				uv   int
				name string
			}{pos: 0x41, uv: 0x0041, name: "A"},
		},
		{
			name: "PDF font encoding line - complex name",
			line: "!E9 U+00E9 eacute",
			expected: struct {
				pos  int
				uv   int
				name string
			}{pos: 0xE9, uv: 0x00E9, name: "eacute"},
		},
		{
			name:        "Invalid format - missing U+",
			line:        "!3F 003F question",
			shouldError: true,
		},
		{
			name:        "Invalid format - missing !",
			line:        "3F U+003F question",
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var pos int
			var uv int
			var name string

			// This is exactly how it's used in the PDF library
			n, err := Sscanf(test.line, "!%x U+%x %s", &pos, &uv, &name)

			if test.shouldError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if n != 3 {
				t.Errorf("Expected 3 items parsed, got %d", n)
			}

			if pos != test.expected.pos {
				t.Errorf("Expected pos %d (0x%X), got %d (0x%X)", test.expected.pos, test.expected.pos, pos, pos)
			}

			if uv != test.expected.uv {
				t.Errorf("Expected uv %d (0x%X), got %d (0x%X)", test.expected.uv, test.expected.uv, uv, uv)
			}

			if name != test.expected.name {
				t.Errorf("Expected name '%s', got '%s'", test.expected.name, name)
			}

			// Additional validation that pos is within valid range for PDF encoding
			if pos < 0 || pos > 255 {
				t.Errorf("Position %d (0x%X) is outside valid range [0-255]", pos, pos)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSscanf(b *testing.B) {
	benchmarks := []struct {
		name   string
		src    string
		format string
		args   []any
	}{
		{
			name:   "Single integer",
			src:    "42",
			format: "%d",
			args:   []any{new(int)},
		},
		{
			name:   "Multiple integers",
			src:    "123 456 789",
			format: "%d %d %d",
			args:   []any{new(int), new(int), new(int)},
		},
		{
			name:   "Mixed types",
			src:    "John 30 5.9",
			format: "%s %d %f",
			args:   []any{new(string), new(int), new(float64)},
		},
		{
			name:   "Hex parsing",
			src:    "deadbeef",
			format: "%x",
			args:   []any{new(uint64)},
		},
		{
			name:   "Complex pattern",
			src:    "!3F U+003F question",
			format: "!%x U+%x %s",
			args:   []any{new(int), new(int), new(string)},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = Sscanf(bm.src, bm.format, bm.args...)
			}
		})
	}
}
