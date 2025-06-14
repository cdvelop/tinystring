package tinystring

import (
	"testing"
	"unsafe"
)

func TestDebugPointerEncoding(t *testing.T) {
	type SimpleStruct struct {
		Value int
		Name  string
	}

	s := SimpleStruct{
		Value: 42,
		Name:  "test",
	}

	// Test pointer to struct
	// Note: Convert() automatically dereferences pointers for convenience
	ptrS := &s
	conv2 := Convert(ptrS)
	t.Logf("After Convert(&struct): vTpe=%s, refKind=%s", conv2.vTpe.String(), conv2.refKind().String())
	t.Logf("refIsValid: %v", conv2.refIsValid())
	t.Logf("ptr nil: %v", conv2.ptr == nil)

	// Debug: Convert() automatically dereferences pointers
	// So conv2 now contains the struct, not the pointer
	t.Logf("Step 1: After Convert(), kind is struct (auto-dereferenced): %v", conv2.refKind() == tpStruct)
	// Test that we can access struct fields
	if conv2.refKind() == tpStruct && conv2.refIsValid() {
		numFields := conv2.refNumField()
		t.Logf("Step 2: Number of struct fields: %d", numFields)

		// Access individual fields
		tt := (*refStructType)(unsafe.Pointer(conv2.Type()))
		for i := 0; i < numFields; i++ {
			field := conv2.refField(i)
			fieldInfo := tt.fields[i]
			fieldName := fieldInfo.name.Name()
			t.Logf("Step 3.%d: Field %s (type %s)", i, fieldName, field.refKind().String())
		}
	}

	// Test JSON encoding of the dereferenced struct
	jsonBytes, err := conv2.JsonEncode()
	t.Logf("JSON encode err: %v", err)
	t.Logf("JSON result: %s", string(jsonBytes))

	// Verify JSON contains expected values
	expectedJSON := `{"Value":42,"Name":"test"}`
	if err == nil && string(jsonBytes) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(jsonBytes))
	}
}
