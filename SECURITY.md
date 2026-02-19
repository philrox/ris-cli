# Security Policy

## Supported Versions

Only the latest release is supported with security updates.

| Version | Supported |
|---------|-----------|
| latest  | Yes       |
| < latest | No       |

## Attack Surface

This CLI has a minimal attack surface:

- **Read-only** — all operations are HTTP GET requests to the public RIS API
- **No authentication** — no secrets, tokens, or credentials are stored or transmitted
- **No state** — the CLI is stateless, no local databases or caches
- **SSRF protection** — document fetching is restricted to allowed hosts only: `data.bka.gv.at`, `www.ris.bka.gv.at`, `ris.bka.gv.at` (HTTPS only)
- **Input validation** — document numbers are validated against a strict pattern before any URL construction

## Reporting a Vulnerability

If you discover a security issue, please [open a GitHub Issue](https://github.com/philrox/ris-cli/issues/new) with the label **security**.

For sensitive issues that should not be disclosed publicly, please use [GitHub's private vulnerability reporting](https://github.com/philrox/ris-cli/security/advisories/new).

We aim to respond within 7 days.
