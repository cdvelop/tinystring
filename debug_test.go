package tinystring

import (
	"testing"
)

func TestDebugComplexUserEncode(t *testing.T) {
	clearRefStructsCache()

	// Create a simple ComplexUser to test
	user := ComplexUser{
		ID:       "test_id",
		Username: "test_user",
		Email:    "test@example.com",
	}

	// Try to encode it
	jsonBytes, err := Convert(user).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode failed: %v", err)
	}

	t.Logf("Encoded JSON: %s", string(jsonBytes))
}
