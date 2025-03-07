package tinystring_test

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestCountOccurrences(t *testing.T) {
	var testCases = map[string]struct {
		text     string
		search   string
		expected int
	}{
		"Caso1": {
			text:     "Hola, mundo!",
			search:   "mundo",
			expected: 1,
		},
		"Caso2": {
			text:     "Hola, mundo!",
			search:   "golang",
			expected: 0,
		},
		"Caso3": {
			text:     "Hola, mundo!",
			search:   "",
			expected: 0,
		},
		"Caso4": {
			text:     "Hola",
			search:   "Hola, mundo!",
			expected: 0,
		},
		"Caso5": {
			text:     "abracadabra",
			search:   "abra",
			expected: 2,
		},
		"Caso6": {
			text:     "abracadabra",
			search:   "bra",
			expected: 2,
		},
		"Caso7": {
			text:     "abra,cadabra",
			search:   ",",
			expected: 1,
		},
		"Caso8": {
			text:     "(abraLol,*?¡¡",
			search:   "Lol",
			expected: 1,
		},
		"Caso9 ": {
			text:     "(abraLol,*?¡¡",
			search:   "LoL",
			expected: 0,
		},
		"Caso10 ": {
			text:     "(¡ab¡raLol,*?¡¡",
			search:   "¡",
			expected: 4,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := CountOccurrences(tc.text, tc.search)
			if result != tc.expected {
				t.Errorf("Error: Se esperaba %v, pero se obtuvo %v. Texto: %s, Búsqueda: %s", tc.expected, result, tc.text, tc.search)
			}
		})
	}
}

func TestContains(t *testing.T) {
	var testCases = map[string]struct {
		text     string
		search   string
		expected bool
	}{
		"Encontrado": {
			text:     "Hola, mundo!",
			search:   "mundo",
			expected: true,
		},
		"No encontrado": {
			text:     "Hola, mundo!",
			search:   "golang",
			expected: false,
		},
		"Búsqueda vacía": {
			text:     "Hola, mundo!",
			search:   "",
			expected: false,
		},
		"Texto más corto que búsqueda": {
			text:     "Hola",
			search:   "Hola, mundo!",
			expected: false,
		},
		"Múltiples ocurrencias": {
			text:     "abracadabra",
			search:   "abra",
			expected: true,
		},
		"Sensible a mayúsculas": {
			text:     "(abraLol,*?¡¡",
			search:   "LoL",
			expected: false,
		},
		"Búsqueda de caracteres especiales": {
			text:     "(¡ab¡raLol,*?¡¡",
			search:   "¡",
			expected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := Contains(tc.text, tc.search)
			if result != tc.expected {
				t.Errorf("Error: Se esperaba %v, pero se obtuvo %v. Texto: %s, Búsqueda: %s", tc.expected, result, tc.text, tc.search)
			}
		})
	}
}
