package scanner

import (
	"os"
	"path/filepath"
	"time"

	"github.com/deluxesande/disk-cleaner/internal/models"
)

// RunTempSweep scans the OS default temporary directory and groups items
// by top-level files or folders to prevent UI clutter.
func RunTempSweep() []models.SpaceWaster {
	var tempItems []models.SpaceWaster
	tempDir := os.TempDir()

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return tempItems
	}

	for _, entry := range entries {
		path := filepath.Join(tempDir, entry.Name())
		var size int64
		var lastAccess time.Time

		if entry.IsDir() {
			// Reuse the calculateDirMetrics function from walk.go
			size, lastAccess = calculateDirMetrics(path)
		} else {
			info, err := entry.Info()
			if err == nil {
				size = info.Size()
				lastAccess = info.ModTime()
			}
		}

		// Only append items that are actually consuming disk space
		if size > 0 {
			tempItems = append(tempItems, models.SpaceWaster{
				Path:       path,
				Size:       size,
				LastAccess: lastAccess,
				Category:   "SystemTemp",
			})
		}
	}

	return tempItems
}
