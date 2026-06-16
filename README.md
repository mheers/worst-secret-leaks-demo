# Worst Secret Leaks Demo — TESTING / SECURITY RESEARCH ONLY

> ⚠️ **WARNING — FOR TESTING PURPOSES ONLY ⚠️**
>
> This repository is an **intentionally vulnerable** fixture project. It
> exists **solely** to support **security research** (e.g. training and
> tuning secret-detection rules, evaluating scanner false-positive /
> false-negative rates, and demonstrating the impact of hardcoded
> credentials).
>
> **Every "leaked" secret in this repo is fake / non-functional and is
> placed here on purpose.** No real credentials, API keys, tokens, or
> private keys are present. **Nothing in this repository should ever be
> used in any non-test environment, against any real service, or for
> any purpose other than security research / scanner testing.** Do not
> copy, paste, or extract any of the values shown here into real code,
> configs, or tooling.

## About

This directory is a **fixture project** used to support general
security research, including the development, testing, and evaluation
of secret-scanning tools. It contains a realistic mix of clean code,
hardcoded secrets, env-loaded secrets, false-positive-prone strings,
and various config formats.

It is **not** a runnable Go service — files are arranged only to exercise
secret-scanning tooling. Do not start it. The keys, tokens, and
credentials that appear throughout this repo are deliberate canaries:
they are constructed to look like real leaks (so scanners have
something realistic to find) but are **not valid** and have **no
associated account, service, or data**.

## Layout

```
demo/
├── README.md                       # this file
├── go.mod                          # module declaration
├── .env                            # DIRTY: hardcoded AWS/GitHub/Slack/OpenAI keys
├── .env.example                    # clean template
├── .env.production.example         # clean template
├── .gitleaks.toml                  # allowlist for known false positives
├── cmd/
│   ├── server/main.go              # DIRTY: hardcoded ghp_ token
│   └── hello/main.go               # CLEAN: control case
├── internal/
│   ├── auth/
│   │   ├── middleware.go           # high-entropy non-secret (false-positive test)
│   │   ├── jwt_test.go             # embeds testdata/*.pem
│   │   ├── users.go                # CLEAN: loads from os.Getenv
│   │   └── testdata/
│   │       ├── test_rsa.pem        # RSA private key (flagged as `private-key`)
│   │       └── empty_pem.pem       # empty PEM (NOT flagged)
│   └── config/
│       ├── config.go               # CLEAN: env-only
│       └── config.toml             # DIRTY: heroku_api_key, sendgrid, password in URL
├── infra/
│   ├── k8s-secrets.yaml            # DIRTY: AWS keys, ghp_, slack webhook
│   └── firebase.json               # DIRTY: AIza... API key
└── scripts/
    └── bootstrap.sh                # DIRTY: inline curl creds, export GITHUB_TOKEN
```

## What each file is testing

The table below lists the **actual** findings this project is designed to
trigger. Results vary by scanner backend:

