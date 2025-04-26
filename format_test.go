package tinystring

import "testing"

func TestRoundDecimals(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		decimals int
		want     string
	}{
		{
			name:     "Round to 2 decimals",
			input:    3.12221,
			decimals: 2,
			want:     "3.12",
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
		},
		{
			name:     "Round negative to 2 decimals",
			input:    -3.12221,
			decimals: 2,
			want:     "-3.12",
		},
		{
			name:     "Round without decimal",
			input:    3,
			decimals: 2,
			want:     "3.00",
		},
		{
			name:     "Round string input",
			input:    "3.12221",
			decimals: 2,
			want:     "3.12",
		},
		{
			name:     "Non-numeric input",
			input:    "hello",
			decimals: 2,
			want:     "hello",
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
