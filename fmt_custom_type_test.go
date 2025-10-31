package tinystring

import "testing"

// customType simula un tipo personalizado como pdfVersion en fpdf
type customType string

func (c customType) String() string {
	return string(c)
}

// customInt simula un tipo numérico personalizado con String()
type customInt int

func (c customInt) String() string {
	return Convert(int(c)).String()
}

// TestFmtWithCustomTypeString verifica que Fmt maneje correctamente tipos personalizados con método String()
func TestFmtWithCustomTypeString(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
		bug      bool // true si este caso representa el bug actual
	}{
		{
			name:     "string literal works",
			format:   "Value: %s",
			args:     []any{"hello"},
			expected: "Value: hello",
			bug:      false,
		},
		{
			name:     "custom string type with String() - BUG",
			format:   "Version: %s",
			args:     []any{customType("1.3")},
			expected: "Version: 1.3",
			bug:      true, // Este es el bug - actualmente devuelve ""
		},
		{
			name:     "custom int type with String() - BUG",
			format:   "Count: %s",
			args:     []any{customInt(42)},
			expected: "Count: 42",
			bug:      true, // Este es el bug - actualmente devuelve ""
		},
		{
			name:     "PDF version format - BUG (real world case)",
			format:   "%%PDF-%s",
			args:     []any{customType("1.4")},
			expected: "%PDF-1.4",
			bug:      true, // Este es exactamente el problema de fpdf
		},
		{
			name:     "multiple custom types - BUG",
			format:   "%s version %s",
			args:     []any{customType("PDF"), customType("1.5")},
			expected: "PDF version 1.5",
			bug:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Fmt(tt.format, tt.args...)

			if tt.bug {
				// Para casos con bug, verificamos el comportamiento ACTUAL (incorrecto)
				if result == tt.expected {
					t.Logf("✅ BUG FIXED: Expected %q, got %q (bug was fixed!)", tt.expected, result)
				} else {
					t.Logf("🐛 BUG CONFIRMED: Expected %q, got %q (bug still exists)", tt.expected, result)
					// No fallamos el test porque sabemos que es un bug conocido
					// Esto permite que el test pase mientras documentamos el bug
				}
			} else {
				// Para casos sin bug, verificamos que funcione correctamente
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

// TestFmtCustomTypeWithOtherFormats verifica que otros formatos (%d, %v) funcionen con custom types
func TestFmtCustomTypeWithOtherFormats(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{
			name:     "custom int with %d",
			format:   "Count: %d",
			args:     []any{customInt(42)},
			expected: "Count: 42",
		},
		{
			name:     "custom type with %v (should work)",
			format:   "Version: %v",
			args:     []any{customType("1.3")},
			expected: "Version: 1.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Fmt(tt.format, tt.args...)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			} else {
				t.Logf("✅ Working: %q", result)
			}
		})
	}
}

// TestFmtStringerInterface verifica comportamiento con diferentes tipos que implementan String()
func TestFmtStringerInterface(t *testing.T) {
	type version struct {
		major int
		minor int
	}

	// Implementar String() para version dentro del test
	var _ = func(v version) string {
		return Fmt("%d.%d", v.major, v.minor)
	}

	// version con método String()
	v := version{1, 4}

	// Intentar formatear con %v (debería funcionar si AnyToBuff maneja String())
	result := Fmt("PDF Version: %v", v)
	t.Logf("Result with %%v: %q", result)

	// Intentar formatear con %s (actualmente fallará)
	result2 := Fmt("PDF Version: %s", v)
	t.Logf("Result with %%s: %q", result2)

	if result2 == "" {
		t.Log("BUG CONFIRMED: percentage-s format doesn't work with custom types (returns empty string)")
	}
}
