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
