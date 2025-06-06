package tinystring

import "testing"

func TestFixedNegativeNumbers(t *testing.T) {
	// Test negative numbers in s2Int
	result, err := Convert("-123").ToInt()
	if err != nil {
		t.Errorf("ToInt(-123) failed: %v", err)
	}
	if result != -123 {
		t.Errorf("ToInt(-123) = %d, want -123", result)
	}

	// Test negative numbers in s2Int64
	result64, err := Convert("-9223372036854775807").ToInt64()
	if err != nil {
		t.Errorf("ToInt64(-9223372036854775807) failed: %v", err)
	}
	if result64 != -9223372036854775807 {
		t.Errorf("ToInt64(-9223372036854775807) = %d, want -9223372036854775807", result64)
	}

	// Test negative numbers should fail for non-decimal bases
	_, err = Convert("-123").ToInt(16)
	if err == nil {
		t.Error("ToInt(-123, base 16) should have failed but didn't")
	}
}
