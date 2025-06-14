package tinystring

import (
	"testing"
)

// Test complete ComplexUser structure decoding (encode-decode cycle)
func TestJsonDecodeComplexUser_DISABLED(t *testing.T) {
	// Generate test data and encode it
	testUsers := GenerateComplexTestData(1)
	originalUser := testUsers[0]

	// Encode to JSON
	jsonBytes, err := Convert(originalUser).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode(ComplexUser) failed: %v", err)
	}

	jsonStr := string(jsonBytes)
	t.Logf("Generated JSON length: %d bytes", len(jsonStr))

	// Decode back to struct
	var decodedUser ComplexUser
	err = Convert(jsonStr).JsonDecode(&decodedUser)
	if err != nil {
		t.Fatalf("JsonDecode(ComplexUser) returned error: %v", err)
	}

	// Validate complete structure
	validateComplexUserDecoding(t, originalUser, decodedUser)
}

func validateComplexUserDecoding(t *testing.T, expected, actual ComplexUser) {
	// Top-level fields
	assertEqual(t, expected.ID, actual.ID, "ID")
	assertEqual(t, expected.Username, actual.Username, "Username")
	assertEqual(t, expected.Email, actual.Email, "Email")
	assertEqual(t, expected.CreatedAt, actual.CreatedAt, "CreatedAt")
	assertEqual(t, expected.LastLogin, actual.LastLogin, "LastLogin")
	assertEqual(t, expected.IsActive, actual.IsActive, "IsActive")

	// Profile validation
	validateComplexProfile(t, expected.Profile, actual.Profile)

	// Permissions slice
	assertSliceEqual(t, expected.Permissions, actual.Permissions, "Permissions")

	// Metadata
	validateMetadata(t, expected.Metadata, actual.Metadata)

	// Stats
	validateComplexStats(t, expected.Stats, actual.Stats)
}

func validateComplexProfile(t *testing.T, expected, actual ComplexProfile) {
	assertEqual(t, expected.FirstName, actual.FirstName, "Profile.FirstName")
	assertEqual(t, expected.LastName, actual.LastName, "Profile.LastName")
	assertEqual(t, expected.DisplayName, actual.DisplayName, "Profile.DisplayName")
	assertEqual(t, expected.Bio, actual.Bio, "Profile.Bio")
	assertEqual(t, expected.AvatarURL, actual.AvatarURL, "Profile.AvatarURL")
	assertEqual(t, expected.BirthDate, actual.BirthDate, "Profile.BirthDate")

	// Phone numbers
	if len(expected.PhoneNumbers) != len(actual.PhoneNumbers) {
		t.Errorf("PhoneNumbers length mismatch: expected %d, got %d", len(expected.PhoneNumbers), len(actual.PhoneNumbers))
		return
	}

	for i, expectedPhone := range expected.PhoneNumbers {
		if i >= len(actual.PhoneNumbers) {
			t.Errorf("Missing phone number at index %d", i)
			continue
		}
		actualPhone := actual.PhoneNumbers[i]
		assertEqual(t, expectedPhone.ID, actualPhone.ID, Format("PhoneNumbers[%d].ID", i).String())
		assertEqual(t, expectedPhone.Type, actualPhone.Type, Format("PhoneNumbers[%d].Type", i).String())
		assertEqual(t, expectedPhone.Number, actualPhone.Number, Format("PhoneNumbers[%d].Number", i).String())
		assertEqual(t, expectedPhone.Extension, actualPhone.Extension, Format("PhoneNumbers[%d].Extension", i).String())
		assertEqual(t, expectedPhone.IsPrimary, actualPhone.IsPrimary, Format("PhoneNumbers[%d].IsPrimary", i).String())
		assertEqual(t, expectedPhone.IsVerified, actualPhone.IsVerified, Format("PhoneNumbers[%d].IsVerified", i).String())
	}

	// Addresses with coordinates
	if len(expected.Addresses) != len(actual.Addresses) {
		t.Errorf("Addresses length mismatch: expected %d, got %d", len(expected.Addresses), len(actual.Addresses))
		return
	}

	for i, expectedAddr := range expected.Addresses {
		if i >= len(actual.Addresses) {
			t.Errorf("Missing address at index %d", i)
			continue
		}
		actualAddr := actual.Addresses[i]
		validateComplexAddress(t, expectedAddr, actualAddr, i)
	}

	// Social links
	if len(expected.SocialLinks) != len(actual.SocialLinks) {
		t.Errorf("SocialLinks length mismatch: expected %d, got %d", len(expected.SocialLinks), len(actual.SocialLinks))
		return
	}

	for i, expectedLink := range expected.SocialLinks {
		if i >= len(actual.SocialLinks) {
			t.Errorf("Missing social link at index %d", i)
			continue
		}
		actualLink := actual.SocialLinks[i]
		assertEqual(t, expectedLink.Platform, actualLink.Platform, Format("SocialLinks[%d].Platform", i).String())
		assertEqual(t, expectedLink.URL, actualLink.URL, Format("SocialLinks[%d].URL", i).String())
		assertEqual(t, expectedLink.Username, actualLink.Username, Format("SocialLinks[%d].Username", i).String())
		assertEqual(t, expectedLink.Verified, actualLink.Verified, Format("SocialLinks[%d].Verified", i).String())
	}

	// Preferences
	validateComplexPreferences(t, expected.Preferences, actual.Preferences)

	// Custom fields
	validateCustomFields(t, expected.CustomFields, actual.CustomFields)
}

