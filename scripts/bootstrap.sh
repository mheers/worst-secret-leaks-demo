#!/usr/bin/env bash
# bootstrap.sh — INTENTIONALLY DIRTY.
#
# This script demonstrates several common ways secrets leak into shell:
#   1. Hardcoded curl with `-u user:password` (matches generic-api-key).
#   2. Heredoc with a token assigned inline.
#   3. `export FOO=...` of a real-looking PAT.
#
# All values below are PUBLIC test strings (AWS example key, gitleaks
# README examples) so the file is safe to commit.
set -euo pipefail

# 1. Hardcoded credentials in a curl command.
# The AKIA value uses a high-entropy suffix (no EXAMPLE tail) so
# the `curl-auth-user` and `aws-access-token` detectors both fire.
curl -fsSL -u "demo:AKIAIOSFODNN7ABCDEFGH" https://api.example.com/v1/whoami

# 2. Heredoc with embedded secret (matching the SendGrid SG. format:
#    20-24 word-characters before the first dot, then 39-50 after).
SENDGRID_KEY=$(cat <<'EOF'
SG.AbCdEfGhIjKlMnOpQrStUv._aBcDeFgHiJkLmNoPqRsTuVwXyZ0123456789ABCDEFGHIJKLM
EOF
)
export SENDGRID_KEY

# 3. Inline-export of a personal access token (high-entropy suffix).
export GITHUB_TOKEN="ghp_M7p9Lq3RtV34X7K2H8Q5N1B6J0Z9D4Y7S2P8"

# 4. Clean: parameter expansion only, no value.
export DATABASE_URL="${DATABASE_URL:-}"
export PORT="${PORT:-8080}"

echo "bootstrap complete"
