package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetCustomExclusionsFile returns the path to the user's custom ignore file.
func GetCustomExclusionsFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".disk-cleaner-ignore"
	}
	return filepath.Join(home, ".disk-cleaner-ignore")
}

// GetSystemExclusions returns hardcoded, OS-specific system directories to protect.
func GetSystemExclusions() []string {
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
		"/System", "/bin", "/sbin", "/usr/bin", "/usr/sbin",
		"/dev", "/proc", "/sys", "/var/run", "/var/lock",
	}
}

// GetCustomExclusions reads the user's custom paths from the ignore file.
func GetCustomExclusions() []string {
	var exclusions []string
	customFile := GetCustomExclusionsFile()
	file, err := os.Open(customFile)
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				exclusions = append(exclusions, line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Warning: error reading custom exclusions file: %v\n", err)
		}
	}
	return exclusions
}

// GetAllExclusions combines system protections and custom user rules.
func GetAllExclusions() []string {
	return append(GetSystemExclusions(), GetCustomExclusions()...)
}

// IsExcluded checks if a given path falls within an excluded directory structure.
func IsExcluded(targetPath string) bool {
	target := strings.ToLower(filepath.Clean(targetPath))
	exclusions := GetAllExclusions()

	for _, exclusion := range exclusions {
		excl := strings.ToLower(filepath.Clean(exclusion))
		if target == excl || strings.HasPrefix(target, excl+string(filepath.Separator)) {
			return true
		}
	}
	return false
}

// AddCustomExclusion appends a new path directly to the user's ignore file.
func AddCustomExclusion(newPath string) error {
	customFile := GetCustomExclusionsFile()
	f, err := os.OpenFile(customFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(newPath + "\n"); err != nil {
		return err
	}
	return nil
}

// RemoveCustomExclusion deletes a specific path from the ignore file by rewriting it.
func RemoveCustomExclusion(targetPath string) error {
	rules := GetCustomExclusions()
	customFile := GetCustomExclusionsFile()

	// Open file with O_TRUNC to wipe it clean before rewriting
	f, err := os.OpenFile(customFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, rule := range rules {
		// Write back all rules EXCEPT the one we are deleting
		if rule != targetPath {
			if _, err := f.WriteString(rule + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}
