//go:build wasm

package tinystring

import (
	"strings"
	"syscall/js"
)

// getSystemLang detects browser language from navigator.language
func getSystemLang() lang {
	// Get browser language
	navigator := js.Global().Get("navigator")
	if navigator.IsUndefined() {
		return EN
	}

	language := navigator.Get("language")
	if language.IsUndefined() {
		return EN
	}

	langCode := language.String()
	if langCode == "" {
		return EN
	}

	// Parse language code (e.g., "es-ES" -> "es")
	code := strings.Split(langCode, "-")[0]
	code = strings.ToLower(code)

	// Map to lang enum
	switch code {
	case "es":
		return ES
	case "pt":
		return PT
	case "fr":
		return FR
	case "ru":
		return RU
	case "de":
		return DE
	case "it":
		return IT
	case "hi":
		return HI
	case "bn":
		return BN
	case "id":
		return ID
	case "ar":
		return AR
	case "ur":
		return UR
	case "zh":
		return ZH
	default:
		return EN
	}
}
