package tinystring

import "testing"

func TestKindString(t *testing.T) {
	// Build a slice of all Kind values from the Kind struct
	kindVals := []Kind{
		K.Array, K.Bool, K.Bytes, K.Chan, K.Complex128, K.Complex64, K.Float32, K.Float64, K.Func,
		K.Int, K.Int16, K.Int32, K.Int64, K.Int8, K.Interface, K.Invalid, K.Map, K.Pointer, K.Slice,
		K.String, K.Struct, K.Uint, K.Uint16, K.Uint32, K.Uint64, K.Uint8, K.Uintptr, K.UnsafePtr,
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
