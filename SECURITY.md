# Security Policy

Disk Cleaner is intended to help users inspect and remove wasted disk space. Because the project may interact with local files, directories, caches, and deletion workflows, security and data safety issues are treated seriously.

The project owner and security contact is `deluxesande`.

## Supported Versions

This project is currently in an early documentation stage. Once releases are available, security support will generally apply to the latest released version.

| Version | Supported |
| --- | --- |
| Latest release | Yes |
| Older releases | No, unless explicitly stated |
| Unreleased local changes | No |

## Reporting a Vulnerability

Do not open a public GitHub issue for security vulnerabilities.

Report suspected vulnerabilities privately to `deluxesande` using the private contact method available on the repository, account profile, or hosting platform.

When reporting a vulnerability, include as much detail as you safely can:

- A clear description of the issue
- Steps to reproduce it
- Affected operating system
- Affected version or commit, if known
- Whether the issue can cause data loss, unexpected deletion, privilege escalation, path traversal, denial of service, or information disclosure
- Any proof-of-concept commands or files needed to demonstrate the issue

Do not include sensitive personal files, credentials, private tokens, or data from systems you do not own or have permission to test.

## Security Scope

Examples of security issues that should be reported privately include:

- Deleting files that were not selected or confirmed by the user
- Bypassing configured exclusions
- Following symlinks or junctions into unsafe locations unexpectedly
- Path traversal bugs
- Unsafe handling of protected operating system directories
- Incorrect duplicate detection that could lead to accidental data loss
- Permission handling that exposes private data
- Crashes or hangs caused by malicious directory structures
- Terminal UI behavior that misrepresents what will be deleted

Issues that are usually not security vulnerabilities:

- General feature requests
- Documentation typos
- Performance problems without a security or data-loss impact
- Expected permission errors while scanning protected folders
- Duplicate files being reported correctly but undesirably

## Response Expectations

After a vulnerability is reported, `deluxesande` will aim to:

1. Acknowledge the report.
2. Review the impact and reproduction steps.
3. Decide whether the issue is accepted as a security vulnerability.
4. Prepare a fix or mitigation when needed.
5. Credit the reporter if requested and appropriate.

Response times may vary because this is an independently maintained project.

## Disclosure Guidelines

Please do not publicly disclose a vulnerability until maintainers have had reasonable time to investigate and prepare a fix.

Coordinated disclosure helps protect users from accidental data loss and unsafe cleanup behavior.

## Security Design Principles

Security-sensitive contributions should follow these principles:

- Never delete files without explicit user confirmation.
- Show users exactly what will be deleted before deletion begins.
- Treat symlinks, junctions, mount points, and shortcuts carefully.
- Keep protected system paths excluded by default.
- Fail closed when path safety is uncertain.
- Prefer exact duplicate detection over assumptions or fuzzy matching.
- Avoid logging private file contents.
- Handle permission errors gracefully.
