package tinystring

import "testing"

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		count    int
		expected string
	}{
		{
			name:     "Repeat a single character",
			input:    "x",
			count:    3,
			expected: "xxx",
		},
		{
			name:     "Repeat a word",
			input:    "hello ",
			count:    2,
			expected: "hello hello ",
		},
		{
			name:     "Zero repetitions",
			input:    "test",
			count:    0,
			expected: "",
		},
		{
			name:     "Negative repetitions",
			input:    "test",
			count:    -1,
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			count:    5,
			expected: "",
		},
		{
			name:     "Repeat once",
			input:    "once",
			count:    1,
			expected: "once",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Convert(tt.input).Repeat(tt.count).String()
			if result != tt.expected {
				t.Errorf("Convert(%q).Repeat(%d) = %q, want %q",
					tt.input, tt.count, result, tt.expected)
			}
		})
	}
}

// TestRepeatChain tests the chaining capabilities of the Repeat method with other methods
func TestRepeatChain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		function func(*conv) *conv
	}{
		{
			name:  "Repeat and convert to upper",
			input: "hello",
			want:  "HELLOHELLOHELLO",
			function: func(t *conv) *conv {
				return t.Repeat(3).ToUpper()
			},
		},
		{
			name:  "Repeat and convert to lower",
			input: "WORLD",
			want:  "worldworld",
			function: func(t *conv) *conv {
				return t.Repeat(2).ToLower()
			},
		},
		{
			name:  "Multiple operations with repeat",
			input: "Test",
			want:  "testtesttest",
			function: func(t *conv) *conv {
				return t.ToLower().Repeat(3)
			},
		},
		{
			name:  "Repeat with CamelCase",
			input: "hello world",
			want:  "helloWorldhelloWorld",
			function: func(t *conv) *conv {
				return t.CamelCaseLower().Repeat(2)
			},
		},
		{
			name:  "Empty after repeat zero",
			input: "conv",
			want:  "",
			function: func(t *conv) *conv {
				return t.Repeat(0).ToUpper()
			},
		},
		{
			name:  "Repeat with accents and remove tildes",
			input: "ñandú",
			want:  "nandunandunandu",
			function: func(t *conv) *conv {
				return t.RemoveTilde().Repeat(3)
			},
		},
		{
			name:  "SnakeCase and Repeat",
			input: "Hello World Example",
			want:  "hello_world_examplehello_world_example",
			function: func(t *conv) *conv {
				return t.ToSnakeCaseLower().Repeat(2)
			},
		},
		{
			name:  "Complex chaining",
			input: "Él Múrcielago",
			want:  "ELMURCIELAGOELMURCIELAGO",
			function: func(t *conv) *conv {
				return t.RemoveTilde().CamelCaseLower().ToUpper().Repeat(2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.function(Convert(tt.input)).String()
			if got != tt.want {
				t.Fatalf("\n🎯Test: %q\ninput: %q\n   got: %q\n  want: %q", tt.name, tt.input, got, tt.want)
			}
		})
	}
}
