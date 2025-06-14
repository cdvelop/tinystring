package tinystring

import (
	"testing"
)

func TestJsonDecodeSimpleStruct(t *testing.T) {
	clearRefStructsCache()

	// Simple struct for testing
	type SimpleUser struct {
		ID   string
		Name string
		Age  int
	}

	// Simple JSON to decode
	jsonStr := `{"ID":"123","Name":"John","Age":30}`

	// Try to decode it
	var user SimpleUser
	err := Convert(jsonStr).JsonDecode(&user)
	if err != nil {
		t.Fatalf("JsonDecode failed: %v", err)
	}

	// Validate the result
	if user.ID != "123" {
		t.Errorf("ID mismatch: expected %q, got %q", "123", user.ID)
	}
	if user.Name != "John" {
		t.Errorf("Name mismatch: expected %q, got %q", "John", user.Name)
	}
	if user.Age != 30 {
		t.Errorf("Age mismatch: expected %d, got %d", 30, user.Age)
	}
}