func validateComplexAddress(t *testing.T, expected, actual ComplexAddress, index int) {
	prefix := Format("Addresses[%d]", index).String()
	assertEqual(t, expected.ID, actual.ID, prefix+".ID")
	assertEqual(t, expected.Type, actual.Type, prefix+".Type")
	assertEqual(t, expected.Street, actual.Street, prefix+".Street")
	assertEqual(t, expected.Street2, actual.Street2, prefix+".Street2")
	assertEqual(t, expected.City, actual.City, prefix+".City")
	assertEqual(t, expected.State, actual.State, prefix+".State")
	assertEqual(t, expected.Country, actual.Country, prefix+".Country")
	assertEqual(t, expected.PostalCode, actual.PostalCode, prefix+".PostalCode")
	assertEqual(t, expected.IsPrimary, actual.IsPrimary, prefix+".IsPrimary")
	assertEqual(t, expected.IsVerified, actual.IsVerified, prefix+".IsVerified")

	// Coordinates pointer handling
	if expected.Coordinates == nil && actual.Coordinates == nil {
		return // Both nil, OK
	}
	if expected.Coordinates == nil && actual.Coordinates != nil {
		t.Errorf("%s.Coordinates: expected nil, got non-nil", prefix)
		return
	}
	if expected.Coordinates != nil && actual.Coordinates == nil {
		t.Errorf("%s.Coordinates: expected non-nil, got nil", prefix)
		return
	}

	// Both non-nil, compare values
	assertEqual(t, expected.Coordinates.Latitude, actual.Coordinates.Latitude, prefix+".Coordinates.Latitude")
	assertEqual(t, expected.Coordinates.Longitude, actual.Coordinates.Longitude, prefix+".Coordinates.Longitude")
	assertEqual(t, expected.Coordinates.Accuracy, actual.Coordinates.Accuracy, prefix+".Coordinates.Accuracy")
}

