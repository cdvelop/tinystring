package jsoncomparison

// Estructuras y datos de prueba para los benchmarks
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

// generateTestData genera datos de prueba para los benchmarks
func generateTestData(count int) []ComplexUser {
	users := make([]ComplexUser, count)
	for i := 0; i < count; i++ {
		users[i] = ComplexUser{
			ID:        "user_" + string(rune(i)),
			Username:  "user_" + string(rune(i)) + "_2024",
			Email:     "user" + string(rune(i)) + "@example.com",
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

// generarDatosInvalidos genera datos JSON inválidos para pruebas de error
func generateInvalidData() []string {
	return []string{
		// JSON mal formado
		`{"id": "user_1", "username": "john_doe", "email": "john@example.com"`,
		// Tipos incorrectos
		`{"id": 123, "username": true, "email": ["not", "an", "email"]}`,
		// Valores nulos inesperados
		`{"id": null, "username": null, "email": null, "profile": null}`,
		// JSON truncado
		`{"id": "user_1", "profile": {"firstName": "John", "lastName":`,
		// Estructura incompleta
		`{}`,
	}
}