| File                                   | trufflehog 3.95.5 findings | gitleaks v8.x findings       | Notes                                                                                                                                                                                                |
| -------------------------------------- | -------------------------- | ---------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/server/main.go`                   | `Github`                   | `github-pat` (v8.18+)        | Hardcoded `ghp_` with a high-entropy suffix. Low-entropy suffixes (e.g. `0123456789` tail) cause both scanners to skip the match.                                                                    |
| `cmd/hello/main.go`                    | (none)                     | (none)                       | Control case — no secrets.                                                                                                                                                                           |
| `internal/auth/middleware.go`          | (none)                     | (none)                       | High-entropy constant that is _not_ flagged — confirms the scanner does not over-match.                                                                                                              |
| `internal/auth/users.go`               | (none)                     | (none)                       | Empty env reads.                                                                                                                                                                                     |
| `internal/auth/testdata/test_rsa.pem`  | `PrivateKey`               | `private-key`                | Public jwt-go README sample PEM.                                                                                                                                                                     |
| `internal/auth/testdata/empty_pem.pem` | (none)                     | (none)                       | Body too short to match.                                                                                                                                                                             |
| `internal/config/config.go`            | (none)                     | (none)                       | Env-only — no hardcoded secrets.                                                                                                                                                                     |
| `internal/config/config.toml`          | `Postgres`                 | `heroku-api-key`, `url-cred` | The postgres URL with embedded credentials fires the `Postgres` rule. Heroku and SendGrid keys use high-entropy hex / SG. format.                                                                    |
| `.env`                                 | `Github`, `GoogleGemini…`  | `slack-bot-token`, `jwt`     | Every line in `.env` now carries a high-entropy token. The `AKIA…` value no longer ends in `EXAMPLE` (which would be silently allow-listed by trufflehog's built-in false-positive filter).          |
| `.env.example`                         | (none)                     | (none)                       | Clean template — no secrets.                                                                                                                                                                         |
| `.env.production.example`              | (none)                     | (none)                       | Clean template — no secrets.                                                                                                                                                                         |
| `infra/k8s-secrets.yaml`               | `Github`, `SlackWebhook`   | `slack-webhook-url`          | AWS values use high-entropy suffixes (no `EXAMPLE` tail). Git PAT uses the same high-entropy format as `.env`.                                                                                       |
| `infra/firebase.json`                  | `GoogleGeminiAPIKey`       | (none)                       | The `AIzaSy…` key now has a full 35-character high-entropy suffix so the `GoogleGeminiAPIKey` detector fires. Older gitleaks versions match it via `generic-api-key`.                                |
| `scripts/bootstrap.sh`                 | `Github`                   | `curl-auth-user`             | The inline `-u demo:AKIA…` matches the `curl-auth-user` rule (gitleaks). The exported `GITHUB_TOKEN` and `SENDGRID_KEY` use high-entropy values. Trufflehog's `Github` detector fires on the former. |

**Key design principles for the fake secrets in this demo:**

1. Every value that is **meant to be found** uses a high-entropy, real-looking
   payload that matches the detector's regex and passes the scanner's built-in
   entropy / false-positive filters.
2. Values that end in `EXAMPLE` or `0123456789` or that omit required
   substrings (like `T3BlbkFJ` for the OpenAI detector) will be silently
   skipped by the scanner — they are useful only as **negative controls** or
   regression tests for allow-list behaviour.
3. Because detector coverage varies between scanner versions and between
   tools (trufflehog vs gitleaks), a value that fires in one may not fire
   in the other. The table above documents both.

## How to scan it manually

```sh
# Scan the entire demo folder with trufflehog
trufflehog filesystem --no-verification --no-update --no-color --json \
  --log-level=-1 --results=unverified,unknown . 2>/dev/null | jq -c '{detector: .DetectorName, file: .SourceMetadata.Data.Filesystem.file, raw: (.Raw | .[0:40])}'

# Scan with gitleaks
gitleaks detect --no-git --source . --report-format json --exit-code 0 --report-path /tmp/gl.json
cat /tmp/gl.json | jq -c '.[] | {rule: .RuleID, file: .File, secret: (.Secret | .[0:40])}'
```

Expected counts with the current fixtures:

| Scanner         | Findings | Files hit | Distinct detectors                                                                               |
| --------------- | -------- | --------- | ------------------------------------------------------------------------------------------------ |
| trufflehog 3.95 | 9        | 7         | `Github` (×4), `Postgres` (×2), `PrivateKey`, `SlackWebhook`, `GoogleGeminiAPIKey`               |
| gitleaks v8.x   | 6        | 5         | `private-key`, `slack-webhook-url`, `heroku-api-key`, `curl-auth-user`, `jwt`, `slack-bot-token` |

## How secret-scanning tools can use this

A typical secret-scanning workflow against this fixture looks like:

1. Point a secret scanner (e.g. `gitleaks`, `trufflehog`, `detect-secrets`,
   or any custom rule engine) at the contents of this folder.
2. Observe which files are flagged, which rules fire, and how the
   scanner handles false-positive-prone strings (high-entropy constants,
   empty PEMs, example placeholders, etc.).
3. Tune rules, evaluate false-positive / false-negative rates, and
   build regressions against the expected findings listed above.

For example, the expected behaviour of any well-configured scanner on
the files in this repo is that every `ghp_…`, `AKIA…`, `xoxb-…`,
`sk-…`, `AIza…`, and PEM block is identified and replaced with a
redaction placeholder (e.g. `[REDACTED:<rule-id>]`) before the
content is sent to any downstream consumer such as an LLM.

---

## Disclaimer

This repository is a **deliberately insecure test fixture** for security
research and secret-scanner validation. It is published for **testing
purposes only**.

- All "leaked" secrets, API keys, tokens, webhooks, and private keys
  appearing anywhere in this repository are **fake canaries** placed
  here **on purpose**.
- They are constructed to look like real leaks so that secret-scanning
  tooling has realistic input to match against, but they are **not
  valid credentials** and grant access to **nothing**.
- **Do not** use any value found in this repository against any real
  service, account, or system.
- **Do not** copy values from this repository into real code,
  configuration, CI/CD pipelines, or production tooling.
- **Do not** treat the patterns in this repository as a model for
  handling real secrets in your own projects — this is the opposite of
  best practice and exists only to prove that scanners can catch them.

If you are a scanner author, a security researcher, or a reviewer
evaluating secret-scanning tools, you are in the right place.
Everyone else: please look away and use a `.env.example` instead.
