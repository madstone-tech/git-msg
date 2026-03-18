# Security Policy

## Supported Versions

Only the latest release receives security fixes.

| Version | Supported |
| ------- | --------- |
| latest  | ✅        |
| older   | ❌        |

## Reporting a Vulnerability

Please **do not** open a public GitHub issue for security vulnerabilities.

Report vulnerabilities privately via GitHub's
[Security Advisories](https://github.com/madstone-tech/git-msg/security/advisories/new)
feature (Security → Report a vulnerability).

Include:
- A description of the vulnerability and its potential impact
- Steps to reproduce or a proof-of-concept
- The version of `git-msg` affected (`git-msg --version`)
- Your suggested fix, if any

You can expect an acknowledgement within 48 hours and a resolution or
update within 14 days depending on severity.

## Scope

Areas of particular interest:

- **Credential handling** — API keys must only be stored in the system
  keychain and must never appear in config files, logs, or CLI output
- **Hook script injection** — the installed `prepare-commit-msg` script
  must not be exploitable via crafted branch names or diff content
- **Dependency vulnerabilities** — CVEs in any of the direct dependencies
  listed in `go.mod`
