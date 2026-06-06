# Disk Cleaner

Disk Cleaner is a planned high-performance, cross-platform command-line tool for finding files and folders that commonly waste disk space. It is designed for developers and power users who want a safer way to review large caches, build artifacts, temporary files, and duplicate files before removing them.

The project goal is simple: scan a target directory, group cleanup candidates into understandable categories, show how much space can be recovered, and let the user decide what should be deleted.

## What It Will Clean

Disk Cleaner is intended to identify common sources of wasted storage, including:

- Development dependency folders such as `node_modules`
- Build output directories such as `dist`, `build`, `target`, and similar generated folders
- Package manager and tool caches
- Temporary files created by applications or development tools
- Duplicate files that are byte-for-byte identical

The tool is designed to focus on generated or recoverable files instead of personal documents. Deletion should always happen only after review.

## Key Features

### Categorized Scanning

Disk Cleaner groups results by type so they are easier to understand before deletion. Instead of showing one flat list of paths, it separates cleanup candidates such as dependency folders, build output, cache directories, temporary files, and duplicate groups.

### Duplicate Detection

The planned duplicate finder uses a staged approach to avoid unnecessary work:

1. Group files by size, because files with different sizes cannot be exact duplicates.
2. Compare a small initial chunk of files with matching sizes to filter out obvious differences.
3. Compute a full cryptographic hash, such as SHA-256, only for files that still look identical.

This approach keeps duplicate detection accurate while reducing expensive full-file reads.

### Interactive Review

Disk Cleaner is intended to provide an interactive terminal interface where users can:

- Browse cleanup candidates in a tree-style view
- Expand and collapse categories
- Select individual files, folders, or duplicate groups
- Review estimated space savings
- Confirm before deleting anything

### Safety-First Defaults

The application should avoid dangerous system paths by default and use explicit confirmation before destructive actions. Expected safeguards include:

- Skipping protected operating system folders
- Ignoring paths that are likely to cause permission or stability issues
- Showing clear deletion summaries before anything is removed
- Supporting configurable exclusions for user-defined safe zones

### Cross-Platform Design

The project is written with cross-platform behavior in mind. The target platforms are:

- Windows
- macOS
- Linux

Platform-specific paths and exclusions should be handled carefully so scans behave predictably on each operating system.

## Installation

This repository currently contains the project documentation. When the Go source code is added, the expected installation flow will be:

```bash
git clone https://github.com/yourusername/disk-cleaner.git
cd disk-cleaner
go mod download
go build -o disk-cleaner main.go
```

To install the binary into your Go toolchain path:

```bash
go install .
```

## Usage

After the CLI implementation is available, expected usage will look like this:

```bash
# Scan the current directory with default settings
disk-cleaner

# Scan a specific directory
disk-cleaner --dir "D:\Projects"

# Increase scan worker count
disk-cleaner --dir "D:\Projects" --workers 12

# Run duplicate detection only
disk-cleaner --mode dedupe --dir "C:\Users\Public\Documents"
```

## Terminal Controls

The planned interactive terminal UI should support controls similar to:

| Key | Action |
| --- | --- |
| `Up` / `Down` or `k` / `j` | Move through the results list |
| `Right` / `Left` or `l` / `h` | Expand or collapse a category |
| `Space` | Mark or unmark an item for deletion |
| `Enter` | Confirm deletion for selected items |
| `q` or `Ctrl+C` | Quit without deleting anything |

## Intended Project Structure

The expected Go project structure is:

```text
disk-cleaner/
├── main.go
├── go.mod
├── internal/
│   ├── config/     # Defaults, platform rules, and excluded paths
│   ├── models/     # Shared data structures
│   ├── scanner/    # Directory traversal and cleanup candidate detection
│   ├── dedupe/     # Duplicate file detection and hashing pipeline
│   └── ui/         # Interactive terminal UI
├── README.md
└── LICENSE
```

## Build Targets

Once implementation is added, release builds can be created with:

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/disk-cleaner-windows-amd64.exe main.go

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/disk-cleaner-darwin-arm64 main.go

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/disk-cleaner-linux-amd64 main.go
```

## Development Notes

Important implementation goals for the project:

- Keep deletion behavior explicit and reviewable.
- Prefer clear categories and summaries over noisy raw file listings.
- Use concurrency carefully so scans are fast without overwhelming slower disks.
- Treat duplicate detection as exact matching, not fuzzy matching.
- Keep platform-specific path handling isolated and testable.

## License

Disk Cleaner is licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.
