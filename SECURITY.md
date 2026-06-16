# Security Policy

## ⚠️ This repository is for TESTING / SECURITY RESEARCH only ⚠️

The `worst-secret-leaks-demo` repository is an **intentionally
vulnerable** fixture project. Every "leaked" secret, API key, token,
webhook URL, and private key contained in this repository is:

- **Fake / non-functional** — the values are constructed to look like
  real leaks so that secret-scanning tooling has realistic input, but
  they are not valid credentials and grant access to **nothing**.
- **Placed here deliberately** — for the sole purpose of supporting
  security research (e.g. tuning detection rules, measuring
  false-positive / false-negative rates, and demonstrating the impact
  of hardcoded credentials).

### What this means for you

- **DO NOT** use any value found in this repository against any real
  service, account, organization, or system. There is nothing to
  authenticate against — attempting to do so would be pointless and
  potentially abusive to whoever owns the prefixes (e.g. AWS, GitHub,
  Slack, OpenAI, SendGrid, Heroku, Google/Firebase).
- **DO NOT** copy values from this repository into real code,
  configuration, CI/CD pipelines, container images, or production
  tooling. Even though the values are fake, copying them risks leaking
  *your own* secrets through the same patterns.
- **DO NOT** treat the patterns shown here as a model for handling
  real secrets in your own projects. This is the **opposite** of best
  practice and exists only to prove that scanners can catch such
  patterns.
- **DO** use this repository for security research, scanner
  evaluation, and educational demonstrations of why secrets should
  never be hardcoded.

### Reporting a real secret

If you believe you have found a **real, valid** credential that
accidentally landed in this repository, please open a private issue
or contact the maintainers directly. Note that the project is a test
fixture, so this should not happen by design — but if it does, it
should be revoked at the source and removed from the repo.

### Acceptable use

Acceptable use of this repository is limited to:

1. Running secret-scanning tools against it to evaluate detection
   coverage and false-positive behavior.
2. Demonstrating the impact of hardcoded credentials in security
   training, blog posts, talks, or academic work — with appropriate
   context that all values are fake.
3. Developing, testing, and validating secret-scanning tools and
   similar security tooling.

Any other use is out of scope.
