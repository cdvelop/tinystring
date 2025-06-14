package tinystring

import (
	"testing"
)

// TestJsonEncodeDecode tests basic JSON encode/decode cycle for coordinates
func TestJsonEncodeDecode(t *testing.T) {
	clearRefStructsCache()

	// Test the specific structure that's failing: Coordinates
	coords := ComplexCoordinates{
		Latitude:  37.7749,
		Longitude: -122.4194,
		Accuracy:  10,
	}

	// Test encode-decode cycle directly
	jsonBytes, err := Convert(coords).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode failed: %v", err)
	}

	jsonStr := string(jsonBytes)
	expectedJson := `{"Latitude":37.7749,"Longitude":-122.4194,"Accuracy":10}`
	if jsonStr != expectedJson {
		t.Errorf("JSON mismatch: expected %s, got %s", expectedJson, jsonStr)
	}

	var decoded ComplexCoordinates
	err = Convert(jsonStr).JsonDecode(&decoded)
	if err != nil {
		t.Fatalf("JsonDecode failed: %v", err)
	} // Validate decoded values
	if decoded.Latitude != coords.Latitude {
		t.Errorf("Latitude mismatch: expected %f, got %f", coords.Latitude, decoded.Latitude)
	}
	if decoded.Longitude != coords.Longitude {
		t.Errorf("Longitude mismatch: expected %f, got %f", coords.Longitude, decoded.Longitude)
	}
	if coords.Accuracy != decoded.Accuracy {
		t.Errorf("Accuracy mismatch: expected %d, got %d", coords.Accuracy, decoded.Accuracy)
	}
}

// TestJsonPointerEncodeDecode tests JSON encode/decode with pointer to struct
func TestJsonPointerEncodeDecode(t *testing.T) {
	clearRefStructsCache()

	coords := ComplexCoordinates{
		Latitude:  37.7749,
		Longitude: -122.4194,
		Accuracy:  10,
	}

	// Test with pointer to coordinates
	ptrCoords := &coords
	jsonBytes, err := Convert(ptrCoords).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode pointer failed: %v", err)
	}

	jsonStr := string(jsonBytes)
	expectedJson := `{"Latitude":37.7749,"Longitude":-122.4194,"Accuracy":10}`
	if jsonStr != expectedJson {
		t.Errorf("Pointer JSON mismatch: expected %s, got %s", expectedJson, jsonStr)
	}

	var decodedPtr ComplexCoordinates
	err = Convert(jsonStr).JsonDecode(&decodedPtr)
	if err != nil {
		t.Fatalf("JsonDecode from pointer failed: %v", err)
	}
	// Validate decoded values
	if decodedPtr.Latitude != coords.Latitude {
		t.Errorf("Pointer Latitude mismatch: expected %f, got %f", coords.Latitude, decodedPtr.Latitude)
	}
	if decodedPtr.Longitude != coords.Longitude {
		t.Errorf("Pointer Longitude mismatch: expected %f, got %f", coords.Longitude, decodedPtr.Longitude)
	}
	if coords.Accuracy != decodedPtr.Accuracy {
		t.Errorf("Pointer Accuracy mismatch: expected %d, got %d", coords.Accuracy, decodedPtr.Accuracy)
	}
}

// TestJsonNestedStructDecode tests nested struct decoding
func TestJsonNestedStructDecode(t *testing.T) {
	clearRefStructsCache()

	// Test simple struct with embedded coordinates
	type SimpleStruct struct {
		Coords ComplexCoordinates
	}

	jsonStr := `{"Coords":{"Latitude":37.7749,"Longitude":-122.4194,"Accuracy":10}}`

	var simple SimpleStruct
	err := Convert(jsonStr).JsonDecode(&simple)
	if err != nil {
		t.Fatalf("Simple struct decode failed: %v", err)
	}

	// Validate values
	if simple.Coords.Latitude != 37.7749 {
		t.Errorf("Simple Latitude mismatch: expected %f, got %f", 37.7749, simple.Coords.Latitude)
	}
	if simple.Coords.Longitude != -122.4194 {
		t.Errorf("Simple Longitude mismatch: expected %f, got %f", -122.4194, simple.Coords.Longitude)
	}
	if simple.Coords.Accuracy != 10 {
		t.Errorf("Simple Accuracy mismatch: expected %d, got %d", 10, simple.Coords.Accuracy)
	}
}

