package tinystring

import (
	"testing"
)

// Basic JSON decoding tests
func TestJsonDecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`""`, ""},
		{`"hello\nworld"`, "hello\nworld"},
	}

	for _, test := range tests {
		var result string
		err := Convert(test.input).JsonDecode(&result)
		if err != nil {
			t.Errorf("JsonDecode(%s) returned error: %v", test.input, err)
			continue
		}

		if result != test.expected {
			t.Errorf("JsonDecode(%s) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestJsonDecodeInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"42", 42},
		{"-123", -123},
		{"0", 0},
	}

	for _, test := range tests {
		var result int64
		err := Convert(test.input).JsonDecode(&result)
		if err != nil {
			t.Errorf("JsonDecode(%s) returned error: %v", test.input, err)
			continue
		}

		if result != test.expected {
			t.Errorf("JsonDecode(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestJsonDecodeBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, test := range tests {
		var result bool
		err := Convert(test.input).JsonDecode(&result)
		if err != nil {
			t.Errorf("JsonDecode(%s) returned error: %v", test.input, err)
			continue
		}

		if result != test.expected {
			t.Errorf("JsonDecode(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestJsonDecodeInvalidJson(t *testing.T) {
	tests := []string{
		"invalid",
		`"unterminated string`,
		`{"invalid": json}`,
		`[1, 2, 3`,
	}

	for _, test := range tests {
		var result interface{}
		err := Convert(test).JsonDecode(&result)
		if err == nil {
			t.Errorf("JsonDecode(%s) should return error for invalid JSON", test)
		}
	}
}

// Struct JSON decoding tests
func TestJsonDecodeStruct(t *testing.T) {
	jsonStr := `{"id":"123","name":"John Doe","birth_date":"1990-01-01","gender":"male","phone":"+1234567890","addresses":[]}`

	var person Person
	err := Convert(jsonStr).JsonDecode(&person)
	if err != nil {
		t.Errorf("JsonDecode(Person) returned error: %v", err)
		return
	}

	if person.Id != "123" {
		t.Errorf("JsonDecode(Person).Id = %s, expected 123", person.Id)
	}
	if person.Name != "John Doe" {
		t.Errorf("JsonDecode(Person).Name = %s, expected John Doe", person.Name)
	}
	if person.BirthDate != "1990-01-01" {
		t.Errorf("JsonDecode(Person).BirthDate = %s, expected 1990-01-01", person.BirthDate)
	}
}

func TestJsonDecodeNestedStruct(t *testing.T) {
	jsonStr := `{"id":"123","name":"John Doe","birth_date":"1990-01-01","gender":"male","phone":"+1234567890","addresses":[{"id":"addr1","street":"123 Main St","city":"New York","zip_code":"10001"}]}`

	var person Person
	err := Convert(jsonStr).JsonDecode(&person)
	if err != nil {
		t.Errorf("JsonDecode(nested Person) returned error: %v", err)
		return
	}

	if len(person.Addresses) != 1 {
		t.Errorf("JsonDecode(nested Person).Addresses length = %d, expected 1", len(person.Addresses))
		return
	}

	addr := person.Addresses[0]
	if addr.Street != "123 Main St" {
		t.Errorf("JsonDecode(nested Person).Addresses[0].Street = %s, expected 123 Main St", addr.Street)
	}
	if addr.City != "New York" {
		t.Errorf("JsonDecode(nested Person).Addresses[0].City = %s, expected New York", addr.City)
	}
}

func TestJsonDecodeStructSlice(t *testing.T) {
	jsonStr := `[{"id":"1","street":"Main St","city":"NYC","zip_code":"10001"},{"id":"2","street":"Oak Ave","city":"LA","zip_code":"90210"}]`

	var addresses []Address
	err := Convert(jsonStr).JsonDecode(&addresses)
	if err != nil {
		t.Errorf("JsonDecode([]Address) returned error: %v", err)
		return
	}

	if len(addresses) != 2 {
		t.Errorf("JsonDecode([]Address) length = %d, expected 2", len(addresses))
		return
	}

	if addresses[0].Street != "Main St" {
		t.Errorf("JsonDecode([]Address)[0].Street = %s, expected Main St", addresses[0].Street)
	}
	if addresses[1].Street != "Oak Ave" {
		t.Errorf("JsonDecode([]Address)[1].Street = %s, expected Oak Ave", addresses[1].Street)
	}
}

// Safe test to debug struct slice decoding without memory issues
func TestSafeStructSliceDecoding(t *testing.T) {
	clearObjectCache() // Clear cache to avoid interference

	// Simple test with very basic JSON
	jsonStr := `[{"id":"1","street":"Main St"}]`

	var addresses []Address
	err := Convert(jsonStr).JsonDecode(&addresses)

	if err != nil {
		t.Logf("JsonDecode returned error: %v", err)
		return
	}

	// Check basic results without printing potentially corrupted values
	if len(addresses) != 1 {
		t.Errorf("Expected 1 address, got %d", len(addresses))
		return
	}

	// Only check length of strings to avoid printing corrupted data
	addr := addresses[0]
	idLen := len(addr.Id)
	streetLen := len(addr.Street)

	t.Logf("Address parsed successfully: Id length=%d, Street length=%d", idLen, streetLen)

	if idLen == 0 {
		t.Error("Address.Id is empty")
	}
	if streetLen == 0 {
		t.Error("Address.Street is empty")
	}
}

// Debug test for struct slice decoding
func TestDebugStructSliceDecoding(t *testing.T) {
	clearObjectCache() // Clear cache to avoid interference

	// Test with a simple Address struct
	addr := Address{}
	addressValue := refValueOf(addr)
	addressType := addressValue.Type()
	t.Logf("Address type kind: %v", addressType.Kind())

	// Test creating zero value
	zeroAddr := refZero(addressType)
	t.Logf("Zero Address kind: %v, valid: %v", zeroAddr.Kind(), zeroAddr.IsValid())

	// Test the struct slice parsing function
	elements := []string{
		`{"id":"1","street":"Main St","city":"NYC","zip_code":"10001"}`,
		`{"id":"2","street":"Oak Ave","city":"LA","zip_code":"90210"}`,
	}

	target := refValueOf(&[]Address{})
	targetElem := target.Elem()
	t.Logf("Target elem kind: %v", targetElem.Kind())

	c := Convert("")
	err := c.parseStructSlice(elements, targetElem, addressType)
	if err != nil {
		t.Errorf("parseStructSlice returned error: %v", err)
		return
	}

	t.Logf("Parsing completed successfully")
}

// Debug test for struct slice decoding
func TestJsonDecodeDebugStructSlice(t *testing.T) {
	clearObjectCache() // Clear cache

	// Single address test
	jsonStr := `{"id":"addr1","street":"123 Main St","city":"New York","zip_code":"10001"}`

	var addr Address
	err := Convert(jsonStr).JsonDecode(&addr)
	if err != nil {
		t.Errorf("Single address decode error: %v", err)
		return
	}

	t.Logf("Single address decoded: Street=%s, City=%s", addr.Street, addr.City)

	// Array of addresses test
	arrayJsonStr := `[{"id":"addr1","street":"123 Main St","city":"New York","zip_code":"10001"}]`

	var addresses []Address
	err = Convert(arrayJsonStr).JsonDecode(&addresses)
	if err != nil {
		t.Errorf("Address array decode error: %v", err)
		return
	}

	if len(addresses) != 1 {
		t.Errorf("Expected 1 address, got %d", len(addresses))
		return
	}

	addr = addresses[0]
	t.Logf("Array address decoded: Street=%s, City=%s", addr.Street, addr.City)
}

func TestStructSliceDecodingDebug(t *testing.T) {
	clearObjectCache() // Clear cache to avoid pollution
	// Simple test with one address
	jsonStr := `[{"Street":"123 Main St","City":"Anytown"}]`

	var addresses []Address
	err := Convert(jsonStr).JsonDecode(&addresses)

	if err != nil {
		t.Logf("Decode error: %v", err)
		t.Fatal("Failed to decode address slice")
	}

	t.Logf("Decoded addresses: %+v", addresses)
	t.Logf("Number of addresses: %d", len(addresses))

	if len(addresses) != 1 {
		t.Fatalf("Expected 1 address, got %d", len(addresses))
	}

	addr := addresses[0]
	t.Logf("Address Street: '%s'", addr.Street)
	t.Logf("Address City: '%s'", addr.City)

	if addr.Street != "123 Main St" {
		t.Errorf("Expected Street '123 Main St', got '%s'", addr.Street)
	}

	if addr.City != "Anytown" {
		t.Errorf("Expected City 'Anytown', got '%s'", addr.City)
	}
}

// Debug test for field name mapping
func TestFieldNameMappingDebug(t *testing.T) {
	clearObjectCache() // Clear cache

	// Test Address struct
	addr := Address{}
	rv := refValueOf(addr)

	t.Logf("Address type kind: %v", rv.Kind())

	structInfo := getStructInfo(rv.Type())
	if structInfo != nil {
		t.Logf("Found struct info with %d fields", len(structInfo.fields))
		for i, f := range structInfo.fields {
			t.Logf("Field %d: name='%s', jsonName='%s'", i, f.name, f.jsonName)
		}
	} else {
		t.Log("structInfo is nil")
	}
	// Now test parsing with the actual field names that are in the JSON
	jsonStr := `{"Street":"123 Main St","City":"Anytown"}`
	t.Logf("Testing JSON: %s", jsonStr)

	var testAddr Address
	err := Convert(jsonStr).JsonDecode(&testAddr)
	if err != nil {
		t.Logf("Decode error: %v", err)
	} else {
		t.Logf("Decoded successfully: Street='%s', City='%s'", testAddr.Street, testAddr.City)
	}
}
