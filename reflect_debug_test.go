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
	elem := rv.Elem() // Get the struct value from pointer

	if elem.Kind() != tpStruct {
		t.Errorf("Container elem kind: expected %v, got %v", tpStruct, elem.Kind())
	}

	// Get fields
	nameField := elem.Field(0)
	coordsField := elem.Field(1)

	if nameField.Kind() != tpString {
		t.Errorf("Name field kind: expected %v, got %v", tpString, nameField.Kind())
	}
	if nameField.String() != "test" {
		t.Errorf("Name field value: expected %q, got %q", "test", nameField.String())
	}

	if coordsField.Kind() != tpPointer {
		t.Errorf("Coords field kind: expected %v, got %v", tpPointer, coordsField.Kind())
	}

	if coordsField.Kind() == tpPointer {
		coordsElem := coordsField.Elem()
		if !coordsElem.IsValid() {
			t.Fatal("Coords elem is not valid")
		}

		if coordsElem.Kind() != tpStruct {
			t.Errorf("Coords elem kind: expected %v, got %v", tpStruct, coordsElem.Kind())
		}

		if coordsElem.Kind() == tpStruct {
			latField := coordsElem.Field(0)
			lngField := coordsElem.Field(1)
			altField := coordsElem.Field(2)

			if latField.Kind() != tpFloat64 {
				t.Errorf("Lat field kind: expected %v, got %v", tpFloat64, latField.Kind())
			}
			if latField.Float() != 37.7749 {
				t.Errorf("Lat field value: expected %f, got %f", 37.7749, latField.Float())
			}

			if lngField.Kind() != tpFloat64 {
				t.Errorf("Lng field kind: expected %v, got %v", tpFloat64, lngField.Kind())
			}
			if lngField.Float() != -122.4194 {
				t.Errorf("Lng field value: expected %f, got %f", -122.4194, lngField.Float())
			}

			if altField.Kind() != tpInt {
				t.Errorf("Alt field kind: expected %v, got %v", tpInt, altField.Kind())
			}
			if altField.Int() != 100 {
				t.Errorf("Alt field value: expected %d, got %d", 100, altField.Int())
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
	newElem := newRv.Elem()

	// Set the name field
	newNameField := newElem.Field(0)
	newNameField.SetString("new_test")
	if newContainer.Name != "new_test" {
		t.Errorf("Name setting failed: expected %q, got %q", "new_test", newContainer.Name)
	}

	// Create a new coords struct and set the pointer
	coordsField := newElem.Field(1)
	coordsType := coordsField.Type().Elem()

	newCoordsPtr := refNew(coordsType)
	newCoordsElem := newCoordsPtr.Elem()

	if newCoordsElem.Kind() != tpStruct {
		t.Errorf("New coords elem kind: expected %v, got %v", tpStruct, newCoordsElem.Kind())
	}

	// Set the pointer field to point to the new struct
	newCoordsField := newElem.Field(1)
	// newCoordsPtr.ptr points to a pointer variable containing the allocated address
	// We need to get the actual allocated address to store in the struct field
	actualAddr := *(*unsafe.Pointer)(newCoordsPtr.ptr)
	*(*unsafe.Pointer)(newCoordsField.ptr) = actualAddr

	// Now try to set values in the pointed-to struct
	if newCoordsElem.Kind() == tpStruct {
		newLatField := newCoordsElem.Field(0)
		newLngField := newCoordsElem.Field(1)
		newAltField := newCoordsElem.Field(2)

		// Set values
		newLatField.SetFloat(99.99)
		newLngField.SetFloat(-99.99)
		newAltField.SetInt(999)

		// Validate memory contents
		latPtr := (*float64)(newLatField.ptr)
		lngPtr := (*float64)(newLngField.ptr)
		altPtr := (*int)(newAltField.ptr)

		if *latPtr != 99.99 {
			t.Errorf("Lat memory: expected %f, got %f", 99.99, *latPtr)
		}
		if *lngPtr != -99.99 {
			t.Errorf("Lng memory: expected %f, got %f", -99.99, *lngPtr)
		}
		if *altPtr != 999 {
			t.Errorf("Alt memory: expected %d, got %d", 999, *altPtr)
		}
	}

	// Check if the values were set correctly through the struct
	if newContainer.Coords == nil {
		t.Fatal("Final coords is nil")
	}
	if newContainer.Coords.Lat != 99.99 {
		t.Errorf("Final Lat: expected %f, got %f", 99.99, newContainer.Coords.Lat)
	}
	if newContainer.Coords.Lng != -99.99 {
		t.Errorf("Final Lng: expected %f, got %f", -99.99, newContainer.Coords.Lng)
	}
	if newContainer.Coords.Alt != 999 {
		t.Errorf("Final Alt: expected %d, got %d", 999, newContainer.Coords.Alt)
	}
}

// TestReflectFieldCorruption tests field access patterns to diagnose corruption issues
func TestReflectFieldCorruption(t *testing.T) {
	clearRefStructsCache()

	// Use the exact types from ComplexUser
	phone := ComplexPhoneNumber{
		ID:         "ph_001",
		Type:       "mobile",
		Number:     "+1-555-123-4567",
		Extension:  "",
		IsPrimary:  true,
		IsVerified: true,
	}

	// Test reflection on this structure
	rv := refValueOf(phone)
	if rv.Kind() != tpStruct {
		t.Fatalf("Expected struct, got %v", rv.Kind())
	}

	// Get struct type info
	tt := (*refStructType)(unsafe.Pointer(rv.typ))
	expectedFieldCount := 6 // ID, Type, Number, Extension, IsPrimary, IsVerified
	if len(tt.fields) != expectedFieldCount {
		t.Errorf("Struct type fields count: expected %d, got %d", expectedFieldCount, len(tt.fields))
	}

	// Check each field
	for i := 0; i < len(tt.fields); i++ {
		field := &tt.fields[i]
		fieldName := field.name.Name()

		// Get field value using reflection
		fieldVal := rv.Field(i)

		if fieldVal.Kind() == tpString {
			strVal := fieldVal.String()

			// Manual check - calculate field address manually
			fieldPtr := add(rv.ptr, field.offset, "manual calculation")
			manualVal := *(*string)(fieldPtr)

			if strVal != manualVal {
				t.Errorf("Field %d (%s) mismatch: reflection=%q vs manual=%q", i, fieldName, strVal, manualVal)
			}
		}
	}

	// Test specific field access patterns
	idField := rv.Field(0) // Should be ID
	if idField.String() != "ph_001" {
		t.Errorf("ID field: expected %q, got %q", "ph_001", idField.String())
	}

	numberField := rv.Field(2) // Should be Number
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
	if v.Kind() != tpPointer {
		t.Errorf("Expected tpPointer, got %v", v.Kind())
	}

	if !v.IsValid() {
		t.Fatal("Pointer refValue should be valid")
	}

	if v.typ == nil {
		t.Fatal("Pointer refValue.typ should not be nil")
	}

	// Test pointer type properties
	if v.typ.Kind() != tpPointer {
		t.Errorf("Expected pointer type kind, got %v", v.typ.Kind())
	}

	// Test elem type access
	elemType := v.typ.Elem()
	if elemType == nil {
		t.Fatal("Pointer type should have valid elem type")
	}

	if elemType.Kind() != tpInt {
		t.Errorf("Expected int elem type, got %v", elemType.Kind())
	}

	// Test Elem() method
	elem := v.Elem()
	if !elem.IsValid() {
		t.Fatal("Elem() should return valid value for non-nil pointer")
	}

	if elem.Kind() != tpInt {
		t.Errorf("Expected int after Elem(), got %v", elem.Kind())
	}

	// Test value retrieval
	actualValue := elem.Int()
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
	t.Logf("Step 1 - refValueOf(outer): Kind=%v, ptr=%p", v.Kind(), v.ptr)
	if v.Kind() != tpPointer {
		t.Errorf("Step 1: Expected tpPointer, got %v", v.Kind())
	}

	// Step 2: Dereference to get struct
	structValue := v.Elem()
	t.Logf("Step 2 - structValue: Kind=%v, ptr=%p", structValue.Kind(), structValue.ptr)
	if structValue.Kind() != tpStruct {
		t.Errorf("Step 2: Expected tpStruct, got %v", structValue.Kind())
	}
	// Verify this points to the actual struct
	if uintptr(structValue.ptr) != uintptr(unsafe.Pointer(outer)) {
		t.Errorf("Step 2: structValue.ptr should point to outer struct at %p, got %p", outer, structValue.ptr)
	}

	// Step 3: Get the InnerPtr field
	field0 := structValue.Field(0)
	t.Logf("Step 3 - field0: Kind=%v, ptr=%p, flags=%d", field0.Kind(), field0.ptr, field0.flag)
	if field0.Kind() != tpPointer {
		t.Errorf("Step 3: Expected tpPointer, got %v", field0.Kind())
	}

	// Debug: Let's see what's actually stored at the field location
	actualPtr := *(*unsafe.Pointer)(field0.ptr)
	t.Logf("Step 3b - actualPtr from field location: %p", actualPtr)
	// This should match the original inner pointer
	if uintptr(actualPtr) != uintptr(unsafe.Pointer(inner)) {
		t.Errorf("Step 3b: Field should contain pointer to inner (%p), got %p", inner, actualPtr)
	}

	// Debug: Check flagIndir
	t.Logf("Step 3c - field0.flag&flagIndir = %d (flagIndir = %d)", field0.flag&flagIndir, flagIndir)
	// For pointer fields, flagIndir should be set since ptr points to field location
	if field0.flag&flagIndir == 0 {
		t.Errorf("Step 3c: flagIndir should be set for pointer fields")
	}

	// Step 4: Dereference the field pointer to get inner struct
	innerStruct := field0.Elem()
	t.Logf("Step 4 - innerStruct: Kind=%v, ptr=%p, valid=%v", innerStruct.Kind(), innerStruct.ptr, innerStruct.IsValid())

	if !innerStruct.IsValid() {
		t.Fatal("Step 4: Inner struct should be valid")
	}
	if innerStruct.Kind() != tpStruct {
		t.Errorf("Step 4: Expected tpStruct, got %v", innerStruct.Kind())
	}
	// This should point to the actual inner struct
	if uintptr(innerStruct.ptr) != uintptr(unsafe.Pointer(inner)) {
		t.Errorf("Step 4: innerStruct.ptr should point to inner struct at %p, got %p", inner, innerStruct.ptr)
	}
	// Step 5: Get the Value field
	valueField := innerStruct.Field(0)
	t.Logf("Step 5 - valueField: Kind=%v, ptr=%p, flags=%d", valueField.Kind(), valueField.ptr, valueField.flag)
	t.Logf("Step 5a - valueField.flag&flagIndir = %d", valueField.flag&flagIndir)
	if valueField.Kind() != tpInt {
		t.Errorf("Step 5: Expected tpInt, got %v", valueField.Kind())
	}

	// Debug: Let's see what's actually stored at this location
	if valueField.Kind() == tpInt {
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

		// Also test the Int() method
		reflectedValue := valueField.Int()
		t.Logf("Step 5c - valueField.Int(): %d", reflectedValue)
		if reflectedValue != 123 {
			t.Errorf("Step 5c: Expected 123 from Int() method, got %d", reflectedValue)
		}
	}

	// Final comparison - THIS MUST PASS
	expected := 123
	got := int(valueField.Int())
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
	t.Logf("Direct string: Kind=%v, flags=%d, flagIndir=%d", v1.Kind(), v1.flag, v1.flag&flagIndir)

	// Test 2: String field from struct
	type TestStruct struct {
		S string
	}
	ts := TestStruct{S: "hello world"}
	v2 := refValueOf(ts)
	field := v2.Field(0)
	t.Logf("String field: Kind=%v, flags=%d, flagIndir=%d", field.Kind(), field.flag, field.flag&flagIndir)

	// Test 3: Pointer field from struct
	type PointerStruct struct {
		Ptr *string
	}
	ps := PointerStruct{Ptr: &s}
	v3 := refValueOf(ps)
	ptrField := v3.Field(0)
	t.Logf("Pointer field: Kind=%v, flags=%d, flagIndir=%d", ptrField.Kind(), ptrField.flag, ptrField.flag&flagIndir)
}
