package tinystring

import (
	"testing"
)

// Test structures for reflection functionality
type TestStruct struct {
	A int
	B string
	C float64
	D bool
}

type BigStruct struct {
	A, B, C, D, E int64
}

type NestedStruct struct {
	Basic TestStruct
	Value int
}

// Test basic refValueOf and refKind detection
func TestRefValueOfBasicTypes(t *testing.T) {
	tests := []struct {
		value    any
		expected kind
		name     string
	}{
		{int(42), tpInt, "int"},
		{int8(42), tpInt8, "int8"},
		{int16(42), tpInt16, "int16"},
		{int32(42), tpInt32, "int32"},
		{int64(42), tpInt64, "int64"},
		{uint(42), tpUint, "uint"},
		{uint8(42), tpUint8, "uint8"},
		{uint16(42), tpUint16, "uint16"},
		{uint32(42), tpUint32, "uint32"},
		{uint64(42), tpUint64, "uint64"},
		{float32(3.14), tpFloat32, "float32"},
		{float64(3.14), tpFloat64, "float64"},
		{true, tpBool, "bool"},
		{"hello", tpString, "string"},
		{[]int{1, 2, 3}, tpSlice, "slice"},
		{TestStruct{}, tpStruct, "struct"},
	}

	for _, test := range tests {
		v := refValueOf(test.value)
		if v.refKind() != test.expected {
			t.Errorf("%s: got kind %v, want %v", test.name, v.refKind(), test.expected)
		}
		if !v.refIsValid() {
			t.Errorf("%s: value should be valid", test.name)
		}
	}
}

// Test pointer handling with refElem()
func TestRefValuePointerElem(t *testing.T) {
	// Test with different pointer types
	tests := []struct {
		value    any
		expected kind
		name     string
	}{
		{new(int), tpInt, "*int"},
		{new(string), tpString, "*string"},
		{new(bool), tpBool, "*bool"},
		{&TestStruct{}, tpStruct, "*struct"},
	}

	for _, test := range tests {
		v := refValueOf(test.value)

		// Pointer should be detected
		if v.refKind() != tpPointer {
			t.Errorf("%s: got kind %v, want %v (pointer)", test.name, v.refKind(), tpPointer)
			continue
		}

		// refElem() should give us the pointed-to type
		elem := v.refElem()
		if !elem.refIsValid() {
			t.Errorf("%s: refElem() should be valid", test.name)
			continue
		}

		if elem.refKind() != test.expected {
			t.Errorf("%s: refElem() got kind %v, want %v", test.name, elem.refKind(), test.expected)
		}
	}
}

// Test setting values through reflection
func TestRefValueSetters(t *testing.T) {
	// Test string setting
	t.Run("SetString", func(t *testing.T) {
		var s string
		v := refValueOf(&s).refElem()
		if v.refKind() != tpString {
			t.Fatalf("Expected string kind, got %v", v.refKind())
		}

		v.refSetString("hello world")
		if s != "hello world" {
			t.Errorf("SetString failed: got %q, want %q", s, "hello world")
		}
	})

	// Test integer setting
	t.Run("SetInt", func(t *testing.T) {
		var i int64
		v := refValueOf(&i).refElem()
		if v.refKind() != tpInt64 {
			t.Fatalf("Expected int64 kind, got %v", v.refKind())
		}
		v.refSetInt(42)
		if i != 42 {
			t.Errorf("refSetInt failed: got %d, want %d", i, 42)
		}
	})

	// Test boolean setting
	t.Run("SetBool", func(t *testing.T) {
		var b bool
		v := refValueOf(&b).refElem()
		if v.refKind() != tpBool {
			t.Fatalf("Expected bool kind, got %v", v.refKind())
		}
		v.refSetBool(true)
		if !b {
			t.Errorf("refSetBool failed: got %v, want %v", b, true)
		}
	})

	// Test float setting
	t.Run("SetFloat", func(t *testing.T) {
		var f float64
		v := refValueOf(&f).refElem()
		if v.refKind() != tpFloat64 {
			t.Fatalf("Expected float64 kind, got %v", v.refKind())
		}
		v.refSetFloat(3.14)
		if f != 3.14 {
			t.Errorf("refSetFloat failed: got %f, want %f", f, 3.14)
		}
	})
}

