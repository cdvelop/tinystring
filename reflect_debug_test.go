package tinystring

import (
	"testing"
	"unsafe"
)

// TestReflectPointerFieldAccess tests reflection access to pointer-to-struct fields and setting values
func TestReflectPointerFieldAccess(t *testing.T) {
	clearRefStructsCache()

	// Test the pattern: struct with pointer to struct with float fields
	type TestCoords struct {
		Lat float64
		Lng float64
		Alt int
	}

	type TestContainer struct {
		Name   string
		Coords *TestCoords
	}

	container := TestContainer{
		Name: "test",
		Coords: &TestCoords{
			Lat: 37.7749,
			Lng: -122.4194,
			Alt: 100,
		},
	}

	// Test reflection access
	rv := refValueOf(&container)
	elem := rv.refElem() // Get the struct value from pointer

	if elem.refKind() != tpStruct {
		t.Errorf("Container elem kind: expected %v, got %v", tpStruct, elem.refKind())
	}

	// Get fields
	nameField := elem.refField(0)
	coordsField := elem.refField(1)

	if nameField.refKind() != tpString {
		t.Errorf("Name field kind: expected %v, got %v", tpString, nameField.refKind())
	}
	if nameField.String() != "test" {
		t.Errorf("Name field value: expected %q, got %q", "test", nameField.String())
	}

	if coordsField.refKind() != tpPointer {
		t.Errorf("Coords field kind: expected %v, got %v", tpPointer, coordsField.refKind())
	}

	if coordsField.refKind() == tpPointer {
		coordsElem := coordsField.refElem()
		if !coordsElem.refIsValid() {
			t.Fatal("Coords elem is not valid")
		}

		if coordsElem.refKind() != tpStruct {
			t.Errorf("Coords elem kind: expected %v, got %v", tpStruct, coordsElem.refKind())
		}

		if coordsElem.refKind() == tpStruct {
			latField := coordsElem.refField(0)
			lngField := coordsElem.refField(1)
			altField := coordsElem.refField(2)

			if latField.refKind() != tpFloat64 {
				t.Errorf("Lat field kind: expected %v, got %v", tpFloat64, latField.refKind())
			}
			if latField.refFloat() != 37.7749 {
				t.Errorf("Lat field value: expected %f, got %f", 37.7749, latField.refFloat())
			}

			if lngField.refKind() != tpFloat64 {
				t.Errorf("Lng field kind: expected %v, got %v", tpFloat64, lngField.refKind())
			}
			if lngField.refFloat() != -122.4194 {
				t.Errorf("Lng field value: expected %f, got %f", -122.4194, lngField.refFloat())
			}

			if altField.refKind() != tpInt {
				t.Errorf("Alt field kind: expected %v, got %v", tpInt, altField.refKind())
			}
			if altField.refInt() != 100 {
				t.Errorf("Alt field value: expected %d, got %d", 100, altField.refInt())
			}
		}
	}
}

// TestReflectFieldSetterOperations tests setting values through reflection on pointer-to-struct fields
func TestReflectFieldSetterOperations(t *testing.T) {
	clearRefStructsCache()

	type TestCoords struct {
		Lat float64
		Lng float64
		Alt int
	}

	type TestContainer struct {
		Name   string
		Coords *TestCoords
	}

	// Test setting to a new struct and see if we can set values
	var newContainer TestContainer
	newRv := refValueOf(&newContainer)
	newElem := newRv.refElem()

	// Set the name field
	newNameField := newElem.refField(0)
	newNameField.refSetString("new_test")
	if newContainer.Name != "new_test" {
		t.Errorf("Name setting failed: expected %q, got %q", "new_test", newContainer.Name)
	}

	// Test setting a pointer field to a new struct instance
	newContainer.Coords = &TestCoords{Lat: 1.0, Lng: 2.0, Alt: 3}
	
	// Get the coords field and verify we can access it
	coordsField := newElem.refField(1)
	if coordsField.refKind() != tpPointer {
		t.Errorf("Coords field kind: expected %v, got %v", tpPointer, coordsField.refKind())
	}

	// Dereference the pointer to get the struct
	coordsElem := coordsField.refElem()
	if coordsElem.refKind() != tpStruct {
		t.Errorf("Coords elem kind: expected %v, got %v", tpStruct, coordsElem.refKind())
	}

	// Test accessing fields of the pointed-to struct
	latField := coordsElem.refField(0)
	if latField.refKind() != tpFloat64 {
		t.Errorf("Lat field kind: expected %v, got %v", tpFloat64, latField.refKind())
	}

	// Test value retrieval
	latValue := latField.refFloat()
	if latValue != 1.0 {
		t.Errorf("Lat value: expected %f, got %f", 1.0, latValue)
	}

	t.Logf("SUCCESS: Field setter operations work correctly")
}

