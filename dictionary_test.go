package tinystring

import (
	"reflect"
	"testing"
)

func TestDictionaryBasicFunctionality(t *testing.T) {
	// Test default English
	err := Err(D.Format, D.Invalid).Error()
	expected := T(D.Format, D.Invalid).String()
	if err != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err)
	}
}

func TestDictionarySpanishTranslation(t *testing.T) {
	// Set Spanish as default
	OutLang(ES)

	err := Err(D.Format, D.Invalid).Error()
	expected := T(ES, D.Format, D.Invalid).String()
	if err != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err)
	}

	// Reset to English
	OutLang(EN)
}

func TestDictionaryInlineLanguage(t *testing.T) {
	// Use French inline
	err := Err(FR, D.Empty, D.String).Error()
	expected := "Vide Chaîne"
	if err != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err)
	}
}

func TestDictionaryMixedWithRegularStrings(t *testing.T) {
	// Mix dictionary words with regular strings
	err := Err(D.Invalid, "custom", D.Value).Error()
	expected := "Invalid custom Value"
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

	// Test: errNegativeUnsigned → D.Negative + D.Numbers + D.Not + D.Supported + D.To + D.Unsigned + D.Integer
	err := Err(D.Numbers, D.Negative, D.Not, D.Supported).Error()

	if len(err) == 0 {
		t.Error("Complex error composition should not be empty")
	}

	t.Logf("Complex error out: %s", err)

	// Reset to English
	OutLang(EN)
}

func TestDictionaryConsistency(t *testing.T) {
	typeOfD := reflect.TypeOf(D)
	valueOfD := reflect.ValueOf(D)
	for i := 0; i < typeOfD.NumField(); i++ {
		field := typeOfD.Field(i)
		fieldName := field.Name
		fieldValue := valueOfD.Field(i).Interface().(LocStr)
		lowerFieldName := Convert(fieldName).Low().String()
		eng := Convert(fieldValue[EN]).Low().String()
		// Tomar los 2 primeros y 2 últimos caracteres
		fnLen := len(lowerFieldName)
		engLen := len(eng)
		if fnLen >= 2 && engLen >= 2 {
			fnFirst := lowerFieldName[:2]
			fnLast := lowerFieldName[fnLen-2:]
			engFirst := eng[:2]
			engLast := eng[engLen-2:]
			if fnFirst != engFirst || fnLast != engLast {
				t.Fatalf("Field '%s' value '%s' does not match first/last 2 chars of field name '%s'", fieldName, eng, lowerFieldName)
			}
		}
	}
}
