# Security Policy

## Supported Versions

Only the latest release of golangster receives security updates.

| Version | Supported |
|---------|-----------|
| latest  | yes       |
| older   | no        |

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Please report security issues using [GitHub private vulnerability reporting](https://github.com/idakhno/golangster/security/advisories/new).

Include the following information in your report:

- Description of the vulnerability and its potential impact
- Steps to reproduce (minimal reproducer preferred)
- Affected version(s)
- Any known mitigations or workarounds

## Response Process

| Step | Timeline |
|------|----------|
| Acknowledgement | within 3 business days |
| Initial assessment | within 7 business days |
| Fix and coordinated disclosure | within 90 days |

We will keep you informed throughout the process. If you do not receive an acknowledgement within 3 business days, please follow up.

## Coordinated Disclosure

We follow a coordinated disclosure model:

1. Reporter submits a private report
2. We confirm, triage, and develop a fix
3. We release the fix and publish a security advisory simultaneously
4. CVE is requested after public disclosure if applicable

We ask that you do not publicly disclose the vulnerability until we have released a fix or until 90 days have elapsed, whichever comes first.

## Scope

The following are considered in scope for security reports:

- **Malicious input processing** — a crafted Go source file causes golangster to panic, execute arbitrary code, or produce incorrect output that silently bypasses the `sensitive` rule
- **False negatives in the `sensitive` rule** — sensitive data patterns that should be flagged but are not, in a way that represents a systematic bypass
- **Supply chain issues** — compromised dependencies or build artifacts

The following are **out of scope**:

- Style or correctness of lint rules (use a regular GitHub issue)
- Vulnerabilities in the Go toolchain itself (report to [the Go team](https://go.dev/security))
- Security of applications being linted by golangster

## Security Considerations for Users

golangster runs as part of your build pipeline and analyzes your source code. Keep the following in mind:

- Pin golangster to a specific version in CI to avoid unexpected changes
- Verify checksums when downloading pre-built binaries
- The `sensitive` rule is a best-effort heuristic — it does not replace a proper secrets scanner

## Credits

We will publicly acknowledge reporters in the security advisory unless they request otherwise.
