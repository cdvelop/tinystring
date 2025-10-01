//go:build !wasm
// +build !wasm

package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	publicDir := "public" // Template variable

	// Debug: Print working directory and check if public exists
	_, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
	}

	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		log.Printf("WARNING: Public directory '%s' does not exist!", publicDir)
	}

	// Serve static files with no-cache headers
	fs := http.FileServer(http.Dir(publicDir))

	// Middleware to disable caching for static files (useful in dev/test)
	noCache := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent browser caching
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			h.ServeHTTP(w, r)
		})
	}

	// Use a dedicated ServeMux so we can pass it to an http.Server
	mux := http.NewServeMux()
	mux.Handle("/", noCache(fs))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Server is running 3"))
	})

	// Create http.Server with Addr and Handler set
	server := &http.Server{
		Addr:    ":4430",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
