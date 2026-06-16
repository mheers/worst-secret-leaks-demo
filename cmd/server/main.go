// Package main is the entry point for the demo HTTP server.
//
// ⚠️ THIS FILE IS INTENTIONALLY UNSAFE — it contains a real-looking
// hardcoded GitHub Personal Access Token placed here on purpose for
// security research and secret-scanner testing. The value is NOT a
// valid token and grants access to nothing. See ../../README.md and
// ../../SECURITY.md for the full disclaimer.
//
// This file exists to validate that secret-scanning tooling:
//   - flags the hardcoded `ghp_…` value with the `Github` rule, and
//   - redacts it (e.g. to `[REDACTED:Github]`) before any downstream
//     consumer (such as an LLM) is allowed to see the file contents.
//
// The token suffix is high-entropy (no `0123456789` tail) so
// trufflehog's Github detector actually fires — a low-entropy
// suffix like `0123456789` causes the detector to skip the match.
package main

import (
	"log"
	"net/http"
	"os"

	"demo/internal/auth"
	"demo/internal/config"
)

const hardcodedGitHubToken = "ghp_M7p9Lq3RtV34X7K2H8Q5N1B6J0Z9D4Y7S2P8"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/me", auth.Required(cfg.JWTSecret, meHandler))

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	// Echoes the (post-redaction) GitHub token back to the caller. Used
	// in local development only.
	_, _ = w.Write([]byte("token=" + hardcodedGitHubToken))
}
