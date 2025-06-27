package tinystring

import (
	"testing"
)

func TestDebugRoundWithoutDecimal(t *testing.T) {
	// Test case: Convert(3).RoundDecimals(2) should be "3.00"

	// Step 1: Initial conversion
	c1 := Convert(3)
	initial := c1.String()
	t.Logf("Step 1 - Convert(3): %q", initial)

	// Step 2: RoundDecimals
	c2 := Convert(3)
	c2.RoundDecimals(2)
	rounded := c2.String()
	t.Logf("Step 2 - RoundDecimals(2): %q", rounded)

	// Expected: "3.00"
	expected := "3.00"
	if rounded != expected {
		t.Errorf("Expected %q, got %q", expected, rounded)
	}
}
