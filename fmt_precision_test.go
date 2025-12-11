package fmt

import (
	"testing"
)

func TestRoundDecimalsUnified(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		decimals int
		down     bool // true = truncate, false = round
		want     string
	}{
		{"Int to 2 decimals", 3, 2, false, "3.00"},
		{"Float to 2 decimals", 3.12221, 2, false, "3.12"},
		{"Float to 3 decimals", 3.1415926, 3, false, "3.142"},
		{"Float to 0 decimals", 3.6, 0, false, "4"},
		{"Negative float to 2 decimals", -3.12221, 2, false, "-3.12"},
		{"String float input", "3.12221", 2, false, "3.12"},
		{"Non-numeric string input", "hello", 2, false, "0.00"},
		{"Int to 2 decimals (truncate)", 3, 2, true, "3.00"},
		{"Float to 2 decimals (truncate)", 3.14159, 2, true, "3.14"},
		{"Float to 3 decimals (truncate)", 3.1415926, 3, true, "3.141"},
		{"Float to 0 decimals (truncate)", 3.9, 0, true, "3"},
		{"Negative float to 2 decimals (truncate)", -3.12221, 2, true, "-3.12"},
		{"String float input (truncate)", "3.12221", 2, true, "3.12"},
		{"Non-numeric string input (truncate)", "hello", 2, true, "0.00"},
		// Enhanced/edge cases
		{"Round up default", "3.154", 2, false, "3.15"},
		{"Round down explicit", "3.154", 2, true, "3.15"},
		{"Round up default zero decimals", "3.7", 0, false, "4"},
		{"Round down zero decimals", "3.7", 0, true, "3"},
		{"Negative number round up", "-3.154", 2, false, "-3.15"},
		{"Negative number round down", "-3.154", 2, true, "-3.15"},
		{"Negative number with 5 - round up", "-3.155", 2, false, "-3.16"},
		{"Negative number with 5 - round down", "-3.155", 2, true, "-3.15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Convert(tt.input)
			if tt.down {
				c.Round(tt.decimals, true)
			} else {
				c.Round(tt.decimals)
			}
			out := c.String()
			if out != tt.want {
				t.Errorf("%s: got = %v, want %v", tt.name, out, tt.want)
			}
		})
	}
}
