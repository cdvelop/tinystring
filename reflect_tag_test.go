package tinystring

import (
	"testing"
)

// Test the refStructTag implementation
func TestRefStructTag(t *testing.T) {
	// Test basic tag parsing
	tag := refStructTag("json:\"user_name,omitempty\" xml:\"UserName\"")

	// Test Get method
	jsonValue := tag.Get("json")
	expectedJson := "user_name,omitempty"
	if jsonValue != expectedJson {
		t.Errorf("Expected json tag '%s', got '%s'", expectedJson, jsonValue)
	}

	t.Logf("Basic tag test passed: got '%s'", jsonValue)
}

// Test JSON field mapping with tags
func TestJsonFieldMappingWithTags(t *testing.T) {
	clearRefStructsCache()

	// Test struct with JSON tags
	type TestUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	// Test with JSON that uses the tag names
	jsonStr := `{"id": "test_123", "username": "testuser", "email": "test@example.com"}`

	var user TestUser
	err := Convert(jsonStr).JsonDecode(&user)

	t.Logf("JSON: %s", jsonStr)
	t.Logf("Error: %v", err)
	t.Logf("Decoded: %+v", user)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify fields were populated correctly
	if user.ID != "test_123" {
		t.Errorf("Expected ID 'test_123', got '%s'", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected Email 'test@example.com', got '%s'", user.Email)
	}
}

// Test type validation errors
func TestJsonTypeValidationErrors(t *testing.T) {
	clearRefStructsCache()

	type SimpleUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	// Test with wrong types - should fail
	testCases := []struct {
		name       string
		jsonStr    string
		shouldFail bool
	}{
		{"valid JSON", `{"id": "123", "username": "test", "email": "test@example.com"}`, false},
		{"ID as number", `{"id": 123, "username": "test", "email": "test@example.com"}`, true},
		{"username as boolean", `{"id": "test", "username": true, "email": "test@example.com"}`, true},
		{"email as array", `{"id": "test", "username": "test", "email": ["not", "valid"]}`, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var user SimpleUser
			err := Convert(tc.jsonStr).JsonDecode(&user)

			if tc.shouldFail {
				if err == nil {
					t.Errorf("Expected error for %s, but got none. Result: %+v", tc.name, user)
				} else {
					t.Logf("Correctly rejected %s with error: %v", tc.name, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for %s, but got: %v", tc.name, err)
				} else {
					t.Logf("Correctly accepted %s. Result: %+v", tc.name, user)
				}
			}
		})
	}
}
