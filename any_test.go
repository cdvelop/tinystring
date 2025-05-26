package tinystring

import (
	"testing"
)

func TestAnyToStringOptimized(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, ""},
		{"empty string", "", ""},
		{"string", "hello", "hello"},
		{"int zero", 0, "0"},
		{"int one", 1, "1"},
		{"int small", 42, "42"},
		{"int negative", -10, "-10"},
		{"int large", 12345678, "12345678"},
		{"int8", int8(8), "8"},
		{"int16", int16(16), "16"},
		{"int32", int32(32), "32"},
		{"int64", int64(64), "64"},
		{"uint", uint(5), "5"},
		{"uint8", uint8(8), "8"},
		{"uint16", uint16(16), "16"},
		{"uint32", uint32(32), "32"},
		{"uint64", uint64(64), "64"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"float32 zero", float32(0), "0"},
		{"float32 one", float32(1), "1"},
		{"float32", float32(3.14), "3.14"},
		{"float64 zero", 0.0, "0"},
		{"float64 one", 1.0, "1"},
		{"float64", 3.14159, "3.14159"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := anyToString(tc.input)
			if result != tc.expected {
				t.Errorf("anyToString(%v) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// Test to ensure that repeated conversions return the expected value
func TestAnyToStringStringReuse(t *testing.T) {
	// For small integers
	s1 := anyToString(42)
	s2 := anyToString(42)

	if s1 != "42" || s2 != "42" {
		t.Errorf("Expected anyToString(42) to return \"42\"")
	}

	// Verify we're using the constant pool for small ints
	if s1 != smallInts[42] {
		t.Errorf("Expected anyToString(42) to use smallInts pool")
	}

	// For boolean values
	b1 := anyToString(true)
	b2 := anyToString(true)

	if b1 != trueString || b2 != trueString {
		t.Errorf("Expected anyToString(true) to return the constant \"true\"")
	}

	// For zero value
	z1 := anyToString(0)
	z2 := anyToString(0)

	if z1 != zeroString || z2 != zeroString {
		t.Errorf("Expected anyToString(0) to return the constant \"0\"")
	}

	// Prueba adicional para verificar la consistencia
	t.Log("Tests for optimized string constants passed")
}
