package tinystring

import (
	"testing"
)

func TestBasicTypeReflection(t *testing.T) {
	tests := []struct {
		name          string
		value         any
		expectedKind  kind
		expectedValid bool
	}{
		{"string", "hello world", tpString, true},
		{"int", int(42), tpInt, true},
		{"int64", int64(42), tpInt64, true},
		{"float64", float64(3.14), tpFloat64, true},
		{"bool", true, tpBool, true},
		{"nil", nil, tpInvalid, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := refValueOf(test.value)

			// Test validity
			if got := v.refIsValid(); got != test.expectedValid {
				t.Errorf("refIsValid() = %v, want %v", got, test.expectedValid)
			}

			if !test.expectedValid {
				return // Skip further tests for invalid values
			}

			// Test kind detection
			if got := v.refKind(); got != test.expectedKind {
				t.Errorf("refKind() = %v, want %v", got, test.expectedKind)
			}

			// Test type consistency
			if v.typ == nil {
				t.Error("typ should not be nil for valid value")
				return
			}

			if got := v.refKind(); got != test.expectedKind {
				t.Errorf("typ.refKind() = %v, want %v", got, test.expectedKind)
			}
		})
	}
}

func TestStringValueRetrieval(t *testing.T) {
	original := "hello world"
	v := refValueOf(original)

	// Validate basic properties
	if !v.refIsValid() {
		t.Fatal("refValue should be valid for string")
	}

	if v.refKind() != tpString {
		t.Fatalf("refKind() = %v, want %v", v.refKind(), tpString)
	}

	// Test String() method
	result := v.String()
	if result != original {
		t.Errorf("String() = %q, want %q", result, original)
	}
}

func TestIntValueRetrieval(t *testing.T) {
	original := int64(42)
	v := refValueOf(original)

	// Validate basic properties
	if !v.refIsValid() {
		t.Fatal("refValue should be valid for int64")
	}

	if v.refKind() != tpInt64 {
		t.Fatalf("refKind() = %v, want %v", v.refKind(), tpInt64)
	}

	// Test refInt() method
	result := v.refInt()
	if result != original {
		t.Errorf("refInt() = %d, want %d", result, original)
	}
}

func TestFlagIndirCorrectness(t *testing.T) {
	tests := []struct {
		name              string
		value             any
		expectedFlagIndir bool
		reason            string
	}{
		{
			name:              "string_direct",
			value:             "hello",
			expectedFlagIndir: false,
			reason:            "basic types stored directly in interface should not have flagIndir",
		},
		{
			name:              "int_direct",
			value:             int(42),
			expectedFlagIndir: false,
			reason:            "basic types stored directly in interface should not have flagIndir",
		},
		{
			name:              "large_struct",
			value:             struct{ A, B, C, D, E int64 }{1, 2, 3, 4, 5},
			expectedFlagIndir: true,
			reason:            "large structs are stored indirectly in interface",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := refValueOf(test.value)
			hasFlagIndir := v.flag&flagIndir != 0

			if hasFlagIndir != test.expectedFlagIndir {
				t.Errorf("flagIndir = %v, want %v - %s", hasFlagIndir, test.expectedFlagIndir, test.reason)

				// Additional debug info
				t.Logf("Value: %+v", test.value)
				t.Logf("Type kind: %v", v.refKind())
				if v.typ != nil {
					t.Logf("Type size: %d", v.typ.Size())
					t.Logf("kindDirectIface: %t", v.typ.kind&kindDirectIface != 0)
					t.Logf("ifaceIndir: %t", ifaceIndir(v.typ))
				}
			}
		})
	}
}
