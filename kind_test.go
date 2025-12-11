package fmt

import "testing"

func TestKindString(t *testing.T) {
	// Build a slice of all Kind values from the Kind struct, matching the new order
	kindVals := []Kind{
		K.Invalid, K.Bool, K.Int, K.Int8, K.Int16, K.Int32, K.Int64, K.Uint, K.Uint8, K.Uint16, K.Uint32, K.Uint64, K.Uintptr,
		K.Float32, K.Float64, K.Complex64, K.Complex128, K.Array, K.Chan, K.Func, K.Interface, K.Map, K.Pointer, K.Slice, K.String, K.Struct, K.UnsafePointer,
	}
	for i, want := range kindNames {
		k := kindVals[i]
		got := k.String()
		if got != want {
			t.Errorf("K.%s.String() = %q, want %q", want, got, want)
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
