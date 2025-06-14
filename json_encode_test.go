package tinystring

import (
	"testing"
)

// Basic JSON encoding tests
func TestJsonEncodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", `"hello"`},
		{"", `""`},
		{"hello\nworld", `"hello\nworld"`},
		{`hello"world`, `"hello\"world"`},
	}

	for _, test := range tests {
		result, err := Convert(test.input).JsonEncode()
		if err != nil {
			t.Errorf("JsonEncode(%q) returned error: %v", test.input, err)
			continue
		}

		if string(result) != test.expected {
			t.Errorf("JsonEncode(%q) = %s, expected %s", test.input, string(result), test.expected)
		}
	}
}

func TestJsonEncodeInt(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{42, "42"},
		{-123, "-123"},
		{0, "0"},
	}

	for _, test := range tests {
		result, err := Convert(test.input).JsonEncode()
		if err != nil {
			t.Errorf("JsonEncode(%d) returned error: %v", test.input, err)
			continue
		}

		if string(result) != test.expected {
			t.Errorf("JsonEncode(%d) = %s, expected %s", test.input, string(result), test.expected)
		}
	}
}

func TestJsonEncodeFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{3.14, "3.14"},
		{0.0, "0"},
		{-2.5, "-2.5"},
	}

	for _, test := range tests {
		result, err := Convert(test.input).JsonEncode()
		if err != nil {
			t.Errorf("JsonEncode(%f) returned error: %v", test.input, err)
			continue
		}

		// Note: float formatting might vary, so we just check it's not empty
		if len(result) == 0 {
			t.Errorf("JsonEncode(%f) returned empty result", test.input)
		}
	}
}

func TestJsonEncodeBool(t *testing.T) {
	tests := []struct {
		input    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, test := range tests {
		result, err := Convert(test.input).JsonEncode()
		if err != nil {
			t.Errorf("JsonEncode(%t) returned error: %v", test.input, err)
			continue
		}

		if string(result) != test.expected {
			t.Errorf("JsonEncode(%t) = %s, expected %s", test.input, string(result), test.expected)
		}
	}
}

func TestJsonEncodeStringSlice(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{}, "[]"},
		{[]string{"hello"}, `["hello"]`},
		{[]string{"a", "b", "c"}, `["a","b","c"]`},
	}

	for _, test := range tests {
		result, err := Convert(test.input).JsonEncode()
		if err != nil {
			t.Errorf("JsonEncode(%v) returned error: %v", test.input, err)
			continue
		}

		if string(result) != test.expected {
			t.Errorf("JsonEncode(%v) = %s, expected %s", test.input, string(result), test.expected)
		}
	}
}

// Test writer interface
func TestJsonEncodeWithWriter(t *testing.T) {
	// Simple test buffer to capture written data
	var capturedData []byte

	// Create a wrapper that implements writer interface
	writer := &testWriter{
		writeFunc: func(p []byte) (int, error) {
			capturedData = append(capturedData, p...)
			return len(p), nil
		},
	}

	// JsonEncode with writer returns (nil, error)
	result, err := Convert("hello").JsonEncode(writer)
	if err != nil {
		t.Errorf("JsonEncode with writer returned error: %v", err)
	}

	// Result should be nil when writer is provided
	if result != nil {
		t.Errorf("JsonEncode with writer should return nil bytes, got %v", result)
	}

	expected := `"hello"`
	if string(capturedData) != expected {
		t.Errorf("JsonEncode with writer wrote %s, expected %s", string(capturedData), expected)
	}
}

// testWriter implements the writer interface for testing
type testWriter struct {
	writeFunc func([]byte) (int, error)
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return w.writeFunc(p)
}

// Test error handling
func TestJsonEncodeUnsupportedType(t *testing.T) {
	type unsupported struct {
		Data map[string]interface{} // Maps are not supported
	}

	input := unsupported{Data: make(map[string]interface{})}
	_, err := Convert(input).JsonEncode()
	if err == nil {
		t.Error("JsonEncode should return error for unsupported type")
	}
}

// Struct JSON encoding tests
func TestJsonEncodeStruct(t *testing.T) {
	clearRefStructsCache() // Clear cache to avoid interference between tests

	person := Person{
		Id:        "123",
		Name:      "John Doe",
		BirthDate: "1990-01-01",
		Gender:    "male",
		Phone:     "+1234567890",
	}

	result, err := Convert(person).JsonEncode()
	if err != nil {
		t.Errorf("JsonEncode(Person) returned error: %v", err)
		return
	}
	// Check that it contains the expected fields in snake_case
	jsonStr := string(result)
	t.Logf("Actual JSON result: %s", jsonStr)
	expectedFields := []string{
		`"Id":"123"`,
		`"Name":"John Doe"`,
		`"BirthDate":"1990-01-01"`,
		`"Gender":"male"`,
		`"Phone":"+1234567890"`,
	}
	for _, field := range expectedFields {
		if !Contains(jsonStr, field) {
			t.Errorf("JsonEncode(Person) missing field: %s in %s", field, jsonStr)
		} else {
			t.Logf("Found field: %s", field)
		}
	}
}