func validateComplexPreferences(t *testing.T, expected, actual ComplexPreferences) {
	assertEqual(t, expected.Language, actual.Language, "Preferences.Language")
	assertEqual(t, expected.Timezone, actual.Timezone, "Preferences.Timezone")
	assertEqual(t, expected.Theme, actual.Theme, "Preferences.Theme")
	assertEqual(t, expected.Currency, actual.Currency, "Preferences.Currency")
	assertEqual(t, expected.DateFormat, actual.DateFormat, "Preferences.DateFormat")
	assertEqual(t, expected.TimeFormat, actual.TimeFormat, "Preferences.TimeFormat")

	// Notifications
	assertEqual(t, expected.Notifications.Email, actual.Notifications.Email, "Preferences.Notifications.Email")
	assertEqual(t, expected.Notifications.SMS, actual.Notifications.SMS, "Preferences.Notifications.SMS")
	assertEqual(t, expected.Notifications.Push, actual.Notifications.Push, "Preferences.Notifications.Push")
	assertEqual(t, expected.Notifications.InApp, actual.Notifications.InApp, "Preferences.Notifications.InApp")
	assertEqual(t, expected.Notifications.Marketing, actual.Notifications.Marketing, "Preferences.Notifications.Marketing")

	// Privacy
	assertEqual(t, expected.Privacy.ProfileVisibility, actual.Privacy.ProfileVisibility, "Preferences.Privacy.ProfileVisibility")
	assertEqual(t, expected.Privacy.ShowEmail, actual.Privacy.ShowEmail, "Preferences.Privacy.ShowEmail")
	assertEqual(t, expected.Privacy.ShowPhone, actual.Privacy.ShowPhone, "Preferences.Privacy.ShowPhone")
	assertEqual(t, expected.Privacy.AllowMessaging, actual.Privacy.AllowMessaging, "Preferences.Privacy.AllowMessaging")
	assertSliceEqual(t, expected.Privacy.BlockedUsers, actual.Privacy.BlockedUsers, "Preferences.Privacy.BlockedUsers")

	// Features
	assertEqual(t, expected.Features.BetaFeatures, actual.Features.BetaFeatures, "Preferences.Features.BetaFeatures")
	assertEqual(t, expected.Features.Analytics, actual.Features.Analytics, "Preferences.Features.Analytics")
	assertEqual(t, expected.Features.AdvancedSearch, actual.Features.AdvancedSearch, "Preferences.Features.AdvancedSearch")
}

func validateCustomFields(t *testing.T, expected, actual CustomFields) {
	assertEqual(t, expected.EmployeeID, actual.EmployeeID, "CustomFields.EmployeeID")
	assertEqual(t, expected.Department, actual.Department, "CustomFields.Department")
	assertEqual(t, expected.Team, actual.Team, "CustomFields.Team")
}

func validateMetadata(t *testing.T, expected, actual Metadata) {
	assertEqual(t, expected.Source, actual.Source, "Metadata.Source")
	assertEqual(t, expected.Campaign, actual.Campaign, "Metadata.Campaign")
	assertEqual(t, expected.Referrer, actual.Referrer, "Metadata.Referrer")
	assertEqual(t, expected.Score, actual.Score, "Metadata.Score")
	assertSliceEqual(t, expected.Experiments, actual.Experiments, "Metadata.Experiments")
}

func validateComplexStats(t *testing.T, expected, actual ComplexStats) {
	assertEqual(t, expected.LoginCount, actual.LoginCount, "Stats.LoginCount")
	assertEqual(t, expected.LastActivity, actual.LastActivity, "Stats.LastActivity")
	assertEqual(t, expected.SessionDuration, actual.SessionDuration, "Stats.SessionDuration")
	assertEqual(t, expected.PageViews, actual.PageViews, "Stats.PageViews")
	assertEqual(t, expected.ActionsCount, actual.ActionsCount, "Stats.ActionsCount")
	assertEqual(t, expected.SubscriptionTier, actual.SubscriptionTier, "Stats.SubscriptionTier")
	assertEqual(t, expected.StorageUsed, actual.StorageUsed, "Stats.StorageUsed")
	assertEqual(t, expected.BandwidthUsed, actual.BandwidthUsed, "Stats.BandwidthUsed")
}

// Test multiple ComplexUser array decoding
func TestJsonDecodeComplexUserArray_DISABLED(t *testing.T) {
	t.Skip("Complex user array test disabled due to memory issues - needs investigation")
}

