package tinystring

import "testing"

func TestT_LanguageDetection(t *testing.T) {
	// Use real dictionary words: Format (EN: "format", ES: "formato", FR: "format")
	t.Run("lang constant ES", func(t *testing.T) {
		got := Translate(ES, D.Format).String()
		if got != "Formato" {
			t.Errorf("expected 'Formato', got '%s'", got)
		}
	})

	t.Run("lang string ES", func(t *testing.T) {
		got := Translate("es", D.Format).String()
		if got != "Formato" {
			t.Errorf("expected 'Formato', got '%s'", got)
		}
	})

	t.Run("lang constant FR", func(t *testing.T) {
		got := Translate(FR, D.Format).String()
		if got != "Format" {
			t.Errorf("expected 'Format', got '%s'", got)
		}
	})

	t.Run("lang string FR", func(t *testing.T) {
		got := Translate("FR", D.Format).String()
		if got != "Format" {
			t.Errorf("expected 'Format', got '%s'", got)
		}
	})

	t.Run("default lang EN", func(t *testing.T) {
		got := Translate(D.Format).String()
		if got != "Format" {
			t.Errorf("expected 'Format', got '%s'", got)
		}
	})

	// Test phrase composition
	t.Run("phrase ES", func(t *testing.T) {
		got := Translate("ES", D.Format, D.Invalid).String()
		if got != "Formato Inválido" {
			t.Errorf("expected 'Formato Inválido', got '%s'", got)
		}
	})

	t.Run("phrase EN", func(t *testing.T) {
		got := Translate(D.Format, D.Invalid).String()
		if got != "Format Invalid" {
			t.Errorf("expected 'Format Invalid', got '%s'", got)
		}
	})
}

func TestTranslationFormatting(t *testing.T) {

	t.Run("no leading space, custom format", func(t *testing.T) {
		// Simula una frase compleja con LocStr y strings, sin espacios extra antes de la primera palabra ni después de los separadores
		// Ejemplo: Translate(D.Fields, ":", D.Cancel, ")")
		got := Translate(D.Fields, ":", D.Cancel, ")").String()
		want := "Fields: Cancel)" // Asumiendo EN por defecto
		if got != want {
			t.Errorf("expected '%s', got '%s'", want, got)
		}
	})

	t.Run("no space before colon, phrase with punctuation", func(t *testing.T) {
		// Simula formato con puntuación pegada
		got := Translate(D.Format, ":", D.Invalid).String()
		want := "Format: Invalid"
		if got != want {
			t.Errorf("expected '%s', got '%s'", want, got)
		}
	})

	t.Run("newline with translated field alignment", func(t *testing.T) {
		// Reproduce el caso del shortcuts.go donde después de newline viene D.Fields
		// y debe quedar alineado sin espacio extra
		got := Translate("Tabs:\n", D.Fields, ":").String()
		want := "Tabs:\nFields:"
		if got != want {
			t.Errorf("expected '%s', got '%s'", want, got)
		}
	})
}

func BenchmarkTranslate(b *testing.B) {
	b.ReportAllocs()
	b.Run("Simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := Translate(D.Format)
			c.PutConv()
		}
	})
	b.Run("WithLang", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := Translate(ES, D.Format)
			c.PutConv()
		}
	})
	b.Run("Complex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := Translate(ES, D.Format, ":", D.Invalid)
			c.PutConv()
		}
	})
}

func BenchmarkTranslatePointer(b *testing.B) {
	b.ReportAllocs()
	b.Run("SimplePtr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := Translate(&D.Format)
			c.PutConv()
		}
	})
	b.Run("WithLangPtr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := Translate(ES, &D.Format)
			c.PutConv()
		}
	})
	b.Run("ComplexPtr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c := Translate(ES, &D.Format, ":", &D.Invalid)
			c.PutConv()
		}
	})
}