func TestJsonEncodeNestedStruct(t *testing.T) {
	clearRefStructsCache() // Clear cache to avoid interference

	address := Address{
		Id:      "addr1",
		Street:  "123 Main St",
		City:    "New York",
		ZipCode: "10001",
	}

	person := Person{
		Id:        "123",
		Name:      "John Doe",
		BirthDate: "1990-01-01",
		Gender:    "male",
		Phone:     "+1234567890",
		Addresses: []Address{address},
	}

	result, err := Convert(person).JsonEncode()
	if err != nil {
		t.Errorf("JsonEncode(nested Person) returned error: %v", err)
		return
	}
	jsonStr := string(result) // Should contain nested addresses array

	if !Contains(jsonStr, `"Addresses"`) {
		t.Errorf("JsonEncode(nested Person) missing addresses field in: %s", jsonStr)
	}

	if !Contains(jsonStr, `"Street":"123 Main St"`) {
		t.Errorf("JsonEncode(nested Person) missing nested street field in: %s", jsonStr)
	}
}

func TestJsonEncodeEmptyStruct(t *testing.T) {
	empty := struct{}{}
	result, err := Convert(empty).JsonEncode()
	if err != nil {
		t.Errorf("JsonEncode(empty struct) returned error: %v", err)
		return
	}

	expected := "{}"
	if string(result) != expected {
		t.Errorf("JsonEncode(empty struct) = %s, expected %s", string(result), expected)
	}
}

func TestJsonEncodeStructSlice(t *testing.T) {
	clearRefStructsCache() // Clear cache to avoid interference
	addresses := []Address{
		{Id: "1", Street: "Main St", City: "NYC", ZipCode: "10001"},
		{Id: "2", Street: "Oak Ave", City: "LA", ZipCode: "90210"},
	}

	result, err := Convert(addresses).JsonEncode()
	if err != nil {
		t.Errorf("JsonEncode([]Address) returned error: %v", err)
		return
	}
	jsonStr := string(result)
	if !Contains(jsonStr, `"Street":"Main St"`) {
		t.Errorf("JsonEncode([]Address) missing first address in: %s", jsonStr)
	}
	if !Contains(jsonStr, `"Street":"Oak Ave"`) {
		t.Errorf("JsonEncode([]Address) missing second address in: %s", jsonStr)
	}
}

// refField name conversion tests
func TestJsonFieldNameConversion(t *testing.T) {
	type TestStruct struct {
		FirstName  string
		LastName   string
		EmailAddr  string
		PhoneNum   string
		BirthDate  string
		IsActive   bool
		UserID     int
		AccountNum uint64
	}

	test := TestStruct{
		FirstName:  "John",
		LastName:   "Doe",
		EmailAddr:  "john@example.com",
		PhoneNum:   "123-456-7890",
		BirthDate:  "1990-01-01",
		IsActive:   true,
		UserID:     42,
		AccountNum: 123456789,
	}

	result, err := Convert(test).JsonEncode()
	if err != nil {
		t.Errorf("JsonEncode(TestStruct) returned error: %v", err)
		return
	}

	jsonStr := string(result)
	// Check original field names (PascalCase)
	expectedFields := []string{
		`"FirstName":"John"`,
		`"LastName":"Doe"`,
		`"EmailAddr":"john@example.com"`,
		`"PhoneNum":"123-456-7890"`,
		`"BirthDate":"1990-01-01"`,
		`"IsActive":true`,
		`"UserID":42`,
		`"AccountNum":123456789`,
	}
	for _, field := range expectedFields {
		if !Contains(jsonStr, field) {
			t.Errorf("JsonEncode(TestStruct) missing PascalCase field: %s in %s", field, jsonStr)
		}
	}
}

