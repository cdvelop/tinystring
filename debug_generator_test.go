package tinystring

import (
	"testing"
)

func TestDebugComplexUserFromGenerator(t *testing.T) {
	clearRefStructsCache()

	// Get the data from the generator
	testUsers := GenerateComplexTestData(1)
	user := testUsers[0]

	t.Logf("User ID: %s", user.ID)
	t.Logf("User Username: %s", user.Username)
	t.Logf("User Email: %s", user.Email)

	// Try to encode it
	jsonBytes, err := Convert(user).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode failed: %v", err)
	}

	t.Logf("Encoded JSON length: %d", len(jsonBytes))
	// Don't print the full JSON as it might be very long
}
