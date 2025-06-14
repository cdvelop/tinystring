package tinystring

import (
	"testing"
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
	ptrS := &s
	conv2 := Convert(ptrS)
	t.Logf("Pointer: vTpe=%s, refKind=%s", conv2.vTpe.String(), conv2.refKind().String())
	t.Logf("Pointer refIsValid: %v", conv2.refIsValid())
	t.Logf("Pointer ptr nil: %v", conv2.ptr == nil)

	// Debug step by step
	t.Logf("Step 1: c.ptr == nil? %v", conv2.ptr == nil)
	t.Logf("Step 2: c.refKind() == tpPointer? %v", conv2.refKind() == tpPointer)

	elem := conv2.refElem()
	t.Logf("Step 3: elem.refIsValid()? %v", elem.refIsValid())

	if elem.refIsValid() {
		elemValue := elem.Interface()
		t.Logf("Step 4: elemValue == nil? %v", elemValue == nil)
		t.Logf("Step 4b: elemValue type: %T", elemValue)
		t.Logf("Step 4c: elemValue: %+v", elemValue)
	}

	// Now test JSON encoding
	jsonBytes, err := conv2.JsonEncode()
	t.Logf("JSON encode err: %v", err)
	t.Logf("JSON result: %s", string(jsonBytes))
}
