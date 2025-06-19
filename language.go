package tinystring

// Private global configuration
var defLang lang = EN

// Language enumeration for supported languages
type lang uint8

const (
	// Group 1: Core Essential Languages (Maximum Global Reach)
	EN lang = iota // 0 - English (default)
	ES             // 1 - Spanish
	ZH             // 2 - Chinese
	HI             // 3 - Hindi
	AR             // 4 - Arabic

	// Group 2: Extended Reach Languages (Europe & Americas)
	PT // 5 - Portuguese
	FR // 6 - French
	DE // 7 - German
	RU // 8 - Russian

	// Group 3: Regional Languages (Commented out to reduce binary size)
	// IT             // Italian
	// ID             // Indonesian
	// BN             // Bengali
	// UR             // Urdu
)

// LocStr represents a string with translations for multiple languages.
//
// It is a fixed-size array where each index corresponds to a language constant
// (EN, ES, PT, etc.). This design ensures type safety and efficiency, as the
// compiler can verify that all translations are provided.
//
// The order of translations must match the order of the language constants.
//
// Example of creating a new translatable term for "File":
//
//	var MyDictionary = struct {
//		File LocStr
//	}{
//		File: LocStr{
//			EN: "file",
//			ES: "archivo",
//			ZH: "文件",
//			HI: "फ़ाइल",
//			AR: "ملف",
//			PT: "arquivo",
//			FR: "fichier",
//			DE: "Datei",
//			RU: "файл",
//		},
//	}
//
// Usage in code:
//	err := Err(MyDictionary.File, D.Not, D.Found) // -> "file not found", "archivo no encontrado", etc.
type LocStr [9]string

// get returns translation for specified language with English fallback
func (o LocStr) get(l lang) string {
	if int(l) < len(o) && o[l] != "" {
		return o[l]
	}
	return o[EN] // Fallback to English
}

// OutLang sets the default output language
// OutLang() without parameters auto-detects system language
// OutLang(ES) sets Spanish as default
func OutLang(l ...lang) {
	if len(l) == 0 {
		defLang = getSystemLang() // from env.front.go or env.back.go
	} else {
		defLang = l[0]
	}
}

// langParser processes a list of language strings (e.g., from env vars or browser settings)
// and returns the first valid language found. It centralizes the parsing logic for both
// frontend and backend environments.
func langParser(langStrings ...string) lang {
	for _, langStr := range langStrings {
		if langStr == "" {
			continue
		}

		// Parse language code from the string, handling common formats.
		code := Split(langStr, ".")[0] // Removes encoding, e.g., ".UTF-8"
		code = Split(code, "_")[0]     // Handles locale format, e.g., "en_US"
		code = Split(code, "-")[0]     // Handles standard format, e.g., "en-US"

		if code == "" {
			continue
		}

		// Convert to lowercase and map to the internal lang type.
		code = Convert(code).ToLower().String()
		return mapLangCode(code)
	}

	return EN // Default fallback if no valid language string is found.
}

// mapLangCode maps a language code string (e.g., "es") to the lang enum.
// It's the single source of truth for language code mapping.
func mapLangCode(code string) lang {
	switch code {
	// Group 1
	case "es":
		return ES
	case "zh":
		return ZH
	case "hi":
		return HI
	case "ar":
		return AR
	// Group 2
	case "pt":
		return PT
	case "fr":
		return FR
	case "de":
		return DE
	case "ru":
		return RU
	default:
		return EN
	}
}
