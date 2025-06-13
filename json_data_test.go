package tinystring

// ============================================================================
// SHARED TEST DATA STRUCTURES
// ============================================================================
// This file contains all shared data structures and test data generators
// used by both json_encode_test.go and json_decode_test.go to avoid duplication.

// Metadata contains tracking and analytics information
type Metadata struct {
	Source      string
	Campaign    string
	Referrer    string
	Experiments []string
	Score       float64
}

// CustomFields contains organization-specific fields
type CustomFields struct {
	EmployeeID string
	Department string
	Team       string
}

// Features contains feature flags and settings
type Features struct {
	BetaFeatures   bool
	Analytics      bool
	AdvancedSearch bool
}

// ComplexUser represents a complete user profile with all nested structures
type ComplexUser struct {
	ID          string
	Username    string
	Email       string
	CreatedAt   string
	LastLogin   string
	IsActive    bool
	Profile     ComplexProfile
	Permissions []string
	Metadata    Metadata
	Stats       ComplexStats
}

// ComplexProfile contains detailed user profile information
type ComplexProfile struct {
	FirstName    string
	LastName     string
	DisplayName  string
	Bio          string
	AvatarURL    string
	BirthDate    string
	PhoneNumbers []ComplexPhoneNumber
	Addresses    []ComplexAddress
	SocialLinks  []ComplexSocialLink
	Preferences  ComplexPreferences
	CustomFields CustomFields
}

// ComplexAddress represents a physical address with optional coordinates
type ComplexAddress struct {
	ID          string
	Type        string
	Street      string
	Street2     string
	City        string
	State       string
	Country     string
	PostalCode  string
	Coordinates *ComplexCoordinates
	IsPrimary   bool
	IsVerified  bool
}

// ComplexCoordinates represents GPS coordinates
type ComplexCoordinates struct {
	Latitude  float64
	Longitude float64
	Accuracy  int
}

// ComplexPhoneNumber represents a phone number with metadata
type ComplexPhoneNumber struct {
	ID         string
	Type       string
	Number     string
	Extension  string
	IsPrimary  bool
	IsVerified bool
}

// ComplexSocialLink represents a social media profile link
type ComplexSocialLink struct {
	Platform string
	URL      string
	Username string
	Verified bool
}

// ComplexPreferences contains user preferences and settings
type ComplexPreferences struct {
	Language      string
	Timezone      string
	Theme         string
	Currency      string
	DateFormat    string
	TimeFormat    string
	Notifications ComplexNotificationPrefs
	Privacy       ComplexPrivacySettings
	Features      Features
}

// ComplexNotificationPrefs contains notification preferences
type ComplexNotificationPrefs struct {
	Email     bool
	SMS       bool
	Push      bool
	InApp     bool
	Marketing bool
}

// ComplexPrivacySettings contains privacy and security settings
type ComplexPrivacySettings struct {
	ProfileVisibility string
	ShowEmail         bool
	ShowPhone         bool
	AllowMessaging    bool
	BlockedUsers      []string
}

// ComplexStats contains user activity and usage statistics
type ComplexStats struct {
	LoginCount       int64
	LastActivity     string
	SessionDuration  int64
	PageViews        int64
	ActionsCount     int64
	SubscriptionTier string
	StorageUsed      int64
	BandwidthUsed    int64
}

// Legacy simple structures for compatibility
type Person struct {
	Id        string
	Name      string
	BirthDate string
	Gender    string
	Phone     string
	Addresses []Address
}

type Address struct {
	Id      string
	Street  string
	City    string
	ZipCode string
}

// ============================================================================
// SHARED TEST DATA GENERATORS
// ============================================================================

