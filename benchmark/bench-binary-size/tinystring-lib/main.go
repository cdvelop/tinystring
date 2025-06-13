package main

import "github.com/cdvelop/tinystring"

// Complex data structures for realistic JSON benchmarking
type Metadata struct {
	Source      string
	Campaign    string
	Referrer    string
	Experiments []string
	Score       float64
}

type CustomFields struct {
	EmployeeID string
	Department string
	Team       string
}

type Features struct {
	BetaFeatures   bool
	Analytics      bool
	AdvancedSearch bool
}

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

type ComplexCoordinates struct {
	Latitude  float64
	Longitude float64
	Accuracy  int
}

type ComplexPhoneNumber struct {
	ID         string
	Type       string
	Number     string
	Extension  string
	IsPrimary  bool
	IsVerified bool
}

type ComplexSocialLink struct {
	Platform string
	URL      string
	Username string
	Verified bool
}

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

type ComplexNotificationPrefs struct {
	Email     bool
	SMS       bool
	Push      bool
	InApp     bool
	Marketing bool
}

type ComplexPrivacySettings struct {
	ProfileVisibility string
	ShowEmail         bool
	ShowPhone         bool
	AllowMessaging    bool
	BlockedUsers      []string
}

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

func main() {
	// Realistic complex operations using TinyString's chaining capabilities
	conv := "MÍ téxtO cön AcÉntos Y MÁS TEXTO"
	// Complex chained transformations (TinyString's strength)
	processed := tinystring.Convert(conv).
		ToLower().
		RemoveTilde().
		Replace(" ", "_").
		CamelCaseLower().
		Capitalize().
		String()

	// Number processing with chaining
	prices := []any{1234.567, 9876.54, 42.0}
	formattedPrices := make([]string, len(prices))
	for i, price := range prices {
		formattedPrices[i] = tinystring.Convert(price).
			RoundDecimals(2).
			FormatNumber().
			String()
	}

	// Complex string operations with chaining
	userInput := "  Hello@World#2024!  "
	cleaned := tinystring.Convert(userInput).
		Trim().
		Replace("@", "_at_").
		Replace("#", "_hash_").
		Replace("!", "").
		ToLower().
		String()

	// Advanced formatting and joining
	priceList := tinystring.Convert(formattedPrices).Join(", ")
	finalResult := tinystring.Format(
		"Processed: %s | Cleaned: %s | Prices: %s",
		processed, cleaned, priceList,
	)
	// Additional complex transformations
	mixedText := "José María-González_2024"
	normalized := tinystring.Convert(mixedText).
		RemoveTilde().
		Replace("-", "_").
		ToSnakeCaseLower().
		String()
	// Final comprehensive result
	result := tinystring.Format("%s | Normalized: %s", finalResult, normalized)

	// Complex JSON operations using TinyString
	complexJsonOperations()

	_ = result
}

// Complex JSON operations using TinyString
func complexJsonOperations() {
	// Create complex nested data structures
	users := createComplexUserData()

	// Multiple JSON encoding operations
	for _, user := range users {
		// Encode user profile
		userJson, _ := tinystring.Convert(&user).JsonEncode()

		// Encode just addresses
		addressJson, _ := tinystring.Convert(&user.Profile.Addresses).JsonEncode()

		// Encode preferences
		prefsJson, _ := tinystring.Convert(&user.Profile.Preferences).JsonEncode()

		// Decode operations
		var decodedUser ComplexUser
		tinystring.Convert(userJson).JsonDecode(&decodedUser)

		var decodedAddresses []ComplexAddress
		tinystring.Convert(addressJson).JsonDecode(&decodedAddresses)

		// Mixed operations - encode then decode with modifications
		temp := user.Profile.Preferences
		temp.Theme = "dark"
		temp.Notifications.Email = !temp.Notifications.Email

		modifiedJson, _ := tinystring.Convert(&temp).JsonEncode()
		var newPrefs ComplexPreferences
		tinystring.Convert(modifiedJson).JsonDecode(&newPrefs)

		// Suppress unused variable warnings
		_ = userJson
		_ = addressJson
		_ = prefsJson
		_ = decodedUser
		_ = decodedAddresses
		_ = newPrefs
	}
}