// Test getting values through reflection
func TestRefValueGetters(t *testing.T) {
	// Test string getting
	t.Run("String", func(t *testing.T) {
		s := "hello world"
		v := refValueOf(s)
		if got := v.String(); got != s {
			t.Errorf("String() failed: got %q, want %q", got, s)
		}
	})

	// Test integer getting
	t.Run("Int", func(t *testing.T) {
		i := int64(42)
		v := refValueOf(i)
		if got := v.refInt(); got != i {
			t.Errorf("refInt() failed: got %d, want %d", got, i)
		}
	})

	// Test boolean getting
	t.Run("Bool", func(t *testing.T) {
		b := true
		v := refValueOf(b)
		if got := v.refBool(); got != b {
			t.Errorf("refBool() failed: got %v, want %v", got, b)
		}
	})

	// Test float getting
	t.Run("Float", func(t *testing.T) {
		f := 3.14
		v := refValueOf(f)
		if got := v.refFloat(); got != f {
			t.Errorf("refFloat() failed: got %f, want %f", got, f)
		}
	})
}

// Test struct field access
func TestRefValueStructFields(t *testing.T) {
	s := TestStruct{
		A: 42,
		B: "hello",
		C: 3.14,
		D: true,
	}

	v := refValueOf(s)
	if v.refKind() != tpStruct {
		t.Fatalf("Expected struct kind, got %v", v.refKind())
	}

	// Test refNumField
	numFields := v.refNumField()
	if numFields != 4 {
		t.Errorf("refNumField() got %d, want %d", numFields, 4)
	}

	// Test individual field access
	t.Run("Field0_A", func(t *testing.T) {
		field := v.refField(0)
		if field.refKind() != tpInt {
			t.Errorf("refField(0) kind got %v, want %v", field.refKind(), tpInt)
		}
		if got := int(field.refInt()); got != s.A {
			t.Errorf("refField(0) value got %d, want %d", got, s.A)
		}
	})

	t.Run("Field1_B", func(t *testing.T) {
		field := v.refField(1)
		if field.refKind() != tpString {
			t.Errorf("refField(1) kind got %v, want %v", field.refKind(), tpString)
		}
		if got := field.String(); got != s.B {
			t.Errorf("refField(1) value got %q, want %q", got, s.B)
		}
	})
}

// Test struct field setting through pointers
func TestRefValueStructFieldSetting(t *testing.T) {
	s := &TestStruct{}
	v := refValueOf(s).refElem()

	if v.refKind() != tpStruct {
		t.Fatalf("Expected struct kind, got %v", v.refKind())
	}

	// Set field A (int)
	fieldA := v.refField(0)
	if fieldA.refKind() != tpInt {
		t.Fatalf("refField 0 expected int kind, got %v", fieldA.refKind())
	}
	fieldA.refSetInt(100)
	if s.A != 100 {
		t.Errorf("refField A setting failed: got %d, want %d", s.A, 100)
	}

	// Set field B (string)
	fieldB := v.refField(1)
	if fieldB.refKind() != tpString {
		t.Fatalf("refField 1 expected string kind, got %v", fieldB.refKind())
	}
	fieldB.refSetString("test")
	if s.B != "test" {
		t.Errorf("refField B setting failed: got %q, want %q", s.B, "test")
	}
}

// Test nil pointer handling
func TestRefValueNilPointer(t *testing.T) {
	var p *int
	v := refValueOf(p)

	if v.refKind() != tpPointer {
		t.Fatalf("Expected pointer kind, got %v", v.refKind())
	}

	elem := v.refElem()
	if elem.refIsValid() {
		t.Errorf("refElem() of nil pointer should not be valid")
	}
}

// Test Interface() method
func TestRefValueInterface(t *testing.T) {
	tests := []struct {
		value any
		name  string
	}{
		{int(42), "int"},
		{int64(42), "int64"},
		{float64(3.14), "float64"},
		{true, "bool"},
		{"hello", "string"},
	}

	for _, test := range tests {
		v := refValueOf(test.value)
		got := v.Interface()

		// For basic types, Interface() should return the same value
		switch test.value.(type) {
		case int:
			if got.(int) != test.value.(int) {
				t.Errorf("%s Interface() got %v, want %v", test.name, got, test.value)
			}
		case int64:
			if got.(int64) != test.value.(int64) {
				t.Errorf("%s Interface() got %v, want %v", test.name, got, test.value)
			}
		case float64:
			if got.(float64) != test.value.(float64) {
				t.Errorf("%s Interface() got %v, want %v", test.name, got, test.value)
			}
		case bool:
			if got.(bool) != test.value.(bool) {
				t.Errorf("%s Interface() got %v, want %v", test.name, got, test.value)
			}
		case string:
			if got.(string) != test.value.(string) {
				t.Errorf("%s Interface() got %v, want %v", test.name, got, test.value)
			}
		}
	}
}

