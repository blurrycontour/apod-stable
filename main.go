package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"apod-stable/apod"
	"apod-stable/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	apod.UserAgent = getEnv("USER_AGENT", apod.UserAgent)
	apod.PageURL = getEnv("APOD_URL", apod.PageURL)
	listenAddr := getEnv("LISTEN_ADDR", ":8080")

	mux := http.NewServeMux()
	server.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("APOD scraper listening on %s", listenAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
