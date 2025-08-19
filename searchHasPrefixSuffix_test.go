package tinystring

import "testing"

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		{"empty prefix", "hello", "", true},
		{"exact match", "hello", "hello", true},
		{"short prefix", "hello", "he", true},
		{"not a prefix", "hello", "ello", false},
		{"prefix longer than string", "hi", "hello", false},
		{"single byte", "abc", "a", true},
		{"null byte prefix", "a\x00b", "a\x00", true},
		{"unicode prefix", "ñandú", "ñan", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HasPrefix(tc.s, tc.prefix)
			if got != tc.want {
				t.Fatalf("HasPrefix(%q, %q) = %v; want %v", tc.s, tc.prefix, got, tc.want)
			}
		})
	}
}

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		suffix string
		want   bool
	}{
		{"empty suffix", "hello", "", true},
		{"exact match", "hello", "hello", true},
		{"short suffix", "hello", "lo", true},
		{"not a suffix", "hello", "hel", false},
		{"suffix longer than string", "go", "golang", false},
		{"single byte", "abc", "c", true},
		{"null byte suffix", "a\x00b", "\x00b", true},
		{"unicode suffix", "pingüino", "üino", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HasSuffix(tc.s, tc.suffix)
			if got != tc.want {
				t.Fatalf("HasSuffix(%q, %q) = %v; want %v", tc.s, tc.suffix, got, tc.want)
			}
		})
	}
}
