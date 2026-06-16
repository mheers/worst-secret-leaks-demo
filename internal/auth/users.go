// SAFE Go file: a struct that loads secrets from the environment.
//
// The string literals below are deliberately low-entropy or empty so the
// scanner does not flag this file. The real secrets are passed in via
// `LoadFromEnv` at runtime.
package auth

import "os"

type APIKey struct {
	Provider string
	Value    string
}

func LoadFromEnv() []APIKey {
	return []APIKey{
		{Provider: "openai", Value: os.Getenv("OPENAI_API_KEY")},
		{Provider: "stripe", Value: os.Getenv("STRIPE_SECRET_KEY")},
		{Provider: "slack", Value: os.Getenv("SLACK_BOT_TOKEN")},
	}
}
