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

// Test basic refValueOf and Kind detection
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
		if v.Kind() != test.expected {
			t.Errorf("%s: got kind %v, want %v", test.name, v.Kind(), test.expected)
		}
		if !v.IsValid() {
			t.Errorf("%s: value should be valid", test.name)
		}
	}
}

// Test pointer handling with Elem()
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
		if v.Kind() != tpPointer {
			t.Errorf("%s: got kind %v, want %v (pointer)", test.name, v.Kind(), tpPointer)
			continue
		}

		// Elem() should give us the pointed-to type
		elem := v.Elem()
		if !elem.IsValid() {
			t.Errorf("%s: Elem() should be valid", test.name)
			continue
		}

		if elem.Kind() != test.expected {
			t.Errorf("%s: Elem() got kind %v, want %v", test.name, elem.Kind(), test.expected)
		}
	}
}

// Test setting values through reflection
func TestRefValueSetters(t *testing.T) {
	// Test string setting
	t.Run("SetString", func(t *testing.T) {
		var s string
		v := refValueOf(&s).Elem()
		if v.Kind() != tpString {
			t.Fatalf("Expected string kind, got %v", v.Kind())
		}

		v.SetString("hello world")
		if s != "hello world" {
			t.Errorf("SetString failed: got %q, want %q", s, "hello world")
		}
	})

	// Test integer setting
	t.Run("SetInt", func(t *testing.T) {
		var i int64
		v := refValueOf(&i).Elem()
		if v.Kind() != tpInt64 {
			t.Fatalf("Expected int64 kind, got %v", v.Kind())
		}

		v.SetInt(42)
		if i != 42 {
			t.Errorf("SetInt failed: got %d, want %d", i, 42)
		}
	})

	// Test boolean setting
	t.Run("SetBool", func(t *testing.T) {
		var b bool
		v := refValueOf(&b).Elem()
		if v.Kind() != tpBool {
			t.Fatalf("Expected bool kind, got %v", v.Kind())
		}

		v.SetBool(true)
		if !b {
			t.Errorf("SetBool failed: got %v, want %v", b, true)
		}
	})

	// Test float setting
	t.Run("SetFloat", func(t *testing.T) {
		var f float64
		v := refValueOf(&f).Elem()
		if v.Kind() != tpFloat64 {
			t.Fatalf("Expected float64 kind, got %v", v.Kind())
		}

		v.SetFloat(3.14)
		if f != 3.14 {
			t.Errorf("SetFloat failed: got %f, want %f", f, 3.14)
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
		if got := v.Int(); got != i {
			t.Errorf("Int() failed: got %d, want %d", got, i)
		}
	})

	// Test boolean getting
	t.Run("Bool", func(t *testing.T) {
		b := true
		v := refValueOf(b)
		if got := v.Bool(); got != b {
			t.Errorf("Bool() failed: got %v, want %v", got, b)
		}
	})

	// Test float getting
	t.Run("Float", func(t *testing.T) {
		f := 3.14
		v := refValueOf(f)
		if got := v.Float(); got != f {
			t.Errorf("Float() failed: got %f, want %f", got, f)
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
	if v.Kind() != tpStruct {
		t.Fatalf("Expected struct kind, got %v", v.Kind())
	}

	// Test NumField
	numFields := v.NumField()
	if numFields != 4 {
		t.Errorf("NumField() got %d, want %d", numFields, 4)
	}

	// Test individual field access
	t.Run("Field0_A", func(t *testing.T) {
		field := v.Field(0)
		if field.Kind() != tpInt {
			t.Errorf("Field(0) kind got %v, want %v", field.Kind(), tpInt)
		}
		if got := int(field.Int()); got != s.A {
			t.Errorf("Field(0) value got %d, want %d", got, s.A)
		}
	})

	t.Run("Field1_B", func(t *testing.T) {
		field := v.Field(1)
		if field.Kind() != tpString {
			t.Errorf("Field(1) kind got %v, want %v", field.Kind(), tpString)
		}
		if got := field.String(); got != s.B {
			t.Errorf("Field(1) value got %q, want %q", got, s.B)
		}
	})
}

// Test struct field setting through pointers
func TestRefValueStructFieldSetting(t *testing.T) {
	s := &TestStruct{}
	v := refValueOf(s).Elem()

	if v.Kind() != tpStruct {
		t.Fatalf("Expected struct kind, got %v", v.Kind())
	}

	// Set field A (int)
	fieldA := v.Field(0)
	if fieldA.Kind() != tpInt {
		t.Fatalf("Field 0 expected int kind, got %v", fieldA.Kind())
	}
	fieldA.SetInt(100)
	if s.A != 100 {
		t.Errorf("Field A setting failed: got %d, want %d", s.A, 100)
	}

	// Set field B (string)
	fieldB := v.Field(1)
	if fieldB.Kind() != tpString {
		t.Fatalf("Field 1 expected string kind, got %v", fieldB.Kind())
	}
	fieldB.SetString("test")
	if s.B != "test" {
		t.Errorf("Field B setting failed: got %q, want %q", s.B, "test")
	}
}

// Test nil pointer handling
func TestRefValueNilPointer(t *testing.T) {
	var p *int
	v := refValueOf(p)

	if v.Kind() != tpPointer {
		t.Fatalf("Expected pointer kind, got %v", v.Kind())
	}

	elem := v.Elem()
	if elem.IsValid() {
		t.Errorf("Elem() of nil pointer should not be valid")
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

	if v.Kind() != tpStruct {
		t.Fatalf("Expected struct kind, got %v", v.Kind())
	}

	if numFields := v.NumField(); numFields != 5 {
		t.Errorf("BigStruct NumField() got %d, want %d", numFields, 5)
	}

	// Test getting values from all fields
	for i := 0; i < 5; i++ {
		field := v.Field(i)
		if field.Kind() != tpInt64 {
			t.Errorf("Field %d expected int64 kind, got %v", i, field.Kind())
		}
		expectedValue := int64(i + 1)
		if got := field.Int(); got != expectedValue {
			t.Errorf("Field %d got %d, want %d", i, got, expectedValue)
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
	t.Logf("Step 1 - refValueOf(&s): Kind=%v, IsValid=%v", v1.Kind(), v1.IsValid())

	if v1.Kind() != tpPointer {
		t.Errorf("Expected pointer kind, got %v", v1.Kind())
		return
	}

	v2 := v1.Elem()
	t.Logf("Step 2 - v1.Elem(): Kind=%v, IsValid=%v", v2.Kind(), v2.IsValid())

	if !v2.IsValid() {
		t.Errorf("Elem() should be valid for non-nil pointer")
		return
	}

	if v2.Kind() != tpString {
		t.Errorf("Elem() expected string kind, got %v", v2.Kind())
		return
	}

	// Test setting through the chain
	v2.SetString("test value")
	if s != "test value" {
		t.Errorf("Setting through reflection failed: got %q, want %q", s, "test value")
	}

	t.Logf("SUCCESS: Pointer chain works correctly")
}

// Additional critical tests adapted from Go's internal/reflectlite/all_test.go

// TestNilPtrValueSub tests Elem() behavior with nil pointers
func TestNilPtrValueSub(t *testing.T) {
	var pi *int
	pv := refValueOf(pi)
	if pv.Elem().IsValid() {
		t.Error("refValueOf((*int)(nil)).Elem().IsValid() should be false")
	}
}

// TestPtrSetNil tests setting pointer values to nil
func TestPtrSetNil(t *testing.T) {
	var i int32 = 1234
	ip := &i
	vip := refValueOf(&ip)

	// vip is **int32, vip.Elem() is *int32
	if vip.Kind() != tpPointer {
		t.Fatalf("Expected pointer to pointer, got %v", vip.Kind())
	}

	elemValue := vip.Elem()
	if elemValue.Kind() != tpPointer {
		t.Fatalf("Expected pointer after Elem(), got %v", elemValue.Kind())
	}

	// Set the *int32 to nil (zero value for pointer)
	zeroPtr := refZero(elemValue.Type())
	elemValue.Set(zeroPtr)

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

	if v.Kind() != tpPointer {
		t.Fatalf("Expected pointer, got %v", v.Kind())
	}

	elem := v.Elem()
	if !elem.IsValid() {
		t.Fatal("Elem() should be valid for non-nil pointer")
	}

	if elem.Kind() != tpInt {
		t.Errorf("Expected int after Elem(), got %v", elem.Kind())
	}

	if elem.Int() != 42 {
		t.Errorf("Expected 42, got %d", elem.Int())
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
	if v.Kind() != tpPointer {
		t.Fatalf("Expected pointer to struct, got %v", v.Kind())
	}

	structValue := v.Elem()
	if structValue.Kind() != tpStruct {
		t.Fatalf("Expected struct after Elem(), got %v", structValue.Kind())
	}

	field0 := structValue.Field(0)
	if field0.Kind() != tpPointer {
		t.Fatalf("Expected pointer field, got %v", field0.Kind())
	}

	innerStruct := field0.Elem()
	if !innerStruct.IsValid() {
		t.Fatal("Inner struct should be valid")
	}

	if innerStruct.Kind() != tpStruct {
		t.Fatalf("Expected struct after field Elem(), got %v", innerStruct.Kind())
	}

	valueField := innerStruct.Field(0)
	if valueField.Kind() != tpInt {
		t.Fatalf("Expected int field, got %v", valueField.Kind())
	}

	if valueField.Int() != 123 {
		t.Errorf("Expected 123, got %d", valueField.Int())
	}
}

// TestInterfaceValue tests interface{} dereferencing
func TestInterfaceValue(t *testing.T) {
	var inter struct {
		E any
	}
	inter.E = 123.456
	v1 := refValueOf(&inter)
	v2 := v1.Elem().Field(0)

	// v2 should be interface{} containing float64
	v3 := v2.Elem()
	if v3.Kind() != tpFloat64 {
		t.Errorf("Expected float64 in interface, got %v", v3.Kind())
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
		v := refValueOf(tt.ptr).Elem()
		newVal := refValueOf(tt.newValue)
		v.Set(newVal)

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
			if v.Kind() != tpPointer {
				t.Errorf("Expected pointer type, got %v", v.Kind())
				return
			}

			// Elem() should work
			elem := v.Elem()
			if !elem.IsValid() {
				t.Error("Elem() should be valid for non-nil pointer")
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
			if v.Kind() != tpPointer {
				t.Errorf("Expected pointer, got %v", v.Kind())
				return
			}

			// Test type.Elem()
			elemType := v.Type().Elem()
			if elemType == nil {
				t.Error("Type().Elem() should not be nil for pointer")
				return
			}

			if elemType.Kind() != test.elemKind {
				t.Errorf("Type().Elem().Kind() expected %v, got %v", test.elemKind, elemType.Kind())
			}

			// Test value.Elem()
			elem := v.Elem()
			if !elem.IsValid() {
				t.Error("Elem() should be valid")
				return
			}

			if elem.Kind() != test.elemKind {
				t.Errorf("Elem().Kind() expected %v, got %v", test.elemKind, elem.Kind())
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
	if v.Kind() != tpStruct {
		t.Fatalf("Expected struct, got %v", v.Kind())
	}

	// Get Name field (field 0)
	nameField := v.Field(0)
	if nameField.Kind() != tpString {
		t.Errorf("Expected string field, got %v", nameField.Kind())
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
		t.Error("Field pointer should not be nil")
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
	if v.Kind() != tpStruct {
		t.Fatalf("Expected struct, got %v", v.Kind())
	}

	// Test ReadStatus field
	readStatusField := v.Field(0)
	if readStatusField.Kind() != tpString {
		t.Errorf("ReadStatus field should be string, got %v", readStatusField.Kind())
	}

	readStatus := readStatusField.String()
	if readStatus != user.ReadStatus {
		t.Errorf("ReadStatus mismatch: got %q, want %q", readStatus, user.ReadStatus)
	}

	// Test OpenStat field
	openStatField := v.Field(1)
	if openStatField.Kind() != tpString {
		t.Errorf("OpenStat field should be string, got %v", openStatField.Kind())
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
	stringField := v.Field(0)
	if stringField.flag&flagIndir != 0 {
		t.Errorf("String field should not have flagIndir set")
	}

	// Test int field (should NOT have flagIndir)
	intField := v.Field(1)
	if intField.flag&flagIndir != 0 {
		t.Errorf("Int field should not have flagIndir set")
	}

	// Test pointer field (SHOULD have flagIndir)
	pointerField := v.Field(2)
	if pointerField.flag&flagIndir == 0 {
		t.Errorf("Pointer field should have flagIndir set")
	}

	// Validate values
	if stringField.String() != s.StringField {
		t.Errorf("String field value mismatch: got %q, want %q", stringField.String(), s.StringField)
	}

	if intField.Int() != int64(s.IntField) {
		t.Errorf("Int field value mismatch: got %d, want %d", intField.Int(), s.IntField)
	}

	// Validate pointer field dereferencing
	pointerElem := pointerField.Elem()
	if !pointerElem.IsValid() {
		t.Error("Pointer field Elem() should be valid")
	} else if pointerElem.String() != *s.PointerField {
		t.Errorf("Pointer field dereferenced value mismatch: got %q, want %q", pointerElem.String(), *s.PointerField)
	}
}