// GenerateComplexTestData generates test data for both encoding and decoding tests
// This replaces both generateComplexTestData and generateComplexTestDataForDecode
func GenerateComplexTestData(count int) []ComplexUser {
	users := make([]ComplexUser, count)
	for i := 0; i < count; i++ {
		users[i] = ComplexUser{
			ID:        Format("user_%d", i).String(),
			Username:  Format("user_%d_2024", i).String(),
			Email:     Format("user%d@example.com", i).String(),
			CreatedAt: "2024-06-12T10:00:00Z",
			LastLogin: "2024-06-05T10:00:00Z",
			IsActive:  true,
			Profile: ComplexProfile{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "Johnny D",
				Bio:         "Software engineer passionate about technology and innovation",
				AvatarURL:   "https://cdn.example.com/avatars/john_doe.jpg",
				BirthDate:   "1990-01-01",
				PhoneNumbers: []ComplexPhoneNumber{
					{ID: "ph_001", Type: "mobile", Number: "+1-555-123-4567", IsPrimary: true, IsVerified: true},
					{ID: "ph_002", Type: "work", Number: "+1-555-987-6543", Extension: "1234", IsPrimary: false, IsVerified: false},
				},
				Addresses: []ComplexAddress{
					{
						ID: "addr_001", Type: "home", Street: "123 Main Street", City: "Anytown",
						State: "CA", Country: "USA", PostalCode: "12345", IsPrimary: true, IsVerified: true,
						Coordinates: &ComplexCoordinates{Latitude: 37.7749, Longitude: -122.4194, Accuracy: 10},
					},
				},
				SocialLinks: []ComplexSocialLink{
					{Platform: "twitter", URL: "https://twitter.com/johndoe", Username: "@johndoe", Verified: false},
					{Platform: "linkedin", URL: "https://linkedin.com/in/johndoe", Username: "johndoe", Verified: true},
				},
				Preferences: ComplexPreferences{
					Language: "en-US", Timezone: "America/Los_Angeles", Theme: "light", Currency: "USD",
					DateFormat: "MM/DD/YYYY", TimeFormat: "12h",
					Notifications: ComplexNotificationPrefs{
						Email: true, SMS: false, Push: true, InApp: true, Marketing: false,
					},
					Privacy: ComplexPrivacySettings{
						ProfileVisibility: "friends", ShowEmail: false, ShowPhone: false,
						AllowMessaging: true, BlockedUsers: []string{},
					},
					Features: Features{
						BetaFeatures: true, Analytics: true, AdvancedSearch: false,
					},
				},
				CustomFields: CustomFields{
					EmployeeID: "EMP001", Department: "Engineering", Team: "Backend",
				},
			},
			Permissions: []string{"read", "write", "admin"},
			Metadata: Metadata{
				Source: "web_signup", Campaign: "summer_2024", Referrer: "google",
				Experiments: []string{"new_ui", "faster_search"}, Score: 85.7,
			},
			Stats: ComplexStats{
				LoginCount: 1247, LastActivity: "2024-06-12T08:00:00Z",
				SessionDuration: 3600, PageViews: 15643, ActionsCount: 892,
				SubscriptionTier: "premium", StorageUsed: 2147483648, BandwidthUsed: 10737418240,
			},
		}
	}
	return users
}

// GenerateComplexProfileForTest generates a specific complex profile for individual testing
func GenerateComplexProfileForTest() ComplexProfile {
	return ComplexProfile{
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
}

// GenerateAddressWithNilCoordinates generates address with nil coordinates for testing
func GenerateAddressWithNilCoordinates() ComplexAddress {
	return ComplexAddress{
		ID:          "test_nil",
		Street:      "No GPS Street",
		City:        "Unknown",
		Coordinates: nil,
	}
}

// GenerateAddressWithCoordinates generates address with valid coordinates for testing
func GenerateAddressWithCoordinates() ComplexAddress {
	return ComplexAddress{
		ID:     "test_coords",
		Street: "GPS Street",
		City:   "Located",
		Coordinates: &ComplexCoordinates{
			Latitude:  40.7589,
			Longitude: -73.9851,
			Accuracy:  12,
		},
	}
}

// GenerateEmptyComplexUser generates a user with empty slices for testing
func GenerateEmptyComplexUser() ComplexUser {
	return ComplexUser{
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
}

// GenerateInvalidTestData generates invalid JSON data for error tests
func GenerateInvalidTestData() []string {
	return []string{
		// Malformed JSON
		`{"id": "user_1", "username": "john_doe", "email": "john@example.com"`,
		// Wrong types
		`{"id": 123, "username": true, "email": ["not", "an", "email"]}`,
		// Unexpected nulls
		`{"id": null, "username": null, "email": null, "profile": null}`,
		// Truncated JSON
		`{"id": "user_1", "profile": {"first_name": "John", "last_name":`,
		// Incomplete structure
		`{}`,
		// Missing required structure
		`{"id": "test"}`,
		// Invalid coordinates
		`{"id": "test", "profile": {"addresses": [{"coordinates": "invalid"}]}}`,
	}
}

// GeneratePascalCaseJSON generates PascalCase JSON for field mapping tests
func GeneratePascalCaseJSON() string {
	return `{
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
}

// ============================================================================
// SIMPLE TEST DATA FOR LEGACY COMPATIBILITY
// ============================================================================

// GenerateSimplePersonData generates simple Person data for basic tests
func GenerateSimplePersonData() Person {
	return Person{
		Id:        "person_001",
		Name:      "John Smith",
		BirthDate: "1980-01-15",
		Gender:    "male",
		Phone:     "+1-555-123-4567",
		Addresses: []Address{
			{Id: "addr_001", Street: "123 Main St", City: "Anytown", ZipCode: "12345"},
			{Id: "addr_002", Street: "456 Oak Ave", City: "Other City", ZipCode: "67890"},
		},
	}
}

// GenerateSimplePersonArray generates an array of Person data
func GenerateSimplePersonArray(count int) []Person {
	persons := make([]Person, count)
	for i := 0; i < count; i++ {
		persons[i] = Person{
			Id:        Format("person_%d", i).String(),
			Name:      Format("Person %d", i).String(),
			BirthDate: "1980-01-15",
			Gender:    "male",
			Phone:     Format("+1-555-123-%04d", i).String(),
			Addresses: []Address{
				{
					Id:      Format("addr_%d_1", i).String(),
					Street:  Format("%d Main St", 100+i).String(),
					City:    "City",
					ZipCode: Format("%05d", 10000+i).String(),
				},
			},
		}
	}
	return persons
}