// TestJsonPointerToStructFields tests pointer-to-struct fields in JSON decode
func TestJsonPointerToStructFields(t *testing.T) {
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

	jsonStr := `{"Name":"test","Coords":{"Lat":37.7749,"Lng":-122.4194,"Alt":100}}`

	var container TestContainer
	err := Convert(jsonStr).JsonDecode(&container)
	if err != nil {
		t.Fatalf("JsonDecode failed: %v", err)
	}

	// Validate container
	if container.Name != "test" {
		t.Errorf("Name mismatch: expected %q, got %q", "test", container.Name)
	}
	if container.Coords == nil {
		t.Fatal("Coords is nil")
	}

	// Validate coordinates
	if container.Coords.Lat != 37.7749 {
		t.Errorf("Lat mismatch: expected %f, got %f", 37.7749, container.Coords.Lat)
	}
	if container.Coords.Lng != -122.4194 {
		t.Errorf("Lng mismatch: expected %f, got %f", -122.4194, container.Coords.Lng)
	}
	if container.Coords.Alt != 100 {
		t.Errorf("Alt mismatch: expected %d, got %d", 100, container.Coords.Alt)
	}
}

// TestJsonConvertPointerHandling tests Convert() function pointer handling for JSON
func TestJsonConvertPointerHandling(t *testing.T) {
	clearRefStructsCache()

	coords := ComplexCoordinates{
		Latitude:  37.7749,
		Longitude: -122.4194,
		Accuracy:  10,
	}

	// Test direct struct
	conv1 := Convert(coords)
	if conv1.vTpe != tpStruct {
		t.Errorf("Direct struct conv type: expected %v, got %v", tpStruct, conv1.vTpe)
	}
	if !conv1.refIsValid() {
		t.Error("Direct struct refVal should be valid")
	}

	// Test pointer to struct
	ptrCoords := &coords
	conv2 := Convert(ptrCoords)
	if conv2.vTpe != tpStruct {
		t.Errorf("Pointer conv type: expected %v, got %v", tpStruct, conv2.vTpe)
	}
	if !conv2.refIsValid() {
		t.Error("Pointer conv refVal should be valid")
	}

	// Both should produce same JSON
	json1, err1 := conv1.JsonEncode()
	json2, err2 := conv2.JsonEncode()

	if err1 != nil {
		t.Fatalf("Direct struct encode failed: %v", err1)
	}
	if err2 != nil {
		t.Fatalf("Pointer struct encode failed: %v", err2)
	}
	if string(json1) != string(json2) {
		t.Errorf("JSON output mismatch: direct=%s, pointer=%s", string(json1), string(json2))
	}
}

