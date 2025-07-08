package tinystring

import "testing"

func TestKindString(t *testing.T) {
	for i, want := range kindNames {
		k := kind(i)
		got := k.String()
		if got != want {
			t.Errorf("kind(%d).String() = %q, want %q", i, got, want)
		}
	}

	// Test out-of-range (too large)
	invalids := []kind{kind(len(kindNames)), 255}
	for _, k := range invalids {
		if k.String() != "invalid" {
			t.Errorf("kind(%d).String() = %q, want 'invalid'", int(k), k.String())
		}
	}
}
