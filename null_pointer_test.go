package fmt

import "testing"

// TestNullPointerProtection tests null pointer verification for *string
func TestNullPointerProtection(t *testing.T) {
	// Test null string pointer
	var nullStrPtr *string = nil

	// This should not panic and should set an error
	c := Convert(nullStrPtr)
	if !c.hasContent(BuffErr) {
		t.Error("Convert with null *string should set an error")
	}

	// String() should return empty due to error
	result := c.String()
	if result != "" {
		t.Errorf("Convert with null *string should return empty string, got: %q", result)
	}
}

// TestValidPointerHandling tests that valid pointers still work
func TestValidPointerHandling(t *testing.T) {
	str := "test string"
	strPtr := &str

	// This should work normally
	c := Convert(strPtr)
	if c.hasContent(BuffErr) {
		t.Error("Convert with valid *string should not set an error")
	}

	result := c.String()
	expected := "test string"
	if result != expected {
		t.Errorf("Convert with valid *string: got %q, want %q", result, expected)
	}
}