// TestDebugJSONReflectionIssue - Critical test to diagnose JSON reflection issues
// This test replicates the exact scenario that causes string corruption in JSON encoding
func TestDebugJSONReflectionIssue(t *testing.T) {
	// Reproduce the exact struct and values from the failing test
	type ComplexUser struct {
		ReadStatus string `json:"read_status"`
		OpenStat   string `json:"open_stat"`
	}

	user := ComplexUser{
		ReadStatus: "read",
		OpenStat:   "open",
	}

	// Test the reflection chain exactly as used in JSON encoding
	v := refValueOf(user)
	if v.refKind() != tpStruct {
		t.Fatalf("Expected struct, got %v", v.refKind())
	}

	numFields := v.refNumField()
	if numFields != 2 {
		t.Errorf("Expected 2 fields, got %d", numFields)
	}

	for i := 0; i < numFields; i++ {
		field := v.refField(i)
		if !field.refIsValid() {
			t.Errorf("refField %d should be valid", i)
			continue
		}

		if field.refKind() != tpString {
			t.Errorf("refField %d expected string, got %v", i, field.refKind())
			continue
		}

		// This is where the corruption happens in JSON encoding
		strValue := field.String()

		// Check for corruption by validating string length and content
		if len(strValue) > 100 {
			t.Errorf("refField %d appears corrupted - string too long: %d", i, len(strValue))
		}

		// Validate against expected values
		expected := user.ReadStatus
		if i == 1 {
			expected = user.OpenStat
		}

		if strValue != expected {
			t.Errorf("refField %d mismatch: got %q, want %q", i, strValue, expected)
		}

		// Test direct memory access for comparison
		if field.ptr != nil {
			directValue := *(*string)(field.ptr)
			if strValue != directValue {
				t.Errorf("refField %d reflection vs direct mismatch: reflection=%q, direct=%q", i, strValue, directValue)
			}
		}
	}
}

// Debug test to understand struct encoding issue
func TestJsonDebugStruct(t *testing.T) {
	clearRefStructsCache() // Clear cache to avoid interference

	person := Person{
		Id:        "123",
		Name:      "John Doe",
		BirthDate: "1990-01-01",
		Gender:    "male",
		Phone:     "+1234567890",
	}

	// Test reflection info first
	rv := refValueOf(person)
	t.Logf("refValue kind: %v", rv.refKind())
	t.Logf("refNumField(): %d", rv.refNumField())

	for i := range rv.refNumField() {
		field := rv.refField(i)
		t.Logf("refField %d: refKind=%v, Valid=%v, String=%v", i, field.refKind(), field.refIsValid(), field.String())
	}
	// Test struct info
	var structInfo refStructInfo
	getStructInfo(rv.Type(), &structInfo)
	if structInfo.refType != nil {
		t.Logf("StructInfo fields count: %d", len(structInfo.fields))
		for i, f := range structInfo.fields {
			t.Logf("StructInfo field %d: name=%s, snakeName=%s", i, f.name, f.snakeName)
		}
	} else {
		t.Log("StructInfo.refType is nil!")
	}

	// Test actual encoding
	result, err := Convert(person).JsonEncode()
	if err != nil {
		t.Errorf("JsonEncode(Person) returned error: %v", err)
		return
	}
	t.Logf("JSON result: %s", string(result))

	// Now test Address to see if it gets different cache entry
	address := Address{
		Id:      "addr1",
		Street:  "123 Main St",
		City:    "New York",
		ZipCode: "10001",
	}

	rv2 := refValueOf(address)
	t.Logf("Address refValue kind: %v", rv2.refKind())
	t.Logf("Address refNumField(): %d", rv2.refNumField())
	var structInfo2 refStructInfo
	getStructInfo(rv2.Type(), &structInfo2)
	if structInfo2.refType != nil {
		t.Logf("Address StructInfo fields count: %d", len(structInfo2.fields))
		for i, f := range structInfo2.fields {
			t.Logf("Address StructInfo field %d: name=%s, snakeName=%s", i, f.name, f.snakeName)
		}
	}

	result2, err := Convert(address).JsonEncode()
	if err != nil {
		t.Errorf("JsonEncode(Address) returned error: %v", err)
		return
	}

	t.Logf("Address JSON result: %s", string(result2))
}

func TestDebugComplexCoordinates(t *testing.T) {
	coords := ComplexCoordinates{
		Latitude:  37.7749,
		Longitude: -122.4194,
		Accuracy:  10,
	}

	conv := Convert(coords)
	t.Logf("vTpe: %v", conv.vTpe)
	t.Logf("refKind: %v", conv.refKind())
	t.Logf("refIsValid: %v", conv.refIsValid())

	if conv.refKind() == tpStruct {
		t.Logf("NumFields: %v", conv.refNumField())
	} else {
		t.Errorf("Expected struct, got %v", conv.refKind())
	}
}
