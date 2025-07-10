package tinystring

import "testing"

func TestKindString(t *testing.T) {
	// Build a slice of all kind values from the Kind struct
	kindVals := []kind{
		Kind.Array, Kind.Bool, Kind.Bytes, Kind.Chan, Kind.Complex128, Kind.Complex64, Kind.Float32, Kind.Float64, Kind.Func,
		Kind.Int, Kind.Int16, Kind.Int32, Kind.Int64, Kind.Int8, Kind.Interface, Kind.Invalid, Kind.Map, Kind.Pointer, Kind.Slice,
		Kind.String, Kind.Struct, Kind.Uint, Kind.Uint16, Kind.Uint32, Kind.Uint64, Kind.Uint8, Kind.Uintptr, Kind.UnsafePtr,
	}
	for i, want := range kindNames {
		k := kindVals[i]
		got := k.String()
		if got != want {
			t.Errorf("Kind.%s.String() = %q, want %q", want, got, want)
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
