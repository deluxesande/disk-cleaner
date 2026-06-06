# Contributing to Disk Cleaner

Thank you for your interest in contributing to Disk Cleaner. This project is owned and maintained by `deluxesande`.

Disk Cleaner is intended to be a safe, reliable, cross-platform tool for reviewing and removing wasted disk space. Contributions should keep that goal in mind, especially when changing scanning, deletion, or duplicate-detection behavior.

## Project Owner

The account owner and final decision maker for this project is:

- GitHub/account owner: `deluxesande`

Major product decisions, licensing changes, release approvals, repository settings, and maintainer permissions must be approved by `deluxesande`.

## Contribution Rules

Please follow these rules when contributing:

- Keep user safety first. Any feature that deletes, moves, overwrites, or modifies files must be explicit, reviewable, and reversible where possible.
- Do not add automatic deletion behavior without a clear confirmation step.
- Do not scan or modify protected operating system paths unless there is a deliberate, reviewed reason.
- Keep platform-specific behavior isolated so Windows, macOS, and Linux support remains maintainable.
- Prefer clear, readable code over clever shortcuts.
- Keep changes focused. Avoid unrelated refactors in the same pull request.
- Include tests for changes that affect scanning, duplicate detection, deletion, path filtering, or configuration.
- Update documentation when behavior, commands, flags, configuration, or safety guarantees change.
- Do not commit secrets, tokens, credentials, private paths, or machine-specific configuration.
- Respect the existing license and do not introduce dependencies with incompatible licenses.

## Development Expectations

Before opening a pull request, contributors should:

1. Read the README and understand the intended behavior of the tool.
2. Make a focused branch for the change.
3. Keep commits clear and descriptive.
4. Run formatting and tests once the implementation exists.
5. Verify that destructive actions require confirmation.
6. Document any new flags, commands, configuration options, or safety rules.

## Pull Request Guidelines

Pull requests should include:

- A short description of what changed
- The reason for the change
- Any risks or edge cases
- Testing performed
- Screenshots or terminal output for user-interface changes, when useful

Pull requests may be rejected or sent back for revision if they are too broad, unsafe, undocumented, untested, or inconsistent with the project direction.

## Safety Requirements

Because Disk Cleaner is a cleanup tool, safety is a core project requirement.

Any code that removes files or folders must:

- Show the user what will be deleted before deletion begins
- Require explicit confirmation
- Respect configured exclusions
- Avoid protected system directories by default
- Handle permission errors without crashing
- Report what was deleted and what failed

Duplicate-detection changes must only mark files as duplicates when they are confirmed to be byte-for-byte identical.

## Code Style

When implementation files are added, contributors should follow the style already used in the repository. In general:

- Use clear names for packages, files, functions, and variables.
- Keep small functions focused on one responsibility.
- Prefer structured path handling over manual string manipulation.
- Keep operating-system-specific logic easy to find and test.
- Add comments only when they explain behavior that is not obvious from the code.

## Reporting Issues

When reporting a bug, include:

- Operating system and version
- Tool version or commit, if available
- Command used
- Target directory type, such as project folder, home folder, or external drive
- Expected behavior
- Actual behavior
- Error output, if any

Do not include private file contents, secrets, credentials, or personal directory listings unless you have removed sensitive information.

## Feature Requests

Feature requests should explain:

- The problem being solved
- The proposed behavior
- Why the feature belongs in Disk Cleaner
- Any safety concerns
- Whether the feature should be enabled by default or opt-in

Features that make deletion faster but less reviewable are unlikely to be accepted.

## Maintainer Authority

`deluxesande` may close issues or pull requests that are inactive, unsafe, out of scope, abusive, spammy, or inconsistent with the project goals.

By contributing to this project, you agree that your contributions may be released under the repository license.