// Test individual complex structures
func TestJsonDecodeComplexProfile(t *testing.T) {
	clearRefStructsCache()

	profile := ComplexProfile{
		FirstName:   "Alice",
		LastName:    "Johnson",
		DisplayName: "Alice J.",
		Bio:         "Data scientist and AI researcher",
		AvatarURL:   "https://example.com/alice.jpg",
		BirthDate:   "1988-07-20",
		PhoneNumbers: []ComplexPhoneNumber{
			{ID: "ph_alice_1", Type: "mobile", Number: "+1-555-888-7777", IsPrimary: true, IsVerified: true},
			{ID: "ph_alice_2", Type: "home", Number: "+1-555-666-5555", Extension: "", IsPrimary: false, IsVerified: false},
		},
		Addresses: []ComplexAddress{
			{
				ID: "addr_alice_1", Type: "home", Street: "789 Science Drive", City: "Tech City",
				State: "TX", Country: "USA", PostalCode: "75001", IsPrimary: true, IsVerified: true,
				Coordinates: &ComplexCoordinates{Latitude: 32.7767, Longitude: -96.7970, Accuracy: 8},
			},
		},
		SocialLinks: []ComplexSocialLink{
			{Platform: "researchgate", URL: "https://researchgate.net/alice", Username: "alice_research", Verified: true},
			{Platform: "twitter", URL: "https://twitter.com/alicescience", Username: "@alicescience", Verified: false},
		},
		Preferences: ComplexPreferences{
			Language:      "en-GB",
			Theme:         "auto",
			Notifications: ComplexNotificationPrefs{Email: false, Push: true, SMS: false, InApp: true, Marketing: false},
			Privacy: ComplexPrivacySettings{
				ProfileVisibility: "private",
				ShowEmail:         false,
				ShowPhone:         true,
				AllowMessaging:    false,
				BlockedUsers:      []string{"spammer1", "troll2", "bot3"},
			},
			Features: Features{BetaFeatures: true, Analytics: false, AdvancedSearch: true},
		},
		CustomFields: CustomFields{EmployeeID: "SCI001", Department: "Research", Team: "AI"},
	}

	// Encode
	jsonBytes, err := Convert(profile).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode(ComplexProfile) failed: %v", err)
	}

	// Decode
	var decodedProfile ComplexProfile
	err = Convert(string(jsonBytes)).JsonDecode(&decodedProfile)
	if err != nil {
		t.Fatalf("JsonDecode(ComplexProfile) returned error: %v", err)
	}

	// Validate
	validateComplexProfile(t, profile, decodedProfile)
}

// Test coordinates pointer handling
func TestJsonDecodeCoordinatesPointer(t *testing.T) {
	clearRefStructsCache()

	// Test nil coordinates
	addr1 := ComplexAddress{
		ID:          "test_nil",
		Street:      "No GPS Street",
		City:        "Unknown",
		Coordinates: nil,
	}

	jsonBytes1, err := Convert(addr1).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode(address with nil coordinates) failed: %v", err)
	}

	var decodedAddr1 ComplexAddress
	err = Convert(string(jsonBytes1)).JsonDecode(&decodedAddr1)
	if err != nil {
		t.Fatalf("JsonDecode(address with nil coordinates) failed: %v", err)
	}

	if decodedAddr1.Coordinates != nil {
		t.Errorf("Expected nil coordinates, got: %+v", decodedAddr1.Coordinates)
	}

	// Test valid coordinates
	addr2 := ComplexAddress{
		ID:     "test_coords",
		Street: "GPS Street",
		City:   "Located",
		Coordinates: &ComplexCoordinates{
			Latitude:  40.7589,
			Longitude: -73.9851,
			Accuracy:  12,
		},
	}

	jsonBytes2, err := Convert(addr2).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode(address with coordinates) failed: %v", err)
	}

	var decodedAddr2 ComplexAddress
	err = Convert(string(jsonBytes2)).JsonDecode(&decodedAddr2)
	if err != nil {
		t.Fatalf("JsonDecode(address with coordinates) failed: %v", err)
	}

	if decodedAddr2.Coordinates == nil {
		t.Fatalf("Expected non-nil coordinates, got nil")
	}

	assertEqual(t, addr2.Coordinates.Latitude, decodedAddr2.Coordinates.Latitude, "Coordinates.Latitude")
	assertEqual(t, addr2.Coordinates.Longitude, decodedAddr2.Coordinates.Longitude, "Coordinates.Longitude")
	assertEqual(t, addr2.Coordinates.Accuracy, decodedAddr2.Coordinates.Accuracy, "Coordinates.Accuracy")
}

