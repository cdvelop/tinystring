//go:build !wasm

package tinystring

import (
	"os"
)

// getSystemLang detects system language from environment variables
func getSystemLang() lang {
	// Use the centralized parser with common environment variables.
	return langParser(
		os.Getenv("LANG"),
		os.Getenv("LANGUAGE"),
		os.Getenv("LC_ALL"),
		os.Getenv("LC_MESSAGES"),
	)
}
