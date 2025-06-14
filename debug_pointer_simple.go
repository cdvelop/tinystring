package tinystring

import (
	"testing"
)

func TestDebugPointerTypeDetection(t *testing.T) {
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

	// Now test JSON encoding
	jsonBytes, err := conv2.JsonEncode()
	t.Logf("JSON encode err: %v", err)
	t.Logf("JSON result: %s", string(jsonBytes))

	// Let's also test what happens when we call generateJsonBytes directly
	directBytes, directErr := conv2.generateJsonBytes()
	t.Logf("Direct generateJsonBytes err: %v", directErr)
	t.Logf("Direct result: %s", string(directBytes))
}