// Test big struct to ensure our reflection works with larger structures
func TestRefValueBigStruct(t *testing.T) {
	s := BigStruct{1, 2, 3, 4, 5}
	v := refValueOf(s)

	if v.refKind() != tpStruct {
		t.Fatalf("Expected struct kind, got %v", v.refKind())
	}

	if numFields := v.refNumField(); numFields != 5 {
		t.Errorf("BigStruct refNumField() got %d, want %d", numFields, 5)
	}

	// Test getting values from all fields
	for i := 0; i < 5; i++ {
		field := v.refField(i)
		if field.refKind() != tpInt64 {
			t.Errorf("refField %d expected int64 kind, got %v", i, field.refKind())
		}
		expectedValue := int64(i + 1)
		if got := field.refInt(); got != expectedValue {
			t.Errorf("refField %d got %d, want %d", i, got, expectedValue)
		}
	}
}

// Test that helps debug our current JSON decode issue
func TestRefValueDebugPointerChain(t *testing.T) {
	// This mimics what happens in JSON decode
	var s string
	ptr := &s

	// Test the full chain like in parseJsonIntoTarget
	v1 := refValueOf(ptr)
	t.Logf("Step 1 - refValueOf(&s): refKind=%v, refIsValid=%v", v1.refKind(), v1.refIsValid())

	if v1.refKind() != tpPointer {
		t.Errorf("Expected pointer kind, got %v", v1.refKind())
		return
	}

	v2 := v1.refElem()
	t.Logf("Step 2 - v1.refElem(): refKind=%v, refIsValid=%v", v2.refKind(), v2.refIsValid())

	if !v2.refIsValid() {
		t.Errorf("refElem() should be valid for non-nil pointer")
		return
	}

	if v2.refKind() != tpString {
		t.Errorf("refElem() expected string kind, got %v", v2.refKind())
		return
	}

	// Test setting through the chain
	v2.refSetString("test value")
	if s != "test value" {
		t.Errorf("Setting through reflection failed: got %q, want %q", s, "test value")
	}

	t.Logf("SUCCESS: Pointer chain works correctly")
}

// Additional critical tests adapted from Go's internal/reflectlite/all_test.go

// TestNilPtrValueSub tests refElem() behavior with nil pointers
func TestNilPtrValueSub(t *testing.T) {
	var pi *int
	pv := refValueOf(pi)
	if pv.refElem().refIsValid() {
		t.Error("refValueOf((*int)(nil)).refElem().refIsValid() should be false")
	}
}

// TestPtrSetNil tests setting pointer values to nil
func TestPtrSetNil(t *testing.T) {
	var i int32 = 1234
	ip := &i
	vip := refValueOf(&ip)

	// vip is **int32, vip.refElem() is *int32
	if vip.refKind() != tpPointer {
		t.Fatalf("Expected pointer to pointer, got %v", vip.refKind())
	}

	elemValue := vip.refElem()
	if elemValue.refKind() != tpPointer {
		t.Fatalf("Expected pointer after refElem(), got %v", elemValue.refKind())
	}

	// Set the *int32 to nil (zero value for pointer)
	zeroPtr := refZero(elemValue.Type())
	elemValue.refSet(zeroPtr)

	if ip != nil {
		t.Errorf("got non-nil (%d), want nil", *ip)
	}
}

// TestPointerElem tests various pointer dereferencing scenarios
func TestPointerElem(t *testing.T) {
	// Test basic pointer elem
	i := 42
	pi := &i
	v := refValueOf(pi)

	if v.refKind() != tpPointer {
		t.Fatalf("Expected pointer, got %v", v.refKind())
	}

	elem := v.refElem()
	if !elem.refIsValid() {
		t.Fatal("refElem() should be valid for non-nil pointer")
	}

	if elem.refKind() != tpInt {
		t.Errorf("Expected int after refElem(), got %v", elem.refKind())
	}

	if elem.refInt() != 42 {
		t.Errorf("Expected 42, got %d", elem.refInt())
	}
}

