package tinystring

import (
	"testing"
	"unsafe"
)

// Debug test to understand what's happening
func TestStructPointerChainDebug(t *testing.T) {
	type Inner struct {
		Value int
	}

	type Outer struct {
		InnerPtr *Inner
	}

	inner := &Inner{Value: 123}
	outer := &Outer{InnerPtr: inner}

	t.Logf("inner address: %p, value: %+v", inner, inner)
	t.Logf("outer address: %p, value: %+v", outer, outer)
	t.Logf("outer.InnerPtr address: %p", outer.InnerPtr)

	// Test accessing nested pointer through reflection
	v := refValueOf(outer)
	t.Logf("refValueOf(outer): Kind=%v, ptr=%p, flag=%v", v.Kind(), v.ptr, v.flag)

	structValue := v.Elem()
	t.Logf("v.Elem(): Kind=%v, ptr=%p, flag=%v", structValue.Kind(), structValue.ptr, structValue.flag)

	field0 := structValue.Field(0)
	t.Logf("structValue.Field(0): Kind=%v, ptr=%p, flag=%v", field0.Kind(), field0.ptr, field0.flag)

	// Check what's actually in the field pointer
	if field0.Kind() == tpPointer && field0.ptr != nil {
		actualPtr := *(*unsafe.Pointer)(field0.ptr)
		t.Logf("field0 actual pointer value: %p", actualPtr)
		t.Logf("expected pointer value: %p", inner)
	}

	innerStruct := field0.Elem()
	t.Logf("field0.Elem(): Kind=%v, ptr=%p, flag=%v", innerStruct.Kind(), innerStruct.ptr, innerStruct.flag)

	if innerStruct.IsValid() && innerStruct.ptr != nil {
		valueField := innerStruct.Field(0)
		t.Logf("innerStruct.Field(0): Kind=%v, ptr=%p, flag=%v", valueField.Kind(), valueField.ptr, valueField.flag)

		if valueField.ptr != nil {
			actualValue := *(*int)(valueField.ptr)
			t.Logf("Actual int value: %d", actualValue)
		}
	}
}
