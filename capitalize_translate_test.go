package tinystring

import (
	"testing"
)

func TestCapitalizeWithMultilineTranslation(t *testing.T) {
	tests := []struct {
		name        string
		appName     string
		lang        string
		expected    string
		description string
	}{
		{
			name:        "Simple multiline with Capitalize",
			appName:     "TestApp",
			lang:        "EN",
			expected:    "Testapp Shortcuts Keyboard (\"en\"):\n\nTabs:\n  • Tab/Shift+Tab  - Switch Tabs\n\nFields :\n  • Left/Right     - Navigate Fields\n  • Enter          - Edit/Execute\n  • Esc            - Cancel \n\nLanguage Supported : En, Es, Zh, Hi, Ar, Pt, Fr, De, Ru",
			description: "Test that Capitalize preserves multiline structure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the generateHelpContent method but simplified
			result := generateSimplifiedHelpContent(tt.appName, tt.lang)

			// Check if the result maintains proper formatting
			if result != tt.expected {
				t.Errorf("Test %s failed.\nExpected: %q\nGot:      %q", tt.name, tt.expected, result)
			}
		})
	}
}

// generateSimplifiedHelpContent simulates the PROBLEMATIC method WITH Capitalize
func generateSimplifiedHelpContent(appName, lang string) string {
	// Test the core issue: preserving spaces in mixed translated/non-translated content
	return T(appName, D.Shortcuts, D.Keyboard, "(\""+lang+"\"):\n\nTabs:\n  • Tab/Shift+Tab  -", D.Switch, " tabs\n\n", D.Fields, ":\n  • Left/Right     - Navigate fields\n  • Enter          - Edit/Execute\n  • Esc            - ", D.Cancel, " \n\n", D.Language, D.Supported, ": EN, ES, ZH, HI, AR, PT, FR, DE, RU").Capitalize().String()
}