// TestStructPointerChain tests dereferencing through struct fields
func TestStructPointerChain(t *testing.T) {
	type Inner struct {
		Value int
	}

	type Outer struct {
		InnerPtr *Inner
	}

	inner := &Inner{Value: 123}
	outer := &Outer{InnerPtr: inner}

	// Test accessing nested pointer through reflection
	v := refValueOf(outer)
	if v.refKind() != tpPointer {
		t.Fatalf("Expected pointer to struct, got %v", v.refKind())
	}

	structValue := v.refElem()
	if structValue.refKind() != tpStruct {
		t.Fatalf("Expected struct after refElem(), got %v", structValue.refKind())
	}

	field0 := structValue.refField(0)
	if field0.refKind() != tpPointer {
		t.Fatalf("Expected pointer field, got %v", field0.refKind())
	}

	innerStruct := field0.refElem()
	if !innerStruct.refIsValid() {
		t.Fatal("Inner struct should be valid")
	}

	if innerStruct.refKind() != tpStruct {
		t.Fatalf("Expected struct after field refElem(), got %v", innerStruct.refKind())
	}

	valueField := innerStruct.refField(0)
	if valueField.refKind() != tpInt {
		t.Fatalf("Expected int field, got %v", valueField.refKind())
	}

	if valueField.refInt() != 123 {
		t.Errorf("Expected 123, got %d", valueField.refInt())
	}
}

// TestInterfaceValue tests interface{} dereferencing
func TestInterfaceValue(t *testing.T) {
	var inter struct {
		E any
	}
	inter.E = 123.456
	v1 := refValueOf(&inter)
	v2 := v1.refElem().refField(0)

	// v2 should be interface{} containing float64
	v3 := v2.refElem()
	if v3.refKind() != tpFloat64 {
		t.Errorf("Expected float64 in interface, got %v", v3.refKind())
	}

	i3 := v2.Interface()
	if _, ok := i3.(float64); !ok {
		t.Errorf("Interface() did not return float64, got %T", i3)
	}
}

// TestSetValue tests setting values through reflection (adapted from Go's tests)
func TestSetValue(t *testing.T) {
	tests := []struct {
		ptr      any
		newValue any
		expected any
		name     string
	}{
		{new(int), int(132), int(132), "int"},
		{new(int8), int8(8), int8(8), "int8"},
		{new(int16), int16(16), int16(16), "int16"},
		{new(int32), int32(32), int32(32), "int32"},
		{new(int64), int64(64), int64(64), "int64"},
		{new(uint), uint(132), uint(132), "uint"},
		{new(uint8), uint8(8), uint8(8), "uint8"},
		{new(uint16), uint16(16), uint16(16), "uint16"},
		{new(uint32), uint32(32), uint32(32), "uint32"},
		{new(uint64), uint64(64), uint64(64), "uint64"},
		{new(float32), float32(256.25), float32(256.25), "float32"},
		{new(float64), 512.125, 512.125, "float64"},
		{new(string), "stringy cheese", "stringy cheese", "string"},
		{new(bool), true, true, "bool"},
	}

	for _, tt := range tests {
		v := refValueOf(tt.ptr).refElem()
		newVal := refValueOf(tt.newValue)
		v.refSet(newVal)

		// Verify the value was set correctly
		result := v.Interface()
		if result != tt.expected {
			t.Errorf("%s: Set() failed, got %v, want %v", tt.name, result, tt.expected)
		}
	}
}

// TestPointerChainConsistency - Test to ensure pointer chains work consistently
func TestPointerChainConsistency(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected any
	}{
		{"int_pointer", &[]int{42}[0], int(42)},
		{"string_pointer", &[]string{"hello"}[0], "hello"},
		{"bool_pointer", &[]bool{true}[0], true},
		{"float64_pointer", &[]float64{3.14}[0], 3.14},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := refValueOf(test.value)

			// Should be pointer
			if v.refKind() != tpPointer {
				t.Errorf("Expected pointer type, got %v", v.refKind())
				return
			}

			// refElem() should work
			elem := v.refElem()
			if !elem.refIsValid() {
				t.Error("refElem() should be valid for non-nil pointer")
				return
			}

			// Interface() should return original value
			result := elem.Interface()
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

// TestPointerTypeElemAccess - Test various ways to access pointer element types
func TestPointerTypeElemAccess(t *testing.T) {
	// Test different pointer types
	tests := []struct {
		name     string
		ptr      any
		elemKind kind
	}{
		{"*int", new(int), tpInt},
		{"*string", new(string), tpString},
		{"*bool", new(bool), tpBool},
		{"*float64", new(float64), tpFloat64},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := refValueOf(test.ptr)

			// Validate pointer type
			if v.refKind() != tpPointer {
				t.Errorf("Expected pointer, got %v", v.refKind())
				return
			}

			// Test type.refElem()
			elemType := v.Type().refElem()
			if elemType == nil {
				t.Error("Type().refElem() should not be nil for pointer")
				return
			}

			if elemType.refKind() != test.elemKind {
				t.Errorf("Type().refElem().refKind() expected %v, got %v", test.elemKind, elemType.refKind())
			}

			// Test value.refElem()
			elem := v.refElem()
			if !elem.refIsValid() {
				t.Error("refElem() should be valid")
				return
			}

			if elem.refKind() != test.elemKind {
				t.Errorf("refElem().refKind() expected %v, got %v", test.elemKind, elem.refKind())
			}
		})
	}
}

