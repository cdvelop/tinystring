package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Complex data structures for realistic JSON benchmarking
type Metadata struct {
	Source      string   `json:"source"`
	Campaign    string   `json:"campaign"`
	Referrer    string   `json:"referrer"`
	Experiments []string `json:"experiments"`
	Score       float64  `json:"score"`
}

type CustomFields struct {
	EmployeeID string `json:"employee_id"`
	Department string `json:"department"`
	Team       string `json:"team"`
}

type Features struct {
	BetaFeatures   bool `json:"beta_features"`
	Analytics      bool `json:"analytics"`
	AdvancedSearch bool `json:"advanced_search"`
}

type ComplexUser struct {
	ID          string         `json:"id"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	CreatedAt   string         `json:"created_at"`
	LastLogin   string         `json:"last_login,omitempty"`
	IsActive    bool           `json:"is_active"`
	Profile     ComplexProfile `json:"profile"`
	Permissions []string       `json:"permissions"`
	Metadata    Metadata       `json:"metadata"`
	Stats       ComplexStats   `json:"stats"`
}

type ComplexProfile struct {
	FirstName    string               `json:"first_name"`
	LastName     string               `json:"last_name"`
	DisplayName  string               `json:"display_name"`
	Bio          string               `json:"bio"`
	AvatarURL    string               `json:"avatar_url"`
	BirthDate    string               `json:"birth_date,omitempty"`
	PhoneNumbers []ComplexPhoneNumber `json:"phone_numbers"`
	Addresses    []ComplexAddress     `json:"addresses"`
	SocialLinks  []ComplexSocialLink  `json:"social_links"`
	Preferences  ComplexPreferences   `json:"preferences"`
	CustomFields CustomFields         `json:"custom_fields"`
}

type ComplexAddress struct {
	ID          string              `json:"id"`
	Type        string              `json:"type"` // home, work, billing, etc.
	Street      string              `json:"street"`
	Street2     string              `json:"street2,omitempty"`
	City        string              `json:"city"`
	State       string              `json:"state"`
	Country     string              `json:"country"`
	PostalCode  string              `json:"postal_code"`
	Coordinates *ComplexCoordinates `json:"coordinates,omitempty"`
	IsPrimary   bool                `json:"is_primary"`
	IsVerified  bool                `json:"is_verified"`
}

type ComplexCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  int     `json:"accuracy"`
}

type ComplexPhoneNumber struct {
	ID         string `json:"id"`
	Type       string `json:"type"` // mobile, home, work, fax
	Number     string `json:"number"`
	Extension  string `json:"extension,omitempty"`
	IsPrimary  bool   `json:"is_primary"`
	IsVerified bool   `json:"is_verified"`
}

type ComplexSocialLink struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Verified bool   `json:"verified"`
}

type ComplexPreferences struct {
	Language      string                   `json:"language"`
	Timezone      string                   `json:"timezone"`
	Theme         string                   `json:"theme"`
	Currency      string                   `json:"currency"`
	DateFormat    string                   `json:"date_format"`
	TimeFormat    string                   `json:"time_format"`
	Notifications ComplexNotificationPrefs `json:"notifications"`
	Privacy       ComplexPrivacySettings   `json:"privacy"`
	Features      Features                 `json:"features"`
}

type ComplexNotificationPrefs struct {
	Email     bool `json:"email"`
	SMS       bool `json:"sms"`
	Push      bool `json:"push"`
	InApp     bool `json:"in_app"`
	Marketing bool `json:"marketing"`
}

type ComplexPrivacySettings struct {
	ProfileVisibility string   `json:"profile_visibility"` // public, friends, private
	ShowEmail         bool     `json:"show_email"`
	ShowPhone         bool     `json:"show_phone"`
	AllowMessaging    bool     `json:"allow_messaging"`
	BlockedUsers      []string `json:"blocked_users"`
}

type ComplexStats struct {
	LoginCount       int64  `json:"login_count"`
	LastActivity     string `json:"last_activity"`
	SessionDuration  int64  `json:"session_duration_seconds"`
	PageViews        int64  `json:"page_views"`
	ActionsCount     int64  `json:"actions_count"`
	SubscriptionTier string `json:"subscription_tier"`
	StorageUsed      int64  `json:"storage_used_bytes"`
	BandwidthUsed    int64  `json:"bandwidth_used_bytes"`
}

func main() {
	// Realistic complex operations using standard library (multiple separate calls)
	conv := "MÍ téxtO cön AcÉntos Y MÁS TEXTO"

	// Complex transformations using multiple standard library calls
	step1 := strings.ToLower(conv)
	step2 := removeTildes(step1)
	step3 := strings.ReplaceAll(step2, " ", "_")
	step4 := toCamelCase(step3)
	processed := capitalize(step4)

	// Number processing with multiple operations
	prices := []float64{1234.567, 9876.54, 42.0}
	formattedPrices := make([]string, len(prices))
	for i, price := range prices {
		rounded := roundFloat(price, 2)
		formattedPrices[i] = formatNumber(rounded)
	}

	// Complex string operations with multiple calls
	userInput := "  Hello@World#2024!  "
	trimmed := strings.TrimSpace(userInput)
	replaced1 := strings.ReplaceAll(trimmed, "@", "_at_")
	replaced2 := strings.ReplaceAll(replaced1, "#", "_hash_")
	replaced3 := strings.ReplaceAll(replaced2, "!", "")
	cleaned := strings.ToLower(replaced3)

	// Manual joining and formatting
	priceList := strings.Join(formattedPrices, ", ")
	finalResult := fmt.Sprintf(
		"Processed: %s | Cleaned: %s | Prices: %s",
		processed, cleaned, priceList,
	)

	// Additional complex transformations
	mixedText := "José María-González_2024"
	normalized1 := removeTildes(mixedText)
	normalized2 := strings.ReplaceAll(normalized1, "-", "_")
	normalized := toSnakeCase(normalized2)
	// Final comprehensive result
	result := fmt.Sprintf("%s | Normalized: %s", finalResult, normalized)

	// Complex JSON operations using standard library
	complexJsonOperations()

	_ = result
}

// Helper functions to simulate equivalent operations
func removeTildes(s string) string {
	replacements := map[rune]rune{
		'á': 'a', 'é': 'e', 'í': 'i', 'ó': 'o', 'ú': 'u',
		'Á': 'A', 'É': 'E', 'Í': 'I', 'Ó': 'O', 'Ú': 'U',
		'ñ': 'n', 'Ñ': 'N', 'ü': 'u', 'Ü': 'U',
		'ç': 'c', 'Ç': 'C',
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, ok := replacements[r]; ok {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		return s
	}

	var result strings.Builder
	result.WriteString(words[0])
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result.WriteString(strings.ToUpper(words[i][:1]))
			if len(words[i]) > 1 {
				result.WriteString(words[i][1:])
			}
		}
	}
	return result.String()
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func roundFloat(val float64, precision int) float64 {
	ratio := 1.0
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}

func formatNumber(val float64) string {
	str := strconv.FormatFloat(val, 'f', 2, 64)
	parts := strings.Split(str, ".")

	// Add thousand separators
	intPart := parts[0]
	var result strings.Builder
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(digit)
	}

	if len(parts) > 1 {
		result.WriteString(".")
		result.WriteString(parts[1])
	}

	return result.String()
}

func toSnakeCase(s string) string {
	// Convert to snake_case
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)
	return strings.ToLower(snake)
}

// Complex JSON operations using standard library
func complexJsonOperations() {
	// Create complex nested data structures
	users := createComplexUserData()

	// Multiple JSON encoding operations
	for _, user := range users {
		// Encode user profile
		userJson, _ := json.Marshal(user)

		// Encode just addresses
		addressJson, _ := json.Marshal(user.Profile.Addresses)

		// Encode preferences
		prefsJson, _ := json.Marshal(user.Profile.Preferences)

		// Decode operations
		var decodedUser ComplexUser
		json.Unmarshal(userJson, &decodedUser)

		var decodedAddresses []ComplexAddress
		json.Unmarshal(addressJson, &decodedAddresses)

		// Mixed operations - encode then decode with modifications
		temp := user.Profile.Preferences
		temp.Theme = "dark"
		temp.Notifications.Email = !temp.Notifications.Email

		modifiedJson, _ := json.Marshal(temp)
		var newPrefs ComplexPreferences
		json.Unmarshal(modifiedJson, &newPrefs)

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
	// Usar fechas fijas para evitar dependencia del paquete time
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
			}, Stats: ComplexStats{
				LoginCount: 1247, LastActivity: "2024-06-12T08:00:00Z",
				SessionDuration: 3600, PageViews: 15643, ActionsCount: 892,
				SubscriptionTier: "premium", StorageUsed: 2147483648, BandwidthUsed: 10737418240,
			},
		},
		{
			ID:       "user_002",
			Username: "jane_smith_pro",
			Email:    "jane.smith@company.com", CreatedAt: "2023-12-12T10:00:00Z",
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
			}, Stats: ComplexStats{
				LoginCount: 892, LastActivity: "2024-06-12T09:30:00Z",
				SessionDuration: 2700, PageViews: 8932, ActionsCount: 456,
				SubscriptionTier: "enterprise", StorageUsed: 5368709120, BandwidthUsed: 21474836480,
			},
		},
	}
}
