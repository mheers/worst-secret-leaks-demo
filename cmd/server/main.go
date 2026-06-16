// Package main is the entry point for the demo HTTP server.
//
// ⚠️ THIS FILE IS INTENTIONALLY UNSAFE — it contains a real-looking
// hardcoded GitHub Personal Access Token placed here on purpose for
// security research and secret-scanner testing. The value is NOT a
// valid token and grants access to nothing. See ../../README.md and
// ../../SECURITY.md for the full disclaimer.
//
// This file exists to validate that secret-scanning tooling:
//   - flags the hardcoded `ghp_…` value with the `github-pat` rule, and
//   - redacts it (e.g. to `[REDACTED:github-pat]`) before any downstream
//     consumer (such as an LLM) is allowed to see the file contents.
package main

import (
	"log"
	"net/http"
	"os"

	"demo/internal/auth"
	"demo/internal/config"
)

const hardcodedGitHubToken = "ghp_aBcDeFgHiJkLmNoPqRsTuVwXyZ0123456789"

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
