package tinystring_test

import (
	"testing"

	. "github.com/cdvelop/tinystring"
)

func TestCountOccurrences(t *testing.T) {
	var testCases = map[string]struct {
		conv     string
		search   string
		expected int
	}{
		"Caso1": {
			conv:     "Hola, mundo!",
			search:   "mundo",
			expected: 1,
		},
		"Caso2": {
			conv:     "Hola, mundo!",
			search:   "golang",
			expected: 0,
		},
		"Caso3": {
			conv:     "Hola, mundo!",
			search:   "",
			expected: 0,
		},
		"Caso4": {
			conv:     "Hola",
			search:   "Hola, mundo!",
			expected: 0,
		},
		"Caso5": {
			conv:     "abracadabra",
			search:   "abra",
			expected: 2,
		},
		"Caso6": {
			conv:     "abracadabra",
			search:   "bra",
			expected: 2,
		},
		"Caso7": {
			conv:     "abra,cadabra",
			search:   ",",
			expected: 1,
		},
		"Caso8": {
			conv:     "(abraLol,*?¡¡",
			search:   "Lol",
			expected: 1,
		},
		"Caso9 ": {
			conv:     "(abraLol,*?¡¡",
			search:   "LoL",
			expected: 0,
		},
		"Caso10 ": {
			conv:     "(¡ab¡raLol,*?¡¡",
			search:   "¡",
			expected: 4,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := CountOccurrences(tc.conv, tc.search)
			if result != tc.expected {
				t.Errorf("Error: Se esperaba %v, pero se obtuvo %v. Texto: %s, Búsqueda: %s", tc.expected, result, tc.conv, tc.search)
			}
		})
	}
}

func TestContains(t *testing.T) {
	var testCases = map[string]struct {
		conv     string
		search   string
		expected bool
	}{
		"Encontrado": {
			conv:     "Hola, mundo!",
			search:   "mundo",
			expected: true,
		},
		"No encontrado": {
			conv:     "Hola, mundo!",
			search:   "golang",
			expected: false,
		},
		"Búsqueda vacía": {
			conv:     "Hola, mundo!",
			search:   "",
			expected: false,
		},
		"Texto más corto que búsqueda": {
			conv:     "Hola",
			search:   "Hola, mundo!",
			expected: false,
		},
		"Múltiples ocurrencias": {
			conv:     "abracadabra",
			search:   "abra",
			expected: true,
		},
		"Sensible a mayúsculas": {
			conv:     "(abraLol,*?¡¡",
			search:   "LoL",
			expected: false,
		},
		"Búsqueda de caracteres especiales": {
			conv:     "(¡ab¡raLol,*?¡¡",
			search:   "¡",
			expected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := Contains(tc.conv, tc.search)
			if result != tc.expected {
				t.Errorf("Error: Se esperaba %v, pero se obtuvo %v. Texto: %s, Búsqueda: %s", tc.expected, result, tc.conv, tc.search)
			}
		})
	}
}
