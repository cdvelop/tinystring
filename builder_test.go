package tinystring

import "testing"

// TestConvertVariadicValidation tests Convert() parameter validation
func TestConvertVariadicValidation(t *testing.T) {
	// Valid usage
	c1 := Convert()        // Empty - should work
	c2 := Convert("hello") // Single value - should work

	if c1.err != "" {
		t.Errorf("Convert() should not have error, got: %s", c1.err)
	}
	if c2.err != "" {
		t.Errorf("Convert(value) should not have error, got: %s", c2.err)
	}

	// Clean up
	c1.putConv()
	c2.putConv()

	// Invalid usage - should set error and continue chain
	c3 := Convert("hello", "world") // Multiple values - should set error
	if c3.err == "" {
		t.Error("Convert with multiple parameters should set error")
	}

	// Chain should continue but operations should be omitted due to error
	result := c3.Write(" more").String() // This auto-releases
	if result != "" {
		t.Errorf("Operations after error should be omitted, got: %s", result)
	}
}

// TestWriteMethod tests the unified Write method
func TestWriteMethod(t *testing.T) {
	tests := []struct {
		name     string
		values   []any
		expected string
	}{
		{"String values", []any{"hello", " ", "world"}, "hello world"},
		{"Mixed types", []any{"Count: ", 42, " items"}, "Count: 42 items"},
		{"Boolean values", []any{"Active: ", true, ", Valid: ", false}, "Active: true, Valid: false"},
		{"Float values", []any{"Price: $", 19.99}, "Price: $19.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Convert()
			for _, v := range tt.values {
				c.Write(v)
			}
			result := c.String() // Auto-releases

			if result != tt.expected {
				t.Errorf("Write chain failed: got %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestResetMethod tests the Reset functionality
func TestResetMethod(t *testing.T) {
	c := Convert("initial")
	c.Write(" more")

	// Reset and reuse
	c.Reset()
	c.Write("new").Write(" content")
	result := c.String() // Auto-releases

	expected := "new content"
	if result != expected {
		t.Errorf("Reset failed: got %q, want %q", result, expected)
	}
}

// TestErrorChainInterruption tests error chain interruption behavior
func TestErrorChainInterruption(t *testing.T) {
	// Test normal case first
	c := Convert("valid")
	c.Write("ok")
	result := c.String() // Auto-releases
	expected := "validok"
	if result != expected {
		t.Errorf("Normal chain failed: got %q, want %q", result, expected)
	}

	// Test error case
	c2 := Convert("hello", "world") // This should set error
	if c2.err == "" {
		t.Error("Expected error for multiple parameters, got none")
	}

	c2.Write(" more") // This should be omitted due to error

	result2, err := c2.StringError()
	if err == nil {
		t.Error("Expected error from StringError(), got nil")
	}
	// When there's an error, result should be empty string
	if result2 != "" {
		t.Errorf("Expected empty result due to error, got: %s", result2)
	}
}

// TestBuilderPattern tests the main optimization goal: empty Convert() for loops
func TestBuilderPattern(t *testing.T) {
	items := []string{"  APPLE  ", "  banana  ", "  Cherry  "}

	// Test builder pattern with transformations
	c := Convert() // Empty initialization
	for i, item := range items {
		c.Write(item).Trim().ToLower().Capitalize()
		if i < len(items)-1 {
			c.Write(" - ")
		}
	}
	result := c.String() // Auto-releases

	expected := "Apple - Banana - Cherry"
	if result != expected {
		t.Errorf("Builder pattern failed: got %q, want %q", result, expected)
	}

	// Test simple pattern too
	c2 := Convert() // Empty initialization
	for _, item := range []string{"apple", "banana", "cherry"} {
		c2.Write(item).Write(" ")
	}
	result2 := c2.String() // Auto-releases

	expected2 := "apple banana cherry "
	if result2 != expected2 {
		t.Errorf("Simple builder pattern failed: got %q, want %q", result2, expected2)
	}
}

// TestTFunction tests the T translation function
func TestTFunction(t *testing.T) {
	// Test basic translation
	result := T(D.Invalid, D.Value)
	if result == "" {
		t.Error("T function returned empty string")
	}

	// Test with language
	result2 := T(ES, D.Invalid, D.Value)
	if result2 == "" {
		t.Error("T function with language returned empty string")
	}

	// They should be different (English vs Spanish)
	if result == result2 {
		t.Error("T function should return different translations for different languages")
	}
}

// TestErrFunction tests the refactored Err function
func TestErrFunction(t *testing.T) {
	// Test basic error creation
	err := Err(D.Invalid, D.Fmt)
	if err.err == "" {
		t.Error("Err function should create error message")
	}

	// Test that it uses pool
	if err.vTpe != typeErr {
		t.Error("Err should set type to typeErr")
	}

	// Clean up
	err.putConv()
}
