package server

import (
	"fmt"
	"log"
	"net/http"

	"apod-stable/apod"
)

// RegisterRoutes attaches all handlers to mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/apod", apodHandler)
	mux.HandleFunc("/health", healthHandler)
}

func apodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := apod.Fetch()
	if err != nil {
		log.Printf("ERROR scraping APOD: %v", err)
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(data) //nolint:errcheck
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error":%q}`, msg)
}