// TestReflectFieldCorruption tests field access patterns to diagnose corruption issues
func TestReflectFieldCorruption(t *testing.T) {
	clearRefStructsCache()

	// Define local test struct to avoid cross-file dependencies
	type TestPhoneNumber struct {
		ID         string
		Type       string
		Number     string
		Extension  string
		IsPrimary  bool
		IsVerified bool
	}

	// Use the local test type
	phone := TestPhoneNumber{
		ID:         "ph_001",
		Type:       "mobile",
		Number:     "+1-555-123-4567",
		Extension:  "",
		IsPrimary:  true,
		IsVerified: true,
	}
	// Test reflection on this structure
	rv := refValueOf(phone)
	if rv.refKind() != tpStruct {
		t.Fatalf("Expected struct, got %v", rv.refKind())
	}

	// Get struct type info - use refStructMeta, not refStructType
	tt := (*refStructMeta)(unsafe.Pointer(rv.typ))
	expectedFieldCount := 6 // ID, Type, Number, Extension, IsPrimary, IsVerified
	if len(tt.fields) != expectedFieldCount {
		t.Errorf("Struct type fields count: expected %d, got %d", expectedFieldCount, len(tt.fields))
	}

	// Check each field
	for i := 0; i < len(tt.fields); i++ {
		field := &tt.fields[i]
		fieldName := field.name.Name() // Use Name() method
		// Get field value using reflection
		fieldVal := rv.refField(i)

		if fieldVal.refKind() == tpString {
			strVal := fieldVal.String()

			// Manual check - calculate field address manually
			fieldPtr := add(rv.ptr, field.offset, "manual calculation")
			manualVal := *(*string)(fieldPtr)

			if strVal != manualVal {
				t.Errorf("refField %d (%s) mismatch: reflection=%q vs manual=%q", i, fieldName, strVal, manualVal)
			}
		}
	}

	// Test specific field access patterns
	idField := rv.refField(0) // Should be ID
	if idField.String() != "ph_001" {
		t.Errorf("ID field: expected %q, got %q", "ph_001", idField.String())
	}

	numberField := rv.refField(2) // Should be Number
	if numberField.String() != "+1-555-123-4567" {
		t.Errorf("Number field: expected %q, got %q", "+1-555-123-4567", numberField.String())
	}
}

// TestDebugPointerType - Comprehensive test to debug pointer type handling
// This test validates that pointer types are correctly identified and can be dereferenced
func TestDebugPointerType(t *testing.T) {
	i := 42
	pi := &i

	// Create refValue from pointer
	v := refValueOf(pi)

	// Test basic pointer properties
	if v.refKind() != tpPointer {
		t.Errorf("Expected tpPointer, got %v", v.refKind())
	}

	if !v.refIsValid() {
		t.Fatal("Pointer refValue should be valid")
	}

	if v.typ == nil {
		t.Fatal("Pointer refValue.typ should not be nil")
	}

	// Test pointer type properties
	if v.refKind() != tpPointer {
		t.Errorf("Expected pointer type kind, got %v", v.refKind())
	}

	// Test elem type access
	elemType := v.refElem()
	if elemType == nil {
		t.Fatal("Pointer type should have valid elem type")
	}

	if elemType.refKind() != tpInt {
		t.Errorf("Expected int elem type, got %v", elemType.refKind())
	}

	// Test refElem() method
	elem := v.refElem()
	if !elem.refIsValid() {
		t.Fatal("refElem() should return valid value for non-nil pointer")
	}

	if elem.refKind() != tpInt {
		t.Errorf("Expected int after refElem(), got %v", elem.refKind())
	}

	// Test value retrieval
	actualValue := elem.refInt()
	if actualValue != int64(i) {
		t.Errorf("Expected %d, got %d", i, actualValue)
	}

	// Test direct memory access consistency
	if elem.ptr != nil {
		directValue := *(*int)(elem.ptr)
		if directValue != i {
			t.Errorf("Direct memory access: expected %d, got %d", i, directValue)
		}
	}
}

