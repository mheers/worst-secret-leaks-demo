// Package main is a CLEAN demo entry point. It has no hardcoded secrets
// and exists so the scanner has a known-good control case to compare
// against.
package main

import "fmt"

// A constant that LOOKS high-entropy but is the well-known OAuth2
// authorize endpoint string. The default `gitleaks` ruleset does not
// flag it, which is the correct behavior.
const oauth2AuthEndpoint = "/oauth/authorize"

// Public, well-known port numbers. Scanner should leave these alone.
var wellKnownPorts = []int{80, 443, 8080, 8443}

func main() {
	for _, p := range wellKnownPorts {
		fmt.Printf("would listen on %d%s\n", p, oauth2AuthEndpoint)
	}
}
