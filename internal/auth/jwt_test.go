// jwt_test.go is a test file with a *real* example RSA private key.
//
// The PEM below is the well-known public test key from the `jwt-go`
// README. It is short enough to be obviously non-secret, but gitleaks
// still flags it as `private-key` — a useful test of the scanner.
package auth_test

import (
	_ "embed"
)

//go:embed testdata/test_rsa.pem
var testRSAKey []byte

//go:embed testdata/empty_pem.pem
var emptyPEM []byte

const placeholderJWTSecret = "this-is-not-a-real-secret-please-ignore"
