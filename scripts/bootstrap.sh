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
curl -fsSL -u "demo:AKIAIOSFODNN7EXAMPLE" https://api.example.com/v1/whoami

# 2. Heredoc with embedded secret.
SENDGRID_KEY=$(cat <<'EOF'
SG.AbCdEfGhIjKlMnOpQrStUvWxYz0123456789._aBcDeFgHiJkLmNoPqRsTuVwX
EOF
)
export SENDGRID_KEY

# 3. Inline-export of a personal access token.
export GITHUB_TOKEN="ghp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789"

# 4. Clean: parameter expansion only, no value.
export DATABASE_URL="${DATABASE_URL:-}"
export PORT="${PORT:-8080}"

echo "bootstrap complete"
