package tinystring

import "testing"

func TestDictionaryBasicFunctionality(t *testing.T) {
	// Test default English
	err := Err(D.Format, D.Invalid).Error()
	expected := T(D.Format, D.Invalid)
	if err != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err)
	}
}

func TestDictionarySpanishTranslation(t *testing.T) {
	// Set Spanish as default
	OutLang(ES)

	err := Err(D.Format, D.Invalid).Error()
	expected := T(ES, D.Format, D.Invalid)
	if err != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err)
	}

	// Reset to English
	OutLang(EN)
}

func TestDictionaryInlineLanguage(t *testing.T) {
	// Use French inline
	err := Err(FR, D.Empty, D.String).Error()
	expected := "vide chaîne"
	if err != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err)
	}
}

func TestDictionaryMixedWithRegularStrings(t *testing.T) {
	// Mix dictionary words with regular strings
	err := Err(D.Invalid, "custom", D.Value).Error()
	expected := "invalid custom value"
	if err != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err)
	}
}

func TestOLFallbackToEnglish(t *testing.T) {
	// Test fallback when translation is empty
	testOL := LocStr{"test", "", "", "", "", "", "", "", ""}
	// Inline get logic for testing
	out := func() string {
		if int(ES) < len(testOL) && testOL[ES] != "" {
			return testOL[ES]
		}
		return testOL[EN] // Fallback to English
	}()
	expected := "test"
	if out != expected {
		t.Errorf("Expected fallback to English 'test', got '%s'", out)
	}
}

func TestLanguageDetection(t *testing.T) {
	c := &conv{}
	// Test that getSystemLang returns a valid language
	lang := c.getSystemLang()
	if lang > ZH {
		t.Errorf("Invalid language detected: %d", lang)
	}
}

func TestComplexErrorComposition(t *testing.T) {
	// Test complex error message composition as per design document
	OutLang(ES)

	// Test: errNegativeUnsigned → D.Negative + D.Numbers + D.Not + D.Supported + D.For + D.Unsigned + D.Integer
	err := Err(D.Negative, D.Numbers, D.Not, D.Supported, D.For, D.Unsigned, D.Integer).Error()

	if len(err) == 0 {
		t.Error("Complex error composition should not be empty")
	}

	t.Logf("Complex error out: %s", err)

	// Reset to English
	OutLang(EN)
}
