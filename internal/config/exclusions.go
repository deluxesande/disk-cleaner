package config

import (
	"path/filepath"
	"runtime"
	"strings"
)

// GetDefaultExclusions returns OS-specific system directories that should be skipped
// to prevent permission errors and avoid scanning critical OS files.
func GetDefaultExclusions() []string {
	if runtime.GOOS == "windows" {
		return []string{
			filepath.Join("C:", "Windows"),
			filepath.Join("C:", "Program Files"),
			filepath.Join("C:", "Program Files (x86)"),
			filepath.Join("C:", "$Recycle.Bin"),
			filepath.Join("C:", "System Volume Information"),
			filepath.Join("C:", "Recovery"),
			filepath.Join("C:", "PerfLogs"),
			filepath.Join("C:", "ProgramData"),
		}
	}

	// Unix-like system exclusions (Linux/macOS)
	return []string{
		"/System",
		"/bin",
		"/sbin",
		"/usr/bin",
		"/usr/sbin",
		"/dev",
		"/proc",
		"/sys",
		"/var/run",
		"/var/lock",
	}
}

// IsExcluded checks if a given path falls within an excluded directory structure.
func IsExcluded(targetPath string) bool {
	target := strings.ToLower(filepath.Clean(targetPath))
	exclusions := GetDefaultExclusions()

	for _, exclusion := range exclusions {
		excl := strings.ToLower(filepath.Clean(exclusion))

		// Check if the target is exactly the exclusion, or a subdirectory of it
		if target == excl || strings.HasPrefix(target, excl+string(filepath.Separator)) {
			return true
		}
	}

	return false
}