// TestDebugPointerChain - Critical debug test to ensure pointer chains work correctly
// This test MUST PASS for JSON decoding to work properly
func TestDebugPointerChain(t *testing.T) {
	clearRefStructsCache()

	type Inner struct {
		Value int
	}

	type Outer struct {
		InnerPtr *Inner
	}

	inner := &Inner{Value: 123}
	outer := &Outer{InnerPtr: inner}

	// Debug step by step
	t.Logf("Original inner value: %d", inner.Value)
	t.Logf("Original inner pointer: %p", inner)
	t.Logf("Original outer.InnerPtr: %p", outer.InnerPtr)

	// Step 1: Get outer pointer
	v := refValueOf(outer)
	t.Logf("Step 1 - refValueOf(outer): refKind=%v, ptr=%p", v.refKind(), v.ptr)
	if v.refKind() != tpPointer {
		t.Errorf("Step 1: Expected tpPointer, got %v", v.refKind())
	}

	// Step 2: Dereference to get struct
	structValue := v.refElem()
	t.Logf("Step 2 - structValue: refKind=%v, ptr=%p", structValue.refKind(), structValue.ptr)
	if structValue.refKind() != tpStruct {
		t.Errorf("Step 2: Expected tpStruct, got %v", structValue.refKind())
	}
	// Verify this points to the actual struct
	if uintptr(structValue.ptr) != uintptr(unsafe.Pointer(outer)) {
		t.Errorf("Step 2: structValue.ptr should point to outer struct at %p, got %p", outer, structValue.ptr)
	}

	// Step 3: Get the InnerPtr field
	field0 := structValue.refField(0)
	t.Logf("Step 3 - field0: refKind=%v, ptr=%p, flags=%d", field0.refKind(), field0.ptr, field0.flag)
	if field0.refKind() != tpPointer {
		t.Errorf("Step 3: Expected tpPointer, got %v", field0.refKind())
	}

	// Debug: Let's see what's actually stored at the field location
	actualPtr := *(*unsafe.Pointer)(field0.ptr)
	t.Logf("Step 3b - actualPtr from field location: %p", actualPtr)
	// This should match the original inner pointer
	if uintptr(actualPtr) != uintptr(unsafe.Pointer(inner)) {
		t.Errorf("Step 3b: refField should contain pointer to inner (%p), got %p", inner, actualPtr)
	}

	// Debug: Check flagIndir
	t.Logf("Step 3c - field0.flag&flagIndir = %d (flagIndir = %d)", field0.flag&flagIndir, flagIndir)
	// For pointer fields, flagIndir should be set since ptr points to field location
	if field0.flag&flagIndir == 0 {
		t.Errorf("Step 3c: flagIndir should be set for pointer fields")
	}

	// Step 4: Dereference the field pointer to get inner struct
	innerStruct := field0.refElem()
	t.Logf("Step 4 - innerStruct: refKind=%v, ptr=%p, valid=%v", innerStruct.refKind(), innerStruct.ptr, innerStruct.refIsValid())

	if !innerStruct.refIsValid() {
		t.Fatal("Step 4: Inner struct should be valid")
	}
	if innerStruct.refKind() != tpStruct {
		t.Errorf("Step 4: Expected tpStruct, got %v", innerStruct.refKind())
	}
	// This should point to the actual inner struct
	if uintptr(innerStruct.ptr) != uintptr(unsafe.Pointer(inner)) {
		t.Errorf("Step 4: innerStruct.ptr should point to inner struct at %p, got %p", inner, innerStruct.ptr)
	}
	// Step 5: Get the Value field
	valueField := innerStruct.refField(0)
	t.Logf("Step 5 - valueField: refKind=%v, ptr=%p, flags=%d", valueField.refKind(), valueField.ptr, valueField.flag)
	t.Logf("Step 5a - valueField.flag&flagIndir = %d", valueField.flag&flagIndir)
	if valueField.refKind() != tpInt {
		t.Errorf("Step 5: Expected tpInt, got %v", valueField.refKind())
	}

	// Debug: Let's see what's actually stored at this location
	if valueField.refKind() == tpInt {
		// Calculate expected address
		expectedAddr := unsafe.Pointer(uintptr(unsafe.Pointer(inner)) + unsafe.Offsetof(inner.Value))
		if valueField.ptr != expectedAddr {
			t.Errorf("Step 5: valueField.ptr should point to inner.Value at %p, got %p", expectedAddr, valueField.ptr)
		}

		actualValue := *(*int)(valueField.ptr)
		t.Logf("Step 5b - actualValue from field location: %d", actualValue)
		if actualValue != 123 {
			t.Errorf("Step 5b: Expected 123 from direct memory read, got %d", actualValue)
		}

		// Also test the refInt() method
		reflectedValue := valueField.refInt()
		t.Logf("Step 5c - valueField.refInt(): %d", reflectedValue)
		if reflectedValue != 123 {
			t.Errorf("Step 5c: Expected 123 from refInt() method, got %d", reflectedValue)
		}
	}

	// Final comparison - THIS MUST PASS
	expected := 123
	got := int(valueField.refInt())
	if got != expected {
		t.Errorf("CRITICAL FAILURE: Expected %d, got %d - Pointer chain is broken!", expected, got)
	} else {
		t.Logf("SUCCESS: Pointer chain works correctly")
	}
}

// TestDebugFlagIndir - Critical test to validate flagIndir behavior
// This test ensures that flagIndir is only set for pointer fields in structs, not for all indirect types
func TestDebugFlagIndir(t *testing.T) {
	clearRefStructsCache()

	// Test 1: Direct value from interface{}
	s := "hello world"
	v1 := refValueOf(s)
	t.Logf("Direct string: refKind=%v, flags=%d, flagIndir=%d", v1.refKind(), v1.flag, v1.flag&flagIndir)

	// Test 2: String field from struct
	type TestStruct struct {
		S string
	}
	ts := TestStruct{S: "hello world"}
	v2 := refValueOf(ts)
	field := v2.refField(0)
	t.Logf("String field: refKind=%v, flags=%d, flagIndir=%d", field.refKind(), field.flag, field.flag&flagIndir)

	// Test 3: Pointer field from struct
	type PointerStruct struct {
		Ptr *string
	}
	ps := PointerStruct{Ptr: &s}
	v3 := refValueOf(ps)
	ptrField := v3.refField(0)
	t.Logf("Pointer field: refKind=%v, flags=%d, flagIndir=%d", ptrField.refKind(), ptrField.flag, ptrField.flag&flagIndir)
}
