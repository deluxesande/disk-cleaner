package scanner

import (
	"path/filepath"
	"strings"
)

// IdentifyCategory determines if a given directory matches known space-wasting patterns.
// It returns the target slice name ("DevArtifacts", "AppCaches", "TempFiles") or an empty string.
func IdentifyCategory(path string, isDir bool) string {
	if !isDir {
		return ""
	}

	base := strings.ToLower(filepath.Base(path))

	devDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".next":        true,
		"dist":         true,
		"build":        true,
		"target":       true,
		".gradle":      true,
	}

	if devDirs[base] {
		return "DevArtifacts"
	}

	cacheDirs := map[string]bool{
		".cache":    true,
		"npm-cache": true,
		"cache":     true,
	}

	if cacheDirs[base] {
		return "AppCaches"
	}

	tempDirs := map[string]bool{
		"temp": true,
		"tmp":  true,
	}

	if tempDirs[base] {
		return "TempFiles"
	}

	return ""
}