// Test empty slices and null values
func TestJsonDecodeEmptySlicesAndNulls(t *testing.T) {
	clearRefStructsCache()

	emptyUser := ComplexUser{
		ID:          "empty_test",
		Username:    "empty_user",
		Email:       "empty@test.com",
		Permissions: []string{},
		Profile: ComplexProfile{
			FirstName:    "Empty",
			LastName:     "User",
			PhoneNumbers: []ComplexPhoneNumber{},
			Addresses:    []ComplexAddress{},
			SocialLinks:  []ComplexSocialLink{},
			Preferences: ComplexPreferences{
				Privacy: ComplexPrivacySettings{
					BlockedUsers: []string{},
				},
			},
		},
		Metadata: Metadata{
			Experiments: []string{},
		},
	}

	// Encode
	jsonBytes, err := Convert(emptyUser).JsonEncode()
	if err != nil {
		t.Fatalf("JsonEncode(empty user) failed: %v", err)
	}

	// Decode
	var decodedUser ComplexUser
	err = Convert(string(jsonBytes)).JsonDecode(&decodedUser)
	if err != nil {
		t.Fatalf("JsonDecode(empty user) failed: %v", err)
	}

	// Validate empty slices are preserved
	if len(decodedUser.Permissions) != 0 {
		t.Errorf("Expected empty Permissions slice, got length %d", len(decodedUser.Permissions))
	}
	if len(decodedUser.Profile.PhoneNumbers) != 0 {
		t.Errorf("Expected empty PhoneNumbers slice, got length %d", len(decodedUser.Profile.PhoneNumbers))
	}
	if len(decodedUser.Profile.Addresses) != 0 {
		t.Errorf("Expected empty Addresses slice, got length %d", len(decodedUser.Profile.Addresses))
	}
	if len(decodedUser.Profile.SocialLinks) != 0 {
		t.Errorf("Expected empty SocialLinks slice, got length %d", len(decodedUser.Profile.SocialLinks))
	}
	if len(decodedUser.Profile.Preferences.Privacy.BlockedUsers) != 0 {
		t.Errorf("Expected empty BlockedUsers slice, got length %d", len(decodedUser.Profile.Preferences.Privacy.BlockedUsers))
	}
	if len(decodedUser.Metadata.Experiments) != 0 {
		t.Errorf("Expected empty Experiments slice, got length %d", len(decodedUser.Metadata.Experiments))
	}
}

// Test error handling with invalid JSON
func TestJsonDecodeInvalidComplexJSON(t *testing.T) {
	clearRefStructsCache()

	invalidJSONTests := []string{
		// Malformed JSON
		`{"id": "user_1", "username": "test", "email": "test@example.com"`,
		// Wrong types
		`{"id": 123, "username": true, "email": ["not", "valid"]}`,
		// Missing required structure
		`{"id": "test"}`,
		// Truncated nested structure
		`{"id": "test", "profile": {"first_name": "John", "last_name":`,
		// Invalid coordinates
		`{"id": "test", "profile": {"addresses": [{"coordinates": "invalid"}]}}`,
	}

	for i, invalidJSON := range invalidJSONTests {
		var result ComplexUser
		err := Convert(invalidJSON).JsonDecode(&result)
		if err == nil {
			t.Errorf("Test %d: JsonDecode should return error for invalid JSON: %s", i, invalidJSON)
		} else {
			t.Logf("Test %d: Correctly rejected invalid JSON with error: %v", i, err)
		}
	}
}