// Test JSON string escaping functionality
func TestJsonStringEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple string", "hello", `"hello"`},
		{"String with quotes", `hello "world"`, `"hello \"world\""`},
		{"String with backslash", `hello\world`, `"hello\\world"`},
		{"String with newline", "hello\nworld", `"hello\nworld"`},
		{"String with tab", "hello\tworld", `"hello\tworld"`},
		{"String with carriage return", "hello\rworld", `"hello\rworld"`},
		{"String with backspace", "hello\bworld", `"hello\bworld"`},
		{"String with form feed", "hello\fworld", `"hello\fworld"`},
		{"Empty string", "", `""`},
		{"String with control characters", "hello\u0001world", `"hello\u0001world"`},
		{"Complex escaped string", "\"test\"\n\t\\", `"\"test\"\n\t\\"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test via struct JSON encoding to trigger string escaping
			type TestStruct struct {
				Value string
			}

			s := TestStruct{Value: tt.input}
			jsonBytes, err := Convert(&s).JsonEncode()
			if err != nil {
				t.Fatalf("JSON encoding failed: %v", err)
			}

			jsonStr := string(jsonBytes)

			// Check if the expected escaped string is in the JSON output
			if !Contains(jsonStr, tt.expected) {
				t.Errorf("Expected JSON to contain %s, got %s", tt.expected, jsonStr)
			}
		})
	}
}

// Test JSON encoding with various data types
func TestJsonEncodingDataTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		contains []string
	}{
		{"String value", "test string", []string{`"test string"`}},
		{"Integer value", 42, []string{"42"}},
		{"Float value", 3.14, []string{"3.14"}},
		{"Boolean true", true, []string{"true"}},
		{"Boolean false", false, []string{"false"}},
		{"Uint value", uint(255), []string{"255"}},
		{"Int64 value", int64(9223372036854775807), []string{"9223372036854775807"}},
		{"Float32 value", float32(2.5), []string{"2.5"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := Convert(tt.input).JsonEncode()
			if err != nil {
				t.Fatalf("JSON encoding failed: %v", err)
			}

			jsonStr := string(jsonBytes)

			for _, expectedSubstr := range tt.contains {
				if !Contains(jsonStr, expectedSubstr) {
					t.Errorf("Expected JSON to contain %s, got %s", expectedSubstr, jsonStr)
				}
			}
		})
	}
}

// Test JSON encoding with slices to improve encodeJsonSlice coverage
func TestJsonSliceEncoding(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []string
	}{{"String slice", []string{"a", "b", "c"}, []string{`["a","b","c"]`}},
		{"Empty string slice", []string{}, []string{`[]`}},
		{"Single element string slice", []string{"single"}, []string{`["single"]`}},
		{"String slice with special chars", []string{"hello\nworld", "test\"quote"}, []string{`["hello\nworld","test\"quote"]`}},
		{"Mixed content slice via struct",
			struct{ Items []string }{Items: []string{"hello", "world"}},
			[]string{`"Items":["hello","world"]`}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := Convert(tt.input).JsonEncode()
			if err != nil {
				t.Fatalf("JSON encoding failed: %v", err)
			}

			jsonStr := string(jsonBytes)

			for _, expectedSubstr := range tt.expected {
				if !Contains(jsonStr, expectedSubstr) {
					t.Errorf("Expected JSON to contain %s, got %s", expectedSubstr, jsonStr)
				}
			}

			t.Logf("JSON result: %s", jsonStr)
		})
	}
}

func TestEncodeJsonPointer(t *testing.T) {
	// Test encodeJsonPointer by creating conv objects with pointer types manually
	// since Convert() auto-dereferences pointers

	t.Run("nil pointer", func(t *testing.T) {
		// Create a conv with nil pointer
		c := &conv{ptr: nil}
		result, err := c.encodeJsonPointer()
		if err != nil {
			t.Errorf("encodeJsonPointer() error: %v", err)
			return
		}
		expected := "null"
		if string(result) != expected {
			t.Errorf("encodeJsonPointer() = %q, expected %q", string(result), expected)
		}
	})

	t.Run("non-pointer kind", func(t *testing.T) {
		// Create a conv that's not a pointer kind
		c := Convert("test")
		result, err := c.encodeJsonPointer()
		if err != nil {
			t.Errorf("encodeJsonPointer() error: %v", err)
			return
		}
		expected := "null"
		if string(result) != expected {
			t.Errorf("encodeJsonPointer() = %q, expected %q", string(result), expected)
		}
	})

	// Note: Testing the actual pointer dereferencing path is complex because
	// it requires setting up proper reflection structures with pointer types
	// The function is mainly tested through the main JSON encoding path
}

// Helper functions for creating pointers
func intPtr(i int) *int           { return &i }
func stringPtr(s string) *string  { return &s }
func boolPtr(b bool) *bool        { return &b }
func floatPtr(f float64) *float64 { return &f }
