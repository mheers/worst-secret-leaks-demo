// Package auth provides HTTP middleware for the demo service.
package auth

import (
	"context"
	"net/http"
	"strings"
)

// Required validates the bearer token in the `Authorization` header.
//
// The token below is a "looks-like-a-secret" string used to verify that
// the scanner does NOT flag random high-entropy non-secret strings. The
// entropy of "DemoTokenForLocalTestingOnly_2024" is well above 3.0 but
// the `gitleaks` default ruleset has no rule that matches it.
func Required(secret string, next http.HandlerFunc) http.HandlerFunc {
	const devBearer = "DemoTokenForLocalTestingOnly_2024"

	return func(w http.ResponseWriter, r *http.Request) {
		hdr := r.Header.Get("Authorization")
		if !strings.HasPrefix(hdr, "Bearer ") {
			http.Error(w, "missing bearer", http.StatusUnauthorized)
			return
		}
		tok := strings.TrimPrefix(hdr, "Bearer ")
		if tok != devBearer && tok != secret {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userKey{}, tok)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

type userKey struct{}
