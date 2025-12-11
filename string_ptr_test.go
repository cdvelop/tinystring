package fmt

import (
	std_fmt "fmt"
	"testing"
)

func TestStringPointer(t *testing.T) {
	tests := []struct {
		name          string
		initialValue  string
		transform     func(*Conv) *Conv
		expectedValue string
	}{
		{
			name:         "Remove tildes from string pointer",
			initialValue: "áéíóúÁÉÍÓÚ",
			transform: func(t *Conv) *Conv {
				return t.Tilde()
			},
			expectedValue: "aeiouAEIOU",
		},
		{
			name:         "Convert to lowercase with string pointer",
			initialValue: "HELLO WORLD",
			transform: func(t *Conv) *Conv {
				return t.ToLower()
			},
			expectedValue: "hello world",
		},
		{
			name:         "Convert to camelCase with string pointer",
			initialValue: "hello world example",
			transform: func(t *Conv) *Conv {
				return t.CamelLow()
			},
			expectedValue: "helloWorldExample",
		},
		{
			name:         "Multiple transforms with string pointer",
			initialValue: "Él Múrcielago Rápido",
			transform: func(t *Conv) *Conv {
				return t.Tilde().CamelLow()
			},
			expectedValue: "elMurcielagoRapido",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create string pointer with initial value
			originalPtr := tt.initialValue

			// Convert using string pointer and apply changes
			tt.transform(Convert(&originalPtr)).Apply()

			// Check if original pointer was updated correctly
			if originalPtr != tt.expectedValue {
				t.Errorf("\noriginalPtr = %q\nwant %q", originalPtr, tt.expectedValue)
			}
		})
	}
}

// Estos ejemplos ilustran cómo usar los punteros a strings para evitar asignaciones adicionales
func Example_stringPointerBasic() {
	// Creamos una variable string que queremos modificar
	myText := "héllô wórld"

	// En lugar de crear una nueva variable con el resultado,
	// modificamos directamente la variable original usando Apply()
	Convert(&myText).Tilde().ToLower().Apply()

	// La variable original ha sido modificada
	std_fmt.Println(myText)
	// Output: hello world
}

func Example_stringPointerCamelCase() {
	// Ejemplo de uso con múltiples transformaciones
	originalText := "Él Múrcielago Rápido"

	// Las transformaciones modifican la variable original directamente
	// usando el método Apply() para actualizar el puntero
	Convert(&originalText).Tilde().CamelLow().Apply()

	std_fmt.Println(originalText)
	// Output: elMurcielagoRapido
}

func Example_stringPointerEfficiency() {
	// En aplicaciones de alto rendimiento, reducir asignaciones de memoria
	// puede ser importante para evitar la presión sobre el garbage collector
	// Método tradicional (crea nuevas asignaciones de memoria)
	traditionalText := "Texto con ACENTOS"
	processedText := Convert(traditionalText).Tilde().ToLower().String()
	std_fmt.Println(processedText)

	// Método con punteros (modifica directamente la variable original)
	directText := "Otro TEXTO con ACENTOS"
	Convert(&directText).Tilde().ToLower().Apply()
	std_fmt.Println(directText)

	// Output:
	// texto con acentos
	// otro texto con acentos
}