func createComplexUserData() []ComplexUser {
	return []ComplexUser{
		{
			ID:        "user_001",
			Username:  "john_doe_2024",
			Email:     "john.doe@example.com",
			CreatedAt: "2024-06-12T10:00:00Z",
			LastLogin: "2024-06-05T10:00:00Z",
			IsActive:  true,
			Profile: ComplexProfile{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "Johnny D",
				Bio:         "Software engineer passionate about technology and innovation",
				AvatarURL:   "https://cdn.example.com/avatars/john_doe.jpg",
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
					{
						ID: "addr_002", Type: "work", Street: "456 Business Ave", Street2: "Suite 100",
						City: "Tech City", State: "CA", Country: "USA", PostalCode: "54321", IsPrimary: false, IsVerified: true,
					},
				},
				SocialLinks: []ComplexSocialLink{
					{Platform: "twitter", URL: "https://twitter.com/johndoe", Username: "@johndoe", Verified: false},
					{Platform: "linkedin", URL: "https://linkedin.com/in/johndoe", Username: "johndoe", Verified: true},
					{Platform: "github", URL: "https://github.com/johndoe", Username: "johndoe", Verified: true},
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
			Permissions: []string{"read", "write", "admin", "analytics"},
			Metadata: Metadata{
				Source: "web_signup", Campaign: "summer_2024", Referrer: "google",
				Experiments: []string{"new_ui", "faster_search"}, Score: 85.7,
			},
			Stats: ComplexStats{
				LoginCount: 1247, LastActivity: "2024-06-12T08:00:00Z",
				SessionDuration: 3600, PageViews: 15643, ActionsCount: 892,
				SubscriptionTier: "premium", StorageUsed: 2147483648, BandwidthUsed: 10737418240,
			},
		},
		{
			ID:        "user_002",
			Username:  "jane_smith_pro",
			Email:     "jane.smith@company.com",
			CreatedAt: "2023-12-12T10:00:00Z",
			LastLogin: "2024-06-12T10:00:00Z",
			IsActive:  true,
			Profile: ComplexProfile{
				FirstName:   "Jane",
				LastName:    "Smith",
				DisplayName: "Jane S.",
				Bio:         "Product manager focused on user experience and data-driven decisions",
				AvatarURL:   "https://cdn.example.com/avatars/jane_smith.jpg",
				PhoneNumbers: []ComplexPhoneNumber{
					{ID: "ph_003", Type: "mobile", Number: "+1-555-555-1234", IsPrimary: true, IsVerified: true},
				},
				Addresses: []ComplexAddress{
					{
						ID: "addr_003", Type: "home", Street: "789 Oak Drive", City: "Somewhere",
						State: "NY", Country: "USA", PostalCode: "67890", IsPrimary: true, IsVerified: true,
						Coordinates: &ComplexCoordinates{Latitude: 40.7128, Longitude: -74.0060, Accuracy: 5},
					},
				},
				SocialLinks: []ComplexSocialLink{
					{Platform: "linkedin", URL: "https://linkedin.com/in/janesmith", Username: "janesmith", Verified: true},
				},
				Preferences: ComplexPreferences{
					Language: "en-US", Timezone: "America/New_York", Theme: "dark", Currency: "USD",
					DateFormat: "DD/MM/YYYY", TimeFormat: "24h",
					Notifications: ComplexNotificationPrefs{
						Email: true, SMS: true, Push: false, InApp: true, Marketing: true,
					},
					Privacy: ComplexPrivacySettings{
						ProfileVisibility: "public", ShowEmail: true, ShowPhone: false,
						AllowMessaging: true, BlockedUsers: []string{"user_spam_001"},
					},
					Features: Features{
						BetaFeatures: false, Analytics: true, AdvancedSearch: true,
					},
				},
				CustomFields: CustomFields{
					EmployeeID: "EMP002", Department: "Product", Team: "UX",
				},
			},
			Permissions: []string{"read", "write", "moderate"},
			Metadata: Metadata{
				Source: "invite", Campaign: "referral_program", Referrer: "user_001",
				Experiments: []string{"new_dashboard"}, Score: 92.3,
			},
			Stats: ComplexStats{
				LoginCount: 892, LastActivity: "2024-06-12T09:30:00Z",
				SessionDuration: 2700, PageViews: 8932, ActionsCount: 456,
				SubscriptionTier: "enterprise", StorageUsed: 5368709120, BandwidthUsed: 21474836480,
			},
		},
	}
}