// TestStructStringFieldAccess - Critical test for struct string field access
// This test ensures that string fields in structs are read correctly through reflection
func TestStructStringFieldAccess(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	s := TestStruct{
		Name: "John Doe",
		Age:  30,
	}

	// Test through reflection
	v := refValueOf(s)
	if v.refKind() != tpStruct {
		t.Fatalf("Expected struct, got %v", v.refKind())
	}

	// Get Name field (field 0)
	nameField := v.refField(0)
	if nameField.refKind() != tpString {
		t.Errorf("Expected string field, got %v", nameField.refKind())
	}

	// Validate flagIndir is NOT set for string fields in structs
	if nameField.flag&flagIndir != 0 {
		t.Errorf("flagIndir should not be set for string fields in structs")
	}

	// Test reading the string through reflection
	nameValue := nameField.String()
	if nameValue != s.Name {
		t.Errorf("String field mismatch: got %q, want %q", nameValue, s.Name)
	}
	// Validate memory address consistency by testing that direct access works
	if nameField.ptr != nil {
		directValue := *(*string)(nameField.ptr)
		if directValue != s.Name {
			t.Errorf("Direct memory access failed: got %q, want %q", directValue, s.Name)
		}
	} else {
		t.Error("refField pointer should not be nil")
	}
}

// TestComplexStructStringFields - Test with JSON-like struct to prevent string corruption
// This test replicates the exact scenario that was causing string corruption in JSON encoding
func TestComplexStructStringFields(t *testing.T) {
	type ComplexUser struct {
		ReadStatus string `json:"read_status"`
		OpenStat   string `json:"open_stat"`
	}

	user := ComplexUser{
		ReadStatus: "read",
		OpenStat:   "open",
	}

	v := refValueOf(user)
	if v.refKind() != tpStruct {
		t.Fatalf("Expected struct, got %v", v.refKind())
	}

	// Test ReadStatus field
	readStatusField := v.refField(0)
	if readStatusField.refKind() != tpString {
		t.Errorf("ReadStatus field should be string, got %v", readStatusField.refKind())
	}

	readStatus := readStatusField.String()
	if readStatus != user.ReadStatus {
		t.Errorf("ReadStatus mismatch: got %q, want %q", readStatus, user.ReadStatus)
	}

	// Test OpenStat field
	openStatField := v.refField(1)
	if openStatField.refKind() != tpString {
		t.Errorf("OpenStat field should be string, got %v", openStatField.refKind())
	}

	openStat := openStatField.String()
	if openStat != user.OpenStat {
		t.Errorf("OpenStat mismatch: got %q, want %q", openStat, user.OpenStat)
	}

	// Test string length validity (corruption check)
	if len(readStatus) > 1000 || len(openStat) > 1000 {
		t.Errorf("String corruption detected: readStatus len=%d, openStat len=%d", len(readStatus), len(openStat))
	}
}

// TestStringFieldFlagConsistency - Ensure flagIndir is handled correctly for different field types
func TestStringFieldFlagConsistency(t *testing.T) {
	type MixedStruct struct {
		StringField  string
		IntField     int
		PointerField *string
	}

	str := "test"
	s := MixedStruct{
		StringField:  "hello",
		IntField:     42,
		PointerField: &str,
	}

	v := refValueOf(s)

	// Test string field (should NOT have flagIndir)
	stringField := v.refField(0)
	if stringField.flag&flagIndir != 0 {
		t.Errorf("String field should not have flagIndir set")
	}

	// Test int field (should NOT have flagIndir)
	intField := v.refField(1)
	if intField.flag&flagIndir != 0 {
		t.Errorf("Int field should not have flagIndir set")
	}

	// Test pointer field (SHOULD have flagIndir)
	pointerField := v.refField(2)
	if pointerField.flag&flagIndir == 0 {
		t.Errorf("Pointer field should have flagIndir set")
	}

	// Validate values
	if stringField.String() != s.StringField {
		t.Errorf("String field value mismatch: got %q, want %q", stringField.String(), s.StringField)
	}

	if intField.refInt() != int64(s.IntField) {
		t.Errorf("Int field value mismatch: got %d, want %d", intField.refInt(), s.IntField)
	}

	// Validate pointer field dereferencing
	pointerElem := pointerField.refElem()
	if !pointerElem.refIsValid() {
		t.Error("Pointer field refElem() should be valid")
	} else if pointerElem.String() != *s.PointerField {
		t.Errorf("Pointer field dereferenced value mismatch: got %q, want %q", pointerElem.String(), *s.PointerField)
	}
}

