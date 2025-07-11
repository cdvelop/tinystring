//go:build wasm

package tinystring

import (
	"syscall/js"
)

// getSystemLang detects browser language from navigator.language
func (c *conv) getSystemLang() lang {
	// Get browser language
	navigator := js.Global().Get("navigator")
	if navigator.IsUndefined() {
		return EN
	}

	language := navigator.Get("language")
	if language.IsUndefined() {
		return EN
	}

	// Use the centralized parser.
	return c.langParser(language.String())
}
