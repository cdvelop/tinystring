package tinystring

import "testing"

func TestKindString(t *testing.T) {
	for i, want := range kindNames {
		k := Kind(i)
		got := k.String()
		if got != want {
			t.Errorf("Kind(%d).String() = %q, want %q", i, got, want)
		}
	}

	// Test out-of-range (too large)
	invalids := []Kind{Kind(len(kindNames)), 255}
	for _, k := range invalids {
		if k.String() != "invalid" {
			t.Errorf("Kind(%d).String() = %q, want 'invalid'", int(k), k.String())
		}
	}
}
