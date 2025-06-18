//go:build !wasm

package tinystring

import (
	"os"
	"strings"
)

// getSystemLang detects system language from environment variables
func getSystemLang() lang {
	// Common environment variables that contain language information
	langVars := []string{"LANG", "LANGUAGE", "LC_ALL", "LC_MESSAGES"}

	for _, envVar := range langVars {
		if envValue := os.Getenv(envVar); envValue != "" {
			// Parse language code from environment variable
			code := strings.Split(envValue, ".")[0] // Remove encoding part
			code = strings.Split(code, "_")[0]      // Get language part
			code = strings.Split(code, "-")[0]      // Handle dash format
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
	}

	return EN // Default fallback
}
