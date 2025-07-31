package tinystring

import "testing"

func TestT_LanguageDetection(t *testing.T) {
	// Use real dictionary words: Format (EN: "format", ES: "formato", FR: "format")
	t.Run("lang constant ES", func(t *testing.T) {
		got := T(ES, D.Format)
		if got != "formato" {
			t.Errorf("expected 'formato', got '%s'", got)
		}
	})

	t.Run("lang string ES", func(t *testing.T) {
		got := T("es", D.Format)
		if got != "formato" {
			t.Errorf("expected 'formato', got '%s'", got)
		}
	})

	t.Run("lang constant FR", func(t *testing.T) {
		got := T(FR, D.Format)
		if got != "format" {
			t.Errorf("expected 'format', got '%s'", got)
		}
	})

	t.Run("lang string FR", func(t *testing.T) {
		got := T("FR", D.Format)
		if got != "format" {
			t.Errorf("expected 'format', got '%s'", got)
		}
	})

	t.Run("default lang EN", func(t *testing.T) {
		got := T(D.Format)
		if got != "format" {
			t.Errorf("expected 'format', got '%s'", got)
		}
	})

	// Test phrase composition
	t.Run("phrase ES", func(t *testing.T) {
		got := T("ES", D.Format, D.Invalid)
		if got != "formato inválido" {
			t.Errorf("expected 'formato inválido', got '%s'", got)
		}
	})

	t.Run("phrase EN", func(t *testing.T) {
		got := T(D.Format, D.Invalid)
		if got != "format invalid" {
			t.Errorf("expected 'format invalid', got '%s'", got)
		}
	})
}

func TestTranslationFormatting(t *testing.T) {

	t.Run("no leading space, custom format", func(t *testing.T) {
		// Simula una frase compleja con LocStr y strings, sin espacios extra antes de la primera palabra ni después de los separadores
		// Ejemplo: T(D.Fields, ":", " ", D.Cancel, ")")
		got := T(D.Fields, ":", " ", D.Cancel, ")")
		want := "fields:  cancel)" // Asumiendo EN por defecto
		if got != want {
			t.Errorf("expected '%s', got '%s'", want, got)
		}
	})

	t.Run("no space before colon, phrase with punctuation", func(t *testing.T) {
		// Simula formato con puntuación pegada
		got := T(D.Format, ":", D.Invalid)
		want := "format: invalid"
		if got != want {
			t.Errorf("expected '%s', got '%s'", want, got)
		}
	})

	t.Run("newline with translated field alignment", func(t *testing.T) {
		// Reproduce el caso del shortcuts.go donde después de newline viene D.Fields
		// y debe quedar alineado sin espacio extra
		got := T("Tabs:\n", D.Fields, ":")
		want := "Tabs:\nfields:"
		if got != want {
			t.Errorf("expected '%s', got '%s'", want, got)
		}
	})
}
