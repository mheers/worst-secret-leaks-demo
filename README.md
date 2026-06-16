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

The table below lists the **actual** findings gitleaks emits against this
folder (gitleaks v8.30.x, default ruleset + the `demo-todo-marker` custom
rule). Six files produce findings, the other seven are clean control
cases.

| File                                   | Actual findings          | Notes                                                                                                                                                                                 |
| -------------------------------------- | ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cmd/server/main.go`                   | (none)                   | The hardcoded `ghp_…` token ends in `0123456789` (low entropy tail) and is below the rule's entropy floor in this build. Move the literal to `.env` to exercise the rule.             |
| `cmd/hello/main.go`                    | (none)                   | Control case.                                                                                                                                                                         |
| `internal/auth/middleware.go`          | (none)                   | High-entropy constant that is _not_ flagged — confirms the scanner does not over-match.                                                                                               |
| `internal/auth/users.go`               | (none)                   | Empty env reads.                                                                                                                                                                      |
| `internal/auth/testdata/test_rsa.pem`  | `private-key`            | Public jwt-go README sample.                                                                                                                                                          |
| `internal/auth/testdata/empty_pem.pem` | (none)                   | Body too short to match.                                                                                                                                                              |
| `internal/config/config.go`            | (none)                   | Env-only.                                                                                                                                                                             |
| `internal/config/config.toml`          | `heroku-api-key`         | The `heroku_api_key` value matches the v8.30.x rule. SendGrid key is _not_ matched by `sendgrid-api-token` here (the test value's entropy is too low); that's the realistic behavior. |
| `.env`                                 | `slack-bot-token`, `jwt` | `AKIA…EXAMPLE` and `ghp_…` are silently dropped by gitleaks' built-in `.+EXAMPLE$` allowlist on the `aws-access-token` rule. The xoxb- token and the JWT (entropy > 4.0) both fire.   |
| `.env.example`                         | (none)                   | Clean template.                                                                                                                                                                       |
| `.env.production.example`              | (none)                   | Clean template.                                                                                                                                                                       |
| `infra/k8s-secrets.yaml`               | `slack-webhook-url`      | AWS values are dropped by the same `EXAMPLE$` allowlist.                                                                                                                              |
| `infra/firebase.json`                  | (none — generic-api-key) | The `AIzaSyA-…` value is below the rule's entropy floor in this build. To exercise `generic-api-key`, swap it for a 40+ char random secret.                                           |
| `scripts/bootstrap.sh`                 | `curl-auth-user`         | The inline `-u demo:AKIA…` matches the `curl-auth-user` rule. The exported `GITHUB_TOKEN` and `SENDGRID_KEY` are dropped by allowlists.                                               |

If you want **all** of the obvious rules (`github-pat`, `aws-access-token`,
`openai-api-key`, `sendgrid-api-token`, `generic-api-key`, `firebase`)
to fire, replace the placeholder values with longer, higher-entropy
strings. The current fixtures bias toward realistic-but-known test
strings so the demo is safe to commit.

## How to scan it manually

```sh
# Scan the entire demo folder
gitleaks detect --no-git --source demo --report-format json --exit-code 0 --report-path /tmp/gl.json

# Scan a single file via stdin (note: --pipe is the v8 flag, not --stdin)
cat demo/.env | gitleaks detect --no-git --pipe --report-format json --exit-code 0 --report-path /tmp/gl.json

# Scan with the custom config in this folder (adds the demo-todo-marker rule)
gitleaks detect --no-git --config demo/.gitleaks.toml --source demo --report-format json --exit-code 0 --report-path /tmp/gl.json
```

Expected counts with the current fixtures: **6 findings from 5 files,
hitting 6 distinct rules** (`private-key`, `slack-webhook-url`,
`heroku-api-key`, `curl-auth-user`, `jwt`, `slack-bot-token`).

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