func TestRefSetUint(t *testing.T) {
	tests := []struct {
		name        string
		createVar   func() any
		setValue    uint64
		expectError bool
	}{
		{"uint", func() any { var v uint; return &v }, 42, false},
		{"uint8", func() any { var v uint8; return &v }, 255, false},
		{"uint16", func() any { var v uint16; return &v }, 65535, false},
		{"uint32", func() any { var v uint32; return &v }, 4294967295, false},
		{"uint64", func() any { var v uint64; return &v }, 18446744073709551615, false},
		{"uintptr", func() any { var v uintptr; return &v }, 12345, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a variable and get its reflection
			variable := tt.createVar()
			rv := refValueOf(variable)

			// Get the element (dereferenced pointer)
			elem := rv.refElem()
			if !elem.refIsValid() {
				t.Fatalf("Element not valid for %s", tt.name)
			} // Set the uint value
			elem.refSetUint(tt.setValue)

			// Check if there was an error
			if elem.err != errNone {
				if !tt.expectError {
					t.Errorf("Unexpected error for %s: %v", tt.name, elem.err)
				}
			} else {
				if tt.expectError {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			}
		})
	}
}

func TestRefSetUintInvalidTypes(t *testing.T) {
	// Test setting uint on invalid types (should set error)
	tests := []struct {
		name      string
		createVar func() any
	}{
		{"string", func() any { var v string; return &v }},
		{"int", func() any { var v int; return &v }},
		{"float64", func() any { var v float64; return &v }},
		{"bool", func() any { var v bool; return &v }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variable := tt.createVar()
			rv := refValueOf(variable)
			elem := rv.refElem() // This should set an error
			elem.refSetUint(42)

			if elem.err == errNone {
				t.Errorf("Expected error when setting uint on %s, but got none", tt.name)
			}
		})
	}
}

func TestRefZeroComprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		testDesc string
	}{
		{"int", 42, "zero value of int"},
		{"string", "hello", "zero value of string"},
		{"bool", true, "zero value of bool"},
		{"float64", 3.14, "zero value of float64"},
		{"pointer to int", &[]int{1}[0], "zero value of pointer"},
		{"struct", TestStruct{A: 1, B: "test"}, "zero value of struct"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get the type of the input
			rv := refValueOf(tt.input)
			typ := rv.Type()

			// Create zero value
			zero := refZero(typ)

			if zero.err != errNone {
				t.Errorf("refZero() error: %v", zero.err)
				return
			}

			if !zero.refIsValid() {
				t.Errorf("refZero() returned invalid value")
				return
			}

			// The zero value should have the same type
			if zero.refKind() != rv.refKind() {
				t.Errorf("refZero() kind mismatch: expected %v, got %v",
					rv.refKind(), zero.refKind())
			}

			t.Logf("%s: kind=%v, valid=%v", tt.testDesc, zero.refKind(), zero.refIsValid())
		})
	}
}

func TestRefZeroNilType(t *testing.T) {
	// Test refZero with nil type (should return error)
	zero := refZero(nil)

	if zero.err == errNone {
		t.Error("refZero(nil) should return an error")
	}

	expectedErr := "reflect: Zero(nil)"
	if string(zero.err) != expectedErr {
		t.Errorf("refZero(nil) error: expected %q, got %q", expectedErr, string(zero.err))
	}
}

func TestRefZeroLargeType(t *testing.T) {
	// Test refZero with a reasonably sized type to ensure it works
	type LargeStruct struct {
		Data [100]int
	}

	large := LargeStruct{}
	rv := refValueOf(large)
	typ := rv.Type()

	zero := refZero(typ)

	if zero.err != errNone {
		t.Errorf("refZero() error with large struct: %v", zero.err)
		return
	}

	if !zero.refIsValid() {
		t.Error("refZero() returned invalid value for large struct")
	}
}
