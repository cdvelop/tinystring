package tinystring_test

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestStringSplit(t *testing.T) {
	testCases := []struct {
		data      string
		separator string
		expected  []string
	}{
		// Original test cases
		{"texto1,texto2", ",", []string{"texto1", "texto2"}},
		{"apple,banana,cherry", ",", []string{"apple", "banana", "cherry"}},
		{"one.two.three.four", ".", []string{"one", "two", "three", "four"}},
		{"hello world", " ", []string{"hello", "world"}},
		{"hello. world", ".", []string{"hello", " world"}},
		{"h.", ".", []string{"h."}},

		// New test cases for edge conditions
		{"", ",", []string{""}},                                   // Empty string
		{"abc", "", []string{"a", "b", "c"}},                      // Empty separator
		{"ab", ",", []string{"ab"}},                               // String shorter than 3 chars
		{"a", "a", []string{"a"}},                                 // Single char string equals separator
		{"aaa", "a", []string{"", "", "", ""}},                    // All chars are separators
		{"abc,def,", ",", []string{"abc", "def", ""}},             // Trailing separator
		{",abc,def", ",", []string{"", "abc", "def"}},             // Leading separator
		{"abc,,def", ",", []string{"abc", "", "def"}},             // Adjacent separators
		{"abc", "abc", []string{"", ""}},                          // Separator equals data
		{"abcdefghi", "def", []string{"abc", "ghi"}},              // Separator in the middle
		{"abc:::def:::ghi", ":::", []string{"abc", "def", "ghi"}}, // Multi-char separator
	}

	for _, tc := range testCases {
		result := Split(tc.data, tc.separator)

		if !areStringSlicesEqual(result, tc.expected) {
			t.Errorf("StringSplit(%q, %q) = %v; expected %v", tc.data, tc.separator, result, tc.expected)
		}
	}
}

func areStringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