// Test field name mapping for complex structures
func TestJsonDecodeFieldNameMapping(t *testing.T) {
	clearRefStructsCache()

	// Test with PascalCase JSON (common in APIs)
	pascalCaseJSON := `{
		"ID": "test_mapping",
		"Username": "mapper_user", 
		"Email": "mapper@example.com",
		"CreatedAt": "2024-01-01T00:00:00Z",
		"LastLogin": "2024-01-02T00:00:00Z",
		"IsActive": true,
		"Profile": {
			"FirstName": "Map",
			"LastName": "Test",
			"DisplayName": "Mapper",
			"PhoneNumbers": [
				{
					"ID": "phone_1",
					"Type": "mobile",
					"Number": "+1-555-MAP-TEST",
					"IsPrimary": true,
					"IsVerified": false
				}
			]
		}
	}`

	var user ComplexUser
	err := Convert(pascalCaseJSON).JsonDecode(&user)
	if err != nil {
		t.Fatalf("JsonDecode(PascalCase JSON) failed: %v", err)
	}

	// Validate mapping worked
	assertEqual(t, "test_mapping", user.ID, "ID mapping")
	assertEqual(t, "mapper_user", user.Username, "Username mapping")
	assertEqual(t, "mapper@example.com", user.Email, "Email mapping")
	assertEqual(t, "Map", user.Profile.FirstName, "Profile.FirstName mapping")
	assertEqual(t, "Test", user.Profile.LastName, "Profile.LastName mapping")

	if len(user.Profile.PhoneNumbers) > 0 {
		assertEqual(t, "phone_1", user.Profile.PhoneNumbers[0].ID, "PhoneNumber.ID mapping")
		assertEqual(t, "mobile", user.Profile.PhoneNumbers[0].Type, "PhoneNumber.Type mapping")
		assertEqual(t, true, user.Profile.PhoneNumbers[0].IsPrimary, "PhoneNumber.IsPrimary mapping")
	} else {
		t.Error("PhoneNumbers array is empty after decoding")
	}
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func assertEqual(t *testing.T, expected, actual interface{}, field string) {
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", field, expected, actual)
	}
}

func assertSliceEqual(t *testing.T, expected, actual []string, field string) {
	if len(expected) != len(actual) {
		t.Errorf("%s: slice length mismatch, expected %d, got %d", field, len(expected), len(actual))
		return
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			t.Errorf("%s[%d]: expected %q, got %q", field, i, expected[i], actual[i])
		}
	}
}

// ============================================================================
// LEGACY COMPATIBILITY TESTS (for backward compatibility)
// ============================================================================

// Basic type decoding tests (kept for compatibility)
func TestJsonDecodeBasicString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`""`, ""},
		{`"hello\nworld"`, "hello\nworld"},
	}

	for _, test := range tests {
		var result string
		err := Convert(test.input).JsonDecode(&result)
		if err != nil {
			t.Errorf("JsonDecode(%s) returned error: %v", test.input, err)
			continue
		}

		if result != test.expected {
			t.Errorf("JsonDecode(%s) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestJsonDecodeBasicInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"42", 42},
		{"-123", -123},
		{"0", 0},
	}

	for _, test := range tests {
		var result int64
		err := Convert(test.input).JsonDecode(&result)
		if err != nil {
			t.Errorf("JsonDecode(%s) returned error: %v", test.input, err)
			continue
		}

		if result != test.expected {
			t.Errorf("JsonDecode(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestJsonDecodeBasicBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, test := range tests {
		var result bool
		err := Convert(test.input).JsonDecode(&result)
		if err != nil {
			t.Errorf("JsonDecode(%s) returned error: %v", test.input, err)
			continue
		}

		if result != test.expected {
			t.Errorf("JsonDecode(%s) = %t, expected %t", test.input, result, test.expected)
		}
	}
}

func TestJsonDecodeInvalidJson(t *testing.T) {
	tests := []string{
		"invalid",
		`"unterminated string`,
		`{"invalid": json}`,
		`[1, 2, 3`,
	}

	for _, test := range tests {
		var result interface{}
		err := Convert(test).JsonDecode(&result)
		if err == nil {
			t.Errorf("JsonDecode(%s) should return error for invalid JSON", test)
		}
	}
}
